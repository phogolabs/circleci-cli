package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/phogolabs/circleci-cli/circleci"
	"github.com/phogolabs/cli"
)

// Artifact provides a subcommands to project's build
type Artifact struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *Artifact) CreateCommand() *cli.Command {
	var (
		list     = &ListArtifact{}
		download = &DownloadArtifact{}
	)

	commands := []*cli.Command{
		list.CreateCommand(),
		download.CreateCommand(),
	}

	return &cli.Command{
		Name:        "artifact",
		Usage:       "Contains a subset of commands to work with CircleCI job's artifacts",
		Description: "Contains a subset of commands to work with CircleCI job's artifacts",
		Commands:    commands,
	}
}

// ListArtifact provides a subcommands to project's build
type ListArtifact struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *ListArtifact) CreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List all artifacts for given job",
		Description: "List all artifacts for given job",
		Action:      m.list,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Usage:    "The username or organization name that owns the project",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "project",
				Usage:    "The project name",
				Required: true,
			},
			&cli.IntFlag{
				Name:  "build-number",
				Usage: "A build number",
			},
		},
	}
}

func (m *ListArtifact) list(ctx *cli.Context) error {
	client := ctx.Metadata["client"].(*circleci.Client)

	query := &circleci.ListArtifactInput{
		Username: ctx.String("username"),
		Project:  ctx.String("project"),
	}

	if number := ctx.Int("build-number"); number == 0 {
		query.Build = "latest"
	} else {
		query.Build = fmt.Sprintf("%v", number)
	}

	artifacts, err := client.ListArtifacts(query)
	if err != nil {
		return err
	}

	printer := NewPrinter(ctx.Writer)
	printer.Print(artifacts)
	return nil
}

// Download provides a subcommands to project's build
type DownloadArtifact struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *DownloadArtifact) CreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "download",
		Usage:       "Download all artifacts for given job",
		Description: "Download all artifacts for given job",
		Action:      m.download,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Usage:    "The username or organization name that owns the project",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "project",
				Usage:    "The project name",
				Required: true,
			},
			&cli.IntFlag{
				Name:  "build-number",
				Usage: "A build number",
			},
			&cli.StringFlag{
				Name:     "directory",
				Usage:    "Directory where to store the artifacts",
				Value:    ".",
				Required: true,
			},
		},
	}
}

func (m *DownloadArtifact) download(ctx *cli.Context) error {
	client := ctx.Metadata["client"].(*circleci.Client)

	query := &circleci.ListArtifactInput{
		Username: ctx.String("username"),
		Project:  ctx.String("project"),
	}

	if number := ctx.Int("build-number"); number == 0 {
		query.Build = "latest"
	} else {
		query.Build = fmt.Sprintf("%v", number)
	}

	artifacts, err := client.ListArtifacts(query)
	if err != nil {
		return err
	}

	directory, err := filepath.Abs(ctx.String("directory"))
	if err != nil {
		return err
	}

	for _, artifact := range artifacts {
		fmt.Fprintln(ctx.Writer, "Downloading artifact", filepath.Join(directory, artifact.Path))
		if err := client.DownloadArtifact(artifact, directory); err != nil {
			return err
		}
	}

	return nil
}

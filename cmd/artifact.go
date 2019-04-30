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
	return &cli.Command{
		Name:        "artifact",
		Usage:       "Contains a subset of commands to work with CircleCI job's artifacts",
		Description: "Contains a subset of commands to work with CircleCI job's artifacts",
		Commands: []*cli.Command{
			&cli.Command{
				Name:        "list",
				Usage:       "List all artifacts for given job",
				Description: "List all artifacts for given job",
				Action:      m.list,
				Flags:       m.flags(),
			},
			&cli.Command{
				Name:        "download",
				Usage:       "Download all artifacts for given job",
				Description: "Download all artifacts for given job",
				Flags: append(m.flags(),
					&cli.StringFlag{
						Name:     "directory",
						Usage:    "Directory where to store the artifacts",
						EnvVar:   "CIRCLE_ARTIFACT_DIR",
						Value:    ".",
						Required: true,
					}),
				Action: m.download,
			},
		},
	}
}

func (m *Artifact) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "username",
			Usage:    "The username or organization name that owns the project",
			EnvVar:   "CIRCLE_USERNAME",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "project",
			Usage:    "The project name",
			EnvVar:   "CIRCLE_PROJECT",
			Required: true,
		},
		&cli.StringFlag{
			Name:   "branch",
			Usage:  "A branch name for this project",
			EnvVar: "CIRCLE_BRANCH",
			Value:  "master",
		},
		&cli.StringFlag{
			Name:     "job",
			EnvVar:   "CIRCLE_JOB",
			Usage:    "A job name",
			Required: true,
		},
	}
}

func (m *Artifact) list(ctx *cli.Context) error {
	client := ctx.Metadata["client"].(*circleci.Client)

	build, err := m.find(ctx)
	if err != nil {
		return err
	}

	query := &circleci.ListArtifactInput{
		Username: ctx.String("username"),
		Project:  ctx.String("project"),
		Build:    build,
	}

	artifacts, err := client.ListArtifacts(query)
	if err != nil {
		return err
	}

	printer := NewPrinter(ctx.Writer)
	printer.Print(artifacts)
	return nil
}

func (m *Artifact) download(ctx *cli.Context) error {
	client := ctx.Metadata["client"].(*circleci.Client)

	build, err := m.find(ctx)
	if err != nil {
		return err
	}

	query := &circleci.ListArtifactInput{
		Username: ctx.String("username"),
		Project:  ctx.String("project"),
		Build:    build,
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
		fmt.Fprintln(ctx.Writer, "Downloading artifact at", filepath.Join(directory, artifact.Path))
		if err := client.DownloadArtifact(artifact, directory); err != nil {
			return err
		}
	}

	return nil
}

func (m *Artifact) find(ctx *cli.Context) (int, error) {
	fmt.Fprintf(ctx.Writer, "Searching for builds %s/%s/%s\n",
		ctx.String("username"),
		ctx.String("project"),
		ctx.String("job"))

	client := ctx.Metadata["client"].(*circleci.Client)

	query := &circleci.SearchBuildInput{
		Username: ctx.String("username"),
		Project:  ctx.String("project"),
		Branch:   ctx.String("branch"),
		Job:      ctx.String("job"),
		Offset:   0,
		Limit:    200,
	}

	builds, err := client.SearchBuilds(query)
	if err != nil {
		return -1, err
	}

	for _, build := range builds {
		if build.HasArtifacts {
			return build.BuildNum, nil
		}
	}

	return -1, nil
}

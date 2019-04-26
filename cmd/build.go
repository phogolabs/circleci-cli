package cmd

import (
	"io"

	"github.com/kataras/tablewriter"
	"github.com/landoop/tableprinter"
	"github.com/phogolabs/circleci-cli/circleci"
	"github.com/phogolabs/cli"
)

// Build provides a subcommands to project's build
type Build struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *Build) CreateCommand() *cli.Command {
	var (
		list   = &ListBuild{}
		search = &SearchBuild{}
	)

	commands := []*cli.Command{
		list.CreateCommand(),
		search.CreateCommand(),
	}

	return &cli.Command{
		Name:        "build",
		Usage:       "Contains a subset of commands to work with CircleCI Builds",
		Description: "Contains a subset of commands to work with CircleCI Builds",
		Commands:    commands,
	}
}

// ListBuild provides a subcommands to project's build
type ListBuild struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *ListBuild) CreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List all recent builds",
		Description: "List all recent builds",
		Action:      m.list,
	}
}

func (m *ListBuild) list(ctx *cli.Context) error {
	client := ctx.Metadata["client"].(*circleci.Client)

	builds, err := client.ListRecentBuilds()
	if err != nil {
		return err
	}

	printer := NewPrinter(ctx.Writer)
	printer.Print(builds)
	return nil
}

// SearchBuild provides a subcommands to project's build
type SearchBuild struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *SearchBuild) CreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "search",
		Usage:       "Search for a recent jobs",
		Description: "Search for a recent jobs",
		Action:      m.search,
		Flags: []cli.Flag{
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
				Name:   "job",
				Usage:  "A job name",
				EnvVar: "CIRCLE_JOB",
				Value:  "master",
			},
			&cli.StringFlag{
				Name:   "status",
				EnvVar: "CIRCLE_STATUS",
				Usage:  "Restricts which builds are returned",
			},
			&cli.IntFlag{
				Name:   "offset",
				EnvVar: "CIRCLE_OFFSET",
				Usage:  "The API returns builds starting from this offset",
			},
			&cli.IntFlag{
				Name:   "limit",
				Usage:  "The number of builds to return. Maximum 100",
				EnvVar: "CIRCLE_LIMIT",
				Value:  200,
			},
		},
	}
}

func (m *SearchBuild) search(ctx *cli.Context) error {
	client := ctx.Metadata["client"].(*circleci.Client)

	query := &circleci.SearchBuildInput{
		Username: ctx.String("username"),
		Project:  ctx.String("project"),
		Status:   ctx.String("status"),
		Branch:   ctx.String("branch"),
		Job:      ctx.String("job"),
		Offset:   ctx.Int("offset"),
		Limit:    ctx.Int("limit"),
	}

	builds, err := client.SearchBuilds(query)
	if err != nil {
		return err
	}

	printer := NewPrinter(ctx.Writer)
	printer.Print(builds)
	return nil
}

// NewPrinter creates a new printer
func NewPrinter(writer io.Writer) *tableprinter.Printer {
	printer := tableprinter.New(writer)
	printer.HeaderBgColor = tablewriter.BgBlackColor
	printer.HeaderFgColor = tablewriter.FgGreenColor
	return printer
}

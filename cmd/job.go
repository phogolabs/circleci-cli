package cmd

import (
	"io"

	"github.com/kataras/tablewriter"
	"github.com/landoop/tableprinter"
	"github.com/phogolabs/circleci-cli/circleci"
	"github.com/phogolabs/cli"
)

// Job provides a subcommands to project's build
type Job struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *Job) CreateCommand() *cli.Command {
	var (
		list   = &ListJob{}
		search = &SearchJob{}
	)

	commands := []*cli.Command{
		list.CreateCommand(),
		search.CreateCommand(),
	}

	return &cli.Command{
		Name:        "job",
		Usage:       "Contains a subset of commands to work with CircleCI Jobs",
		Description: "Contains a subset of commands to work with CircleCI Jobs",
		Commands:    commands,
	}
}

// ListJob provides a subcommands to project's build
type ListJob struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *ListJob) CreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List all recent jobs",
		Description: "List all recent jobs",
		Action:      m.list,
	}
}

func (m *ListJob) list(ctx *cli.Context) error {
	client := ctx.Metadata["client"].(*circleci.Client)

	builds, err := client.ListRecentBuilds()
	if err != nil {
		return err
	}

	printer := NewPrinter(ctx.Writer)
	printer.Print(builds)
	return nil
}

// SearchJob provides a subcommands to project's build
type SearchJob struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *SearchJob) CreateCommand() *cli.Command {
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
				Value:  30,
			},
		},
	}
}

func (m *SearchJob) search(ctx *cli.Context) error {
	client := ctx.Metadata["client"].(*circleci.Client)

	query := &circleci.SearchBuildInput{
		Username: ctx.String("username"),
		Project:  ctx.String("project"),
		Status:   ctx.String("status"),
		Branch:   ctx.String("branch"),
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

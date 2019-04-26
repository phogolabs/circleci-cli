package main

import (
	"os"

	"github.com/phogolabs/circleci-cli/circleci"
	"github.com/phogolabs/circleci-cli/cmd"
	"github.com/phogolabs/cli"
)

func main() {
	var (
		job      = &cmd.Job{}
		artifact = &cmd.Artifact{}
	)

	commands := []*cli.Command{
		job.CreateCommand(),
		artifact.CreateCommand(),
	}

	app := &cli.App{
		Name:      "circleci-cli",
		HelpName:  "circleci-cli",
		Usage:     "OpenAPI Viewer and Generator",
		UsageText: "stride [global options]",
		Version:   "1.0-beta-05",
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "token",
				Usage:    "An authentication token",
				EnvVar:   "CIRCLE_TOKEN",
				Required: true,
			},
		},
		AfterInit: afterInit,
		Commands:  commands,
	}

	app.Run(os.Args)
}

func afterInit(ctx *cli.Context) error {
	client := &circleci.Client{
		Token: ctx.String("token"),
	}

	ctx.Metadata["client"] = client
	return nil
}

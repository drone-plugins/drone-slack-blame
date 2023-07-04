package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "slack blame plugin"
	app.Usage = "slack blame plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token",
			Usage:  "slack access token",
			EnvVar: "PLUGIN_TOKEN,SLACK_TOKEN",
		},
		cli.StringFlag{
			Name:   "channel",
			Usage:  "slack channel",
			EnvVar: "PLUGIN_CHANNEL",
		},
		cli.StringFlag{
			Name:   "mapping",
			Usage:  "mapping of authors to slack users",
			EnvVar: "PLUGIN_MAPPING",
		},
		cli.StringFlag{
			Name:   "success_username",
			Usage:  "username for successful builds",
			Value:  "drone",
			EnvVar: "PLUGIN_SUCCESS_USERNAME",
		},
		cli.StringFlag{
			Name:   "success_icon",
			Usage:  "icon for successful builds",
			Value:  ":drone:",
			EnvVar: "PLUGIN_SUCCESS_ICON",
		},
		cli.StringFlag{
			Name:   "success_template",
			Usage:  "template for successful builds",
			EnvVar: "PLUGIN_SUCCESS_TEMPLATE",
		},
		cli.StringSliceFlag{
			Name:   "success_image_attachments",
			Usage:  "image attachments for successful builds",
			EnvVar: "PLUGIN_SUCCESS_IMAGE_ATTACHMENTS",
		},
		cli.StringFlag{
			Name:   "failure_username",
			Usage:  "username for failed builds",
			Value:  "drone",
			EnvVar: "PLUGIN_FAILURE_USERNAME",
		},
		cli.StringFlag{
			Name:   "failure_icon",
			Usage:  "icon for failed builds",
			Value:  ":drone:",
			EnvVar: "PLUGIN_FAILURE_ICON",
		},
		cli.StringFlag{
			Name:   "failure_template",
			Usage:  "template for failed builds",
			EnvVar: "PLUGIN_FAILURE_TEMPLATE",
		},
		cli.StringSliceFlag{
			Name:   "failure_image_attachments",
			Usage:  "image attachments for failed builds",
			EnvVar: "PLUGIN_FAILURE_IMAGE_ATTACHMENTS",
		},
		cli.StringFlag{
			Name:   "repo.fullname",
			Usage:  "repository full name",
			EnvVar: "DRONE_REPO",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "repo.link",
			Usage:  "repository link",
			EnvVar: "DRONE_REPO_LINK",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "git commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "commit.link",
			Usage:  "git commit link",
			EnvVar: "DRONE_COMMIT_LINK",
		},
		cli.StringFlag{
			Name:   "commit.author.name",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.author.email",
			Usage:  "git author email",
			EnvVar: "DRONE_COMMIT_AUTHOR_EMAIL",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.StringFlag{
			Name:   "build.deploy",
			Usage:  "build deployment target",
			EnvVar: "DRONE_DEPLOY_TO",
		},
		cli.IntFlag{
			Name:   "prev.build.number",
			Usage:  "previous build number",
			EnvVar: "DRONE_PREV_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "prev.build.status",
			Usage:  "previous build status",
			EnvVar: "DRONE_PREV_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "prev.commit.sha",
			Usage:  "previous build sha",
			EnvVar: "DRONE_PREV_COMMIT_SHA",
		},
	}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		godotenv.Overload("/run/drone/env")
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Repo: Repo{
			FullName: c.String("repo.fullname"),
			Owner:    c.String("repo.owner"),
			Name:     c.String("repo.name"),
			Link:     c.String("repo.link"),
		},
		Build: Build{
			Commit:    c.String("commit.sha"),
			Branch:    c.String("commit.branch"),
			Ref:       c.String("commit.ref"),
			Link:      c.String("commit.link"),
			Message:   c.String("commit.message"),
			Author:    c.String("commit.author.name"),
			Email:     c.String("commit.author.email"),
			Number:    c.Int("build.number"),
			Status:    c.String("build.status"),
			Event:     c.String("build.event"),
			Deploy:    c.String("build.deploy"),
			BuildLink: c.String("build.link"),
		},
		BuildLast: Build{
			Number: c.Int("prev.build.number"),
			Status: c.String("prev.build.status"),
			Commit: c.String("prev.commit.sha"),
		},
		Config: Config{
			Token:   c.String("token"),
			Channel: c.String("channel"),
			Mapping: c.String("mapping"),
			Success: MessageOptions{
				Username:         c.String("success_username"),
				Icon:             c.String("success_icon"),
				Template:         c.String("success_template"),
				ImageAttachments: c.StringSlice("success_image_attachments"),
			},
			Failure: MessageOptions{
				Username:         c.String("failure_username"),
				Icon:             c.String("failure_icon"),
				Template:         c.String("failure_template"),
				ImageAttachments: c.StringSlice("failure_image_attachments"),
			},
		},
	}

	if plugin.Config.Token == "" {
		return errors.New("Missing authentication token")
	}

	return plugin.Exec()
}

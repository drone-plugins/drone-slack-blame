// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/urfave/cli/v2"

	"github.com/drone-plugins/drone-slack-blame/pkg/slackblame"
)

const (
	tokenFlag   = "token"
	channelFlag = "channel"
	mappingFlag = "mapping"

	successUsernameFlag         = "success.username"
	successIconFlag             = "success.icon"
	successTemplateFlag         = "success.template"
	successImageAttachmentsFlag = "success.image-attachments"

	failureUsernameFlag         = "failure.username"
	failureIconFlag             = "failure.icon"
	failureTemplateFlag         = "failure.template"
	failureImageAttachmentsFlag = "failure.image-attachments"

	defaultUsername = "drone"
	defaultIcon     = ":drone:"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    tokenFlag,
			Usage:   "slack access token",
			EnvVars: []string{"PLUGIN_TOKEN,SLACK_TOKEN"},
		},
		&cli.StringFlag{
			Name:    channelFlag,
			Usage:   "slack channel",
			EnvVars: []string{"PLUGIN_CHANNEL"},
		},
		&cli.StringFlag{
			Name:    mappingFlag,
			Usage:   "mapping of authors to slack users",
			EnvVars: []string{"PLUGIN_MAPPING"},
		},
		&cli.StringFlag{
			Name:    successUsernameFlag,
			Usage:   "username for successful builds",
			Value:   defaultUsername,
			EnvVars: []string{"PLUGIN_SUCCESS_USERNAME"},
		},
		&cli.StringFlag{
			Name:    successIconFlag,
			Usage:   "icon for successful builds",
			Value:   defaultIcon,
			EnvVars: []string{"PLUGIN_SUCCESS_ICON"},
		},
		&cli.StringFlag{
			Name:    successTemplateFlag,
			Usage:   "template for successful builds",
			EnvVars: []string{"PLUGIN_SUCCESS_TEMPLATE"},
		},
		&cli.StringSliceFlag{
			Name:    successImageAttachmentsFlag,
			Usage:   "image attachments for successful builds",
			EnvVars: []string{"PLUGIN_SUCCESS_IMAGE_ATTACHMENTS"},
		},
		&cli.StringFlag{
			Name:    failureUsernameFlag,
			Usage:   "username for failed builds",
			Value:   defaultUsername,
			EnvVars: []string{"PLUGIN_FAILURE_USERNAME"},
		},
		&cli.StringFlag{
			Name:    failureIconFlag,
			Usage:   "icon for failed builds",
			Value:   defaultIcon,
			EnvVars: []string{"PLUGIN_FAILURE_ICON"},
		},
		&cli.StringFlag{
			Name:    failureTemplateFlag,
			Usage:   "template for failed builds",
			EnvVars: []string{"PLUGIN_FAILURE_TEMPLATE"},
		},
		&cli.StringSliceFlag{
			Name:    failureImageAttachmentsFlag,
			Usage:   "image attachments for failed builds",
			EnvVars: []string{"PLUGIN_FAILURE_IMAGE_ATTACHMENTS"},
		},
	}
}

// settingsFromContext creates a plugin.Settings from the cli.Context.
func settingsFromContext(ctx *cli.Context) slackblame.Settings {
	return slackblame.Settings{
		Token:   ctx.String(tokenFlag),
		Channel: ctx.String(channelFlag),
		Mapping: ctx.String(mappingFlag),
		Success: slackblame.MessageOptions{
			Username:         ctx.String(successUsernameFlag),
			Icon:             ctx.String(successIconFlag),
			Template:         ctx.String(successTemplateFlag),
			ImageAttachments: ctx.StringSlice(successImageAttachmentsFlag),
		},
		Failure: slackblame.MessageOptions{
			Username:         ctx.String(failureUsernameFlag),
			Icon:             ctx.String(failureIconFlag),
			Template:         ctx.String(failureTemplateFlag),
			ImageAttachments: ctx.StringSlice(failureImageAttachmentsFlag),
		},
	}
}

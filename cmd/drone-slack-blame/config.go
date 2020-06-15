// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/drone-plugins/drone-slack-blame/plugin"
	"github.com/urfave/cli/v2"
)

const (
	defaultUsername = "drone"
	defaultIcon     = ":drone:"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "token",
			Usage:       "slack access token",
			EnvVars:     []string{"PLUGIN_TOKEN,SLACK_TOKEN"},
			Destination: &settings.Token,
		},
		&cli.StringFlag{
			Name:        "channel",
			Usage:       "slack channel",
			EnvVars:     []string{"PLUGIN_CHANNEL"},
			Destination: &settings.Channel,
		},
		&cli.StringFlag{
			Name:        "mapping",
			Usage:       "mapping of authors to slack users",
			EnvVars:     []string{"PLUGIN_MAPPING"},
			Destination: &settings.Mapping,
		},
		&cli.StringFlag{
			Name:        "success.username",
			Usage:       "username for successful builds",
			Value:       defaultUsername,
			EnvVars:     []string{"PLUGIN_SUCCESS_USERNAME"},
			Destination: &settings.Success.Username,
		},
		&cli.StringFlag{
			Name:        "success.icon",
			Usage:       "icon for successful builds",
			Value:       defaultIcon,
			EnvVars:     []string{"PLUGIN_SUCCESS_ICON"},
			Destination: &settings.Success.Icon,
		},
		&cli.StringFlag{
			Name:        "success.template",
			Usage:       "template for successful builds",
			EnvVars:     []string{"PLUGIN_SUCCESS_TEMPLATE"},
			Destination: &settings.Success.Template,
		},
		&cli.StringSliceFlag{
			Name:        "success.image-attachments",
			Usage:       "image attachments for successful builds",
			EnvVars:     []string{"PLUGIN_SUCCESS_IMAGE_ATTACHMENTS"},
			Destination: &settings.Success.ImageAttachments,
		},
		&cli.StringFlag{
			Name:        "failure.username",
			Usage:       "username for failed builds",
			Value:       defaultUsername,
			EnvVars:     []string{"PLUGIN_FAILURE_USERNAME"},
			Destination: &settings.Failure.Username,
		},
		&cli.StringFlag{
			Name:        "failure.icon",
			Usage:       "icon for failed builds",
			Value:       defaultIcon,
			EnvVars:     []string{"PLUGIN_FAILURE_ICON"},
			Destination: &settings.Failure.Icon,
		},
		&cli.StringFlag{
			Name:        "failure.template",
			Usage:       "template for failed builds",
			EnvVars:     []string{"PLUGIN_FAILURE_TEMPLATE"},
			Destination: &settings.Failure.Template,
		},
		&cli.StringSliceFlag{
			Name:        "failure.image-attachments",
			Usage:       "image attachments for failed builds",
			EnvVars:     []string{"PLUGIN_FAILURE_IMAGE_ATTACHMENTS"},
			Destination: &settings.Failure.ImageAttachments,
		},
	}
}

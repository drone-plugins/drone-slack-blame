// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

// DO NOT MODIFY THIS FILE DIRECTLY

package main

import (
	"os"

	"github.com/drone-plugins/drone-plugin-lib/pkg/urfave"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/drone-plugins/drone-slack-blame/pkg/slackblame"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "slack-blame plugin"
	app.Usage = "targeted build notifications for slack"
	app.Action = run
	app.Flags = append(settingsFlags(), urfave.Flags()...)

	// Run the application
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	urfave.LoggingFromContext(ctx)

	plugin := slackblame.New(
		settingsFromContext(ctx),
		urfave.PipelineFromContext(ctx),
		urfave.NetworkFromContext(ctx),
	)

	// Validate the settings
	if err := plugin.Validate(); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	// Run the plugin
	if err := plugin.Exec(); err != nil {
		return errors.Wrap(err, "exec failed")
	}

	return nil
}

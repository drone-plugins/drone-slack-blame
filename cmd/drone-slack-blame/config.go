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
	// Add all the flag names here as const strings.
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags() []cli.Flag {
	// Replace below with all the flags required for the plugin's specific
	// settings.
	return []cli.Flag{}
}

// settingsFromContext creates a plugin.Settings from the cli.Context.
func settingsFromContext(ctx *cli.Context) slackblame.Settings {
	// Replace below with the parsing of the
	return slackblame.Settings{}
}

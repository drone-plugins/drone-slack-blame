// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package slackblame

import "text/template"

const (
	defaultSuccessTemplate = `

	`

	defaultFailureTemplate = `

	`
)

func parseTemplate(template string) *template.Template {
	return nil
}

// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"encoding/json"
	"strings"

	"github.com/slack-go/slack"
)

// userMapping gets the user mapping file.
func userMapping(str string) (map[string]string, error) {
	values := map[string]string{}
	if str == "" {
		return values, nil
	}

	err := json.Unmarshal([]byte(str), &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

// checkEmail sees if the email is used by the user.
func checkEmail(user *slack.User, email string) bool {
	return strings.EqualFold(user.Profile.Email, email)
}

// checkUsername sees if the username is the same as the user.
func checkUsername(user *slack.User, name string) bool {
	return user.Profile.DisplayName == name || user.RealName == name
}

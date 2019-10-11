// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package slackblame

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/drone-plugins/drone-plugin-lib/pkg/plugin"
	"github.com/drone-plugins/drone-plugin-lib/pkg/urfave"
	"github.com/pkg/errors"
)

type pluginImpl struct {
	settings Settings
	pipeline plugin.Pipeline
	network  urfave.Network
}

// New Plugin from the given Settings, Pipeline, and Network.
func New(settings Settings, pipeline plugin.Pipeline, network urfave.Network) plugin.Plugin {
	return &pluginImpl{
		settings: settings,
		pipeline: pipeline,
		network:  network,
	}
}

func (p *pluginImpl) contents(str string) (string, error) {
	// Check for the empty string
	if str == "" {
		return str, nil
	}

	isFilePath := false

	// See if the string is referencing a URL
	if u, err := url.Parse(str); err == nil {
		switch u.Scheme {
		case "http", "https":
			req, err := http.NewRequestWithContext(p.network.Context, "GET", str, nil)
			if err != nil {
				return "", err
			}

			res, err := p.network.Client.Do(req)
			if err != nil {
				return "", err
			}

			defer res.Body.Close()
			b, _ := ioutil.ReadAll(res.Body)
			return string(b), nil

		case "file":
			// Fall through to file loading
			str = u.Path
			isFilePath = true
		}
	}

	// See if the string is referencing a file
	if _, err := os.Stat(str); err == nil {
		b, err := ioutil.ReadFile(str)
		return string(b), err
	}

	if isFilePath {
		return "", errors.Errorf("could not load file %s", str)
	}

	// Its a regular string
	return str, nil
}

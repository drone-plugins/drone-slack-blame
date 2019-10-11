// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package slackblame

import (
	"context"

	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/drone-plugins/drone-plugin-lib/pkg/plugin"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	// Settings for the Plugin.
	Settings struct {
		Token   string
		Channel string
		Mapping string
		Success MessageOptions
		Failure MessageOptions

		mapping map[string]string
	}

	// MessageOptions contains the slack message.
	MessageOptions struct {
		Icon             string
		Username         string
		Template         string
		ImageAttachments []string

		template *template.Template
	}

	// searchFunc determines how to search for a slack user.
	searchFunc func(*slack.User, string) bool
)

func (p *pluginImpl) Validate() error {
	// REMOVE
	p.network.Client = http.DefaultClient
	p.network.Context = context.Background()
	logrus.SetLevel(logrus.TraceLevel)

	// Check the token
	if p.settings.Token == "" {
		return errors.New("slack token not found")
	}

	// Load mapping if requrested
	m, err := p.contents(p.settings.Mapping)
	if err != nil {
		return errors.Wrapf(err, "mapping could not be loaded from %s", p.settings.Mapping)
	}
	p.settings.mapping, err = userMapping(m)
	logrus.WithField("user-mapping", m).Debug("user mapping contents")
	if err != nil {
		return errors.Wrap(err, "could not load mapping")
	}

	// Load template
	if p.pipeline.Build.Status == "success" {
		// Load success template
		st, err := p.contents(p.settings.Success.Template)
		if err != nil {
			return errors.Wrapf(err, "success template could not be loaded from %s", p.settings.Success.Template)
		}
		if st == "" {
			st = defaultSuccessTemplate
		}
		logrus.WithField("success-template", st).Debug("success template contents")

		tmpl := template.New("success-template")
		p.settings.Success.template, err = tmpl.Parse(st)
		if err != nil {
			return errors.Wrap(err, "could not parse success template")
		}
	} else {
		// Load failure template
		ft, err := p.contents(p.settings.Failure.Template)
		if err != nil {
			return errors.Wrapf(err, "failure template could not be loaded from %s", p.settings.Failure.Template)
		}
		if ft == "" {
			ft = defaultFailureTemplate
		}
		logrus.WithField("failure-template", ft).Debug("failure template contents")

		tmpl := template.New("failure-template")
		p.settings.Failure.template, err = tmpl.Parse(ft)
		if err != nil {
			return errors.Wrap(err, "could not parse failure template")
		}
	}

	return nil
}

func (p *pluginImpl) Exec() error {
	// create the API
	api := slack.New(p.settings.Token)

	// verify the connection
	authResponse, err := api.AuthTestContext(p.network.Context)

	if err != nil {
		return errors.Wrap(err, "failed to test auth")
	}

	logrus.WithFields(logrus.Fields{
		"team": authResponse.Team,
		"user": authResponse.User,
	}).Info("Successfully authenticated with Slack API")

	// get the user
	user, _ := p.findSlackUser(api)

	// get the associated @ string
	messageOptions := p.createMessage()
	var userAt string

	if user != nil {
		userAt = fmt.Sprintf("@%s", user.Name)

		_, _, err := api.PostMessageContext(p.network.Context, userAt, messageOptions)

		if err == nil {
			logrus.WithField("username", user.Name).Info("Notified user")
		} else {
			logrus.WithField("username", user.Name).Error("Could not notify user")
		}
	} else {
		userAt = p.pipeline.Commit.Author
		logrus.WithField("author", userAt).Error("Could not find author")
	}

	if p.settings.Channel != "" {
		if !strings.HasPrefix(p.settings.Channel, "#") {
			p.settings.Channel = "#" + p.settings.Channel
		}
		_, _, err := api.PostMessageContext(p.network.Context, p.settings.Channel, messageOptions)

		if err == nil {
			logrus.WithField("channel", p.settings.Channel).Info("Channel notified")
		} else {
			logrus.WithField("channel", p.settings.Channel).Error("Unable to notify channel")
		}
	}

	return nil
}

// createMessage generates the message to post to Slack.
func (p *pluginImpl) createMessage() slack.MsgOption {
	// This is currently deprecated
	var messageOptions MessageOptions
	var color string
	var messageTitle string

	// Determine if the build was a success
	if p.pipeline.Build.Status == "success" {
		messageOptions = p.settings.Success
		color = "good"
		messageTitle = "Build succeeded"
	} else {
		messageOptions = p.settings.Failure
		color = "danger"
		messageTitle = "Build failed"
	}

	// setup the message
	messageParams := slack.PostMessageParameters{
		Username: messageOptions.Username,
	}

	if strings.HasPrefix(messageOptions.Icon, "http") {
		logrus.WithField("icon", messageOptions.Icon).Debug("Icon is a URL")
		messageParams.IconURL = messageOptions.Icon
	} else {
		logrus.WithField("icon", messageOptions.Icon).Debug("Icon is an emoji")
		messageParams.IconEmoji = messageOptions.Icon
	}

	// Render the template
	messageText := strings.Builder{}
	messageValues := struct {
		Build  plugin.Build
		Repo   plugin.Repo
		Commit plugin.Commit
		Stage  plugin.Stage
		Step   plugin.Step
		SemVer plugin.SemVer
	}{
		p.pipeline.Build,
		p.pipeline.Repo,
		p.pipeline.Commit,
		p.pipeline.Stage,
		p.pipeline.Step,
		p.pipeline.SemVer,
	}

	err := messageOptions.template.Execute(&messageText, messageValues)

	if err != nil {
		logrus.Error("could not render template")
	} else {
		logrus.WithField("rendered", messageText.String()).Debug("rendered template")
	}

	// create the attachment
	attachment := slack.Attachment{
		Color:      color,
		Text:       messageText.String(),
		Title:      messageTitle,
		TitleLink:  p.pipeline.Commit.Link,
		MarkdownIn: []string{"pretext", "text", "fields"},
	}

	// Add image if any are provided
	imageCount := len(messageOptions.ImageAttachments)

	if imageCount > 0 {
		logrus.WithField("count", imageCount).Debug("Choosing from images")
		rand.Seed(time.Now().UTC().UnixNano())
		attachment.ImageURL = messageOptions.ImageAttachments[rand.Intn(imageCount)]
	}

	return slack.MsgOptionCompose(
		slack.MsgOptionPostMessageParameters(messageParams),
		slack.MsgOptionAttachments(attachment),
	)
}

// findSlackUser uses the slack API to find the user who made the commit that
// is being built.
func (p *pluginImpl) findSlackUser(api *slack.Client) (*slack.User, error) {
	// get the mapping
	mapping := p.settings.mapping

	// determine the search function to use
	var search searchFunc
	var find string

	if val, ok := mapping[p.pipeline.Commit.AuthorEmail]; ok {
		logrus.WithField("username", val).Info("Searching for user by name, using build.email as key")
		search = checkUsername
		find = val
	} else if val, ok := mapping[p.pipeline.Commit.Author]; ok {
		logrus.WithField("username", val).Info("Searching for user by name, using build.author as key")
		search = checkUsername
		find = val
	} else {
		logrus.WithField("email", p.pipeline.Commit.AuthorEmail).Info("Searching for user by email")
		search = checkEmail
		find = p.pipeline.Commit.AuthorEmail
	}

	if find == "" {
		return nil, errors.New("No user to search for")
	}

	// search for the user
	users, err := api.GetUsersContext(p.network.Context)

	if err != nil {
		return nil, errors.Wrap(err, "failed to query users")
	}

	var blameUser *slack.User

	for _, user := range users {
		if search(&user, find) {
			logrus.WithFields(logrus.Fields{
				"username": user.Name,
				"realname": user.RealName,
				"email":    user.Profile.Email,
			}).Info("Found user")

			blameUser = &user
			break
		} else {
			logrus.WithFields(logrus.Fields{
				"username": user.Name,
				"email":    user.Profile.Email,
			}).Trace("User")
		}
	}

	return blameUser, nil
}
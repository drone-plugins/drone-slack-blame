package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/drone/drone-template-lib/template"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type (
	// MessageOptions contains the slack message.
	MessageOptions struct {
		Icon             string
		Username         string
		Template         string
		ImageAttachments []string
	}

	// Repo information.
	Repo struct {
		FullName string
		Owner    string
		Name     string
		Link     string
	}

	// Build information.
	Build struct {
		Commit    string
		Branch    string
		Ref       string
		Link      string
		Message   string
		Author    string
		Email     string
		Number    int
		Status    string
		Event     string
		Deploy    string
		BuildLink string
	}

	// Config for the plugin.
	Config struct {
		Token   string
		Channel string
		Mapping string
		Success MessageOptions
		Failure MessageOptions
	}

	// Plugin values.
	Plugin struct {
		Repo      Repo
		Build     Build
		BuildLast Build
		Config    Config
		User      *slack.User
	}

	// searchFunc determines how to search for a slack user.
	searchFunc func(*slack.User, string) bool
)

// Exec executes the plugin.
func (p Plugin) Exec() error {
	// create the API
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 2
	api := slack.New(p.Config.Token, slack.OptionHTTPClient(retryClient.StandardClient()))

	// verify the connection
	authResponse, err := api.AuthTest()

	if err != nil {
		return errors.Wrap(err, "failed to test auth")
	}

	logrus.WithFields(logrus.Fields{
		"team": authResponse.Team,
		"user": authResponse.User,
	}).Info("Successfully authenticated with Slack API")

	// get the user
	p.User, _ = p.findSlackUser(api)

	// get the associated @ string
	messageOptions := p.createMessage()
	var userAt string

	if p.User != nil {
		logrus.WithFields(logrus.Fields{
			"username": p.User.Name,
		}).Info("Found user")

		userAt = fmt.Sprintf("@%s", p.User.Name)

		_, _, err := api.PostMessage(userAt, messageOptions)

		if err == nil {
			logrus.WithFields(logrus.Fields{
				"username": p.User.Name,
			}).Info("Notified user")
		} else {
			logrus.WithFields(logrus.Fields{
				"username": p.User.Name,
			}).Error("Could not notify user")
		}
	} else {
		userAt = p.Build.Author
		logrus.WithFields(logrus.Fields{
			"author": userAt,
		}).Error("Could not find author")
	}

	if p.Config.Channel != "" {
		if !strings.HasPrefix(p.Config.Channel, "#") {
			p.Config.Channel = "#" + p.Config.Channel
		}
		_, _, err := api.PostMessage(p.Config.Channel, messageOptions)

		if err == nil {
			logrus.WithFields(logrus.Fields{
				"channel": p.Config.Channel,
			}).Info("Channel notified")
		} else {
			logrus.WithFields(logrus.Fields{
				"channel": p.Config.Channel,
			}).Error("Unable to notify channel")
		}
	}

	return nil
}

// createMessage generates the message to post to Slack.
func (p Plugin) createMessage() slack.MsgOption {
	// This is currently deprecated
	var messageOptions MessageOptions
	var color string
	var messageTitle string

	// Determine if the build was a success
	if p.Build.Status == "success" {
		messageOptions = p.Config.Success
		color = "good"
		messageTitle = "Build succeeded"
	} else {
		messageOptions = p.Config.Failure
		color = "danger"
		messageTitle = "Build failed"
	}

	// setup the message
	messageParams := slack.PostMessageParameters{
		Username: messageOptions.Username,
	}

	if strings.HasPrefix(messageOptions.Icon, "http") {
		logrus.Info("Icon is a URL")
		messageParams.IconURL = messageOptions.Icon
	} else {
		logrus.Info("Icon is an emoji")
		messageParams.IconEmoji = messageOptions.Icon
	}

	messageText, err := template.Render(messageOptions.Template, &p)

	if err != nil {
		logrus.Error("Could not parse template")
	}

	// create the attachment
	attachment := slack.Attachment{
		Color:     color,
		Text:      messageText,
		Title:     messageTitle,
		TitleLink: p.Build.Link,
	}

	// Add image if any are provided
	imageCount := len(messageOptions.ImageAttachments)

	if imageCount > 0 {
		logrus.WithFields(logrus.Fields{
			"count": imageCount,
		}).Info("Choosing from images")
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
func (p Plugin) findSlackUser(api *slack.Client) (*slack.User, error) {
	// get the mapping
	mapping := userMapping(p.Config.Mapping)

	// determine the search function to use
	var search searchFunc
	var find string

	if val, ok := mapping[p.Build.Email]; ok {
		logrus.WithFields(logrus.Fields{
			"username": val,
		}).Info("Searching for user by name, using build.email as key")
		search = checkUsername
		find = val
	} else if val, ok := mapping[p.Build.Author]; ok {
		logrus.WithFields(logrus.Fields{
			"username": val,
		}).Info("Searching for user by name, using build.author as key")
		search = checkUsername
		find = val
	} else {
		// if using email then we call api.GetUserByEmail directlywhich is more efficient
		logrus.WithFields(logrus.Fields{
			"email": p.Build.Email,
		}).Info("Searching for user by email")

		return api.GetUserByEmail(p.Build.Email)
	}

	if len(find) == 0 {
		return nil, errors.New("No user to search for")
	}

	// search for the user
	users, err := api.GetUsers()

	if err != nil {
		return nil, errors.Wrap(err, "failed to query users")
	}

	var blameUser *slack.User

	for _, user := range users {
		if search(&user, find) {
			logrus.WithFields(logrus.Fields{
				"username": user.Name,
				"email":    user.Profile.Email,
			}).Info("Found user")

			blameUser = &user
			break
		} else {
			logrus.WithFields(logrus.Fields{
				"username": user.Name,
				"email":    user.Profile.Email,
			}).Debug("User")
		}
	}

	return blameUser, nil
}

// userMapping gets the user mapping file.
func userMapping(value string) map[string]string {
	mapping := []byte(contents(value))

	// turn into a map
	values := map[string]string{}
	err := json.Unmarshal(mapping, &values)

	if err != nil {
		if len(mapping) != 0 {
			logrus.WithFields(logrus.Fields{
				"mapping": value,
				"error":   err,
			}).Error("Could not parse mapping")
		}

		values = make(map[string]string)
	}

	return values
}

// contents gets the value referenced either in a local filem, a URL or the
// string value itself.
func contents(s string) string {
	if _, err := os.Stat(s); err == nil {
		o, _ := ioutil.ReadFile(s)
		return os.ExpandEnv(string(o))
	}
	if _, err := url.Parse(s); err == nil {
		resp, err := http.Get(s)
		if err != nil {
			return s
		}
		defer resp.Body.Close()
		o, _ := ioutil.ReadAll(resp.Body)
		return os.ExpandEnv(string(o))
	}
	return os.ExpandEnv(s)
}

// checkUsername sees if the username is the same as the user.
func checkUsername(user *slack.User, name string) bool {
	return user.Profile.DisplayName == name || user.RealName == name
}

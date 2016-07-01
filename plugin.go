package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/drone/drone-go/template"
	"github.com/nlopes/slack"
)

type (
	// Repo information.
	Repo struct {
		Owner string
		Name  string
	}

	// Build information.
	Build struct {
		Event  string
		Number int
		Commit string
		Branch string
		Author string
		Email  string
		Status string
		Link   string
	}

	// Config for the plugin.
	Config struct {
		Token   string
		Channel string
		Mapping string
		Success MessageOptions
		Failure MessageOptions
	}

	// MessageOptions contains the slack message.
	MessageOptions struct {
		Icon             string
		Username         string
		Template         string
		ImageAttachments []string
	}

	// Plugin values
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}

	// templatePayload contains values passed to the template
	templatePayload struct {
		Repo  Repo
		Build Build
		User  *slack.User
	}

	// searchFunc determines how to search for a slack user.
	searchFunc func(*slack.User, string) bool
)

// Exec executes the plugin.
func (p Plugin) Exec() error {
	// create the API
	api := slack.New(p.Config.Token)

	// verify the connection
	authResponse, err := api.AuthTest()

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"team": authResponse.Team,
		"user": authResponse.User,
	}).Info("Successfully authenticated with Slack API")

	// get the user
	blameUser, _ := findSlackUser(api, p)

	// get the associated @ string
	messageParams := createMessage(p, blameUser)
	var userAt string

	if blameUser != nil {
		userAt = fmt.Sprintf("@%s", blameUser.Name)

		_, _, err := api.PostMessage(userAt, "", messageParams)

		if err == nil {
			log.WithFields(log.Fields{
				"username": blameUser.Name,
			}).Info("Notified user")
		} else {
			log.WithFields(log.Fields{
				"username": blameUser.Name,
			}).Error("Could not notify user")
		}
	} else {
		userAt = p.Build.Author
		log.WithFields(log.Fields{
			"author": userAt,
		}).Error("Could not find author")
	}

	if len(p.Config.Channel) != 0 {
		if !strings.HasPrefix(p.Config.Channel, "#") {
			p.Config.Channel = "#" + p.Config.Channel
		}
		_, _, err := api.PostMessage(p.Config.Channel, "", messageParams)

		if err == nil {
			log.WithFields(log.Fields{
				"channel": p.Config.Channel,
			}).Info("Channel notified")
		} else {
			log.WithFields(log.Fields{
				"channel": p.Config.Channel,
			}).Error("Unable to notify channel")
		}
	}

	return nil
}

// findSlackUser uses the slack API to find the user who made the commit that
// is being built.
func findSlackUser(api *slack.Client, p Plugin) (*slack.User, error) {
	// get the mapping
	mapping := userMapping(p.Config.Mapping)

	// determine the search function to use
	var search searchFunc
	var find string

	if val, ok := mapping[p.Build.Email]; ok {
		log.WithFields(log.Fields{
			"username": val,
		}).Info("Searching for user by name")
		search = checkUsername
		find = val
	} else {
		log.WithFields(log.Fields{
			"email": p.Build.Email,
		}).Info("Searching for user by email")
		search = checkEmail
		find = p.Build.Email
	}

	if len(find) == 0 {
		return nil, errors.New("No user to search for")
	}

	// search for the user
	users, err := api.GetUsers()

	if err != nil {
		log.Error("Could not query users")
		return nil, err
	}

	var blameUser *slack.User

	for _, user := range users {
		if search(&user, find) {
			log.WithFields(log.Fields{
				"username": user.Name,
				"email":    user.Profile.Email,
			}).Info("Found user")

			blameUser = &user
			break
		} else {
			log.WithFields(log.Fields{
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
			log.WithFields(log.Fields{
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
		return string(o)
	}
	if _, err := url.Parse(s); err == nil {
		resp, err := http.Get(s)
		if err != nil {
			return s
		}
		defer resp.Body.Close()
		o, _ := ioutil.ReadAll(resp.Body)
		return string(o)
	}
	return s
}

// checkEmail sees if the email is used by the user.
func checkEmail(user *slack.User, email string) bool {
	return user.Profile.Email == email
}

// checkUsername sees if the username is the same as the user.
func checkUsername(user *slack.User, name string) bool {
	return user.Name == name
}

// createMessage generates the message to post to Slack.
func createMessage(p Plugin, user *slack.User) slack.PostMessageParameters {
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
		Username:  messageOptions.Username,
		IconEmoji: messageOptions.Icon,
	}

	// setup the payload
	payload := templatePayload{
		Build: p.Build,
		Repo:  p.Repo,
		User:  user,
	}

	messageText, err := template.Render(messageOptions.Template, &payload)

	if err != nil {
		log.Error("Could not parse template")
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
		log.WithFields(log.Fields{
			"count": imageCount,
		}).Info("Choosing from images")
		rand.Seed(time.Now().UTC().UnixNano())
		attachment.ImageURL = messageOptions.ImageAttachments[rand.Intn(imageCount)]
	}

	messageParams.Attachments = []slack.Attachment{attachment}

	return messageParams
}

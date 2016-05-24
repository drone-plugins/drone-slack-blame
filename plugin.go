package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Event  string
		Number int
		Commit string
		Branch string
		Author string
		Status string
		Link   string
	}

	Config struct {
		Token   string
		Channel string
		Mapping string
		Success MessageOptions
		Failure MessageOptions
	}

	MessageOptions struct {
		Icon             string
		Username         string
		Template         string
		ImageAttachments []string
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}

	searchFunc func(*slack.User, string) bool
)

func (p Plugin) Exec() error {
	// create the API
	api := slack.New(p.Config.Token)

	// verify the connection
	authResponse, err := api.AuthTest()

	if err != nil {
		log.Error("Could not authenticate with Slack API token")
		return err
	} else {
		log.WithFields(log.Fields{
			"team": authResponse.Team,
			"user": authResponse.User,
		}).Info("Successfully authenticated with Slack API")
	}

	// get the user
	blameUser, _ := findSlackUser(api, p)

	// get the associated @ string
	var userAt string

	if blameUser != nil {
		userAt = fmt.Sprintf("@%s", blameUser.Name)
	} else {
		userAt = p.Build.Author
		log.WithFields(log.Fields{
			"author": userAt,
		}).Error("Could not find author")
	}

	return nil
}

func findSlackUser(api *slack.Client, p Plugin) (*slack.User, error) {
	// get the mapping
	mapping := userMapping(p.Config.Mapping)

	// determine the search function to use
	var search searchFunc
	var find string

	if val, ok := mapping[p.Build.Author]; ok {
		log.Info("Searching for user by name")
		search = checkUsername
		find = val
	} else {
		log.Info("Searching for user by email")
		search = checkEmail
		find = p.Build.Author
	}

	if len(find) == 0 {
		log.Error("No user to search for")
		return nil, nil
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

func userMapping(value string) map[string]string {
	mapping := []byte(contents(value))

	// turn into a map
	var values interface{}
	err := json.Unmarshal(mapping, &values)

	if err != nil {
		if len(mapping) != 0 {
			log.WithFields(log.Fields{
				"mapping": mapping,
			}).Error("Could not parse mapping")
		}

		return make(map[string]string)
	} else {
		return values.(map[string]string)
	}
}

func contents(value string) string {
	u, err := url.Parse(value)
	if err == nil {
		switch u.Scheme {
		case "http", "https":
			res, err := http.Get(value)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Could not retrieve contents")
				return ""
			}
			defer res.Body.Close()
			out, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Could not read contents of remote file")
				return ""
			}
			value = string(out)

		case "file":
			out, err := ioutil.ReadFile(u.Path)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Could not read contents of local file")
				return ""
			}
			value = string(out)
		}
	}

	return value
}

func checkEmail(user *slack.User, email string) bool {
	return user.Profile.Email == email
}

func checkUsername(user *slack.User, name string) bool {
	return user.Name == name
}

func color(build Build) string {
	switch build.Status {
	case "success":
		return "good"
	case "failure", "error", "killed":
		return "danger"
	default:
		return "warning"
	}
}

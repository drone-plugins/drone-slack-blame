package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/drone/drone-plugin-go/plugin"
	"github.com/nlopes/slack"
)

type Slack struct {
	Token   string         `json:"token"`
	Channel string         `json:"channel"`
	Success MessageOptions `json:"success"`
	Failure MessageOptions `json:"failure"`
}

type MessageOptions struct {
	Icon             string   `json:"icon"`
	Username         string   `json:"username"`
	Message          string   `json:"message"`
	ImageAttachments []string `json:"image_attachments"`
}

var (
	buildCommit string
)

func main() {
	fmt.Printf("Drone Slack Blame Plugin built from %s\n", buildCommit)

	repo := plugin.Repo{}
	build := plugin.Build{}
	system := plugin.System{}
	vargs := Slack{}

	plugin.Param("build", &build)
	plugin.Param("system", &system)
	plugin.Param("repo", &repo)
	plugin.Param("vargs", &vargs)

	// parse the parameters
	if err := plugin.Parse(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// setup the message
	buildLink := fmt.Sprintf("%s/%s/%d", system.Link, repo.FullName, build.Number)
	var messageOptions MessageOptions
	var color string
	var messageText string
	var channelText string

	// Determine if the build was a success
	if build.Status == "success" {
		messageOptions = vargs.Success
		color = "good"
		messageText = fmt.Sprintf("Build succeeded at %s", buildLink)
		channelText = "Thanks"
	} else {
		messageOptions = vargs.Failure
		color = "danger"
		messageText = fmt.Sprintf("Build failed at %s", buildLink)
		channelText = "Blame"
	}

	// set default values
	if len(messageOptions.Username) == 0 {
		messageOptions.Username = "drone"
	}

	if len(messageOptions.Icon) == 0 {
		messageOptions.Icon = ":drone:"
	}

	if len(messageOptions.ImageAttachments) == 0 {
		messageOptions.ImageAttachments = []string{""}
	}

	// setup the message
	messageParams := slack.PostMessageParameters{
		Username:  messageOptions.Username,
		IconEmoji: messageOptions.Icon,
	}

	imageCount := len(messageOptions.ImageAttachments)
	rand.Seed(time.Now().UTC().UnixNano())

	attachment := slack.Attachment{
		Color:    color,
		Text:     messageText,
		ImageURL: messageOptions.ImageAttachments[rand.Intn(imageCount)],
	}

	messageParams.Attachments = []slack.Attachment{attachment}

	// get the commit author
	commitAuthor := build.Email

	// create the slack api
	api := slack.New(vargs.Token)

	// get the users
	//
	// Slack doesn't let you search by email so just need to get
	// everything and find the user in question
	var blameUser *slack.User

	users, _ := api.GetUsers()

	for _, user := range users {
		if user.Profile.Email == commitAuthor {
			fmt.Printf("%s\n", user.Name)
			fmt.Printf("%s\n", user.Profile.Email)
			blameUser = &user
			break
		}
	}

	// notify the user if possible
	var userAt string

	if blameUser != nil {
		userAt = fmt.Sprintf("@%s", blameUser.Name)

		// send the message to the user's channel
		//
		// this will appear through slackbot
		_, _, err := api.PostMessage(userAt, messageOptions.Message, messageParams)

		if err == nil {
			fmt.Printf("User %s notified\n", userAt)
		} else {
			fmt.Printf("Could not notify user %s!\n", userAt)
		}
	} else {
		userAt = build.Author
		fmt.Print("User could not be found")
	}

	// notify the channel if requested
	if len(vargs.Channel) != 0 {
		if !strings.HasPrefix(vargs.Channel, "#") {
			vargs.Channel = "#" + vargs.Channel
		}

		_, _, err := api.PostMessage(vargs.Channel, fmt.Sprintf("%s %s %s", messageOptions.Message, channelText, userAt), messageParams)

		if err == nil {
			fmt.Printf("Channel notified\n")
		} else {
			fmt.Printf("Could not notify channel!\n")
		}
	}
}

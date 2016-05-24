package main

import(
	log "github.com/Sirupsen/logrus"
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
)

func (p Plugin) Exec() error {
	log.WithFields(log.Fields{
		"username": p.Config.Success.Username,
		"icon": p.Config.Success.Icon,
		"image_attachments": p.Config.Success.ImageAttachments,
	}).Info("Success stuff")
	return nil
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

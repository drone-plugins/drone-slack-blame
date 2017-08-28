# drone-slack-blame

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-slack-blame/status.svg)](http://beta.drone.io/drone-plugins/drone-slack-blame)
[![Join the chat at https://gitter.im/drone/drone](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/drone/drone)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-slack-blame?status.svg)](http://godoc.org/github.com/drone-plugins/drone-slack-blame)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-slack-blame)](https://goreportcard.com/report/github.com/drone-plugins/drone-slack-blame)
[![](https://images.microbadger.com/badges/image/plugins/slack-blame.svg)](https://microbadger.com/images/plugins/slack-blame "Get your own image badge on microbadger.com")

Drone plugin to send build status blames via Slack. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-slack-blame/).

## Build

Build the binary with the following commands:

```
go build
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-slack-blame
docker build --rm -t plugins/slack-blame .
```

## Usage

Execute from the working directory:

```sh
docker run --rm \
  -e PLUGIN_TOKEN=xxxxx \
  -e PLUGIN_CHANNEL=dev \
  -e PLUGIN_SUCCESS_USERNAME="Happy Keanu (on behalf of Drone)" \
  -e PLUGIN_SUCCESS_ICON=":happy_keanu:" \
  -e PLUGIN_SUCCESS_MESSAGE="The build is fixed!" \
  -e PLUGIN_SUCCESS_IMAGE_ATTACHMENTS="http://i.imgur.com/TP4PIxc.jpg" \
  -e PLUGIN_FAILURE_USERNAME="Sad Keanu (on behalf of Drone)" \
  -e PLUGIN_FAILURE_ICON=":sad_keanu:" \
  -e PLUGIN_FAILURE_MESSAGE="The build is broken!" \
  -e PLUGIN_FAILURE_IMAGE_ATTACHMENTS="http://cdn.meme.am/instances/51000361.jpg" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/slack-blame
```

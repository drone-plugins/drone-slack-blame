# drone-slack-blame

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-slack-blame/status.svg)](http://beta.drone.io/drone-plugins/drone-slack-blame)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-slack-blame?status.svg)](http://godoc.org/github.com/drone-plugins/drone-slack-blame)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-slack-blame)](https://goreportcard.com/report/github.com/drone-plugins/drone-slack-blame)
[![Join the chat at https://gitter.im/drone/drone](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/drone/drone)

Drone plugin to send build status blames via Slack. For the usage information
and a listing of the available options please take a look at
[the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t plugins/slack-blame .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-slack-blame' not found or does not exist..
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
  -v $(pwd)/$(pwd) \
  -w $(pwd) \
  plugins/slack-blame
```

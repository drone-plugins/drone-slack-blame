# drone-slack-blame
Drone plugin for sending Slack notifications

## Overview

This plugin is responsible for sending build notifications to developers when
they break the build, either directly through slackbot messages, or to the
designated channel. This lets developers get feedback on their own feature
branches without bombarding the main development channel as well as notifying
the channel on important branches such as release or master.

## Usage

```sh
./drone-slack-blame <<EOF
{
    "repo" : {
        "owner": "foo",
        "name": "bar",
        "full_name": "foo/bar"
    },
    "system": {
        "link_url": "http://mydrone.io"
    },
    "build" : {
        "number": 22,
        "status": "failed",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "commit": "9f2849d5",
        "branch": "master",
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
    },
    "vargs": {
        "token": "xxxxxx",
        "channel": "dev",
        "success": {
            "username": "Happy Keanu (on behalf of Drone)",
            "icon": ":happy_keanu:",
            "message": "The build is fixed!",
            "image_attachments": [
                "http://i.imgur.com/TP4PIxc.jpg"
            ]
        },
        "failure": {
            "username": "Sad Keanu (on behalf of Drone)",
            "icon": ":sad_keanu:",
            "message": "The build is broken!",
            "image_attachments": [
                "http://cdn.meme.am/instances/51000361.jpg"
            ]
        }
    }
}
EOF

```

## Docker

Build the Docker container. Note that we need to use the `-netgo` tag so that
the binary is built without a CGO dependency:

```sh
CGO_ENABLED=0 go build -a -tags netgo
docker build --rm=true -t plugins/drone-slack .
```

Send a Slack notification:

```sh
docker run -i plugins/drone-slack-blame <<EOF
{
    "repo" : {
        "owner": "foo",
        "name": "bar",
        "full_name": "foo/bar"
    },
    "system": {
        "link_url": "http://mydrone.io"
    },
    "build" : {
        "number": 22,
        "status": "failed",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "commit": "9f2849d5",
        "branch": "master",
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
    },
    "vargs": {
        "token": "xxxxxx",
        "channel": "dev",
        "success": {
            "username": "Happy Keanu (on behalf of Drone)",
            "icon": ":happy_keanu:",
            "message": "The build is fixed!",
            "image_attachments": [
                "http://i.imgur.com/TP4PIxc.jpg"
            ]
        },
        "failure": {
            "username": "Sad Keanu (on behalf of Drone)",
            "icon": ":sad_keanu:",
            "message": "The build is broken!",
            "image_attachments": [
                "http://cdn.meme.am/instances/51000361.jpg"
            ]
        }
    }
}
EOF
```

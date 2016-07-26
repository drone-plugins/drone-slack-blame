# drone-slack-blame

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-slack-blame/status.svg)](http://beta.drone.io/drone-plugins/drone-slack-blame)
[![Coverage Status](https://aircover.co/badges/drone-plugins/drone-slack-blame/coverage.svg)](https://aircover.co/drone-plugins/drone-slack-blame)
[![](https://badge.imagelayers.io/plugins/drone-slack-blame:latest.svg)](https://imagelayers.io/?images=plugins/drone-slack-blame:latest 'Get your own badge on imagelayers.io')

Drone plugin to send build status blames via Slack. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Binary

Build the binary using `make`:

```
make deps build
```

### Example

```sh
./drone-slack-blame <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "owner": "drone",
        "name": "drone",
        "full_name": "drone/drone"
    },
    "system": {
        "link_url": "https://beta.drone.io"
    },
    "build": {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com",
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
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

Build the container using `make`:

```
make deps docker
```

### Example

```sh
docker run -i plugins/drone-slack-blame <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "owner": "drone",
        "name": "drone",
        "full_name": "drone/drone"
    },
    "system": {
        "link_url": "https://beta.drone.io"
    },
    "build": {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com",
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
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

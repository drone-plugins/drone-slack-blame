Use the Slack Blame plugin to send a message to a Slack channel or through
direct message when a build completes. You will need to supply Drone with an
access token to the Slack API. You can create a new access token here:
https://api.slack.com/web

## Overview

This plugin is responsible for sending build notifications to developers when
they break the build, either directly through slackbot messages, or to the
designated channel. This lets developers get feedback on their own feature
branches without bombarding the main development channel as well as notifying
the channel on important branches such as release or master.

The following parameters are used to configure the notification

* **token** - the access token.
* **channel** - the channel to post to (if present messages will also be posted
  to the channel).
* **mapping** - a JSON string containing a mapping between the commit author's
  email address and the Slack username. This can be a URL, a local file or a
  string. If present this will be used in place of querying the Slack API.
* **success_icon** - an emoji or image url for the bot announcing a successful
  build (defaults to ":drone:")
* **success_username** - the username for a successful build (defaults to
  "drone").
* **success_template** - the message template for a successful build. This can
  be a URL, a local file or a string.
* **success_image_attachments** - an optional list of image attachments to
  append (only one will be randomly selected per build) to the message when a
  build is successful.
* **failure_icon** - an emoji or image url for the bot announcing a failed
  build (defaults to ":drone:")
* **failure_username** - the username for a failed build (defaults to "drone").
* **failure_template** - the message template for a failed build. This can be
  a URL, a local file or a string.
* **failure_image_attachments** - an optional list of image attachments to
  append (only one will be randomly selected per build) to the message when a
  build fails.

The following secret values can be set to configure the plugin.

  * **SLACK_TOKEN** - corresponds to **token**.

It is highly recommended to put the **SLACK_TOKEN** into a secret so it is not
exposed to users. This can be done using the drone-cli.

```bash
drone secret add --image=slack-blame \
    octocat/hello-world SLACK_TOKEN xxxxxxx
```

Then sign the YAML file after all secrets are added.

```bash
drone sign octocat/hello-world
```

See [secrets](http://readme.drone.io/0.5/usage/secrets/) for additional
information on secrets

## Templates

The Slack Blame plugin uses [mustache templates](https://mustache.github.io/).

The following values are available.
  * **repo.owner** - The repository owner..
  * **repo.name** - The name of the repository.
  * **build.event** - The event triggering the build.
  * **build.number** - The associated build number.
  * **build.commit** - The commit hash.
  * **build.branch** - The branch the build was created from.
  * **build.author** - The author of the commit.
  * **build.email** - The email for the author of the commit.
  * **build.status** - The status of the build.
  * **build.link** - The link to the build.
  * **user.name** - The slack username of the author.

## Examples

The following is a simple Slack Blame configuration in your .drone.yml file:

```yaml
pipeline:
  slack_blame:
    image: plugins/slack-blame
    channel: dev
    success_template: |
      The build is fixed! Thanks @{{slack.name}}
    success_image_attachments:
      - "http://i.imgur.com/TP4PIxc.jpg"
    failure_template: |
      The build is broken! Blame {{slack.name}}
    failure_image_attachments:
      - "http://cdn.meme.am/instances/51000361.jpg"
```

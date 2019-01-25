# escape=`
FROM plugins/base:windows-1803

LABEL maintainer="Drone.IO Community <drone-dev@googlegroups.com>" `
  org.label-schema.name="Drone Slack Blame" `
  org.label-schema.vendor="Drone.IO Community" `
  org.label-schema.schema-version="1.0"

ADD release/windows/amd64/drone-slack-blame.exe C:/bin/drone-slack-blame.exe
ENTRYPOINT [ "C:\\bin\\drone-slack-blame.exe" ]

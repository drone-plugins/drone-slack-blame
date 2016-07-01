# Docker image for Drone's Slack blame notification plugin
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-slack-blame .

FROM alpine:3.4
RUN apk add -U ca-certificates && rm -rf /var/cache/apk/*
ADD drone-slack-blame /bin/
ENTRYPOINT ["/bin/drone-slack-blame"]

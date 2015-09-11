# Docker image for Drone's Slack blame notification plugin
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-slack-blame .

FROM gliderlabs/alpine:3.1
RUN apk-install ca-certificates
ADD drone-slack-blame /bin/
ENTRYPOINT ["/bin/drone-slack-blame"]

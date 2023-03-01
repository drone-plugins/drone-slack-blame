FROM golang:1.19-alpine
WORKDIR /go/src/github.com/quintoandar/drone-slack-blame
ADD . .
RUN GOOS=linux CGO_ENABLED=0 go build -a -tags netgo -trimpath -ldflags='-buildid= -w -s -extldflags "-static"' -o /bin/drone-slack-blame

FROM plugins/base:multiarch
COPY --from=0 /bin/drone-slack-blame /bin/drone-slack-blame

LABEL maintainer="QuintoAndar" \
  org.label-schema.name="Drone Slack Blame" \
  org.label-schema.vendor="QuintoAndar" \
  org.label-schema.schema-version="1.0"

ENTRYPOINT [ "/bin/drone-slack-blame" ]

FROM m2msolutions-docker.jfrog.io/m2m-go-alpine:1.0.0


ARG PROJECT
WORKDIR /go/src/$PROJECT
ENV PROJECT $PROJECT
ENV COMMAND "lazypersistence -config-location=/go/bin/config.yml"


ENV GIT_URL ssh://git-codecommit.us-east-1.amazonaws.com/v1/repos/
ENV ENV_RUN DEV
COPY . .


RUN govendor init
RUN govendor fetch +missing
RUN go install -v ./...
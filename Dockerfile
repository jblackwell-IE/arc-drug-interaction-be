# FROM alpine:3.12.1

# COPY arc-drug-interaction-be /arc-drug-interaction-be

# ENTRYPOINT ["/arc-drug-interaction-be"]

## Local development build
FROM golang:1.13-buster AS dev_environment

WORKDIR /app
RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]

COPY . .

# meant to be used with docker-compose and volume mounted ssh keys.
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
ENTRYPOINT CompileDaemon -log-prefix=false -build="go build ." -command "./arc-drug-interaction-be"
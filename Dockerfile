FROM golang:1.13-buster

RUN mkdir /app
WORKDIR /app
ENV GOPRIVATE=github.com/IdeaEvolver
RUN mkdir ~/.ssh && ssh-keyscan -t rsa github.com > ~/.ssh/known_hosts
RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]
RUN ["go", "get", "github.com/IdeaEvolver/cutter-pkg"]
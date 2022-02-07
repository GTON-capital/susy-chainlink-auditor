# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app

ENV GO111MODULE=on

COPY . .
RUN go mod tidy
RUN go mod download


RUN go build -o /initiator

# EXPOSE 8080

CMD [ "/initiator" ]
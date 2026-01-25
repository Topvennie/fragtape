FROM golang:1.25.5-alpine3.22

WORKDIR /app

RUN apk add --no-cache make gcc musl-dev

RUN go install github.com/air-verse/air@latest
COPY .air.server.toml .

COPY go.mod go.sum ./

ENV APP_ENV=development

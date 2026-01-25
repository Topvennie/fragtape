FROM golang:1.25.5-alpine3.22 AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o recorder.exe ./cmd/recorder/main.go

FROM dockurr/windows:latest

COPY --from=builder /app/docker/recorder-oem/ /oem/
COPY --from=builder /app/recorder.exe /oem/recorder.exe
COPY --from=builder /app/config/production.yml /oem/config/production.yml


FROM golang:1.25.5-alpine3.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o worker ./cmd/worker/main.go

FROM alpine:3.22
WORKDIR /app

COPY --from=builder /app/worker .
COPY --from=builder /app/config/production.toml ./config/production.toml

ENTRYPOINT ["./worker"]

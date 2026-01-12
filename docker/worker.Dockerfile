FROM golang:1.25.5-alpine3.22 AS builder
WORKDIR /app

RUN apk add --no-cache gcc musl-dev libwebp-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o worker ./cmd/worker/main.go

FROM alpine:3.22
WORKDIR /app

RUN apk add --no-cache libwebp ca-certificates tzdata

COPY --from=builder /app/worker .
COPY --from=builder /app/config/production.yml ./config/production.yml

ENTRYPOINT ["./worker"]

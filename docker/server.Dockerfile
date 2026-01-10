# Stage 1: Build Backend
FROM golang:1.25.5-alpine3.22 AS backend-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o migrate ./migrate.go

# Stage 2: Build Frontend
FROM node:24.12.0-alpine3.22 AS frontend-builder
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /frontend/ui

COPY ./ui/package.json ./ui/pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY ./ui .

ARG APP_VERSION

ENV VITE_APP_VERSION=$APP_VERSION
ENV CI=true

RUN pnpm run build

# Stage 3: Combine
FROM alpine:3.22 AS prod
WORKDIR /app

COPY --from=backend-builder /app/server .
COPY --from=backend-builder /app/migrate .
COPY --from=backend-builder /app/config/production.toml ./config/production.toml
COPY --from=frontend-builder /frontend/ui/dist ./public

ENV APP_ENV=production

EXPOSE 8000

ENTRYPOINT ["/bin/sh", "-c" , "./migrate && ./server"]

# syntax=docker/dockerfile:1

FROM node:20-alpine AS web
WORKDIR /build/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.25-bookworm AS backend
WORKDIR /build
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web /build/cmd/cc-go/web-dist ./cmd/cc-go/web-dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /cc-go ./cmd/cc-go/

FROM debian:bookworm-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*
COPY --from=backend /cc-go /usr/local/bin/cc-go
ENV HOME=/root
WORKDIR /root
EXPOSE 18080
CMD ["cc-go"]

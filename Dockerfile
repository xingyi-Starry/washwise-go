# syntax=docker/dockerfile:1

# Create a stage for building the application.
ARG GO_VERSION=1.24.2
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build
WORKDIR /src

RUN apk add --no-cache gcc musl-dev

# Download dependencies
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

# Build
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=1 GOARCH=$TARGETARCH go build -a -ldflags '-extldflags "-static"' -o /bin/server .

# Create a new stage for deployment
FROM alpine:latest AS final

# Install runtime dependencies
RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
    ca-certificates \
    tzdata \
    && \
    update-ca-certificates

# Create a non-privileged user
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

# 设置工作目录
WORKDIR /app

# 创建配置和数据目录
RUN mkdir -p /app/config /app/data && \
    chown -R appuser:appuser /app

USER appuser

# Copy executable from the "build" stage
COPY --from=build /bin/server /app/server

# Copy swagger files directly from build context
COPY --chown=appuser:appuser docs/swagger.json /app/docs/swagger.json
COPY --chown=appuser:appuser docs/swagger.yaml /app/docs/swagger.yaml

EXPOSE 8000

ENTRYPOINT [ "/app/server" ]

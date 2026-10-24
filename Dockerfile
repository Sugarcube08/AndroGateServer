# Build stage
FROM golang:1.26-alpine AS build

ARG APP=sms-gateway
ARG APP_VERSION=1.0.0
ARG APP_RELEASE_ID=1
WORKDIR /go/src

# Cache Go modules
COPY go.* ./
RUN go mod download

# Copy repository source code
COPY . .

# Compile statically linked binary using repository code
RUN CGO_ENABLED=0 go build -a -installsuffix cgo \
    -ldflags="-w -s -X github.com/android-sms-gateway/server/internal/version.AppVersion=${APP_VERSION} -X github.com/android-sms-gateway/server/internal/version.AppRelease=${APP_RELEASE_ID}" \
    -o app ./cmd/${APP}/main.go

# Production runtime stage
FROM alpine:3.21 AS prod

WORKDIR /app

# ca-certificates for outbound TLS; tzdata for timezone calculation
RUN apk add --no-cache ca-certificates tzdata

COPY scripts/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

COPY --from=build /go/src/app /app

EXPOSE 3000

USER guest

ENTRYPOINT ["/docker-entrypoint.sh"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
    CMD [ "/app/app", "health" ]

CMD [ "/app/app" ]

# Building the binary of the App
FROM golang:1.26-alpine AS build

ARG APP=sms-gateway
ARG APP_VERSION=1.0.0
ARG APP_RELEASE_ID=1
WORKDIR /go/src

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go.mod and go.sum
COPY go.* ./

# Downloads all the dependencies in advance
RUN go mod download

# Copy all the Code and stuff to compile everything
COPY . .

# Replace {APP_VERSION} with the actual version
RUN sed -i "s/{APP_VERSION}/${APP_VERSION}/g" ./cmd/${APP}/main.go

# Generate swagger docs
RUN go generate ./...

# Builds the application as a statically linked binary
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-w -s -X github.com/android-sms-gateway/server/internal/version.AppVersion=${APP_VERSION} -X github.com/android-sms-gateway/server/internal/version.AppRelease=${APP_RELEASE_ID}" -o app ./cmd/${APP}/main.go

# Moving the binary to the 'final Image'
FROM alpine:3 AS prod

WORKDIR /app

RUN apk add --no-cache tzdata

COPY scripts/docker-entrypoint.sh /docker-entrypoint.sh

COPY --from=build /go/src/app /app

# Exposes port 3000
EXPOSE 3000

USER guest

ENTRYPOINT ["/docker-entrypoint.sh"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
    CMD [ "/app/app", "health" ]

CMD [ "/app/app" ]

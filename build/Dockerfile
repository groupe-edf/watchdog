FROM node:alpine AS build
WORKDIR /usr/src/app/



FROM golang:alpine AS builder
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
ENV CGO_ENABLED=0
# Create user
ENV USER=watchdog
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
# Install git.
# Git is required for fetching the dependencies.
WORKDIR $GOPATH/src/github.com/groupe-edf/watchdog
COPY . .
# Fetch dependencies.
RUN go mod download
RUN go mod verify
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/watchdog ./cmd/watchdog-server


FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Copy our static executable.
COPY --from=builder /go/bin/watchdog /go/bin/watchdog
# Use an unprivileged user.
USER watchdog:watchdog
# Run the watchdog binary.
ENTRYPOINT ["/go/bin/watchdog"]

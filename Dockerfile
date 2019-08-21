FROM golang:alpine as build-env
MAINTAINER "Martin Ombura Jr. <info@martinomburajr.com>"

# Install git + SSL ca certificates
RUN apk update && apk add --no-cache git ca-certificates tzdata bash && update-ca-certificates

# Create appuser.
#RUN adduser -D -g '' appuser
ENV GO111MODULE=on
WORKDIR /app
# <- COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/gcp-vision-api
#RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/plutoagent
#RUN go build -o /go/bin/gcp-vision-api

FROM alpine@sha256:644fcb1a676b5165371437feaa922943aaf7afcfa8bfee4472f6860aad1ef2a0

RUN apk add --no-cache ca-certificates

RUN mkdir -p /go/bin/app/credentials
COPY --from=build-env /go/bin/gcp-vision-api /go/bin/gcp-vision-api
COPY /app/credentials/credentials.json /go/bin/app/credentials

# Use an unprivileged user.
#USER appuser
WORKDIR /go/bin/

ENV PORT 8080
ENV GOOGLE_APPLICATION_CREDENTIALS /go/bin/app/credentials/credentials.json

ENTRYPOINT ["./gcp-vision-api"]
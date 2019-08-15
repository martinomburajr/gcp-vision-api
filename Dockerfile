FROM golang:alpine as build-env
MAINTAINER "Martin Ombura Jr. <info@martinomburajr.com>"

# Install git + SSL ca certificates
RUN apk update && apk add --no-cache git ca-certificates tzdata bash && update-ca-certificates

# Create appuser.
#RUN adduser -D -g '' appuser

WORKDIR /go/src/github.com/martinomburajr/gcp-vision-api
COPY . .

RUN go get -d -v

# Build the binary
#RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/plutoagent.
#RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/plutoagent
RUN go build -o /go/bin/gcp-vision-api

FROM alpine@sha256:644fcb1a676b5165371437feaa922943aaf7afcfa8bfee4472f6860aad1ef2a0

COPY --from=build-env /go/bin/gcp-vision-api /go/bin/gcp-vision-api
COPY --from=build-env /go/src/github.com/martinomburajr/gcp-vision-api/credentials/ /go/bin/credentials/

# Use an unprivileged user.
#USER appuser

ENV PORT 8080

ENTRYPOINT ["./go/bin/gcp-vision-api"]
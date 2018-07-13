FROM golang:1.10-alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /go/src/github.com/int128/kubesnapshot
COPY . .
RUN go get -v -t -d ./...
RUN go install

FROM alpine:latest
RUN apk update && apk add --no-cache ca-certificates
COPY --from=builder /go/bin/kubesnapshot /kubesnapshot
USER daemon
ENTRYPOINT ["/kubesnapshot"]

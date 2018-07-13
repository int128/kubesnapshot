FROM golang:1.10 AS builder
ADD . /go/src/github.com/int128/kubesnapshot
WORKDIR /go/src/github.com/int128/kubesnapshot
RUN go get -v -t -d ./...
RUN go build

FROM alpine:latest
COPY --from=builder /go/src/github.com/int128/kubesnapshot/kubesnapshot /usr/local/bin/kubesnapshot
USER daemon
CMD ["/usr/local/bin/kubesnapshot"]

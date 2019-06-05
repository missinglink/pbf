FROM golang:1.10-alpine3.8 AS builder

WORKDIR "$GOPATH/src/github.com/tadjik1/pbf"

RUN apk update \
  && apk add git gcc musl-dev

COPY . "$GOPATH/src/github.com/tadjik1/pbf"

RUN go get && go build

FROM alpine:3.8

COPY --from=builder /go/src/github.com/tadjik1/pbf/pbf /bin/
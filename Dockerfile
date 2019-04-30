FROM golang:1.12 AS builder

ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64

ADD . /go/src/github.com/phogolabs/circleci-cli

WORKDIR /go/src/github.com/phogolabs/circleci-cli

RUN go mod download
RUN go build github.com/phogolabs/circleci-cli/cmd/circleci-cli

FROM alpine:3.9

RUN apk add --no-cache ca-certificates=~20190108
RUN update-ca-certificates

RUN mkdir -p /usr/local/bin
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY --from=builder /go/src/github.com/phogolabs/circleci-cli/circleci-cli /usr/local/bin

CMD ["/usr/local/bin/circleci-cli"]

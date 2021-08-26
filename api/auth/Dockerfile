FROM golang:1.17-alpine

RUN mkdir -p /go/src/auth
ADD . /go/src/auth


RUN cd /go/src/auth && go build && cp auth /go

ENTRYPOINT /go/auth

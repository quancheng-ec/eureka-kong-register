FROM golang:1.9

ENV GOROOT=/usr/local/go

RUN go get github.com/tools/godep
RUN go get github.com/zephyrpersonal/eureka-kong-register

WORKDIR /go/src/github.com/zephyrpersonal/eureka-kong-register
RUN godep restore
RUN go build main.go
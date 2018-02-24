FROM golang:1.9

ENV GOROOT=/usr/local/go

ADD . /go/src/github.com/zephyrpersonal/eureka-kong-register

WORKDIR /go/src/github.com/zephyrpersonal/eureka-kong-register

RUN go build main.go

CMD ["./main"]
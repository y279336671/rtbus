FROM golang:1.8.3-onbuild
MAINTAINER Xue Bing <xuebing1110@gmail.com>

RUN cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# move to GOPATH
RUN mkdir -p /go/src/github.com/xuebing1110/rtbus
# RUN cp -r /go/src/app/* $GOPATH/src/github.com/xuebing1110/rtbus
COPY . $GOPATH/src/github.com/xuebing1110/rtbus/
WORKDIR $GOPATH/src/github.com/xuebing1110/rtbus

# example config
RUN mkdir /etc/rtbus
RUN cp server/log.json /etc/rtbus/log.json

# build
RUN mkdir -p /app
RUN go build -o /app/rtbus rtbus/server/main.go

WORKDIR /app
EXPOSE 1318
CMD ["/app/rtbus"]
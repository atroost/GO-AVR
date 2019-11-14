FROM golang:alpine as builder
MAINTAINER Alexander Troost <alexander.troost@kpn.com>

RUN apk update && apk add git && go get gopkg.in/natefinch/lumberjack.v2
RUN mkdir /server
COPY . /server
WORKDIR /server

EXPOSE 2498

RUN go build -o avrserver .

ENTRYPOINT ["./avrserver"]

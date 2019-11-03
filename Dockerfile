FROM golang:alpine as builder
MAINTAINER Alexander Troost <alexander.troost@kpn.com>

RUN mkdir /server
COPY . /server
WORKDIR /server

RUN go build -o avrserver .

CMD ["/server/avrserver 2498"]

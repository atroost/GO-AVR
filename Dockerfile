FROM golang:alpine as builder
MAINTAINER Alexander Troost <alexander.troost@kpn.com> / Gabor de Wit  <gabor.de.wit@accenture.com>

RUN apk update && apk add git && go get gopkg.in/natefinch/lumberjack.v2 && go get github.com/rs/zerolog/log && go get github.com/eclipse/paho.mqtt.golang
RUN mkdir /server
RUN mkdir /server/certs
RUN mkdir /server/config
COPY *.go /server/
COPY ./certs-prod/server.key /server/certs
COPY ./certs-prod/server.pem /server/certs
COPY ./config/serverConfig.json /server/config
WORKDIR /server

EXPOSE 2499

RUN go build -o avrserver .

ENTRYPOINT ["./avrserver"]

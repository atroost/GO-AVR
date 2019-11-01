#!/bin/bash
# call this script with an email address (valid or not).
# like:
# ./makecert.sh alexander.troost@kpn.com
mkdir certs
rm certs/*
echo "make server cert"
openssl req -new -nodes -x509 -out certs/server.pem -keyout certs/server.key -days 365 -subj "/C=NL/ST=ZH/L=Rotterdam/O=KPN N.V./OU=CTDO SP TV/CN=avr.stb.tv.kpn.com/emailAddress=$1"
echo "make client cert"
openssl req -new -nodes -x509 -out certs/client.pem -keyout certs/client.key -days 365 -subj "/C=NL/ST=ZH/L=Rotterdam/O=KPN N.V./CN=avr.stb.tv.kpn.com/emailAddress=$1"

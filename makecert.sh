#!/bin/bash
# just call this script!
# like:
# ./makecert.sh
mkdir certs
rm certs/*
echo "make server cert"
openssl req -newkey rsa:4096 -nodes -x509 -out certs/server.pem -keyout certs/server.key -days 365 -subj "/C=NL/ST=ZH/L=Rotterdam/O=KPN N.V./OU=CTDO SP TV/CN=avr.stb.tv.kpn.com"
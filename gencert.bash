#!/bin/bash


echo "creating server.key"
openssl genrsa -out server.key 2048
echo "creating server.crt"
openssl req -new -x509 -sha256 -key server.key -out server.crt -batch -days 365
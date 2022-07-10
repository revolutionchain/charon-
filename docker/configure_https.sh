#!/bin/sh

if [ -f "/https/cert.pem" ]; then
    echo "cert.pem already exists"
    if [ ! -f "/https/key.pem" ]; then
        echo "key.pem does not exist, both are needed, please delete "`pwd`"/https/key.pem manually"
        exit 1
    fi
    exit 0
fi

mkdir -p /https
if [ ! -d "/https" ]; then
    echo "Failed to mkdir -p /https"
    exit 1
fi

echo "Generating key.pem and cert.pem"
openssl req -nodes  -x509 -newkey rsa:4096 -keyout /https/key.pem -out /https/cert.pem -days 365 -subj "/C=US/ST=ST/L=L/O=Janus Self-signed https/OU=Janus Self-signed https/CN=Janus Self-signed https"
if [ 0 -ne $? ]; then
    echo "Failed to generate server.key"
    exit $?
fi

exit 0

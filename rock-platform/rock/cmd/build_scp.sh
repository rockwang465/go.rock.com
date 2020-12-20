#!/bin/bash

echo -e "go build:"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

echo -e "mv cmd rock"
mv cmd rock

echo -e "scp rock to 10.151.3.87:/rock"
scp rock root@10.151.3.87:/rock

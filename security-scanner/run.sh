#!/bin/sh

docker stop security-scanner
docker rm security-scanner
docker run --rm -p 8080:8080 --name security-scanner rnddockerdev.azurecr.io/ifs/security-scanner 
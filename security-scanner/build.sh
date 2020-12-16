#!/bin/sh

docker rmi security-scanner
docker build -f Dockerfile ../ -t rnddockerdev.azurecr.io/ifs/security-scanner
docker push rnddockerdev.azurecr.io/ifs/security-scanner
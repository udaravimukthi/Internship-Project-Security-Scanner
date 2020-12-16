#!/bin/sh

docker rmi security-scanner-clair
docker build . -t rnddockerdev.azurecr.io/ifs/security-scanner-clair
docker push rnddockerdev.azurecr.io/ifs/security-scanner-clair
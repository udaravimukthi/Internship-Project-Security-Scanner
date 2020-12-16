#!/bin/sh

docker rmi security-scanner-falco
docker build . -t rnddockerdev.azurecr.io/ifs/security-scanner-falco
docker push rnddockerdev.azurecr.io/ifs/security-scanner-falco
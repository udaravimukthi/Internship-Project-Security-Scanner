#!/bin/sh

docker rmi security-scanner-sysdig
docker build . -t rnddockerdev.azurecr.io/ifs/security-scanner-sysdig
docker push rnddockerdev.azurecr.io/ifs/security-scanner-sysdig
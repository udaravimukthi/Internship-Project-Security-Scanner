#!/bin/sh

docker rmi security-scanner-csysdig
docker build . -t rnddockerdev.azurecr.io/ifs/security-scanner-csysdig
docker push rnddockerdev.azurecr.io/ifs/security-scanner-csysdig
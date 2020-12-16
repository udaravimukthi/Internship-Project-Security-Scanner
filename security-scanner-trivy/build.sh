#!/bin/sh

docker rmi security-scanner-trivy
docker build . -t rnddockerdev.azurecr.io/ifs/security-scanner-trivy
docker push rnddockerdev.azurecr.io/ifs/security-scanner-trivy
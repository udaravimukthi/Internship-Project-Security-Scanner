#!/bin/sh

docker stop security-scanner-trivy
docker rm security-scanner-trivy
docker run --rm -it --name security-scanner-trivy rnddockerdev.azurecr.io/ifs/security-scanner-trivy
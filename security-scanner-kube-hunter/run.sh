#!/bin/sh

docker stop security-scanner-kube-hunter
docker rm security-scanner-kube-hunter
docker run --rm -it --name security-scanner-kube-hunter rnddockerdev.azurecr.io/ifs/security-scanner-kube-hunter
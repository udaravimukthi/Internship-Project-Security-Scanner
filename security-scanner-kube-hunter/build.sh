#!/bin/sh

docker rmi security-scanner-kube-hunter
docker build . -t rnddockerdev.azurecr.io/ifs/security-scanner-kube-hunter
docker push rnddockerdev.azurecr.io/ifs/security-scanner-kube-hunter
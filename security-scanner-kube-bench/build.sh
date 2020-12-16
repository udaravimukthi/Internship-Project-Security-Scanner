#!/bin/sh

docker rmi security-scanner-kube-bench
docker build . -t rnddockerdev.azurecr.io/ifs/security-scanner-kube-bench
docker push rnddockerdev.azurecr.io/ifs/security-scanner-kube-bench

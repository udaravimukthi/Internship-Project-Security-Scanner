#!/bin/sh

docker stop security-scanner-kube-bench
docker rm security-scanner-kube-bench
docker run --rm -it --name security-scanner-kube-bench rnddockerdev.azurecr.io/ifs/security-scanner-kube-bench
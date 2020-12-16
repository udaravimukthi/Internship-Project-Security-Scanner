#!/bin/sh

docker stop security-scanner-falco
docker rm security-scanner-falco
docker run --rm -it --name security-scanner-falco --privileged rnddockerdev.azurecr.io/ifs/security-scanner-falco
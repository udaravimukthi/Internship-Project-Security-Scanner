#!/bin/sh

docker stop security-scanner-csysdig
docker rm security-scanner-csysdig
docker run --rm -it --name security-scanner-csysdig --privileged rnddockerdev.azurecr.io/ifs/security-scanner-csysdig
# -v /var/run/docker.sock:/host/var/run/docker.sock \
# -v /dev:/host/dev \
# -v /proc:/host/proc:ro \
# -v /boot:/host/boot:ro \
# -v /lib/modules:/host/lib/modules:ro \
# -v /usr:/host/usr:ro \
# rnddockerdev.azurecr.io/ifs/security-scanner-sysdig 
#docker run -i -t --name security-scanner-sysdig --privileged -v /dev:/host/dev -v /proc:/host/proc:ro -v /boot:/host/boot:ro -v /lib/modules:/host/lib/modules:ro -v /usr:/host/usr:ro rnddockerdev.azurecr.io/ifs/security-scanner-sysdig
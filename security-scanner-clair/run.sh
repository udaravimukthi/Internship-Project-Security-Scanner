#!/bin/sh

docker stop security-scanner-clair
docker rm security-scanner-clair
docker run --rm -it -p 6060-6061:6060-6061 --name security-scanner-clair rnddockerdev.azurecr.io/ifs/security-scanner-clair -config=/config/config.yaml #rnddockerdev.azurecr.io/ifs/security-scanner-clair

#docker run --rm -it -p 6060-6061:6060-6061 -v /root/clair_config:/config quay.io/projectquay/clair:4.0.0-rc.11 -config=/config/config.yaml
#docker run -d -e POSTGRES_USER="postgres" POSTGRES_PASSWORD="password" POSTGRES_DB="postgres" -p 5432:5432 postgres:latest
#docker run --rm -it --name security-scanner-clair rnddockerdev.azurecr.io/ifs/security-scanner-clair

#docker run -d -e POSTGRES_PASSWORD="" -p 5432:5432 postgres:9.6
#docker run --net=host -d -p 6060-6061:6060-6061 -v $PWD/clair_config:/config quay.io/coreos/clair:latest -config=/config/config.yaml

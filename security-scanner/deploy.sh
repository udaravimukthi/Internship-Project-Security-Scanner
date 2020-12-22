#!/bin/sh

kubectl delete secret regcred -n security-scanner
kubectl apply -f ./
( cd ~ &&
kubectl create secret generic regcred -n security-scanner \
    --from-file=.dockerconfigjson=.docker/config.json \
    --type=kubernetes.io/dockerconfigjson ) 


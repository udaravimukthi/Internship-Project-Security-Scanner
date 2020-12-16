#!/bin/sh

echo "inside kube-bench"

#sleep 60

#sh

mkdir -p /scanresults

kube-bench | tee /scanresults/security-scanner-kube-bench-job-scan-results.txt

sleep 10

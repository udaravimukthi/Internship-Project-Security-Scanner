#!/bin/sh

echo "inside kube-hunter"

# echo 2 | kube-hunter --statistics
# kube-hunter --remote 10.1.72.74

mkdir -p /scanresults

kube-hunter --pod | tee /scanresults/security-scanner-kube-hunter-job-scan-results.txt

#sh

sleep 60

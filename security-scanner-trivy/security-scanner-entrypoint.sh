#!/bin/sh

curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/master/contrib/install.sh | sh -s -- -b /usr/local/bin

trivy --version

mkdir -p /scanresults
$COMMAND | tee /scanresults/security-scanner-trivy-job-scan-results.txt

sleep 60

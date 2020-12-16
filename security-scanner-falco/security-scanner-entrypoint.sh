#!/bin/sh

# curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/master/contrib/install.sh | sh -s -- -b /usr/local/bin

/docker-entrypoint.sh 
# touch ./events.txt
falco --version
falco -c /etc/falco/falco.yaml --disable-cri-async
mkdir -p /scanresults
$COMMAND | tee /scanresults/security-scanner-falco-job-scan-results.txt

sleep 60

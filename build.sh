#!/bin/sh

(cd security-scanner && ./build.sh)
(cd security-scanner-falco && ./build.sh)
(cd security-scanner-sysdig && ./build.sh)
(cd security-scanner-csysdig && ./build.sh)
(cd security-scanner-trivy && ./build.sh)
(cd security-scanner-clair && ./build.sh)
(cd security-scanner-kube-hunter && ./build.sh)
(cd security-scanner-kube-bench && ./build.sh)
#!/bin/sh

# apt-get update && apt-get upgrade -y
# apt-get install dialog        
# apt-get install -y --no-install-recommends apt-utils
# apt-get install -y curl
# apt-get install -y sudo
# apt --fix-broken install -y 
# #apt-get install linux-headers-3.10.0-1127.13.1.el7.x86_64
# apt-get -y install linux-headers-$(uname -r)
# # apt-get -y install linux-headers-$(3.10.0-1127.13.1.el7.x86_64)
# apt-get install sysdig -y
# curl -s https://s3.amazonaws.com/download.draios.com/stable/install-sysdig | sudo bash
# #apt-get install dkms

apt -y --fix-broken install
curl -s https://s3.amazonaws.com/download.draios.com/DRAIOS-GPG-KEY.public | apt-key add -  
curl -s -o /etc/apt/sources.list.d/draios.list https://s3.amazonaws.com/download.draios.com/stable/deb/draios.list  
apt-get update
apt-get -y install sysdig

/docker-entrypoint.sh

sysdig --version

mkdir -p /scanresults

touch /scanresults/capture.scap
chmod -R 777 /scanresults

sysdig -M $SCANSECONDS -w /scanresults/capture.scap
$COMMAND -r /scanresults/capture.scap | tee /scanresults/security-scanner-sysdig-job-scan-results.txt
sleep 60


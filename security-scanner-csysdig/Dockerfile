FROM rnddockerdev.azurecr.io/ifs/security-scanner-sysdig


# RUN apt-get update
# RUN apt-get install -y curl

# RUN apt --fix-broken install -y 
# RUN apt-get -y install linux-headers-$(uname -a)
# RUN apt-get -y install linux-headers-$(3.10.0-1127.13.1.el7.x86_64)
# RUN apt-get install sysdig -y
# RUN apt-get update && apt-get install -y --no-install-recommends apt-utils

WORKDIR /root/
COPY security-scanner-entrypoint.sh security-scanner-entrypoint.sh 
ENTRYPOINT ["./security-scanner-entrypoint.sh"]


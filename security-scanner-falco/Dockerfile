FROM falcosecurity/falco 
# RUN apk --no-cache add curl
# FROM ubuntu
# RUN apk update && apk add curl
# RUN curl -s https://falco.org/repo/falcosecurity-3672BA8F.asc | apt-key add -
# RUN echo "deb https://dl.bintray.com/falcosecurity/deb stable main" | tee -a /etc/apt/sources.list.d/falcosecurity.list
# # RUN apt --fix-broken install
# RUN apt-get update -y
# RUN apt --fix-broken install
# # RUN apt-get -y install linux-headers-$(uname -r)
# RUN apt-get install -y falco
RUN apt-get update
RUN apt --fix-broken install -y
RUN apt-get install mailutils -y

WORKDIR /root/
COPY security-scanner-entrypoint.sh security-scanner-entrypoint.sh 

COPY falco.yaml /etc/falco/falco.yaml

COPY falco_rules.yaml /etc/falco/falco_rules.yaml
COPY falco_rules.local.yaml /etc/falco/falco_rules.local.yaml
COPY k8s_audit_rules.yaml /etc/falco/k8s_audit_rules.yaml
COPY rules.d/ /etc/falco/rules.d 
COPY rules/ /etc/falco/rules

ENTRYPOINT ["./security-scanner-entrypoint.sh"]









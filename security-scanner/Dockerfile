FROM golang:latest
WORKDIR /go/src/ifs.com/ifs/security-scanner
COPY security-scanner/security-scanner.go security-scanner.go
RUN go get -d -v golang.org/x/net/html  
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o security-scanner .


FROM alpine:latest  
RUN apk --no-cache add ca-certificates curl  
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.18.8/bin/linux/amd64/kubectl \
    && mv ./kubectl /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl
    
ENV PATH="$PATH:/usr/local/bin"

WORKDIR /root/    

COPY --from=0 /go/src/ifs.com/ifs/security-scanner/security-scanner .
COPY security-scanner/security-scanner-deployment-pv.yaml.template security-scanner-deployment-pv.yaml
COPY security-scanner-sysdig/security-scanner-sysdig-job.yaml security-scanner-sysdig-job.yaml 
COPY security-scanner-trivy/security-scanner-trivy-job.yaml security-scanner-trivy-job.yaml
COPY security-scanner-clair/security-scanner-clair-job-service.yaml security-scanner-clair-job-service.yaml
COPY security-scanner-clair/security-scanner-clair-postgres-service.yaml security-scanner-clair-postgres-service.yaml
COPY security-scanner-clair/security-scanner-clair-postgres-deployment.yaml security-scanner-clair-postgres-deployment.yaml
COPY security-scanner-clair/security-scanner-clair-postgres-deployment-pv.yaml security-scanner-clair-postgres-deployment-pv.yaml
COPY security-scanner-clair/security-scanner-clair-postgres-pvc.yaml security-scanner-clair-postgres-pvc.yaml
COPY security-scanner-clair/security-scanner-clair-server-deployment.yaml security-scanner-clair-server-deployment.yaml
COPY security-scanner-clair/security-scanner-clair-job.yaml security-scanner-clair-job.yaml 
COPY security-scanner-csysdig/security-scanner-csysdig-deployment.yaml security-scanner-csysdig-deployment.yaml 
COPY security-scanner-kube-bench/security-scanner-kube-bench-job.yaml security-scanner-kube-bench-job.yaml 
COPY security-scanner-kube-hunter/security-scanner-kube-hunter-job.yaml security-scanner-kube-hunter-job.yaml 
COPY security-scanner-falco/security-scanner-falcosidekick-service.yaml security-scanner-falcosidekick-service.yaml
COPY security-scanner-falco/security-scanner-falcosidekick-deployment.yaml security-scanner-falcosidekick-deployment.yaml
COPY security-scanner-falco/security-scanner-falco-daemonset.yaml security-scanner-falco-daemonset.yaml
COPY security-scanner/static static
EXPOSE 8080
CMD ["./security-scanner"] 

FROM quay.io/coreos/clair:v2.1.6

# RUN apt-get update -y \
#     && apt-get install -y wget \
#     && apt-get install -y curl \
#     && apt install -y postgresql postgresql-contrib \
#     && apt install -y postgresql-client \
#     && apt install -y net-tools
#     # && apt-get install -y git \
#     # && apt-get install -y golang 

# WORKDIR /root/

RUN apk add curl \
    && apk add wget \
    && apk add docker 
    # && docker login -u udvilk -p udvilk123
    # && usermod -aG docker $USER
    # && addgroup -S appgroup && adduser -S appuser -G appgroup

# RUN wget https://github.com/jgsqware/clairctl/releases/download/v1.2.8/clairctl-linux-amd64 -O /usr/local/bin/clairctl \
#     && chmod 777 /usr/local/bin/clairctl
# RUN apk add docker-engine
RUN apk add go
RUN apk add --no-cache musl-dev

RUN export GOPATH=/usr/local/go && \
    git clone https://github.com/jgsqware/clairctl.git $GOPATH/src/github.com/jgsqware/clairctl && \
    cd $GOPATH/src/github.com/jgsqware/clairctl && \
    go get -u github.com/jteeuwen/go-bindata/... && \
    go generate ./clair && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build && \
    cp -v $GOPATH/src/github.com/jgsqware/clairctl/clairctl /usr/local/bin

RUN mkdir /clairctl \
    && cd /clairctl \   
    && wget https://raw.githubusercontent.com/jgsqware/clairctl/master/clairctl.yml.default -O clairctl.yaml \
    && sed -i 's~http://clair~http://clair-svc~g' /clairctl/clairctl.yaml \
    && chmod +x clairctl.yaml \
    && mkdir reports \
    && chmod +x reports

# RUN mkdir clair_config \
#     && mkdir reports

# RUN mkdir clair_config \
#     && curl -L https://raw.githubusercontent.com/coreos/clair/master/config.yaml.sample -o clair_config/config.yaml
# COPY config.yaml clair_config/config.yaml
COPY config.yaml /config/config.yaml
RUN chmod 777 /config/config.yaml
COPY security-scanner-entrypoint.sh /usr/local/bin/security-scanner-entrypoint.sh

EXPOSE 6060 6061

# ENV PATH="$PATH:/usr/local/go/bin"
# COPY config.go config.go

ENTRYPOINT ["/usr/local/bin/security-scanner-entrypoint.sh"]

# CMD ["-config=/config/config.yaml"]

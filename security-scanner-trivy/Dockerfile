FROM alpine:3.12  
RUN apk --no-cache add curl
WORKDIR /root/
COPY security-scanner-entrypoint.sh security-scanner-entrypoint.sh
CMD ["./security-scanner-entrypoint.sh"]









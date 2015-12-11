FROM ubuntu:trusty

ADD https://files.tutum.co/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ["ngrok", "tutum-agent", "/run.sh", "/"]
CMD ["/run.sh"]


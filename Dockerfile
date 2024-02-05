FROM ubuntu:22.04

RUN apt update
RUN apt install -y ca-certificates

COPY dns_resolver /dns_resolver

RUN chmod +x /dns_resolver

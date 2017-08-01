FROM alpine:latest

MAINTAINER Alex Wauck "alexwauck@exosite.com"
EXPOSE 8080

ENV GIN_MODE release

RUN apk add --no-cache curl

COPY build/linux-amd64/httpdebug /usr/bin/httpdebug
RUN chmod 0755 /usr/bin/httpdebug

CMD ["/usr/bin/httpdebug"]

FROM ubuntu:22.04

RUN apt update && apt upgrade -y
RUN apt install -y wget

ENV PATH=$PATH:/go/go/bin
ENV PATH=$PATH:/go/bin
ENV GOPATH=/go
ENV GOROOT=/go/go
ENV GOCACHE=/go/.cache/

WORKDIR /app

COPY build.sh /
RUN tr -d '\r' < /build.sh > /build.Unix.sh
RUN rm /build.sh
RUN mv /build.Unix.sh /build.sh
RUN chmod +x /build.sh

CMD [ "/build.sh" ]

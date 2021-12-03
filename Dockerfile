FROM golang:1.16-buster

RUN mkdir /src
WORKDIR /src
ENV GOPROXY https://goproxy.io

ENTRYPOINT  ["/src/build.sh"]

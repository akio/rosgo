FROM ros:kinetic-ros-base

ARG GOROOT

RUN apt-get update && apt-get install -y wget

#RUN wget -q https://dl.google.com/go/go1.11.linux-amd64.tar.gz
#RUN tar -C /usr/local -xzf go1.11.linux-amd64.tar.gz

VOLUME /usr/local/go

RUN mkdir -p src/github.com/akio/rosgo
COPY . src/github.com/akio/rosgo
COPY .travis/entrypoint.sh ./entrypoint.sh

CMD ./entrypoint.sh

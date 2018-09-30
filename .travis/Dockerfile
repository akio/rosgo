FROM ros:kinetic-ros-base

RUN wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz
RUN tar -C /usr/local https://dl.google.com/go/go1.11.linux-amd64.tar.gz


RUN mkdir -p src/github.com/akio/rosgo
COPY . src/github.com/akio/rosgo
COPY .travis/entrypoint.sh ./entrypoint.sh

CMD ./entrypoint.sh





FROM ros:kinetic-ros-base

RUN mkdir -p src/github.com/akio/rosgo
COPY . src/github.com/akio/rosgo
COPY ./travis/entrypoint.sh ./entrypoint.sh

CMD ./entrypoint.sh





ARG ROS_DOCKER=ros:melodic-ros-base
FROM $ROS_DOCKER

RUN apt-get update && apt-get install -y wget

VOLUME /usr/local/go

RUN mkdir -p src/github.com/fetchrobotics/rosgo
COPY . src/github.com/fetchrobotics/rosgo
COPY .travis/entrypoint.sh ./entrypoint.sh

CMD ./entrypoint.sh

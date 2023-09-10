FROM ubuntu:20.04

WORKDIR /opt/workdir

RUN --mount=type=bind,target=. \
    ./install.sh

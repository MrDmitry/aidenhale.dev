FROM ubuntu:20.04

RUN apt-get update && apt-get install -y expect && apt-get autoclean

WORKDIR /opt

COPY . /tmp/installer
RUN cd /tmp/installer && \
    ./expect_install && \
    rm -rf /tmp/installer

FROM ubuntu:20.04

RUN apt-get update && apt-get install -y expect && apt-get autoclean

WORKDIR /opt

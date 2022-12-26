# syntax=docker/dockerfile:1

FROM ubuntu:22.04

RUN apt-get update && apt install -y git libc6 libc-bin
RUN apt install -y gcc g++ make golang-go

WORKDIR /app

COPY bin/npv /bin/
COPY config.yaml ./

EXPOSE 3000

CMD [ "/bin/sh" ]
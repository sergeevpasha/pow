FROM golang:1.23

RUN curl -sSL https://taskfile.dev/install.sh | sh -s -- -d -b /usr/local/bin

RUN useradd -m -s /bin/bash admin

WORKDIR /home/app

RUN chown -R admin:admin /home/app

USER admin


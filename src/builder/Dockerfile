FROM lscr.io/linuxserver/code-server

ENV PUID=0
ENV PGID=0
ENV DEFAULT_WORKSPACE=/workspace

# install latest version of go
RUN apt-get update && \
    apt-get install -y software-properties-common wget && \
    rm -rf /var/lib/apt/lists/*
RUN add-apt-repository ppa:longsleep/golang-backports && \
    apt-get update
RUN apt-get install -y golang
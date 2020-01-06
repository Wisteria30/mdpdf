FROM ubuntu:18.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && \
        apt install wget texlive-luatex texlive-lang-cjk lmodern texlive-xetex git gcc make -y && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/*

RUN wget https://github.com/jgm/pandoc/releases/download/2.7.2/pandoc-2.7.2-1-amd64.deb && \
        dpkg -i pandoc-2.7.2-1-amd64.deb && \
        wget https://github.com/lierdakil/pandoc-crossref/releases/download/v0.3.4.1/linux-pandoc_2_7_2.tar.gz && \
        tar -zxvf linux-pandoc_2_7_2.tar.gz && \
        mv pandoc-crossref /usr/local/bin

RUN wget https://dl.google.com/go/go1.13.5.linux-amd64.tar.gz && \
        tar -C /usr/local -xzf go1.13.5.linux-amd64.tar.gz && \
        rm go1.13.5.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin

COPY scripts/ /usr/local/bin/
COPY crossref_config.yaml /config/crossref_config.yaml

WORKDIR /workdir

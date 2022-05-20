FROM debian:unstable 

RUN mkdir -p /root/rhad && \
    ln -fs /root/rhad /usr/local/rhad

WORKDIR /root/rhad

RUN apt-get update && apt-get install -y \
        golang \
        make && \
    rm -rf /var/cache/apt/*

COPY Makefile .
COPY *.go .
COPY go.mod .
# COPY go.sum .
COPY scripts/ .
COPY linters/ .
COPY tests/ .

RUN ls -l && make test clean

RUN make build && ln -fs build/linux-amd64/rhad ./rhad

RUN /usr/local/rhad/rhad sysinit

RUN mkdir -p /root/rhad/src
WORKDIR /root/rhad/src

ENTRYPOINT ["/usr/local/rhad/rhad"]
CMD ["run", "."]

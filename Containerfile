FROM debian:unstable

# Go's staticcheck linter complains unless this is set
ENV GOFLAGS -buildvcs=false

# rhad has a sysinit subcommand, but all the needed COPY directives are a lot to
# keep up with when files change, so just call the same script directly, early
COPY ./scripts/sysinit.sh ./scripts/sysinit.sh
RUN bash ./scripts/sysinit.sh

RUN adduser --gecos "" rhad
USER rhad
WORKDIR /home/rhad

COPY . .

RUN make test clean

RUN make build && \
    ln -fs build/linux-amd64/rhad ./rhad

RUN mkdir -p /home/rhad/src
WORKDIR /home/rhad/src

ENTRYPOINT ["/home/rhad/rhad"]
CMD ["run", "all"]

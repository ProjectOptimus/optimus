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

# Sets up the rest of the non-root-needed installs; the script checks if the runner is root or not
RUN bash ./scripts/sysinit.sh

RUN make test clean

RUN make build && \
    ln -fs build/linux-amd64/rhad ./rhad

RUN mkdir -p /home/rhad/src /home/rhad/.local/bin
WORKDIR /home/rhad/src

ENV PATH /home/rhad/.local/bin:${PATH}

ENTRYPOINT ["/home/rhad/rhad"]
CMD ["run", "all"]

FROM debian:unstable

# Go's staticcheck linter complains unless this is set
ENV GOFLAGS -buildvcs=false

# oscar will check for this to set its focus accordingly
ENV OSCAR_SRC '/home/oscar/oscar-src'

WORKDIR /root

COPY ./scripts/sysinit.sh ./scripts/sysinit.sh
RUN bash ./scripts/sysinit.sh

RUN useradd --create-home oscar
USER oscar
RUN mkdir -p \
      /home/oscar/oscar-src \
      /home/oscar/src \
      /home/oscar/.local/bin
WORKDIR /home/oscar/oscar-src

ENV OSCAR_SRC=/home/oscar/oscar-src

# Set up PATH correctly for oscar user (I can't find a better way to do this)
ENV PATH="/home/oscar/.local/bin:/home/oscar/go/bin:${PATH}"

# Sets up the rest of the non-root-needed installs; the script checks if the runner is root or not
COPY ./scripts/sysinit.sh ./scripts/sysinit.sh
RUN bash ./scripts/sysinit.sh

# Ok now hopefully we're all cached up
COPY --chown=oscar:oscar . .

RUN go mod tidy
RUN make test clean

RUN make build && \
    ln -fs build/linux-amd64/oscar ./oscar

WORKDIR /home/oscar/src

ENTRYPOINT ["/home/oscar/oscar-src/oscar"]

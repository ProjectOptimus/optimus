FROM debian:unstable

# Go's staticcheck linter complains unless this is set
ENV GOFLAGS -buildvcs=false

# rhad will check for this to set its focus accordingly
ENV RHAD_SRC '/home/rhad/rhad-src'

WORKDIR /root

COPY ./scripts/sysinit.sh ./scripts/sysinit.sh
RUN bash ./scripts/sysinit.sh

RUN useradd --create-home rhad
USER rhad
RUN mkdir -p \
      /home/rhad/rhad-src \
      /home/rhad/src \
      /home/rhad/.local/bin
WORKDIR /home/rhad/rhad-src

ENV RHAD_SRC=/home/rhad/rhad-src

# Set up PATH correctly for rhad user (I can't find a better way to do this)
ENV PATH="/home/rhad/.local/bin:/home/rhad/go/bin:${PATH}"

# Sets up the rest of the non-root-needed installs; the script checks if the runner is root or not
COPY ./scripts/sysinit.sh ./scripts/sysinit.sh
RUN bash ./scripts/sysinit.sh

# Ok now hopefully we're all cached up
COPY --chown=rhad:rhad . .

RUN go mod tidy
RUN make test clean

RUN make build && \
    ln -fs build/linux-amd64/rhad ./rhad

WORKDIR /home/rhad/src

ENTRYPOINT ["/home/rhad/rhad-src/rhad"]

FROM debian:unstable

# Go's staticcheck linter complains unless this is set
ENV GOFLAGS -buildvcs=false

COPY ./scripts/sysinit.sh ./scripts/sysinit.sh
RUN bash ./scripts/sysinit.sh

RUN useradd --create-home rhad
USER rhad
WORKDIR /home/rhad

# Set up PATH correctly for rhad user (I can't find a better way to do this)
ENV PATH="/home/rhad/.local/bin:/home/rhad/go/bin:${PATH}"

# Sets up the rest of the non-root-needed installs; the script checks if the runner is root or not
COPY ./scripts/sysinit.sh ./scripts/sysinit.sh
RUN bash ./scripts/sysinit.sh

# Ok now hopefully we're all cached up
COPY . .

RUN make test clean

RUN make build && \
    ln -fs build/linux-amd64/rhad ./rhad

RUN mkdir -p /home/rhad/src /home/rhad/.local/bin
WORKDIR /home/rhad/src

ENTRYPOINT ["/home/rhad/rhad"]
CMD ["run", "all"]

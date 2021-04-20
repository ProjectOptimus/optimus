FROM debian:unstable

WORKDIR /root

# init via script instead of Dockerfile lines, so rhad can still be installed
# via other targets/patterns
COPY init.sh init.sh
COPY tests/test-init.bats tests/test-init.bats
RUN bash init.sh

COPY . .
RUN make test

ENTRYPOINT ["/root/rhad"]
CMD ["."]

FROM debian:unstable

WORKDIR /root

# init via script instead of Dockerfile lines, so rhad can still be installed
# via other targets/patterns
COPY init.sh init.sh
COPY tests/test-init.bats tests/test-init.bats

RUN bash init.sh \
    && rm -rf /var/cache/apt/* \
    && rm -rf .cache/*

COPY . .
RUN make test \
    && rm -rf .mypy_cache

ENTRYPOINT ["/root/rhad"]
CMD ["."]

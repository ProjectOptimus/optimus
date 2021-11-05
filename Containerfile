FROM debian:unstable 

WORKDIR /root

# init via scripts instead of Containerfile lines, so rhad can still be
# installed via other targets/patterns
COPY scripts/init.sh scripts/init.sh
RUN bash scripts/init.sh && \
    rm -rf /var/cache/apt/* && \
    rm -rf *cache* .*cache*

COPY Makefile Makefile
COPY scripts/ scripts/
COPY linters/ linters/
COPY tests/ tests/

# rhad needs to be on PATH for its `realpath` call to work
RUN chmod +x /root/scripts/rhad && \
    ln -fs /root/scripts/rhad /usr/bin/rhad
RUN make test && \
    rm -rf *cache* .*cache*

RUN mkdir -p /root/src
WORKDIR /root/src

ENTRYPOINT ["/root/scripts/rhad"]
CMD ["."]

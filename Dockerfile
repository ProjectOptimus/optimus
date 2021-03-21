FROM debian:unstable

WORKDIR /root

COPY . .

# init via script instead of Dockerfile lines, so optimus can still be installed
# via other targets/patterns
RUN bash init.sh

ENTRYPOINT ["/root/optimus"]
CMD ["."]

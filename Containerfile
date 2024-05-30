FROM debian:12

# oscar will check for this to set its focus accordingly
ENV OSCAR_SRC='/home/oscar/oscar-src'

WORKDIR /root

RUN apt-get update && apt-get install -y \
      ansible-core \
      ansible-lint

COPY ./scripts/ansible-playbooks/ ./scripts/ansible-playbooks
RUN ansible-lint ./scripts/ansible-playbooks/main.yaml && \
    ansible-playbook ./scripts/ansible-playbooks/main.yaml

# Note: oscar user and the src dir are created in the init script
USER oscar
WORKDIR /home/oscar/src
COPY --chown=oscar:oscar . .
RUN go mod tidy
RUN make test clean
RUN make build && \
    ln -fs ./build/linux-amd64/oscar /oscar

ENTRYPOINT ["/oscar"]

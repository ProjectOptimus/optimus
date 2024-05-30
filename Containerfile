FROM debian:12

# oscar will check for this to set its focus accordingly
ENV OSCAR_SRC='/home/oscar/oscar-src'

WORKDIR /root

RUN apt-get update && apt-get install -y \
      ansible-core \
      ansible-lint

COPY ./scripts/ansible-playbooks/ ./scripts/ansible-playbooks
RUN ansible-lint ./scripts/ansible-playbooks/main.yaml

# Run based on tags, just so we can have Docker cache the layers
RUN ansible-playbook --tags core ./scripts/ansible-playbooks/main.yaml
RUN ansible-playbook --tags testers ./scripts/ansible-playbooks/main.yaml

# Note: oscar user and the src dir are created above
USER oscar
WORKDIR /home/oscar/oscar-src
COPY --chown=oscar:oscar . .
RUN go mod tidy
RUN make test clean
RUN make build

USER root
RUN ln -fs /home/oscar/oscar-src/build/linux-amd64/oscar /oscar

USER oscar
WORKDIR /home/oscar/src

ENTRYPOINT ["/oscar"]

#!/usr/bin/env bash
set -euo pipefail

export DEBIAN_FRONTEND=noninteractive

errorf() {
  printf "ERROR: %s\n" "$@" > /dev/stderr
  exit 1
}

init-sys() {
  apt-get update
  apt-get install -y \
    curl \
    git \
    golang \
    make \
    npm \
    python3 \
    python3-pip \
    python3-venv \
    ruby-full \
    shellcheck \
  || errorf "Could not init system packages for rhad!"
}

init-bats() {
  git clone https://github.com/bats-core/bats-core.git /tmp/bats
  bash /tmp/bats/install.sh /usr/local
}

init-go() {
  pkgs=(
    # golang.org/x/lint/golint
    honnef.co/go/tools/cmd/staticcheck
  )
  for pkg in "${pkgs[@]}"; do
    go install \
      "${pkg}"@latest
  done
  mkdir -p "${HOME}"/.local/bin/
  ln -fs "$(go env GOPATH)"/bin/* "${HOME}"/.local/bin/
}

init-python() {
  pip3 install --user \
    black \
    mypy \
    pylint \
    pytest \
    pytest-cov \
  || errorf "Could not init Python packages for rhad!"
}

init-ruby() {
  gem install \
    mdl \
  || errorf "Could not init Ruby packages for rhad!"
}

main() {
  if [[ $(id -u) -eq 0 ]]; then
    init-sys
    init-bats
  else  
    init-go
    init-python
    init-ruby
  fi
}

main || errorf "Failed to initialize rhad host!"

exit 0

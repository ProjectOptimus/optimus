#!/usr/bin/env bash
set -euo pipefail

# This script prepares (and validates) the rhad host

export DEBIAN_FRONTEND=noninteractive

errorf() {
  printf "ERROR: %s\n" "$@" > /dev/stderr
  exit 1
}

init-sys() {
  local pkgs=(
    build-essential
    curl
    golang
    git
    make
    python3
    python3-pip
  )
  {
    apt-get update
    apt-get install -y "${pkgs[@]}"
  } || errorf "Could not init system packages for rhad!"
}

init-bats() {
  git clone https://github.com/bats-core/bats-core.git /tmp/bats
  bash /tmp/bats/install.sh /usr/local
  rm -rf /tmp/bats
}

init-go() {
  local pkgs=(
    github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  )
  for pkg in "${pkgs[@]}"; do
    go install "${pkg}"
  done

  mkdir -p "${HOME}"/.local/bin/
  ln -fs "$(go env GOPATH)"/bin/* "${HOME}"/.local/bin/
}

init-python() {
  local pkgs=(
    pytest
    pytest-cov
  )
  pip3 install --user "${pkgs[@]}" \
  || errorf "Could not init Python packages for rhad!"
}

test-sysinit() {
  local failed=""
  local cmds=(
    curl
    git
    go
    golangci-lint
    make
    python3
    pytest
  )
  for cmd in "${cmds[@]}"; do
    command -v "${cmd}" >/dev/null || {
      printf 'ERROR: Command "%s" not found on PATH\n' "${cmd}"
      failed=true
    }
  done
  if [[ -n "${failed}" ]]; then
    errorf '^ Above command(s) not found on PATH -- did you run the sysinit script for rhad?'
  fi
}

main() {
  if [[ $(id -u) -eq 0 ]]; then
    init-sys
    init-bats
  else  
    init-go
    init-python
    
    # Also run tests as nonroot, so setup is confirmed for the least-privileged user
    test-sysinit
  fi
}

# Allow sourcing the file to run e.g. test-sysinit() by itself
if [[ "${1:-}" == "test" ]]; then
  test-sysinit
elif [[ "${BASH_SOURCE[0]:-}" == "${0}" ]]; then
  main || errorf "Failed to initialize rhad host!"
fi

exit 0

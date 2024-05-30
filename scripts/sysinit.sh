#!/usr/bin/env bash
set -euo pipefail

# This script prepares (and validates) the oscar host

export DEBIAN_FRONTEND=noninteractive

errorf() {
  printf "ERROR: %s\n" "$@" > /dev/stderr
  exit 1
}

apt-get update
apt-get install -y \
  ansible-core \
  ansible-lint

###

ansible-playbook ./ansible-playbooks/main.yaml

# init-sys() {
#   local pkgs=(
#     build-essential
#     curl
#     golang
#     git
#     make
#     python3
#     python3-pip
#     terraform
#   )
#   {
#     apt-get update
#     apt-get install -y "${pkgs[@]}"
#   } || errorf "Could not init system packages for oscar!"
# }

# init-bats() {
#   git clone https://github.com/bats-core/bats-core.git /tmp/bats
#   bash /tmp/bats/install.sh /usr/local
#   rm -rf /tmp/bats
# }

# init-go() {
#   local pkgs=(
#     github.com/golangci/golangci-lint/cmd/golangci-lint@latest
#   )
#   for pkg in "${pkgs[@]}"; do
#     go install "${pkg}"
#   done

#   mkdir -p "${HOME}"/.local/bin/
#   ln -fs "$(go env GOPATH)"/bin/* "${HOME}"/.local/bin/
# }

# init-python() {
#   local pkgs=(
#     python3-pytest
#     python3-pytest-cov
#   )
#   apt-get install -y "${pkgs[@]}" \
#   || errorf "Could not init Python packages for oscar!"
# }

# test-sysinit() {
#   local failed=""
#   local cmds=(
#     curl
#     git
#     go
#     make
#     python3
#     pytest
#   )
#   for cmd in "${cmds[@]}"; do
#     command -v "${cmd}" >/dev/null || {
#       printf 'ERROR: Command "%s" not found on PATH\n' "${cmd}"
#       failed=true
#     }
#   done
#   if [[ -n "${failed}" ]]; then
#     errorf '^ Above command(s) not found on PATH -- did you run the sysinit script for oscar?'
#   fi
# }

# main() {
#   if [[ $(id -u) -eq 0 ]]; then
#     init-sys
#     init-bats
#     init-python
#   else
#     init-go

#     # Also run tests as nonroot, so setup is confirmed for the least-privileged user
#     test-sysinit
#   fi
# }

# # Allow sourcing the file to run e.g. test-sysinit() by itself
# if [[ "${1:-}" == "test" ]]; then
#   test-sysinit
# elif [[ "${BASH_SOURCE[0]:-}" == "${0}" ]]; then
#   main || errorf "Failed to initialize oscar host!"
# fi

# exit 0

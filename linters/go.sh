#!/usr/bin/env bash
set -euo pipefail

# shellcheck disable=SC1091
source "${RHAD_ROOT:-}"/scripts/utils.sh

# Fail if gofmt finds any diff from its own rules
printf '>>>>>>>> Running Go formatter diff check...\n'
gofmt_out=$(gofmt -d "$1")
if [[ "${#gofmt_out}" -gt 0 ]]; then
  printf '%s\n' "${gofmt_out}"
  mark-failed-lint 'go-gofmt'
fi

printf '>>>>>>>> Running Go linter...\n'
if [[ -d "$1" ]]; then
  staticcheck "$1"/... || mark-failed-lint 'go-staticcheck'
else
  staticcheck "$1" || mark-failed-lint 'go-staticcheck'
fi

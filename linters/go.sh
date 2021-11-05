#!/usr/bin/env bash
set -euo pipefail

# Fail if gofmt finds any diff from its own rules
printf '>>>>>>>> Running Go formatter diff check...\n'
gofmt_out=$(gofmt -d "$1")
printf '%s\n' "${gofmt_out}"
if [[ "${#gofmt_out}" -gt 0 ]]; then
  exit 1
fi

printf '>>>>>>>> Running Go linter...\n'
if [[ -d "$1" ]]; then
  staticcheck "$1"/...
else
  staticcheck "$1"
fi

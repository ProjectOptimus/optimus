#!/usr/bin/env bash
set -euo pipefail

tests_dir="$(realpath "$(dirname "${BASH_SOURCE[0]}")")"

cd "${tests_dir}" || {
  printf 'ERROR: Could not resolve path to test files (%s) on host!\n' "${tests_dir}" > /dev/stderr
  exit 1
}

bats ./

#!/usr/bin/env bash
set -euo pipefail

file_list=$(find "$1" -type f -name '*.py')

for f in ${file_list}; do
  printf "Running Python linters against %s...\n" "${f}"
  pylint "${f}"
  mypy "${f}"
done

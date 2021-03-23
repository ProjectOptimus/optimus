#!/usr/bin/env bash
set -euo pipefail

file_list=$(find "$1" -type f -name '*.md')

for f in ${file_list}; do
  printf "Running Markdown linter against %s...\n" "${f}"
  mdl "${f}"
done

#!/usr/bin/env bash
set -euo pipefail

file_list=$(find "$1" -type f -regextype posix-extended -regex '.*\.sh|.*\.bash|.*\.bats')

for f in ${file_list}; do
  printf ">>>>>>>> Running shell linter against %s...\n" "${f}"
  shellcheck "${f}"
done

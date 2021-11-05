#!/usr/bin/env bash
set -euo pipefail

# TODO: Find a better way to do a full directory scan for pylint, since it's mad
# if a dir arg isn't inside a formal package
file_list=$(find "$1" -type f -name '*.py')
for f in ${file_list}; do
  printf '>>>>>>>> Running Python linter against %s...\n' "${f}"
  pylint "${f}"
done 

printf '>>>>>>>> Running Python typechecker...\n'
mypy "$1"

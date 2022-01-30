#!/usr/bin/env bash
set -euo pipefail

# shellcheck disable=SC1091
source "${RHAD_ROOT:-}"/scripts/utils.sh

# TODO: Find a better way to do a full directory scan for pylint, since it's mad
# if a dir arg isn't inside a formal package
file_list=$(find "$1" -type f -name '*.py')

if [[ "${#file_list}" -eq 0 ]]; then
  exit 0
else

  for f in ${file_list}; do
    printf '>>>>>>>> Running Python linter against %s...\n' "${f}"
    pylint \
      --disable=import-error,invalid-name \
      "${f}" \
    || mark-failed-lint 'python-pylint'
  done

  printf '>>>>>>>> Running Python typechecker...\n'
  mypy "$1" || mark-failed-lint 'python-mypy'

  # Fail if Black finds any diff from its own rules
  printf '>>>>>>>> Running Python Black formatter diff check...\n'
  black_diff_out=$(black --diff "$1") > /dev/null 2>&1
  if [[ "${#black_diff_out}" -gt 0 ]]; then
    printf '%s\n' "${black_diff_out}"
    mark-failed-lint 'python-black'
  fi

fi

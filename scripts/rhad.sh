#!/usr/bin/env bash
set -euo pipefail

# Script to kick off rhad's utilities. Serves as both a reference as well as a
# means to run the rhad subcommands on a local machine.

cmd="${1}"
if [[ -z "${cmd}" ]]; then
  printf 'ERROR: must specify a valid rhad subcommand as script arg\n'
fi

github-super-linter() {
  # Clean up after any last runs, because it'll creep on its own files
  rm -rf ./super-linter.log ./*cache*

  local result=0

  # Lots of arg setting, like e.g. disabling Go linting because super-linter's
  # Go lint runs don't respect packages/multiple files
  docker run \
    --rm \
    -it \
    --user "$(id -u):$(id -g)" \
    -e RUN_LOCAL=true \
    -e USE_FIND_ALGORITHM=true \
    -e IGNORE_GITIGNORED_FILES=true \
    -e FILTER_REGEX_EXCLUDE='testdata.*' \
    -e FILTER_REGEX_EXCLUDE='.*\.vmdk|.*\.ovf|.*\.box|.*\.iso' \
    -e VALIDATE_GO=false \
    -e VALIDATE_NATURAL_LANGUAGE=false \
    -v "${PWD}":/tmp/lint \
    docker.io/github/super-linter:slim-v4 \
  || result="$?"

  # Don't want log hanging around locally if everything was fine
  if [[ "${result}" -eq 0 ]]; then
    rm super-linter.log
  fi

  return "${result}"
}

rhad-lint() {
  local result=0
  docker run \
    --rm -it \
    -v "${PWD}":/home/rhad/src \
    ghcr.io/opensourcecorp/rhad:latest \
  lint \
  || result="$?"
  return "${result}"
}

rhad-test() {
  printf 'Not implemented.\n'
  return 1
}

rhad-build() {
  printf 'Not implemented.\n'
  return 1
}

case "${cmd}" in
  lint)
    github-super-linter
    rhad-lint
  ;;
  test)
    rhad-test
  ;;
  build)
    rhad-build
  ;;
  *)
    printf 'ERROR: must specify a valid rhad subcommand as script arg\n'
  ;;
esac

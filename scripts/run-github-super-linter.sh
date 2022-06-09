#!/usr/bin/env bash
set -euo pipefail

# This script serves as a reference call for how GitHub's Super-Linter is run.
# rhad will focus its linting efforts on any gaps in the Super-Linter

# Clean up after any last runs, because it'll creep on its own files
rm -rf ./super-linter.log ./*cache*

result=0

# Lots of arg setting, like e.g. disabling Go linting because super-linter's Go
# lint runs don't respect packages/multiple files
docker run \
  --rm \
  -it \
  --user "$(id -u):$(id -g)" \
  -e RUN_LOCAL=true \
  -e USE_FIND_ALGORITHM=true \
  -e IGNORE_GITIGNORED_FILES=true \
  -e FILTER_REGEX_EXCLUDE='testdata.*' \
  -e VALIDATE_GO=false \
  -e VALIDATE_NATURAL_LANGUAGE=false \
  -v "${PWD}":/tmp/lint \
  docker.io/github/super-linter:slim-v4 \
|| result="$?"

# Don't want log hanging around locally if everything was fine
if [[ "${result}" -eq 0 ]]; then
  rm super-linter.log
fi

exit "${result}"

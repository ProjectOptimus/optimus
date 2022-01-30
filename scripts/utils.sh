#!/usr/bin/env bash
set -euo pipefail

# Need to know where rhad root ACTUALLY lives on disk -- this looks gross, but
RHAD_ROOT="$(dirname "$(dirname "$(realpath "$(command -v rhad)")")")"
export RHAD_ROOT

failed_lint_logfile='/tmp/rhad-failed-lints'

mark-failed-lint() {
  printf "%s\n" "$1" >> "${failed_lint_logfile}"
}

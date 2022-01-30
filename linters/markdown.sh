#!/usr/bin/env bash
set -euo pipefail

# shellcheck disable=SC1091
source "${RHAD_ROOT:-}"/scripts/utils.sh

printf '>>>>>>>> Running Markdown linter...\n'
mdl --style /root/linters/.mdlrc.style.rb "$1" || mark-failed-lint 'markdown-markdownlint'

#!/usr/bin/env bash
set -euo pipefail

# shellcheck disable=SC1091
source "${RHAD_ROOT:-}"/scripts/utils.sh

printf '>>>>>>>> Running Markdown linter...\n'

# Find list of files to lint, but ignore the ones that should be... well,
# ignored; e.g. terraform-docs-generated .md files
find "$1" -type f -name '*.md' -or -name '*.markdown' > /tmp/rhad-mdl-files
grep -vE '.*providers/.*/docs/README.md' /tmp/rhad-mdl-files > /tmp/rhad-mdl-files-cleaned

# shellcheck disable=SC2046
mdl \
  --style /root/linters/.mdlrc.style.rb \
  $(paste -s -d' ' /tmp/rhad-mdl-files-cleaned) \
|| mark-failed-lint 'markdown-markdownlint'

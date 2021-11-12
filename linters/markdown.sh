#!/usr/bin/env bash
set -euo pipefail

printf '>>>>>>>> Running Markdown linter...\n'
mdl --style /root/linters/.mdlrc.style.rb "$1"

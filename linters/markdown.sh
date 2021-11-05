#!/usr/bin/env bash
set -euo pipefail

printf '>>>>>>>> Running Markdown linter...\n'
mdl "$1"

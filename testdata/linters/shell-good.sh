#!/usr/bin/env bash
set -euo pipefail

echo "${var1:-unset}" > /dev/null

SHELL = /usr/bin/env bash -euo pipefail

test:
	@printf "\n=== Running optimus tests ===\n\n"
	@bats ./tests/test-optimus-linters.bats
	@printf "============ DONE ===========\n\n"

docker-build:
	@docker build -t optimus:latest .

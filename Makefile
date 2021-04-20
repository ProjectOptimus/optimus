SHELL = /usr/bin/env bash -euo pipefail

test:
	@printf "\n=== Running rhad tests ===\n\n"
	@bats ./tests/test-rhad-linters.bats
	@printf "============ DONE ===========\n\n"

docker-build:
	@docker build -t rhad:latest .

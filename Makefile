SHELL = /usr/bin/env bash -euo pipefail
DOCKER ?= docker

test:
	@printf "\n======== Running rhad tests\n\n"
	@bash ./tests/test-all.sh
	@printf "======== DONE\n\n"

image-build:
	@$(DOCKER) build -f Containerfile -t opensourcecorp/rhadamanthus:latest .

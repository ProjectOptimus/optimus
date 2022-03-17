SHELL = /usr/bin/env bash -euo pipefail
DOCKER ?= docker
REGISTRY ?= sauce.opensourcecorp.org

test:
	@printf "\n================ Running rhad tests\n\n"
	@bash ./tests/test-all.sh
	@printf "\n================ DONE\n\n"

image-build:
	@$(DOCKER) build -f Containerfile -t $(REGISTRY)/opensourcecorp/rhadamanthus:latest .

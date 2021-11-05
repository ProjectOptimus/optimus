SHELL = /usr/bin/env bash -euo pipefail

test:
	@printf "\n======== Running rhad tests\n\n"
	@bash ./tests/test-all.sh
	@printf "======== DONE\n\n"

image-build:
	@docker build -f Containerfile -t opensourcecorp/rhad:latest .

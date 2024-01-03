SHELL = /usr/bin/env bash -euo pipefail

all: test

.PHONY: %

test:
	@OSC_IS_TESTING=true go test -v -cover ./...

SHELL = /usr/bin/env bash -euo pipefail

PKGNAME := oscar
BINNAME := oscar

DOCKER ?= docker
OCI_REGISTRY ?= ghcr.io# ociregistry.opensourcecorp.org
OCI_REGISTRY_OWNER ?= opensourcecorp# library

all: test clean

.PHONY: %

test: clean
	@printf 'Running go vet first...\n' && go vet ./... && printf 'Ok.\n'
	@OSC_IS_TESTING=true go test -cover ./...

build: clean
	@mkdir -p build/$$(go env GOOS)-$$(go env GOARCH)
	@go build -o build/$$(go env GOOS)-$$(go env GOARCH)/$(BINNAME)

xbuild: clean
	@for target in \
		darwin-amd64 \
		linux-amd64 \
		linux-arm \
		linux-arm64 \
		windows-amd64 \
	; \
	do \
		GOOS=$$(echo "$${target}" | cut -d'-' -f1) ; \
		GOARCH=$$(echo "$${target}" | cut -d'-' -f2) ; \
		outdir=build/"$${GOOS}-$${GOARCH}" ; \
		mkdir -p "$${outdir}" ; \
		printf "Building for %s-%s into build/ ...\n" "$${GOOS}" "$${GOARCH}" ; \
		GOOS="$${GOOS}" GOARCH="$${GOARCH}" go build -o "$${outdir}"/$(BINNAME) ; \
	done

package: xbuild
	@mkdir -p dist
	@cd build || exit 1; \
	for built in * ; do \
		printf 'Packaging for %s into dist/ ...\n' "$${built}" ; \
		cd $${built} && tar -czf ../../dist/$(PKGNAME)_$${built}.tar.gz * && cd - >/dev/null ; \
	done

clean:
	@rm -rf \
		/tmp/$(PKGNAME)-tests \
		*cache* \
		.*cache* \
		build/ \
		dist/

# clean as a dep because of mounted volume permissions issues when it tries to
# run 'make clean' within the container
image-build: clean
	@$(DOCKER) build -f Containerfile -t $(OCI_REGISTRY)/$(OCI_REGISTRY_OWNER)/$(PKGNAME):latest .

# Some targets that help set up local workstations with oscar tooling. Assumes
# ~/.local/bin is on $PATH
add-local-symlinks:
	@mkdir -p "${HOME}"/.local/bin
	@ln -fs $(realpath ./scripts/oscar.sh) "${HOME}"/.local/bin/oscar
	@printf 'Symlinked oscar runner script to %s\n' "${HOME}"/.local/bin/oscar

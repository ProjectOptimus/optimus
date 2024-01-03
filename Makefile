SHELL = /usr/bin/env bash -euo pipefail

PKGNAME := oscar
BINNAME := oscar

DOCKER ?= docker
# OCI_REGISTRY ?= ociregistry.opensourcecorp.org
# OCI_REGISTRY_OWNER ?= library
OCI_REGISTRY ?= ghcr.io
OCI_REGISTRY_OWNER ?= opensourcecorp

all: test clean

.PHONY: %

test: clean
	@printf 'Running go vet first...\n' && go vet ./... && printf 'Ok.\n'
	@go test -cover ./...

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

oscar
====

<!-- badges: start -->
![Github Actions](https://github.com/opensourcecorp/oscar/actions/workflows/main.yaml/badge.svg)

[![Support OpenSourceCorp on Ko-Fi!](https://img.shields.io/badge/Ko--fi-F16061?style=for-the-badge&logo=ko-fi&logoColor=white)](https://ko-fi.com/ryapric)
<!-- badges: end -->

`oscar` ("`oscaramanthus`") is a CI/CD task runner used in [OpenSourceCorp's CI/CD
subsystem](https://github.com/opensourcecorp/osc-infra/tree/main/cicd). It does
not orchestrate CI/CD tasks -- that's the subsystem's job. Rather, it is a set
of utilities that are designed to be easily ported between any CI/CD platform of
your choosing. `oscar` comprises all the CI/CD logic that your platform would
normally run as steps in that process -- lint, test, build, push, deploy, etc.
In this way, you can think of `oscar` being much like a Jenkins shared library.

How to use
----------

To build & run `oscar` locally, you can clone this repo, and run:

    make image-build

from the repo root. This will build (by default) an image tagged as
`ghcr.io/opensourcecorp/oscar:latest`. Please be patient, as `oscar` has a lot of
build-time dependencies that it needs to fetch; and note that the resulting
image will be quite large!

To actually run `oscar`, you will need to run the image's container with your
local folder mounted to it:

    docker run --rm -it -v "${PWD}":/home/oscar/src ghcr.io/opensourcecorp/oscar:latest <subcommand> # e.g. 'lint'

Note that `oscar`'s instructions are provided within a container runtime context
only. As `oscar` depends on many system & CLI utilities being present at runtime,
it is an unfair assumption that someone will have their host machine configured
with all of these disparate tools (there are a LOT, mostly for the linters).

If you really want to get `oscar` working on a dedicated machine, it's certainly
easy enough to do -- just fire up a Debian-based machine, grab the
`scripts/sysinit.sh` script, and take note of the order of things specified in
the top-level `Containerfile`. Note that `oscar`'s container image is built off
of Debian's unstable/"Sid" branch, and has not been tested on a stable release.

rhad
====

<!-- badges: start -->
![Github Actions](https://github.com/opensourcecorp/rhadamanthus/actions/workflows/rhad.yaml/badge.svg)

[![Support OpenSourceCorp on Ko-Fi!](https://img.shields.io/badge/Ko--fi-F16061?style=for-the-badge&logo=ko-fi&logoColor=white)](https://ko-fi.com/ryapric)
<!-- badges: end -->

---

*King Rhadamanthus has found you worthy -- King Minos, God of War 3*

---

`rhadamanthus` ("`rhad`") is a CI/CD task runner used in OpenSourceCorp's CI/CD
subsystem. It does not orchestrate CI/CD tasks -- rather, it is a set of
utilities that are designed to be easily ported between any CI/CD platform of
your choosing. You can think of `rhad` as all the CI/CD logic that your platform
would normally run as steps in that process -- lint, test, build, push, deploy,
etc. In this way, `rhad` is much like a Jenkins shared library.

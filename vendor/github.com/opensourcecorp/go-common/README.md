Common Go library for OpenSourceCorp
====================================

<!-- badges: start -->
![Github Actions](https://github.com/opensourcecorp/go-common/actions/workflows/main.yaml/badge.svg)

[![Support OpenSourceCorp on Ko-Fi!](https://img.shields.io/badge/Ko--fi-F16061?style=for-the-badge&logo=ko-fi&logoColor=white)](https://ko-fi.com/ryapric)
<!-- badges: end -->

This repository houses common functionality to be used across Go packages
authored by OpenSourceCorp.

Suppressing log messages in your own program's tests
----------------------------------------------------

If you want or need to suppress log messages from this package's logging
facilities, set the env var `OSC_IS_TESTING=true` for your test runs. This will
suppress all but `FatalLog()` messages.

If you need to trick the loggers into thinking you're *not* testing for some
reason, you can temporarily set `osc.IsTesting = false` during any test run, and
then re-set it to `true`.

#!/usr/bin/env bats

function make-test() {
  ./optimus tests/testfiles/"$1"
  run ./optimus tests/testfiles/"$2"
  [ "${status}" -ne 0 ]
  if [[ "${output}" == *"such file or directory"* ]]; then
    printf "%s\n" "${output}"
    exit 1
  fi
}

@test "can lint shell" {
  make-test shell-{good,bad}.sh
}

@test "can lint Python" {
  make-test python{_good,-bad}.py
}

@test "can lint markdown" {
  make-test markdown-{good,bad}.md
}

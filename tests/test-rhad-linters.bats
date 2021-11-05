#!/usr/bin/env bats

make-test() {
  # Direct invocation will fail on tests that SHOULD pass
  ../scripts/rhad ../tests/testfiles/"$1" "$3"
  # Invocation via `run` in bats will pass on tests that SHOULD fail
  run ../scripts/rhad ../tests/testfiles/"$2" "$3"

  [[ "${status}" -ne 0 ]]
  if [[ "${output}" == *"such file or directory"* ]]; then
    printf "%s\n" "${output}"
    exit 1
  fi
}

@test "Can lint Shell" {
  make-test shell-{good,bad}.sh shell
}

@test "Can lint Go" {
  make-test go_{good,bad}.go go
}

@test "Can lint Python" {
  make-test python{_good,-bad}.py python
}

@test "Can lint Markdown" {
  make-test markdown-{good,bad}.md markdown
}

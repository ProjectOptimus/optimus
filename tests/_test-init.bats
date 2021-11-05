#!/usr/bin/env bats

@test "init-sys() installed everything you expected, and it's all on the \$PATH" {
  cmds=(
    curl
    git
    go
    make
    npm
    python3
    pip3
    ruby
    shellcheck
  )
  for cmd in "${cmds[@]}"; do
    command -v "${cmd}" >/dev/null || {
      printf "Command '%s' not found\n" "${cmd}"
      exit 1
    }
  done

  python3 -m venv -h > /dev/null
}

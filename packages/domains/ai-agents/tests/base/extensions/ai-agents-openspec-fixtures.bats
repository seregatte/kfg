#!/usr/bin/env bats

FIXTURES="$BATS_TEST_DIRNAME/../../../tests/fixtures/openspec"

setup() {
  if ! command -v kfg >/dev/null 2>&1; then
    skip "kfg not installed"
  fi
}

@test "openspec local fixture builds" {
  run kfg build "$FIXTURES/local"
  [ "$status" -eq 0 ]
}

@test "openspec store fixture builds" {
  run kfg build "$FIXTURES/store"
  [ "$status" -eq 0 ]
}

@test "openspec references fixture builds" {
  run kfg build "$FIXTURES/references"
  [ "$status" -eq 0 ]
}

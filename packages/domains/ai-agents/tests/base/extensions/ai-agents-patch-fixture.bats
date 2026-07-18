#!/usr/bin/env bats

FIXTURE_BASE="$BATS_TEST_DIRNAME/../../../tests/fixtures/patch-test"

setup() {
  if ! command -v kfg >/dev/null 2>&1; then
    skip "kfg not installed"
  fi
}

@test "patch fixture builds successfully" {
  run kfg build "$FIXTURE_BASE"
  [ "$status" -eq 0 ]
}

@test "patch fixture applies model override" {
  run kfg build "$FIXTURE_BASE"
  [ "$status" -eq 0 ]
  echo "$output" | grep -q "overridden-model"
}

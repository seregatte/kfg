#!/usr/bin/env bats

DISCOVERY_FIXTURE="$BATS_TEST_DIRNAME/../../../tests/fixtures/discovery"

@test "discovery fixture capability has kustomization" {
  [ -f "$DISCOVERY_FIXTURE/capabilities/new-tool/kustomization.yaml" ]
}

@test "discovery fixture capability has README" {
  [ -f "$DISCOVERY_FIXTURE/capabilities/new-tool/README.md" ]
}

@test "discovery fixture capability README has required sections" {
  readme="$DISCOVERY_FIXTURE/capabilities/new-tool/README.md"
  grep -qi "Purpose" "$readme"
  grep -qi "Entrypoints" "$readme"
  grep -qi "Exported Resources" "$readme"
}

@test "discovery fixture capability builds" {
  if ! command -v kfg >/dev/null 2>&1; then
    skip "kfg not installed"
  fi
  run kfg build "$DISCOVERY_FIXTURE/capabilities/new-tool"
  [ "$status" -eq 0 ]
}

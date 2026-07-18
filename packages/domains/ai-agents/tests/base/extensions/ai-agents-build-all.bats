#!/usr/bin/env bats

MANIFESTS_BASE="$BATS_TEST_DIRNAME/../../../manifests"
DOMAIN_ROOT="$BATS_TEST_DIRNAME/../../../"
OVERLAYS_AI="$BATS_TEST_DIRNAME/../../../overlays/ai"
OVERLAYS_DEV="$BATS_TEST_DIRNAME/../../../overlays/dev"

setup() {
  if ! command -v kfg >/dev/null 2>&1; then
    skip "kfg not installed"
  fi
}

BUILD_CMD="kfg build"

# Category aggregates
@test "agents aggregate builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/agents"
  [ "$status" -eq 0 ]
}

@test "cmds aggregate builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/cmds"
  [ "$status" -eq 0 ]
}

@test "steps aggregate builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/steps"
  [ "$status" -eq 0 ]
}

@test "prompts aggregate builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/prompts"
  [ "$status" -eq 0 ]
}

@test "subagents aggregate builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/subagents"
  [ "$status" -eq 0 ]
}

@test "converters aggregate builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/converters"
  [ "$status" -eq 0 ]
}

# Selective entrypoints
@test "opencode agent selective builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/agents/opencode"
  [ "$status" -eq 0 ]
}

@test "pi agent selective builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/agents/pi"
  [ "$status" -eq 0 ]
}

@test "chrome-devtools selective builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/chrome-devtools"
  [ "$status" -eq 0 ]
}

@test "ctx7 selective builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/ctx7"
  [ "$status" -eq 0 ]
}

@test "playwright selective builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/playwright"
  [ "$status" -eq 0 ]
}

@test "notebooklm selective builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/notebooklm"
  [ "$status" -eq 0 ]
}

@test "gws selective builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/gws"
  [ "$status" -eq 0 ]
}

@test "openspec selective builds" {
  run $BUILD_CMD "$MANIFESTS_BASE/openspec"
  [ "$status" -eq 0 ]
}

# Domain aggregate
@test "domain manifests aggregate builds" {
  run $BUILD_CMD "$MANIFESTS_BASE"
  [ "$status" -eq 0 ]
}

# Overlay entrypoints
@test "ai overlay builds" {
  run $BUILD_CMD "$OVERLAYS_AI"
  [ "$status" -eq 0 ]
}

@test "dev overlay builds" {
  run $BUILD_CMD "$OVERLAYS_DEV"
  [ "$status" -eq 0 ]
}

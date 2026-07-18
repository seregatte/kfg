#!/usr/bin/env bats

MANIFESTS_BASE="$BATS_TEST_DIRNAME/../../../manifests"

@test "openspec cmd forwards arguments directly" {
  cmd_file="$MANIFESTS_BASE/cmds/openspec/openspec.yaml"
  [ -f "$cmd_file" ]
  # Verify the command does not contain OPENSPEC_ROOT_DIR logic
  ! grep -q "OPENSPEC_ROOT_DIR" "$cmd_file"
  # Verify it runs openspec with arguments
  grep -q 'command openspec "\$@"' "$cmd_file"
}

@test "openspec cmd sets AGENT env" {
  cmd_file="$MANIFESTS_BASE/cmds/openspec/openspec.yaml"
  grep -q 'AGENT: "openspec"' "$cmd_file"
}

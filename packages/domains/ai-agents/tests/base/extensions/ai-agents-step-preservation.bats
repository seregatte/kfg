#!/usr/bin/env bats

MANIFESTS_BASE="$BATS_TEST_DIRNAME/../../../manifests"

@test "openspec install step does not delete project openspec directory" {
  step_file="$MANIFESTS_BASE/openspec/steps/install.yaml"
  [ -f "$step_file" ]
  # The step should NOT contain rm -rf openspec (which would delete the directory)
  ! grep -q 'rm -rf openspec' "$step_file"
  # It should only remove the temporary binary: rm -f openspec
  grep -q 'rm -f openspec' "$step_file"
}

@test "playwright install step does not delete all of .opencode" {
  step_file="$MANIFESTS_BASE/playwright/steps/install.yaml"
  [ -f "$step_file" ]
  # Should not rm -rf ".opencode" (only the specific subdirectory)
  ! grep -q 'rm -rf "\.opencode"' "$step_file"
  # Should only remove the playwright-cli subdirectory
  grep -q 'rm -rf "$TEMP_SKILL_DIR"' "$step_file"
}

@test "ctx7 install step cleans up only .agents directory" {
  step_file="$MANIFESTS_BASE/ctx7/steps/install.yaml"
  [ -f "$step_file" ]
  # Should remove .agents (temp directory created by ctx7 setup)
  grep -q 'rm -rf .agents' "$step_file"
}

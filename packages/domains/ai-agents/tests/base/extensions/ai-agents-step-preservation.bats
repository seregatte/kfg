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

# --- Regression tests for openspec/ root cleanup (prevent-openspec-root-cleanup) ---

@test "cleanup step does not unconditionally delete openspec directory" {
  # The framework cleanup step should use rm -rf on artifacts, but the engine
  # must never register openspec/ as an artifact in the first place.
  # This test verifies that the cleanup step's run code iterates KFG_ARTIFACTS
  # and does not contain a hardcoded rm -rf for openspec.
  step_file="$BATS_TEST_DIRNAME/../../../../../framework/manifests/steps/cleanup.yaml"
  [ -f "$step_file" ]
  ! grep -q 'rm -rf openspec' "$step_file"
}

@test "cleanup step only removes artifacts listed in KFG_ARTIFACTS" {
  # Verify the cleanup step iterates KFG_ARTIFACTS and removes each entry.
  # This ensures that if openspec/ is never registered as an artifact, it will
  # never be deleted by cleanup.
  step_file="$BATS_TEST_DIRNAME/../../../../../framework/manifests/steps/cleanup.yaml"
  [ -f "$step_file" ]
  grep -q 'for artifact in' "$step_file"
  grep -q 'rm -rf "$artifact"' "$step_file"
}

@test "openspec install step does not register openspec directory as artifact" {
  # The openspec install step should not call __kfg_add_artifact with the
  # openspec/ directory path. It should only register file-level artifacts.
  step_file="$MANIFESTS_BASE/openspec/steps/install.yaml"
  [ -f "$step_file" ]
  ! grep -q '__kfg_add_artifact.*openspec/' "$step_file"
}

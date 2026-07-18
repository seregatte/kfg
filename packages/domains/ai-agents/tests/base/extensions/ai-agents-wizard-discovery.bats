#!/usr/bin/env bats

WIZARD_PROMPT="$BATS_TEST_DIRNAME/../../../overlays/ai/assets/prompts/wizard.yaml"

@test "wizard prompt exists" {
  [ -f "$WIZARD_PROMPT" ]
}

@test "wizard prompt does not contain hardcoded resource catalog" {
  # The wizard should NOT contain hardcoded lists of resource names
  # It should reference documentation discovery instead
  ! grep -q "ai.steps.detect.*ctx7.steps.install.*openspec.steps.install" "$WIZARD_PROMPT"
}

@test "wizard prompt references documentation discovery" {
  grep -qi "read.*README\|capability.*index\|documentation\|discover" "$WIZARD_PROMPT"
}

@test "wizard prompt references Kustomization verification" {
  grep -qi "kustomiz.*build\|verify.*kustomiz\|adjacent.*kustomiz" "$WIZARD_PROMPT"
}

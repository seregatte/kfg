#!/usr/bin/env bats

WIZARD_PROMPT="$BATS_TEST_DIRNAME/../../../overlays/ai/assets/prompts/wizard.yaml"

@test "wizard prompt exists" {
  [ -f "$WIZARD_PROMPT" ]
}

@test "wizard prompt contains resource catalog for building blocks" {
  # The wizard should contain a catalog of available building blocks
  # This is used for the Catalog Check (Mental Inventory) phase
  grep -q "ai.steps.detect" "$WIZARD_PROMPT"
  grep -q "ctx7.steps.install" "$WIZARD_PROMPT"
}

@test "wizard prompt references documentation discovery" {
  grep -qi "read.*README\|capability.*index\|documentation\|discover" "$WIZARD_PROMPT"
}

@test "wizard prompt references kfg build verification" {
  grep -qi "kfg build\|kfg apply" "$WIZARD_PROMPT"
}

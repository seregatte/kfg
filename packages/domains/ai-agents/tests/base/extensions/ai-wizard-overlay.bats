#!/usr/bin/env bats

# AI Wizard Overlay Tests
# Tests for the AI wizard overlay structure and build output

MANIFESTS_OVERLAY="$BATS_TEST_DIRNAME/../../../overlays/ai"

setup() {
    if ! command -v kfg >/dev/null 2>&1; then
        skip "kfg not installed"
    fi
}

@test "kfg build overlays/ai succeeds" {
    run kfg build "$MANIFESTS_OVERLAY"
    [ "$status" -eq 0 ]
}

@test "AI overlay build output contains kfg.ai.workflow" {
    run kfg build "$MANIFESTS_OVERLAY"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "kfg.ai.workflow"
}

@test "AI overlay build output contains wizard prompt asset" {
    run kfg build "$MANIFESTS_OVERLAY"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "ai.prompts.wizard"
}

@test "AI overlay build output contains pi and opencode commands" {
    run kfg build "$MANIFESTS_OVERLAY"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "ai.pi.cmd.main"
    echo "$output" | grep -q "ai.opencode.cmd.main"
}

@test "AI overlay build output contains wizard-command prompt asset" {
    run kfg build "$MANIFESTS_OVERLAY"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "ai.prompts.wizard-command"
}

@test "AI overlay build output contains wizard skill steps for opencode" {
    run kfg build "$MANIFESTS_OVERLAY"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "ai.wizard.skill.opencode"
    echo "$output" | grep -q "weight: -44"
    echo "$output" | grep -q ".opencode/skills/wizard/SKILL.md"
}

@test "AI overlay build output contains wizard skill steps for pi" {
    run kfg build "$MANIFESTS_OVERLAY"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "ai.wizard.skill.pi"
    echo "$output" | grep -q "weight: -44"
    echo "$output" | grep -q ".pi/skills/wizard/SKILL.md"
}

@test "AI overlay build output references wizard-command in workflow steps" {
    run kfg build "$MANIFESTS_OVERLAY"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "ASSETS: ai.prompts.wizard-command"
    echo "$output" | grep -q ".opencode/commands/wizard.md"
    echo "$output" | grep -q ".pi/prompts/wizard.md"
}

@test "AI overlay directory structure is correct" {
    [ -f "$MANIFESTS_OVERLAY/kustomization.yaml" ]
    [ -f "$MANIFESTS_OVERLAY/ai-workflow.yaml" ]
    [ -f "$MANIFESTS_OVERLAY/assets/prompts/wizard.yaml" ]
    [ -f "$MANIFESTS_OVERLAY/assets/prompts/wizard-command.yaml" ]
}

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
    echo "$output" | grep -q "ai.prompts.kfg-wizard"
}

@test "AI overlay build output contains pi and opencode commands" {
    run kfg build "$MANIFESTS_OVERLAY"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "ai.pi.cmd.main"
    echo "$output" | grep -q "ai.opencode.cmd.main"
}

@test "AI overlay directory structure is correct" {
    [ -f "$MANIFESTS_OVERLAY/kustomization.yaml" ]
    [ -f "$MANIFESTS_OVERLAY/ai-workflow.yaml" ]
    [ -f "$MANIFESTS_OVERLAY/assets/prompts/kfg-wizard.yaml" ]
}

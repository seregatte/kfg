#!/usr/bin/env bats

# Validate kfg build succeeds for base and overlay.
# Skip if kfg is not available.

MANIFESTS_BASE="../manifests"
MANIFESTS_OVERLAY="../overlays/dev"

setup() {
    if ! command -v kfg >/dev/null 2>&1; then
        skip "kfg not installed"
    fi
}

@test "kfg build base succeeds" {
    run kfg build "$MANIFESTS_BASE"
    [ "$status" -eq 0 ]
}

@test "kfg build overlay/dev succeeds" {
    run kfg build "$MANIFESTS_OVERLAY"
    [ "$status" -eq 0 ]
}

@test "base build output contains ai-agents resources" {
    run kfg build "$MANIFESTS_BASE"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "ai.claude.asset.settings"
    echo "$output" | grep -q "ai.steps.detect"
    echo "$output" | grep -q "ai.conv.to-json"
}

@test "base build output contains renamed extension resources" {
    run kfg build "$MANIFESTS_BASE"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "ctx7.steps.install"
    echo "$output" | grep -q "ctx7.assets.mcp"
    echo "$output" | grep -q "openspec.steps.install"
    echo "$output" | grep -q "chrome.assets.mcp"
    echo "$output" | grep -q "playwright.assets.mcp"
}

@test "base build output does not contain old names" {
    run kfg build "$MANIFESTS_BASE"
    [ "$status" -eq 0 ]
    ! echo "$output" | grep -q "kfg.agent"
    ! echo "$output" | grep -q "kfg.extension.self"
    ! echo "$output" | grep -q "kfg.detect-agent"
    ! echo "$output" | grep -q "kfg.inject-ctx7-context"
}

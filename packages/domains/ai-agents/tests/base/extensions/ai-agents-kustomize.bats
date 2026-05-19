#!/usr/bin/env bats

# Validate kustomize build succeeds for base and overlay.
# Skip if kustomize is not installed.

MANIFESTS_BASE=".manifests/base"
MANIFESTS_OVERLAY=".manifests/overlay/dev"

setup() {
    if ! command -v kustomize >/dev/null 2>&1; then
        skip "kustomize not installed"
    fi
}

@test "kustomize build base succeeds" {
    run kustomize build "$MANIFESTS_BASE"
    [ "$status" -eq 0 ]
}

@test "kustomize build overlay/dev succeeds" {
    run kustomize build "$MANIFESTS_OVERLAY"
    [ "$status" -eq 0 ]
}

@test "base build output contains ai-agents resources" {
    run kustomize build "$MANIFESTS_BASE"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "ai.claude.asset.settings"
    echo "$output" | grep -q "ai.steps.detect"
    echo "$output" | grep -q "ai.conv.to-json"
}

@test "base build output contains renamed extension resources" {
    run kustomize build "$MANIFESTS_BASE"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "ctx7.steps.install"
    echo "$output" | grep -q "ctx7.assets.mcp"
    echo "$output" | grep -q "openspec.steps.install"
    echo "$output" | grep -q "chrome.assets.mcp"
    echo "$output" | grep -q "playwright.assets.mcp"
}

@test "base build output does not contain old names" {
    run kustomize build "$MANIFESTS_BASE"
    [ "$status" -eq 0 ]
    ! echo "$output" | grep -q "kfg.agent"
    ! echo "$output" | grep -q "kfg.extension.self"
    ! echo "$output" | grep -q "kfg.detect-agent"
    ! echo "$output" | grep -q "kfg.inject-ctx7-context"
}

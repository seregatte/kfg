#!/usr/bin/env bats

# Deprecation Warning Step Tests
# Tests for the deprecation-warn.yaml step

MANIFESTS_STEPS="$BATS_TEST_DIRNAME/../../../manifests/steps"

@test "deprecation-warn.yaml exists in steps directory" {
    [ -f "$MANIFESTS_STEPS/deprecation-warn.yaml" ]
}

@test "deprecation-warn.yaml is valid YAML" {
    run grep -q "apiVersion: kfg.dev/v1alpha1" "$MANIFESTS_STEPS/deprecation-warn.yaml"
    [ "$status" -eq 0 ]
}

@test "deprecation-warn.yaml contains ai.steps.deprecation-warn metadata name" {
    run grep -q "name: ai.steps.deprecation-warn" "$MANIFESTS_STEPS/deprecation-warn.yaml"
    [ "$status" -eq 0 ]
}

@test "deprecation-warn.yaml is registered in kustomization.yaml" {
    run grep -q "deprecation-warn.yaml" "$MANIFESTS_STEPS/kustomization.yaml"
    [ "$status" -eq 0 ]
}

@test "deprecation-warn step is included in dev overlay workflow" {
    OVERLAY_WORKFLOW="$BATS_TEST_DIRNAME/../../../overlays/dev/agents-workflow.yaml"
    run grep -q "ai.steps.deprecation-warn" "$OVERLAY_WORKFLOW"
    [ "$status" -eq 0 ]
}

@test "deprecation-warn step has weight -100" {
    OVERLAY_WORKFLOW="$BATS_TEST_DIRNAME/../../../overlays/dev/agents-workflow.yaml"
    run grep -A 3 "ai.steps.deprecation-warn" "$OVERLAY_WORKFLOW"
    [ "$status" -eq 0 ]
    echo "$output" | grep -q "weight: -100"
}

#!/usr/bin/env bats

# kfg run command terminology tests
# Validates that the run command uses generic "cmd/command" terminology
# instead of AI-specific "agent" terminology

load '../test_helper'

@test "kfg run --help shows 'cmd' in usage string (not 'agent')" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # Check that usage shows [cmd] not [agent]
    [[ "$output" =~ "run [cmd]" ]]
    # Verify it does NOT contain agent terminology
    [[ ! "$output" =~ "run [agent]" ]]
}

@test "kfg run --help shows 'Run a command' in short description (not 'Run an agent')" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # Check for command terminology in short description
    [[ "$output" =~ "Run a command" ]]
    # Verify it does NOT contain agent terminology
    [[ ! "$output" =~ "Run an agent" ]]
}

@test "kfg run --help shows 'Available commands' (not 'Available agents')" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # Check examples section uses commands terminology
    [[ "$output" =~ "lists available commands" ]]
    # Verify it does NOT contain agent terminology
    [[ ! "$output" =~ "lists available agents" ]]
}

@test "kfg run --help does not contain 'agent' anywhere" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # The help output should not contain any agent references
    [[ ! "$output" =~ "agent" ]]
}

@test "kfg run with no args shows 'Available commands:' (not 'Available agents:')" {
    # Create empty manifests directory with proper kustomization.yaml
    mkdir -p "${TEST_TMPDIR}/empty_manifests"
    cat > "${TEST_TMPDIR}/empty_manifests/kustomization.yaml" <<EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources: []
EOF
    
    run "${KFG_BIN}" run -k "${TEST_TMPDIR}/empty_manifests"
    [ "$status" -eq 1 ]
    # Check for command terminology - when empty, shows "No commands found"
    [[ "$output" =~ "No commands found" ]]
}

@test "kfg run with no args shows 'No commands found' (not 'No agents found')" {
    # Create empty manifests directory with proper kustomization.yaml
    mkdir -p "${TEST_TMPDIR}/empty_manifests"
    cat > "${TEST_TMPDIR}/empty_manifests/kustomization.yaml" <<EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources: []
EOF
    
    run "${KFG_BIN}" run -k "${TEST_TMPDIR}/empty_manifests"
    [ "$status" -eq 1 ]
    # Check for command terminology
    [[ "$output" =~ "No commands found" ]]
    # Verify it does NOT contain agent terminology
    [[ ! "$output" =~ "No agents found" ]]
}

@test "kfg run error message uses 'command' when cmd not found" {
    # Create minimal manifests with proper kustomization.yaml
    mkdir -p "${TEST_TMPDIR}/minimal_manifests"
    cat > "${TEST_TMPDIR}/minimal_manifests/kustomization.yaml" <<EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources: []
EOF
    
    run "${KFG_BIN}" run -k "${TEST_TMPDIR}/minimal_manifests" nonexistent-cmd
    [ "$status" -eq 1 ]
    # Error should use command terminology
    [[ "$output" =~ "command 'nonexistent-cmd' not found" ]]
    # Verify it does NOT use agent terminology
    [[ ! "$output" =~ "agent 'nonexistent-cmd' not found" ]]
}

@test "kfg run error message uses 'command' when workflow specified but cmd not found" {
    # Create minimal manifests with proper kustomization.yaml
    mkdir -p "${TEST_TMPDIR}/minimal_manifests"
    cat > "${TEST_TMPDIR}/minimal_manifests/kustomization.yaml" <<EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources: []
EOF
    
    run "${KFG_BIN}" run -k "${TEST_TMPDIR}/minimal_manifests" -w test-workflow nonexistent-cmd
    [ "$status" -eq 1 ]
    # Error should use command terminology
    [[ "$output" =~ "command" ]]
    # Verify it does NOT use agent terminology
    [[ ! "$output" =~ "agent" ]]
}
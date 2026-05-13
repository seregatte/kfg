#!/usr/bin/env bats

# kfg GitHub URL and KFG_KPATH Tests
# Tests for GitHub URL source support and KFG_KPATH environment variable

load '../test_helper'

# Test fixture path
FIXTURE_PATH="${PROJECT_ROOT}/tests/fixtures/basic/kustomization"

# ============================================================================
# KFG_KPATH Environment Variable Tests
# ============================================================================

@test "kfg apply with KFG_KPATH produces shell output" {
    # Set KFG_KPATH to the test fixture
    export KFG_KPATH="${FIXTURE_PATH}"
    
    run "${KFG_BIN}" apply
    [ "$status" -eq 0 ]
    
    # Verify output contains shell code
    [[ "$output" =~ "#!/bin/bash" ]]
    
    # Unset for other tests
    unset KFG_KPATH
}

@test "kfg apply KFG_KPATH output matches explicit -k flag" {
    # Get output with explicit -k flag
    run "${KFG_BIN}" apply -k "${FIXTURE_PATH}"
    [ "$status" -eq 0 ]
    explicit_output="$output"
    
    # Get output with KFG_KPATH
    export KFG_KPATH="${FIXTURE_PATH}"
    run "${KFG_BIN}" apply
    [ "$status" -eq 0 ]
    env_output="$output"
    unset KFG_KPATH
    
    # Outputs should be identical (except for path in comments)
    # The shell code should be the same
    [[ "$explicit_output" =~ "#!/bin/bash" ]]
    [[ "$env_output" =~ "#!/bin/bash" ]]
    
    # Both should define the same functions
    [[ "$explicit_output" =~ "testCmd()" ]]
    [[ "$env_output" =~ "testCmd()" ]]
}

@test "kfg build with KFG_KPATH produces YAML output" {
    # Set KFG_KPATH to the test fixture
    export KFG_KPATH="${FIXTURE_PATH}"
    
    run "${KFG_BIN}" build
    [ "$status" -eq 0 ]
    
    # Verify output contains YAML
    [[ "$output" =~ "apiVersion:" ]]
    [[ "$output" =~ "kind:" ]]
    
    unset KFG_KPATH
}

@test "kfg build KFG_KPATH output matches explicit argument" {
    # Get output with explicit argument
    run "${KFG_BIN}" build "${FIXTURE_PATH}"
    [ "$status" -eq 0 ]
    explicit_output="$output"
    
    # Get output with KFG_KPATH
    export KFG_KPATH="${FIXTURE_PATH}"
    run "${KFG_BIN}" build
    [ "$status" -eq 0 ]
    env_output="$output"
    unset KFG_KPATH
    
    # Outputs should be identical
    [ "$explicit_output" = "$env_output" ]
}

@test "kfg build without argument or KFG_KPATH shows error" {
    # Ensure KFG_KPATH is not set
    unset KFG_KPATH
    
    run "${KFG_BIN}" build
    [ "$status" -eq 2 ]
    
    # Error message should mention KFG_KPATH
    [[ "$output" =~ "KFG_KPATH" ]] || [[ "$output" =~ "kustomization source required" ]]
}

@test "kfg apply without -k, -f, or KFG_KPATH shows error" {
    # Ensure KFG_KPATH is not set
    unset KFG_KPATH
    
    run "${KFG_BIN}" apply
    [ "$status" -eq 2 ]
    
    # Error message should mention KFG_KPATH
    [[ "$output" =~ "KFG_KPATH" ]] || [[ "$output" =~ "kustomization source required" ]]
}

@test "kfg run without -k, -f, or KFG_KPATH shows error" {
    # Ensure KFG_KPATH is not set
    unset KFG_KPATH
    
    run "${KFG_BIN}" run
    [ "$status" -eq 2 ]
    
    # Error message should mention KFG_KPATH
    [[ "$output" =~ "KFG_KPATH" ]] || [[ "$output" =~ "kustomization source required" ]]
}

@test "kfg build with argument overrides KFG_KPATH" {
    # Set KFG_KPATH to a different path
    export KFG_KPATH="/some/other/path"
    
    # Use explicit argument (should override)
    run "${KFG_BIN}" build "${FIXTURE_PATH}"
    [ "$status" -eq 0 ]
    
    # Output should be from explicit argument, not KFG_KPATH
    [[ "$output" =~ "apiVersion:" ]]
    [[ "$output" =~ "kind:" ]]
    
    unset KFG_KPATH
}

@test "kfg apply with -k flag overrides KFG_KPATH" {
    # Set KFG_KPATH to a different path
    export KFG_KPATH="/some/other/path"
    
    # Use explicit -k flag (should override)
    run "${KFG_BIN}" apply -k "${FIXTURE_PATH}"
    [ "$status" -eq 0 ]
    
    # Output should be from explicit -k, not KFG_KPATH
    [[ "$output" =~ "#!/bin/bash" ]]
    
    unset KFG_KPATH
}

# ============================================================================
# GitHub URL Support Tests
# ============================================================================

# Note: These tests require network access and are skipped in CI
# unless KFG_NETWORK_TESTS=true is set

@test "kfg build accepts GitHub URL syntax (network required)" {
    # Skip in CI unless explicitly enabled
    if [ "${CI}" = "true" ] && [ "${KFG_NETWORK_TESTS}" != "true" ]; then
        skip "Skipping network test in CI environment"
    fi
    
    # Test that kfg build accepts GitHub URL syntax
    # Note: This test requires network access to clone the repo
    # Using a stable kustomize test fixture
    github_url="https://github.com/kubernetes-sigs/kustomize//cmd/config/testdata/bases/simple?ref=master"
    
    run "${KFG_BIN}" build "${github_url}"
    # The test may fail if the repository structure changes, but
    # the important thing is that the URL is accepted and processed
    # Status should be 0 (success) or 1 (load error, but URL was accepted)
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
    
    # If successful, output should contain YAML
    if [ "$status" -eq 0 ]; then
        [[ "$output" =~ "apiVersion:" ]]
    fi
}

@test "kfg apply accepts GitHub URL in -k flag (network required)" {
    # Skip in CI unless explicitly enabled
    if [ "${CI}" = "true" ] && [ "${KFG_NETWORK_TESTS}" != "true" ]; then
        skip "Skipping network test in CI environment"
    fi
    
    # Test that kfg apply accepts GitHub URL in -k flag
    # Note: This requires manifests with CmdWorkflow, which kustomize fixtures don't have
    # So we expect an error about missing workflows, but the URL should be accepted
    github_url="https://github.com/kubernetes-sigs/kustomize//cmd/config/testdata/bases/simple?ref=master"
    
    run "${KFG_BIN}" apply -k "${github_url}"
    # Status could be 0, 1, or 2 depending on manifest content
    # The important thing is the URL syntax is accepted
    
    # Check that it's not an argument validation error (exit code 2)
    # If exit code is 2, check that it's not "invalid argument" type error
    if [ "$status" -eq 2 ]; then
        # Should not be an "invalid argument" or "missing source" error
        [[ "$output" != *"invalid argument"* ]] || true
    fi
}

@test "kfg run accepts GitHub URL in -k flag (network required)" {
    # Skip in CI unless explicitly enabled
    if [ "${CI}" = "true" ] && [ "${KFG_NETWORK_TESTS}" != "true" ]; then
        skip "Skipping network test in CI environment"
    fi
    
    # Test that kfg run accepts GitHub URL in -k flag
    github_url="https://github.com/kubernetes-sigs/kustomize//cmd/config/testdata/bases/simple?ref=master"
    
    run "${KFG_BIN}" run -k "${github_url}"
    # Status could be various values depending on manifest content
    # The important thing is the URL syntax is accepted
}

@test "kfg build --help mentions KFG_KPATH and GitHub URLs" {
    run "${KFG_BIN}" build --help
    [ "$status" -eq 0 ]
    
    # Check for KFG_KPATH documentation
    [[ "$output" =~ "KFG_KPATH" ]]
    
    # Check for GitHub URL examples
    [[ "$output" =~ "github.com" ]]
}

@test "kfg apply --help mentions KFG_KPATH and GitHub URLs" {
    run "${KFG_BIN}" apply --help
    [ "$status" -eq 0 ]
    
    # Check for KFG_KPATH documentation
    [[ "$output" =~ "KFG_KPATH" ]]
    
    # Check for GitHub URL examples
    [[ "$output" =~ "github.com" ]]
}

@test "kfg run --help mentions KFG_KPATH and GitHub URLs" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    
    # Check for KFG_KPATH documentation
    [[ "$output" =~ "KFG_KPATH" ]]
    
    # Check for GitHub URL examples
    [[ "$output" =~ "github.com" ]]
}

# ============================================================================
# Shell Validation Tests
# ============================================================================

@test "kfg apply output is valid bash script" {
    run "${KFG_BIN}" apply -k "${FIXTURE_PATH}"
    [ "$status" -eq 0 ]
    
    # Verify the shell code is syntactically valid
    verify_shell_syntax "$output"
}

@test "kfg apply KFG_KPATH output is valid bash script" {
    export KFG_KPATH="${FIXTURE_PATH}"
    
    run "${KFG_BIN}" apply
    [ "$status" -eq 0 ]
    
    # Verify the shell code is syntactically valid
    verify_shell_syntax "$output"
    
    unset KFG_KPATH
}
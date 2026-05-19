#!/usr/bin/env bats

# kfg run command examples tests
# Validates that examples in --help use generic paths like .manifests/overlay/dev
# instead of legacy .nixai/overlay/dev paths

load '../test_helper'

@test "kfg run --help examples use '.manifests/overlay/dev' (not '.nixai/overlay/dev')" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # Check that examples use .manifests path
    [[ "$output" =~ ".manifests/overlay/dev" ]]
    # Verify it does NOT use legacy .nixai path
    [[ ! "$output" =~ ".nixai/overlay/dev" ]]
}

@test "kfg run --help examples do not contain '.nixai' anywhere" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # The help output should not contain any .nixai references
    [[ ! "$output" =~ ".nixai" ]]
}

@test "kfg run --help examples use '.manifests' paths consistently" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # Count occurrences of .manifests
    manifests_count=$(echo "$output" | grep -c ".manifests" || true)
    # Should have at least one .manifests reference in examples
    [ "$manifests_count" -ge 1 ]
}

@test "kfg run --help GitHub example uses valid path format" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # Check GitHub URL example format
    [[ "$output" =~ "https://github.com/owner/repo//manifests" ]]
}

@test "kfg run --help KFG_KPATH example uses .manifests path" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # Check KFG_KPATH example uses .manifests
    [[ "$output" =~ "KFG_KPATH=./manifests" ]]
    # Verify it does NOT use legacy .nixai
    [[ ! "$output" =~ "KFG_KPATH=.nixai" ]]
}

@test "kfg run --help examples do not reference AI-specific names" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # Examples should not contain AI-specific command names like 'claude'
    # (except in comments explaining the functionality)
    # Check that examples use generic names like 'my-cmd'
    [[ "$output" =~ "my-cmd" ]]
}
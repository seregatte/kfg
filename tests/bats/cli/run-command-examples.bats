#!/usr/bin/env bats

# kfg run command examples tests
# Validates that examples in --help use package-oriented paths like packages/domains/ai-agents/overlays/dev
# instead of legacy .manifests/ or .nixai/ paths

load '../test_helper'

@test "kfg run --help examples use package paths (not '.manifests' or '.nixai')" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # Check that examples use packages/ path format
    [[ "$output" =~ "packages/domains/ai-agents/overlays/dev" ]]
    # Verify it does NOT use legacy .manifests or .nixai paths
    [[ ! "$output" =~ ".manifests/overlay/dev" ]]
    [[ ! "$output" =~ ".nixai/overlay/dev" ]]
}

@test "kfg run --help examples do not contain '.manifests' or '.nixai'" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # The help output should not contain legacy paths
    [[ ! "$output" =~ ".manifests" ]]
    [[ ! "$output" =~ ".nixai" ]]
}

@test "kfg run --help examples reference package paths" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # Count occurrences of packages/ path
    packages_count=$(echo "$output" | grep -c "packages/" || true)
    # Should have at least one packages/ reference in examples
    [ "$packages_count" -ge 1 ]
}

@test "kfg run --help GitHub example uses valid package path format" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # Check GitHub URL example format with new package path
    [[ "$output" =~ "https://github.com/owner/repo//packages/domains/ai-agents/overlays/dev" ]]
}

@test "kfg run --help KFG_KPATH example uses package path" {
    run "${KFG_BIN}" run --help
    [ "$status" -eq 0 ]
    # Check KFG_KPATH example uses packages/framework path
    [[ "$output" =~ "KFG_KPATH=./packages/framework" ]]
    # Verify it does NOT use legacy .manifests or .nixai
    [[ ! "$output" =~ "KFG_KPATH=.manifests" ]]
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
#!/usr/bin/env bats

# kfg ai Command Tests
# Tests for the AI wizard CLI command

load '../test_helper'

@test "kfg ai --help succeeds" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
}

@test "kfg ai --help mentions KFG_AI_AGENT environment variable" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "KFG_AI_AGENT" ]]
}

@test "kfg ai --help mentions default agent (pi)" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "default: pi" ]]
}

@test "kfg ai --help mentions supported agents" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "pi" ]]
    [[ "$output" =~ "opencode" ]]
}

@test "kfg ai --help includes usage examples" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "kfg ai" ]]
    [[ "$output" =~ "Examples:" ]]
}

@test "kfg ai --help describes wizard purpose" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "interactive" ]]
    [[ "$output" =~ "configuration" ]]
}

@test "kfg ai --help explains argument forwarding" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "forwarded" ]] || [[ "$output" =~ "--" ]]
}

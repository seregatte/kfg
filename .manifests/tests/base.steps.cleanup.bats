#!/usr/bin/env bats

load 'test_helper'

@test "removes artifacts from KFG_ARTIFACTS" {
    mkdir -p "$TEST_TMPDIR/.claude"
    mkdir -p "$TEST_TMPDIR/.opencode"
    
    cd "$TEST_TMPDIR"
    
    local code
    code=$(printf 'declare -a KFG_ARTIFACTS=(.claude .opencode)\nWORKSPACE_DIR="%s"\n%s\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test' "$TEST_TMPDIR" "$(mock_kfg_log)" "$(step_run_code cleanup)")
    
    bash -c "$code"
    
    [ ! -d "$TEST_TMPDIR/.claude" ]
    [ ! -d "$TEST_TMPDIR/.opencode" ]
}

@test "does not fail when KFG_ARTIFACTS is empty" {
    cd "$TEST_TMPDIR"
    KFG_ARTIFACTS=() WORKSPACE_DIR="$TEST_TMPDIR" run_step_with_env "cleanup"
    
    [ $? -eq 0 ]
}

@test "does not fail when artifact does not exist" {
    cd "$TEST_TMPDIR"
    KFG_ARTIFACTS=("nonexistent-dir") WORKSPACE_DIR="$TEST_TMPDIR" run_step_with_env "cleanup"
    
    [ $? -eq 0 ]
}
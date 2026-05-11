#!/usr/bin/env bats

load 'test_helper'

@test "creates directories from DIRECTORIES env" {
    run_step "materialize-scaffold" "foo:bar/baz:deep/nested/path"
    
    [ -d "$TEST_TMPDIR/foo" ]
    [ -d "$TEST_TMPDIR/bar/baz" ]
    [ -d "$TEST_TMPDIR/deep/nested/path" ]
}

@test "handles empty DIRECTORIES gracefully" {
    run_step "materialize-scaffold" ""
    
    [ $? -eq 0 ]
}

@test "skips empty segments in colon-separated paths" {
    run_step "materialize-scaffold" "foo::bar"
    
    [ -d "$TEST_TMPDIR/foo" ]
    [ -d "$TEST_TMPDIR/bar" ]
}

@test "creates nested directories for claude agent" {
    run_step "materialize-scaffold" ".claude:.claude/skills:.claude/commands:.claude/subagents"
    
    [ -d "$TEST_TMPDIR/.claude" ]
    [ -d "$TEST_TMPDIR/.claude/skills" ]
    [ -d "$TEST_TMPDIR/.claude/commands" ]
    [ -d "$TEST_TMPDIR/.claude/subagents" ]
}

@test "creates nested directories for opencode agent" {
    run_step "materialize-scaffold" ".opencode:.opencode/skills:.opencode/commands:.opencode/subagents"
    
    [ -d "$TEST_TMPDIR/.opencode" ]
    [ -d "$TEST_TMPDIR/.opencode/skills" ]
    [ -d "$TEST_TMPDIR/.opencode/commands" ]
    [ -d "$TEST_TMPDIR/.opencode/subagents" ]
}

@test "creates nested directories for gemini agent" {
    run_step "materialize-scaffold" ".gemini:.gemini/skills:.gemini/commands"
    
    [ -d "$TEST_TMPDIR/.gemini" ]
    [ -d "$TEST_TMPDIR/.gemini/skills" ]
    [ -d "$TEST_TMPDIR/.gemini/commands" ]
}

@test "creates nested directories for pi agent" {
    run_step "materialize-scaffold" ".pi:.pi/skills:.pi/commands"
    
    [ -d "$TEST_TMPDIR/.pi" ]
    [ -d "$TEST_TMPDIR/.pi/skills" ]
    [ -d "$TEST_TMPDIR/.pi/commands" ]
}

@test "does not fail on missing DIRECTORIES env" {
    cd "$TEST_TMPDIR"
    local code
    code=$(printf '%s\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test' "$(mock_kfg_log)" "$(step_run_code materialize-scaffold)")
    bash -c "$code"
    
    [ $? -eq 0 ]
}

@test "registers artifacts for created directories" {
    local code
    code=$(printf '%s\n%s\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\necho "${KFG_ARTIFACTS[@]}"' "$(mock_kfg_log)" "$(mock_artifact_tracking)" "$(step_run_code materialize-scaffold)")
    
    cd "$TEST_TMPDIR"
    local output
    output=$(DIRECTORIES=".claude:.claude/skills" bash -c "$code")
    
    [ -d "$TEST_TMPDIR/.claude" ]
    [ -d "$TEST_TMPDIR/.claude/skills" ]
    echo "$output" | grep -q ".claude"
    echo "$output" | grep -q ".claude/skills"
}
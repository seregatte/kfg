#!/usr/bin/env bats

load 'test_helper'

@test "copies file with SRC and DEST" {
    mkdir -p "$TEST_TMPDIR/src"
    echo "test content" > "$TEST_TMPDIR/src/test.txt"
    
    cd "$TEST_TMPDIR/src"
    SRC="test.txt" DEST="output.txt" OUTPUT_DIR="$TEST_TMPDIR" run_step_with_env "copy-context"
    
    [ -f "$TEST_TMPDIR/output.txt" ]
    [ "$(cat "$TEST_TMPDIR/output.txt")" = "test content" ]
}

@test "copies file using basename when DEST is empty" {
    mkdir -p "$TEST_TMPDIR/src"
    echo "test content" > "$TEST_TMPDIR/src/test.txt"
    
    cd "$TEST_TMPDIR/src"
    SRC="test.txt" OUTPUT_DIR="$TEST_TMPDIR" run_step_with_env "copy-context"
    
    [ -f "$TEST_TMPDIR/test.txt" ]
}

@test "copies directory recursively" {
    mkdir -p "$TEST_TMPDIR/src/subdir"
    echo "file1" > "$TEST_TMPDIR/src/file1.txt"
    echo "file2" > "$TEST_TMPDIR/src/subdir/file2.txt"
    
    cd "$TEST_TMPDIR"
    SRC="src" DEST="dest" OUTPUT_DIR="$TEST_TMPDIR" run_step_with_env "copy-context"
    
    [ -d "$TEST_TMPDIR/dest" ]
    [ -f "$TEST_TMPDIR/dest/file1.txt" ]
    [ -f "$TEST_TMPDIR/dest/subdir/file2.txt" ]
}

@test "does not fail when SRC does not exist" {
    cd "$TEST_TMPDIR"
    SRC="nonexistent.txt" OUTPUT_DIR="$TEST_TMPDIR" run_step_with_env "copy-context"
    
    [ $? -eq 0 ]
}

@test "does not fail when SRC is empty" {
    cd "$TEST_TMPDIR"
    SRC="" OUTPUT_DIR="$TEST_TMPDIR" run_step_with_env "copy-context"
    
    [ $? -eq 0 ]
}

@test "creates parent directories for DEST" {
    mkdir -p "$TEST_TMPDIR/src"
    echo "test content" > "$TEST_TMPDIR/src/test.txt"
    
    cd "$TEST_TMPDIR/src"
    SRC="test.txt" DEST="deep/nested/path/output.txt" OUTPUT_DIR="$TEST_TMPDIR" run_step_with_env "copy-context"
    
    [ -f "$TEST_TMPDIR/deep/nested/path/output.txt" ]
}

run_step_with_env() {
    local step="$1"
    local code
    code=$(printf '%s\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test' "$(mock_kfg_log)" "$(step_run_code "$step")")
    bash -c "$code"
}

@test "registers artifact for copied file" {
    mkdir -p "$TEST_TMPDIR/src"
    echo "test content" > "$TEST_TMPDIR/src/test.txt"
    
    local code
    code=$(printf '%s\n%s\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\necho "${KFG_ARTIFACTS[@]}"' "$(mock_kfg_log)" "$(mock_artifact_tracking)" "$(step_run_code copy-context)")
    
    cd "$TEST_TMPDIR/src"
    local output
    output=$(SRC="test.txt" DEST="output.txt" OUTPUT_DIR="$TEST_TMPDIR" bash -c "$code")
    
    [ -f "$TEST_TMPDIR/output.txt" ]
    echo "$output" | grep -q "output.txt"
}
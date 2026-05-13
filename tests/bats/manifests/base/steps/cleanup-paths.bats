#!/usr/bin/env bats

load '../../../helpers/common'
load '../../../helpers/manifests'

@test "removes __kfg_* paths from TMPDIR" {
    mkdir -p "$TMPDIR/__kfg_test_dir"
    touch "$TMPDIR/__kfg_test_file"
    
    cd "$TEST_TMPDIR"
    run_step_with_env "cleanup-paths"
    
    [ ! -d "$TMPDIR/__kfg_test_dir" ]
    [ ! -f "$TMPDIR/__kfg_test_file" ]
}

@test "removes __kfg_* paths from /tmp" {
    mkdir -p "/tmp/__kfg_test_dir"
    touch "/tmp/__kfg_test_file"
    
    cd "$TEST_TMPDIR"
    run_step_with_env "cleanup-paths"
    
    [ ! -d "/tmp/__kfg_test_dir" ]
    [ ! -f "/tmp/__kfg_test_file" ]
}

@test "does not fail when no __kfg_* paths exist" {
    cd "$TEST_TMPDIR"
    run_step_with_env "cleanup-paths"
    
    [ $? -eq 0 ]
}
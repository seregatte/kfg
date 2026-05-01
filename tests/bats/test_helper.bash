#!/usr/bin/env bash

# Test helper functions for Bats tests

# Get the project root directory (absolute path)
# BATS_TEST_DIRNAME is the directory containing the test file (tests/bats)
# We need to go up two levels to get to the project root
PROJECT_ROOT="$(cd "${BATS_TEST_DIRNAME}/../.." && pwd)"

# KFG binary path
KFG_BIN="${PROJECT_ROOT}/bin/kfg"

# Setup function - runs before each test
setup() {
    # Nothing to do - binary should already exist
    :
}

# Teardown function - runs after each test
teardown() {
    # Cleanup any temporary files
    rm -f /tmp/test_kfg_*.sh
}

# Helper function to check if a function exists in generated shell code
function_exists_in_output() {
    local output="$1"
    local function_name="$2"
    
    if [[ "$output" =~ "${function_name}()" ]]; then
        return 0
    else
        return 1
    fi
}

# Helper function to count function definitions in output
count_functions_in_output() {
    local output="$1"
    echo "$output" | grep -c "^[a-zA-Z_][a-zA-Z0-9_]*()"
}

# Helper function to verify shell syntax
verify_shell_syntax() {
    local shell_code="$1"
    echo "$shell_code" > /tmp/test_syntax.sh
    bash -n /tmp/test_syntax.sh
    local status=$?
    rm -f /tmp/test_syntax.sh
    return $status
}

# Helper function to measure execution time
measure_time() {
    local command="$1"
    local start=$(date +%s%N)
    eval "$command"
    local end=$(date +%s%N)
    local elapsed=$(( (end - start) / 1000000 ))
    echo "$elapsed"
}
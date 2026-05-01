#!/usr/bin/env bats

# Integration tests for workspace start/stop
# Tests artifact-scoped backup and cleanup behavior

load 'test_helper'

setup() {
    # Create a temporary workspace directory for each test
    WORK_DIR="$(mktemp -d)"
    # Use test number to create unique tags (images are immutable, so can't reuse tags)
    TEST_TAG="workspace-test-${BATS_TEST_NUMBER}-$$"
    # Create isolated store directory for test isolation
    STORE_DIR="${WORK_DIR}/store"
    mkdir -p "$STORE_DIR"
    cd "$WORK_DIR"
}

teardown() {
    # Cleanup workspace
    cd /
    rm -rf "$WORK_DIR"
}

@test "workspace with unrelated files, start creates scoped backup" {
    # Task 6.1: Test scoped backup with unrelated files
    
    # Create image with specific file
    cat > Imagefile <<EOF
FROM scratch
COPY artifact.txt ./
TAG ${TEST_TAG}:v1
EOF
    
    echo "artifact content" > artifact.txt
    
    OUTPUT_DIR="$WORK_DIR/build-output"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    # Create workspace with unrelated files and conflicting file
    WORKSPACE_DIR="$WORK_DIR/workspace"
    mkdir -p "$WORKSPACE_DIR"
    echo "unrelated content" > "$WORKSPACE_DIR/unrelated.txt"
    echo "existing artifact" > "$WORKSPACE_DIR/artifact.txt"
    
    # Start workspace
    run "${KFG_BIN}" --store "$STORE_DIR" workspace start "${TEST_TAG}:v1" --root "$WORKSPACE_DIR" --name "test-instance"
    [ "$status" -eq 0 ]
    
    # Verify only conflicting file was backed up
    BACKUP_DIR="$STORE_DIR/.workspace/test-instance/backup/data"
    [ -f "$BACKUP_DIR/artifact.txt" ]
    [ ! -f "$BACKUP_DIR/unrelated.txt" ]
    
    # Stop workspace
    run "${KFG_BIN}" --store "$STORE_DIR" workspace stop --name "test-instance"
    [ "$status" -eq 0 ]
}

@test "stop removes only image artifacts, unrelated files remain" {
    # Task 6.2: Test scoped cleanup preserving unrelated files
    
    # Create image
    cat > Imagefile <<EOF
FROM scratch
COPY artifact.txt ./
TAG ${TEST_TAG}:v1
EOF
    
    echo "artifact content" > artifact.txt
    
    OUTPUT_DIR="$WORK_DIR/build-output"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    # Create workspace with unrelated files
    WORKSPACE_DIR="$WORK_DIR/workspace"
    mkdir -p "$WORKSPACE_DIR"
    echo "unrelated content" > "$WORKSPACE_DIR/unrelated.txt"
    
    # Start workspace (no conflicts, no backup)
    run "${KFG_BIN}" --store "$STORE_DIR" workspace start "${TEST_TAG}:v1" --root "$WORKSPACE_DIR" --name "test-instance-2"
    [ "$status" -eq 0 ]
    
    # Verify artifact was materialized
    [ -f "$WORKSPACE_DIR/artifact.txt" ]
    
    # Stop workspace
    run "${KFG_BIN}" --store "$STORE_DIR" workspace stop --name "test-instance-2"
    [ "$status" -eq 0 ]
    
    # Verify artifact was removed but unrelated file remains
    [ ! -f "$WORKSPACE_DIR/artifact.txt" ]
    [ -f "$WORKSPACE_DIR/unrelated.txt" ]
}

@test "instance metadata contains materialized paths after start" {
    # Task 6.3: Test materialized paths tracking
    
    # Create image with multiple files
    cat > Imagefile <<EOF
FROM scratch
COPY file1.txt ./
COPY file2.txt ./
TAG ${TEST_TAG}:v1
EOF
    
    echo "content1" > file1.txt
    echo "content2" > file2.txt
    
    OUTPUT_DIR="$WORK_DIR/build-output"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    # Create empty workspace
    WORKSPACE_DIR="$WORK_DIR/workspace"
    mkdir -p "$WORKSPACE_DIR"
    
    # Start workspace
    run "${KFG_BIN}" --store "$STORE_DIR" workspace start "${TEST_TAG}:v1" --root "$WORKSPACE_DIR" --name "test-instance-3"
    [ "$status" -eq 0 ]
    
    # Verify instance metadata exists
    INSTANCE_FILE="$STORE_DIR/.workspace/test-instance-3/instance.json"
    [ -f "$INSTANCE_FILE" ]
    
    # Verify materialized_paths field contains expected paths
    run cat "$INSTANCE_FILE"
    [[ "$output" =~ "materialized_paths" ]]
    
    # Stop workspace
    run "${KFG_BIN}" --store "$STORE_DIR" workspace stop --name "test-instance-3"
    [ "$status" -eq 0 ]
}

@test "empty parent directories removed after cleanup" {
    # Task 6.4: Test empty directory removal
    
    # Create image with nested file
    mkdir -p nested_dir
    cat > Imagefile <<EOF
FROM scratch
COPY nested_dir/file.txt ./nested_dir/
TAG ${TEST_TAG}:v1
EOF
    
    echo "nested content" > nested_dir/file.txt
    
    OUTPUT_DIR="$WORK_DIR/build-output"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    # Create empty workspace
    WORKSPACE_DIR="$WORK_DIR/workspace"
    mkdir -p "$WORKSPACE_DIR"
    
    # Start workspace
    run "${KFG_BIN}" --store "$STORE_DIR" workspace start "${TEST_TAG}:v1" --root "$WORKSPACE_DIR" --name "test-instance-4"
    [ "$status" -eq 0 ]
    
    # Verify nested directory was created
    [ -f "$WORKSPACE_DIR/nested_dir/file.txt" ]
    
    # Stop workspace
    run "${KFG_BIN}" --store "$STORE_DIR" workspace stop --name "test-instance-4"
    [ "$status" -eq 0 ]
    
    # Verify file and empty parent directories are removed
    [ ! -f "$WORKSPACE_DIR/nested_dir/file.txt" ]
    [ ! -d "$WORKSPACE_DIR/nested_dir" ]
}
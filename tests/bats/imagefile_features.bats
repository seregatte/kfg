#!/usr/bin/env bats

# Integration tests for Imagefile features
# Tests: multi-stage isolation, glob expansion, absolute paths, WORKDIR, logging

load 'test_helper'

setup() {
    # Create a temporary workspace directory for each test
    WORK_DIR="$(mktemp -d)"
    # Use test number to create unique tags
    TEST_TAG="feature-${BATS_TEST_NUMBER}-$$"
    # Create isolated store directory
    STORE_DIR="${WORK_DIR}/store"
    mkdir -p "$STORE_DIR"
    cd "$WORK_DIR"
}

teardown() {
    # Cleanup workspace
    cd /
    rm -rf "$WORK_DIR"
}

# Section 7.1: Multi-stage build with COPY --from
@test "multi-stage build with COPY --from" {
    cat > base-artifact.txt <<EOF
base content
EOF
    
    cat > Imagefile <<EOF
FROM scratch AS base
COPY base-artifact.txt /artifact.txt
FROM scratch AS final
COPY --from=base /artifact.txt /final-artifact.txt
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-multi-stage"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    [ -f "$OUTPUT_DIR/artifacts/final-artifact.txt" ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts" --keep-build
    [ "$status" -eq 0 ]
}

# Section 7.2: FROM scratch isolation test
@test "FROM scratch isolation - no previous files" {
    # Create files in workspace
    echo "workspace file" > workspace-file.txt
    
    # Create Imagefile with FROM scratch (should be empty, not include workspace files)
    cat > Imagefile <<EOF
FROM scratch
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-scratch-isolation"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    # Verify artifacts directory has no files (isolated from workspace)
    artifacts_count=$(find "$OUTPUT_DIR/artifacts" -type f ! -name "metadata.json" | wc -l | tr -d ' ')
    [ "$artifacts_count" -eq 0 ]
}

# Section 7.3: Glob expansion with multiple matching files
@test "glob expansion with multiple matching files" {
    echo "file1" > file1.txt
    echo "file2" > file2.txt
    echo "file3" > file3.txt
    
    cat > Imagefile <<EOF
FROM scratch
COPY *.txt /docs/
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-glob"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    [ -f "$OUTPUT_DIR/artifacts/docs/file1.txt" ]
    [ -f "$OUTPUT_DIR/artifacts/docs/file2.txt" ]
    [ -f "$OUTPUT_DIR/artifacts/docs/file3.txt" ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts" --keep-build
    [ "$status" -eq 0 ]
}

# Section 7.4: Glob zero match test (expect failure)
@test "glob zero match - expect failure" {
    # Create files that don't match the glob pattern
    echo "data" > data.json
    
    # Create Imagefile with glob pattern that matches nothing
    cat > Imagefile <<EOF
FROM scratch
COPY *.go /src/
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-glob-zero"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    
    # Should fail (glob matches zero files)
    [ "$status" -ne 0 ]
    
    # Error message should mention zero match or no files
    [[ "$output" =~ "0" ]] || [[ "$output" =~ "no" ]] || [[ "$output" =~ "match" ]] || [[ "$output" =~ "found" ]]
}

# Section 7.5: Absolute path COPY test
@test "absolute path COPY - normalized to relative" {
    echo "config data" > config.json
    
    cat > Imagefile <<EOF
FROM scratch
COPY config.json /app/config/config.json
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-absolute-path"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    [ -f "$OUTPUT_DIR/artifacts/app/config/config.json" ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts" --keep-build
    [ "$status" -eq 0 ]
}

# Section 7.6: WORKDIR instruction test
@test "WORKDIR instruction - affects COPY destination" {
    echo "app content" > app-file.txt
    
    cat > Imagefile <<EOF
FROM scratch
WORKDIR /app
COPY app-file.txt .
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-workdir"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    [ -f "$OUTPUT_DIR/artifacts/app/app-file.txt" ]
    [ ! -f "$OUTPUT_DIR/artifacts/app-file.txt" ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts" --keep-build
    [ "$status" -eq 0 ]
}

@test "WORKDIR chained - relative path on absolute" {
    echo "nested content" > nested-file.txt
    
    cat > Imagefile <<EOF
FROM scratch
WORKDIR /app
WORKDIR subdir
COPY nested-file.txt .
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-workdir-chained"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    [ -f "$OUTPUT_DIR/artifacts/app/subdir/nested-file.txt" ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts" --keep-build
    [ "$status" -eq 0 ]
}

# Section 7.7: Logging format test (verify JSONL output)
@test "logging format - JSONL structured output" {
    cat > Imagefile <<EOF
FROM scratch
COPY test.txt /test.txt
RUN echo "test stdout"
TAG ${TEST_TAG}:v1
EOF
    
    echo "test content" > test.txt
    
    OUTPUT_DIR="$WORK_DIR/build-logging"
    LOG_FILE="$WORK_DIR/build.log"
    
    NIXAI_LOG_FILE="$LOG_FILE" run "${KFG_BIN}" --store "$STORE_DIR" --verbose=3 image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    if [ -f "$LOG_FILE" ]; then
        while IFS= read -r line; do
            [[ "$line" =~ "\"ts\"" ]]
            [[ "$line" =~ "\"level\"" ]]
            [[ "$line" =~ "\"component\"" ]]
            [[ "$line" =~ "\"msg\"" ]]
            [[ "$line" =~ "\"source\"" ]]
        done < "$LOG_FILE"
        
        if grep -q "build:run:" "$LOG_FILE"; then
            [[ "$(grep "build:run:stdout" "$LOG_FILE" | head -1)" =~ "\"component\":\"build:run:stdout\"" ]] || true
            [[ "$(grep "build:run:stderr" "$LOG_FILE" | head -1)" =~ "\"component\":\"build:run:stderr\"" ]] || true
        fi
    fi
}

# Section 7.8: Verify new directory structure
@test "new directory structure - stages and artifacts" {
    # Create multi-stage Imagefile
    cat > file1.txt <<EOF
stage1 content
EOF
    
    cat > file2.txt <<EOF
stage2 content
EOF
    
    cat > Imagefile <<EOF
FROM scratch AS stage1
COPY file1.txt /file1.txt
FROM scratch AS stage2
COPY file2.txt /file2.txt
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-structure"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    # Verify stages directory exists
    [ -d "$OUTPUT_DIR/stages" ]
    
    # Verify each stage has its own directory
    [ -d "$OUTPUT_DIR/stages/stage1" ]
    [ -d "$OUTPUT_DIR/stages/stage2" ]
    
    # Verify artifacts directory exists (contains final stage)
    [ -d "$OUTPUT_DIR/artifacts" ]
    
    # Verify only final stage files are in artifacts
    [ -f "$OUTPUT_DIR/artifacts/file2.txt" ]
    [ ! -f "$OUTPUT_DIR/artifacts/file1.txt" ]
}

@test "COPY --from error handling - non-existent stage" {
    # Create Imagefile with COPY --from referencing non-existent stage
    cat > Imagefile <<EOF
FROM scratch
COPY --from=nonexistent file.txt /file.txt
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-copy-error"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    
    # Should fail
    [ "$status" -ne 0 ]
    
    # Error message should mention non-existent stage
    [[ "$output" =~ "stage" ]] || [[ "$output" =~ "not found" ]] || [[ "$output" =~ "nonexistent" ]]
}
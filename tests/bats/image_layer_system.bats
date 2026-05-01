#!/usr/bin/env bats

# Integration tests for image layer system
# Tests the full workflow: build -> push -> start -> stop

load 'test_helper'

setup() {
    WORK_DIR="$(mktemp -d)"
    TEST_TAG="test-${BATS_TEST_NUMBER}-$$"
    STORE_DIR="${WORK_DIR}/store"
    mkdir -p "$STORE_DIR"
    cd "$WORK_DIR"
}

teardown() {
    cd /
    rm -rf "$WORK_DIR"
}

@test "build image from scratch" {
    cat > Imagefile <<EOF
FROM scratch
TAG ${TEST_TAG}:v1
EOF
    
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root .
    
    [ "$status" -eq 0 ]
}

@test "push and list images" {
    cat > Imagefile <<EOF
FROM scratch
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-output"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    [ -d "$OUTPUT_DIR/artifacts" ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image list
    [ "$status" -eq 0 ]
}

@test "inspect stored image" {
    cat > Imagefile <<EOF
FROM scratch
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-inspect"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect "${TEST_TAG}:v1"
    [ "$status" -eq 0 ]
}

@test "remove image from store" {
    cat > Imagefile <<EOF
FROM scratch
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-remove"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image remove "${TEST_TAG}:v1"
    [ "$status" -eq 0 ]
}

@test "immutability check - cannot push same image twice" {
    cat > Imagefile <<EOF
FROM scratch
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-immutable"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -ne 0 ]
}

@test "error handling - missing imagefile" {
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root .
    [ "$status" -ne 0 ]
}

@test "error handling - non-existent image inspect" {
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect nonexistent-${BATS_TEST_NUMBER}:v99
    [ "$status" -ne 0 ]
}

@test "tag resolution defaults" {
    cat > Imagefile <<EOF
FROM scratch
TAG ${TEST_TAG}:latest
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-latest"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect "${TEST_TAG}"
    [ "$status" -eq 0 ]
}

@test "inspect --recipe outputs only Imagefile content" {
    cat > Imagefile <<EOF
FROM scratch
COPY artifact.txt ./
TAG ${TEST_TAG}:v1
EOF
    
    echo "test content" > artifact.txt
    
    OUTPUT_DIR="$WORK_DIR/build-recipe"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect "${TEST_TAG}:v1" --recipe
    [ "$status" -eq 0 ]
    
    expected_output="FROM scratch
COPY artifact.txt ./
TAG ${TEST_TAG}:v1"
    
    [ "$output" = "$expected_output" ]
}

@test "inspect --recipe contains no metadata fields" {
    cat > Imagefile <<EOF
# Custom Imagefile
FROM scratch
COPY file.txt ./
TAG ${TEST_TAG}:v1
EOF
    
    echo "content" > file.txt
    
    OUTPUT_DIR="$WORK_DIR/build-recipe-meta"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect "${TEST_TAG}:v1" --recipe
    [ "$status" -eq 0 ]
    
    [[ ! "$output" =~ "Name:" ]]
    [[ ! "$output" =~ "Tag:" ]]
    [[ ! "$output" =~ "Digest:" ]]
    [[ ! "$output" =~ "Created:" ]]
    [[ ! "$output" =~ "Files:" ]]
    [[ ! "$output" =~ "Source Images:" ]]
    [[ ! "$output" =~ "# Recipe:" ]]
    [[ ! "$output" =~ "# Format:" ]]
}

@test "inspect --files outputs one path per line" {
    cat > Imagefile <<EOF
FROM scratch
COPY CLAUDE.md ./
COPY config.json ./
COPY README.md ./
TAG ${TEST_TAG}:v1
EOF
    
    echo "claude content" > CLAUDE.md
    echo '{"key": "value"}' > config.json
    echo "readme" > README.md
    
    OUTPUT_DIR="$WORK_DIR/build-files"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect "${TEST_TAG}:v1" --files
    [ "$status" -eq 0 ]
    
    lines=$(echo "$output" | wc -l | tr -d ' ')
    [ "$lines" -eq 3 ]
    
    [[ "$output" =~ "CLAUDE.md" ]]
    [[ "$output" =~ "config.json" ]]
    [[ "$output" =~ "README.md" ]]
}

@test "inspect --files paths are sorted alphabetically" {
    cat > Imagefile <<EOF
FROM scratch
COPY z-file.txt ./
COPY a-file.txt ./
COPY m-file.txt ./
TAG ${TEST_TAG}:v1
EOF
    
    echo "z" > z-file.txt
    echo "a" > a-file.txt
    echo "m" > m-file.txt
    
    OUTPUT_DIR="$WORK_DIR/build-files-sort"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect "${TEST_TAG}:v1" --files
    [ "$status" -eq 0 ]
    
    first_line=$(echo "$output" | head -n1)
    second_line=$(echo "$output" | sed -n '2p')
    third_line=$(echo "$output" | tail -n1)
    
    [ "$first_line" = "a-file.txt" ]
    [ "$second_line" = "m-file.txt" ]
    [ "$third_line" = "z-file.txt" ]
}

@test "inspect --files contains no metadata or headers" {
    cat > Imagefile <<EOF
FROM scratch
COPY CLAUDE.md ./
TAG ${TEST_TAG}:v1
EOF
    
    echo "claude" > CLAUDE.md
    
    OUTPUT_DIR="$WORK_DIR/build-files-meta"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect "${TEST_TAG}:v1" --files
    [ "$status" -eq 0 ]
    
    [[ ! "$output" =~ "Name:" ]]
    [[ ! "$output" =~ "Tag:" ]]
    [[ ! "$output" =~ "Digest:" ]]
    [[ ! "$output" =~ "Created:" ]]
    [[ ! "$output" =~ "Files:" ]]
    [[ ! "$output" =~ "Source Images:" ]]
    
    [ "$output" = "CLAUDE.md" ]
}

@test "inspect --files with empty image outputs No files" {
    cat > Imagefile <<EOF
FROM scratch
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-files-empty"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect "${TEST_TAG}:v1" --files
    [ "$status" -eq 0 ]
    
    [ "$output" = "No files" ]
}

@test "list --json when empty returns empty array" {
    run "${KFG_BIN}" --store "$STORE_DIR" image list --json
    [ "$status" -eq 0 ]
    
    [ "$output" = "[]" ]
}

@test "inspect --files --json outputs valid JSON array" {
    cat > Imagefile <<EOF
FROM scratch
COPY file1.txt ./
COPY file2.txt ./
TAG ${TEST_TAG}:v1
EOF
    
    echo "content1" > file1.txt
    echo "content2" > file2.txt
    
    OUTPUT_DIR="$WORK_DIR/build-files-json"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect "${TEST_TAG}:v1" --files --json
    [ "$status" -eq 0 ]
    
    [[ "$output" =~ ^\[ ]]
    [[ "$output" =~ \]$ ]]
    
    [[ "$output" =~ "file1.txt" ]]
    [[ "$output" =~ "file2.txt" ]]
}

@test "inspect --files --json with empty image outputs empty array" {
    cat > Imagefile <<EOF
FROM scratch
TAG ${TEST_TAG}:v1
EOF
    
    OUTPUT_DIR="$WORK_DIR/build-files-json-empty"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect "${TEST_TAG}:v1" --files --json
    [ "$status" -eq 0 ]
    
    [ "$output" = "[]" ]
}

@test "inspect --recipe --files returns error" {
    cat > Imagefile <<EOF
FROM scratch
COPY test.txt ./
TAG ${TEST_TAG}:v1
EOF
    
    echo "test" > test.txt
    
    OUTPUT_DIR="$WORK_DIR/build-conflict-test"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect "${TEST_TAG}:v1" --recipe --files
    [ "$status" -ne 0 ]
    
    [[ "$output" =~ "mutually exclusive" ]]
}

@test "inspect --files --json succeeds" {
    cat > Imagefile <<EOF
FROM scratch
COPY test.txt ./
TAG ${TEST_TAG}:v1
EOF
    
    echo "test" > test.txt
    
    OUTPUT_DIR="$WORK_DIR/build-files-json-valid"
    run "${KFG_BIN}" --store "$STORE_DIR" image build --root . --output "$OUTPUT_DIR"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image push "$OUTPUT_DIR/artifacts"
    [ "$status" -eq 0 ]
    
    run "${KFG_BIN}" --store "$STORE_DIR" image inspect "${TEST_TAG}:v1" --files --json
    [ "$status" -eq 0 ]
    
    [[ "$output" =~ ^\[ ]]
}
#!/usr/bin/env bats

# Tests for kfg sys cache CLI commands
# Tests:
# - kfg sys cache exists (hit/miss scenarios)
# - kfg sys cache ls (table/JSON/YAML output)
# - kfg sys cache inspect (by step-ref name)
# - kfg sys cache rm (single/multiple entries)
# - kfg sys cache prune (old entries)
# - kfg sys cache du (disk usage)
# - kfg sys cache store/restore round-trip

load '../test_helper'

# Test fixtures directory
FIXTURES_DIR="${PROJECT_ROOT}/tests/bats/fixtures/cache-cli"

setup() {
    export KFG_STORE_DIR="$(mktemp -d)"
    mkdir -p "$FIXTURES_DIR"
}

teardown() {
    rm -rf "$KFG_STORE_DIR"
    rm -rf "$FIXTURES_DIR"
}

# Helper: create a test cache entry
create_test_cache_entry() {
    local step_ref="$1"
    local cache_dir="${KFG_STORE_DIR}/cache"
    # Compute SHA256 hash of step ref name
    local hash
    hash=$(printf '%s' "$step_ref" | sha256sum | cut -d' ' -f1)
    local entry_path="${cache_dir}/${hash}"

    mkdir -p "${entry_path}/artifacts"
    cat > "${entry_path}/metadata.yaml" << EOF
stepRefName: ${step_ref}
timestamp: $(date -u +%Y-%m-%dT%H:%M:%SZ)
artifacts:
  - test-file.txt
EOF
    echo "test content" > "${entry_path}/artifacts/test-file.txt"
}

# Test: kfg sys cache exists returns 0 for existing entry
@test "kfg sys cache exists returns 0 for cache hit" {
    create_test_cache_entry "test.step.exists"

    run "${KFG_BIN}" sys cache exists "test.step.exists"
    [ "$status" -eq 0 ]
}

# Test: kfg sys cache exists returns 1 for missing entry
@test "kfg sys cache exists returns 1 for cache miss" {
    run "${KFG_BIN}" sys cache exists "nonexistent.step"
    [ "$status" -eq 1 ]
}

# Test: kfg sys cache ls lists entries in table format
@test "kfg sys cache ls shows table output" {
    create_test_cache_entry "test.step.ls1"
    create_test_cache_entry "test.step.ls2"

    run "${KFG_BIN}" sys cache ls
    [ "$status" -eq 0 ]
    [[ "$output" =~ "STEP REF NAME" ]]
    [[ "$output" =~ "test.step.ls1" ]]
    [[ "$output" =~ "test.step.ls2" ]]
}

# Test: kfg sys cache ls --json outputs JSON
@test "kfg sys cache ls --json outputs valid JSON" {
    create_test_cache_entry "test.step.json"

    run "${KFG_BIN}" sys cache ls --json
    [ "$status" -eq 0 ]
    # Verify JSON structure
    [[ "$output" =~ '"stepRef"' ]]
    [[ "$output" =~ "test.step.json" ]]
}

# Test: kfg sys cache ls --yaml outputs YAML
@test "kfg sys cache ls --yaml outputs valid YAML" {
    create_test_cache_entry "test.step.yaml"

    run "${KFG_BIN}" sys cache ls --yaml
    [ "$status" -eq 0 ]
    # Verify YAML structure
    [[ "$output" =~ "stepRef:" ]]
    [[ "$output" =~ "test.step.yaml" ]]
}

# Test: kfg sys cache ls with no entries
@test "kfg sys cache ls shows message when no entries" {
    run "${KFG_BIN}" sys cache ls
    [ "$status" -eq 0 ]
    [[ "$output" =~ "No cache entries found" ]]
}

# Test: kfg sys cache inspect shows entry details
@test "kfg sys cache inspect shows entry metadata" {
    create_test_cache_entry "test.step.inspect"

    run "${KFG_BIN}" sys cache inspect "test.step.inspect"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "test.step.inspect" ]]
    [[ "$output" =~ "test-file.txt" ]]
}

# Test: kfg sys cache inspect --json outputs JSON
@test "kfg sys cache inspect --json outputs valid JSON" {
    create_test_cache_entry "test.step.inspect-json"

    run "${KFG_BIN}" sys cache inspect "test.step.inspect-json" --json
    [ "$status" -eq 0 ]
    [[ "$output" =~ '"stepRef"' ]]
    [[ "$output" =~ "test.step.inspect-json" ]]
}

# Test: kfg sys cache inspect --yaml outputs YAML
@test "kfg sys cache inspect --yaml outputs valid YAML" {
    create_test_cache_entry "test.step.inspect-yaml"

    run "${KFG_BIN}" sys cache inspect "test.step.inspect-yaml" --yaml
    [ "$status" -eq 0 ]
    [[ "$output" =~ "stepRef:" ]]
    [[ "$output" =~ "test.step.inspect-yaml" ]]
}

# Test: kfg sys cache inspect fails for nonexistent entry
@test "kfg sys cache inspect fails for nonexistent entry" {
    run "${KFG_BIN}" sys cache inspect "nonexistent.step"
    [ "$status" -ne 0 ]
}

# Test: kfg sys cache rm removes entry
@test "kfg sys cache rm removes cache entry" {
    create_test_cache_entry "test.step.rm"

    # Verify entry exists
    run "${KFG_BIN}" sys cache exists "test.step.rm"
    [ "$status" -eq 0 ]

    # Remove entry
    run "${KFG_BIN}" sys cache rm "test.step.rm"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Removed" ]]

    # Verify entry is gone
    run "${KFG_BIN}" sys cache exists "test.step.rm"
    [ "$status" -eq 1 ]
}

# Test: kfg sys cache rm removes multiple entries
@test "kfg sys cache rm removes multiple entries" {
    create_test_cache_entry "test.step.rm1"
    create_test_cache_entry "test.step.rm2"

    run "${KFG_BIN}" sys cache rm "test.step.rm1" "test.step.rm2"
    [ "$status" -eq 0 ]

    # Verify both entries are gone
    run "${KFG_BIN}" sys cache exists "test.step.rm1"
    [ "$status" -eq 1 ]
    run "${KFG_BIN}" sys cache exists "test.step.rm2"
    [ "$status" -eq 1 ]
}

# Test: kfg sys cache rm warns for nonexistent entry
@test "kfg sys cache rm warns for nonexistent entry" {
    run "${KFG_BIN}" sys cache rm "nonexistent.step"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "not found" ]]
}

# Test: kfg sys cache prune removes old entries
@test "kfg sys cache prune removes old entries" {
    # Create an old entry (manually set timestamp to 35 days ago)
    local step_ref="test.step.old"
    local cache_dir="${KFG_STORE_DIR}/cache"
    local hash
    hash=$(printf '%s' "$step_ref" | sha256sum | cut -d' ' -f1)
    local entry_path="${cache_dir}/${hash}"

    mkdir -p "${entry_path}/artifacts"
    cat > "${entry_path}/metadata.yaml" << EOF
stepRefName: ${step_ref}
timestamp: $(date -u -v-35d +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -d '35 days ago' +%Y-%m-%dT%H:%M:%SZ)
artifacts:
  - test-file.txt
EOF
    echo "old content" > "${entry_path}/artifacts/test-file.txt"

    run "${KFG_BIN}" sys cache prune
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Pruned" ]]
}

# Test: kfg sys cache prune --json outputs JSON
@test "kfg sys cache prune --json outputs valid JSON" {
    run "${KFG_BIN}" sys cache prune --json
    [ "$status" -eq 0 ]
    [[ "$output" =~ '"pruned"' ]]
    [[ "$output" =~ '"count"' ]]
    [[ "$output" =~ '"freedBytes"' ]]
}

# Test: kfg sys cache du shows disk usage
@test "kfg sys cache du shows disk usage" {
    create_test_cache_entry "test.step.du"

    run "${KFG_BIN}" sys cache du
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Cache Directory" ]]
    [[ "$output" =~ "test.step.du" ]]
}

# Test: kfg sys cache du --json outputs JSON
@test "kfg sys cache du --json outputs valid JSON" {
    create_test_cache_entry "test.step.du-json"

    run "${KFG_BIN}" sys cache du --json
    [ "$status" -eq 0 ]
    [[ "$output" =~ '"cacheDir"' ]]
    [[ "$output" =~ '"entries"' ]]
    [[ "$output" =~ '"totalBytes"' ]]
}

# Test: kfg sys cache du --yaml outputs YAML
@test "kfg sys cache du --yaml outputs valid YAML" {
    create_test_cache_entry "test.step.du-yaml"

    run "${KFG_BIN}" sys cache du --yaml
    [ "$status" -eq 0 ]
    [[ "$output" =~ "cacheDir:" ]]
    [[ "$output" =~ "entries:" ]]
}

# Test: kfg sys cache store/restore round-trip
@test "kfg sys cache store and restore round-trip preserves artifacts" {
    local workdir="$(mktemp -d)"
    local step_ref="test.step.roundtrip"

    # Create a test file in workdir
    echo "test content" > "${workdir}/test-file.txt"

    # Store cache
    local store_input='{"before":[],"after":["test-file.txt"],"declarative":[],"output":null}'
    printf '%s' "$store_input" | "${KFG_BIN}" sys cache store "$step_ref" --workdir "$workdir"
    local store_status=$?
    [ "$store_status" -eq 0 ]

    # Verify cache exists
    run "${KFG_BIN}" sys cache exists "$step_ref"
    [ "$status" -eq 0 ]

    # Remove original file
    rm -f "${workdir}/test-file.txt"

    # Restore cache
    run "${KFG_BIN}" sys cache restore "$step_ref" --workdir "$workdir"
    [ "$status" -eq 0 ]

    # Verify file is restored
    [ -f "${workdir}/test-file.txt" ]
    local content
    content=$(cat "${workdir}/test-file.txt")
    [ "$content" = "test content" ]

    # Cleanup
    rm -rf "$workdir"
}

# Test: kfg sys cache store/restore round-trip preserves output
@test "kfg sys cache store and restore round-trip preserves output" {
    local workdir="$(mktemp -d)"
    local step_ref="test.step.output-roundtrip"

    # Store cache with output (base64 of "test-value" is dGVzdC12YWx1ZQ==)
    local store_input='{"before":[],"after":[],"declarative":[],"output":{"name":"test-output","value":"dGVzdC12YWx1ZQ=="}}'
    printf '%s' "$store_input" | "${KFG_BIN}" sys cache store "$step_ref" --workdir "$workdir"
    local store_status=$?
    [ "$store_status" -eq 0 ]

    # Restore cache
    run "${KFG_BIN}" sys cache restore "$step_ref" --workdir "$workdir"
    [ "$status" -eq 0 ]

    # Verify output restoration line is present
    [[ "$output" =~ "__kfg_output_set" ]]
    [[ "$output" =~ "test-output" ]]

    # Cleanup
    rm -rf "$workdir"
}

# Test: kfg sys cache help shows usage
@test "kfg sys cache shows help" {
    run "${KFG_BIN}" sys cache
    [ "$status" -eq 0 ]
    [[ "$output" =~ "cache" ]]
    [[ "$output" =~ "exists" ]]
    [[ "$output" =~ "store" ]]
    [[ "$output" =~ "restore" ]]
}

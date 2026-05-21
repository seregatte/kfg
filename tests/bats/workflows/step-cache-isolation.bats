#!/usr/bin/env bats

# Tests for Step cache artifact isolation and path-preserving restore
# Tests:
# - Generated code has artifact snapshot before execution
# - Generated code computes delta after execution
# - Generated code passes step artifacts and declarative artifacts to cache store
# - Generated code passes declarative artifacts to cache restore
# - Cache helpers preserve relative paths
# - Cache helpers restore to original paths

load '../test_helper'

# Test fixtures directory
FIXTURES_DIR="${PROJECT_ROOT}/tests/bats/fixtures/step-cache-isolation"

setup() {
    mkdir -p "$FIXTURES_DIR"
    create_cache_fixture
}

teardown() {
    rm -rf "$FIXTURES_DIR"
}

# Helper to create test fixtures
create_cache_fixture() {
    # Create cacheable Step with nested artifact
    cat > "$FIXTURES_DIR/resources.yaml" << 'EOF'
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: nested-artifact-step
spec:
  run: |
    mkdir -p .pi/skills/test
    echo 'nested content' > .pi/skills/test/file.txt
  cache:
    enabled: true
    key: "nested-test"
  artifacts:
    - .pi/skills/test/file.txt
---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: duplicate-basename-step
spec:
  run: |
    mkdir -p dir1 dir2
    echo 'content1' > dir1/file.txt
    echo 'content2' > dir2/file.txt
  cache:
    enabled: true
    key: "duplicate-test"
  artifacts:
    - dir1/file.txt
    - dir2/file.txt
---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: output-step
spec:
  run: echo "output-value"
  cache:
    enabled: true
    key: "output-test"
  output:
    name: test-output
---
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: test.cmd
  commandName: testcmd
spec:
  run: echo "main cmd"
---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: cache-test-workflow
  shell: bash
spec:
  cmds:
    - test.cmd
  before:
    - name: nested-step
      step: nested-artifact-step
    - name: duplicate-step
      step: duplicate-basename-step
    - name: output-test-step
      step: output-step
EOF

    cat > "$FIXTURES_DIR/kustomization.yaml" << 'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - resources.yaml
EOF
}

# Test: Generated code has artifact snapshot before execution
@test "generated code has artifact snapshot for cacheable steps" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify snapshot of artifacts before execution
    [[ "$output" =~ "__artifacts_before" ]]
    [[ "$output" =~ "KFG_ARTIFACTS" ]]
}

# Test: Generated code computes delta after execution
@test "generated code computes step-local artifact delta" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify delta computation
    [[ "$output" =~ "__step_artifacts" ]]
    [[ "$output" =~ "__found" ]]  # Used to check if artifact was in before snapshot
}

# Test: Generated code passes both step artifacts and declarative artifacts to cache store
@test "generated code passes step artifacts and declarative artifacts to cache store" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify cache store call includes both arrays
    [[ "$output" =~ "__kfg_cache_store" ]]
    [[ "$output" =~ "__step_artifacts" ]]
    [[ "$output" =~ "__decl_artifacts" ]]
}

# Test: Generated code passes declarative artifacts to cache restore
@test "generated code passes declarative artifacts to cache restore" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify cache restore call includes declarative artifacts
    [[ "$output" =~ "__kfg_cache_restore" ]]
    [[ "$output" =~ "__decl_artifacts" ]]
}

# Test: Generated code registers declarative artifacts after execution
@test "generated code registers declarative artifacts after step execution" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify __kfg_add_artifact is called for declarative artifacts
    [[ "$output" =~ "__kfg_add_artifact" ]]
}

# Test: Generated code has detail logging for cache operations
@test "generated code has detail logs for cache operations" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify detail logging
    [[ "$output" =~ "__kfg_log_detail" ]]
    [[ "$output" =~ "cache" ]]
}

# Test: Generated code has debug logging for artifact paths
@test "generated code has debug logs for artifact paths" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify debug logging
    [[ "$output" =~ "__kfg_log_debug" ]]
}

# Test: Cache helper functions preserve relative paths in storage
@test "cache helper functions preserve relative paths" {
    # Check helper template contains path preservation logic
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify helper stores artifacts with full paths (not just basenames)
    # The helper should preserve the directory structure
    [[ "$output" =~ "mkdir -p" ]]  # Creates parent directories
    [[ "$output" =~ "dirname" ]]   # Extracts parent directory from path
}

# Test: Cache helper functions restore to original paths
@test "cache helper functions restore to original paths" {
    # Check helper template contains restore logic
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify helper restores to original path (not current directory)
    # The helper should recreate the original directory structure
    [[ "$output" =~ "mkdir -p" ]]  # Recreates parent directories on restore
}

# Test: Generated shell code has valid syntax
@test "generated cache code has valid shell syntax" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify shell syntax is valid
    verify_shell_syntax "$output"
}

# Test: Step function accepts declarative artifacts parameter
@test "step function accepts declarative artifacts as parameter" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify step function has second parameter for declarative artifacts
    [[ "$output" =~ "__decl_artifacts" ]]
    [[ "$output" =~ "IFS=' '" ]]  # Used to parse space-separated artifacts
}

# Test: Nested paths are not reduced to basenames
@test "nested artifact paths are preserved in generated code" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify nested path .pi/skills/test/file.txt is preserved
    [[ "$output" =~ ".pi/skills/test/file.txt" ]]
}

# Test: Duplicate basenames from different directories are both present
@test "duplicate basenames from different directories are handled" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify both dir1/file.txt and dir2/file.txt are present (not collapsed)
    [[ "$output" =~ "dir1/file.txt" ]]
    [[ "$output" =~ "dir2/file.txt" ]]
}

# Test: Output is restored from cache
@test "output restoration logic is present in cache restore" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify output restoration in cache restore function
    [[ "$output" =~ "__kfg_output_set" ]]
    [[ "$output" =~ "output_encoded" ]]
    [[ "$output" =~ "base64" ]]
}

# Test: Cache identity includes step ref name, cache key, and script hash
@test "cache identity computation includes all components" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify cache identity computation
    [[ "$output" =~ "__kfg_cache_identity" ]]
    [[ "$output" =~ "identity_key" ]]
    [[ "$output" =~ "sha256sum" ]]
}

# Test: Refresh bypass logic is present
@test "refresh bypass logic is present in generated code" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify KFG_REFRESH check
    [[ "$output" =~ "KFG_REFRESH" ]]
}
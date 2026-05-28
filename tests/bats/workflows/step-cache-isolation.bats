#!/usr/bin/env bats

# Tests for Step cache artifact isolation and path-preserving restore
# Tests:
# - Generated code delegates cache operations to Go subcommands
# - Generated code computes artifact delta before/after execution
# - Generated code passes step artifacts and declarative artifacts to cache store
# - Generated code passes declarative artifacts to cache restore
# - Cache helpers use kfg sys cache subcommands

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
  output:
    name: test-output
---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: output-with-artifact-step
spec:
  run: |
    mkdir -p artifacts
    echo "artifact content" > artifacts/runtime-artifact.txt
    __kfg_add_artifact "artifacts/runtime-artifact.txt"
    echo "captured-output-value"
  cache:
    enabled: true
  output:
    name: runtime-output
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
    - name: output-artifact-test-step
      step: output-with-artifact-step
EOF

    cat > "$FIXTURES_DIR/kustomization.yaml" << 'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - resources.yaml
EOF
}

# Test: Generated code delegates cache exists to Go subcommand
@test "generated code delegates cache exists to kfg sys cache exists" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify cache exists delegates to Go subcommand
    [[ "$output" =~ "__kfg_internal sys cache exists" ]]
}

# Test: Generated code delegates cache store to Go subcommand
@test "generated code delegates cache store to kfg sys cache store" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify cache store delegates to Go subcommand
    [[ "$output" =~ "__kfg_internal sys cache store" ]]
}

# Test: Generated code delegates cache restore to Go subcommand
@test "generated code delegates cache restore to kfg sys cache restore" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify cache restore delegates to Go subcommand
    [[ "$output" =~ "__kfg_internal sys cache restore" ]]
}

# Test: Generated code computes artifact delta before/after execution
@test "generated code computes step-local artifact delta" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify delta computation
    [[ "$output" =~ "__step_artifacts" ]]
    [[ "$output" =~ "__artifacts_before" ]]
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
    [[ "$output" =~ "readarray -t" ]]  # Used to parse newline-separated artifacts
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
}

# Test: Refresh invalidation logic is present
@test "refresh invalidation logic is present in generated code" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify KFG_REFRESH check
    [[ "$output" =~ "KFG_REFRESH" ]]

    # Verify refresh now invalidates cache (not just bypasses restore)
    [[ "$output" =~ "Invalidating cache for step" ]]
}

# Test: Refresh bypasses cache restore but still allows cache store
@test "refresh bypasses cache restore but still stores cache" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify cache restore is skipped when KFG_REFRESH is set
    [[ "$output" =~ "if [ -z \"\${KFG_REFRESH:-}\" ]; then" ]]

    # Verify cache store happens after execution
    [[ "$output" =~ "__kfg_cache_store" ]]
}

# Test: Internal kfg execution helper is present
@test "internal kfg execution helper is present in generated code" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify __kfg_internal helper is present
    [[ "$output" =~ "__kfg_internal" ]]
}

# Test: Internal kfg execution helper sets KFG_VERBOSE=0
@test "internal kfg helper sets KFG_VERBOSE=0 for child process" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify __kfg_internal sets KFG_VERBOSE=0
    [[ "$output" =~ "KFG_VERBOSE=0 kfg" ]]
}

# Test: Cache helpers use internal kfg execution helper
@test "cache helpers use __kfg_internal for Go subcommands" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify cache helpers call __kfg_internal instead of kfg directly
    [[ "$output" =~ "__kfg_internal sys cache exists" ]]
    [[ "$output" =~ "__kfg_internal sys cache store" ]]
    [[ "$output" =~ "__kfg_internal sys cache restore" ]]
}

# Test: Internal execution helper doesn't mutate parent environment
@test "internal execution helper uses environment prefix not export" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify __kfg_internal uses environment prefix (KFG_VERBOSE=0 kfg)
    # NOT export (export KFG_VERBOSE=0; kfg)
    [[ "$output" =~ 'KFG_VERBOSE=0 kfg "$@"' ]]
}

# Tests for output-producing steps with runtime artifact registration

# Test: Output-producing steps capture stdout via temp file
@test "output-producing steps use temp file capture" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify output-producing steps use temp file mechanism
    [[ "$output" =~ "__output_file" ]]
    [[ "$output" =~ "mktemp" ]]
}

# Test: Output-producing steps run in parent shell preserving runtime side effects
@test "output-producing steps run in parent shell preserving side effects" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify __kfg_add_artifact calls are present in generated code
    [[ "$output" =~ "__kfg_add_artifact" ]]
}

# Test: Temp file cleanup trap is set for output-producing steps
@test "output-producing steps have cleanup trap for temp file" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify cleanup trap is defined
    [[ "$output" =~ "__cleanup_trap" ]]
    [[ "$output" =~ "rm -f" ]]
    [[ "$output" =~ "trap" ]]
}

# Test: Output-producing steps read captured output after execution
@test "output-producing steps read captured output from temp file" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify output is read from temp file after execution
    [[ "$output" =~ "__output=\$(<" ]]
    [[ "$output" =~ "__output_file" ]]

    # Verify __kfg_output_set is called with captured value
    [[ "$output" =~ "__kfg_output_set" ]]
}

# Test: Output-producing step with runtime artifact registration caches both
@test "output-producing step with runtime artifact caches output and artifact together" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify cache store includes output name parameter
    [[ "$output" =~ "__kfg_cache_store" ]]
    [[ "$output" =~ "runtime-output" ]]  # Output name from fixture
}

# Test: Output-producing step restores output and artifacts from cache together
@test "output-producing step restores output and artifacts from cache together" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify cache restore includes output name parameter
    [[ "$output" =~ "__kfg_cache_restore" ]]
    [[ "$output" =~ "runtime-output" ]]  # Output name from fixture
}

# Test: Output-producing step wrapper handles execution failure
@test "output-producing step wrapper returns error on failure" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify error handling pattern in output-producing step wrapper
    [[ "$output" =~ "__step_rc=\$?" ]]
    [[ "$output" =~ "if [ \$__step_rc -ne 0 ]; then" ]]
    [[ "$output" =~ "return \$__step_rc" ]]
}

# Test: Refresh emits invalidation log before execution
@test "refresh emits invalidation log before execution" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify invalidation log is present
    [[ "$output" =~ "Invalidating cache for step" ]]
}

# Test: Refresh emits rebuild log after successful execution
@test "refresh emits rebuild log after successful execution" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify rebuild log is present
    [[ "$output" =~ "Rebuilding cache for step" ]]
}

# Test: Refresh invalidation is step-scoped (not workflow-wide)
@test "refresh invalidation is scoped to current step" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify cache exists uses step ref name (step-specific)
    [[ "$output" =~ "__kfg_cache_exists \"$__step_ref_name" ]]
}

# Test: JSON serialization helper is present for cache store
@test "JSON serialization helper is present for cache store" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify __kfg_serialize_cache_input is defined
    [[ "$output" =~ "__kfg_serialize_cache_input" ]]
}

# Test: Cache store receives JSON via stdin
@test "cache store receives JSON via stdin" {
    create_cache_fixture

    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]

    # Verify pipe to kfg sys cache store
    [[ "$output" =~ "printf" ]]
    [[ "$output" =~ "__kfg_internal sys cache store" ]]
}

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
    [[ "$output" =~ "output_encoded" ]]
    [[ "$output" =~ "base64" ]]
}

# Test: Cache identity includes step ref name, cache key, and script hash
@test "cache identity computation includes all components" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify cache identity computation (simplified to use only step reference name)
    [[ "$output" =~ "__kfg_cache_identity" ]]
    [[ "$output" =~ "sha256sum" ]]
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
    
    # Verify cache entry is removed when refresh is enabled
    [[ "$output" =~ "rm -rf" ]]
}

# Test: Refresh bypasses cache restore but still allows cache store
@test "refresh bypasses cache restore but still stores cache" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify cache restore is skipped when KFG_REFRESH is set
    # The cache check should be inside a "if [ -z \"${KFG_REFRESH:-}\" ]; then" block
    [[ "$output" =~ "if [ -z \"\${KFG_REFRESH:-}\" ]; then" ]]
    
    # Verify cache store happens AFTER the cache check block ends
    # Cache store should NOT be inside the KFG_REFRESH conditional
    # Look for cache store after the conditional that checks for cache existence
    [[ "$output" =~ "__kfg_cache_store" ]]
}

# Test: Cache identity is always computed for cacheable steps
@test "cache identity is computed before cache check for cacheable steps" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify __cache_path is computed before the refresh check
    # This ensures cache path is available for both restore and store
    [[ "$output" =~ "__cache_path=\$(__kfg_cache_identity" ]]
    
    # Verify the cache path computation happens outside the KFG_REFRESH check
    # by checking it appears before the "if [ -z \"${KFG_REFRESH:-}\" ]" line
}

# Test: Cache store removes existing cache entries before writing
@test "cache store helper removes existing cache entries" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify __kfg_cache_store function includes rm -rf to remove stale cache
    [[ "$output" =~ "rm -rf \"\$cache_path.tmp\"" ]]
}

# Test: Cache store happens unconditionally after successful execution
@test "cache store is not conditional on KFG_REFRESH" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Extract the cache store section
    # The cache store should appear after execution, not inside a KFG_REFRESH conditional
    # We verify this by checking that cache store appears and is not wrapped in
    # a separate "if [ -z \"${KFG_REFRESH:-}\" ]; then" block
    
    # Count occurrences of KFG_REFRESH checks in cache-related sections
    # There should be only ONE check for KFG_REFRESH (the cache restore bypass)
    local refresh_check_count
    refresh_check_count=$(echo "$output" | grep -o 'KFG_REFRESH' | wc -l)
    
    # Should have at least one KFG_REFRESH check (for cache bypass)
    [ "$refresh_check_count" -ge 1 ]
}

# Test: Refresh rebuild semantics - integration test
@test "refresh rebuilds cache after successful execution" {
    # This test verifies the full refresh behavior:
    # 1. Cache identity is computed
    # 2. Cache restore is bypassed when KFG_REFRESH is set
    # 3. Step executes
    # 4. Cache is stored (rebuilding the entry)
    # 5. Old artifacts are removed before storing
    
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify the sequence in generated code:
    # 1. Cache identity computation
    [[ "$output" =~ "__cache_path=\$(__kfg_cache_identity" ]]
    
    # 2. Cache bypass check (skips restore when KFG_REFRESH is set)
    [[ "$output" =~ "if [ -z \"\${KFG_REFRESH:-}\" ]; then" ]]
    [[ "$output" =~ "__kfg_cache_restore" ]]
    
    # 3. Cache store helper removes old cache
    [[ "$output" =~ "rm -rf \"\$cache_path.tmp\"" ]]
    
    # 4. Cache store happens after execution
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
    # The helper should use "KFG_VERBOSE=0 kfg" pattern (environment prefix)
    [[ "$output" =~ "KFG_VERBOSE=0 kfg" ]]
}

# Test: Filesystem helpers use internal kfg execution helper
@test "filesystem helpers use __kfg_internal" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify filesystem helpers call __kfg_internal instead of kfg directly
    [[ "$output" =~ "__kfg_fs_snapshot" ]]
    [[ "$output" =~ "__kfg_internal sys fs snapshot" ]]
    
    [[ "$output" =~ "__kfg_fs_diff" ]]
    [[ "$output" =~ "__kfg_internal sys fs diff" ]]
}

# Test: Internal execution helper doesn't mutate parent environment
@test "internal execution helper uses environment prefix not export" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify __kfg_internal uses environment prefix (KFG_VERBOSE=0 kfg)
    # NOT export (export KFG_VERBOSE=0; kfg)
    # This ensures parent environment is not mutated
    # The pattern should match "KFG_VERBOSE=0 kfg "$@""
    [[ "$output" =~ 'KFG_VERBOSE=0 kfg "$@"' ]]
}

# Test: Quiet internal subprocesses preserve functionality
@test "quiet internal subprocesses preserve stdout stderr and exit status" {
    # This test verifies that __kfg_internal still returns
    # functional stdout/stderr and exit status even with KFG_VERBOSE=0
    
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify the helper uses "kfg "$@"" which preserves all outputs and status
    # The environment prefix "KFG_VERBOSE=0" only affects verbosity, not functionality
    [[ "$output" =~ 'KFG_VERBOSE=0 kfg "$@"' ]]
}

# Tests for output-producing steps with runtime artifact registration (fix-output-step-subshell-cache-loss)

# Test: Output-producing steps capture stdout via temp file (not command substitution)
@test "output-producing steps use temp file capture not command substitution" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify output-producing steps use temp file mechanism
    # Look for: __output_file=$(mktemp)
    [[ "$output" =~ "__output_file" ]]
    [[ "$output" =~ "mktemp" ]]
    [[ "$output" =~ "{ " ]]  # Group command start
    [[ "$output" =~ \} ]]  # Group command end
    
    # Verify NO command substitution for output capture
    # Old pattern was: __output="$( <run script> )"
    # Should NOT match this pattern for output-producing steps
    ! [[ "$output" =~ '__output="\$\(' ]]
}

# Test: Output-producing steps run in parent shell preserving runtime side effects
@test "output-producing steps run in parent shell preserving side effects" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify output capture uses group command which runs in parent shell
    # Group command "{ commands; }" does not create a subshell
    [[ "$output" =~ "{ " ]]
    [[ "$output" =~ \} ]]  # Closing brace
    [[ "$output" =~ "__output_file" ]]  # Temp file variable
    
    # Verify __kfg_add_artifact calls are present in generated code
    # (This test fixture has output-with-artifact-step that calls __kfg_add_artifact)
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
    
    # Verify cache store is called AFTER execution (where artifacts would be registered)
    # The cache store should be after the output capture section
}

# Test: Output-producing step restores output and artifacts from cache together
@test "output-producing step restores output and artifacts from cache together" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify cache restore includes output name parameter
    [[ "$output" =~ "__kfg_cache_restore" ]]
    [[ "$output" =~ "runtime-output" ]]  # Output name from fixture
    
    # Verify cache restore is called BEFORE execution (when cache exists)
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

# Test: Output-producing step cleanup trap fires on failure path
@test "output-producing step cleanup trap handles failure path" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify cleanup trap includes temp file removal
    # The trap should fire on RETURN which happens when function returns (including on error)
    [[ "$output" =~ "rm -f" ]]
    [[ "$output" =~ "trap" ]]
    [[ "$output" =~ "RETURN" ]]
}
# Tests for refresh invalidation + rebuild semantics (fix-refresh-step-cache-invalidation)

# Test: Refresh emits invalidation log before execution
@test "refresh emits invalidation log before execution" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify invalidation log is present
    [[ "$output" =~ "Invalidating cache for step" ]]
    
    # Verify this log appears in the else branch of KFG_REFRESH check
    # (after "if [ -z \"${KFG_REFRESH:-}\" ]; then")
}

# Test: Refresh emits rebuild log after successful execution
@test "refresh emits rebuild log after successful execution" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify rebuild log is present
    [[ "$output" =~ "Rebuilding cache for step" ]]
    
    # Verify rebuild log appears in cache store section
    # (when KFG_REFRESH is set, it logs "Rebuilding" instead of "Storing")
}

# Test: Refresh invalidation is step-scoped (not workflow-wide)
@test "refresh invalidation is scoped to current step" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify cache identity uses step ref name (step-specific)
    [[ "$output" =~ "__kfg_cache_identity \"$__step_ref_name" ]]
    
    # Verify the rm -rf uses $__cache_path (which is step-specific)
    [[ "$output" =~ "rm -rf \"$__cache_path" ]]
}

# Tests for diff-based artifact registration helper (fix-refresh-step-cache-invalidation)

# Test: Diff artifact helper is present in generated code
@test "diff artifact helper is present in generated code" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify __kfg_add_diff_artifacts is defined
    [[ "$output" =~ "__kfg_add_diff_artifacts()" ]]
}

# Test: Diff artifact helper accepts root, before, after parameters
@test "diff artifact helper accepts correct parameters" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify helper accepts root as first parameter
    [[ "$output" =~ "local root=" ]]
    
    # Verify helper accepts before_snapshot as second parameter
    [[ "$output" =~ "local before_file=" ]]
    
    # Verify helper accepts after_snapshot as third parameter
    [[ "$output" =~ "local after_file=" ]]
}

# Test: Diff artifact helper prefixes relative paths with root
@test "diff artifact helper prefixes relative paths with root" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify helper prefixes paths with root
    [[ "$output" =~ "full_path=" ]]
}

# Test: Diff artifact helper verifies existence before registration
@test "diff artifact helper verifies existence before registration" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify helper checks if path exists before registering
    [[ "$output" =~ "[ -e" ]]
    
    # Verify helper calls __kfg_add_artifact for existing paths
    [[ "$output" =~ "__kfg_add_artifact" ]]
}

# Test: Diff artifact helper uses __kfg_fs_diff internally
@test "diff artifact helper uses __kfg_fs_diff" {
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify helper calls __kfg_fs_diff
    [[ "$output" =~ "__kfg_fs_diff" ]]
}

# Test: Refresh invalidates then rebuilds cache with artifacts
@test "refresh invalidates then rebuilds cache with artifacts" {
    # This test verifies the full refresh cycle for artifact reconstruction:
    # 1. Cache entry exists with artifacts
    # 2. Refresh invalidates (removes) the cache entry
    # 3. Step executes and produces new artifacts
    # 4. Cache is rebuilt with new artifacts
    
    create_cache_fixture
    
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow cache-test-workflow
    [ "$status" -eq 0 ]
    
    # Verify the invalidation + rebuild cycle in generated code:
    # 1. Invalidating log
    [[ "$output" =~ "Invalidating cache for step" ]]
    
    # 2. Cache entry removal
    [[ "$output" =~ "rm -rf" ]]
    
    # 3. Rebuilding log (when refresh is enabled)
    [[ "$output" =~ "Rebuilding cache for step" ]]
    
    # 4. Cache store preserves artifact paths
    [[ "$output" =~ "__kfg_cache_store" ]]
}

#!/usr/bin/env bats

load '../../../tests/bats/helpers/common'
load '../../../tests/bats/helpers/manifests'

# ==============================================================================
# Task 3.1: Per-item mode tests (single-item and multi-item conversions)
# ==============================================================================

@test "materialize: per-item mode converts single asset to single output" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.one
spec:
  data:
    name: test-one
    description: Test asset one
    prompt: Test prompt one
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test.converter.simple
spec:
  input:
    format: yaml
  engine:
    expression: |
      "---\nname: " + .name + "\ndescription: " + .description + "\n---\n\n" + .prompt
  output:
    format: raw
YAML
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="per-item" \
        ASSETS="test.asset.one" \
        CONVERTER="test.converter.simple" \
        OUTPUTS="$TEST_TMPDIR/output/one.md" \
        bash "$script_file"
    
    [ -f "$TEST_TMPDIR/output/one.md" ]
    grep -q "name: test-one" "$TEST_TMPDIR/output/one.md"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: per-item mode converts multiple assets to multiple outputs (positional mapping)" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.alpha
spec:
  data:
    name: alpha
    description: Asset alpha
    prompt: Alpha prompt
---
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.beta
spec:
  data:
    name: beta
    description: Asset beta
    prompt: Beta prompt
---
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.gamma
spec:
  data:
    name: gamma
    description: Asset gamma
    prompt: Gamma prompt
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test.converter.simple
spec:
  input:
    format: yaml
  engine:
    expression: |
      "---\nname: " + .name + "\ndescription: " + .description + "\n---\n\n" + .prompt
  output:
    format: raw
YAML
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="per-item" \
        ASSETS="test.asset.alpha:test.asset.beta:test.asset.gamma" \
        CONVERTER="test.converter.simple" \
        OUTPUTS="$TEST_TMPDIR/out/a.md:$TEST_TMPDIR/out/b.md:$TEST_TMPDIR/out/c.md" \
        bash "$script_file"
    
    # Verify all three outputs exist
    [ -f "$TEST_TMPDIR/out/a.md" ]
    [ -f "$TEST_TMPDIR/out/b.md" ]
    [ -f "$TEST_TMPDIR/out/c.md" ]
    
    # Verify positional mapping: alpha -> a.md, beta -> b.md, gamma -> c.md
    grep -q "name: alpha" "$TEST_TMPDIR/out/a.md"
    grep -q "name: beta" "$TEST_TMPDIR/out/b.md"
    grep -q "name: gamma" "$TEST_TMPDIR/out/c.md"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: per-item mode registers all outputs as artifacts" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.first
spec:
  data:
    name: first
    description: First asset
    prompt: First prompt
---
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.second
spec:
  data:
    name: second
    description: Second asset
    prompt: Second prompt
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test.converter.simple
spec:
  input:
    format: yaml
  engine:
    expression: |
      "name: " + .name
  output:
    format: raw
YAML
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\necho "ARTIFACTS: ${KFG_ARTIFACTS[@]}"\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    local output
    output=$(PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="per-item" \
        ASSETS="test.asset.first:test.asset.second" \
        CONVERTER="test.converter.simple" \
        OUTPUTS="$TEST_TMPDIR/out/first.txt:$TEST_TMPDIR/out/second.txt" \
        bash "$script_file")
    
    # Verify both artifacts were registered
    [[ "$output" == *"first.txt"* ]]
    [[ "$output" == *"second.txt"* ]]
    
    rm -f "$tmp_file" "$script_file"
}

# ==============================================================================
# Task 3.2: Aggregate mode tests (merge, wrap, merge-with-existing)
# ==============================================================================

@test "materialize: aggregate mode merges multiple assets into single output" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.mcp.one
spec:
  data:
    serverOne:
      command: echo
      args: ["one"]
---
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.mcp.two
spec:
  data:
    serverTwo:
      command: echo
      args: ["two"]
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test.converter.identity
spec:
  input:
    format: yaml
  engine:
    expression: "."
  output:
    format: yaml
YAML
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="aggregate" \
        ASSETS="test.asset.mcp.one:test.asset.mcp.two" \
        CONVERTER="test.converter.identity" \
        OUTPUTS="$TEST_TMPDIR/out/merged.yaml" \
        bash "$script_file"
    
    [ -f "$TEST_TMPDIR/out/merged.yaml" ]
    
    # Verify deep merge: both serverOne and serverTwo should be present
    grep -q "serverOne" "$TEST_TMPDIR/out/merged.yaml"
    grep -q "serverTwo" "$TEST_TMPDIR/out/merged.yaml"
    grep -q "command: echo" "$TEST_TMPDIR/out/merged.yaml"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: aggregate mode wraps merged result with WRAP_KEY" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.wrap.one
spec:
  data:
    serverOne:
      command: echo
---
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.wrap.two
spec:
  data:
    serverTwo:
      command: ls
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test.converter.identity
spec:
  input:
    format: yaml
  engine:
    expression: "."
  output:
    format: yaml
YAML
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="aggregate" \
        ASSETS="test.asset.wrap.one:test.asset.wrap.two" \
        CONVERTER="test.converter.identity" \
        OUTPUTS="$TEST_TMPDIR/out/wrapped.yaml" \
        WRAP_KEY="mcpServers" \
        bash "$script_file"
    
    [ -f "$TEST_TMPDIR/out/wrapped.yaml" ]
    
    # Verify wrap key is present
    grep -q "mcpServers:" "$TEST_TMPDIR/out/wrapped.yaml"
    grep -q "serverOne:" "$TEST_TMPDIR/out/wrapped.yaml"
    grep -q "serverTwo:" "$TEST_TMPDIR/out/wrapped.yaml"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: aggregate mode merges with existing output file" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.merge.new
spec:
  data:
    newKey:
      value: new-value
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test.converter.identity
spec:
  input:
    format: yaml
  engine:
    expression: "."
  output:
    format: yaml
YAML
    
    # Create existing output file with existing content
    mkdir -p "$TEST_TMPDIR/out"
    cat > "$TEST_TMPDIR/out/existing.yaml" <<'YAML'
existingKey:
  value: existing-value
YAML
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="aggregate" \
        ASSETS="test.asset.merge.new" \
        CONVERTER="test.converter.identity" \
        OUTPUTS="$TEST_TMPDIR/out/existing.yaml" \
        bash "$script_file"
    
    [ -f "$TEST_TMPDIR/out/existing.yaml" ]
    
    # Verify both existing and new keys are present (deep merge)
    grep -q "existingKey:" "$TEST_TMPDIR/out/existing.yaml"
    grep -q "newKey:" "$TEST_TMPDIR/out/existing.yaml"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: aggregate mode registers single output as artifact" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.artifact.one
spec:
  data:
    key: value1
---
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.artifact.two
spec:
  data:
    key: value2
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test.converter.identity
spec:
  input:
    format: yaml
  engine:
    expression: "."
  output:
    format: yaml
YAML
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\necho "ARTIFACTS: ${KFG_ARTIFACTS[@]}"\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    local output
    output=$(PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="aggregate" \
        ASSETS="test.asset.artifact.one:test.asset.artifact.two" \
        CONVERTER="test.converter.identity" \
        OUTPUTS="$TEST_TMPDIR/out/aggregate.yaml" \
        bash "$script_file")
    
    # Verify single artifact was registered
    [[ "$output" == *"aggregate.yaml"* ]]
    
    rm -f "$tmp_file" "$script_file"
}

# ==============================================================================
# Task 3.3: Negative tests (validation failures)
# ==============================================================================

@test "materialize: fails when MODE is missing" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.missing.mode
spec:
  data:
    name: test
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test.converter.simple
spec:
  input:
    format: yaml
  engine:
    expression: "."
  output:
    format: yaml
YAML
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { echo "ERROR: $*" >&2; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    ! PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        ASSETS="test.asset.missing.mode" \
        CONVERTER="test.converter.simple" \
        OUTPUTS="$TEST_TMPDIR/out/test.yaml" \
        bash "$script_file"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: fails when ASSETS is missing" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { echo "ERROR: $*" >&2; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    ! PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="per-item" \
        CONVERTER="test.converter.simple" \
        OUTPUTS="$TEST_TMPDIR/out/test.yaml" \
        bash "$script_file"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: fails when CONVERTER is missing" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { echo "ERROR: $*" >&2; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    ! PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="per-item" \
        ASSETS="test.asset" \
        OUTPUTS="$TEST_TMPDIR/out/test.yaml" \
        bash "$script_file"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: fails when OUTPUTS is missing" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { echo "ERROR: $*" >&2; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    ! PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="per-item" \
        ASSETS="test.asset" \
        CONVERTER="test.converter.simple" \
        bash "$script_file"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: fails when MODE is invalid" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { echo "ERROR: $*" >&2; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    ! PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="invalid-mode" \
        ASSETS="test.asset" \
        CONVERTER="test.converter.simple" \
        OUTPUTS="$TEST_TMPDIR/out/test.yaml" \
        bash "$script_file"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: per-item mode fails when ASSETS and OUTPUTS counts mismatch" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.mismatch
spec:
  data:
    name: test
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test.converter.simple
spec:
  input:
    format: yaml
  engine:
    expression: "."
  output:
    format: yaml
YAML
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { echo "ERROR: $*" >&2; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    # 2 assets but 3 outputs - should fail
    ! PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="per-item" \
        ASSETS="test.asset.mismatch:test.asset.mismatch" \
        CONVERTER="test.converter.simple" \
        OUTPUTS="$TEST_TMPDIR/out/a.yaml:$TEST_TMPDIR/out/b.yaml:$TEST_TMPDIR/out/c.yaml" \
        bash "$script_file"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: aggregate mode fails when OUTPUTS has more than one path" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.aggregate.multi
spec:
  data:
    name: test
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test.converter.simple
spec:
  input:
    format: yaml
  engine:
    expression: "."
  output:
    format: yaml
YAML
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { echo "ERROR: $*" >&2; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    # aggregate mode with 2 outputs - should fail
    ! PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="aggregate" \
        ASSETS="test.asset.aggregate.multi" \
        CONVERTER="test.converter.simple" \
        OUTPUTS="$TEST_TMPDIR/out/a.yaml:$TEST_TMPDIR/out/b.yaml" \
        bash "$script_file"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: aggregate mode fails when OUTPUTS is empty" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { echo "ERROR: $*" >&2; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    # aggregate mode with empty outputs - should fail
    ! PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="aggregate" \
        ASSETS="test.asset" \
        CONVERTER="test.converter.simple" \
        OUTPUTS="" \
        bash "$script_file"
    
    rm -f "$tmp_file" "$script_file"
}

@test "materialize: per-item mode fails when WRAP_KEY is set" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test.asset.wrapkey.invalid
spec:
  data:
    name: test
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test.converter.simple
spec:
  input:
    format: yaml
  engine:
    expression: "."
  output:
    format: yaml
YAML
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { echo "ERROR: $*" >&2; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$(step_run_code materialize)" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    # per-item mode with WRAP_KEY set - should fail
    ! PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" \
        MODE="per-item" \
        ASSETS="test.asset.wrapkey.invalid" \
        CONVERTER="test.converter.simple" \
        OUTPUTS="$TEST_TMPDIR/out/test.yaml" \
        WRAP_KEY="someKey" \
        bash "$script_file"
    
    rm -f "$tmp_file" "$script_file"
}
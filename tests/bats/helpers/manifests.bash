#!/usr/bin/env bash

# Manifest Bats test helpers
# This file provides helpers for executing and validating manifest resources.

# Manifests directory (default overlay)
MANIFESTS_DIR="${PROJECT_ROOT}/.manifests/overlay/dev"

# Helper: run kfg build and capture output
kfg_build() {
    "$KFG_BIN" build "$MANIFESTS_DIR" 2>/dev/null
}

# Helper: run kfg apply with converter
kfg_convert() {
    local asset="$1"
    local converter="$2"
    "$KFG_BIN" apply -k "$MANIFESTS_DIR" --convert "$asset" --use "$converter" 2>/dev/null
}

# Helper: run kfg apply for workflow
kfg_workflow() {
    local workflow="$1"
    "$KFG_BIN" apply -k "$MANIFESTS_DIR" --workflow "$workflow" 2>/dev/null
}

# Helper: count resources of a kind in build output
count_kind() {
    local kind="$1"
    local output="$2"
    echo "$output" | grep -c "^kind: ${kind}$"
}

# Helper: get step file path (accepts step name or full path)
step_file() {
    local step="$1"
    if [[ "$step" == *".yaml" ]]; then
        echo "${PROJECT_ROOT}/${step}"
    else
        echo "${PROJECT_ROOT}/.manifests/base/steps/${step}.yaml"
    fi
}

# Helper: extract step run code from yaml
step_run_code() {
    local step="$1"
    yq eval '.spec.run' "$(step_file "$step")" 2>/dev/null
}

# Helper: mock _kfg.log functions
mock_kfg_log() {
    cat <<'EOF'
_kfg.log.info() { :; }
_kfg.log.warn() { :; }
_kfg.log.error() { :; }
_kfg.log.debug() { :; }
EOF
}

# Helper: mock artifact tracking functions
mock_artifact_tracking() {
    cat <<'EOF'
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() {
    KFG_ARTIFACTS+=("$1")
}
EOF
}

# Helper: run step in sandbox directory
run_step() {
    local step="$1"
    local dirs="$2"
    cd "$TEST_TMPDIR"
    local code
    code=$(printf '%s\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test' "$(mock_kfg_log)" "$(step_run_code "$step")")
    DIRECTORIES="$dirs" bash -c "$code"
}

# Helper: run step in sandbox directory with env vars
run_step_with_env() {
    local step="$1"
    local code
    code=$(printf '%s\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test' "$(mock_kfg_log)" "$(step_run_code "$step")")
    bash -c "$code"
}

# Helper: run step with artifact tracking mock
run_step_with_artifacts() {
    local step="$1"
    local code
    code=$(printf '%s\n%s\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test' "$(mock_kfg_log)" "$(mock_artifact_tracking)" "$(step_run_code "$step")")
    bash -c "$code"
}
#!/usr/bin/env bash

# Manifest Bats test helpers
# This file provides helpers for executing and validating manifest resources.
# Updated to support package-oriented repository structure.

# Default manifests directory (can be overridden per test)
# For backward compatibility, check if old .manifests path exists, otherwise use new package path
if [[ -d "${PROJECT_ROOT}/.manifests/overlay/dev" ]]; then
    DEFAULT_MANIFESTS_DIR="${PROJECT_ROOT}/.manifests/overlay/dev"
elif [[ -d "${PROJECT_ROOT}/packages/domains/ai-agents/overlays/dev" ]]; then
    DEFAULT_MANIFESTS_DIR="${PROJECT_ROOT}/packages/domains/ai-agents/overlays/dev"
else
    DEFAULT_MANIFESTS_DIR="${PROJECT_ROOT}/packages/domains/ai-agents/overlays/dev"
fi

# Allow tests to override manifest directory via MANIFESTS_DIR
MANIFESTS_DIR="${MANIFESTS_DIR:-$DEFAULT_MANIFESTS_DIR}"

# Helper: run kfg build and capture output
kfg_build() {
    local path="${1:-$MANIFESTS_DIR}"
    "$KFG_BIN" build "$path" 2>/dev/null
}

# Helper: run kfg apply with converter
kfg_convert() {
    local asset="$1"
    local converter="$2"
    local path="${3:-$MANIFESTS_DIR}"
    "$KFG_BIN" apply -k "$path" --convert "$asset" --use "$converter" 2>/dev/null
}

# Helper: run kfg apply for workflow
kfg_workflow() {
    local workflow="$1"
    local path="${2:-$MANIFESTS_DIR}"
    "$KFG_BIN" apply -k "$path" --workflow "$workflow" 2>/dev/null
}

# Helper: count resources of a kind in build output
count_kind() {
    local kind="$1"
    local output="$2"
    echo "$output" | grep -c "^kind: ${kind}$"
}

# Helper: find step file path - supports new package structure
# First tries framework package steps, then falls back to old .manifests path
step_file() {
    local step="$1"
    local step_name="${step%.yaml}"  # Remove .yaml extension if present
    
    if [[ "$step" == *".yaml" ]]; then
        # Full path provided
        echo "${PROJECT_ROOT}/${step}"
    elif [[ -f "${PROJECT_ROOT}/packages/framework/manifests/steps/${step_name}.yaml" ]]; then
        # Framework step
        echo "${PROJECT_ROOT}/packages/framework/manifests/steps/${step_name}.yaml"
    elif [[ -f "${PROJECT_ROOT}/.manifests/base/steps/${step_name}.yaml" ]]; then
        # Fallback to old location (for backward compatibility)
        echo "${PROJECT_ROOT}/.manifests/base/steps/${step_name}.yaml"
    else
        # Return the expected new path anyway (test will fail if file doesn't exist)
        echo "${PROJECT_ROOT}/packages/framework/manifests/steps/${step_name}.yaml"
    fi
}

# Helper: extract step run code from yaml
step_run_code() {
    local step="$1"
    yq eval '.spec.run' "$(step_file "$step")" 2>/dev/null
}

# Helper: mock __kfg_log functions
mock_kfg_log() {
    cat <<'EOF'
__kfg_log_info() { :; }
__kfg_log_warn() { :; }
__kfg_log_error() { :; }
__kfg_log_debug() { :; }
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
    code=$(printf '%s\n%s\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test' "$(mock_kfg_log)" "$(mock_artifact_tracking)" "$(step_run_code "$step")")
    DIRECTORIES="$dirs" bash -c "$code"
}

# Helper: run step in sandbox directory with env vars
run_step_with_env() {
    local step="$1"
    local code
    code=$(printf '%s\n%s\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test' "$(mock_kfg_log)" "$(mock_artifact_tracking)" "$(step_run_code "$step")")
    bash -c "$code"
}

# Helper: run step with artifact tracking mock
run_step_with_artifacts() {
    local step="$1"
    local code
    code=$(printf '%s\n%s\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test' "$(mock_kfg_log)" "$(mock_artifact_tracking)" "$(step_run_code "$step")")
    bash -c "$code"
}

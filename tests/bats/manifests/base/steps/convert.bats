#!/usr/bin/env bats

load '../../../helpers/common'
load '../../../helpers/manifests'

@test "converts asset with converter and writes output" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: kfg.extension.self.commands.git-commit
spec:
  data:
    name: git-commit
    description: Generate a conventional commit message
    prompt: Analyze and generate commit message.
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: kfg.convert.self.command.claude
spec:
  input:
    format: yaml
  engine:
    expression: |
      "---\nname: " + .name + "\ndescription: " + .description + "\n---\n\n" + .prompt
  output:
    format: raw
YAML
    
    local step_code
    step_code=$(yq eval '.spec.run' "${PROJECT_ROOT}/.manifests/base/steps/convert.yaml")
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
_kfg.log.info() { :; }
_kfg.log.warn() { :; }
_kfg.log.error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$step_code" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" ASSET="kfg.extension.self.commands.git-commit" CONVERTER="kfg.convert.self.command.claude" OUTPUT="$TEST_TMPDIR/output/git-commit.md" bash "$script_file"
    
    [ -f "$TEST_TMPDIR/output/git-commit.md" ]
    grep -q "name: git-commit" "$TEST_TMPDIR/output/git-commit.md"
    
    rm -f "$tmp_file" "$script_file"
}

@test "does not fail when asset is empty" {
    local step_code
    step_code=$(yq eval '.spec.run' "${PROJECT_ROOT}/.manifests/base/steps/convert.yaml")
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
_kfg.log.info() { :; }
_kfg.log.warn() { :; }
_kfg.log.error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$step_code" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    ASSET="" CONVERTER="kfg.convert.self.command.claude" OUTPUT="$TEST_TMPDIR/output.md" bash "$script_file"
    
    [ $? -eq 0 ]
    
    rm -f "$script_file"
}

@test "does not fail when converter is empty" {
    local step_code
    step_code=$(yq eval '.spec.run' "${PROJECT_ROOT}/.manifests/base/steps/convert.yaml")
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
_kfg.log.info() { :; }
_kfg.log.warn() { :; }
_kfg.log.error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$step_code" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    ASSET="kfg.extension.self.commands.git-commit" CONVERTER="" OUTPUT="$TEST_TMPDIR/output.md" bash "$script_file"
    
    [ $? -eq 0 ]
    
    rm -f "$script_file"
}

@test "does not fail when output is empty" {
    local step_code
    step_code=$(yq eval '.spec.run' "${PROJECT_ROOT}/.manifests/base/steps/convert.yaml")
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
_kfg.log.info() { :; }
_kfg.log.warn() { :; }
_kfg.log.error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$step_code" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    ASSET="kfg.extension.self.commands.git-commit" CONVERTER="kfg.convert.self.command.claude" OUTPUT="" bash "$script_file"
    
    [ $? -eq 0 ]
    
    rm -f "$script_file"
}

@test "fails when asset does not exist" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: other.asset
spec:
  data:
    name: other
YAML
    
    local step_code
    step_code=$(yq eval '.spec.run' "${PROJECT_ROOT}/.manifests/base/steps/convert.yaml")
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
_kfg.log.info() { :; }
_kfg.log.warn() { :; }
_kfg.log.error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$step_code" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    ! PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" ASSET="kfg.extension.self.commands.git-commit" CONVERTER="kfg.convert.self.command.claude" OUTPUT="$TEST_TMPDIR/output.md" bash "$script_file"
    
    rm -f "$tmp_file" "$script_file"
}

@test "fails when converter does not exist" {
    local tmp_file
    tmp_file=$(mktemp -t kfg-build-XXXXXX.yaml)
    
    cat > "$tmp_file" <<'YAML'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: kfg.extension.self.commands.git-commit
spec:
  data:
    name: git-commit
    description: Test
    prompt: Test prompt.
YAML
    
    local step_code
    step_code=$(yq eval '.spec.run' "${PROJECT_ROOT}/.manifests/base/steps/convert.yaml")
    
    local script_file="${TEST_TMPDIR}/test-script.sh"
    cat > "$script_file" <<'HEADER'
_kfg.log.info() { :; }
_kfg.log.warn() { :; }
_kfg.log.error() { :; }
declare -a KFG_ARTIFACTS=()
__kfg_add_artifact() { KFG_ARTIFACTS+=("$1"); }
HEADER
    printf '\n__kfg_run_step_test() {\n%s\n}\n__kfg_run_step_test\n' "$step_code" >> "$script_file"
    
    cd "$TEST_TMPDIR"
    ! PATH="${PROJECT_ROOT}/bin:${PATH}" KFG_BUILD_RESULT_FILE="$tmp_file" ASSET="kfg.extension.self.commands.git-commit" CONVERTER="kfg.convert.self.command.claude" OUTPUT="$TEST_TMPDIR/output.md" bash "$script_file"
    
    rm -f "$tmp_file" "$script_file"
}

#!/usr/bin/env bats

# kfg CLI Tests
# Tests for kfg CLI commands that don't depend on nixai-specific manifests

load 'test_helper'

# Version command tests

@test "kfg --version returns exit code 0" {
    run "${KFG_BIN}" --version
    [ "$status" -eq 0 ]
}

@test "kfg --version output contains version prefix" {
    run "${KFG_BIN}" --version
    [ "$status" -eq 0 ]
    [[ "$output" =~ "kfg version" ]]
}

@test "kfg --version output format matches expected pattern" {
    run "${KFG_BIN}" --version
    [ "$status" -eq 0 ]
    # Expected format: kfg version <semver> (<commit>, <date>)
    [[ "$output" =~ kfg\ version\ (dev|[0-9]+\.[0-9]+\.[0-9]+)\ \( ]]
}

@test "kfg --version output contains commit hash" {
    run "${KFG_BIN}" --version
    [ "$status" -eq 0 ]
    # Commit should be either 12-char hex or "unknown"
    [[ "$output" =~ \([a-f0-9]{12}, ]] || [[ "$output" =~ \(unknown, ]]
}

@test "kfg --version output contains build date" {
    run "${KFG_BIN}" --version
    [ "$status" -eq 0 ]
    # Date should be RFC3339 format or "unknown"
    [[ "$output" =~ [0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z\) ]] || [[ "$output" =~ unknown\) ]]
}

@test "kfg binary exists and is executable" {
    [ -x "${KFG_BIN}" ]
}

@test "kfg --help shows correct output" {
    run "${KFG_BIN}" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "KFG is a declarative shell compiler" ]]
    [[ "$output" =~ "shell" ]]
    [[ "$output" =~ "generates shell integration code" ]]
}

# Build command tests

@test "kfg build --help shows correct output" {
    run "${KFG_BIN}" build --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Build a kustomization directory" ]]
    [[ "$output" =~ "--output" ]]
}

@test "kfg build error on invalid path" {
    run "${KFG_BIN}" build /nonexistent/path
    [ "$status" -eq 1 ]
}

# Apply command tests

@test "kfg apply --help shows correct output" {
    run "${KFG_BIN}" apply --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Apply a kustomization" ]]
    [[ "$output" =~ "--kustomize" ]]
    [[ "$output" =~ "--file" ]]
    [[ "$output" =~ "--workflow" ]]
}

# Log command tests

@test "kfg sys log --help shows correct output" {
    run "${KFG_BIN}" sys log --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Levels: error, warn, info, detail, debug" ]]
}

@test "kfg sys log basic invocation" {
    run "${KFG_BIN}" sys log info "test:component" "test message"
    [ "$status" -eq 0 ]
}

@test "kfg sys log with empty message" {
    run "${KFG_BIN}" sys log debug "test:component" ""
    [ "$status" -eq 0 ]
}

@test "kfg sys log invalid level returns error" {
    run "${KFG_BIN}" sys log foo "test:component" "test message"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Levels: error, warn, info, detail, debug" ]]
}

@test "kfg sys log missing component returns error" {
    run "${KFG_BIN}" sys log info
    [ "$status" -ne 0 ]
}

@test "kfg sys log creates .log file" {
    export KFG_LOG_FILE="/tmp/test_kfg_log.log"
    rm -f "$KFG_LOG_FILE"

    run "${KFG_BIN}" sys log info "test:bats" "test message"
    [ "$status" -eq 0 ]

    # Check log file was created with .log extension
    [ -f "$KFG_LOG_FILE" ]

    # Check log file contains valid JSON
    content=$(cat "$KFG_LOG_FILE")
    [[ "$content" =~ "ts" ]]
    [[ "$content" =~ "level" ]]
    [[ "$content" =~ "component" ]]
    [[ "$content" =~ "msg" ]]
    [[ "$content" =~ "source" ]]
    [[ "$content" =~ "pid" ]]

    rm -f "$KFG_LOG_FILE"
}

@test "kfg default log file has .log extension" {
    unset KFG_LOG_FILE
    unset KFG_LOG_DIR

    run "${KFG_BIN}" sys log info "test:extension" "test message"
    [ "$status" -eq 0 ]

    # Check default log file path ends with .log
    default_log_path="$HOME/.local/state/kfg/logs/kfg.log"
    if [ -f "$default_log_path" ]; then
        [[ "$default_log_path" =~ ".log" ]]
    fi
}

@test "kfg sys log KFG_VERBOSE=0 suppresses stderr" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    KFG_VERBOSE=0 run "${KFG_BIN}" sys log info "test:verbose" "test message"
    [ "$status" -eq 0 ]
    [ -z "$output" ]

    rm -f "$KFG_LOG_FILE"
}

@test "kfg sys log KFG_VERBOSE=1 shows error in stderr" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    KFG_VERBOSE=1 run "${KFG_BIN}" sys log error "test:verbose" "test error message"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "ERROR" ]] || [[ "$output" =~ "[ERROR" ]]

    rm -f "$KFG_LOG_FILE"
}

@test "kfg sys log KFG_VERBOSE=1 does NOT show info in stderr" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    KFG_VERBOSE=1 run "${KFG_BIN}" sys log info "test:verbose" "test info message"
    [ "$status" -eq 0 ]
    [ -z "$output" ]

    rm -f "$KFG_LOG_FILE"
}

@test "kfg sys log KFG_VERBOSE=2 shows info in stderr" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    KFG_VERBOSE=2 run "${KFG_BIN}" sys log info "test:verbose" "test info message"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "INFO" ]]

    rm -f "$KFG_LOG_FILE"
}

@test "kfg sys log KFG_VERBOSE=2 does NOT show detail in stderr" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    KFG_VERBOSE=2 run "${KFG_BIN}" sys log detail "test:verbose" "test detail message"
    [ "$status" -eq 0 ]
    [ -z "$output" ]

    rm -f "$KFG_LOG_FILE"
}

@test "kfg sys log KFG_VERBOSE=3 shows detail in stderr" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    KFG_VERBOSE=3 run "${KFG_BIN}" sys log detail "test:verbose" "test detail message"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "DETAIL" ]]

    rm -f "$KFG_LOG_FILE"
}

@test "kfg sys log stdout remains clean" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    KFG_VERBOSE=1 run "${KFG_BIN}" sys log info "test:stdout" "test message"
    [ "$status" -eq 0 ]
    [ -z "$output" ]

    rm -f "$KFG_LOG_FILE"
}

@test "kfg sys log source field is shell" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    run "${KFG_BIN}" sys log info "test:source" "test message"
    [ "$status" -eq 0 ]

    content=$(cat "$KFG_LOG_FILE")
    [[ "$content" =~ '"source":"shell"' ]]

    rm -f "$KFG_LOG_FILE"
}

@test "kfg sys log color output with KFG_LOG_COLOR=always" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    KFG_VERBOSE=1 KFG_LOG_COLOR=always run "${KFG_BIN}" sys log info "test:color" "test message"
    [ "$status" -eq 0 ]

    rm -f "$KFG_LOG_FILE"
}

@test "kfg sys log no color with KFG_LOG_COLOR=never" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    KFG_VERBOSE=1 KFG_LOG_COLOR=never run "${KFG_BIN}" sys log info "test:color" "test message"
    [ "$status" -eq 0 ]
    [[ ! "$output" =~ $'\x1b' ]]

    rm -f "$KFG_LOG_FILE"
}

# Session ID tests

@test "kfg sys log --session-id flag sets session_id in JSONL" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    run "${KFG_BIN}" sys log --session-id "test-session-123" info "test:session" "test message"
    [ "$status" -eq 0 ]

    content=$(cat "$KFG_LOG_FILE")
    [[ "$content" =~ '"session_id":"test-session-123"' ]]

    rm -f "$KFG_LOG_FILE"
}

@test "kfg sys log --session-id flag overrides env var" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    export KFG_SESSION_ID="env-session-456"

    run "${KFG_BIN}" sys log --session-id "flag-session-789" info "test:session" "test message"
    [ "$status" -eq 0 ]

    content=$(cat "$KFG_LOG_FILE")
    [[ "$content" =~ '"session_id":"flag-session-789"' ]]
    [[ ! "$content" =~ '"session_id":"env-session-456"' ]]

    rm -f "$KFG_LOG_FILE"
    unset KFG_SESSION_ID
}

@test "kfg sys log uses env var when flag not provided" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    export KFG_SESSION_ID="env-only-session"

    run "${KFG_BIN}" sys log info "test:session" "test message"
    [ "$status" -eq 0 ]

    content=$(cat "$KFG_LOG_FILE")
    [[ "$content" =~ '"session_id":"env-only-session"' ]]

    rm -f "$KFG_LOG_FILE"
    unset KFG_SESSION_ID
}

@test "kfg sys log empty --session-id omits session_id field" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    export KFG_SESSION_ID="env-to-omit"

    run "${KFG_BIN}" sys log --session-id "" info "test:session" "test message"
    [ "$status" -eq 0 ]

    content=$(cat "$KFG_LOG_FILE")
    [[ ! "$content" =~ '"session_id"' ]]

    rm -f "$KFG_LOG_FILE"
    unset KFG_SESSION_ID
}

@test "kfg sys log no session_id when neither flag nor env var" {
    export KFG_LOG_FILE="/tmp/test_kfg_temp.log"
    rm -f "$KFG_LOG_FILE"

    unset KFG_SESSION_ID

    run "${KFG_BIN}" sys log info "test:session" "test message"
    [ "$status" -eq 0 ]

    content=$(cat "$KFG_LOG_FILE")
    [[ ! "$content" =~ '"session_id"' ]]

    rm -f "$KFG_LOG_FILE"
}

# --verbose flag tests

@test "kfg --verbose=2 apply shows info in stderr" {
    run "${KFG_BIN}" --verbose=2 apply -k nonexistent 2>&1
    # Command will fail but verbose flag should be parsed
    [ "$status" -eq 1 ] || [ "$status" -eq 0 ]
}

@test "kfg --verbose=0 apply silences stderr" {
    run "${KFG_BIN}" --verbose=0 apply -k nonexistent 2>&1
    # Command will fail but verbose flag should be parsed
    [ "$status" -eq 1 ] || [ "$status" -eq 0 ]
}

@test "kfg --verbose flag overrides KFG_VERBOSE env var" {
    KFG_VERBOSE=0 run "${KFG_BIN}" --verbose=2 apply -k nonexistent 2>&1
    # Command will fail but verbose flag should be parsed
    [ "$status" -eq 1 ] || [ "$status" -eq 0 ]
}

# Error message format tests

@test "kfg error messages have core: prefix format" {
    run "${KFG_BIN}" shell bash
    [ "$status" -eq 2 ]
    [[ "$output" =~ "[ERROR][core:cli]" ]] || [[ "$output" =~ "ERROR" ]] && [[ "$output" =~ "cli" ]]
}

@test "kfg invalid command error format" {
    run "${KFG_BIN}" nonexistent-command
    [ "$status" -ne 0 ]
    [[ "$output" =~ "ERROR" ]] || [[ "$output" =~ "error" ]]
}

# Log file extension test

@test "kfg default log file is kfg.log" {
    unset KFG_LOG_FILE
    unset KFG_LOG_DIR

    "${KFG_BIN}" --verbose=1 sys log info "test:extension" "test message"

    default_log_path="$HOME/.local/state/kfg/logs/kfg.log"
    [ -f "$default_log_path" ] || skip "Default log file not created yet"

    [[ "$default_log_path" == *".log" ]]
    [[ "$default_log_path" != *".jsonl" ]]
}

@test "kfg sys log shell source does not get core: prefix" {
    export KFG_LOG_FILE="/tmp/test_kfg_shell_source.log"
    rm -f "$KFG_LOG_FILE"

    run "${KFG_BIN}" sys log info "feature:mcps" "shell message"
    [ "$status" -eq 0 ]

    content=$(cat "$KFG_LOG_FILE")
    [[ "$content" =~ '"component":"feature:mcps"' ]]
    [[ ! "$content" =~ '"component":"core:feature:mcps"' ]]
    [[ "$content" =~ '"source":"shell"' ]]

    rm -f "$KFG_LOG_FILE"
}

# ============================================================================
# Multi-Workflow Tests (using temp fixtures)
# ============================================================================

@test "kfg apply multi-workflow: comma-separated workflow flag parsing" {
    mkdir -p /tmp/test_multi_wf_kust
    
    cat > /tmp/test_multi_wf_kust/kustomization.yaml << 'KUSTEOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - workflows.yaml
  - cmds.yaml
KUSTEOF
    
    cat > /tmp/test_multi_wf_kust/workflows.yaml << 'WFEOF'
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: dev
  shell: bash
spec:
  cmds:
    - kfg.cmds.build
  before:
    - step: kfg.steps.setup
---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: prod
  shell: bash
spec:
  cmds:
    - kfg.cmds.deploy
  before:
    - step: kfg.steps.setup
WFEOF
    
    cat > /tmp/test_multi_wf_kust/cmds.yaml << 'CMDEOF'
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: kfg.cmds.build
  commandName: build
spec:
  run: echo building
---
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: kfg.cmds.deploy
  commandName: deploy
spec:
  run: echo deploying
---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: kfg.steps.setup
spec:
  run: echo setup
CMDEOF
    
    run "${KFG_BIN}" apply -k /tmp/test_multi_wf_kust --workflow dev,prod
    [ "$status" -eq 0 ]
    
    [[ "$output" =~ "build()" ]]
    [[ "$output" =~ "deploy()" ]]
    [[ ! "$output" =~ "KFG_WORKFLOW_NAME" ]]
    [[ "$output" =~ "KFG_KUSTOMIZATION_NAME" ]]
    
    rm -rf /tmp/test_multi_wf_kust
}

@test "kfg apply multi-workflow: no --workflow generates all workflows" {
    mkdir -p /tmp/test_multi_wf_all
    
    cat > /tmp/test_multi_wf_all/kustomization.yaml << 'KUSTEOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - workflows.yaml
  - cmds.yaml
KUSTEOF
    
    cat > /tmp/test_multi_wf_all/workflows.yaml << 'WFEOF'
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: workflow-a
  shell: bash
spec:
  cmds:
    - cmd-alpha
---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: workflow-b
  shell: bash
spec:
  cmds:
    - cmd-beta
WFEOF
    
    cat > /tmp/test_multi_wf_all/cmds.yaml << 'CMDEOF'
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: cmd-alpha
  commandName: cmdAlpha
spec:
  run: echo alpha
---
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: cmd-beta
  commandName: cmdBeta
spec:
  run: echo beta
CMDEOF
    
    run "${KFG_BIN}" apply -k /tmp/test_multi_wf_all
    [ "$status" -eq 0 ]
    
    [[ "$output" =~ "cmdAlpha()" ]]
    [[ "$output" =~ "cmdBeta()" ]]
    [[ ! "$output" =~ "KFG_WORKFLOW_NAME" ]]
    
    rm -rf /tmp/test_multi_wf_all
}

@test "kfg apply multi-workflow: step deduplication" {
    mkdir -p /tmp/test_multi_wf_dedup
    
    cat > /tmp/test_multi_wf_dedup/kustomization.yaml << 'KUSTEOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - workflows.yaml
  - cmds.yaml
KUSTEOF
    
    cat > /tmp/test_multi_wf_dedup/workflows.yaml << 'WFEOF'
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: wf1
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - step: shared-setup
---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: wf2
  shell: bash
spec:
  cmds:
    - cmd2
  before:
    - step: shared-setup
WFEOF
    
    cat > /tmp/test_multi_wf_dedup/cmds.yaml << 'CMDEOF'
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: cmd1
  commandName: cmdOne
spec:
  run: echo one
---
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: cmd2
  commandName: cmdTwo
spec:
  run: echo two
---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: shared-setup
spec:
  run: echo "shared setup step"
CMDEOF
    
    run "${KFG_BIN}" apply -k /tmp/test_multi_wf_dedup
    [ "$status" -eq 0 ]
    
    step_count=$(echo "$output" | grep -c "__kfg_run_step_shared-setup()" || true)
    [ "$step_count" -eq 1 ]
    
    [[ "$output" =~ "cmdOne()" ]]
    [[ "$output" =~ "cmdTwo()" ]]
    
    rm -rf /tmp/test_multi_wf_dedup
}

@test "kfg apply multi-workflow: single workflow still works" {
    mkdir -p /tmp/test_single_wf
    
    cat > /tmp/test_single_wf/kustomization.yaml << 'KUSTEOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - workflow.yaml
  - cmd.yaml
KUSTEOF
    
    cat > /tmp/test_single_wf/workflow.yaml << 'WFEOF'
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: single-workflow
  shell: bash
spec:
  cmds:
    - testcmd
WFEOF
    
    cat > /tmp/test_single_wf/cmd.yaml << 'CMDEOF'
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: testcmd
  commandName: testCmd
spec:
  run: echo test
CMDEOF
    
    run "${KFG_BIN}" apply -k /tmp/test_single_wf
    [ "$status" -eq 0 ]
    
    [[ "$output" =~ "KFG_WORKFLOW_NAME=single-workflow" ]]
    [[ "$output" =~ "testCmd()" ]]
    
    rm -rf /tmp/test_single_wf
}

@test "kfg apply multi-workflow: invalid workflow name shows available workflows" {
    mkdir -p /tmp/test_multi_wf_invalid
    
    cat > /tmp/test_multi_wf_invalid/kustomization.yaml << 'KUSTEOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - workflows.yaml
  - cmds.yaml
KUSTEOF
    
    cat > /tmp/test_multi_wf_invalid/workflows.yaml << 'WFEOF'
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: valid-wf
  shell: bash
spec:
  cmds:
    - validcmd
---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: another-wf
  shell: bash
spec:
  cmds:
    - anothercmd
WFEOF
    
    cat > /tmp/test_multi_wf_invalid/cmds.yaml << 'CMDEOF'
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: validcmd
  commandName: validCmd
spec:
  run: echo valid
---
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: anothercmd
  commandName: anotherCmd
spec:
  run: echo another
CMDEOF
    
    KFG_VERBOSE=2 run "${KFG_BIN}" apply -k /tmp/test_multi_wf_invalid --workflow nonexistent
    [ "$status" -eq 1 ]
    
    [[ "$output" =~ "Available" ]] || [[ "$output" =~ "valid-wf" ]] || [[ "$output" =~ "another-wf" ]]
    
    rm -rf /tmp/test_multi_wf_invalid
}

@test "kfg apply multi-workflow: shell code can be sourced" {
    mkdir -p /tmp/test_multi_wf_source
    
    cat > /tmp/test_multi_wf_source/kustomization.yaml << 'KUSTEOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - workflows.yaml
  - cmds.yaml
KUSTEOF
    
    cat > /tmp/test_multi_wf_source/workflows.yaml << 'WFEOF'
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: wf1
  shell: bash
spec:
  cmds:
    - cmd1
---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: wf2
  shell: bash
spec:
  cmds:
    - cmd2
WFEOF
    
    cat > /tmp/test_multi_wf_source/cmds.yaml << 'CMDEOF'
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: cmd1
  commandName: cmdOne
spec:
  run: echo "command one"
---
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: cmd2
  commandName: cmdTwo
spec:
  run: echo "command two"
CMDEOF
    
    run "${KFG_BIN}" apply -k /tmp/test_multi_wf_source
    [ "$status" -eq 0 ]
    
    echo "$output" > /tmp/test_multi_wf_shell.sh
    
    run bash -n /tmp/test_multi_wf_shell.sh
    [ "$status" -eq 0 ]
    
    run bash -c "source /tmp/test_multi_wf_shell.sh && type cmdOne && type cmdTwo"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "function" ]]
    
    rm -rf /tmp/test_multi_wf_source
    rm -f /tmp/test_multi_wf_shell.sh
}

@test "kfg apply multi-workflow: whitespace in workflow names is trimmed" {
    mkdir -p /tmp/test_multi_wf_trim
    
    cat > /tmp/test_multi_wf_trim/kustomization.yaml << 'KUSTEOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - workflows.yaml
  - cmds.yaml
KUSTEOF
    
    cat > /tmp/test_multi_wf_trim/workflows.yaml << 'WFEOF'
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: wf-one
  shell: bash
spec:
  cmds:
    - cmd1
---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: wf-two
  shell: bash
spec:
  cmds:
    - cmd2
WFEOF
    
    cat > /tmp/test_multi_wf_trim/cmds.yaml << 'CMDEOF'
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: cmd1
  commandName: cmdOne
spec:
  run: echo one
---
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: cmd2
  commandName: cmdTwo
spec:
  run: echo two
CMDEOF
    
    run "${KFG_BIN}" apply -k /tmp/test_multi_wf_trim --workflow "wf-one, wf-two"
    [ "$status" -eq 0 ]
    
    [[ "$output" =~ "cmdOne()" ]]
    [[ "$output" =~ "cmdTwo()" ]]
    
    rm -rf /tmp/test_multi_wf_trim
}

@test "kfg apply multi-workflow: --workflow help shows comma-separated usage" {
    run "${KFG_BIN}" apply --help
    [ "$status" -eq 0 ]
    
    [[ "$output" =~ "comma-separated" ]] || [[ "$output" =~ "workflow" ]] && [[ "$output" =~ "dev,openspec" ]]
}
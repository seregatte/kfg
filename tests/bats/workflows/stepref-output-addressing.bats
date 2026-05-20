#!/usr/bin/env bats

# Test stepref-output-addressing feature:
# - Required name on StepReference
# - Runtime execution identity for outputs
# - $kfg.output(...) expansion in env

load '../test_helper'

# Test fixtures directory
FIXTURES_DIR="${PROJECT_ROOT}/tests/bats/fixtures/stepref-output-addressing"

setup() {
    # Create fixtures directory if it doesn't exist
    mkdir -p "$FIXTURES_DIR"
    
    # Create a simple step with output
    cat > "$FIXTURES_DIR/resources.yaml" << 'EOF'
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: test.detect
spec:
  run: echo "claude"
  output:
    name: AGENT
    type: string
---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: test.setup
spec:
  run: echo "setup done"
---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test.workflow
  shell: bash
spec:
  cmds:
    - test.cmd
  before:
    - name: detect-agent
      step: test.detect
      weight: -90
    - name: setup-claude
      step: test.setup
      weight: -80
      when:
        output:
          step: detect-agent
          name: AGENT
          equals: "claude"
---
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: test.cmd
  commandName: testcmd
spec:
  run: echo "test cmd"
EOF

    # Create kustomization.yaml that references resources
    cat > "$FIXTURES_DIR/kustomization.yaml" << 'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - resources.yaml
EOF
}

teardown() {
    rm -rf "$FIXTURES_DIR"
}

@test "workflow with named step references generates shell code" {
    # Generate shell code from the workflow
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow test.workflow
    [ "$status" -eq 0 ]
}

@test "step functions have runtime identity parameter" {
    # Generate shell code
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow test.workflow
    [ "$status" -eq 0 ]
    
    # Verify step function has __step_ref_name parameter
    [[ "$output" =~ "__step_ref_name" ]]
}

@test "output is stored under step reference name" {
    # Generate shell code
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow test.workflow
    [ "$status" -eq 0 ]
    
    # Verify output is stored using __kfg_output_set with step_ref_name
    [[ "$output" =~ "__kfg_output_set" ]]
}

@test "when condition uses step reference name" {
    # Generate shell code
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow test.workflow
    [ "$status" -eq 0 ]
    
    # Verify when condition uses detect-agent (step reference name)
    [[ "$output" =~ "detect-agent" ]]
}

@test "generated shell code has valid syntax" {
    # Generate shell code
    run "${KFG_BIN}" apply -k "$FIXTURES_DIR" --workflow test.workflow
    [ "$status" -eq 0 ]
    
    # Verify shell syntax
    verify_shell_syntax "$output"
}
#!/usr/bin/env bats

# kfg ai Command Tests
# Tests for the AI wizard CLI command

load '../test_helper'

@test "kfg ai --help succeeds" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
}

@test "kfg ai --help mentions KFG_AI_AGENT environment variable" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "KFG_AI_AGENT" ]]
}

@test "kfg ai --help mentions default agent (pi)" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "default: pi" ]]
}

@test "kfg ai --help mentions supported agents" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "pi" ]]
    [[ "$output" =~ "opencode" ]]
}

@test "kfg ai --help includes usage examples" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "kfg ai" ]]
    [[ "$output" =~ "Examples:" ]]
}

@test "kfg ai --help describes wizard purpose" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "interactive" ]]
    [[ "$output" =~ "configuration" ]]
}

@test "kfg ai --help explains argument forwarding" {
    run "${KFG_BIN}" ai --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "forwarded" ]] || [[ "$output" =~ "--" ]]
}

@test "kfg ai uses local overlay from current directory" {
    # Create minimal valid AI overlay in temp directory
    local overlay_dir="${TEST_TMPDIR}/packages/domains/ai-agents/overlays/ai"
    mkdir -p "$overlay_dir"

    cat > "${overlay_dir}/kustomization.yaml" <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - manifest.yaml
EOF

    cat > "${overlay_dir}/manifest.yaml" <<'EOF'
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: test.cmd.pi
  commandName: pi
spec:
  run: |
    echo "AI_WIZARD_OK"
---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: kfg.ai.workflow
  shell: bash
spec:
  cmds:
    - test.cmd.pi
EOF

    # Run from TEST_TMPDIR so resolveAIOverlay finds the local overlay
    cd "${TEST_TMPDIR}"
    run "${KFG_BIN}" ai
    [ "$status" -eq 0 ]
    [[ "$output" =~ "AI_WIZARD_OK" ]]
}

@test "kfg ai ignores KFG_KPATH when local overlay exists" {
    local overlay_dir="${TEST_TMPDIR}/packages/domains/ai-agents/overlays/ai"
    mkdir -p "$overlay_dir"

    cat > "${overlay_dir}/kustomization.yaml" <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - manifest.yaml
EOF

    cat > "${overlay_dir}/manifest.yaml" <<'EOF'
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: test.cmd.pi
  commandName: pi
spec:
  run: |
    echo "AI_WIZARD_OK"
---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: kfg.ai.workflow
  shell: bash
spec:
  cmds:
    - test.cmd.pi
EOF

    cd "${TEST_TMPDIR}"
    KFG_KPATH=/unrelated/path run "${KFG_BIN}" ai
    [ "$status" -eq 0 ]
    [[ "$output" =~ "AI_WIZARD_OK" ]]
}

@test "kfg ai real overlay builds with kfg.ai.workflow" {
    run "${KFG_BIN}" run -k packages/domains/ai-agents/overlays/ai
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Available commands" ]]
    [[ "$output" =~ "pi" ]]
    [[ "$output" =~ "opencode" ]]
}

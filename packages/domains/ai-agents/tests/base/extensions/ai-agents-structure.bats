#!/usr/bin/env bats

# Validate the ai-agents extension structure: all expected files exist.

MANIFESTS_BASE="packages/domains/ai-agents/manifests"

@test "ai-agents root kustomization exists" {
    [ -f "$MANIFESTS_BASE/kustomization.yaml" ]
}

@test "ai-agents agents directory exists" {
    [ -d "$MANIFESTS_BASE/agents" ]
    [ -f "$MANIFESTS_BASE/agents/kustomization.yaml" ]
}

@test "ai-agents cmds directory exists" {
    [ -d "$MANIFESTS_BASE/cmds" ]
    [ -f "$MANIFESTS_BASE/cmds/kustomization.yaml" ]
}

@test "ai-agents steps directory exists" {
    [ -d "$MANIFESTS_BASE/steps" ]
    [ -f "$MANIFESTS_BASE/steps/kustomization.yaml" ]
}

@test "ai-agents prompts directory exists" {
    [ -d "$MANIFESTS_BASE/prompts" ]
    [ -f "$MANIFESTS_BASE/prompts/kustomization.yaml" ]
}

@test "ai-agents subagents directory exists" {
    [ -d "$MANIFESTS_BASE/subagents" ]
    [ -f "$MANIFESTS_BASE/subagents/kustomization.yaml" ]
}

@test "ai-agents converters directory exists" {
    [ -d "$MANIFESTS_BASE/converters" ]
    [ -f "$MANIFESTS_BASE/converters/kustomization.yaml" ]
}



@test "opencode agent structure exists" {
    [ -d "$MANIFESTS_BASE/agents/opencode" ]
    [ -f "$MANIFESTS_BASE/agents/opencode/kustomization.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/opencode/assets/kustomization.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/opencode/assets/settings.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/opencode/converters/kustomization.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/opencode/converters/command.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/opencode/converters/mcp.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/opencode/converters/subagent.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/opencode/converters/cfg.yaml" ]
}

@test "pi agent structure exists" {
    [ -d "$MANIFESTS_BASE/agents/pi" ]
    [ -f "$MANIFESTS_BASE/agents/pi/kustomization.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/pi/assets/kustomization.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/pi/assets/settings.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/pi/converters/kustomization.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/pi/converters/command.yaml" ]
}

@test "shared resource files exist" {
    [ -f "$MANIFESTS_BASE/cmds/agents/opencode.yaml" ]
    [ -f "$MANIFESTS_BASE/cmds/agents/pi.yaml" ]
    [ -f "$MANIFESTS_BASE/cmds/openspec/openspec.yaml" ]
    [ -f "$MANIFESTS_BASE/steps/detect/detect.yaml" ]
    [ -f "$MANIFESTS_BASE/prompts/git-commit/git-commit.yaml" ]
    [ -f "$MANIFESTS_BASE/prompts/refactor-pure/refactor-pure.yaml" ]
    [ -f "$MANIFESTS_BASE/prompts/review-code/review-code.yaml" ]
    [ -f "$MANIFESTS_BASE/prompts/review-search/review-search.yaml" ]
    [ -f "$MANIFESTS_BASE/subagents/review-minimal/review-minimal.yaml" ]
    [ -f "$MANIFESTS_BASE/converters/to-json/to-json.yaml" ]
}

@test "old directories are removed" {
    [ ! -d ".manifests/base/agents" ]
    [ ! -d ".manifests/base/cmds" ]
    [ ! -d ".manifests/base/extensions/self" ]
}

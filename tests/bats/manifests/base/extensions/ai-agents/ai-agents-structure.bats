#!/usr/bin/env bats

# Validate the ai-agents extension structure: all expected files exist.

MANIFESTS_BASE=".manifests/base/extensions/ai-agents"

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

@test "claude agent structure exists" {
    [ -d "$MANIFESTS_BASE/agents/claude" ]
    [ -f "$MANIFESTS_BASE/agents/claude/kustomization.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/claude/assets/kustomization.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/claude/assets/settings.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/claude/converters/kustomization.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/claude/converters/command.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/claude/converters/mcp.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/claude/converters/subagent.yaml" ]
}

@test "gemini agent structure exists" {
    [ -d "$MANIFESTS_BASE/agents/gemini" ]
    [ -f "$MANIFESTS_BASE/agents/gemini/kustomization.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/gemini/assets/kustomization.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/gemini/assets/settings.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/gemini/converters/kustomization.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/gemini/converters/command.yaml" ]
    [ -f "$MANIFESTS_BASE/agents/gemini/converters/mcp.yaml" ]
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
    [ -f "$MANIFESTS_BASE/cmds/agents.yaml" ]
    [ -f "$MANIFESTS_BASE/cmds/openspec.yaml" ]
    [ -f "$MANIFESTS_BASE/steps/detect.yaml" ]
    [ -f "$MANIFESTS_BASE/prompts/git-commit.yaml" ]
    [ -f "$MANIFESTS_BASE/prompts/refactor-pure.yaml" ]
    [ -f "$MANIFESTS_BASE/prompts/review-code.yaml" ]
    [ -f "$MANIFESTS_BASE/prompts/review-search.yaml" ]
    [ -f "$MANIFESTS_BASE/subagents/review-minimal.yaml" ]
    [ -f "$MANIFESTS_BASE/converters/to-json.yaml" ]
}

@test "old directories are removed" {
    [ ! -d ".manifests/base/agents" ]
    [ ! -d ".manifests/base/cmds" ]
    [ ! -d ".manifests/base/extensions/self" ]
}

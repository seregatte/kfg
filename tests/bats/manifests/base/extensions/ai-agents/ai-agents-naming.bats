#!/usr/bin/env bats

# Validate all resource metadata.name fields follow the new naming convention.

MANIFESTS_BASE=".manifests/base/extensions/ai"

_extract_name() {
    grep '^\s*name:' "$1" | head -1 | sed 's/.*name:\s*//' | tr -d ' '
}

# --- Per-agent assets ---

@test "claude asset settings name is ai.claude.asset.settings" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/claude/assets/settings.yaml")
    [ "$name" = "ai.claude.asset.settings" ]
}

@test "gemini asset settings name is ai.gemini.asset.settings" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/gemini/assets/settings.yaml")
    [ "$name" = "ai.gemini.asset.settings" ]
}

@test "opencode asset settings name is ai.opencode.asset.settings" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/opencode/assets/settings.yaml")
    [ "$name" = "ai.opencode.asset.settings" ]
}

@test "pi asset settings name is ai.pi.asset.settings" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/pi/assets/settings.yaml")
    [ "$name" = "ai.pi.asset.settings" ]
}

# --- Per-agent command converters ---

@test "claude command converter name is ai.claude.conv.command" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/claude/converters/command.yaml")
    [ "$name" = "ai.claude.conv.command" ]
}

@test "gemini command converter name is ai.gemini.conv.command" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/gemini/converters/command.yaml")
    [ "$name" = "ai.gemini.conv.command" ]
}

@test "opencode command converter name is ai.opencode.conv.command" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/opencode/converters/command.yaml")
    [ "$name" = "ai.opencode.conv.command" ]
}

@test "pi command converter name is ai.pi.conv.command" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/pi/converters/command.yaml")
    [ "$name" = "ai.pi.conv.command" ]
}

# --- Per-agent MCP converters ---

@test "claude MCP converter name is ai.claude.conv.mcp" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/claude/converters/mcp.yaml")
    [ "$name" = "ai.claude.conv.mcp" ]
}

@test "gemini MCP converter name is ai.gemini.conv.mcp" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/gemini/converters/mcp.yaml")
    [ "$name" = "ai.gemini.conv.mcp" ]
}

@test "opencode MCP converter name is ai.opencode.conv.mcp" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/opencode/converters/mcp.yaml")
    [ "$name" = "ai.opencode.conv.mcp" ]
}

# --- Per-agent subagent converters ---

@test "claude subagent converter name is ai.claude.conv.subagent" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/claude/converters/subagent.yaml")
    [ "$name" = "ai.claude.conv.subagent" ]
}

@test "opencode subagent converter name is ai.opencode.conv.subagent" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/opencode/converters/subagent.yaml")
    [ "$name" = "ai.opencode.conv.subagent" ]
}

# --- Per-agent config converter ---

@test "opencode config converter name is ai.opencode.conv.cfg" {
    name=$(_extract_name "$MANIFESTS_BASE/agents/opencode/converters/cfg.yaml")
    [ "$name" = "ai.opencode.conv.cfg" ]
}

# --- Shared converters ---

@test "shared to-json converter name is ai.conv.to-json" {
    name=$(_extract_name "$MANIFESTS_BASE/converters/to-json.yaml")
    [ "$name" = "ai.conv.to-json" ]
}

# --- Shared cmds ---

@test "agents.yaml contains ai.*.cmd.main names" {
    grep -q 'name: ai.claude.cmd.main' "$MANIFESTS_BASE/cmds/agents.yaml"
    grep -q 'name: ai.gemini.cmd.main' "$MANIFESTS_BASE/cmds/agents.yaml"
    grep -q 'name: ai.opencode.cmd.main' "$MANIFESTS_BASE/cmds/agents.yaml"
    grep -q 'name: ai.pi.cmd.main' "$MANIFESTS_BASE/cmds/agents.yaml"
}

@test "openspec cmd name is ai.cmds.openspec" {
    name=$(_extract_name "$MANIFESTS_BASE/cmds/openspec.yaml")
    [ "$name" = "ai.cmds.openspec" ]
}

# --- Shared steps ---

@test "detect step name is ai.steps.detect" {
    name=$(_extract_name "$MANIFESTS_BASE/steps/detect.yaml")
    [ "$name" = "ai.steps.detect" ]
}

# --- Shared prompts ---

@test "git-commit prompt name is ai.prompts.git-commit" {
    name=$(_extract_name "$MANIFESTS_BASE/prompts/git-commit.yaml")
    [ "$name" = "ai.prompts.git-commit" ]
}

@test "refactor-pure prompt name is ai.prompts.refactor-pure" {
    name=$(_extract_name "$MANIFESTS_BASE/prompts/refactor-pure.yaml")
    [ "$name" = "ai.prompts.refactor-pure" ]
}

@test "review-code prompt name is ai.prompts.review-code" {
    name=$(_extract_name "$MANIFESTS_BASE/prompts/review-code.yaml")
    [ "$name" = "ai.prompts.review-code" ]
}

@test "review-search prompt name is ai.prompts.review-search" {
    name=$(_extract_name "$MANIFESTS_BASE/prompts/review-search.yaml")
    [ "$name" = "ai.prompts.review-search" ]
}

# --- Shared subagents ---

@test "review-minimal subagent name is ai.subagents.review-minimal" {
    name=$(_extract_name "$MANIFESTS_BASE/subagents/review-minimal.yaml")
    [ "$name" = "ai.subagents.review-minimal" ]
}

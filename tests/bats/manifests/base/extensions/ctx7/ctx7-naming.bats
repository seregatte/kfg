#!/usr/bin/env bats

# Validate ctx7 extension resource names follow ctx7.<kind>.<name> convention.

CTX7_EXT=".manifests/base/extensions/ai/ctx7"

_extract_name() {
    grep '^\s*name:' "$1" | head -1 | sed 's/.*name:\s*//' | tr -d ' '
}

@test "ctx7 install step name is ctx7.steps.install" {
    name=$(_extract_name "$CTX7_EXT/steps/install.yaml")
    [ "$name" = "ctx7.steps.install" ]
}

@test "ctx7 inject step name is ctx7.steps.inject" {
    name=$(_extract_name "$CTX7_EXT/steps/inject-ctx7-context.yaml")
    [ "$name" = "ctx7.steps.inject" ]
}

@test "ctx7 MCP asset name is ctx7.assets.mcp" {
    name=$(_extract_name "$CTX7_EXT/assets/mcp.yaml")
    [ "$name" = "ctx7.assets.mcp" ]
}

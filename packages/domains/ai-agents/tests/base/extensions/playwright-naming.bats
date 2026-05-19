#!/usr/bin/env bats

# Validate playwright extension resource names follow playwright.<kind>.<name> convention.

PLAYWRIGHT_EXT="packages/domains/ai-agents/manifests/playwright"

_extract_name() {
    grep '^\s*name:' "$1" | head -1 | sed 's/.*name:\s*//' | tr -d ' '
}

@test "playwright MCP asset name is playwright.assets.mcp" {
    name=$(_extract_name "$PLAYWRIGHT_EXT/assets/mcp.yaml")
    [ "$name" = "playwright.assets.mcp" ]
}

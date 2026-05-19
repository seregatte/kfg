#!/usr/bin/env bats

# Validate chrome-devtools extension resource names follow chrome.<kind>.<name> convention.

CHROME_EXT=".manifests/base/extensions/chrome-devtools"

_extract_name() {
    grep '^\s*name:' "$1" | head -1 | sed 's/.*name:\s*//' | tr -d ' '
}

@test "chrome MCP asset name is chrome.assets.mcp" {
    name=$(_extract_name "$CHROME_EXT/assets/mcp.yaml")
    [ "$name" = "chrome.assets.mcp" ]
}

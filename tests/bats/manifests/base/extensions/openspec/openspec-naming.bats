#!/usr/bin/env bats

# Validate openspec extension resource names follow openspec.<kind>.<name> convention.

OPENSPEC_EXT=".manifests/base/extensions/ai/openspec"

_extract_name() {
    grep '^\s*name:' "$1" | head -1 | sed 's/.*name:\s*//' | tr -d ' '
}

@test "openspec install step name is openspec.steps.install" {
    name=$(_extract_name "$OPENSPEC_EXT/steps/install.yaml")
    [ "$name" = "openspec.steps.install" ]
}

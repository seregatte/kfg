#!/usr/bin/env bats

MANIFESTS_BASE="$BATS_TEST_DIRNAME/../../../manifests"

# List of all public capabilities with expected README locations
CAPABILITIES=(
  "agents/opencode:OpenCode"
  "agents/pi:Pi"
  "chrome-devtools:Chrome DevTools"
  "ctx7:Context7"
  "gws:Google Workspace"
  "notebooklm:NotebookLM"
  "openspec:OpenSpec"
  "playwright:Playwright"
)

REQUIRED_SECTIONS=("Purpose" "Prerequisites" "Entrypoints" "Exported Resources" "Agent Compatibility" "Configuration" "JSON Patch" "Workflow Usage" "Generated Artifacts" "Cache" "Limitations" "Validation" "Canonical Specs")

@test "capability READMEs exist for all public capabilities" {
  for cap in "${CAPABILITIES[@]}"; do
    IFS=':' read -r path name <<< "$cap"
    readme="$MANIFESTS_BASE/$path/README.md"
    [ -f "$readme" ]
  done
}

@test "capability READMEs contain all required sections" {
  for cap in "${CAPABILITIES[@]}"; do
    IFS=':' read -r path name <<< "$cap"
    readme="$MANIFESTS_BASE/$path/README.md"
    if [ -f "$readme" ]; then
      for section in "${REQUIRED_SECTIONS[@]}"; do
        grep -qi "$section" "$readme" || {
          echo "Missing section '$section' in $path/README.md"
          return 1
        }
      done
    fi
  done
}

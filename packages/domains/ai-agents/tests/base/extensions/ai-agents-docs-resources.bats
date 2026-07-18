#!/usr/bin/env bats

MANIFESTS_BASE="$BATS_TEST_DIRNAME/../../../manifests"

@test "documented resource names match built manifests" {
  for readme in "$MANIFESTS_BASE"/agents/*/README.md \
                "$MANIFESTS_BASE"/chrome-devtools/README.md \
                "$MANIFESTS_BASE"/ctx7/README.md \
                "$MANIFESTS_BASE"/gws/README.md \
                "$MANIFESTS_BASE"/notebooklm/README.md \
                "$MANIFESTS_BASE"/openspec/README.md \
                "$MANIFESTS_BASE"/playwright/README.md; do
    if [ -f "$readme" ]; then
      grep -qE '(ai\.[a-z]+\.[a-z]+\.[a-z]+|[a-z]+-?[a-z]+\.assets\.|steps\.)' "$readme" || {
        echo "No resource names found in $readme"
        return 1
      }
    fi
  done
}

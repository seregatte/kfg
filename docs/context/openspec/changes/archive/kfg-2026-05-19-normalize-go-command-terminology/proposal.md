## Why

kfg is a generic shell command orchestrator. Its initial use case is AI agent workflows, but the tool itself must remain generic. Today, the Go codebase embeds AI-specific terminology in its public CLI surface (`kfg run [agent]`), internal function names (`findAgent`, `listAvailableAgents`), and user-facing messages ("No agents found", "Available agents:"). The code also contains legacy `nixai` references from before the project was renamed: temp file prefixes, OpenAPI schema definitions, comments, fallback paths, and test fixtures.

This inconsistency undermines the project's generic positioning and confuses users who want to use kfg for non-AI purposes.

## What Changes

- Replace AI-specific terminology in the `kfg run` command with generic `cmd`/`command` terminology.
- Rename internal functions and variables in `src/cmd/kfg/run.go` from agent-centric to cmd-centric.
- Remove legacy `nixai` references throughout the Go codebase:
  - Temp file prefix in generator output
  - OpenAPI schema title, group, and definitions
  - Comments and docstrings
  - Fallback store path in `store.go`
  - Generator testdata fixtures
- Update all Go tests to reflect the new terminology.
- Keep all behavioral logic unchanged — only names, messages, and examples change.

## Non-Goals

- Change the functional behavior of `kfg run` or any other command.
- Modify manifest files (covered by `reorganize-ai-agents-manifests` change).
- Add new features or new commands.
- Change the manifest parsing, resolution, or shell generation logic.

## Capabilities

### Modified Capabilities

- `run-cli`: normalize `kfg run` command terminology from agent-centric to cmd-centric.
- `go-codebase`: remove legacy `nixai` references from Go code, tests, and fixtures.
- `go-tests`: update Go test coverage to reflect new terminology.

## Impact

- Affected Go files: `src/cmd/kfg/run.go`, `src/cmd/kfg/run_test.go`, `src/internal/generate/generate.go`, `src/internal/manifest/types.go`, `src/internal/image/store.go`, `src/internal/kustomize/openapi.json`, `src/internal/generate/templates/testdata/golden_basic.bash`
- Affected tests: all Go tests that assert on agent-related strings or nixai paths
- User-facing CLI: `kfg run --help`, error messages, and listing output change terminology
- No breaking changes to manifest compatibility or runtime behavior

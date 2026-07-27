## Why

`kfg ai` currently depends on a repository-relative overlay, so installed binaries fail when run outside the kfg source tree. Always using the hosted overlay would fix consumers but would bypass uncommitted manifest changes during local development.

## What Changes

- Make `kfg ai` use `packages/domains/ai-agents/overlays/ai` when that path exists relative to the current working directory.
- Fall back to `https://github.com/seregatte/kfg.git//packages/domains/ai-agents/overlays/ai?ref=main` when the local overlay is absent.
- Keep the selected source explicit through `kfg run -k`, so the generic `KFG_KPATH` setting cannot redirect the wizard.
- Add deterministic tests for local development, installed usage, and `KFG_KPATH` isolation without requiring network access.
- Update the command specification and user-facing documentation to describe source selection.
- Do not add ancestor traversal, a dedicated override variable, manifest model changes, or fallback from an invalid local overlay.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `ai-wizard-command`: Define local-first overlay resolution with an online fallback and independence from `KFG_KPATH`.

## Impact

- Affected engine code: `src/cmd/kfg/ai.go` and focused command tests.
- Affected CLI behavior: `kfg ai` works outside the kfg checkout while continuing to consume local overlay edits when run from a checkout or worktree root.
- Affected documentation: the AI wizard command specification and command documentation that describes its manifest source.
- Domain impact: the existing AI overlay remains the entrypoint; no resources, defaults, compatibility rules, or generated artifacts change.
- Framework impact: no framework steps or runtime contracts change.
- Dependencies: the remote fallback uses the existing Kustomize Git source support and requires network access only when no local overlay is selected.
- Compatibility: existing development usage from the repository root remains unchanged. Installed usage that previously failed gains a working default. Generic `KFG_KPATH` behavior for `build`, `apply`, and `run` remains unchanged.

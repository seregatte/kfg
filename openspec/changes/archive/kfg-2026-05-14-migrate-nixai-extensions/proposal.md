## Why

The `.nixai/` extensions currently use shell scripts (`build.sh`, `generate-*.sh`) and jq templates for skill installation and config generation. This approach:
- Requires `.self/` scripts to be materialized in the workspace at runtime
- Mixes build-time logic (Imagefile) with runtime logic (Steps)
- Makes extensions tightly coupled to the shell-based generation pipeline

The new declarative model (Assets/Converters + Steps) already exists in `.manifests/` for MCP configs and commands. We need to complete the migration by adding skill installation Steps and organizing the manifest structure.

Additionally, `base/kustomization.yaml` lists 49 lines of individual files — it should reference directories instead.

## What Changes

- **New Steps**: Add generic skill installation Steps for 6 extensions (ctx7, chrome-devtools, playwright, gws, notebooklm, openspec)
- **No shell logic in Steps**: All agent-specific behavior resolved via `spec.env` variables — no `case`/`if` in Step code
- **Kustomization reorganization**: Add `kustomization.yaml` to every directory in `manifests/base/` so the root kustomization references only 4 directories
- **No changes to `.nixai/`**: The source of truth remains untouched
- **No changes to existing Assets/Converters**: ctx7, chrome-devtools, playwright MCP assets already migrated

## Capabilities

### New Capabilities

- `skill-installation-steps`: Generic Steps for installing agent skills via external CLIs (ctx7, npx skills, notebooklm, openspec). Each Step receives all agent-specific config via env vars.

### Modified Capabilities

- `manifest-organization`: Add hierarchical kustomization files to `agents/`, `cmds/`, `steps/`, `extensions/` and all subdirectories. Root `kustomization.yaml` simplified to reference 4 directories.

## Impact

- **Files created**: 17 kustomization.yaml files + 6 install Step files = 23 files
- **Files modified**: `base/kustomization.yaml` (simplified from 49 lines to 4 lines)
- **No breaking changes**: Existing workflow references remain valid
- **No dependency changes**: No new Go code or external tools required

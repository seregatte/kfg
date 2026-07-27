## Why

Users need to customize AI agent configurations (skills, commands, MCP settings) without
modifying the generated workflow manifests. Currently, all agent configs are materialized
from KFG manifests — there is no escape hatch for project-specific overrides. The nixai-v1-freeze
project had a concept of merging a local folder into generated configs per agent, and KFG
lacks an equivalent mechanism.

## What Changes

1. **New framework Step**: `kfg.merge-local` — copies/deep-merges files from a `SRC` directory
   into a `DEST` directory, with configurable `MODE` (deep, overwrite, skip) and optional `FILTER`.
2. **Wizard integration**: The `kfg ai` wizard gains a new Phase 3.5 question asking the user
   if they want local customization. If yes, the generated workflow includes the `kfg.merge-local`
   step with the user's chosen source path.
3. **Wizard prompt update**: `wizard.yaml` is extended to handle the local customization
   conversation, detect existing agent configs, and generate the merge step in workflows.

## Impact

- **Framework**: 1 new Step manifest, 1 kustomization update (no Go code changes)
- **Wizard**: 1 prompt update in `wizard.yaml` (the wizard itself is the implementation vehicle)
- **Backward compatible**: The step is opt-in; workflows without it behave identically

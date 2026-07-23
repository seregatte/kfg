## Why

The current extension model is split across incomplete `base/extensions/*` manifests, legacy image-based implementations under `$HOME/.nixai/images/extensions/*`, and a `dev` overlay that still consumes transitional `kfg.extension.self.*` assets. This makes extension behavior hard to reason about, blocks reuse in `lifeos`, and leaves installation and MCP wiring inconsistent across agents.

The change is needed now because `lifeos` is the next concrete consumer of the extension system. Before migrating its overlay, `kfg` needs a single extension contract that preserves the working behavior of the legacy images while moving installation and MCP generation into reusable `Step` and `Assets` resources.

## What Changes

- Normalize the extension model so each extension owns its reusable manifests under `base/extensions/<name>/` and exposes `Step` and `Assets` resources instead of relying on legacy image builds.
- Update extension installation Steps for `chrome-devtools`, `gws`, `notebooklm`, and `playwright` to match the working behavior currently implemented in `$HOME/.nixai/images/extensions/*`, while keeping `$HOME/.nixai` as read-only reference material.
- Standardize MCP-first extensions so canonical MCP assets live in extension namespaces such as `kfg.extension.playwright.mcp` and `kfg.extension.chrome-devtools.mcp`.
- Modify the dev overlay workflow to consume canonical extension namespaces instead of transitional `kfg.extension.self.*` MCP assets where extension-owned assets exist.
- Define how agent-specific install inputs, temporary agent homes, output directories, copied artifacts, and deterministic outputs are handled by extension install Steps.
- Migrate the `lifeos` overlay to the normalized extension pattern, reusing existing `Steps`, `Assets`, and `Converters` from `kfg` wherever possible and modeling LifeOS-specific commands, subagents, and MCPs as local assets.
- **BREAKING** Remove the transitional assumption that legacy image outputs or `kfg.extension.self.*` MCP assets are the canonical extension interface for normalized extensions.

## Capabilities

### New Capabilities
- `extension-mcp-assets`: Define the canonical manifest contract for extension-owned MCP assets and their consumption by overlays.
- `lifeos-overlay-migration`: Define the normalized LifeOS overlay structure and how it composes shared extension Steps, Assets, and Converters.

### Modified Capabilities
- `skill-installation-steps`: Change extension install Step requirements to preserve legacy image behavior while using explicit Step contracts, agent-specific inputs, artifacts, and outputs.
- `dev-workflow`: Change the agent development workflow to consume normalized extension namespaces and phase ordering for extension install, MCP aggregation, and cleanup.
- `ctx7-install`: Align the ctx7 install capability with the normalized extension contract and extension-owned MCP asset usage.

## Impact

- Affected manifests under `.manifests/base/extensions/` for `ctx7`, `openspec`, `chrome-devtools`, `gws`, `notebooklm`, and `playwright`.
- Affected overlay orchestration under `.manifests/overlay/dev/agents-workflow.yaml`.
- New or updated OpenSpec capability specs for extension MCP assets and LifeOS overlay migration.
- Follow-on implementation work in the `lifeos` repository, executed from `~/Sites/lifeos`, but using `$HOME/.nixai` only as read-only behavioral reference.
- Test impact for manifest validation, workflow generation, and integration coverage around extension install and MCP aggregation behavior.

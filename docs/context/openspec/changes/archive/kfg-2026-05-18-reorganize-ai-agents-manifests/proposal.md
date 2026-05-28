## Why

kfg is a generic shell command orchestrator based on YAML manifests. Its initial use case is AI agent workflows, but the tool itself must remain generic in conception. Today, AI-agent-specific resources are fragmented across multiple directories in the base manifests (`agents/`, `cmds/`, `steps/`, `extensions/self/`), making it unclear what is core infrastructure vs. domain-specific. The naming conventions also embed AI-specific terminology (`kfg.agent.*`, `kfg.detect-agent`) into what should be a generic manifest model.

This change consolidates all AI-agent resources into a single `ai-agents` extension, reorganizes per-agent resources into subdirectories, and renames all resource metadata to a short, consistent convention (`ai.<agent>.<kind>.<name>`).

## What Changes

- Consolidate all AI-agent-specific resources into `.manifests/base/extensions/ai-agents/`.
- Reorganize per-agent resources into `agents/<agent>/assets/` and `agents/<agent>/converters/` subdirectories.
- Move shared resources (cmds, steps, prompts, subagents, converters) into dedicated directories under `ai-agents/`.
- Rename all `metadata.name` fields to a new short convention:
  - Per-agent: `ai.<agent>.asset.settings`, `ai.<agent>.cmd.main`, `ai.<agent>.conv.<type>`
  - Shared: `ai.cmds.<name>`, `ai.steps.<name>`, `ai.conv.<name>`, `ai.prompts.<name>`, `ai.subagents.<name>`
- Rename non-AI extension names for consistency: `ctx7.<kind>.<name>`, `openspec.<kind>.<name>`, `chrome.<kind>.<name>`, `playwright.<kind>.<name>`.
- Update the overlay workflow (`agents-workflow.yaml`) to reference all new names.
- Remove the old fragmented directories (`base/agents/`, `base/cmds/`, `extensions/self/`).
- Update documentation examples to use generic terminology (non-AI examples).
- Add Bats tests validating the new manifest structure and resource name resolution.

## Non-Goals

- Change the kfg CLI behavior or Go code logic (this is a manifest-only reorganization).
- Add new functionality or new agent types.
- Rename core generic steps (`kfg.cleanup`, `kfg.materialize`, etc.) - those remain as-is.
- Change the `apiVersion` from `kfg.dev/v1alpha1`.

## Capabilities

### Modified Capabilities

- `ai-agents-extension`: consolidate all AI-agent resources under `extensions/ai-agents/` with per-agent subdirectories and new naming convention.
- `dev-workflow`: update overlay workflow to use new resource names and extension references.
- `manifest-model`: update documentation to reflect new naming convention and generic examples.

## Impact

- Affected manifests: all files under `.manifests/base/agents/`, `.manifests/base/cmds/`, `.manifests/base/steps/detect-agent.yaml`, `.manifests/base/extensions/self/`, `.manifests/overlay/dev/agents-workflow.yaml`
- Affected tests: new Bats tests for manifest structure validation; existing Bats tests that reference old names
- Affected docs: `docs/manifest-model.md`, `docs/AGENTS.md`
- User-facing shell UX: resource names change but workflow behavior is identical

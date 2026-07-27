# OpenSpec Capability

OpenSpec planning integration for kfg projects.

## Purpose

Provides OpenSpec workflow support for structured change management including proposals, designs, specs, and tasks.

## Prerequisites

- OpenSpec CLI installed

## Entrypoints

- `local/` — Local planning configuration
- `store/` — External store pointer
- `references/` — Read-only references
- `steps/` — Installation steps
- `converters/` — YAML converters
- `assets/` — Configuration assets

## Exported Resources

- `openspec.steps.install` — Installation step
- `openspec.conv.to-yaml` — YAML identity converter
- `openspec.assets.config.local` — Local planning config
- `openspec.assets.config.store` — Store pointer config
- `openspec.assets.config.references` — References config

## Agent Compatibility

All agents (openspec is agent-agnostic).

## Configuration

### Local Config (`openspec.assets.config.local`)

```yaml
schema: spec-driven
context:
  rules: []
artifacts:
  proposal:
    outputPath: proposal.md
  specs:
    outputPath: "specs/**/*.md"
  design:
    outputPath: design.md
  tasks:
    outputPath: tasks.md
```

### Store Config (`openspec.assets.config.store`)

```yaml
store: ""
```

### References Config (`openspec.assets.config.references`)

```yaml
references: []
```

## JSON Patch Examples

### Change schema

```json
[
  { "op": "replace", "path": "/spec/data/schema", "value": "full-cycle" }
]
```

### Add rules

```json
[
  { "op": "add", "path": "/spec/data/context/rules/-", "value": "no-emoji" }
]
```

### Set store ID

```json
[
  { "op": "replace", "path": "/spec/data/store", "value": "my-store-id" }
]
```

### Add references

```json
[
  { "op": "add", "path": "/spec/data/references/-", "value": "https://example.com/ref" }
]
```

## Workflow Usage

Used in dev overlay for openspec install.

## Generated Artifacts

- `openspec/config.yaml` (local mode) — planning configuration
- Skill files under `.opencode/skills/openspec/` (for supported tools)

## Cache and Cleanup

The `openspec.steps.install` step has caching enabled. The engine uses a
leaf-only filesystem snapshot to detect new files created during the step.
Only files and symlinks are registered as artifacts; directories are never
included in the automatic discovery snapshot.

**Cache miss (first run):** The step installs OpenSpec skills and emits
`openspec/config.yaml` (and any skill files) as cached artifacts. The
`openspec/` planning root is created during the session but is *not*
registered as an artifact because the leaf-only snapshot excludes
directories. Cleanup removes only the explicitly cached files, leaving
the planning root intact.

**Cache hit (subsequent runs):** The engine restores only the concrete
files from cache (e.g. `openspec/config.yaml`). The `openspec/` root
directory is not restored or registered. Cleanup operates on the same
file-level artifacts.

**Limitations:** The engine cache never manages or cleans up the
`openspec/` directory tree itself. If you need to remove planning data,
delete `openspec/` manually or use `openspec cleanup`.

## Limitations

- Store and references are beta
- Local and store are mutually exclusive

## Validation

```bash
kustomize build packages/domains/ai-agents/manifests/openspec/config/local
kustomize build packages/domains/ai-agents/manifests/openspec/config/store
kustomize build packages/domains/ai-agents/manifests/openspec/config/references
```

## Canonical Specs

See `docs/context/openspec/specs/` for authoritative behavior specifications.

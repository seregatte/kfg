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

- `openspec/config.yaml` (local mode)
- Skill files

## Cache and Cleanup

Install step has cache.

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

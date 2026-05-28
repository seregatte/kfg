# kfg Consolidated OpenSpec Root

This directory contains the consolidated OpenSpec artifacts for the kfg project, organized by layer:

- **Engine** (`specs/kfg-`) - Core kfg CLI behavior, manifest model, and store operations
- **Framework** (`specs/framework-`) - Shared manifest primitives and reusable steps
- **Domain** (`specs/domain-ai-agents-`) - AI-specific manifests and capabilities

## Structure

- `config.yaml` - Consolidated OpenSpec configuration covering all layers
- `specs/` - Layer-scoped capability specifications
  - `kfg/` - Engine-level specs
  - `framework/` - Framework package specs
  - `domain-ai-agents/` - AI agents domain specs
- `changes/` - Active and archived change proposals and implementations
  - `kfg-*` - Engine layer changes
  - `framework-*` - Framework layer changes
  - `domain-ai-agents-*` - Domain layer changes
  - `archive/` - Completed changes

## Using This Root

Set `OPENSPEC_ROOT_DIR=docs/context/openspec` to point tools at the consolidated root.

### Engine Specifications

Engine capability specifications document core CLI behavior, manifest model, shell runtime, and cross-layer contracts. Key specs include:

- `kfg-project-structure/spec.md` - Repository layout
- `kfg-manifest-model/spec.md` - Resource kinds and composition
- `kfg-shell-runtime-api/spec.md` - Engine-to-framework runtime contract
- `kfg-cli-conventions/spec.md` - Command and flag standards

### Framework Specifications

Framework specifications document shared primitives and exported steps:

- `framework-reusable-framework-steps/spec.md` - Exported framework steps
- `framework-artifact-scoped-cleanup/spec.md` - Artifact lifecycle

### Domain Specifications

Domain specifications document AI-specific capabilities and resources.

## Change Conventions

Changes use layer-prefixed slugs for unambiguous identification:

- Engine changes: `kfg-<slug>`
- Framework changes: `framework-<slug>`
- Domain changes: `domain-ai-agents-<slug>`

Cross-layer changes that affect multiple layers use sibling changes with matching slugs (e.g., `kfg-improve-cache` and `framework-improve-cache`).

## Known Limitations

Specs are organized with layer prefixes: `specs/kfg-*`, `specs/framework-*`, `specs/domain-ai-agents-*`.

To view all available specs, use:
```bash
openspec list --specs
```

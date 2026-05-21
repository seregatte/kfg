# Engine & Core kfg OpenSpec Root

This directory contains engine-level and core kfg OpenSpec artifacts: CLI behavior, manifest model, store operations, and cross-layer architectural contracts.

## Structure

- `config.yaml` - OpenSpec configuration for the engine/core root
- `specs/` - Durable engine and core capability specifications
- `changes/` - Active engine change proposals and implementations

## Engine Responsibilities

The engine (`src/`) implements the core kfg CLI: manifest loading, shell generation, store operations, and workspace management. It defines the shell runtime API that framework and domain packages depend on.

### Core Capabilities

- `build` - Build kustomizations and output merged YAML
- `apply` - Apply kustomizations and generate shell code
- `image` - Manage configuration images
- `workspace` - Manage workspace instances
- `sys log` - Structured logging for shell scripts
- `assets` - Transform assets using converters
- Internal helpers - Shell generation templates, JSONL handling, source organization

### Public API

CLI commands are documented in `docs/cli-reference.md`. The shell runtime API contract is defined in `specs/shell-runtime-api/spec.md`.

## Engine Specs

Engine capability specifications document core behavior, manifest model, CLI conventions, and cross-layer contracts. Specs use normative language (MUST, SHALL) and provide usage examples.

### Key Specs

- Project structure - defines repository layout
- Manifest model - resource kinds, naming, and composition
- Shell runtime API - engine-to-framework runtime contract
- Framework package contract - cross-layer architecture spec
- Domain package contract - cross-layer architecture spec
- CLI conventions - command and flag standards

## Engine Changes

Engine-specific changes are tracked here. Changes affecting both the engine and framework/domain packages are coordinated as sibling changes with matching slugs across:

- `docs/context/kfg/openspec/changes/<slug>/` (engine-level)
- `docs/context/framework/openspec/changes/<slug>/` (framework-level)
- `docs/context/domains/ai-agents/openspec/changes/<slug>/` (domain-level, if affected)

## Running Engine Tests

```bash
# Run all tests
make test-bats

# Unit tests
make test

# Or directly with bats
bats tests/bats/**/*.bats
```

Engine tests are discovered by the canonical `make test-bats` target alongside framework and domain tests.

# AI Agents Domain OpenSpec Root

This directory contains domain-specific OpenSpec artifacts: capability specifications and active changes for the AI agents domain.

## Structure

- `config.yaml` - OpenSpec configuration for the AI agents domain
- `specs/` - Durable domain capability specifications
- `changes/` - Active domain change proposals and implementations

## Domain Responsibilities

The AI agents domain package (`packages/domains/ai-agents/`) provides AI-specific manifests, converters, assets, and resources. It composes the framework package and adds domain-specific behavior.

### Public Entrypoints

- `packages/domains/ai-agents/kustomization.yaml` - domain public entrypoint (includes framework)
- `packages/domains/ai-agents/overlays/dev/` - development overlay

Consumers should reference these files, not internal domain paths.

## Domain Content

The domain provides:

- Agent resources (Claude, Gemini, Pi, OpenCode, etc.)
- AI-specific converters and assets
- Domain-specific MCP integrations
- Prompts and subagent definitions

## Domain Specs

Domain capability specifications document domain-specific behavior. Specs use normative language (MUST, SHALL) and provide usage examples.

### Key Specs

- AI agent resources - documents provided agents
- Domain package contract - (cross-link to engine-level spec)

## Domain Changes

Domain-specific changes are tracked here. Changes affecting both the engine/framework and domain are coordinated as sibling changes with matching slugs across:

- `docs/context/openspec/changes/<slug>/` (engine-level)
- `packages/framework/openspec/changes/<slug>/` (framework-level, if affected)
- `packages/domains/ai-agents/openspec/changes/<slug>/` (domain-level)

## Running Domain Tests

```bash
# Run all domain tests
make test-bats | grep "packages/domains/ai-agents/tests"

# Or directly with bats
bats packages/domains/ai-agents/tests/**/*.bats
```

Domain tests are discovered by the canonical `make test-bats` target alongside engine and framework tests.

## Building the Domain

```bash
# Build the domain package
./bin/kfg build packages/domains/ai-agents/kustomization.yaml

# Build the development overlay
./bin/kfg build packages/domains/ai-agents/overlays/dev
```

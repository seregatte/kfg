# Framework Package OpenSpec Root

This directory contains framework-specific OpenSpec artifacts: capability specifications and active changes.

## Structure

- `config.yaml` - OpenSpec configuration for the framework package
- `specs/` - Durable framework capability specifications
- `changes/` - Active framework change proposals and implementations

## Framework Responsibilities

The framework package (`packages/framework/`) is the kfg shared manifest layer. It exports reusable steps and primitives that domain packages compose.

### Exported Primitives

- `kfg.materialize` - Generate shell code from manifests
- `kfg.cleanup` - Clean up generated artifacts
- `kfg.ensure-gitignore` - Manage repository `.gitignore`
- `kfg.copy-context` - Copy context files into artifacts
- `kfg.materialize-scaffold` - Generate scaffolding from templates

Each exported step is documented by a specification in `specs/`.

### Public Entrypoint

`packages/framework/kustomization.yaml` is the stable public kustomization. Consumers should reference this file, not internal framework paths.

## Framework Specs

Framework capability specifications document the behavior and contract of exported primitives. Specs use normative language (MUST, SHALL) and provide usage examples.

### Key Specs

- Framework step contract - documents exported steps
- Shell runtime API - defines the engine-to-framework runtime contract
- Framework package contract - (cross-link to engine-level spec)

## Framework Changes

Framework-specific changes are tracked here. Changes affecting both the engine and framework are coordinated as sibling changes with matching slugs across:

- `docs/context/openspec/changes/<slug>/` (engine-level)
- `packages/framework/openspec/changes/<slug>/` (framework-level)
- `packages/domains/*/openspec/changes/<slug>/` (domain-level, if affected)

## Running Framework Tests

```bash
# Run all framework tests
make test-bats | grep "packages/framework/tests"

# Or directly with bats
bats packages/framework/tests/**/*.bats
```

Framework tests are discovered by the canonical `make test-bats` target alongside engine and other package tests.

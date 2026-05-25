# AGENTS.md

This file provides guidance to AI agents when working
with code in the kfg repository.

kfg is a standalone CLI for processing YAML manifests
into shell functions. It's written in Go using Cobra for
CLI framework and Viper for configuration.

## Quick Reference

All commands run through the Nix dev shell (`flake.nix`):

```bash
# Build
nix develop --command make build    # → ./bin/kfg

# Test
nix develop --command make test
nix develop --command make test-bats

# Development
nix develop --command make fmt
nix develop --command make lint
nix develop --command make vet
```

## OpenSpec Commands

Always run `openspec` through `kfg run` with the AI
agents dev overlay instead of invoking `openspec`
directly:

```bash
nix develop --command kfg \
  -k packages/domains/ai-agents/overlays/dev \
  run openspec -- <openspec-args>
```

Example:

```bash
nix develop --command kfg \
  -k packages/domains/ai-agents/overlays/dev \
  run openspec -- view
```

## Testing

### Go Unit Tests

Run Go unit tests with `make test`. Tests live in
`src/internal/*_test.go` files.

### Bats Integration Tests

Run integration tests with `make test-bats`. That is the
canonical entrypoint and it discovers tests from engine
and package roots.

Key test roots:

- `tests/bats/`
- `packages/framework/tests/`
- `packages/domains/ai-agents/tests/`

## Canonical Specs

For detailed design and authoritative behavior, refer to
the OpenSpec specs for the layer you are changing:

- `../context/openspec/specs/kfg/` - engine/core kfg specs
- `../context/openspec/specs/framework/` - framework package specs
- `../context/openspec/specs/domain-ai-agents/` - AI agents domain specs

All specs are consolidated in a single OpenSpec root: `../context/openspec/`

### OpenSpec Root Structure

The single consolidated root at `../context/openspec/` organizes specs by layer:

- `specs/kfg/` - engine layer specifications and implementation details
- `specs/framework/` - framework package specifications
- `specs/domain-ai-agents/` - AI agents domain specifications
- `changes/kfg-*` - engine layer changes (prefixed with `kfg-`)
- `changes/framework-*` - framework layer changes (prefixed with `framework-`)
- `changes/domain-ai-agents-*` - domain layer changes (prefixed with `domain-ai-agents-`)

**When working on changes or proposals, refer to the appropriate layer specs:**

- Engine changes: Update specs in `../context/openspec/specs/kfg/`
- Framework changes: Update specs in `../context/openspec/specs/framework/`
- Domain changes: Update specs in `../context/openspec/specs/domain-ai-agents/`
- Cross-layer changes: Create sibling changes with matching slugs across relevant layers (e.g., `kfg-improve-cache` and `framework-improve-cache`)

### Especially Relevant Engine Specs

- `../context/openspec/specs/kfg/project-structure/spec.md`
- `../context/openspec/specs/kfg/manifest-model/spec.md`
- `../context/openspec/specs/kfg/manifest-placeholder/spec.md`
- `../context/openspec/specs/kfg/cli-conventions/spec.md`
- `../context/openspec/specs/kfg/bats-test-layout/spec.md`
- `../context/openspec/specs/kfg/run-command/spec.md`
- `../context/openspec/specs/kfg/store-imagefile/spec.md`
- `../context/openspec/specs/kfg/store-workspace/spec.md`
- `../context/openspec/specs/kfg/session-system/spec.md`

### Framework & Domain Spec Roots

- `../context/openspec/specs/framework/`
- `../context/openspec/specs/domain-ai-agents/`

## Language Policy

All repository-facing written content MUST be in en-US.

This applies to:

- All files under `docs/`, including `docs/context/`
- All OpenSpec content in `../context/openspec/`
- Code comments in all source files
- User-facing strings in source files
- Examples, guides, and agent instructions

Do not introduce Portuguese or mixed-language content in
new or updated files unless the file intentionally records
external third-party content or localized product strings.

## Local State & Gotchas

- Version is injected via ldflags at build time
  (see Makefile)
- The binary must be built before running Bats tests
- Store directory is created on first use
- Log files use `.log` extension (not `.jsonl`)

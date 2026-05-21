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
the OpenSpec roots corresponding to the layer you are
changing, not just `docs/context/kfg/openspec/specs/`.

### OpenSpec Roots & Syncing

When working on changes or proposals, always sync specs
across all relevant OpenSpec roots:

- `docs/context/kfg/openspec/` - engine/core kfg specs and
  changes
- `docs/context/framework/openspec/` - framework package
  specs and changes
- `docs/context/domains/ai-agents/openspec/` - AI agents
  domain specs and changes

**Sync behavior:**

- Read `specs/` and `changes/` from every affected root
- Keep sibling changes with the same slug aligned across
  engine, framework, and domain roots when a change spans
  multiple layers
- Treat `docs/context/kfg/openspec/specs/` as canonical for
  engine behavior
- Treat `docs/context/framework/openspec/specs/` as canonical
  for shared manifest/framework capabilities
- Treat `docs/context/domains/ai-agents/openspec/specs/` as
  canonical for AI agents domain capabilities
- Check `config.yaml` and `changes/` even when `specs/` is
  still empty, as they define scope and context

### Especially Relevant Engine Specs

- `project-structure/spec.md`
- `manifest-model/spec.md`
- `manifest-placeholder/spec.md`
- `cli-conventions/spec.md`
- `bats-test-layout/spec.md`
- `run-command/spec.md`
- `store-imagefile/spec.md`
- `store-workspace/spec.md`
- `session-system/spec.md`

### Package/Domain Spec Roots

- `docs/context/framework/openspec/specs/`
- `docs/context/domains/ai-agents/openspec/specs/`

## Language Policy

All repository-facing written content MUST be in en-US.

This applies to:

- All files under `docs/`, including `docs/context/`
- All OpenSpec content in `docs/context/kfg/openspec/`,
  `docs/context/framework/openspec/`, and
  `docs/context/domains/ai-agents/openspec/`
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

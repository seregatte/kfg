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

For detailed design and authoritative behavior, prefer
the OpenSpec specs in `docs/context/openspec/specs/`.

Especially relevant specs:

- `project-structure/spec.md`
- `manifest-model/spec.md`
- `manifest-placeholder/spec.md`
- `cli-conventions/spec.md`
- `bats-test-layout/spec.md`
- `run-command/spec.md`
- `store-imagefile/spec.md`
- `store-workspace/spec.md`
- `session-system/spec.md`

## Local State & Gotchas

- Version is injected via ldflags at build time
  (see Makefile)
- The binary must be built before running Bats tests
- Store directory is created on first use
- Log files use `.log` extension (not `.jsonl`)

# AI Agent Developer Context

kfg is a declarative shell compiler (Go + Cobra/Viper). See [`README.md`](../README.md) for overview and installation.

## Quick Dev Commands

All commands via Nix dev shell:

```bash
make build          # → ./bin/kfg
make test           # Go unit tests
make test-bats      # Integration tests
make fmt lint vet   # Code quality
```

## Git Worktree Workflow

**CRITICAL RULE: All code changes ONLY in git worktrees. Worktree first, always.**

### Worktree Setup

Worktrees at: `../wkt/kfg/<branch-name>`

**Branch naming:**
- Recognized prefixes: `feature/`, `fix/`, `chore/`, `hotfix/`, `docs/`, `release/`, `main`
- Unrecognized: auto-prefix with `feature/`

Example: `nixai-absort` → `feature/nixai-absort`

### Workflow Steps

1. Determine normalized branch name
2. Create/switch worktree:
   ```bash
   git worktree add ../wkt/kfg/<branch> -b <branch>  # New branch
   git worktree add ../wkt/kfg/<branch> <branch>     # Existing
   ```
3. Push to remote: `git push -u origin <branch>`
4. Create draft PR (base: `main` for feature/fix/docs, ask for release branches)
5. Work in worktree; all commands execute inside `../wkt/kfg/<branch>`

**Important:** NEVER modify files outside a worktree. Main repo stays untouched.

## OpenSpec Commands

Run via `kfg run` with AI agents dev overlay:

```bash
nix develop --command kfg \
  -k packages/domains/ai-agents/overlays/dev \
  run openspec -- view
```

## Testing

- **Go tests:** `make test` → `src/internal/*_test.go`
- **Bats tests:** `make test-bats` (canonical entrypoint)
- **Key roots:** `tests/bats/`, `packages/framework/tests/`, `packages/domains/ai-agents/tests/`

## Canonical Specifications

Refer to OpenSpec specs for authoritative behavior:

- **Engine specs** (`docs/context/openspec/specs/kfg-*`):
  - `kfg-project-structure` - Repository layout
  - `kfg-manifest-model` - Resource kinds
  - `kfg-manifest-placeholder` - Placeholder resolution
  - `kfg-cli-conventions` - Command standards
  - `kfg-bats-test-layout` - Test organization
  - `kfg-shellgen-run-command` - Run command spec
  - `kfg-logging-session-system` - Logging API
  - `kfg-cache-step` - Cache behavior

- **Framework specs** (`docs/context/openspec/specs/framework-*`)
- **Domain specs** (`docs/context/openspec/specs/domain-ai-agents-*`)

All specs consolidated at: `docs/context/openspec/`

## Language Policy

All repository-facing content in **en-US**:
- Files under `docs/` and `docs/context/`
- All OpenSpec content
- Code comments and user-facing strings
- Examples and guides

No Portuguese/mixed-language unless recording third-party content.

## Release Process

### Workflow

1. Stabilize features on `release/<version>` branch
2. Create tag: `git tag v0.0.4 && git push origin v0.0.4`
3. CI runs: GoReleaser builds + publishes, then updates `flake.nix`
4. Create PR: `release/<version>` → `main`
5. Merge to main

### Agent Responsibilities

When user says "release vX.Y.Z":
1. Tag on release branch
2. Push tag (triggers workflow)
3. Wait for `update-flake` job to complete
4. Create PR to main

**Notes:**
- Workflow only triggers on tag pushes (`v*`)
- Never manually modify `flake.nix` hashes/version
- CI runs on `update-flake` push (expected, harmless)

## Local State & Gotchas

- Version injected via ldflags at build time (Makefile)
- Binary must be built before running Bats tests
- Store directory created on first use
- Log files use `.log` extension (not `.jsonl`)

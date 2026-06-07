# AI Agent Developer Context

kfg is a declarative shell compiler (Go + Cobra/Viper). See [`README.md`](../README.md) for overview and installation.

## Quick Dev Commands

All commands via Nix **development shell** (explicit `.#dev`):

```bash
nix develop .#dev --command make build          # → ./bin/kfg
nix develop .#dev --command make test           # Go unit tests
nix develop .#dev --command make test-bats      # Integration tests
nix develop .#dev --command make fmt lint vet   # Code quality
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
5. Work in worktree; all commands execute inside `../wkt/kfg/<branch>` using the **development shell**:
   ```bash
   nix develop .#dev --command make build
   ```

**Important:** NEVER modify files outside a worktree. Main repo stays untouched. Always use the explicit `.#dev` shell for kfg development (the default shell is for consumers).

## Versioning Policy

**CRITICAL RULE: All version changes (in `flake.nix`) MUST go through `release/<version>` branches ONLY.**

Version bumps are **never** allowed in feature/fix branches. This ensures:
1. Version and code changes are tracked separately
2. Release process is deterministic and auditable
3. Main branch always has correct version

### Mandatory Workflow for Version Changes

1. **Feature/Fix branch** → Changes WITHOUT version bump
2. **PR feature/fix → main** → Merge feature changes
3. **Create release branch** → `release/v<NEW.VERSION>`
4. **On release branch:**
   - Cherry-pick or apply feature changes (if needed)
   - Bump version in `flake.nix`
   - Commit with message: `chore: bump version to X.Y.Z`
5. **Tag on release** → `git tag vX.Y.Z && git push origin vX.Y.Z`
6. **CI runs:** GoReleaser + flake.nix auto-update
7. **PR release → main** → Merge release branch

**Non-negotiable:** Agents must enforce this. If version change detected in non-release branch, reject and explain the mandatory workflow.

## OpenSpec Commands

Run via `kfg run` with AI agents dev overlay (use **development shell**):

```bash
nix develop .#dev --command kfg \
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

**See "Versioning Policy" above for mandatory version change workflow.**

### Workflow

1. Stabilize features on `release/<version>` branch (with version bump per versioning policy)
2. Create tag: `git tag v0.0.8 && git push origin v0.0.8`
3. CI runs: GoReleaser builds + publishes, then updates `flake.nix`
4. Create PR: `release/<version>` → `main`
5. Merge to main

### Agent Responsibilities

When user says "release vX.Y.Z" or makes version change request:

**Enforce versioning policy:**
1. Reject version changes in feature/fix branches
2. Create `release/vX.Y.Z` branch (if needed)
3. Cherry-pick feature changes to release branch
4. Bump version only on release branch
5. Tag on release: `git tag vX.Y.Z && git push origin vX.Y.Z`
6. Wait for `update-flake` job to complete
7. Create PR: `release/<version>` → `main`

**Example workflow:**
```
# Feature branch (NO version bump)
git worktree add ../wkt/kfg/feature/my-changes -b feature/my-changes
# ... make changes to flake.nix WITHOUT touching version line
git push && gh pr create --draft --base main

# Release branch (version bump ONLY)
git worktree add ../wkt/kfg/release/v0.0.8 -b release/v0.0.8
git cherry-pick <feature-commits>
# ... update version in flake.nix
git tag v0.0.8 && git push -u origin release/v0.0.8 && git push origin v0.0.8
```

**Notes:**
- Workflow only triggers on tag pushes (`v*`)
- Never manually modify `flake.nix` hashes/version
- CI runs on `update-flake` push (expected, harmless)

## Local State & Gotchas

- Version injected via ldflags at build time (Makefile)
- Binary must be built before running Bats tests
- Store directory created on first use
- Log files use `.log` extension (not `.jsonl`)

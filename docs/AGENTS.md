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

## Git Worktree Workflow

All code changes must be developed in a Git worktree. This ensures
isolation between branches and prevents conflicts with the main
repository state.

### Worktree Path Structure

Worktrees are created at:
```
../wkt/kfg/<branch-name>
```

The project name (`kfg`) is derived from the repository root directory name.

### Branch Naming & Normalization

**Recognized prefixes (passthrough without modification):**
- `feature/` — new features
- `fix/` — bug fixes
- `chore/` — maintenance tasks
- `hotfix/` — critical production fixes
- `docs/` — documentation changes
- `main` — primary development branch
- `release/` — release branches (e.g., `release/1.2.0`)

**Unnaming rule:** If a branch name lacks a recognized prefix, the
agent automatically prefixes it with `feature/`.

Examples:
- `nixai-absort` → `feature/nixai-absort`
- `feature/nixai-absort` → `feature/nixai-absort` (no duplication)
- `fix/login-crash` → `fix/login-crash` (passthrough)
- `main` → `main` (passthrough)
- `release/1.2.0` → `release/1.2.0` (passthrough)

### Agent Workflow

1. **Determine the branch name**  
   From OpenSpec change slug, user instruction, or task context.

2. **Normalize the branch name**  
   Apply prefix rules if needed.

3. **Check if already in worktree**  
   If CWD is already `../wkt/kfg/<branch>`, skip to step 9.

4. **Check if worktree exists**  
   If `../wkt/kfg/<branch>` exists, switch into it and skip to step 8.

5. **Create the worktree**  
   ```bash
   # For a new branch
   git worktree add ../wkt/kfg/<branch> -b <branch>
   
   # For an existing local branch
   git worktree add ../wkt/kfg/<branch> <branch>
   ```

6. **Push the branch to remote**  
   ```bash
   git push -u origin <branch>
   ```

7. **Determine PR base branch**  
   - For `feature/`, `fix/`, `chore/`, `hotfix/`, `docs/`: base is `main`
   - For `release/*`: **ask the user** which branch is the parent before opening the PR

8. **Create a draft PR**  
   ```bash
   gh pr create --draft --base <parent-branch> \
     --title "<branch-name or descriptive title>" \
     --body "<minimal context, link to OpenSpec change if applicable>"
   ```

9. **Work within the worktree**  
   All subsequent commands (build, test, code edits) execute inside
   `../wkt/kfg/<branch>`.

### Important Notes

- Never work in the main repository root when developing features;
  always use the worktree.
- If a worktree already exists and you are not inside it, switch into
  it rather than creating a new one.
- The draft PR serves as early visibility; it can be marked ready for
  review once the work is complete.

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

- Engine specs with prefix: `../context/openspec/specs/kfg-*`
- Framework specs with prefix: `../context/openspec/specs/framework-*`
- Domain specs with prefix: `../context/openspec/specs/domain-ai-agents-*`

All specs are consolidated in a single OpenSpec root: `../context/openspec/`

### OpenSpec Root Structure

The single consolidated root at `../context/openspec/` organizes specs by layer prefix:

- `specs/kfg-*` - engine layer specifications (cross-cutting and core)
- `specs/kfg-shellgen-*` - shell generation feature specs
- `specs/kfg-transform-*` - data transformation feature specs
- `specs/kfg-kustomize-*` - kustomization processing feature specs
- `specs/kfg-cache-*` - step cache feature specs
- `specs/kfg-logging-*` - logging infrastructure feature specs
- `specs/kfg-runtime-*` - workflow runtime feature specs
- `specs/kfg-cli-*` - CLI framework feature specs
- `specs/kfg-build-*` - build and release feature specs
- `specs/framework-*` - framework package specifications
- `specs/domain-ai-agents-*` - AI agents domain specifications
- `changes/kfg-*` - engine layer changes (prefixed with `kfg-`)
- `changes/framework-*` - framework layer changes (prefixed with `framework-`)
- `changes/domain-ai-agents-*` - domain layer changes (prefixed with `domain-ai-agents-`)

**When working on changes or proposals, refer to the appropriate layer specs:**

- Engine changes: Update specs with `kfg-` prefix in `../context/openspec/specs/`
- Framework changes: Update specs with `framework-` prefix in `../context/openspec/specs/`
- Domain changes: Update specs with `domain-ai-agents-` prefix in `../context/openspec/specs/`
- Cross-layer changes: Create sibling changes with matching slugs across relevant layers (e.g., `kfg-improve-cache` and `framework-improve-cache`)

### Especially Relevant Engine Specs

- `../context/openspec/specs/kfg-project-structure/spec.md`
- `../context/openspec/specs/kfg-manifest-model/spec.md`
- `../context/openspec/specs/kfg-manifest-placeholder/spec.md`
- `../context/openspec/specs/kfg-cli-conventions/spec.md`
- `../context/openspec/specs/kfg-bats-test-layout/spec.md`
- `../context/openspec/specs/kfg-shellgen-run-command/spec.md`
- `../context/openspec/specs/kfg-logging-session-system/spec.md`
- `../context/openspec/specs/kfg-cache-step/spec.md`

### Framework & Domain Spec Prefixes

- Framework specs: `../context/openspec/specs/framework-*`
- Domain specs: `../context/openspec/specs/domain-ai-agents-*`

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

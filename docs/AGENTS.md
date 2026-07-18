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
2. Create/switch worktree from the **originating branch**:
   ```bash
   git worktree add ../wkt/kfg/<branch> -b <branch> <base>  # New branch from base
   git worktree add ../wkt/kfg/<branch> <branch>             # Existing
   ```
3. Develop the feature inside the worktree following the **Feature Development Workflow** (see below)
4. Push to remote: `git push -u origin <branch>`
5. Create draft PR targeting the **originating branch** (the base the worktree was created from):
   ```bash
   gh pr create --draft --base <originating-branch>
   ```
   Typically `main`, or `release/vX.Y.Z` during a release cycle.

**Important:** NEVER modify files outside a worktree. Main repo stays untouched. Always use the explicit `.#dev` shell for kfg development (the default shell is for consumers).

## Feature Development Workflow

**CRITICAL RULE: Every feature SHALL follow the OpenSpec cycle inside its own git worktree.**

The OpenSpec cycle creates structured artifacts (proposal, design, tasks) that guide implementation and provide traceability.

### Per-Feature Cycle

For each feature, follow this exact sequence:

1. **Create worktree** from the originating branch (see Git Worktree Workflow above).

2. **Create OpenSpec change** directory INSIDE the worktree at `openspec/changes/<change-name>/` with:
   - `.openspec.yaml` — schema and creation date
   - `proposal.md` — what changes and why
   - `design.md` — how to implement
   - `tasks.md` — implementation checklist

3. **Implement** following `/opsx-propose` → `/opsx-apply` → `/opsx-archive`:
   - **`/opsx-propose`**: Generates proposal, design, tasks artifacts from the change description
   - **`/opsx-apply`**: Implements each task, marking `- [ ]` → `- [x]` as work is completed
   - **`/opsx-archive`**: Moves the change to `archive/YYYY-MM-DD-<name>/` after all tasks are done

4. **Commit and push** the implementation along with the archived change:
   ```bash
   git add -A
   git commit -m "feat: <description>"
   git push -u origin <branch>
   ```

5. **Create draft PR** back to the originating branch:
   ```bash
   gh pr create --draft --base <originating-branch>
   ```

### OpenSpec Change Placement

Each OpenSpec change MUST live at `openspec/changes/<name>/` **INSIDE its feature worktree**, not in the main repo. This ensures the change artifacts are versioned alongside the implementation.

After all tasks are complete and the change is archived, commit everything together. The PR will contain both the implementation code and the archived OpenSpec artifacts.

> **Note:** The `openspec new change` CLI command auto-creates changes at `openspec/changes/` (auto-detected root at CWD). Agents SHOULD create the change directory and artifacts manually using `mkdir` + file writes at `openspec/changes/<name>/` to keep changes in the correct root.

### Multiple Features

When developing multiple features against the same originating branch:

- Each feature gets its **own independent worktree** branching from the same base.
- **Execute features sequentially**, not in parallel. This avoids merge conflicts in shared files (e.g., `manifest/cmds/agents.yaml` when removing multiple agents).
- Features can be merged in any order once complete. Each PR is independent.
- After all feature PRs are merged into the originating branch, the complete set of changes is consolidated.

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

## AI Wizard Command

The `kfg ai` command starts an interactive AI wizard that generates project-specific kfg configurations.

```bash
kfg ai                                          # Start wizard with pi (default)
kfg ai -- "create a new project"                # Start wizard with initial prompt
KFG_AI_AGENT=opencode kfg ai                    # Start wizard with opencode
kfg ai -- --model sonnet "create project"       # Forward args to agent
```

### Environment Variables

- **KFG_AI_AGENT**: Agent to use for the wizard (default: `pi`)
  - Supported values: `pi`, `opencode`
  - Invalid values fall back to `pi` with a warning

### What the Wizard Does

1. Asks about your project type, language, and requirements
2. Selects appropriate building blocks from the kfg domain catalog
3. Presents a summary of what will be generated
4. Asks for confirmation before writing any files
5. Generates complete, project-specific kfg manifests

### Implementation Details

- CLI command: `src/cmd/kfg/ai.go`
- Overlay: `packages/domains/ai-agents/overlays/ai/`
- Wizard skill prompt: `packages/domains/ai-agents/overlays/ai/assets/prompts/wizard.yaml`
- The command wraps `kfg run -k packages/domains/ai-agents/overlays/ai <agent> -- <args>` via subprocess

### Deprecation Notice

The `overlays/dev/` overlay is deprecated. When used, it displays a warning directing users to `kfg ai` instead. The dev overlay remains functional for backward compatibility.

## Testing

- **Go tests:** `make test` → `src/internal/*_test.go`
- **Bats tests:** `make test-bats` (canonical entrypoint)
- **Key roots:** `tests/bats/`, `packages/framework/tests/`, `packages/domains/ai-agents/tests/`

## Canonical Specifications

Refer to OpenSpec specs for authoritative behavior:

- **Engine specs** (`openspec/specs/kfg-*`):
  - `kfg-project-structure` - Repository layout
  - `kfg-manifest-model` - Resource kinds
  - `kfg-manifest-placeholder` - Placeholder resolution
  - `kfg-cli-conventions` - Command standards
  - `kfg-bats-test-layout` - Test organization
  - `kfg-shellgen-run-command` - Run command spec
  - `kfg-logging-session-system` - Logging API
  - `kfg-cache-step` - Cache behavior

- **Framework specs** (`openspec/specs/framework-*`)
- **Domain specs** (`openspec/specs/domain-ai-agents-*`)

All specs consolidated at: `openspec/`

## Language Policy

**CRITICAL RULE: All code and documentation must be written in en-US (American English).**

This applies to:
- Files under `docs/` and `docs/context/`
- All OpenSpec content
- Code comments and user-facing strings
- Examples and guides
- README.md, CONTRIBUTING.md, and all Markdown files
- Commit messages (following Conventional Commits)
- Issue/PR descriptions
- Variable names, function names, and error messages

**No Portuguese or mixed-language content** unless recording third-party content verbatim.

When contributing translations or localizations, keep the canonical documentation in en-US and add localized versions only as separate files (e.g., `docs/pt-BR/getting-started.md`).

Agents must proactively translate any Portuguese content encountered to en-US before committing changes.

## AI Agents Capability Documentation Contract

**CRITICAL RULE: Every public AI agents domain capability MUST have an adjacent README with standardized sections.**

When changing any of the following in `packages/domains/ai-agents/`, the adjacent capability README MUST be updated in the same commit:

- Public resource names or kinds
- Default configuration values
- Configurable fields or JSON Patch targets
- Agent compatibility rules
- Generated artifacts or output paths
- Cache or cleanup behavior
- Kustomization entrypoints
- Workflow usage references

### Capability Documentation Locations

Each independently consumable capability has a `README.md` adjacent to its `kustomization.yaml`:

- `packages/domains/ai-agents/manifests/agents/opencode/README.md`
- `packages/domains/ai-agents/manifests/agents/pi/README.md`
- `packages/domains/ai-agents/manifests/chrome-devtools/README.md`
- `packages/domains/ai-agents/manifests/ctx7/README.md`
- `packages/domains/ai-agents/manifests/gws/README.md`
- `packages/domains/ai-agents/manifests/notebooklm/README.md`
- `packages/domains/ai-agents/manifests/openspec/README.md`
- `packages/domains/ai-agents/manifests/playwright/README.md`

### Required README Sections

Every capability README MUST include: Purpose, Prerequisites, Entrypoints, Exported Resources, Agent Compatibility, Configuration, JSON Patch Examples, Workflow Usage, Generated Artifacts, Cache and Cleanup, Limitations, Validation, and Canonical Specs.

See [CAPABILITY-README-TEMPLATE.md](../packages/domains/ai-agents/CAPABILITY-README-TEMPLATE.md) for the template.

### Enforcement

- Domain Bats tests verify README presence, required sections, advertised resource names, and entrypoint builds
- The wizard discovers capabilities through README files and adjacent Kustomizations
- Never hardcode capability details in the wizard prompt; always reference documentation

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

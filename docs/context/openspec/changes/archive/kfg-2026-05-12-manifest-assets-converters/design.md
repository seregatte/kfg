## Context

kfg currently processes YAML manifests (Step, Cmd, CmdWorkflow) into shell functions. The user's workflows depend on `~/.nixai/images/` — a collection of 14 Dockerfile-like Imagefiles that compose agent configurations (claude, gemini, opencode, pi) with extensions (ctx7, playwright, self, etc.) into immutable images, which are then materialized into workspaces via `.self/` directories.

The image system works but has fundamental problems:
- Imagefiles are versioned outside manifests (not in Git)
- `.self/` materialization is opaque — steps assume files exist without declarative source
- 7 of 8 extensions follow near-identical Imagefile templates (WET)
- No reuse across overlays — each profile rebuilds the same image from scratch

Assets and Converters are already implemented in the codebase (`src/internal/manifest/types.go`, `src/internal/converter/`) but only used via `kfg apply --convert` in isolation. They are "source kinds" explicitly skipped by the resolution pipeline.

## Goals / Non-Goals

**Goals:**
- Eliminate Imagefiles entirely; all data lives in YAML manifests
- Eliminate `.self/` materialization; Steps read `$KFG_BUILD_RESULT_FILE` directly
- Make Steps parameterized (env vars) with zero conditional logic in `run` blocks
- Move all conditional execution to `when` conditions in CmdWorkflow references
- Organize `.manifests/` as a self-contained package with base/overlay composition
- Co-locate Bats tests in `.manifests/tests/`

**Non-Goals:**
- Changing the core manifest model (Step, Cmd, CmdWorkflow remain unchanged)
- Adding new Go code for Assets/Converters (already implemented)
- Backward compatibility with Imagefile-based manifests
- Changing kustomize loading semantics

## Decisions

### D1: Assets as data payloads, Converters as transformations

**Decision:** Assets declare structured data (commands, MCP configs, skills). Converters declare yq-go expressions that transform that data into agent-specific output.

**Why over alternatives:**
- Alternative: bake transformations into Step `run` blocks → would duplicate jq/yq logic across steps, not reusable
- Alternative: keep `.self/` scripts → defeats the purpose of declarative manifests
- yq-go is already a dependency; its expression language covers the jq patterns in use (map, select, dynamic keys, string concatenation)

### D2: `.manifests/` mirrors `images/` structure

**Decision:** `.manifests/` follows the same 3-layer hierarchy as `~/.nixai/images/`:
```
base/agents/          → base image configs (settings.json, etc.)
base/extensions/      → extension data + converters + steps
base/steps/           → core steps (detect, scaffold, cleanup)
overlay/dev/          → dev profile (cmds + workflow)
```

**Why:** Preserves mental model. Users know where to find what. The directory layout communicates dependency: extensions depend on base agents, overlay depends on base.

### D3: Steps are parameterized, not conditional

**Decision:** Steps receive all variation via env vars (`NIXAI_AGENT`, `OUTPUT_DIR`, `FILE_EXT`, etc.). The workflow's `when` conditions decide which parameterization fires.

**Why over alternatives:**
- Alternative: `case $AGENT` inside step → logic hidden, harder to test, WET
- Alternative: one step per agent → DRY violation, combinatorial explosion

### D4: `kfg.cleanup` kept in after phase

**Decision:** The `$KFG_BUILD_RESULT_FILE` cleanup step remains in the workflow's `after` section to prevent stale-file errors on consecutive runs.

**Why:** User already experienced this bug. TMPDIR is not guaranteed to clean between shell invocations.

### D5: MCP assets live in extension, not self

**Decision:** `kfg.extension.ctx7.mcp` lives in `extensions/ctx7/assets/`, not in `extensions/self/assets/mcp/`. The `self/` MCP files are duplicates — each extension owns its MCP declaration.

**Why:** Eliminates duplication. The `self/src/mcp/` files are identical copies of what each extension declares.

### D6: `{env:HOME}` placeholder for path-sensitive configs

**Decision:** Use `{env:HOME}` in Asset data (e.g., `Read(/{env:HOME}/Downloads/**)`) instead of hardcoded paths.

**Why:** kfg's placeholder resolver converts `{env:VAR}` to `$VAR` at generation time, expanded at runtime. Makes manifests portable across machines.

### D7: Bats tests co-located in `.manifests/tests/`

**Decision:** Tests live alongside the manifests they validate, not in `tests/bats/`.

**Why:** Manifests are "user data" — their tests belong with them. Running `make test-manifests` is opt-in and doesn't slow down core CI.

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| yq-go expression complexity for multi-line output (commands, subagents) | Use `output: raw` with string concatenation — tested with simple expressions first |
| Step `run` blocks read `$KFG_BUILD_RESULT_FILE` via yq repeatedly | Performance acceptable for dev workflows; optimize later if needed with temp file extraction |
| Large `.manifests/` directory with 40+ files | Kustomize composition keeps individual files small and focused |
| `when` conditions multiply in large workflows | Workflow sections with comments (phase headers) keep organization clear |

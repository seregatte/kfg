## Context

The kfg manifest system loads YAML resources from `KFG_MANIFEST_PATH` and merges them via kustomize. The current `base/kustomization.yaml` lists 49 individual file references, making it hard to maintain.

The `.nixai/` extensions use a shell-based approach for skill installation (build.sh + jq templates). The `.manifests/` already has declarative Assets/Converters for MCP configs and commands, but lacks Steps for skill installation.

Current state:
- `self` extension: fully migrated (assets + converters)
- `ctx7`, `chrome-devtools`, `playwright`: MCP assets migrated, no install Steps
- `gws`, `notebooklm`, `openspec`: empty directories
- `base/kustomization.yaml`: flat list of 49 files

## Goals / Non-Goals

**Goals:**
- Add generic skill installation Steps for all 6 extensions
- Reorganize manifest structure with hierarchical kustomization files
- Keep Steps agent-agnostic (all specifics via env vars)

**Non-Goals:**
- Modify `.nixai/` source files
- Change existing Assets/Converters
- Add new Go code or external dependencies
- Modify workflow references (they continue to work)

## Decisions

### Decision 1: Generic Steps with env vars

Each install Step receives all agent-specific configuration via `spec.env` variables. No `case`/`if` logic in Step code.

**Rationale**: Keeps Steps simple, testable, and reusable across agents. The workflow layer is responsible for setting the correct env vars per agent invocation.

**Alternatives considered**:
- Agent-specific Steps (e.g., `ctx7.install.claude`): Rejected — creates N×M Step explosion
- Shell functions with parameters: Rejected — still requires shell logic in the Step

### Decision 2: Hierarchical kustomization

Each directory gets its own `kustomization.yaml`. Root references 4 top-level directories: `agents`, `cmds`, `extensions`, `steps`.

**Rationale**: 
- Each directory is self-contained and can be loaded independently
- Adding new resources only requires editing the local kustomization
- Root kustomization stays clean (4 lines vs 49)

**Alternatives considered**:
- Keep flat structure: Rejected — doesn't scale
- Use kustomize `resources` with glob patterns: Not supported by kustomize

### Decision 3: Install Steps use `eval` for command execution

The `INSTALL_CMD` and `INSTALL_FLAGS` env vars are combined and executed via `eval`.

**Rationale**: Allows the workflow to pass different flags per agent (e.g., `--claude`, `--opencode`) without the Step needing to know about agent types.

**Alternatives considered**:
- Direct command execution: Less flexible for flag composition
- Shell function wrapper: Adds unnecessary indirection

## Risks / Trade-offs

- **Risk**: `eval` with user-provided env vars could be a security concern → **Mitigation**: Env vars are set by the workflow, not user input
- **Risk**: Missing env vars cause silent failures → **Mitigation**: Steps validate required vars and log errors
- **Trade-off**: More kustomization files to maintain → **Benefit**: Each directory is self-documenting

## Migration Plan

1. Create kustomization files bottom-up (leaves first, then parents)
2. Create install Steps for each extension
3. Update root `kustomization.yaml` last
4. Verify with `kfg build` that manifests still load correctly

## Open Questions

None — all decisions resolved during planning.

## Context

Three current steps cover variations of the same workflow concern:

- `kfg.agents.steps.settings` converts a single agent asset to one settings file.
- `kfg.convert` converts one asset to one output file today, with the planned direction of plural inputs and outputs.
- `kfg.aggregate-mcp` converts multiple assets, merges the results, wraps the merged object under an agent-specific key, and writes one target file.

All three are built around the same engine call:

```text
Assets -> kfg apply --convert/--use -> generated document(s) -> file artifact(s)
```

The difference lies in output cardinality and post-processing. A single shared step is justified if those differences can be represented by a small, stable contract rather than domain-specific one-off steps.

## Goals / Non-Goals

**Goals:**
- Create a single public manifest step that covers all converter-driven materialization cases in the repo.
- Keep the contract small and explicit, with only two modes: `per-item` and `aggregate`.
- Use one output field, `OUTPUTS`, across both modes.
- Preserve deterministic file generation and artifact registration behavior.
- Make workflow manifests more uniform by grouping materialization by agent and type.

**Non-Goals:**
- Introduce user-defined reducer expressions, custom write strategies, or arbitrary pipelines.
- Retain old step names as aliases.
- Change converter lookup, raw input fallback, or other `kfg apply` CLI behavior.

## Decisions

### Introduce `kfg.materialize` as the single public step

The shared step will be named `kfg.materialize` to describe its responsibility in neutral terms: turn manifest assets into concrete artifacts. This name covers both direct file emission and aggregate file generation better than `kfg.convert`.

Alternative considered: overload `kfg.convert` with aggregate behavior. Rejected because the name becomes too narrow for merge-and-wrap behavior and keeps the old single-item mental model attached to the step.

### Use two explicit modes only

The step contract will have a required `MODE` environment variable with exactly two supported values:

- `per-item`
- `aggregate`

`per-item` writes one file per converted asset using positional mapping.

```text
ASSETS[i] -> CONVERTER -> OUTPUTS[i]
```

`aggregate` converts each asset, merges the converted documents into one object, optionally wraps the result under `WRAP_KEY`, and writes a single output path.

```text
ASSETS[*] -> CONVERTER -> merge -> optional wrap -> OUTPUTS[0]
```

Alternative considered: infer behavior from which environment variables are present. Rejected because explicit mode selection makes validation clearer and avoids ambiguous contracts.

### Use `OUTPUTS` in both modes

The step will use a single output field, `OUTPUTS`, across both modes.

- In `per-item`, `OUTPUTS` is a colon-separated list that must match the `ASSETS` count.
- In `aggregate`, `OUTPUTS` must contain exactly one path.

This keeps the step interface uniform and avoids forcing workflow authors to remember a separate `TARGET` field for aggregate cases.

Alternative considered: keep `TARGET` for aggregate mode. Rejected because it introduces a naming difference without changing the underlying concept of a destination path.

### Make aggregate merge semantics fixed and predictable

Aggregate mode will use the existing deep-merge behavior already implemented by `kfg.aggregate-mcp`:

- convert each asset with the same converter
- deep-merge all converted documents in input order
- if `WRAP_KEY` is set, wrap the merged object under that key
- if the output file already exists, deep-merge the existing file with the newly generated object before writing back

This keeps `kfg.materialize` expressive enough for MCP generation while avoiding a generic reducer API.

Alternative considered: expose a reducer expression through environment variables. Rejected because it would turn the step into an unbounded shell-level DSL and weaken its role as a stable platform primitive.

### Make validation strict and fail-fast

The new shared step will reject incomplete or mismatched configuration:

- `MODE`, `ASSETS`, `CONVERTER`, and `OUTPUTS` are required in all cases
- `MODE=per-item` requires `len(ASSETS) == len(OUTPUTS)`
- `MODE=aggregate` requires `len(OUTPUTS) == 1`
- `WRAP_KEY` is only meaningful in `aggregate`

This replaces the older no-op behavior of `kfg.convert` and aligns better with `kfg.aggregate-mcp` as a workflow contract.

Alternative considered: preserve tolerant no-op behavior for missing values. Rejected because the new step is intended to be the canonical platform primitive, and silent misconfiguration would make workflow failures harder to diagnose.

### Migrate workflow organization by agent and type

The dev workflow should use `kfg.materialize` in batches grouped by agent and materialization type.

Examples:

- `agents.settings.claude`
- `agents.commands.opencode`
- `agents.subagents.claude`
- `agents.mcp.gemini`

This preserves the current phase ordering while making step usage more uniform.

Alternative considered: group all generated artifacts for an agent into one materialization step. Rejected because command, MCP, and subagent outputs currently belong to different workflow phases and use different converters.

## Data Contract

The new shared step contract is:

```yaml
env:
  MODE: ""        # per-item | aggregate
  ASSETS: ""      # colon-separated asset names
  CONVERTER: ""   # converter metadata.name
  OUTPUTS: ""     # colon-separated output paths; exactly one path in aggregate mode
  WRAP_KEY: ""    # optional, aggregate mode only
```

Examples:

```yaml
- name: agents.commands.claude
  step: kfg.materialize
  env:
    MODE: "per-item"
    ASSETS: "kfg.extension.self.commands.git-commit:kfg.extension.self.commands.pr-review"
    CONVERTER: "kfg.convert.self.command.claude"
    OUTPUTS: ".claude/commands/git-commit.md:.claude/commands/pr-review.md"
```

```yaml
- name: agents.mcp.claude
  step: kfg.materialize
  env:
    MODE: "aggregate"
    ASSETS: "kfg.extension.ctx7.mcp:kfg.extension.chrome-devtools.mcp:kfg.extension.playwright.mcp"
    CONVERTER: "kfg.convert.self.mcp.claude"
    OUTPUTS: ".mcp.json"
    WRAP_KEY: "mcpServers"
```

## Risks / Trade-offs

- [Single step becomes too abstract] -> Keep the public contract restricted to two modes and fixed aggregate semantics.
- [Workflow YAML becomes more verbose for trivial single-file cases] -> Accept the small verbosity increase in exchange for a stable shared primitive and consistent mental model.
- [Removing old steps breaks current workflows] -> Update all repository call sites in the same change and avoid compatibility shims.
- [Aggregate merge semantics may not fit future cases] -> Treat new behavior beyond current deep-merge/wrap needs as a separate design decision rather than expanding this contract opportunistically.

## Migration Plan

1. Add `kfg.materialize` with the new strict contract.
2. Migrate current settings, convert, and aggregate workflow call sites to the new step.
3. Remove or replace the old specialized step manifests from the shared base.
4. Add regression coverage for both modes and updated workflow groupings.
5. Update OpenSpec capabilities and developer-facing docs to reference the single step.

## Open Questions

- Should the repository remove the old step files immediately or leave them absent from workflows first and delete them in the same implementation change after tests pass?
- Does any non-dev overlay in the repository rely on the current step names and need to be migrated in the same change?

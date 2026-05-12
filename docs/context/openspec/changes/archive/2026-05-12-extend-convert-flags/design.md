## Context

The current `kfg apply --convert/--use` conversion mode works by looking up an Asset and a Converter resource by `metadata.name` from the loaded manifests, then running the yq-go engine. This is a two-resource lookup that produces a single output file via `os.WriteFile`.

The MCP aggregation problem in `.nixai/manifests/` exposes a gap: individual MCP assets need to be converted and merged incrementally into a single agent config file (e.g., `.claude/settings.local.json`). The current model has no way to:
1. Convert raw JSON/YAML input without an Asset resource
2. Apply an inline yq expression without a Converter resource
3. Feed the output of one conversion as input to another (pipeline)

## Goals / Non-Goals

**Goals:**
- Enable `--convert` to accept raw string input (JSON/YAML literal) when no Asset is found by name
- Add `--with` flag for inline yq expressions, bypassing Converter lookup
- Support stdin (`-f -`) with `--with` for multi-document merge pipelines
- Eliminate external dependency on `jq` for MCP aggregation workflows
- Maintain full backward compatibility with existing `--convert asset --use converter` usage

**Non-Goals:**
- No changes to the converter engine itself (yq-go logic unchanged)
- No changes to Asset or Converter manifest resource types
- No changes to shell generation or workflow resolution
- No automatic merge/aggregation built into the CLI — the step author controls the pipeline

## Decisions

### D1: Asset name fallback to raw string

**Decision**: When `--convert` value does not match any Asset `metadata.name`, treat it as raw input data.

**Rationale**: This avoids adding a new flag. The lookup is deterministic: try name resolution first, then fall back to raw string. The error message improves by listing available assets before falling back.

**Alternatives considered**:
- `--convert-raw` flag: More explicit but adds CLI surface area
- `--input` flag: Too generic, conflicts with existing semantics
- Pipe-only via stdin: Works but loses the ability to inline small literals in Step definitions

### D2: `--with` flag for inline yq expressions

**Decision**: Add `--with "<expression>"` as a new flag that bypasses Converter resource lookup entirely.

**Rationale**: The Converter resource exists for reusable, versionable transformations. But ad-hoc transformations (like JSON merge in a shell step) don't need a resource — just an expression. `--with` is short, clear, and parallels the existing `--use` pattern.

**Alternatives considered**:
- `--expr`: Too generic
- `--yq`: Ties to implementation detail
- Reuse `--use` with `expression:` prefix: Overloads semantics

### D3: Stdin multi-document merge with `-f -` and `--with`

**Decision**: When `-f -` is used with `--with`, stdin content is passed directly to the yq-go engine with the inline expression. No manifest parsing occurs.

**Rationale**: This enables the classic multi-document merge pattern:
```bash
cat existing.json new.json | kfg apply -f - --with 'select(fi == 0) * select(fi == 1)'
```
The yq-go engine supports multi-document input natively via `fileIndex`/`fi`.

**Alternatives considered**:
- New subcommand `kfg convert`: Would duplicate apply logic
- `--merge` flag: Too specific to one use case

### D4: Step-level `kfg.aggregate-mcp` uses read-modify-write in shell

**Decision**: The aggregation step is implemented as a Step (shell code), not as a built-in CLI feature. It reads the target file, converts the asset, merges via `--with`, and writes back.

**Rationale**: Keeps the CLI generic. The Step controls the merge semantics (deep merge via `*`, file existence check, artifact tracking). Different aggregation strategies can be implemented as different Steps.

**Alternatives considered**:
- Built-in `kfg apply --merge`: Would hardcode one merge strategy
- New `kfg merge` command: Adds CLI surface for something that's just shell + convert

### D5: Mutual exclusivity — `--with` cannot be combined with `--use`

**Decision**: `--with` and `--use` are mutually exclusive. `--with` implies no Converter lookup; `--use` requires a Converter resource.

**Rationale**: Clear separation of concerns. Combining them would be ambiguous (which expression wins?).

## Risks / Trade-offs

- **[Risk]** Raw string in `--convert` could be ambiguous if an Asset has the same name as valid JSON. **Mitigation**: Name lookup always takes precedence; the error message lists available assets before falling back.
- **[Risk]** yq-go expression errors in `--with` mode are less debuggable than named Converter errors. **Mitigation**: Include the expression in the error message.
- **[Risk]** Multi-document merge via `-f - --with` requires users to understand yq's `fileIndex` semantics. **Mitigation**: Document the pattern in step examples; the `kfg.aggregate-mcp` step abstracts it for the common MCP use case.
- **[Trade-off]** Using shell-level read-modify-write for aggregation means non-atomic writes (read file → merge → write). **Acceptable**: Generated files are ephemeral per-invocation; race conditions are not a concern.

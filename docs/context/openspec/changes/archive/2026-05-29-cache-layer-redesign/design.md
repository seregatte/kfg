## Context

The kfg cache layer persists Step execution results (artifacts and outputs) across invocations. Today, all cache logic — identity computation, existence checks, store/restore, filesystem diff, metadata serialization — lives in shell helper functions embedded in generated bash templates (`bash_helper.tmpl`, `bash_step.tmpl`). The Go layer only handles admin operations (`kfg sys gc ls/inspect/rm/prune/du`).

This architecture has fundamental problems:
- Shell cache logic is untestable without running full integration flows
- Metadata serialization bugs (multi-line base64, double slashes) recur because shell string manipulation is fragile
- The `__kfg_fs_snapshot` / `__kfg_fs_diff` / `__kfg_add_diff_artifacts` helpers exist solely to support cache artifact detection, adding complexity to the generated shell
- `kfg sys gc` exposes opaque SHA256 hash IDs instead of human-readable step names
- The `kfg sys fs` command group exists only as an implementation detail of cache, not as a user-facing feature
- `src/internal/store/` (v1 legacy) is dead code

The generated shell script must remain sourçable (user sources it to get shell functions), so the shell is the execution environment. But the shell does not need to contain cache logic — it can delegate to Go subcommands.

## Goals / Non-Goals

**Goals:**
- Move all cache logic (identity, exists, store, restore, fs diff, atomic write, metadata) into a Go `cache` package
- Expose cache operations through `kfg sys cache` command namespace
- Shell helpers become thin wrappers calling `kfg sys cache` subcommands
- All CLI operations use `StepReference.name` as identifier (no hash IDs exposed)
- Simplify domain manifests that manually do fs snapshot/diff for artifact detection
- Support `--json` and `--yaml` output flags on all admin subcommands
- Remove all legacy code (`sys_gc.go`, `sys_fs.go`, `store/` package)

**Non-Goals:**
- Changing cache identity algorithm (remains SHA256 of StepReference.name only)
- Changing the cache persistence format (`metadata.yaml` + `artifacts/` directory)
- Adding new cache features beyond the architectural restructuring
- Modifying the shell generation model (scripts remain sourçable)
- Changing how `__kfg_output_set` and `__kfg_add_artifact` work (these stay in shell)

## Decisions

### Decision 1: Shell wrappers delegate to Go subcommands

**Choice:** The `__kfg_cache_exists`, `__kfg_cache_store`, `__kfg_cache_restore` helpers in `bash_helper.tmpl` become thin wrappers that call `kfg sys cache` subcommands.

**Rationale:**
- Keeps the `bash_step.tmpl` interface unchanged — it still calls the same helper names
- All logic (fs diff, metadata, atomic write) moves to Go where it's testable
- Shell helpers are ~3-5 lines each instead of 30-50

**Alternatives considered:**
- Go orchestrates step execution directly (breaks sourçable shell model)
- Shell calls Go via shared library/FFI (unnecessary complexity)

### Decision 2: `kfg sys cache store` reads JSON from stdin

**Choice:** The `store` subcommand reads a JSON object from stdin containing `before` (registered artifacts before step), `after` (registered artifacts after step), `declarative` (manifest-declared artifacts), and `output` (name + base64 value).

**Rationale:**
- JSON handles complex structured data (lists of paths, nested output object)
- Stdin avoids shell escaping issues with paths containing spaces
- Go parses JSON reliably — no shell string manipulation bugs
- The shell wrapper uses `printf` to build the JSON from shell variables

**Alternatives considered:**
- Repeated `--artifact` flags (doesn't scale, hard to pass before/after lists)
- Temp files (cleanup burden, more failure modes)
- YAML stdin (JSON is simpler for shell generation)

### Decision 3: `kfg sys cache restore` emits shell eval-safe output

**Choice:** The `restore` subcommand emits lines like `__kfg_add_artifact "path"` and `__kfg_output_set "ref" "name" "value"` to stdout. The shell wrapper applies them with `eval "$(kfg sys cache restore ...)"`.

**Rationale:**
- No JSON parsing needed in shell — just `eval`
- Reuses existing shell functions (`__kfg_add_artifact`, `__kfg_output_set`)
- Output is human-readable for debugging

**Alternatives considered:**
- JSON output (requires `jq` or complex shell parsing)
- Direct environment variable export (impossible from subprocess)

### Decision 4: Hybrid artifact detection (registered + fs diff)

**Choice:** `kfg sys cache store` receives registered artifact lists (before/after delta) via JSON AND performs its own fs diff on `--workdir` to detect unregistered artifacts.

**Rationale:**
- Steps may forget to call `__kfg_add_artifact` — fs diff catches those
- Steps that explicitly register artifacts get those too
- Union of both approaches provides the most complete artifact set
- Domain manifests (`ctx7`, `openspec`) can be simplified by removing manual snapshot/diff code

**Alternatives considered:**
- Only registered artifacts (misses unregistered files)
- Only fs diff (can't distinguish step artifacts from pre-existing files)

### Decision 5: All operations use StepReference.name as identifier

**Choice:** Every `kfg sys cache` subcommand accepts `StepReference.name` as the key. The SHA256 hash is used only internally as the directory name on disk.

**Rationale:**
- Human-readable — users see `ctx7.steps.install` not `a3f8c1...`
- Consistent — same string used in manifests, logs, and CLI
- Simple — no need to maintain a name→hash lookup table

### Decision 6: Remove `kfg sys fs` entirely

**Choice:** The `kfg sys fs snapshot` and `kfg sys fs diff` subcommands are removed. Fs diff logic moves into the `cache` package as an internal function.

**Rationale:**
- `sys fs` exists only to support cache — no independent use case
- Moving fs diff into Go eliminates two shell helper functions and two subprocess calls per step
- Simplifies the generated shell

### Decision 7: Output format flags on all admin subcommands

**Choice:** `ls`, `inspect`, `du`, `prune` all support `--json` and `--yaml` flags. Default is human-readable text/table.

**Rationale:**
- Consistent CLI experience
- Enables programmatic consumption
- `inspect` defaults to YAML (structured, complete, human-readable)

## Risks / Trade-offs

**Risk: Performance — subprocess overhead per cache operation**
→ Each `kfg sys cache *` call spawns a Go process. For a workflow with N cacheable steps, this adds N subprocess calls for exists + 1 for store/restore.
→ Mitigation: Go CLI startup is fast (~10ms). The overhead is negligible compared to step execution time. Can optimize later with a persistent daemon if needed.

**Risk: Shell variable scope — `eval` of restore output**
→ `eval "$(kfg sys cache restore ...)"` runs in the current shell context, which is correct for setting variables. But if the restore output contains unexpected content, it could execute arbitrary code.
→ Mitigation: The restore output is generated by our Go code from trusted metadata. The risk is equivalent to sourcing the generated script itself.

**Risk: JSON construction in shell — escaping issues**
→ Building JSON with `printf` in shell can break on paths with quotes or special characters.
→ Mitigation: Use `jq`-style escaping in the shell template, or use a Go helper to serialize shell variables to JSON. Alternatively, pass artifacts as NUL-delimited via a temp file.

**Risk: Breaking existing cache entries**
→ Existing entries on disk use the same format (metadata.yaml + artifacts/). The new Go code reads the same format.
→ Mitigation: No migration needed. Existing entries work as-is.

**Trade-off: More subprocess calls vs. simpler shell**
→ Accept trade-off: shell simplicity and Go testability outweigh the minor subprocess overhead.

## Migration Plan

**Phase 1: Implement Go cache package**
- Create `src/internal/cache/` with all modules
- Create `src/cmd/kfg/sys_cache.go` with all subcommands
- Unit tests for all cache operations

**Phase 2: Update shell templates**
- Replace shell helper bodies with Go delegation wrappers
- Update `bash_step.tmpl` cache block if needed

**Phase 3: Remove legacy code**
- Delete `sys_gc.go`, `sys_gc_test.go`, `sys_fs.go`
- Delete `src/internal/store/`
- Update `sys.go` to register new command group

**Phase 4: Update manifests and tests**
- Simplify `ctx7/steps/install.yaml` and `openspec/steps/install.yaml`
- Rewrite `step-cache-isolation.bats` for new architecture

**Phase 5: Update specs**
- Update `kfg-cache-step` spec for Go runtime
- Create `kfg-cache-sys-cache-command` spec
- Remove `kfg-cache-sys-gc-command` spec
- Update `kfg-cache-internal-atomic-write` spec for Go responsibility

**Rollback:** Keep old shell helpers commented out temporarily. If Go cache has issues, uncomment and revert template changes.

## Open Questions

1. **JSON escaping in shell**: Building JSON with `printf` can break on special characters in artifact paths. Should we use a small Go helper (`kfg sys cache encode-json`) to safely serialize shell variables, or handle escaping in the shell template?
2. **`--workdir` flag**: The `store` and `restore` subcommands need `--workdir` to resolve relative artifact paths. Should this be required or default to `$PWD` (which Go reads from environment)?

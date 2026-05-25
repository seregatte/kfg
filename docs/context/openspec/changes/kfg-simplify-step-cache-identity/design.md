## Context

The step cache currently uses a three-part identity: `SHA256(StepReference.name : cache.key : SHA256(spec.run)[:16])`. This means the manifest must carry a `cache.key` field whose purpose overlaps with `StepReference.name`, and any change to the step script automatically invalidates the cache even when the user intentionally wants the old result. A code review also found that artifact paths are kept in a separate `artifact_paths.txt` file next to `metadata.yaml`, splitting the cache entry metadata across two files.

Several correctness bugs were identified at the same time:

- Paths with spaces in declarative artifact lists are silently corrupted (space-separated serialization through shell `IFS`).
- A partial cache write (e.g., interrupted between `mkdir` and `cp`) is treated as a valid hit because `__kfg_cache_exists` only checks for `metadata.yaml`.
- `printf '%b'` in the store helper interprets backslash sequences in user-controlled strings, potentially corrupting `metadata.yaml`.
- `sys gc inspect` only lists top-level entries in the `artifacts/` directory, missing nested paths.
- `isCacheEnabled` is dead code — the identical expression is inlined at every call site.

## Goals / Non-Goals

**Goals:**

- Reduce cache identity to `SHA256(StepReference.name)` — one authoritative input, no redundant `key` field.
- Remove `cache.key` from the manifest model (`Step.spec.cache` and `StepReference.cache`).
- Consolidate artifact path metadata into `metadata.yaml` under an `artifacts:` list, eliminating `artifact_paths.txt`.
- Make cache writes atomic so that a partial write can never masquerade as a valid hit.
- Fix space-in-path corruption for declarative artifacts.
- Fix `printf '%b'` YAML corruption in `__kfg_cache_store`.
- Expose full artifact paths in `sys gc inspect`.
- Activate `isCacheEnabled` at all call sites.
- Maintain backward read compatibility with cache entries created before this change.

**Non-Goals:**

- Automatic cache invalidation on script changes (intentionally removed).
- Multi-input cache keys composed from env vars or file checksums.
- Distributed or shared cache backends.
- Fixing the `sha256sum` cross-platform issue (BUG 1) or `base64` line-wrapping (BUG 2) — those are separate concerns.

## Decisions

### Decision: Identity = SHA256(StepReference.name) only

**Chosen:** Hash only `StepReference.name`.

**Rationale:** `StepReference.name` is already required to be unique within a workflow, making it a natural stable identifier. Adding `cache.key` on top creates a two-level identity where the second level has no clear semantic advantage. Removing `script_hash` means cache invalidation becomes explicit (via `--refresh` or renaming the step reference) rather than implicit, which is consistent with how other build caches (e.g., Nix, Bazel) work when the user controls the key.

**Alternative considered:** Keep `script_hash` as a tie-breaker. Rejected because it makes the cache silently miss on every script edit, forcing users to add `key:` to get stable behavior anyway.

**Alternative considered:** Replace `cache.key` with a stable user-controlled string as the sole identity. Rejected because this is exactly what `StepReference.name` already is.

---

### Decision: Artifact paths move into metadata.yaml

**Chosen:** Add an `artifacts:` YAML list to `metadata.yaml` and stop writing `artifact_paths.txt`.

**Rationale:** A cache entry is a single logical unit. Splitting its index across two files requires both the shell restore path and the Go GC commands to open two files and handle inconsistencies between them. YAML sequences handle paths with spaces natively, whereas the text file needed custom escaping. Having the full path list in `metadata.yaml` also lets `sys gc inspect` display correct paths without a separate recursive directory walk.

**Alternative considered:** Keep `artifact_paths.txt` but fix its format (null-delimited). Rejected because two-file split still complicates the GC reader and the newline-delimited format is already sufficient for the bash side.

---

### Decision: Atomic store via temp directory + rename

**Chosen:** Write to `<cache_path>.tmp`, then `mv` to final path.

**Rationale:** `mv` within the same filesystem is atomic on POSIX systems. This guarantees that `__kfg_cache_exists` either sees a complete entry or nothing — eliminating the window where a partial write is treated as a hit.

**Alternative considered:** Write a sentinel file last. Simpler, but requires changing `__kfg_cache_exists` to check the sentinel and risks leaving orphaned partial directories on crash.

---

### Decision: Newline-separated declarative artifacts + readarray

**Chosen:** Join artifact paths with `\n` in Go; parse with `readarray -t` in bash.

**Rationale:** Newlines cannot appear in POSIX file paths, making them safe as a separator. `readarray -t` (bash 4+) natively handles newline-delimited input into an array, requiring no IFS manipulation.

**Alternative considered:** Null-delimited with `read -d ''`. More robust but harder to debug and less readable in generated code.

---

### Decision: Backward compatibility on restore

**Chosen:** `__kfg_cache_restore` reads artifacts from `metadata.yaml` if the `artifacts:` key is present, otherwise falls back to `artifact_paths.txt`. `readCacheEntry` in Go does the same.

**Rationale:** Existing cache entries in user environments should not become unreadable after upgrading. The fallback adds two lines of shell and a conditional in Go at the cost of a short-lived compatibility window.

## Risks / Trade-offs

- **Script changes no longer auto-invalidate cache** → Users must run `--refresh` or rename the step reference when they want a clean rebuild after editing `spec.run`. This is a deliberate trade-off for simplicity. Mitigation: document this in the step-cache spec and in `kfg sys gc` help text.
- **Existing `cache.key` values in manifests will cause parse errors after upgrade** → The six affected manifests in this repo are updated as part of the change. External users must remove `key:` lines from their own manifests. Mitigation: the field removal is a BREAKING change; document in release notes.
- **`readarray` requires bash 4+** → macOS ships bash 3.2. The shell helper already assumes bash (shebang) and `local -a` (bash 4+ for associative arrays) is used elsewhere. Mitigation: verify that the runtime environment provides bash 4+ (already a de-facto requirement from `declare -A` usage).

## Open Questions

- None. All decisions have been made based on the constraints and code review findings above.

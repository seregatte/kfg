## Why

The step cache identity currently combines three inputs — `StepReference.name`, a user-supplied `cache.key`, and a hash of `spec.run` — which makes the contract confusing and forces users to manage a redundant `key` field when the step reference name already provides a unique identity within the workflow. A code review also revealed that artifact paths are stored in a separate `artifact_paths.txt` file alongside `metadata.yaml`, splitting information that belongs together and complicating both restore logic and the `sys gc inspect` command.

## What Changes

- **BREAKING**: Remove `cache.key` from the `Step.spec.cache` and `StepReference.cache` manifest fields. Cache identity is now determined solely by `StepReference.name`.
- Remove `script_hash` from cache identity. Script changes no longer automatically invalidate the cache; users must use `--refresh` or change the step reference name.
- Cache artifact paths are stored inside `metadata.yaml` under an `artifacts:` list, eliminating the separate `artifact_paths.txt` file.
- Cache store is made atomic: a temporary directory is written and renamed into place, so a partial write can never appear as a valid cache hit.
- Declarative artifact paths are serialized with a newline separator (replacing the space separator) so paths containing spaces are handled correctly.
- `__kfg_cache_store` switches from `printf '%b'` to `echo`-based line construction, preventing backslash escape interpretation from corrupting `metadata.yaml`.
- `sys gc inspect` reads artifact paths directly from `metadata.yaml`, showing the full stored paths rather than only top-level directory entries.
- `isCacheEnabled` replaces the duplicated inline expression at all call sites.
- Backward compatibility: `__kfg_cache_restore` and `readCacheEntry` fall back to `artifact_paths.txt` when reading cache entries created before this change.

## Capabilities

### New Capabilities

- `step-cache-atomic-write`: Atomic cache store via temp-dir-then-rename, guaranteeing no partial cache entries are treated as valid hits.

### Modified Capabilities

- `step-cache`: Cache identity simplified to `StepReference.name` only; `cache.key` field removed from manifest model; artifact paths consolidated into `metadata.yaml`.

## Impact

- `src/internal/manifest/types.go` — `CacheConfig.Key` field removed.
- `src/internal/generate/generate.go` — `computeScriptHash`, `getCacheKey` removed; `isCacheEnabled` activated at all call sites; `CacheKey`/`ScriptHash` removed from template data population; newline separator for declarative artifacts.
- `src/internal/generate/templates/data.go` — `CacheKey` and `ScriptHash` fields removed from `StepData` and `WorkflowStepData`.
- `src/internal/generate/templates/bash_helper.tmpl` — `__kfg_cache_identity` signature simplified; `__kfg_cache_store` uses atomic write and writes artifacts into `metadata.yaml`; `__kfg_cache_restore` reads artifacts from `metadata.yaml` with `artifact_paths.txt` fallback.
- `src/internal/generate/templates/bash_step.tmpl` — `__kfg_cache_identity` call updated; declarative artifact parsing switches to `readarray -t`.
- `src/cmd/kfg/sys_gc.go` — `CacheMetadata` gains `Artifacts []string`; `listArtifacts` replaced by reading from metadata; `base64.StdEncoding` tolerates line-wrapped output.
- Six manifests under `packages/domains/ai-agents/manifests/` — `key:` lines removed from `cache:` blocks.
- Test files — `golden_test.go`, `resolve_test.go`, `parser_test.go` updated to remove `Key`/`ScriptHash` references.

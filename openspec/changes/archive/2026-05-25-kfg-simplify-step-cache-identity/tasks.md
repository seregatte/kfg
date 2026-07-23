## 1. Manifest Model

- [x] 1.1 Remove `Key string` field from `CacheConfig` struct in `src/internal/manifest/types.go`
- [x] 1.2 Remove `key:` from all six package manifests: `ctx7`, `openspec`, `notebooklm`, `playwright`, `chrome-devtools`, and `gws` install steps under `packages/domains/ai-agents/manifests/`

## 2. Go Data Structures

- [x] 2.1 Remove `CacheKey string` and `ScriptHash string` fields from `StepData` in `src/internal/generate/templates/data.go`
- [x] 2.2 Remove `CacheKey string` and `ScriptHash string` fields from `WorkflowStepData` in `src/internal/generate/templates/data.go`
- [x] 2.3 Add `Artifacts []string` field to `CacheMetadata` struct in `src/cmd/kfg/sys_gc.go` (with `yaml:"artifacts,omitempty"`)
- [x] 2.4 Add `Artifacts []string` field to `CacheEntry` struct in `src/cmd/kfg/sys_gc.go`

## 3. Code Generation — Go Side

- [x] 3.1 Delete `computeScriptHash()` function from `src/internal/generate/generate.go`
- [x] 3.2 Delete `getCacheKey()` function from `src/internal/generate/generate.go`
- [x] 3.3 Replace all five inline `cache != nil && (cache.Enabled == nil || *cache.Enabled)` expressions with calls to `isCacheEnabled(cache)` in `src/internal/generate/generate.go`
- [x] 3.4 Remove `CacheKey:` and `ScriptHash:` assignments from all four `WorkflowStepData` population blocks (global before, global after, per-cmd before, per-cmd after) in `src/internal/generate/generate.go`
- [x] 3.5 Remove `CacheKey:` and `ScriptHash:` assignments from `convertStepToTemplateData()` in `src/internal/generate/generate.go`
- [x] 3.6 Change declarative artifact serialization in `generateStepCall` and `generateAfterStepCall` from `strings.Join(uniqueDecl, " ")` to `strings.Join(uniqueDecl, "\n")` in `src/internal/generate/generate.go`

## 4. Shell Templates

- [x] 4.1 Simplify `__kfg_cache_identity` in `bash_helper.tmpl` to accept only `step_ref_name` and compute `SHA256(step_ref_name)` — remove `cache_key` and `script_hash` parameters
- [x] 4.2 Rewrite `__kfg_cache_store` in `bash_helper.tmpl` to write to a `<cache_path>.tmp` temporary directory and rename to final path atomically (`mv`)
- [x] 4.3 In `__kfg_cache_store`, write artifact paths as a YAML `artifacts:` list inside `metadata.yaml` instead of writing `artifact_paths.txt`
- [x] 4.4 Replace `printf '%b'` with line-by-line `echo` (or `printf '%s\n'`) for writing `metadata.yaml` in `__kfg_cache_store` to prevent backslash interpretation
- [x] 4.5 Update `__kfg_cache_restore` in `bash_helper.tmpl` to read artifact paths from the `artifacts:` list in `metadata.yaml` when present
- [x] 4.6 Add `artifact_paths.txt` fallback in `__kfg_cache_restore` for legacy cache entries
- [x] 4.7 Update the `__kfg_cache_identity` call in `bash_step.tmpl` from three arguments to one (`"$__step_ref_name"` only)
- [x] 4.8 Replace `IFS=' ' read -ra __decl_artifacts <<< "$__decl_str"` with `readarray -t __decl_artifacts <<< "$__decl_str"` in `bash_step.tmpl`

## 5. GC Command — Go Side

- [x] 5.1 Update `readCacheEntry()` in `sys_gc.go` to read artifact paths from `metadata.Artifacts` when populated
- [x] 5.2 Add `artifact_paths.txt` fallback in `readCacheEntry()` for legacy cache entries (read file, split by newline)
- [x] 5.3 Update `gcInspectCmd` to print artifact paths from `CacheEntry.Artifacts` instead of calling `listArtifacts(entry.ArtifactsDir)`
- [x] 5.4 Remove or deprecate the `listArtifacts()` function if no longer used after 5.3

## 6. Tests — Go Unit Tests

- [x] 6.1 Remove `CacheKey` and `ScriptHash` from `StepData` literals in `golden_test.go`; remove assertions that check for `test-cache-key`, `ctx7-install`, and `abcd1234` in generated output
- [x] 6.2 Update `TestMergeCache_*` tests in `resolve_test.go` to remove `Key:` from all `CacheConfig` literals and `assert.Equal` checks on `.Key`
- [x] 6.3 Update `TestResolveStepReferences_Cache*` tests in `resolve_test.go` to remove `Key:` from `CacheConfig` literals and `.Key` assertions
- [x] 6.4 Remove the `"step with cache key only"` test case from `TestParseStepWithCache` in `parser_test.go`
- [x] 6.5 Update the remaining `TestParseStepWithCache` cases in `parser_test.go` to remove `Key:` from expected `CacheConfig` values and the `assert.Equal(t, tt.expected.Key, ...)` line

## 7. Validation

- [x] 7.1 Run `nix develop --command make build` and confirm the binary compiles with no errors
- [x] 7.2 Run `nix develop --command make vet` and confirm no vet issues
- [x] 7.3 Run `nix develop --command make test` and confirm all Go unit tests pass
- [x] 7.4 Run `nix develop --command make test-bats` and confirm all Bats integration tests pass
- [x] 7.5 Run `nix develop --command make lint` and confirm no lint issues

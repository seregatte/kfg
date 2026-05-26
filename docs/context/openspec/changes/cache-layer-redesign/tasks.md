## 1. Go Cache Package Foundation

- [ ] 1.1 Create `src/internal/cache/identity.go` — SHA256(StepReference.name) function
- [ ] 1.2 Create `src/internal/cache/metadata.go` — CacheMetadata struct, YAML serialization/deserialization
- [ ] 1.3 Create `src/internal/cache/artifacts.go` — file copy with path preservation, directory walk
- [ ] 1.4 Create `src/internal/cache/fsdiff.go` — filesystem snapshot and diff logic (moved from sys_fs.go)
- [ ] 1.5 Create `src/internal/cache/store.go` — store logic: JSON parse, delta calc, fs diff, artifact copy, atomic write, metadata.yaml
- [ ] 1.6 Create `src/internal/cache/restore.go` — restore logic: read metadata, copy artifacts, emit shell eval-safe output
- [ ] 1.7 Create `src/internal/cache/cache_test.go` — unit tests for identity, metadata, store, restore, fsdiff

## 2. CLI Commands

- [ ] 2.1 Create `src/cmd/kfg/sys_cache.go` — register `kfg sys cache` command group
- [ ] 2.2 Implement `kfg sys cache exists <step-ref>` subcommand — exit 0/1
- [ ] 2.3 Implement `kfg sys cache store <step-ref> --workdir <path>` subcommand — read JSON stdin, delegate to cache package
- [ ] 2.4 Implement `kfg sys cache restore <step-ref> --workdir <path>` subcommand — emit shell eval-safe stdout
- [ ] 2.5 Implement `kfg sys cache ls [--json|--yaml]` subcommand — table/JSON/YAML output
- [ ] 2.6 Implement `kfg sys cache inspect <step-ref> [--json|--yaml]` subcommand — YAML default, full value
- [ ] 2.7 Implement `kfg sys cache rm <step-ref>...` subcommand — remove by name
- [ ] 2.8 Implement `kfg sys cache prune [--json|--yaml]` subcommand — remove entries >30 days
- [ ] 2.9 Implement `kfg sys cache du [--json|--yaml]` subcommand — disk usage report

## 3. Shell Template Updates

- [ ] 3.1 Update `bash_helper.tmpl` — remove `__kfg_cache_identity`, `__kfg_cache_exists`, `__kfg_cache_store`, `__kfg_cache_restore` function bodies; replace with thin Go wrappers
- [ ] 3.2 Update `bash_helper.tmpl` — remove `__kfg_fs_snapshot`, `__kfg_fs_diff`, `__kfg_add_diff_artifacts` functions
- [ ] 3.3 Update `bash_helper.tmpl` — add `__kfg_serialize_artifacts` helper for JSON array construction from shell variables
- [ ] 3.4 Update `bash_step.tmpl` — cache block calls `__kfg_cache_exists` / `__kfg_cache_restore` / `__kfg_cache_store` (interface unchanged, bodies now delegate to Go)

## 4. Legacy Code Removal

- [ ] 4.1 Remove `src/cmd/kfg/sys_gc.go` — replaced by sys_cache.go
- [ ] 4.2 Remove `src/cmd/kfg/sys_gc_test.go` — replaced by cache_test.go
- [ ] 4.3 Remove `src/cmd/kfg/sys_fs.go` — fs diff absorbed into cache package
- [ ] 4.4 Remove `src/internal/store/store.go` — legacy v1 store, dead code
- [ ] 4.5 Remove `src/internal/store/artifacts.go` — dead code
- [ ] 4.6 Update `src/cmd/kfg/sys.go` — register `sys_cache` instead of `sys_gc` and `sys_fs`

## 5. Domain Manifest Updates

- [ ] 5.1 Update `packages/domains/ai-agents/manifests/ctx7/steps/install.yaml` — remove `__kfg_fs_snapshot` before/after and `__kfg_add_diff_artifacts` calls
- [ ] 5.2 Update `packages/domains/ai-agents/manifests/openspec/steps/install.yaml` — remove `__kfg_fs_snapshot` before/after and `__kfg_add_diff_artifacts` calls

## 6. Test Rewrites

- [ ] 6.1 Rewrite `tests/bats/workflows/step-cache-isolation.bats` — update generated code pattern tests for new architecture (shell wrappers calling `kfg sys cache`)
- [ ] 6.2 Add bats tests for `kfg sys cache ls/inspect/rm/prune/du` CLI behavior
- [ ] 6.3 Add bats tests for `kfg sys cache exists` hit/miss scenarios
- [ ] 6.4 Add bats tests for `kfg sys cache store/restore` round-trip

## 7. Spec Updates

- [ ] 7.1 Update `docs/context/openspec/specs/kfg-cache-step/spec.md` — apply delta specs from change
- [ ] 7.2 Create `docs/context/openspec/specs/kfg-cache-sys-cache-command/spec.md` — apply new spec from change
- [ ] 7.3 Remove `docs/context/openspec/specs/kfg-cache-sys-gc-command/spec.md` — replaced by new spec
- [ ] 7.4 Update `docs/context/openspec/specs/kfg-cache-internal-atomic-write/spec.md` — apply delta specs from change

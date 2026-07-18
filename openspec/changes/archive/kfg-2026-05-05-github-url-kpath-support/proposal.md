## Why

Users want to use kfg with GitHub repositories directly without manually cloning them first. Currently, users must clone a repo locally before running `kfg build/apply/run`. Additionally, users want to configure a default kustomization source via environment variable to avoid passing `-k` or a path argument every time.

## What Changes

- Add `KFG_KPATH` environment variable for default kustomization source
- Support GitHub repository URLs as kustomization paths in `build`, `apply`, and `run` commands
- Make the path argument optional in `build` command when `KFG_KPATH` is set
- Make the `-k` flag optional in `apply` and `run` commands when `KFG_KPATH` is set
- Kustomize loader receives GitHub URLs directly and clones them transparently (`--depth=1` shallow clone)
- Update CLI help and examples to document GitHub URL usage and `KFG_KPATH`

## Capabilities

### New Capabilities

- `github-url-source`: Support for GitHub repository URLs as kustomization source paths (github.com, raw.githubusercontent.com)
- `kpath-env-var`: `KFG_KPATH` environment variable for default kustomization source with fallback priority

### Modified Capabilities

- `build-command`: Path argument becomes optional when `KFG_KPATH` is set; GitHub URLs accepted
- `apply-command`: `-k` flag becomes optional when `KFG_KPATH` is set; GitHub URLs accepted in `-k`

## Impact

**Affected code:**
- `src/cmd/kfg/build.go` — optional argument, env var fallback, URL detection
- `src/cmd/kfg/apply.go` — env var fallback in `runApplyPipeline`, URL detection
- `src/cmd/kfg/run.go` — env var fallback, URL detection
- `src/internal/config/config.go` — add `KFG_KPATH` binding and getter
- `src/internal/urlresolve/resolver.go` — new package for URL detection

**Affected tests:**
- `src/internal/urlresolve/resolver_test.go` — new tests for URL detection
- `src/internal/config/config_test.go` — tests for `KFG_KPATH` binding
- `src/cmd/kfg/build_test.go` — tests for env var fallback
- `src/cmd/kfg/apply_test.go` — tests for env var fallback

**Affected documentation:**
- CLI help for `build`, `apply`, `run` commands
- `docs/AGENTS.md` — add `KFG_KPATH` to environment variables table

**Dependencies:**
- No new external dependencies (kustomize already handles git cloning)
- Uses existing `sigs.k8s.io/kustomize/api` git repo cloning functionality

**Backward compatibility:**
- Fully backward compatible — existing workflows continue to work unchanged
- `KFG_KPATH` is optional; commands work without it
- GitHub URL support is additive; existing local paths work unchanged
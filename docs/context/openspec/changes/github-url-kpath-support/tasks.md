## 1. Setup

- [ ] 1.1 Create `src/internal/urlresolve/` package directory
- [ ] 1.2 Create `src/internal/urlresolve/resolver.go` with package declaration

## 2. URL Detection Implementation

- [ ] 2.1 Implement `IsGitHubURL(arg string) bool` function in resolver.go
- [ ] 2.2 Add unit tests for `IsGitHubURL` in `resolver_test.go`
- [ ] 2.3 Test cases: HTTPS GitHub URL, GitHub URL with path separator, GitHub URL with ref parameter, non-GitHub URL, local path

## 3. Config Binding for KFG_KPATH

- [ ] 3.1 Add `viper.BindEnv("kpath", "KFG_KPATH")` in `src/internal/config/config.go` Initialize()
- [ ] 3.2 Add `GetKPath() string` getter function in config.go
- [ ] 3.3 Add unit tests for KFG_KPATH binding in `config_test.go`
- [ ] 3.4 Test cases: env var set, env var empty, viper override

## 4. Build Command Modifications

- [ ] 4.1 Change `cobra.ExactArgs(1)` to `cobra.MaximumNArgs(1)` in build.go
- [ ] 4.2 Add source resolution logic: arg[0] > GetKPath() > error
- [ ] 4.3 Add error message when no source available
- [ ] 4.4 Pass GitHub URLs directly to kustomize loader (no preprocessing needed)
- [ ] 4.5 Update command `Use` to `build [path-or-url]`
- [ ] 4.6 Update `Long` description to mention KFG_KPATH and GitHub URLs
- [ ] 4.7 Update examples to include GitHub URL and KFG_KPATH usage

## 5. Apply Command Modifications

- [ ] 5.1 Add KFG_KPATH fallback in `runApplyPipeline()` when kustomizePath is empty
- [ ] 5.2 Add error message when no source available (no -k, no -f, no KFG_KPATH)
- [ ] 5.3 Pass GitHub URLs directly to kustomize loader
- [ ] 5.4 Update `Long` description to mention KFG_KPATH and GitHub URLs
- [ ] 5.5 Update examples to include GitHub URL and KFG_KPATH usage

## 6. Run Command Modifications

- [ ] 6.1 Add KFG_KPATH fallback when runKustomizePath is empty
- [ ] 6.2 Add error message when no source available
- [ ] 6.3 Pass GitHub URLs directly to kustomize loader
- [ ] 6.4 Update `Long` description to mention KFG_KPATH and GitHub URLs
- [ ] 6.5 Update examples to include GitHub URL and KFG_KPATH usage

## 7. Unit Tests

- [ ] 7.1 Add tests for build command with KFG_KPATH fallback in `build_test.go`
- [ ] 7.2 Add tests for build command without argument or KFG_KPATH (error case)
- [ ] 7.3 Add tests for apply command with KFG_KPATH fallback in `apply_test.go`
- [ ] 7.4 Add tests for run command with KFG_KPATH fallback in `run_test.go`
- [ ] 7.5 Verify all existing tests still pass (backward compatibility)

## 8. Integration Tests

- [ ] 8.1 Add integration test for `kfg build https://github.com/owner/repo//path` (tag with `//go:build integration`)
- [ ] 8.2 Add integration test for `kfg build https://github.com/owner/repo//path?ref=v1.0.0`
- [ ] 8.3 Add integration test for `kfg apply -k https://github.com/owner/repo//path`
- [ ] 8.4 Add integration test for `kfg run -k https://github.com/owner/repo//path claude`
- [ ] 8.5 Add integration test for KFG_KPATH with local path
- [ ] 8.6 Add integration test for KFG_KPATH with GitHub URL

## 9. Shell Validation Tests (Bats)

- [ ] 9.1 Add Bats test for shell generation from GitHub URL source
- [ ] 9.2 Verify generated shell code is valid when source is GitHub URL
- [ ] 9.3 Verify KFG_KPATH produces same shell output as explicit path

## 10. Documentation

- [ ] 10.1 Update `docs/AGENTS.md` environment variables table to include KFG_KPATH
- [ ] 10.2 Add KFG_KPATH description: default kustomization source path or URL
- [ ] 10.3 Update CLI reference documentation if exists

## 11. Validation

- [ ] 11.1 Run `make test` and verify all unit tests pass
- [ ] 11.2 Run `make test-bats` and verify all Bats tests pass
- [ ] 11.3 Run `make lint` and verify no lint errors
- [ ] 11.4 Manual test: `kfg build https://github.com/owner/repo//manifests`
- [ ] 11.5 Manual test: `KFG_KPATH=./manifests kfg build`
- [ ] 11.6 Manual test: `kfg build` without KFG_KPATH (verify error message)
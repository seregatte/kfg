## Why

The current CI workflow (`ci.yml`) builds and tests only the PR branch in isolation — it does not verify that the branch integrates cleanly with its target. This means merge conflicts and integration regressions are discovered only after merge. Additionally, the CI environment uses a bare Go toolchain rather than the Nix dev shell, diverging from the local development environment and bypassing the Bats integration test suite entirely.

## What Changes

- Replace the existing `ci.yml` with a Nix-based workflow that installs Nix via `DeterminateSystems/nix-installer-action` and caches the dev shell via `magic-nix-cache-action`
- Add a pre-test merge step that fetches the PR target branch and performs a trial `git merge --no-commit --no-ff` to detect conflicts before any tests run (skipped for direct pushes to `main`)
- Run all test targets through the Nix dev shell: `make build`, `make test`, `make vet`, and `make test-bats`
- Support any PR target branch (not just `main`) using `github.base_ref`

## Capabilities

### New Capabilities

- `ci-nix-pr-tests`: GitHub Actions CI workflow that runs the full Nix-based test suite (Go unit tests + Bats integration tests) on every PR update, after performing a trial merge with the PR target branch to detect integration issues early

### Modified Capabilities

- `kfg-build`: CI now builds via `nix develop --command make build` instead of bare `go build`, ensuring the binary under test matches the Nix-packaged environment

## Impact

- `.github/workflows/ci.yml`: full replacement
- No source code changes
- First run will be slower due to Nix store population; subsequent runs benefit from `magic-nix-cache` layer caching
- Bats tests (`packages/framework/tests/`, `packages/domains/ai-agents/tests/`, `tests/bats/`) now run in CI

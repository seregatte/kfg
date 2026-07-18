## Context

The existing `.github/workflows/ci.yml` runs on `pull_request` events targeting `main` and on direct pushes to `main`. It installs Go directly, builds with `go build`, and tests with `go test ./...`. It does not:

- Use the Nix dev shell defined in `flake.nix`
- Run the Bats integration test suite (`make test-bats`)
- Simulate the post-merge state before testing

The local development workflow (documented in `AGENTS.md`) runs all commands through `nix develop --command make <target>`. The CI environment diverges from this, meaning regressions in Bats tests or Nix-specific build behavior go undetected until a developer runs them locally.

## Goals / Non-Goals

**Goals:**
- Run the full test suite (Go unit tests + Bats integration tests) on every PR update, using the Nix dev shell
- Detect merge conflicts and integration failures against the PR's target branch before the tests run
- Support PRs targeting any branch, not just `main`
- Keep CI fast through Nix store caching

**Non-Goals:**
- Running CI on direct pushes to feature branches (only on PRs and pushes to `main`)
- Replacing the release workflow
- Publishing or deploying artifacts from CI

## Decisions

### Decision: Use DeterminateSystems Nix actions

**Choice**: `DeterminateSystems/nix-installer-action` + `DeterminateSystems/magic-nix-cache-action`

**Rationale**: These are the de-facto standard for Nix in GitHub Actions. The installer enables flakes and nix-command automatically. The cache action transparently proxies the Nix store through GitHub Actions cache, cutting dev shell rebuild time on subsequent runs from minutes to seconds.

**Alternative considered**: `cachix/install-nix-action` — less opinionated, requires manual flake feature flags, no built-in caching story. Rejected for being more setup work with no benefit for this use case.

### Decision: Trial merge with `--no-commit --no-ff`

**Choice**: After checkout, fetch the target branch and run `git merge --no-commit --no-ff origin/${{ github.base_ref }}`.

**Rationale**: This is the simplest way to detect merge conflicts without creating a merge commit. If the merge fails, Git exits non-zero and the workflow fails immediately with a conflict error — before wasting time on the Nix setup and test run. The `--no-commit` flag means the working tree is modified in place but nothing is committed, so subsequent steps test the integrated code.

**Alternative considered**: Using a separate "merge check" job — adds complexity and latency. Rejected in favor of an early step in the single job.

**Alternative considered**: GitHub's merge queue feature — requires repo-level configuration and branch protection rules. Rejected as out of scope for this change.

### Decision: Single job, sequential steps

**Choice**: One job (`test`) with steps in order: checkout → merge → nix install → nix cache → build → unit tests → vet → bats.

**Rationale**: The Bats tests depend on the binary produced by `make build`. The build depends on the merged working tree. Sequential steps naturally express these dependencies. A matrix or parallel jobs would complicate the merge-first requirement and add overhead for a test suite that runs in under a minute locally.

### Decision: Skip merge step on `push` to `main`

**Choice**: Add `if: github.event_name == 'pull_request'` to the trial merge step.

**Rationale**: On a direct `push` to `main`, `github.base_ref` is empty. The merge step would fail unconditionally. Direct pushes to `main` represent already-integrated code; testing the branch as-is is the correct behavior.

## Risks / Trade-offs

- **First-run slowness** → Mitigation: `magic-nix-cache-action` populates on the first run and subsequent runs are fast. The `nixai` flake input is an external dependency that will be fetched; its size is bounded.
- **Nix installer adds ~30s per run** → Acceptable. No mitigation needed; this is the cost of environment fidelity.
- **Bats tests may have environment assumptions** (e.g., paths, tools) that differ from Linux runners → Mitigation: local dev is macOS but runners are `ubuntu-latest`. If Bats tests have macOS-only assumptions, they will surface here and should be fixed separately.
- **Trial merge leaves dirty working tree** → Intended. Steps after the merge step operate on the merged code. Git status will show the merge in progress; this does not affect `go test` or `make test-bats`.

## Migration Plan

1. Replace `.github/workflows/ci.yml` with the new workflow file
2. No branch protection rule changes required — the existing `pull_request` trigger already feeds GitHub's required status checks
3. On the first PR after deployment, the Nix cache will be cold; expect a longer run (~3-5 min). Subsequent runs will be faster (~1-2 min)
4. Rollback: revert the `ci.yml` change to the previous version

## Open Questions

- None. The approach is fully determined.

## 1. Worktree Setup

- [x] 1.1 Create or switch to worktree `../wkt/kfg/kfg-ci-nix-pr-tests`
- [x] 1.2 Push branch to remote and open draft PR targeting `main`

## 2. Replace CI Workflow

- [x] 2.1 Replace `.github/workflows/ci.yml` trigger block to fire on `pull_request` (any target) and `push` to `main`
- [x] 2.2 Add checkout step with `fetch-depth: 0` so the full git history and all remote refs are available
- [x] 2.3 Add trial merge step: fetch `origin/${{ github.base_ref }}` and run `git merge --no-commit --no-ff origin/${{ github.base_ref }}` with `if: github.event_name == 'pull_request'`
- [x] 2.4 Add `DeterminateSystems/nix-installer-action@v4` step to install Nix with flakes enabled
- [x] 2.5 Add `DeterminateSystems/magic-nix-cache-action@v2` step to cache the Nix store
- [x] 2.6 Replace build step with `nix develop --command make build`
- [x] 2.7 Replace unit test step with `nix develop --command make test`
- [x] 2.8 Replace vet step with `nix develop --command make vet`
- [x] 2.9 Add Bats integration test step: `nix develop --command make test-bats`

## 3. Verification

- [x] 3.1 Open a test PR and confirm the trial merge step runs and the workflow passes all steps
- [x] 3.2 Confirm Bats tests appear in the workflow run logs
- [ ] 3.3 Confirm a PR with an intentional conflict fails at the merge step before reaching Nix setup
- [ ] 3.4 Confirm a second run on the same branch shows Nix cache hits in the logs

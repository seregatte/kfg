## ADDED Requirements

### Requirement: CI runs on every PR update against any target branch
The CI workflow SHALL trigger on all `pull_request` events regardless of the target branch, and on direct `push` events to `main`.

#### Scenario: PR opened or updated
- **WHEN** a pull request is opened, synchronized, or reopened against any branch
- **THEN** the CI workflow starts and runs all steps

#### Scenario: Direct push to main
- **WHEN** a commit is pushed directly to `main`
- **THEN** the CI workflow runs without the trial merge step

### Requirement: Trial merge detects integration conflicts before tests run
The CI workflow SHALL perform a trial merge of the PR branch with the target branch before any build or test step. If the merge produces conflicts, the workflow SHALL fail immediately with a non-zero exit code.

#### Scenario: Clean merge
- **WHEN** the PR branch merges cleanly with `origin/${{ github.base_ref }}`
- **THEN** the trial merge step succeeds and subsequent steps operate on the merged working tree

#### Scenario: Merge conflict
- **WHEN** the PR branch has conflicts with `origin/${{ github.base_ref }}`
- **THEN** `git merge --no-commit --no-ff` exits non-zero and the workflow fails at the merge step without running tests

#### Scenario: Trial merge skipped on push to main
- **WHEN** the event is a direct `push` (not a `pull_request`)
- **THEN** the trial merge step is skipped (`if: github.event_name == 'pull_request'`)

### Requirement: Build and tests run inside the Nix dev shell
All build and test commands SHALL be executed via `nix develop --command <cmd>` to ensure environment parity with local development.

#### Scenario: Build via Nix
- **WHEN** the CI workflow runs the build step
- **THEN** it executes `nix develop --command make build` and produces `./bin/kfg`

#### Scenario: Go unit tests via Nix
- **WHEN** the CI workflow runs the unit test step
- **THEN** it executes `nix develop --command make test` and all Go tests pass

#### Scenario: Go vet via Nix
- **WHEN** the CI workflow runs the vet step
- **THEN** it executes `nix develop --command make vet` with zero findings

#### Scenario: Bats integration tests via Nix
- **WHEN** the CI workflow runs the integration test step
- **THEN** it executes `nix develop --command make test-bats` which discovers and runs all Bats test roots

### Requirement: Nix store is cached between CI runs
The CI workflow SHALL use `DeterminateSystems/magic-nix-cache-action` to cache the Nix store, reducing cold-start time on subsequent runs.

#### Scenario: Cache hit on repeated run
- **WHEN** a CI run follows a previous run on the same branch
- **THEN** the Nix store is restored from cache and the dev shell does not need to be rebuilt from scratch

#### Scenario: Cache miss on first run
- **WHEN** no cache exists for the current Nix configuration
- **THEN** the Nix store is built from scratch and stored in cache for future runs

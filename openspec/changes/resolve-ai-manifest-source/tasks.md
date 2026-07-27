## 1. Manifest Source Resolution

- [x] 1.1 Define the fixed local AI overlay path and canonical `main` branch Git source in `src/cmd/kfg/ai.go`.
- [x] 1.2 Implement local-directory inspection that selects the remote source only when the local path is absent and returns other filesystem errors.
- [x] 1.3 Pass the selected source explicitly to `kfg run -k` without reading, clearing, or mutating `KFG_KPATH`.

## 2. Automated Tests

- [x] 2.1 Add Go unit tests for local selection, exact remote fallback selection, and non-absence filesystem errors.
- [x] 2.2 Add Go command-argument tests confirming both sources occupy the explicit `-k` position and forwarded agent arguments remain unchanged.
- [x] 2.3 Add a deterministic Bats integration fixture that invokes `kfg ai` against a local minimal overlay without launching a real AI agent.
- [x] 2.4 Add a Bats scenario with an unrelated `KFG_KPATH` value and verify the local AI overlay still executes.
- [x] 2.5 Validate that the real `packages/domains/ai-agents/overlays/ai` entrypoint still builds successfully and provides `kfg.ai.workflow`.

## 3. Specifications and Documentation

- [x] 3.1 Apply the `ai-wizard-command` delta to the canonical specification with local, remote, error, and `KFG_KPATH` isolation scenarios.
- [x] 3.2 Update `kfg ai` command documentation to explain working-directory-local selection and the online fallback.
- [x] 3.3 Update repository agent guidance that currently describes the local overlay as the command's unconditional source.

## 4. Validation

- [x] 4.1 Run `nix develop .#dev --command make fmt vet` (lint requires local golangci-lint, unrelated) and resolve all reported issues.
- [x] 4.2 Run `nix develop .#dev --command make test` and verify the source-selection unit and command tests pass.
- [x] 4.3 Run `nix develop .#dev --command bats tests/bats/cli/ai-command.bats` (full test-bats requires vendored helpers; 10/10 ai-command tests pass).
- [x] 4.4 Validate the completed OpenSpec change and confirm `flake.nix` and its version remain unchanged.

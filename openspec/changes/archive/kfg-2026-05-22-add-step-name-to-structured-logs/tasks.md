## 1. Engine logging contract

- [ ] 1.1 Extend `src/internal/logger` so `KFG_STEP_NAME` is enriched as `step_name` for both normal and explicit-session log writes.
- [ ] 1.2 Update `src/cmd/kfg/log.go` and related logging paths to normalize legacy `step:<name>` shell components into structured Step attribution.
- [ ] 1.3 Add or update unit tests covering `step_name` enrichment, legacy normalization, and unchanged behavior for non-Step logs.

## 2. Generated runtime and message normalization

- [ ] 2.1 Update generated shell runtime templates so Step wrappers export scoped `KFG_STEP_NAME` and restore prior context on exit.
- [ ] 2.2 Preserve compatibility for current `__kfg_log_*` calling styles while enabling message-only Step logs to emit structured Step attribution.
- [ ] 2.3 Normalize touched Go log messages to the agreed message style and update any affected assertions.

## 3. Documentation and validation

- [ ] 3.1 Add OpenSpec deltas for `shell-runtime-api` and `log-command` describing `step_name`, legacy compatibility, and `kfg sys log` behavior.
- [ ] 3.2 Update user-facing logging documentation to match the implemented logger behavior and structured Step attribution.
- [ ] 3.3 Run `nix develop --command make test` and relevant Bats coverage for logging-related engine behavior.

## 1. Logger and CLI compatibility

- [x] 1.1 Extend logger enrichment so `KFG_STEP_NAME` is persisted as `step_name` in JSONL output and preserved across the explicit-session logging path.
- [x] 1.2 Update `kfg sys log` handling to normalize legacy `step:<name>` components into structured `step_name` data without breaking existing shell callers.
- [x] 1.3 Add unit tests for `step_name` enrichment, legacy component normalization, and unchanged behavior for non-Step logs.

## 2. Generated runtime support

- [x] 2.1 Update generated Step wrappers to export scoped `KFG_STEP_NAME` for the active Step reference and restore prior state on exit.
- [x] 2.2 Update `__kfg_log_*` helpers to preserve both current calling styles, including one-argument message-only calls and explicit component plus message calls.
- [x] 2.3 Add generator and runtime tests covering Step-scoped logs, cache logs with `step_name`, and compatibility with legacy helper usage.

## 3. Manifest migration and message normalization

- [x] 3.1 Refactor framework Step manifests under `packages/framework/manifests/steps/` to stop encoding Step identity in `component` strings and to follow the standardized message style.
- [x] 3.2 Refactor AI-agent Step manifests under `packages/domains/ai-agents/manifests/` to stop encoding Step identity in `component` strings and to follow the standardized message style.
- [x] 3.3 Normalize touched Go log messages to the new message style where they are currently out of pattern.

## 4. Specs, docs, and end-to-end validation

- [x] 4.1 Add OpenSpec delta files for the modified engine, framework, and AI-agent capabilities affected by structured Step logging.
- [x] 4.2 Update user-facing logging documentation to describe `kfg sys log`, `step_name`, and the current verbosity and JSONL behavior.
- [x] 4.3 Run `nix develop --command make test` and `nix develop --command make test-bats`, fixing any logging-related regressions in engine, framework, or domain scenarios.

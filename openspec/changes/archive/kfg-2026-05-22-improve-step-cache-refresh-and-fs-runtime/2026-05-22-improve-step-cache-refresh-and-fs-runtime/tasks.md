## 1. Engine cache refresh overwrite

- [x] 1.1 Update generated Step wrappers so refresh bypasses cache restore but still stores cache after successful execution.
- [x] 1.2 Update cache store helpers to remove any existing cache entry before writing the replacement artifact tree and metadata.
- [x] 1.3 Add or update generator and runtime tests covering refresh-driven cache rebuilds and removal of stale cached artifacts.

## 2. Internal filesystem CLI and runtime wrappers

- [x] 2.1 Add `kfg sys fs snapshot` with normalized relative path output and `--maxdepth` support where `0` means unlimited depth.
- [x] 2.2 Add `kfg sys fs diff` that reports only paths newly present in the `after` snapshot.
- [x] 2.3 Add Go unit tests for filesystem snapshot/diff behavior, including depth limits, invalid depth values, normalization, and empty/no-change cases.
- [x] 2.4 Add generated shell runtime wrappers for `sys fs` and document their contract in the touched specs and CLI help text.

## 3. Quiet internal kfg subprocesses

- [x] 3.1 Add a generated shell helper for nested internal `kfg` execution that forces child-scoped `KFG_VERBOSE=0` without mutating the parent environment.
- [x] 3.2 Route runtime-owned `kfg` subprocesses, including the new `sys fs` wrappers, through the internal execution helper.
- [x] 3.3 Add tests covering `KFG_VERBOSE=3` parent execution so nested internal `kfg` calls do not emit child startup logs while still returning their functional stdout/stderr and exit status.

## 4. Framework and domain adoption

- [x] 4.1 Update `packages/framework/manifests/steps/materialize.yaml` to use the internal `kfg` execution helper for nested conversion calls.
- [x] 4.2 Update `packages/domains/ai-agents/manifests/ctx7/steps/install.yaml` to snapshot `OUTPUT_DIR`, diff it at `--maxdepth 1`, and register newly created children as artifacts.
- [x] 4.3 Remove redundant ctx7 `artifacts:` declarations from `packages/domains/ai-agents/overlays/dev/agents-workflow.yaml` after the Step owns artifact discovery.

## 5. Validation and documentation

- [x] 5.1 Add or update OpenSpec delta files across engine, framework, and AI-agent roots for refresh overwrite semantics, internal filesystem commands, quiet internal subprocesses, materialize, and ctx7 install.
- [x] 5.2 Update user-facing CLI help and documentation so refresh wording reflects bypass-plus-rebuild behavior and the internal `sys fs` command surface is accurately described.
- [x] 5.3 Run `nix develop --command make test` and `nix develop --command make test-bats`, fixing regressions across engine, framework, and domain scenarios.

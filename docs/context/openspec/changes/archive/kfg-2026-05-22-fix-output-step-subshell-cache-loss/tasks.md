## 1. Runtime wrapper fix

- [x] 1.1 Update generated Step wrappers for `spec.output` so they capture stdout without running the full Step body in command substitution.
- [x] 1.2 Ensure temporary output capture state is cleaned up correctly on success and failure paths.

## 2. Cache regression coverage

- [x] 2.1 Add unit or golden coverage for output-producing Step wrappers preserving runtime side effects.
- [x] 2.2 Add Bats coverage for a cacheable Step with `spec.output` that also calls `__kfg_add_artifact`, verifying both output and artifacts are persisted and restored.

## 3. Validation

- [x] 3.1 Validate the ctx7 install flow or an equivalent fixture as the real-world regression case.
- [x] 3.2 Run `nix develop --command make test` and `nix develop --command make test-bats`, fixing any regressions in runtime output/cache behavior.

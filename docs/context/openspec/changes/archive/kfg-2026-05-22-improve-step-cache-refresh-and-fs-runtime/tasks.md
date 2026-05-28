## 1. Engine runtime and CLI

- [ ] 1.1 Implement refresh-driven cache overwrite in generated runtime templates and cache helpers.
- [ ] 1.2 Add `kfg sys fs snapshot` and `kfg sys fs diff` plus the underlying Go helpers and tests.
- [ ] 1.3 Add generated runtime wrappers for quiet internal `kfg` execution and `sys fs` delegation, with generator and integration coverage.

## 2. Docs and verification

- [ ] 2.1 Update engine specs and CLI help text for refresh rebuild wording and the new internal filesystem command surface.
- [ ] 2.2 Run engine-focused unit and Bats coverage for cache overwrite, filesystem snapshot/diff, and nested internal `kfg` verbosity behavior.

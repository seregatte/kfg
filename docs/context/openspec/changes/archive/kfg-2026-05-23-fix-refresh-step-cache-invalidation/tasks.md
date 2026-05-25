## 1. Refresh semantics

- [x] 1.1 Update generated step wrappers so `KFG_REFRESH` invalidates only the current step cache entry before execution.
- [x] 1.2 Replace refresh-specific runtime logs and help text so they describe cache invalidation and rebuild semantics explicitly.
- [x] 1.3 Add or update generator and runtime tests covering step-scoped refresh invalidation and rebuilt cache entries.

## 2. Diff artifact helper

- [x] 2.1 Add a generated shell runtime helper that registers artifacts from snapshot diffs rooted at a provided directory.
- [x] 2.2 Update `packages/domains/ai-agents/manifests/ctx7/steps/install.yaml` to use the shared diff artifact helper instead of the inline loop.
- [x] 2.3 Add tests covering relative-path diff registration and refreshed cache reconstruction of artifact lists.

## 3. Specs and validation

- [x] 3.1 Add or update OpenSpec delta files for `step-cache`, `shell-runtime-api`, `cli-conventions`, `apply-command`, and `run-command`.
- [x] 3.2 Update user-facing CLI documentation so `KFG_REFRESH` is described as step cache invalidation plus rebuild.
- [x] 3.3 Run `nix develop --command make test` and `nix develop --command make test-bats`.

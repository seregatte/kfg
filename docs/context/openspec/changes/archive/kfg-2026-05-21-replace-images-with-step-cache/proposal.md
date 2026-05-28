## Why

Install-oriented Steps currently perform network downloads and tool bootstrapping during every `kfg run`, which makes agent startup slow and unpredictable. The existing image/workspace feature set solves a different problem, adds separate CLI surface area, and does not help the dominant workflow where generated agent files and Step outputs need fast reuse across invocations.

## What Changes

- Add Step-level cache configuration so reusable Steps can restore previously generated artifacts and outputs instead of re-running expensive setup logic.
- Allow cache configuration at both `Step` and `StepReference`, with `StepReference` overrides taking precedence for per-invocation control.
- Persist Step outputs along with artifacts so cached Steps continue to satisfy `when` conditions and `$kfg.output(...)` consumers.
- Add `--refresh` to `kfg apply` and `kfg run`, and introduce `KFG_REFRESH` as the public runtime override for bypassing Step cache.
- Add internal `kfg sys gc` commands to inspect, remove, and prune cached runtime entries.
- Normalize generated logging helpers to the `__kfg_*` naming convention used by the rest of the shell runtime API.
- **BREAKING** Remove the `kfg image` and `kfg workspace` command families and retire the associated image/workspace storage model.
- **BREAKING** Remove the image/workspace specifications and replace their retained runtime use case with Step cache behavior.

## Capabilities

### New Capabilities
- `step-cache`: declarative Step caching, runtime restore semantics, refresh controls, and operational GC for cached entries.

### Modified Capabilities
- `manifest-model`: Step and StepReference schemas gain cache configuration and cached output semantics.
- `shell-runtime-api`: generated shell runtime gains cache helpers, cache-aware output restore, and `__kfg_log_*` logging helpers.
- `cli-conventions`: CLI surface loses `image` and `workspace`, adds refresh flags, `sys gc`, and public environment variable help.
- `run-command`: `kfg run` gains cache refresh behavior and executes cached Steps transparently.
- `apply-command`: `kfg apply` gains refresh-aware shell generation semantics.
- `store-imagefile`: imagefile-based image composition is removed.
- `store-image-build`: image build behavior is removed.
- `store-image-persistence`: image persistence behavior is removed.
- `store-image-metadata`: image metadata behavior is removed.
- `store-workspace`: workspace materialization and restore behavior is removed.
- `cli-store-isolation`: store CLI isolation rules must be revised to reflect `sys gc` and the absence of image/workspace commands.

## Impact

- Affects manifest parsing, validation, and resolution in `src/internal/manifest` and `src/internal/resolve`.
- Affects shell generation and runtime helpers in `src/internal/generate`.
- Adds cache storage and GC behavior under the existing `KFG_STORE_DIR` root.
- Removes command implementations under `src/cmd/kfg/image.go` and `src/cmd/kfg/workspace.go` plus supporting internal packages.
- Requires migration of install Steps and log helper usage in `packages/framework/` and `packages/domains/ai-agents/`.
- Requires updates to engine-level OpenSpec specs, CLI help text, unit tests, and Bats coverage.

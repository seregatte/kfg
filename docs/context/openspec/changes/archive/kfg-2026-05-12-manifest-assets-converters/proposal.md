## Why

The current workflow relies on Dockerfile-like Imagefiles (`~/.nixai/images/`) that are opaque, versioned outside manifests, and materialized via `.self/` directories at runtime. This creates a maintenance burden: 14 Imagefiles with near-identical patterns, scripts hidden from the manifest system, and no declarative description of what constitutes each image's content.

Replacing images with **Assets** (declarative data payloads) and **Converters** (transformation rules) embedded in the manifest pipeline eliminates the image concept entirely, makes all content versionable in YAML, and enables reuse through kustomize composition.

## What Changes

- **New resource kind: Assets** — declares data payloads with format metadata (YAML input)
- **New resource kind: Converter** — declares transformations via yq-go expressions (input → output)
- **New directory: `.manifests/`** — self-contained manifest package at project root with base + overlay structure
- **Deprecation: Imagefiles** — `~/.nixai/images/` replaced by declarative Assets/Converters
- **Deprecation: `.self/` materialization** — no runtime intermediate scripts; all logic via Steps + `when` conditions
- **Step model update** — Steps become parameterized via env vars (zero conditional logic in `run` blocks)
- **Workflow model update** — `when` conditions (`equals`, `in`, `anyOf`) drive all conditional execution
- **Testing** — Bats tests co-located in `.manifests/tests/`, run via `make test-manifests`

## Capabilities

### New Capabilities
- `manifest-source-layer`: Assets and Converter resource kinds — declaration and validation semantics
- `manifest-directory-layout`: `.manifests/` directory structure — base/overlay composition, co-located tests
- `converter-execution`: Converter application via yq-go engine — how expressions transform Asset data

### Modified Capabilities
- `manifest-model`: Adds Source Layer to the resource kind taxonomy (Step, Cmd, CmdWorkflow + Assets, Converter). Updates reference validation to include Asset→Converter resolution.
- `project-structure`: Adds `.manifests/` as a valid manifest source path alongside `.kfg/manifests/` and `~/.config/kfg/manifests`. Documents base/overlay pattern and co-located tests.

## Impact

- **Manifest model** (`src/internal/manifest/`): Assets/Converter types already exist — no Go changes needed for parsing
- **Generation pipeline** (`src/internal/generate/`): Steps that read `$KFG_BUILD_RESULT_FILE` for Asset data — no generator changes needed
- **No breaking changes** to existing Step/Cmd/CmdWorkflow semantics
- **Backward compatibility NOT required** per project conventions — old Imagefile-based manifests can be retired
- **New Makefile target**: `make test-manifests` runs Bats tests in `.manifests/tests/`
- **CLI**: `kfg apply --convert <asset> --use <converter>` already supported — no CLI changes needed

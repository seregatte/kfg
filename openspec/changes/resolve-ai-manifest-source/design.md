## Context

`kfg ai` invokes the current executable as `kfg run -k <source> <agent>`. Its source is currently the repository-relative `packages/domains/ai-agents/overlays/ai`, which supports editing the overlay in a checkout but is unavailable to installed binaries running in other projects. The `run` command already accepts Kustomize Git sources, and an explicit `-k` value already takes precedence over the generic `KFG_KPATH` environment setting.

The source decision must preserve the contributor workflow from a checkout or worktree root without introducing configuration that can accidentally point the wizard at unrelated manifests.

## Goals / Non-Goals

**Goals:**

- Select the working-tree AI overlay when its repository-relative directory exists.
- Select the canonical `main` branch Git source when the local directory does not exist.
- Pass either selection explicitly to `kfg run -k` so `KFG_KPATH` cannot affect the wizard.
- Keep source-selection tests deterministic and independent of the network and agent executables.
- Surface local filesystem errors instead of silently changing behavior.

**Non-Goals:**

- Searching parent directories for a checkout.
- Resolving manifests relative to the kfg executable or `KFG_DIR`.
- Adding `KFG_AI_KPATH` or another source override.
- Falling back to the remote source when a selected local overlay is invalid.
- Changing the AI overlay composition, manifest merge semantics, generated shell, framework steps, or domain resources.

## Decisions

### Resolve only the fixed working-directory path

The command will check `packages/domains/ai-agents/overlays/ai` relative to its current working directory. If the directory exists, that same relative path is selected. If it does not exist, the command selects `https://github.com/seregatte/kfg.git//packages/domains/ai-agents/overlays/ai?ref=main`.

This keeps the development convention explicit: contributors run the command from a checkout or worktree root. Parent traversal was rejected because it adds discovery rules that are unnecessary for the established workflow and could unexpectedly select a different checkout.

### Treat only absence as a remote fallback

The filesystem check will distinguish a missing path from other inspection failures. A missing local directory selects the remote source. Permission and other filesystem errors are returned with context. If the local directory exists but its Kustomization is invalid, the normal `kfg run` pipeline reports that error; the command does not hide development defects by retrying remotely.

### Keep the source explicit and independent of generic configuration

The selected source remains the value of the existing `-k` argument. `kfg ai` will not read, clear, or mutate `KFG_KPATH`. Because an explicit flag is authoritative in `kfg run`, an inherited generic setting cannot redirect the wizard, while all other kfg commands retain their current configuration behavior.

An AI-specific environment override was rejected because there is no current customization requirement and it would expand the public configuration surface.

### Isolate the filesystem decision for unit testing

The local-or-remote decision will be kept in a small resolver used by `aiCmd`. Unit tests can change into temporary directories and assert the returned source without spawning an agent, loading Kustomize, or accessing GitHub. Existing Bats help coverage remains responsible for the public command surface.

The manifest loading algorithm is otherwise unchanged: Kustomize receives exactly one explicit local directory or Git source, builds the existing AI overlay, and the current resolver selects `kfg.ai.workflow` and the requested agent command.

## Risks / Trade-offs

- [Running from a nested checkout directory uses the remote source] -> Document that local development selection is relative to the checkout or worktree root; avoid adding implicit traversal.
- [The mutable `main` ref can change behavior between invocations] -> This is intentional so installed `kfg ai` uses the current wizard manifests; local development remains available for controlled testing.
- [Remote fallback requires network access] -> Preserve local selection when the overlay is present and allow Kustomize to report remote loading failures normally.
- [A local directory can be present but invalid] -> Fail through the existing build path instead of masking the defect with a remote retry.

## Migration Plan

No data or manifest migration is required. Release the CLI behavior with its specification and documentation updates. Rollback consists of restoring the fixed local source in `ai.go`; no generated or persisted format changes are involved.

## Open Questions

None.

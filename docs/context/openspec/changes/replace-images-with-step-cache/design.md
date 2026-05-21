## Context

`kfg` currently has two unrelated runtime strategies for generated assets. The first is the Step execution model, where each invocation runs shell code inline and can register outputs and artifacts that exist only for the current command. The second is the image/workspace system, which persists immutable image snapshots and materializes them into a working directory with backup/restore semantics. The slow agent startup problem sits squarely in the Step path: install Steps for ctx7, openspec, playwright, chrome-devtools, gws, and notebooklm download packages or generate files every time `kfg run` executes.

The current runtime cannot reuse those Step results because outputs live only in `__kfg_outputs`, artifact tracking is only a runtime array, and there is no cache store or restore path. At the same time, image/workspace commands introduce a separate object model, parser, storage metadata, and CLI surface that are not part of the main manifest-to-shell workflow.

This change is cross-cutting. It modifies the manifest schema, resolved workflow model, generated shell runtime, CLI contracts, and several engine specs. It also removes the image/workspace path completely, which simplifies the product but requires a coherent replacement for persisted generated artifacts.

## Goals / Non-Goals

**Goals:**
- Make expensive Steps cacheable by declaration, with per-reference overrides for workflow-specific behavior.
- Persist both artifacts and Step outputs so downstream `when` conditions and env expansions continue to work on cache hits.
- Add a user-controlled refresh path through `--refresh` and `KFG_REFRESH`.
- Provide internal operational commands to inspect and clean cached runtime entries.
- Remove image/workspace commands, implementation packages, and specs in the same change.
- Align generated logging helper names with the `__kfg_*` runtime naming convention.

**Non-Goals:**
- Preserve backward compatibility with the image/workspace feature set.
- Introduce multi-kind cache management in the initial `sys gc` interface.
- Add a separate public cache command group.
- Cache arbitrary external side effects that are not represented by registered artifacts or declared Step outputs.

## Decisions

### Decision: Cache is configured on Step definitions and Step references

The model adds `cache` to `Step.spec` and `StepReference`, using the same override pattern already used for `env` and artifacts. `StepReference.cache` overrides `Step.spec.cache` so shared Steps can declare a default while workflows can refine or disable caching per invocation.

Why this approach:
- It preserves reusable Step definitions.
- It allows the same Step to be cached differently in different workflows.
- It keeps the cache contract declarative and colocated with execution semantics.

Alternatives considered:
- Step-only cache config: rejected because per-invocation tuning is required for shared Steps.
- Separate cache manifests: rejected as unnecessary schema expansion.

### Decision: Cache identity is per resolved Step invocation

Each cache entry is keyed from the resolved workflow invocation, not only from `Step.metadata.name`. The effective key combines `ResolvedStep.Name`, the resolved cache key string, and a hash of `spec.run`. This ensures that different workflow references to the same Step can diverge safely.

Why this approach:
- Prevents collisions between multiple references to the same reusable Step.
- Automatically invalidates cache when the Step script changes.
- Keeps the user-facing key simple while still providing deterministic invalidation.

Alternatives considered:
- Hash only `StepReference.name`: rejected because env-sensitive variations would be stale by default.
- Hash the entire resolved Step struct: rejected because it makes cache behavior harder to reason about and explain.

### Decision: Persist artifacts from declarative and runtime sources

Cache restore must reproduce the observable effects of a Step. The runtime will capture artifacts from three sources: `Step.spec.artifacts`, `StepReference.artifacts`, and the runtime delta of `KFG_ARTIFACTS` registered during the Step invocation. The cache writer persists the union of those paths.

Why this approach:
- Existing install Steps use a mix of declarative artifacts and dynamic `__kfg_add_artifact` calls.
- It avoids adding another manifest field just to describe cacheable outputs.
- It keeps the runtime aligned with what Steps already declare as artifacts.

Alternatives considered:
- Only use `spec.artifacts`: rejected because several existing install Steps would not cache anything useful.
- Add a new `cache.outputs` path list: rejected because it duplicates artifact intent.

### Decision: Persist outputs as part of cache metadata

Any Step with `spec.output` must restore that output on a cache hit. The cache metadata will store output values in a shell-safe serialized form, using base64 so multiline output and arbitrary content survive round trips. On restore, the runtime repopulates `__kfg_outputs` before later Steps execute.

Why this approach:
- `when` conditions and `$kfg.output(...)` already depend on runtime output state.
- Install Steps like ctx7 use multiline generated content as a Step output.
- Restoring only artifacts would break downstream workflow logic.

Alternatives considered:
- Recompute outputs from cached files: rejected because not every Step has a deterministic file representation of its output.
- Ignore outputs on cache hits: rejected because it breaks workflow semantics.

### Decision: Add refresh controls to apply and run

`kfg run --refresh` forces cache bypass during immediate execution. `kfg apply --refresh` generates shell code that exports refresh state so a later sourced invocation also bypasses cache. `KFG_REFRESH` is the public environment variable contract for both code paths.

Why this approach:
- It gives users an explicit invalidation escape hatch without adding public cache management commands.
- It keeps apply and run consistent.
- It matches the current CLI model where shell runtime behavior is influenced by generated environment.

Alternatives considered:
- Run-only refresh: rejected because generated shell should expose the same control.
- Hidden refresh env only: rejected because users expect a flag.

### Decision: Manage cache operationally through `kfg sys gc`

The initial operational surface will live under `kfg sys gc` with `ls`, `inspect`, `rm`, `prune`, and `du`. The interface starts without `kind` because Step cache is the only managed object type in v1, but metadata should remain extensible for future cache kinds.

Why this approach:
- Keeps cache management out of the main user-facing CLI.
- Provides operational tooling similar to image management workflows without exposing cache as a top-level product concept.
- Leaves room to expand later without committing to a public `kfg cache` interface.

Alternatives considered:
- No operational commands: rejected because the cache would grow without inspection or cleanup tools.
- Public `kfg cache` commands: rejected because the desired surface is internal.

### Decision: Remove image/workspace systems completely

The change removes the `image` and `workspace` command trees, their internal implementation packages, the imagefile parser, and their associated specs. Cache storage remains rooted under `KFG_STORE_DIR`, but the new subdirectory uses a cache-specific path rather than reusing `images/`.

Why this approach:
- It eliminates a separate persistence model that no longer serves the primary runtime use case.
- It reduces CLI and implementation complexity.
- It avoids carrying legacy path names that imply image semantics after the feature is removed.

Alternatives considered:
- Keep images/workspaces and add Step cache alongside them: rejected because the stated direction is to migrate away from them.
- Reuse `images/` for cache entries: rejected because the name would misrepresent the new model.

### Decision: Rename generated logging helpers to `__kfg_log_*`

Generated logging helpers will adopt the same naming convention as the rest of the shell runtime API, while still delegating to `kfg sys log` under the hood. All manifests and tests that call the old `_kfg.log.*` names will be migrated in the same change.

Why this approach:
- It aligns the logging API with the documented `__kfg_*` helper family.
- It avoids introducing another special-case naming convention in shell runtime.
- It preserves the current CLI logging backend.

Alternatives considered:
- Keep `_kfg.log.*`: rejected because it is inconsistent with the rest of the runtime API.
- Change CLI logging commands too: rejected because the requested scope is helpers only.

## Risks / Trade-offs

- [Cache entries can grow without automatic retention] -> Mitigation: provide `kfg sys gc` commands in the first version and document that there is no automatic GC.
- [Cache restore may miss externally visible side effects] -> Mitigation: define cache safety in spec terms as artifacts plus outputs only, and limit initial rollout to suitable Steps.
- [Removing image/workspace is a breaking CLI and spec change] -> Mitigation: mark all removals as breaking, update help/docs in the same change, and remove tests/specs together.
- [Output serialization bugs could break downstream `when` conditions] -> Mitigation: use base64 metadata, add unit tests for multiline and special-character outputs, and add integration coverage for cache-hit output consumers.
- [Renaming logging helpers can break domain manifests if done partially] -> Mitigation: migrate all helper call sites in the same PR and add generator tests for the new helper names.

## Migration Plan

1. Add cache fields to manifest types and resolved workflow structures.
2. Implement shell runtime cache helpers, output serialization, and `__kfg_log_*` helper names.
3. Integrate cache restore/store into generated Step execution and add refresh propagation in `apply` and `run`.
4. Add `kfg sys gc` commands against the new cache metadata format.
5. Migrate expensive install Steps and other manifests to declare cache and use the renamed log helpers.
6. Update CLI help and public environment variable documentation.
7. Remove image/workspace CLI commands, internal packages, specs, and tests.
8. Run unit tests, generator tests, and Bats coverage for cache hits, refresh, GC commands, and removal of legacy commands.

Rollback during development is straightforward because backward compatibility is not required. Partial rollout is not recommended; cache, refresh, helper renames, and image/workspace removal should land as one coherent change.

## Open Questions

- None. The change intentionally chooses manual operational GC, no compatibility mode for images/workspaces, and no multi-kind GC selector in v1.

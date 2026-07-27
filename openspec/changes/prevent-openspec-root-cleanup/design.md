## Context

Cacheable Steps take a filesystem snapshot before execution and compare it with the workdir after execution. `SnapshotDirectory` and the generated shell snapshot currently include both directories and leaf paths. A new tree therefore produces artifacts such as `openspec`, `openspec/config.yaml`, and generated agent integration files.

Cache restore reproduces every cached artifact and emits `__kfg_add_artifact` for each path. The framework `kfg.cleanup` Step intentionally treats a registered directory as recursively disposable and runs `rm -rf` on it. If the agent creates `openspec/changes/...` after an `openspec` directory artifact is restored, final cleanup removes both the generated installation state and the new persistent planning data.

The engine must distinguish automatically discovered leaf outputs from directories that a manifest author deliberately declares as disposable. Existing cache entries do not retain artifact provenance, so they cannot be safely reinterpreted after discovery semantics change.

## Goals / Non-Goals

**Goals:**

- Prevent automatically discovered parent directories from entering cache metadata and later becoming recursive cleanup targets.
- Preserve explicit directory artifacts declared by Steps, StepReferences, or runtime artifact registration.
- Ensure cache entries produced with the old directory-discovery behavior are never restored by the new runtime.
- Preserve repository OpenSpec roots and planning data across installation cache misses, cache hits, and final cleanup.
- Keep shell snapshots and Go filesystem snapshots behaviorally aligned.

**Non-Goals:**

- Add an OpenSpec-specific exclusion to the framework cleanup Step.
- Change the public Step, StepReference, Cmd, or CmdWorkflow manifest schemas.
- Prevent cleanup of a directory explicitly registered by a manifest or Step implementation.
- Manage OpenSpec store registration, synchronization, worksets, or remote state.
- Recover planning data already deleted by an older session.

## Decisions

### 1. Automatic filesystem discovery records leaf paths only

The generated pre-Step snapshot and the Go post-Step snapshot will contain files and symbolic links but omit directories. Filesystem diff therefore records `openspec/config.yaml` or generated skill files, not `openspec` or `.opencode` merely because they are parent directories.

Explicit artifact sources remain unchanged. A directory listed in `spec.artifacts`, StepReference `artifacts`, or passed to `__kfg_add_artifact` remains a directory artifact and retains recursive cleanup semantics. This preserves the framework contract for disposable scaffolds while making best-effort automatic discovery conservative.

**Alternative considered:** Add `openspec/` to a cleanup denylist. This would leak domain knowledge into the framework, leave every other persistent directory vulnerable, and make intentionally disposable OpenSpec fixtures impossible to declare.

**Alternative considered:** Change `kfg.cleanup` to remove only empty directories. That would weaken explicit directory artifacts and leave generated scaffold content behind when a Step does not register every child.

### 2. Version the cache identity namespace

`ComputeIdentity` will include an engine-controlled cache schema prefix before the StepReference name, for example `SHA256("v2\x00" + stepRefName)`. All cache commands continue to derive paths through this function. Existing entries remain on disk for normal cache garbage collection but cannot satisfy `exists` or `restore` under the new namespace.

This one-time invalidation is required because current metadata stores only artifact paths. It cannot distinguish an automatically discovered `openspec` directory from an explicitly declared directory, so selectively migrating legacy entries would be unsafe.

**Alternative considered:** Add a metadata version while retaining the same cache path. A version mismatch would produce a miss, but storing the replacement would collide with the existing non-empty entry unless the store path also gained replacement and rollback logic. Namespacing identity is smaller and keeps atomic store behavior unchanged.

**Alternative considered:** Rename only the OpenSpec workflow StepReferences. That would not protect custom workflows or other persistent directories affected by the same engine behavior.

### 3. Preserve the existing cleanup and manifest contracts

No fields or resources are added. Manifest composition, merge behavior, dependency resolution, `when` evaluation, and Step execution order are unchanged. The generated shell template changes only the filesystem paths written to the temporary pre-Step snapshot. The Go cache store continues to merge runtime delta, declarative artifacts, and filesystem-diff artifacts before copying them atomically.

The OpenSpec capability remains composed through its existing Kustomizations and `openspec.steps.install`. Its adjacent README will state that repository planning roots are persistent and that cleanup is limited to concrete generated integration artifacts unless a consumer explicitly declares a directory artifact.

## Risks / Trade-offs

- [Empty directories created only by a cached Step are no longer automatically cached or cleaned] -> Treat empty directories as non-observable scaffolding unless the Step explicitly declares or registers them.
- [All existing Step cache entries miss once after the namespace change] -> Rebuild entries automatically on first use and leave old entries available to existing GC commands.
- [A Step may still explicitly register a persistent directory incorrectly] -> Preserve explicit semantics by design and cover the OpenSpec workflow so it never declares the planning root as an artifact.
- [Shell and Go path selection may drift] -> Add generated-template and Go snapshot tests for nested files, parent directory exclusion, and symbolic links.
- [A newly generated OpenSpec configuration file may still be cleaned] -> Cleanup may remove the exact file generated during the session, but it cannot recursively remove sibling specs or changes; pre-existing files are absent from the diff and remain untouched.

## Migration Plan

1. Introduce the versioned cache identity and leaf-only filesystem snapshot semantics together so no legacy directory-bearing entry can be restored by the new runtime.
2. Rebuild affected cache entries naturally on their next invocation.
3. Validate OpenSpec installation through both miss and hit paths before running final cleanup.
4. Roll back by reverting the engine and domain changes together. A rollback runtime will continue to see its original unversioned cache entries; users should clear those entries before rollback if they may contain persistent directory artifacts.

## Open Questions

None.

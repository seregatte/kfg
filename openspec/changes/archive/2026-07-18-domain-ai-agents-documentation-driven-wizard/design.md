## Context

The AI agents domain currently has one broad aggregate entrypoint plus several integration directories, but only two package-local READMEs. The wizard compensates by embedding a catalog of known agents, Steps, Assets, converters, and example paths in its prompt. That catalog already diverges from the manifests: some installation Steps are not wired into either overlay, converter support differs by agent, paths and resource names have changed, and new capabilities require editing the wizard even when the underlying domain contract is otherwise complete.

OpenSpec is affected by the same coupling. The repository stores canonical planning artifacts under `docs/context/openspec/`, exports `OPENSPEC_ROOT_DIR=docs/context`, and wraps the CLI with a directory change. Meanwhile, `openspec init` creates an ignored `openspec/` root at the repository root. OpenSpec 1.6 resolves the nearest standard root natively and distinguishes a normal repository root from standalone stores and read-only references, so the current wrapper can select different planning data depending on how it is invoked.

The change spans domain packaging, agent instructions, repository documentation, OpenSpec planning data, Nix shell behavior, and Bats validation. Existing aggregate Kustomizations remain public and must continue to build while selective consumption becomes the preferred generated form.

## Goals / Non-Goals

**Goals:**

- Make capability documentation the discoverable human and agent-facing catalog for the AI agents domain.
- Keep Kustomizations authoritative for composition and manifests authoritative for defaults.
- Give every independently consumable capability a documented public entrypoint and stable JSON Patch targets.
- Make the wizard discover and verify capabilities instead of duplicating their resource inventory.
- Generate minimal, compatibility-aware compositions that consumers can maintain through overlays.
- Represent OpenSpec local planning, external store pointers, and read-only references without managing machine-local store state.
- Move kfg's canonical OpenSpec data to the standard repository-root location and rely on native CLI root resolution.
- Enforce documentation synchronization through repository instructions and automated domain tests.

**Non-Goals:**

- Define a new kfg manifest kind for OpenSpec stores.
- Automatically clone, register, synchronize, push, or remove OpenSpec stores.
- Set OpenSpec's global `defaultStore` on behalf of a consumer.
- Give every configuration value its own variant or overlay.
- Preserve `docs/context/openspec/` as a compatibility alias or symlink.
- Guarantee every capability works with every supported agent; unsupported combinations are documented and reported.
- Replace canonical OpenSpec requirements with README documentation.

## Decisions

### 1. Use a three-layer capability contract

Each public capability uses an adjacent README for selection and usage, a Kustomization for the exported resource graph, and manifests for defaults. The domain root README indexes capability READMEs and differentiates the complete package, the domain-only aggregate, selective capability entrypoints, and operational overlays.

This avoids parsing implementation files as the wizard's primary user interface while preventing documentation from becoming the composition source of truth. The wizard reads documentation to understand intent and compatibility, then verifies resource names and availability against the adjacent Kustomization before proposing output.

**Alternative considered:** Generate a machine-readable central catalog. This would create another artifact that must remain synchronized and would continue the current duplication problem.

### 2. Add selective entrypoints without removing aggregates

Every independently useful integration and agent exposes a directory-level `kustomization.yaml`. Independently selectable prompts, subagents, commands, and similar building blocks receive selective entrypoints, while their existing category Kustomizations continue to aggregate all resources. The root domain entrypoint remains available for consumers that explicitly want the full catalog.

Named variants are introduced only when resource composition or semantics differ. Ordinary values such as models, enabled flags, arguments, output paths, and reference IDs are customized with Kustomize JSON Patch.

**Alternative considered:** Move all consumers immediately to selective entrypoints and remove aggregates. Preserving aggregates limits migration risk and keeps existing downstream Kustomizations valid.

### 3. Standardize capability README content

Each capability README documents purpose, prerequisites, entrypoints, exported resource names and kinds, agent compatibility, composition, configurable defaults, JSON Patch examples, workflow references, generated artifacts, cache and cleanup behavior, limitations, validation, and canonical specs. Category and domain READMEs act as indexes rather than repeating detailed contracts.

`docs/AGENTS.md` makes updating the adjacent README mandatory whenever a public resource name, default, configurable field, compatibility rule, artifact, cleanup behavior, workflow usage, or entrypoint changes. Tests verify README presence, required sections, advertised resources, entrypoint builds, and representative patches.

**Alternative considered:** Maintain only a single comprehensive domain README. Per-capability documentation keeps changes local and lets the wizard load only relevant context.

### 4. Redesign wizard discovery around the domain index

The wizard locates the installed domain through `KFG_DIR`, reads the domain index, follows README links for candidate capabilities, and verifies selected Kustomizations. Its prompt retains interaction and confirmation safeguards but removes hardcoded resource lists and assumed compatibility.

Generated configurations prefer selective entrypoints. Before writing, the wizard presents selected entrypoints, patches, artifacts, prerequisites, unsupported combinations, and every output path. After confirmation, it validates the generated Kustomization and reports failures without inventing replacement resources.

**Alternative considered:** Refresh the embedded catalog during each domain change. This still requires synchronized edits and prevents third-party or newly added capabilities from being discovered naturally.

### 5. Model OpenSpec as two exclusive roots plus an additive reference set

The OpenSpec package exposes common command and installation resources plus three configuration Kustomizations:

- `local` materializes a normal `openspec/config.yaml` for planning beside code.
- `store` materializes a config-only `store: <id>` pointer for fully externalized planning.
- `references` contributes read-only store references and is primarily composed with `local`.

Local and store-pointer modes are mutually exclusive. References are additive and do not change command write targets. Assets expose stable patch paths for project context, rules, store IDs, reference IDs, and remotes. Materialization deep-merges selected configuration Assets into `openspec/config.yaml` through an identity YAML converter and the framework materialization contract.

Store identity remains in the standalone store's `.openspec-store/store.yaml`; registry and workset data remain machine-local. The README provides explicit `openspec store setup` and `register` prerequisites instead of executing them from a cached Step or shell hook.

**Alternative considered:** Materialize `.openspec-store/store.yaml` in every consumer. That incorrectly changes a code repository into a standalone planning store and conflates versioned declarations with local registration.

### 6. Migrate kfg to one standard OpenSpec root

All tracked contents of `docs/context/openspec/`, including active changes, archives, configuration, and specs, move to repository-root `openspec/`. The `/openspec` ignore rule is removed. Repository documentation, scripts, tests, and canonical specs are updated atomically.

The feature change is initially proposed under the current root to comply with the repository's pre-migration workflow. During implementation the active change moves with the full root and is ultimately archived under `openspec/changes/archive/`.

No symlink or fallback path is retained because two qualifying roots make OpenSpec behavior dependent on the current directory. Rollback is a Git revert of the atomic migration.

**Alternative considered:** Keep the nested root and rely on `OPENSPEC_ROOT_DIR`. This preserves custom wrapper behavior that native OpenSpec commands, agents, and external tooling do not consistently honor.

### 7. Make the OpenSpec wrapper and installer root-neutral

`ai.cmds.openspec` delegates directly to `command openspec "$@"`; native resolution selects an explicit `--store`, the nearest local root, a declared store pointer, or configured fallback according to OpenSpec semantics. Development and CI shells stop exporting `OPENSPEC_ROOT_DIR`.

The OpenSpec installation Step remains responsible for agent integration files but must not delete `openspec/` or leave a disposable root that shadows canonical data. Its inputs and documented artifacts are aligned with what the current CLI actually generates for OpenCode and Pi.

**Alternative considered:** Keep the wrapper only for kfg development. A wrapper-specific root makes generated shells behave differently from direct OpenSpec usage and obscures failures in consumer configurations.

### 8. Test documentation as part of the public API

Domain Bats tests build aggregate and selective entrypoints, apply representative JSON Patches from fixtures, assert adjacent README coverage and required sections, and compare documented resource names with built manifests. OpenSpec tests cover local, pointer, and referenced configurations plus native `doctor`, `context`, and validation behavior. Wizard tests assert discovery instructions and absence of a hardcoded capability catalog.

The test suite validates contracts rather than exact prose so documentation can improve without brittle snapshot churn.

**Alternative considered:** Review documentation manually. The number of public capabilities and prior drift make manual enforcement insufficient.

## Risks / Trade-offs

- [Risk] The scope touches many files and may expose unrelated stale resource names or unsupported converter paths. → Mitigation: inventory capabilities first, preserve aggregates, and update one capability package at a time with focused build tests.
- [Risk] README-based discovery could trust stale documentation. → Mitigation: require adjacent Kustomization verification and test documented resource names and patch targets.
- [Risk] Moving the OpenSpec root creates a large rename diff and can conflict with concurrent planning changes. → Mitigation: perform the move once after planning artifacts are complete, preserve all active changes, and avoid parallel feature work touching the old root.
- [Risk] OpenSpec stores and references are beta and may change format. → Mitigation: isolate their Assets and entrypoints, document the beta status, and avoid managing the local registry.
- [Risk] Selective entrypoints increase the number of public paths. → Mitigation: use a consistent directory contract, retain category indexes, and test every advertised entrypoint.
- [Risk] JSON Patch examples can become invalid after manifest restructuring. → Mitigation: represent examples as test fixtures or equivalent build assertions.
- [Trade-off] Per-capability READMEs add maintenance work. → The explicit maintenance cost replaces hidden wizard and consumer breakage caused by undocumented changes.
- [Trade-off] The wizard performs more filesystem reads. → It loads only the domain index and selected capability documentation, reducing prompt context compared with embedding the full catalog.

## Migration Plan

1. Add the documentation contract specs and test helpers while the current OpenSpec root remains authoritative.
2. Inventory and normalize public capability entrypoints, then add domain and capability READMEs without removing aggregate paths.
3. Add OpenSpec common resources and local, store-pointer, and references entrypoints with fixture-based patch tests.
4. Rewrite the wizard instructions to discover the domain index and selected capability contracts.
5. Update `docs/AGENTS.md` with the documentation synchronization rule and the pending standard OpenSpec path.
6. Move `docs/context/openspec/` atomically to `openspec/`, remove the ignore rule and environment redirects, and update all repository references.
7. Simplify the OpenSpec wrapper and installer, then validate native root discovery from the repository and nested directories.
8. Run domain, Bats, Go, formatting, lint, and vet checks; archive this change in the migrated root.

Rollback uses Git to revert the complete change. Partial rollback of only the root move is not supported because documentation, shell behavior, and wizard output will reference the new location.

## Open Questions

None. The public contract, OpenSpec mode boundaries, wizard discovery source, and migration direction are defined for implementation.

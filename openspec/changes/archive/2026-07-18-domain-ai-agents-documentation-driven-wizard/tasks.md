## 1. Establish the Domain Documentation Contract

- [x] 1.1 Inventory every public AI agents capability, resource name, Kustomization entrypoint, supported agent, configurable field, generated artifact, cache setting, and cleanup behavior against the current manifests.
- [x] 1.2 Define and add the required capability README section template used by domain documentation tests.
- [x] 1.3 Create `packages/domains/ai-agents/README.md` as the public index for complete, aggregate, selective, and overlay entrypoints.
- [x] 1.4 Update `packages/domains/ai-agents/manifests/README.md` to describe the current package structure and link to capability documentation instead of obsolete extension paths.
- [x] 1.5 Add the mandatory AI agents capability documentation synchronization rule to `docs/AGENTS.md`.

## 2. Document Agents and Shared Building Blocks

- [x] 2.1 Add OpenCode capability documentation covering settings, command, converters, supported materialization targets, JSON Patch examples, artifacts, and limitations.
- [x] 2.2 Add Pi capability documentation covering settings, command, converter support, JSON Patch examples, artifacts, and documented compatibility gaps.
- [x] 2.3 Add selective Kustomization entrypoints and adjacent documentation for independently consumable commands.
- [x] 2.4 Add selective Kustomization entrypoints and adjacent documentation for independently consumable shared Steps.
- [x] 2.5 Add selective Kustomization entrypoints and adjacent documentation for each prompt while retaining the all-prompts aggregate.
- [x] 2.6 Add selective Kustomization entrypoints and adjacent documentation for each subagent while retaining the all-subagents aggregate.
- [x] 2.7 Add selective Kustomization entrypoints and adjacent documentation for shared and agent-specific converters while retaining existing aggregates.
- [x] 2.8 Verify converter documentation reflects actual emitted fields and explicitly reports unsupported tools, permissions, MCP metadata, or Pi conversion targets.

## 3. Document and Normalize Integration Capabilities

- [x] 3.1 Add Chrome DevTools documentation for MCP Assets, skill installation, workflow usage, JSON Patch targets, artifacts, and cleanup.
- [x] 3.2 Add Context7 documentation for MCP Assets, skill installation, context output and injection, workflow usage, JSON Patch targets, and cleanup.
- [x] 3.3 Add Google Workspace documentation for installation inputs, supported agents, artifacts, cache behavior, workflow usage, and JSON Patch targets.
- [x] 3.4 Add NotebookLM documentation for installation inputs, supported agents, output paths, artifacts, cache behavior, and JSON Patch targets.
- [x] 3.5 Add Playwright documentation for MCP Assets, skill installation, workflow usage, JSON Patch targets, artifacts, and cleanup.
- [x] 3.6 Add README documentation for `overlays/ai` and deprecated `overlays/dev`, including composition, generated artifacts, cleanup, and migration guidance.
- [x] 3.7 Correct installation Steps that assume undocumented agent homes or remove files outside their registered artifacts.
- [x] 3.8 Align installation Step inputs, outputs, resource names, cache behavior, and workflow references with their documented contracts.

## 4. Add Modern OpenSpec Composition

- [x] 4.1 Restructure the OpenSpec capability into common resources plus documented converter, Step, and configuration entrypoints.
- [x] 4.2 Add an identity YAML converter and `openspec.assets.config.local` with patchable schema, context, and artifact rules.
- [x] 4.3 Add the local planning Kustomization and materialization example for `openspec/config.yaml`.
- [x] 4.4 Add `openspec.assets.config.store` and the mutually exclusive external-store pointer Kustomization.
- [x] 4.5 Add `openspec.assets.config.references` and the additive read-only references Kustomization with patchable IDs and remotes.
- [x] 4.6 Add the OpenSpec README covering root precedence, local planning, stores, pointers, references, versioned versus machine-local state, beta limitations, JSON Patch, setup prerequisites, and validation.
- [x] 4.7 Simplify `ai.cmds.openspec` to forward arguments directly to the OpenSpec binary without `OPENSPEC_ROOT_DIR` or directory changes.
- [x] 4.8 Update `openspec.steps.install` to install documented agent integration files without deleting or shadowing a project `openspec/` root.

## 5. Redesign the Wizard Around Capability Discovery

- [x] 5.1 Replace the wizard's hardcoded domain catalog with instructions to locate `KFG_DIR` and read the domain capability index.
- [x] 5.2 Require the wizard to read selected capability READMEs and verify their adjacent Kustomizations before proposing resources.
- [x] 5.3 Add compatibility-aware discovery questions and explicit handling for unsupported agent and capability combinations.
- [x] 5.4 Make generated projects prefer selective entrypoints and documented JSON Patches over complete-domain composition or resource duplication.
- [x] 5.5 Add OpenSpec discovery questions for local planning, external store pointers, and read-only references without automatic store registration.
- [x] 5.6 Expand the pre-write summary to include entrypoints, patches, prerequisites, artifacts, cleanup behavior, limitations, and output paths.
- [x] 5.7 Require the wizard to build generated Kustomizations after writing and report validation failure accurately.
- [x] 5.8 Update the wizard command and skill documentation to describe the documentation-driven discovery flow.

## 6. Migrate the Repository OpenSpec Root

- [x] 6.1 Normalize `docs/context/openspec/config.yaml` artifact rules to the current OpenSpec array schema while preserving project context.
- [x] 6.2 Move all tracked configuration, specs, active changes, and archives from `docs/context/openspec/` to repository-root `openspec/`, including this active change.
- [x] 6.3 Remove `/openspec` from `.gitignore` and ensure no nested compatibility root or symlink remains.
- [x] 6.4 Remove `OPENSPEC_ROOT_DIR` from `devShells.dev` and `devShells.ci` without changing the flake version.
- [x] 6.5 Update `docs/AGENTS.md`, root documentation, tests, manifests, and specs to use the canonical `openspec/` paths and workflow.
- [x] 6.6 Preserve unrelated active and archived OpenSpec changes unchanged during the root move.
- [x] 6.7 Verify `openspec doctor`, `openspec context`, list, status, and validate resolve the repository root from both the root and nested directories.

## 7. Add Contract and Composition Tests

- [x] 7.1 Add Bats coverage requiring adjacent READMEs and mandatory sections for all advertised public capabilities.
- [x] 7.2 Add Bats coverage comparing documented exported resource names with each capability's built manifests.
- [x] 7.3 Add Bats coverage that builds every selective, category aggregate, domain aggregate, and overlay Kustomization.
- [x] 7.4 Add fixture overlays validating representative documented JSON Patch targets and paths.
- [x] 7.5 Add OpenSpec fixtures for local, external-store pointer, and local-with-references composition and materialization.
- [x] 7.6 Add command tests proving `ai.cmds.openspec` forwards normal and `--store` arguments without changing directories.
- [x] 7.7 Add installation tests proving OpenSpec and other normalized Steps preserve unrelated project and agent files.
- [x] 7.8 Add wizard tests requiring documentation discovery and Kustomization verification while rejecting a hardcoded authoritative catalog.
- [x] 7.9 Add a discovery fixture proving a newly indexed capability can be selected without editing the wizard prompt.

## 8. Verify and Finalize

- [x] 8.1 Run `nix develop .#dev --command make build`.
- [x] 8.2 Run `nix develop .#dev --command make test`.
- [x] 8.3 Run `nix develop .#dev --command make test-bats`.
- [x] 8.4 Run `nix develop .#dev --command make fmt lint vet`.
- [x] 8.5 Run OpenSpec validation for the migrated root, all durable specs, and this change.
- [x] 8.6 Review the final diff for stale `docs/context/openspec`, `OPENSPEC_ROOT_DIR`, obsolete resource names, undocumented public capabilities, and unintended version changes.

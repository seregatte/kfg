## 1. OpenSpec And Manifest Contracts

- [x] 1.1 Review the new change artifacts and confirm the final capability scope for extension MCP assets, install steps, dev workflow, ctx7 normalization, and LifeOS overlay migration.
- [x] 1.2 Update `.manifests/base/extensions/ctx7/` so `kfg.extension.ctx7.install` validates explicit inputs, copies skills into `OUTPUT_DIR`, and aligns with the normalized extension contract.
- [x] 1.3 Update `.manifests/base/extensions/openspec/` so `kfg.extension.openspec.install` makes copied outputs explicit and aligns with the normalized install-step contract.

## 2. Normalize MCP-First Extensions

- [x] 2.1 Update `.manifests/base/extensions/chrome-devtools/assets/mcp.yaml` to match the canonical MCP command and metadata from the legacy image reference behavior.
- [x] 2.2 Update `.manifests/base/extensions/chrome-devtools/steps/install.yaml` to preserve the legacy skill-install behavior with explicit inputs, output handling, and deterministic success output.
- [x] 2.3 Update `.manifests/base/extensions/playwright/assets/mcp.yaml` to match the canonical MCP command and metadata from the legacy image reference behavior.
- [x] 2.4 Update `.manifests/base/extensions/playwright/steps/install.yaml` to preserve the legacy Playwright skill-install behavior, including temporary agent-home handling and Pi-specific MCP exclusion.

## 3. Normalize Install-Only Extensions

- [x] 3.1 Update `.manifests/base/extensions/gws/steps/install.yaml` to preserve the legacy GWS skill-install behavior with explicit validation, output handling, and deterministic success output.
- [x] 3.2 Update `.manifests/base/extensions/notebooklm/steps/install.yaml` to preserve the legacy NotebookLM skill-install behavior, including copying generated skill files into overlay-owned output paths.
- [x] 3.3 Review extension kustomizations and resource names so each normalized extension is consumed through `kfg.extension.<name>.*` resources only.

## 4. Update Shared Overlay Consumption

- [x] 4.1 Update `.manifests/overlay/dev/agents-workflow.yaml` to aggregate canonical extension-owned MCP assets instead of transitional `kfg.extension.self.*` MCP assets.
- [x] 4.2 Update `.manifests/overlay/dev/agents-workflow.yaml` to keep extension install phases ordered before context injection, command conversion, subagent conversion, and MCP aggregation.
- [x] 4.3 Remove or stop referencing transitional extension assets and assumptions that are no longer canonical after the normalization.

## 5. Migrate LifeOS Overlay

- [x] 5.1 Create `~/Sites/lifeos/.manifests/overlay/lifeos/` with `kustomization.yaml`, `cmds.yaml`, and `agents-workflow.yaml` following the normalized overlay pattern.
- [x] 5.2 Add overlay-local LifeOS assets for commands, subagents, and MCP definitions based on the current LifeOS content while keeping `project/` access policy intact.
- [x] 5.3 Wire the LifeOS workflow to shared `kfg` steps, canonical extension-owned MCP assets, and shared converters.
- [x] 5.4 Update the LifeOS manifest entrypoint so execution from `~/Sites/lifeos` resolves to the normalized overlay.

## 6. Validation And Documentation

- [x] 6.1 Add or update unit and integration coverage for normalized install-step behavior and MCP aggregation paths where repository tests exist.
- [x] 6.2 Add or update Bats coverage for the dev workflow and manifest resources affected by extension normalization.
- [x] 6.3 Validate manifest build and shell generation behavior in `kfg` for the normalized extensions and overlay workflow.
- [x] 6.4 Validate the migrated LifeOS overlay from `~/Sites/lifeos` without modifying anything under `$HOME/.nixai`.
- [x] 6.5 Update relevant documentation to reflect the canonical extension contract, namespace usage, and LifeOS overlay migration pattern.

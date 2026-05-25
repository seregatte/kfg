## Why

The ctx7 install Step currently depends on repeated workflow-level artifact declarations to make cacheable installs restore the generated skill paths correctly. That makes the development overlay noisy and ties artifact tracking to manual path lists even though the Step already knows which output directory it populates.

## What Changes

- Update the ctx7 install Step to discover newly created skill directories under `OUTPUT_DIR` dynamically.
- Register those discovered paths as artifacts through the runtime API so cache persistence and cleanup continue to work.
- Remove redundant workflow-level artifact declarations for ctx7 install references once the Step owns artifact discovery.
- Add or update domain tests for ctx7 install artifact discovery and cache behavior.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `ctx7-install`: ctx7 install dynamically discovers produced skill artifacts from `OUTPUT_DIR` instead of relying on repeated workflow-level artifact lists.

## Impact

- Affects `packages/domains/ai-agents/manifests/ctx7/steps/install.yaml`.
- Affects ctx7 references in `packages/domains/ai-agents/overlays/dev/agents-workflow.yaml`.
- Depends on the engine runtime exposing portable filesystem snapshot/diff helpers.
- Requires domain-level validation that ctx7 install still produces the expected skills and cacheable artifacts.

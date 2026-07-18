# Google Workspace

## Purpose

Google Workspace skills integration for Docs, Sheets, Slides, and Gmail.

## Prerequisites

- **Node.js** — Required for `npx` execution
- **npx** — Package runner for installing skills

## Exported Resources

| Resource Name | Kind | Description |
|---------------|------|-------------|
| `gws.steps.install` | Step | Installs Google Workspace skills via `npx skills add` |

## Selective Entrypoints

| Entrypoint | Description |
|------------|-------------|
| `steps/` | Installation step only |

## Aggregate Entrypoint

| Entrypoint | Description |
|------------|-------------|
| `kustomization.yaml` | All Google Workspace resources |

## Agent Compatibility

| Agent | Supported |
|-------|-----------|
| opencode | Yes |
| pi | No |

The install step hardcodes `--agent opencode` in the `AGENT_FLAG` environment
variable, limiting this capability to OpenCode agents only.

## Configuration

| Input | Default | Description |
|-------|---------|-------------|
| `SKILL_NAME` | `googleworkspace/cli` | npm package name of the skill to install |
| `AGENT_FLAG` | `opencode` | Target agent for skill installation |

Both inputs are set as environment variables in the step spec and are not
configurable via JSON Patch defaults.

### No Configurable Defaults

The step has no configurable `spec.defaults` field. The environment variables
are hardcoded in the step definition. To customize, override the step's `env`
block via JSON Patch:

```yaml
# kustomization.yaml
patches:
  - target:
      kind: Step
      name: gws.steps.install
    patch: |
      - op: replace
        path: /spec/env/SKILL_NAME
        value: "googleworkspace/custom-skill"
```

## Workflow Usage

Available but **not wired into either overlay by default**. To use in a
workflow, reference the step directly:

```yaml
steps:
  - name: Install Google Workspace
    ref: gws.steps.install
    weight: 10
```

## Generated Artifacts

| Artifact | Description |
|----------|-------------|
| `.agents/` | Directory containing installed skill files |

The step creates an `.agents` directory in the working directory with the
installed Google Workspace skill definitions.

## Cache and Cleanup

| Feature | Status | Description |
|---------|--------|-------------|
| Cache | Enabled | Results are cached to avoid re-installing unchanged skills |
| Cleanup | Automatic | Temporary files from `npx` execution are removed after install |

## Limitations

- **OpenCode only** — The `AGENT_FLAG` input is hardcoded to `opencode`,
  making this capability incompatible with other agents.
- **Google Workspace account required** — Skills require an active Google
  Workspace account for API access.

## Validation

```bash
nix develop .#dev --command kustomize build packages/domains/ai-agents/manifests/gws
```

## Canonical Specs

- [kfg-domain-package-contract](../../../../../docs/context/openspec/specs/kfg-domain-package-contract/spec.md)

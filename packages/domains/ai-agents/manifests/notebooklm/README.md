# NotebookLM

## Purpose

Google NotebookLM integration for AI-powered research and note-taking.

## Prerequisites

- `notebooklm` CLI installed and available on `$PATH`

## Selective Entrypoints

| Entrypoint | Description |
|------------|-------------|
| `kustomization.yaml` | NotebookLM package (selective) |

## Aggregate Entrypoint

| Entrypoint | Description |
|------------|-------------|
| `steps/kustomization.yaml` | All NotebookLM steps |

## Exported Resources

| Resource Name | Kind | Description |
|---------------|------|-------------|
| `notebooklm.steps.install` | Step | Installs NotebookLM skill files into the output directory |

## Agent Compatibility

| Agent | Supported |
|-------|-----------|
| opencode | Yes |
| pi | Yes |

## Configuration

### Step Inputs (`notebooklm.steps.install`)

| Field | Default | Description |
|-------|---------|-------------|
| `SKILL_NAME` | — | Name of the skill to install (optional) |
| `AGENT_HOME` | — | Path to the agent home directory (optional) |
| `OUTPUT_DIR` | `.opencode/skills/` | Output directory for generated skill files |

## JSON Patch Examples

N/A — the step has no configurable defaults exposed via JSON Patch.

## Workflow Usage

Available but not wired into either the `dev` or `ai` overlay by default.
To include it, add the step to a workflow YAML:

```yaml
steps:
  - name: notebooklm.steps.install
    weight: 0
```

## Generated Artifacts

Skill files copied from `~/.opencode/skills/notebooklm/` into `$OUTPUT_DIR`
(typically `.opencode/skills/`).

## Cache and Cleanup

| Behavior | Detail |
|----------|--------|
| Cache | Enabled (`cache.enabled: true`) |
| Cleanup | Removes temporary directories after copy |

## Limitations

- Requires the `notebooklm` CLI to be installed separately outside of kfg.
- The step is a pass-through installer; no validation of CLI version or
  capability is performed.

## Validation

```bash
kustomize build packages/domains/ai-agents/manifests/notebooklm/
```

## Canonical Specs

- [kfg-domain-package-contract](../../../../docs/context/openspec/specs/domain-ai-agents-package-contract.md)

# Shared Steps

Workflow steps shared across AI agents and overlays.

## Resources

| Resource Name | Description |
|---------------|-------------|
| `ai.steps.detect` | Detects the active AI agent by echoing `$AGENT` |
| `ai.steps.deprecation-warn` | Displays a deprecation warning for `overlays/dev/` |

## Selective Entrypoints

| Entrypoint | Description |
|------------|-------------|
| `detect/` | Agent detection step only |
| `deprecation-warn/` | Deprecation warning step only |

## Aggregate Entrypoint

| Entrypoint | Description |
|------------|-------------|
| `kustomization.yaml` | All shared steps |

## Configuration

### Agent Detection (`ai.steps.detect`)

| Field | Default | Description |
|-------|---------|-------------|
| Output name | `AGENT` | Name of the output variable |
| Output type | `string` | Type of the output value |

The step reads the `AGENT` environment variable set by the Cmd wrapper and
outputs it for conditional branching in workflows.

### Deprecation Warning (`ai.steps.deprecation-warn`)

No configurable fields. Logs a warning message directing users to `kfg ai`.

## Workflow Usage

- `overlays/dev/agents-workflow.yaml` — Uses `ai.steps.deprecation-warn` (weight -100) and `ai.steps.detect` (weight -70) for conditional agent setup

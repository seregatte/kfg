# Shared Converters

Format converters shared across AI agents and overlays.

## Resources

| Resource Name | Input Format | Output Format | Expression Description |
|---------------|--------------|---------------|------------------------|
| `ai.conv.to-json` | YAML | JSON | Identity expression (`.`); passes input through unchanged |

## Agent Compatibility

Shared converters work with **all agents**. Unlike per-agent converters
(e.g., `ai.opencode.conv.command`), shared converters are not tied to a
specific agent's output format.

## Selective Entrypoints

| Entrypoint | Description |
|------------|-------------|
| `to-json/` | YAML-to-JSON converter only |

## Aggregate Entrypoint

| Entrypoint | Description |
|------------|-------------|
| `kustomization.yaml` | All shared converters |

## Configuration

Both fields are defined in the Converter spec:

| Field | Default | Description |
|-------|---------|-------------|
| `spec.engine.expression` | `.` | Expression applied to input; identity returns input unchanged |
| `spec.output.format` | `json` | Target output format |

### Changing the Expression

Use JSON Patch to override the expression at build time:

```yaml
# kustomization.yaml
patches:
  - target:
      kind: Converter
      name: ai.conv.to-json
    patch: |
      - op: replace
        path: /spec/engine/expression
        value: |
          .converted = true
```

## Limitations

- **Identity expression only** — the default `.` expression passes input
  through unchanged; no field mapping or renaming is performed.
- **No complex transformations** — the engine does not support conditional
  logic, aggregation, or multi-step pipelines. For advanced transformations,
  use per-agent converters that emit agent-specific formats.

## Workflow Usage

- `overlays/dev/agents-workflow.yaml` — Uses `ai.conv.to-json` for settings
  materialization; converts YAML agent settings to JSON before writing to
  agent config files

# AI Agents Domain

This package provides AI agents, integrations, workflows, and building blocks for
kfg-based projects. It composes the framework package and all domain-specific manifests.

## Quick Start

Apply the complete domain package:

```yaml
resources:
  - https://github.com/seregatte/kfg.git//packages/domains/ai-agents?ref=v0.0.12
```

Or apply the domain-only aggregate (without framework resources):

```yaml
resources:
  - https://github.com/seregatte/kfg.git//packages/domains/ai-agents/manifests?ref=v0.0.12
```

## Available Agents

| Agent | Documentation | Selective Entrypoint |
|-------|--------------|---------------------|
| OpenCode | [OpenCode](manifests/agents/opencode/README.md) | `manifests/agents/opencode/kustomization.yaml` |
| Pi | [Pi](manifests/agents/pi/README.md) | `manifests/agents/pi/kustomization.yaml` |

## Available Integrations

| Integration | Documentation | Selective Entrypoint |
|-------------|--------------|---------------------|
| Chrome DevTools | [Chrome DevTools](manifests/chrome-devtools/README.md) | `manifests/chrome-devtools/kustomization.yaml` |
| Context7 | [Context7](manifests/ctx7/README.md) | `manifests/ctx7/kustomization.yaml` |
| Google Workspace | [Google Workspace](manifests/gws/README.md) | `manifests/gws/kustomization.yaml` |
| NotebookLM | [NotebookLM](manifests/notebooklm/README.md) | `manifests/notebooklm/kustomization.yaml` |
| Playwright | [Playwright](manifests/playwright/README.md) | `manifests/playwright/kustomization.yaml` |
| OpenSpec | [OpenSpec](manifests/openspec/README.md) | `manifests/openspec/kustomization.yaml` |

## Available Building Blocks

| Category | Documentation | Selective Entrypoint |
|----------|--------------|---------------------|
| Commands | [Commands](manifests/cmds/README.md) | `manifests/cmds/kustomization.yaml` |
| Steps | [Steps](manifests/steps/README.md) | `manifests/steps/kustomization.yaml` |
| Prompts | [Prompts](manifests/prompts/README.md) | `manifests/prompts/kustomization.yaml` |
| Subagents | [Subagents](manifests/subagents/README.md) | `manifests/subagents/kustomization.yaml` |
| Converters | [Converters](manifests/converters/README.md) | `manifests/converters/kustomization.yaml` |

## Category Aggregates

| Aggregate | Entrypoint |
|-----------|------------|
| All agents | `manifests/agents/kustomization.yaml` |
| All commands | `manifests/cmds/kustomization.yaml` |
| All shared steps | `manifests/steps/kustomization.yaml` |
| All prompts | `manifests/prompts/kustomization.yaml` |
| All subagents | `manifests/subagents/kustomization.yaml` |
| All converters | `manifests/converters/kustomization.yaml` |

## Overlays

| Overlay | Status | Entrypoint | Documentation |
|---------|--------|------------|---------------|
| AI Wizard | Active | `overlays/ai/kustomization.yaml` | [README](overlays/ai/README.md) |
| Development | Deprecated | `overlays/dev/kustomization.yaml` | [README](overlays/dev/README.md) |

## Complete Package

The root Kustomization at `kustomization.yaml` composes the framework package and all
domain manifests into a single entrypoint. This provides every agent, integration,
command, step, prompt, subagent, and converter defined in this domain.

```yaml
resources:
  - kustomization.yaml  # Complete package
```

## Configuration

Most capabilities support customization via Kustomize JSON Patch. Each capability's
README documents its configurable fields, defaults, and patch examples.

## Documentation Contract

Every independently consumable capability has an adjacent README with standardized
sections covering purpose, prerequisites, entrypoints, exported resources, agent
compatibility, configuration, JSON Patch examples, workflow usage, generated artifacts,
cache and cleanup behavior, limitations, validation, and canonical specs.

The wizard discovers capabilities through these README files and their adjacent
Kustomizations rather than maintaining a hardcoded catalog.

## Adding New Capabilities

1. Create the capability directory under `manifests/`
2. Add a `kustomization.yaml` entrypoint
3. Add an adjacent `README.md` following the [template](CAPABILITY-README-TEMPLATE.md)
4. Register the entrypoint in the parent Kustomization
5. Add to this index

See [manifests/README.md](manifests/README.md) for the manifest structure conventions.

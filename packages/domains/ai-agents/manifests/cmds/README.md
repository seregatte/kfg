# Commands

Shared command wrappers that launch AI agents or external tools.

## Resources

| Resource Name | Command | Description |
|---------------|---------|-------------|
| `ai.opencode.cmd.main` | `opencode` | Launches the OpenCode agent |
| `ai.pi.cmd.main` | `pi` | Launches the Pi agent |
| `ai.cmds.openspec` | `openspec` | Wraps the OpenSpec CLI |

## Selective Entrypoints

| Entrypoint | Description |
|------------|-------------|
| `agents/opencode/` | OpenCode command only |
| `agents/pi/` | Pi command only |
| `openspec/` | OpenSpec command only |

## Aggregate Entrypoint

| Entrypoint | Description |
|------------|-------------|
| `kustomization.yaml` | All commands |

## Configuration

Commands do not have configurable fields. Each command delegates to the underlying
binary and sets the `AGENT` environment variable.

## JSON Patch Examples

Commands do not support JSON Patch. To customize command behavior, modify the
agent settings through the agent-specific entrypoints.

## Workflow Usage

Commands are used by:
- `overlays/ai/ai-workflow.yaml` — Invokes `ai.pi.cmd.main` or `ai.opencode.cmd.main`
- `overlays/dev/agents-workflow.yaml` — Invokes all three commands

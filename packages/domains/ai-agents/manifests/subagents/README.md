# Subagents

Shared subagent definitions materialized as agent-specific formats during overlay
processing.

## Resources

| Resource Name | Description | Model |
|---------------|-------------|-------|
| `ai.subagents.review-minimal` | Lightweight code review for quick pull request feedback | `claude-sonnet-4-20250514` |

## Agent Compatibility

| Agent | Supported | Notes |
|-------|-----------|-------|
| OpenCode | Yes | Converted via `ai.opencode.conv.subagent` |
| Pi | No | No subagent converter available |

## Selective Entrypoints

| Entrypoint | Description |
|------------|-------------|
| `review-minimal/` | Review-minimal subagent only |

## Aggregate Entrypoint

| Entrypoint | Description |
|------------|-------------|
| `kustomization.yaml` | All subagents |

## Configuration

All subagent configuration (model, tools, permissions) lives in the asset's `spec.data`
field and is configurable via JSON Patch on the resource.

| Field | Default | Description |
|-------|---------|-------------|
| `model` | `claude-sonnet-4-20250514` | Model to use for the subagent |
| `tools` | `["Read", "Bash"]` | Tools available to the subagent (documentation only) |
| `permissions.allow` | `["Read(/**)", "Bash(git diff*)"]` | Allowed permission patterns |

### Permission Materialization

The `ai.opencode.conv.subagent` converter translates `permissions.allow` entries
into OpenCode's `permission` frontmatter block. Each `Tool(pattern)` entry
produces a broad `"*": "deny"` rule followed by the specific allow pattern:

```yaml
# Input permissions.allow:
#   - Read(/**)
#   - Bash(git diff*)
#
# Produces:
permission: {
  "Read": {
    "*": "deny",
    "/**": "allow"
  },
  "Bash": {
    "*": "deny",
    "git diff*": "allow"
  }
}
```

The deny-* rule is emitted before the allow rule (deny all, then allow specific).
When no permissions are defined, `permission: {}` is emitted.

> **Note:** The `tools` field is reserved for future use and currently serves as
> documentation. The effective permission rules are driven entirely by
> `permissions.allow`.

## JSON Patch Examples

### Change model

```yaml
# patch.yaml
- op: replace
  path: /spec/data/model
  value: opencode/deepseek-v4-flash-free
```

### Add a tool

```yaml
# patch.yaml
- op: add
  path: /spec/data/tools/-
  value: Grep
```

## Workflow Usage

Used by `overlays/dev/` during materialization to convert subagent definitions into
agent-specific formats (e.g., OpenCode Markdown frontmatter files in `.opencode/agents/`).

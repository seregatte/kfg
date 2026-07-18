# Capability README Template

Every independently consumable AI agents capability MUST have an adjacent `README.md`
with the sections described below. This template defines the required structure.

## Required Sections

### 1. Title and Purpose

```markdown
# <Capability Name>

<One-paragraph description of what this capability provides and why it exists.>
```

### 2. Prerequisites

```markdown
## Prerequisites

<List of external tools, packages, or configurations required before using this capability.>

- Required binary or CLI
- Required environment variables
- Required kfg version
```

### 3. Entrypoints

```markdown
## Entrypoints

<List the Kustomization entrypoints for this capability.>

| Entrypoint | Type | Description |
|------------|------|-------------|
| `path/to/kustomization.yaml` | Selective | What this entrypoint provides |
```

### 4. Exported Resources

```markdown
## Exported Resources

<List every resource name and kind exported by this capability's Kustomization.>

| Resource Name | Kind | Description |
|---------------|------|-------------|
| `resource.name` | Assets/Step/Converter/Cmd | What it does |
```

### 5. Agent Compatibility

```markdown
## Agent Compatibility

<List which agents are supported and any limitations.>

| Agent | Supported | Notes |
|-------|-----------|-------|
| OpenCode | Yes | Any special behavior |
| Pi | Yes/No | Any limitations |
```

### 6. Configuration

```markdown
## Configuration

<Describe configurable fields with defaults.>

| Field | Default | Description |
|-------|---------|-------------|
| `field.name` | `default` | What it controls |
```

### 7. JSON Patch Examples

```markdown
## JSON Patch Examples

<Show how to customize the default configuration using Kustomize JSON Patch.>

### Example 1: <Description>

\`\`\`yaml
# patch.yaml
<JSON Patch document>
\`\`\`
```

### 8. Workflow Usage

```markdown
## Workflow Usage

<Describe how this capability is used within CmdWorkflow resources.>

This capability is used in:
- `overlays/ai/` — <how it is used>
- `overlays/dev/` — <how it is used>
```

### 9. Generated Artifacts

```markdown
## Generated Artifacts

<List files or directories created when this capability is applied.>

| Artifact | Path | Description |
|----------|------|-------------|
| File/Directory | `target/path` | What it contains |
```

### 10. Cache and Cleanup

```markdown
## Cache and Cleanup

<Describe caching behavior and what happens when the capability is removed.>

- **Cache**: Yes/No — <cache location and behavior>
- **Cleanup**: <what gets removed on cleanup>
```

### 11. Limitations

```markdown
## Limitations

<List known limitations, unsupported features, or caveats.>

- Limitation 1
- Limitation 2
```

### 12. Validation

```markdown
## Validation

<How to verify the capability is correctly installed/configured.>

\`\`\`bash
# Verify capability
<validation command>
\`\`\`
```

### 13. Canonical Specs

```markdown
## Canonical Specs

<Links to related OpenSpec specifications.>

- [spec-name](../../docs/context/openspec/specs/spec-name/spec.md)
```

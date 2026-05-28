## Context

AI-agent-specific resources are currently fragmented across multiple locations in the base manifests:

- `.manifests/base/agents/` - Agent settings assets (claude, gemini, opencode, pi) + to-json converter
- `.manifests/base/cmds/agents.yaml` - Agent command wrappers
- `.manifests/base/cmds/openspec.yaml` - OpenSpec command wrapper
- `.manifests/base/steps/detect-agent.yaml` - Agent detection step
- `.manifests/base/extensions/self/assets/` - Shared prompts and subagent definitions
- `.manifests/base/extensions/self/converters/` - Per-agent format converters (commands, mcp, subagents, config)

The naming conventions embed AI-specific terminology (`kfg.agent.*`, `kfg.detect-agent`) and use an ambiguous `self` namespace (`kfg.convert.self.*`, `kfg.extension.self.*`).

Other extensions (ctx7, openspec, chrome-devtools, playwright) follow a clearer pattern: `kfg.extension.<name>.<type>`.

## Goals / Non-Goals

**Goals:**
- Consolidate all AI-agent resources into one logical extension directory.
- Organize per-agent resources into subdirectories with assets/converters separation.
- Establish a short, consistent naming convention for all resource metadata.
- Make the base manifests generic (only truly generic steps remain).
- Update overlay workflow to use new names.
- Add Bats tests validating the new structure.

**Non-Goals:**
- Change Go code or CLI behavior (manifest-only change).
- Add new agents or functionality.
- Rename core generic steps (`kfg.cleanup`, `kfg.materialize`, `kfg.copy-context`, `kfg.ensure-gitignore`, `kfg.materialize-scaffold`).
- Change `apiVersion` from `kfg.dev/v1alpha1`.

## Decisions

### Consolidate into `extensions/ai-agents/` extension

All AI-agent resources move into `.manifests/base/extensions/ai-agents/`. This follows the existing extension pattern (ctx7, openspec, etc.) and makes it clear that AI agents are an extension, not core infrastructure.

The extension is included by default in `base/extensions/kustomization.yaml`, so existing workflows continue to work without changes to overlay kustomization files.

### Per-agent subdirectory structure

Each agent gets its own directory under `agents/` with:
- `assets/settings.yaml` - Agent configuration (kind: Assets)
- `converters/` - Per-agent format converters (kind: Converter)

```
agents/claude/
├── assets/settings.yaml
└── converters/
    ├── command.yaml
    ├── mcp.yaml
    └── subagent.yaml
```

Each agent directory has ONE `kustomization.yaml` at its root referencing `assets` and `converters`.

### Naming convention: `ai.<scope>.<kind>.<name>`

Per-agent resources:
- `ai.<agent>.asset.settings` - Agent settings
- `ai.<agent>.cmd.main` - Agent command wrapper
- `ai.<agent>.conv.<type>` - Agent converter (command, mcp, subagent, cfg)

Shared resources:
- `ai.cmds.<name>` - Command wrappers (openspec)
- `ai.steps.<name>` - Workflow steps (detect)
- `ai.conv.<name>` - Shared converters (to-json)
- `ai.prompts.<name>` - Shared prompts (git-commit, refactor-pure, etc.)
- `ai.subagents.<name>` - Shared subagent definitions (review-minimal)

Non-AI extensions:
- `ctx7.<kind>.<name>` - ctx7 extension resources
- `openspec.<kind>.<name>` - openspec extension resources
- `chrome.<kind>.<name>` - chrome-devtools extension resources
- `playwright.<kind>.<name>` - playwright extension resources

### Move `inject-ctx7-context` to ctx7 extension

The step `kfg.inject-ctx7-context` is ctx7-specific and moves from `base/steps/` to `extensions/ctx7/steps/`. Renamed to `ctx7.steps.inject`.

### Remove `extensions/self/` completely

After moving prompts, subagents, and converters to `ai-agents/`, the `self` extension has no remaining content and is removed.

### Remove `base/agents/` and `base/cmds/` completely

After moving all content to `ai-agents/`, these directories are removed.

### Update documentation with generic examples

Replace AI-specific examples in `docs/manifest-model.md` with generic examples (myapp.deploy, build, test).

## Data Contract

### New naming convention

```yaml
# Per-agent
metadata:
  name: ai.claude.asset.settings     # kind: Assets
  name: ai.claude.cmd.main           # kind: Cmd
  name: ai.claude.conv.command        # kind: Converter
  name: ai.claude.conv.mcp            # kind: Converter
  name: ai.claude.conv.subagent       # kind: Converter
  name: ai.opencode.conv.cfg          # kind: Converter

# Shared
metadata:
  name: ai.cmds.openspec              # kind: Cmd
  name: ai.steps.detect               # kind: Step
  name: ai.conv.to-json               # kind: Converter
  name: ai.prompts.git-commit         # kind: Assets
  name: ai.subagents.review-minimal   # kind: Assets

# Non-AI extensions
metadata:
  name: ctx7.steps.install            # kind: Step
  name: ctx7.steps.inject             # kind: Step
  name: ctx7.assets.mcp               # kind: Assets
  name: openspec.steps.install        # kind: Step
  name: chrome.assets.mcp             # kind: Assets
  name: playwright.assets.mcp         # kind: Assets
```

### Workflow reference updates

```yaml
# cmds
cmds:
  - ai.claude.cmd.main
  - ai.gemini.cmd.main
  - ai.opencode.cmd.main
  - ai.pi.cmd.main
  - ai.cmds.openspec

# steps
step: ai.steps.detect
step: ctx7.steps.install
step: ctx7.steps.inject
step: openspec.steps.install

# materialize references
ASSETS: "ai.claude.asset.settings"
CONVERTER: "ai.conv.to-json"
ASSETS: "ai.prompts.git-commit"
CONVERTER: "ai.claude.conv.command"
ASSETS: "ctx7.assets.mcp:chrome.assets.mcp:playwright.assets.mcp"
CONVERTER: "ai.claude.conv.mcp"
ASSETS: "ai.subagents.review-minimal"
CONVERTER: "ai.claude.conv.subagent"
```

## Risks / Trade-offs

- [Breaking existing overlays] -> Only the dev overlay exists; it is updated in this change.
- [Long naming convention] -> Chose short prefixes (`ai`, `ctx7`, `chrome`) to minimize length.
- [Removing `self` namespace] -> All content is redistributed; nothing is lost.
- [Documentation churn] -> Generic examples improve long-term clarity.

## Migration Plan

1. Create `extensions/ai-agents/` directory structure.
2. Move and rename all agent resources.
3. Move shared resources (prompts, subagents, cmds, steps, converters).
4. Move `inject-ctx7-context` to ctx7 extension.
5. Update all kustomization.yaml references.
6. Update overlay workflow with new names.
7. Remove old directories.
8. Rename non-AI extension resources.
9. Update documentation.
10. Add Bats tests.
11. Validate with `kustomize build` and `kfg apply`.

# Development Overlay (Deprecated)

**This overlay is deprecated.** Use `kfg ai` with the AI wizard overlay instead.

## Composition

This overlay composes:
- Framework base package
- Domain manifests (agents, commands, steps, prompts, subagents, converters, integrations)
- Development workflow (`agents-workflow.yaml`)

## Entrypoint

| Entrypoint | Description |
|------------|-------------|
| `kustomization.yaml` | Complete development overlay |

## Workflow

The `kfg.workflow.agents` CmdWorkflow provides the full agent development setup:

1. **Deprecation warning** (weight -100) — Logs deprecation notice
2. **Ensure gitignore** (weight -90) — Adds agent directories to `.gitignore`
3. **Detect agent** (weight -70) — Identifies active agent
4. **Scaffold** (weight -65) — Creates agent directories
5. **Materialize settings** (weight -63) — Generates agent config files
6. **Copy context** (weight -60) — Copies AGENTS.md to project root
7. **Install ctx7** (weight -55) — Installs Context7 skills
8. **Install openspec** (weight -53) — Installs OpenSpec skills
9. **Materialize commands** (weight -45) — Generates git-commit command
10. **Aggregate MCP** (weight -40) — Combines ctx7, chrome-devtools, playwright MCP configs
11. **Materialize subagents** (weight -35) — Generates review-minimal subagent
12. **Cleanup** — Removes temporary files

## Generated Artifacts

| Artifact | Agent | Path |
|----------|-------|------|
| Gitignore entries | All | `.gitignore` (appended) |
| Agent directories | OpenCode | `.opencode/` and subdirectories |
| Agent directories | Pi | `.pi/` and subdirectories |
| Settings | OpenCode | `opencode.json` |
| Settings | Pi | `.pi/settings.json` |
| Context | All | `AGENTS.md` |
| ctx7 skills | All | Agent-specific skills directory |
| OpenSpec skills | All | Agent-specific skills directory |
| Commands | OpenCode | `.opencode/commands/git-commit.md` |
| Commands | Pi | `.pi/prompts/git-commit.md` |
| MCP config | OpenCode | `opencode.json` (mcp section) |
| Subagents | OpenCode | `.opencode/agents/review-minimal.md` |

## Cleanup

The `kfg.cleanup` step runs in the `after` phase and removes all temporary artifacts.

## Migration to AI Wizard

To migrate from the dev overlay to the AI wizard:

1. Remove generated artifacts (`.opencode/`, `.pi/`, `opencode.json`, `AGENTS.md`)
2. Run `kfg ai` instead of `kfg run -k overlays/dev`
3. The wizard will generate a new, project-specific configuration

The dev overlay remains functional for backward compatibility but will display a
deprecation warning on each invocation.

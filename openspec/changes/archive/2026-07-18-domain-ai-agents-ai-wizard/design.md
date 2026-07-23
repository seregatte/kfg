## Context

The `ai-agents` domain currently provides a development overlay (`overlays/dev/`) that projects inherit via kustomize `resources:`. This inheritance pattern requires users to understand overlay composition and kustomize patching — a significant barrier. The existing pattern also means projects load more manifests than they need, and customization requires kustomize patches with JSON merge semantics.

Projects using kfg (allungare-infra, anafood, bling-api, homelab, nixai, nixai_v2, lifeos) all follow the same basic flow: `flake.nix` → `kfg apply` → agent environment. But each project's implementation varies in complexity from 1 file (`anafood`) to 59 files (`nixai_v2`). A wizard that generates tailored manifests eliminates this variance by producing exactly what each project needs.

The kfg engine already provides the building blocks: `kfg run` dispatches agents with workflows, `kfg.materialize` generates files from assets and converters, and `ai.steps.detect` detects which agent is active. What's missing is a discovery mechanism — a way for users to generate their initial kfg configuration without understanding the internal architecture.

## Goals / Non-Goals

**Goals:**
- Provide an interactive CLI command `kfg ai` that starts an AI agent pre-configured to generate kfg configurations
- Create a lightweight overlay (`overlays/ai/`) that serves as the wizard agent's workspace
- Create a wizard skill prompt that teaches agents how to compose manifests from domain building blocks
- Deprecate `overlays/dev/` with a visible warning, without removing it
- Support both pi and opencode agents initially
- Use existing domain building blocks (steps, cmds, converters, assets) before creating custom ones

**Non-Goals:**
- Not removing or modifying `overlays/dev/` functionality (backward compatible)
- Not adding new agent support (claude, gemini) — only pi and opencode initially
- Not creating custom converters — reuse existing `ai.{agent}.conv.command`
- Not adding MCP servers, subagents, or settings to the wizard overlay (those are project concerns, not wizard concerns)
- Not modifying the kfg engine's execution model or shell generation

## Decisions

1. **CLI command wraps `kfg run` via subprocess (not in-process)**: The `kfg ai` command resolves `KFG_AI_AGENT` (default `pi`), constructs `kfg run -k overlays/ai <agent> -- <args>`, and executes it. Subprocess is simpler than refactoring `run.go`'s execution pipeline and avoids coupling the wizard to the engine's internal dispatch.

2. **Wizard skill as a prompt asset (not system prompt)**: The wizard lives as a command in `.opencode/commands/wizard.md` or `.pi/prompts/wizard.md`. This keeps the wizard available on-demand without polluting the agent's system prompt. The agent can invoke `/wizard` when the user expresses intent to generate configurations.

3. **ctx7 integration for library documentation**: The wizard overlay installs ctx7 so the agent can look up library documentation while generating configurations. This follows the same pattern established by `overlays/dev/`.

4. **Deprecation via a separate step (not inline)**: A new step `ai.steps.deprecation-warn` is defined as a standalone resource and added to the dev overlay workflow. This makes the deprecation composable and testable, rather than inline bash in the workflow YAML.

5. **`KFG_AI_AGENT` env var (not KFG_WIZARD_AGENT)**: Named to match the `kfg ai` command name. Defaults to `pi`. Can be set to `opencode` to use the opencode agent.

## Risks / Trade-offs

- **[Agent availability]** The wizard depends on the agent binary being installed (pi or opencode). If neither is installed, the wizard simply fails with the agent's own error message → Mitigation: document requirements in README/help text
- **[Skill prompt staleness]** If new building blocks are added to the domain but the wizard skill is not updated, the wizard may create custom manifests for things that already exist → Mitigation: skill prompt instructs agent to explore the domain manifests directory before creating anything
- **[Subprocess overhead]** Each `kfg ai` invocation starts a subprocess for `kfg run`, which means two binaries are loaded → Acceptable: the overhead is negligible for interactive use

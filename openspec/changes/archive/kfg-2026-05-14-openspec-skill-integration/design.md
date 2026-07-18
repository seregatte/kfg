## Context

The kfg agents workflow (`kfg.workflow.agents`) installs skills from multiple extensions (ctx7, chrome-devtools, playwright) per-agent. The openspec extension has an install step (`kfg.extension.openspec.install`) defined in `.manifests/base/extensions/openspec/steps/install.yaml` but it is not wired into the agents workflow. The openspec Cmd is also misplaced in `.manifests/overlay/dev/cmds.yaml` instead of `.manifests/base/cmds/`.

Current state:
- `.manifests/base/extensions/openspec/` has only `steps/install.yaml`
- `.manifests/base/cmds/` has `agents.yaml` with claude, opencode, gemini, pi Cmds
- `.manifests/overlay/dev/cmds.yaml` has a placeholder openspec Cmd with extra env var logic
- The openspec install step only copies skills from `$AGENT_HOME/skills/`, not commands from `$AGENT_HOME/commands/`
- `openspec init --tools=<agent> --force` generates 4 skills and 4 commands per agent

## Goals / Non-Goals

**Goals:**
- Move the openspec Cmd to base layer following the standard agent cmd pattern
- Update the openspec install step to copy both skills and commands
- Integrate openspec install into the agents workflow with per-agent steps
- Delete the misplaced overlay dev cmds.yaml

**Non-Goals:**
- Adding openspec MCP assets or converters (openspec generates skills directly, no MCP server)
- Changing the openspec install step's default env vars (keep agent-specific overrides in workflow)
- Modifying other extension install steps

## Decisions

### D1: Cmd naming follows `kfg.agent.cmd.<name>` pattern

The openspec Cmd uses `kfg.agent.cmd.openspec` to match `kfg.agent.cmd.claude`, `kfg.agent.cmd.opencode`, etc. The `commandName` is `openspec` and the run script uses `command openspec "$@"` for direct passthrough.

### D2: Weight -53 for openspec install steps

Openspec install runs after ctx7 install (-55) but before ctx7 inject (-50). This ensures ctx7 skills are installed first (they may be referenced by openspec), and openspec skills are available before context injection.

### D3: Copy commands to `$OUTPUT_DIR/../commands/`

The openspec init generates commands in `$AGENT_HOME/commands/`. The install step copies them to the agent's commands directory relative to the skills output. For opencode: `.opencode/skills/` and `.opencode/commands/`.

### D4: No MCP asset for openspec

Unlike ctx7 which provides an MCP server, openspec generates skill files directly. No Assets or Converter resources are needed.

## Risks / Trade-offs

- **[Risk] openspec CLI not available** → The install step uses `failurePolicy: Ignore` so the workflow continues if openspec is missing
- **[Trade-off] 4 invocations of openspec.install** → Verbose but explicit. Each agent gets its own `when` condition and env vars, matching the pattern used by ctx7 install

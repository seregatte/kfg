## Why

The openspec extension has an install step (`kfg.extension.openspec.install`) but it is not wired into the agents workflow. Developers must manually run `openspec init` to get openspec skills installed per-agent. The openspec Cmd is also misplaced in the overlay dev directory instead of the base layer where it belongs as a generic CLI wrapper.

## What Changes

- **New base Cmd `kfg.agent.cmd.openspec`**: Move the openspec CLI wrapper from `.manifests/overlay/dev/cmds.yaml` to `.manifests/base/cmds/openspec.yaml`, following the same pattern as `kfg.agent.cmd.claude`, `kfg.agent.cmd.opencode`, etc.
- **Delete overlay dev cmds.yaml**: Remove `.manifests/overlay/dev/cmds.yaml` and update the overlay kustomization.
- **Modify openspec install step**: Update `kfg.extension.openspec.install` to copy both skills AND commands from `$AGENT_HOME/` to `$OUTPUT_DIR/`. Currently only copies skills.
- **Modified workflow `kfg.workflow.agents`**: Add 4 openspec install steps (one per agent) at weight -53, between ctx7 install (-55) and ctx7 inject (-50). Each step passes agent-specific `TOOLS_FLAG`, `AGENT_HOME`, and `OUTPUT_DIR`.

## Capabilities

### New Capabilities

- `openspec-cmd`: Base Cmd resource that wraps the `openspec` CLI binary as a shell function, following the standard agent cmd pattern.

### Modified Capabilities

- `skill-installation-steps`: The openspec install step now copies both skills and commands from `$AGENT_HOME/skills/` and `$AGENT_HOME/commands/` to the output directory.

## Impact

- **Manifest files**: 1 new file in `.manifests/base/cmds/`, 1 deleted file in `.manifests/overlay/dev/`, 1 modified step in `.manifests/base/extensions/openspec/steps/`, 1 modified workflow in `.manifests/overlay/dev/`
- **Shell generation**: The `openspec` function will be available in all workflows that include `kfg.agent.cmd.openspec`
- **Agent UX**: Developers get openspec skills and commands automatically installed when running the agents workflow, without manual `openspec init`
- **Dependencies**: Requires `openspec` CLI available in PATH for the install step

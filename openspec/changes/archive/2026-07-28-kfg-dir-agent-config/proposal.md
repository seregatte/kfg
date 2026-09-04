## Why

AI coding agents (pi, OpenCode) store configuration, sessions, auth, and other
runtime state in agent-specific directories. Running them through KFG should use
isolated, disposable agent state by default, while allowing users to opt into a
persistent directory.

## What Changes

- Add `KFG_DIR` setup logic to `ai.pi.cmd.main` and `ai.opencode.cmd.main` Cmd manifests
- When `KFG_DIR` is unset, the Cmd creates a temporary directory with `mktemp -d`
- When `KFG_DIR` is set by the user, use it as-is (persistent, no cleanup)
- Export the agent-specific config-dir variable from `KFG_DIR`:
  - `PI_CODING_AGENT_DIR` for pi
  - `OPENCODE_CONFIG_DIR` for opencode
- Clean up the temp directory on EXIT only when KFG created it

## Capabilities

### Modified Capabilities

- `pi-cmd`: The pi Cmd now manages `PI_CODING_AGENT_DIR` from `KFG_DIR`
- `opencode-cmd`: The opencode Cmd now manages `OPENCODE_CONFIG_DIR` from `KFG_DIR`

## Impact

- Affected manifests: `manifests/cmds/agents/pi.yaml`, `manifests/cmds/agents/opencode.yaml`
- Affected documentation: `manifests/cmds/README.md`
- Default behavior changes: `kfg run pi` now uses a fresh temp directory cleaned up on exit
- Agent runtime state is ephemeral by default; opt into persistence with `KFG_DIR=...`

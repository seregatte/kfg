## Context

Each Cmd manifest has a `spec.run` field containing a bash script that wraps the
agent binary. The generated shell code calls `formatCmdEnv` for cmd-level env vars,
which emits `export KEY="VALUE"` as literal assignment — without `${VAR:-default}`
expansion. Therefore, bash expansion for `KFG_DIR` must live inside `spec.run`
instead of `spec.env`.

## Goals / Non-Goals

**Goals:**

- Provide isolated agent state by default (temporary directory)
- Support persistent state via user-provided `KFG_DIR`
- Implement per-Cmd, not in global Go generation or bash helpers
- Clean up KFG-created directories only

**Non-Goals:**

- Adding the same pattern to Claude Code or Codex (future)
- Modifying the Go generator, bash templates, or framework steps
- Adding new manifest resources or kinds

## Decisions

### Inline bash in spec.run

The KFG_DIR setup is a bash preamble prepended to each Cmd's `spec.run`. This keeps
agent-specific behavior in declarative manifests instead of hardcoding agent knowledge
into `generate.go` or `bash_helper.tmpl`.

```bash
__kfg_dir_created=0
if [ -z "${KFG_DIR:-}" ]; then
  KFG_DIR="$(mktemp -d)"
  __kfg_dir_created=1
fi
export KFG_DIR
export AGENT_CONFIG_DIR_VAR="$KFG_DIR"
if [ "$__kfg_dir_created" = "1" ]; then
  trap 'rm -rf "$KFG_DIR"' EXIT
fi
command agent "$@"
```

### Per-agent mapping

Each agent maps `KFG_DIR` to its specific variable name:
- pi → `PI_CODING_AGENT_DIR`
- OpenCode → `OPENCODE_CONFIG_DIR`

Future agents follow the same pattern:
- Claude Code → `CLAUDE_CONFIG_DIR`
- Codex → `CODEX_HOME`

### Cleanup responsibility

The trap only fires when `__kfg_dir_created=1`. If the user provided `KFG_DIR`,
the directory is never deleted by KFG.

## Risks / Trade-offs

- The `$HOME/.config/opencode` style defaults from the original issue proposal
  are not implemented because `spec.env` does not support bash expansion syntax.
  When `KFG_DIR` is unset, a temp dir is always used — there is no fallback to
  a persistent default directory.

## Migration Plan

No migration required. Agents without `KFG_DIR` set previously used their own
internal default paths (e.g. `~/.pi/agent`, `$XDG_CONFIG_HOME/opencode`). With
this change, they use a temp directory instead when `KFG_DIR` is unset.

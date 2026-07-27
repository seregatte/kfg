## Context

The current devShell PATH ordering (showing only relevant entries) looks like:

```
/nix/store/...-coreutils/bin
/nix/store/...-go/bin
/nix/store/...-node/bin
/nix/store/...-kfg-bundle/bin
/etc/profiles/per-user/...
/run/current-system/sw/bin
/opt/homebrew/bin
/usr/local/bin
/usr/bin
```

All Nix store entries are listed before any host path, so every command resolves to the Nix version. This design adds a mechanism to selectively prefer host-installed commands without reordering the entire PATH.

## Decisions

### 1. Temporary symlink directory, not PATH reorder

A symlink directory at a well-known temp path is created, containing only the commands found in the host path. This directory is prepended to PATH, keeping the Nix entries as immediate fallback for those specific commands while leaving all other tools unchanged.

**Alternative considered:** Prepend the full host PATH. Rejected because it would override coreutils (sed, grep, awk) with BSD variants on macOS and break CI reproducibility.

### 2. Command-level allowlist, not package-level

The resolution is by command name (binary), not by Nix package. This allows selecting `go` while keeping `bats` from Nix, even though both might come from the same or different packages.

**Alternative considered:** Build-level exclusion from buildInputs. Rejected because it conflicts with Nix evaluation — packages must be realized before shellHook runs.

### 3. Existence-based, not version-based

If the command exists on the host PATH, it wins. No version comparison.

**Alternative considered:** Version comparison. Rejected because there is no reliable cross-command way to compare versions, and the user explicitly chose "use anyway."

### 4. External PATH defined as "not /nix/store"

A path entry is considered "external" if it does not contain `/nix/store`. This covers Homebrew, system profiles, user profiles, Nix profiles, `/usr/local/bin`, `/usr/bin`, etc.

### 5. ./bin priority preserved in dev shell

The dev shell exports `PATH="./bin:$PATH"` before the system preference logic, so `./bin/kfg` always wins.

## Implementation

### Shared shellHook function

A single shell hook snippet shared across all three devShells:

```bash
# If KFG_PREFER_SYSTEM is explicitly set to 0, skip this logic.
if [ "${KFG_PREFER_SYSTEM:-1}" != "0" ]; then
  # Default command list (space-separated)
  DEFAULT_COMMANDS="go node npm npx corepack uv uvx bats openspec pi ctx7 chrome-devtools-mcp gws notebooklm nblm opencode playwright"

  # Use env var override if set, otherwise default
  COMMANDS="${KFG_PREFER_SYSTEM_COMMANDS:-$DEFAULT_COMMANDS}"

  # Build a PATH that excludes /nix/store entries
  HOST_PATH=""
  IFS=':' read -ra PATH_ENTRIES <<< "$PATH"
  for entry in "${PATH_ENTRIES[@]}"; do
    case "$entry" in
      */nix/store/*) ;;
      *) HOST_PATH="${HOST_PATH:+$HOST_PATH:}$entry" ;;
    esac
  done

  # Create temp directory for system-preferred binaries
  KFG_SYSTEM_BIN_DIR="$(mktemp -d /tmp/kfg-system-bin-XXXXXX 2>/dev/null)"
  if [ -n "$KFG_SYSTEM_BIN_DIR" ]; then
    for cmd in $COMMANDS; do
      host_bin="$(PATH="$HOST_PATH" command -v "$cmd" 2>/dev/null)"
      if [ -n "$host_bin" ] && [ -x "$host_bin" ]; then
        ln -sf "$host_bin" "$KFG_SYSTEM_BIN_DIR/$cmd"
      fi
    done
    export PATH="$KFG_SYSTEM_BIN_DIR:$PATH"
  fi
fi
```

### Placement in flake.nix

A `sharedShellHook` variable is defined in the `devShells` let block and referenced by all three shells. Each shell appends its own specific logic after the shared block.

### Shells

- `default`: shared hook only
- `dev`: shared hook + `PATH="./bin:$PATH"` (before shared hook, so ./bin wins) + other dev-specific vars + `source <(go run ...)`
- `ci`: shared hook + CI-specific logic (bats vendor setup)

## Risks / Trade-offs

- [Temp directory not cleaned up] → Mitigation: Use `trap` or rely on `/tmp` cleanup. Acceptable for interactive shells.
- [Host command might be incompatible version] → Mitigation: User controls the allowlist and can always set `KFG_PREFER_SYSTEM=0`.
- [Race condition on temp dir name] → Mitigation: `mktemp` ensures uniqueness.
- [Increase shellHook complexity] → Mitigation: Encapsulated in a single block; minimal conditional overhead when disabled.

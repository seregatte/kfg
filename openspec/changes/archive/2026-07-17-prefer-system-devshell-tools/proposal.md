## Why

When a developer runs `nix develop`, the devShell always prefers executables from `/nix/store` over those already installed on the host system (Homebrew, system profiles, `/usr/local/bin`, etc.). This forces the Nix-managed version even when the developer already has the same tool installed and wants to use their system-configured version — for example, Go installed via Homebrew with custom GOROOT, or node installed via nvm with project-specific npm packages.

The user's system tools may have configurations, plugins, or versions that differ from what the flake pins. While the Nix devShell guarantees reproducibility, there is no easy way to opt into using the host's version of specific tools when entering the shell.

## What Changes

- **Shared shellHook logic** in `flake.nix` that filters the current `PATH` to find host-installed versions of a configurable list of commands
- **`KFG_PREFER_SYSTEM_COMMANDS`** environment variable: space-separated list of commands to prefer from the host path
- **`KFG_PREFER_SYSTEM=0`** environment variable: disable system preference entirely, restoring standard Nix behavior
- **Default allowlist**: `go node npm npx corepack uv uvx bats openspec pi ctx7 chrome-devtools-mcp gws notebooklm nblm opencode playwright`
- **Temporary symlink directory**: a `/tmp/kfg-system-bin-*` directory containing only the resolved host commands, prepended to `PATH` — never the full host `PATH`
- **Applies to all three shells**: `devShells.default`, `devShells.dev`, and `devShells.ci`
- **Dev shell preserves `./bin` priority**: `./bin` remains first in `PATH` for the `dev` shell

## Capabilities

### New Capabilities

- `devshell-prefer-system-tools`: Mechanism to prefer host-installed executables over Nix devShell versions for a configurable set of commands

## Impact

- `flake.nix` — shared `shellHook` logic for PATH filtering across all three devShells
- `docs/context/openspec/specs/kfg-devshell-prefer-system-tools/spec.md` — new spec for the behavior
- `README.md` — document the environment variables and behavior
- No version bump (feature branch against release/v0.0.12)
- No change to `buildInputs` — packages are still realized by Nix, only the resolved command changes

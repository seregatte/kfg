## Context

kfg provides three devShells in `flake.nix`: `default`, `dev`, and `ci`. Currently, `default` and `dev` are identical — both provide the full development environment (Go, Node.js, bats, yq, jq, gomplate, kfg-bundle) and both run a shellHook that sources `go run ./src/cmd/kfg apply`. When external projects consume kfg via flake inputs and use `inputsFrom = [kfg.devShells.${system}.default]`, they inherit the Go development shellHook, which fails in non-Go projects.

Current shellHook (lines 218-233 in flake.nix):
```bash
export KFG_DIR=${self.outPath}
export PATH="./bin:$PATH"
export OPENSPEC_ROOT_DIR=docs/context
if [ "$COLUMNS" -lt 45 ]; then
  export STARSHIP_CONFIG=${self.outPath}/assets/starship/mobile.toml
else
  export STARSHIP_CONFIG=${self.outPath}/assets/starship/full.toml
fi
source <(go run ./src/cmd/kfg apply -k packages/domains/ai-agents/overlays/dev)
```

## Goals / Non-Goals

**Goals:**

- `devShells.default` is consumer-friendly: provides tools and environment setup without Go development shellHook
- `devShells.dev` is development-only: maintains Go workflow with `source <(go run ./src/cmd/kfg apply)`
- Consumers using `inputsFrom = [kfg.devShells.${system}.default]` see no Go errors
- kfg developers use explicit `nix develop .#dev` for development workflow
- Version bump 0.0.6 → 0.0.7 reflects breaking change
- Documentation clearly distinguishes consumer vs development usage

**Non-Goals:**

- Changing buildInputs tool list (devInputs unchanged)
- Changing kfg-bundle composition
- Changing `ci` shell behavior
- Modifying consumer projects (dotfiles, homelab, toolbox-stack, lifeos)

## Decisions

### 1. Swap roles: default → consumer, dev → development

Rename semantics: `default` becomes the "consumer shell" (tools + env, no Go), `dev` becomes the "development shell" (Go + shellHook).

**Alternative considered:** Add new `devShells.consumer` and keep `default` as dev.
**Rejected** because:
- Semantically confusing: "default" should be the safe, common case
- Most users of kfg.devShells are consumers, not developers
- Breaking change is acceptable (kfg repo is solo use, documented clearly)

### 2. Remove Go and Node.js from default buildInputs

Consumers don't need `pkgs.go` or `pkgs.nodejs` in buildInputs:
- kfg-bundle already contains Node.js-based tools (openspec, ctx7, etc.)
- Consumers don't compile Go code
- Reduces unnecessary dependencies

**Alternative considered:** Keep Go/Node in default for completeness.
**Rejected** because:
- Consumers never use them
- Larger closure size
- Confusing (why inherit Go if not developing?)

### 3. Remove OPENSPEC_ROOT_DIR from default shellHook

Consumers set their own `OPENSPEC_ROOT_DIR` in their shellHook.

**Alternative considered:** Keep `OPENSPEC_ROOT_DIR=docs/context` in default.
**Rejected** per user feedback: consumers set their own.

### 4. Keep STARSHIP_CONFIG in both shells

Dynamic starship config (mobile vs full) is useful for both developers and consumers. Consumers can override if desired.

**Approved by user:** Keep in both default and dev.

### 5. Version bump 0.0.6 → 0.0.7

Semver: this is a breaking change for kfg development workflow. Even though consumers benefit, the version bump signals that kfg's devShell contract has changed.

**Alternative considered:** Keep 0.0.6, just document.
**Rejected** per user request: bump version.

## Risks / Trade-offs

- [Breaking change for kfg developers] → Mitigation: Clear documentation in README.md and AGENTS.md; developers use explicit `.#dev` shell
- [Consumers unaware of new behavior] → Mitigation: No action needed — consumers benefit automatically (no Go error)
- [Future kfg contributors confused by explicit .#dev] → Mitigation: AGENTS.md clearly documents development workflow
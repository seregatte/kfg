## Why

kfg was designed to replace nixai's imperative agent-environment system, and the
behavioral migration (manifests, extensions, steps) is complete. However, kfg's
dev shell still has an external hard dependency on `github:seregatte/nixai` via
`inputsFrom` — meaning developers cannot build or develop kfg without nixai being
reachable as a Nix flake input. Absorbing nixai's devShells and packages into kfg
makes the project fully self-contained.

## What Changes

- Remove `nixai` from `flake.nix` inputs; declare all dev shell packages inline
- Migrate `devInputs` toolchain (`yq-go`, `jq`, `yajsv`, `gomplate`, `coreutils`,
  `findutils`, `gnused`, `gnugrep`, `bash`, `bats`, `google-cloud-sdk`, `uv`) into kfg's flake
- Migrate `npmGlobal` derivation (`@fission-ai/openspec`, `@mariozechner/pi-coding-agent`,
  `ctx7`, `chrome-devtools-mcp`) into kfg's flake
- Migrate `gws-bin` derivation (Google Workspace CLI v0.22.5, all four platforms) into kfg's flake
- Migrate `notebooklmWrapper` and `nblmWrapper` shell script wrappers into kfg's flake
- Create `kfg-bundle` package (`symlinkJoin`) replacing nixai's bundle, without
  the circular `kfg-bin` reference (kfg binary is already `packages.${system}.default`)
- Rename `NIXAI_DIR` env var to `KFG_DIR` in `shellHook`
- Copy Starship prompt assets (`assets/starship/full.toml`, `assets/starship/mobile.toml`)
  from nixai into kfg
- Update all `shellHook` references to use `${self.outPath}/assets/starship/`
- Regenerate `flake.lock` (removes `nixai` and `nixai_2` entries)
- Purge residual `.nixai/overlay/dev` path references from Go source help text
  (`apply.go`, `build.go`, `loader.go`)
- Fix stale `~/.config/nixai/store` references in `config_test.go` (actual default
  is already `~/.kfg/store`)
- Update `README.md`: remove non-existent `kfg image` and `kfg workspace` commands;
  add `kfg run` and `kfg sys` to command reference; add `nix develop` dev shell docs;
  fix repository structure paths
- Update `docs/AGENTS.md`: trim to minimal development instructions
- Update `docs/cli-reference.md`: remove non-existent `kfg assets` and `kfg sys gc`
  commands; add correct `kfg sys cache` subcommands; fix flag and env var accuracy

## Capabilities

### New Capabilities

- `kfg-self-contained-dev-shell`: kfg's Nix dev shell declares all required
  tooling inline without inheriting from any external flake input

### Modified Capabilities

- `kfg-build-nix-packaging`: The flake no longer declares a `nixai` input; the
  `kfg-bundle` package and dev shell toolchain are declared entirely within kfg's
  own `flake.nix`

## Impact

- `flake.nix` — primary change; rewritten to inline all nixai content
- `flake.lock` — regenerated (nixai entries removed)
- `assets/starship/full.toml`, `assets/starship/mobile.toml` — new files
- `src/cmd/kfg/apply.go` — help text path examples
- `src/cmd/kfg/build.go` — help text path examples
- `src/internal/kustomize/loader.go` — code comment
- `src/cmd/kfg/config_test.go` — stale store path references
- `README.md` — command reference, dev shell docs, repo structure
- `docs/AGENTS.md` — trimmed to minimal dev instructions
- `docs/cli-reference.md` — command and flag accuracy

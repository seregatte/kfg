## 1. Starship Assets

- [ ] 1.1 Create `assets/starship/` directory in kfg repo
- [ ] 1.2 Copy `full.toml` from `nixai/assets/starship/full.toml` to `assets/starship/full.toml`
- [ ] 1.3 Copy `mobile.toml` from `nixai/assets/starship/mobile.toml` to `assets/starship/mobile.toml`

## 2. Rewrite flake.nix

- [ ] 2.1 Remove `nixai.url` from `inputs`; keep only `nixpkgs`
- [ ] 2.2 Update `outputs` signature to remove `nixai` binding
- [ ] 2.3 Declare `devInputs` list inline (`yq-go`, `jq`, `yajsv`, `gomplate`, `coreutils`, `findutils`, `gnused`, `gnugrep`, `bash`, `bats`, `google-cloud-sdk`, `uv`)
- [ ] 2.4 Declare `npmGlobal` derivation inline (rename `pname` to `"kfg-global-npm"`; packages: `@fission-ai/openspec`, `@mariozechner/pi-coding-agent`, `ctx7`, `chrome-devtools-mcp`)
- [ ] 2.5 Declare `gws-bin` derivation inline (version 0.22.5, all four platform hashes)
- [ ] 2.6 Declare `notebooklmWrapper` and `nblmWrapper` shell script wrappers inline
- [ ] 2.7 Declare `kfg-bundle` as `pkgs.symlinkJoin` of `npmGlobal`, `gws-bin`, `notebooklmWrapper`, `nblmWrapper`, `pkgs.claude-code`, `pkgs.gemini-cli-bin`, `pkgs.opencode`, `pkgs.playwright-test` (no `kfg-bin` — avoid circular reference)
- [ ] 2.8 Expose `kfg-bundle` under `packages.${system}.kfg-bundle`
- [ ] 2.9 Update `devShells.default` and `devShells.dev`: `buildInputs = devInputs ++ [ pkgs.nodejs pkgs.go kfg-bundle ]`
- [ ] 2.10 Update `shellHook`: replace `NIXAI_DIR` with `KFG_DIR`; update `STARSHIP_CONFIG` paths to use `${self.outPath}/assets/starship/`; merge existing kfg `shellHook` lines (`PATH`, `OPENSPEC_ROOT_DIR`, `kfg apply` source)

## 3. Regenerate flake.lock

- [ ] 3.1 Remove `nixai` and `nixai_2` entries from `flake.lock` (or run `nix flake lock --update-input nixai`)
- [ ] 3.2 Run `nix flake lock` to regenerate the lock file cleanly
- [ ] 3.3 Verify `flake.lock` contains no `nixai` entries

## 4. Verify Dev Shell

- [ ] 4.1 Run `nix develop --command go version` and confirm Go is available
- [ ] 4.2 Run `nix develop --command bats --version` and confirm bats is available
- [ ] 4.3 Run `nix develop --command yq --version` and confirm yq is available
- [ ] 4.4 Run `nix develop --command make build` and confirm `./bin/kfg` builds successfully
- [ ] 4.5 Run `nix develop --command make test` (Go unit tests pass)
- [ ] 4.6 Run `nix develop --command make test-bats` (Bats integration tests pass)

## 5. Fix Go Source References

- [ ] 5.1 Update `src/cmd/kfg/apply.go` help text: replace all `.nixai/overlay/dev` occurrences with `packages/domains/ai-agents/overlays/dev`
- [ ] 5.2 Update `src/cmd/kfg/build.go` help text: replace `.nixai/overlay/dev` and `.nixai/base` occurrences with appropriate `packages/domains/...` paths
- [ ] 5.3 Update `src/internal/kustomize/loader.go` line 70: update comment example path
- [ ] 5.4 Update `src/cmd/kfg/config_test.go`: fix `expectedDefault` (line 94) from `~/.config/nixai/store` to `~/.kfg/store`; rename all `test-nixai-store` fixtures to `test-kfg-store`; update stale comments

## 6. Update README.md

- [ ] 6.1 Remove `kfg image` section (command does not exist)
- [ ] 6.2 Remove `kfg workspace` section (command does not exist)
- [ ] 6.3 Update Command Reference table: add `kfg run`, `kfg sys log`, `kfg sys cache`; remove `kfg image` and `kfg workspace` rows
- [ ] 6.4 Add Development section entry for `nix develop` as entrypoint to full dev shell (Go, Node.js, bats, AI agent tools)
- [ ] 6.5 Fix Repository Structure paths for OpenSpec (`docs/context/openspec/` unified root)
- [ ] 6.6 Fix Environment Variables table: `KFG_STORE_DIR` default `~/.kfg/store`; remove stale vars

## 7. Update docs/AGENTS.md

- [ ] 7.1 Trim to minimal dev instructions: keep Quick Reference commands, OpenSpec `kfg run` usage, OpenSpec root reference, language policy, and essential gotchas
- [ ] 7.2 Remove the long OpenSpec root structure breakdown and list of spec prefixes (this level of detail belongs in specs, not AGENTS.md)

## 8. Update docs/cli-reference.md

- [ ] 8.1 Remove `kfg assets` section entirely (command does not exist)
- [ ] 8.2 Replace `kfg sys gc` section with accurate `kfg sys cache` section covering all subcommands: `ls`, `inspect`, `rm`, `prune`, `du`, `exists`, `store`, `restore`
- [ ] 8.3 Fix Global Flags table: remove non-existent `--debug` flag; correct `-v/--verbose` type (`int`, not boolean); remove non-existent `--store` and `--session-id` global flags
- [ ] 8.4 Fix `KFG_STORE_DIR` default in Environment Variables table: `~/.config/kfg/store` → `~/.kfg/store`
- [ ] 8.5 Fix `KFG_VERBOSE` scale: update to `0`=quiet through `5`=debug (currently shows 0-3)

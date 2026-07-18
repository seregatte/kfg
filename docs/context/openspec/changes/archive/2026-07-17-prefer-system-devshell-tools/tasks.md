## 1. Create spec for system tools preference

- [x] 1.1 Create `docs/context/openspec/specs/kfg-devshell-prefer-system-tools/spec.md` with requirements

## 2. Implement shared shellHook in flake.nix

- [x] 2.1 Add `sharedShellHook` variable in the `devShells` let block with PATH filtering logic
- [x] 2.2 Reference shared hook in `devShells.default` shellHook
- [x] 2.3 Reference shared hook in `devShells.dev` shellHook (ensure `./bin` remains first)
- [x] 2.4 Reference shared hook in `devShells.ci` shellHook
- [x] 2.5 Verify `KFG_PREFER_SYSTEM=0` disables the logic
- [x] 2.6 Verify `KFG_PREFER_SYSTEM_COMMANDS` overrides the default allowlist

## 3. Test locally

- [x] 3.1 Run `nix develop --command bash -c "command -v go"` with Go on host → verify host Go is used
- [x] 3.2 Run `nix develop --command bash -c "command -v node"` with Node on host → verify host Node is used
- [x] 3.3 Run `nix develop --command bash -c "command -v bats"` without bats on host → verify Nix bats is used
- [x] 3.4 Run `nix develop .#dev --command bash -c "command -v kfg"` → verify `./bin/kfg` wins over host and Nix
- [x] 3.5 Run `KFG_PREFER_SYSTEM=0 nix develop --command bash -c "command -v go"` → verify Nix go is used
- [x] 3.6 Run `KFG_PREFER_SYSTEM_COMMANDS="opencode" nix develop --command bash -c "command -v node"` → verify Nix node is used (not in allowlist)
- [x] 3.7 Run `nix develop --command make build` → verify build succeeds
- [x] 3.8 Run `nix develop .#dev --command make test` → verify Go unit tests pass
- [x] 3.9 Run `nix develop .#dev --command make test-bats` → verify integration tests pass

## 4. Update documentation

- [x] 4.1 Update `README.md` with devShell system preference section
- [x] 4.2 Add examples for KFG_PREFER_SYSTEM_COMMANDS and KFG_PREFER_SYSTEM=0

## 5. Code quality and commit

- [x] 5.1 Run `nix develop .#dev --command make fmt lint vet`
- [x] 5.2 Commit: `git add -A && git commit -m "feat(devShells): prefer host-installed tools via configurable allowlist"`
- [x] 5.3 Push: `git push -u origin feature/prefer-system-devshell-tools`
- [x] 5.4 Create draft PR: `gh pr create --draft --base release/v0.0.12`

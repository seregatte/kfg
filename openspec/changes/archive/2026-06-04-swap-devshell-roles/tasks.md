## 1. Update flake.nix devShells

- [x] 1.1 Edit `devShells.default` buildInputs: remove `pkgs.go` and `pkgs.nodejs`; change to `devInputs ++ [kfg-bundle]`
- [x] 1.2 Edit `devShells.default` shellHook: remove `PATH="./bin:$PATH"`, remove `OPENSPEC_ROOT_DIR=docs/context`, remove `source <(go run ./src/cmd/kfg apply ...)`
- [x] 1.3 Verify `devShells.default` shellHook retains: `export KFG_DIR=${self.outPath}` and STARSHIP_CONFIG dynamic setup
- [x] 1.4 Edit `devShells.dev` shellHook: verify `PATH="./bin:$PATH"` present, verify `OPENSPEC_ROOT_DIR=docs/context` present
- [x] 1.5 Verify `devShells.dev` maintains: `pkgs.go`, `pkgs.nodejs` in buildInputs, `source <(go run ./src/cmd/kfg apply)` in shellHook
- [x] 1.6 Verify `devShells.ci` unchanged (no changes to minimal CI shell)

## 2. Update version

- [x] 2.1 Update `version = "0.0.7"` in flake.nix (line 10)
- [x] 2.2 Verify version reflects breaking change for kfg development workflow

## 3. Test locally

- [x] 3.1 Run `nix develop --command bash -c "which kfg"` → verify kfg binary available in default shell
- [x] 3.2 Run `nix develop --command bash -c "which go"` → verify Go NOT available in default shell
- [x] 3.3 Run `nix develop --command bash -c "echo \$STARSHIP_CONFIG"` → verify STARSHIP_CONFIG set in default shell
- [x] 3.4 Run `nix develop --command bash -c "echo \$OPENSPEC_ROOT_DIR"` → verify OPENSPEC_ROOT_DIR NOT set in default shell
- [x] 3.5 Run `nix develop .#dev --command bash -c "which go"` → verify Go available in dev shell
- [x] 3.6 Run `nix develop .#dev --command bash -c "go run ./src/cmd/kfg version"` → verify Go workflow works in dev shell
- [x] 3.7 Run `nix develop .#dev --command make build` → verify build succeeds in dev shell
- [x] 3.8 Run `nix develop .#dev --command make test` → verify Go unit tests pass in dev shell
- [x] 3.9 Run `nix develop .#dev --command make test-bats` → verify Bats integration tests pass in dev shell
- [x] 3.10 Run `nix develop .#ci --command bash -c "which go"` → verify Go available in CI shell

## 4. Update documentation

- [x] 4.1 Update `README.md`: add DevShells section explaining `default` (consumer) vs `dev` (development) vs `ci` (minimal)
- [x] 4.2 Update `README.md`: Development section now uses `nix develop .#dev` instead of implicit `nix develop`
- [x] 4.3 Update `docs/AGENTS.md`: Quick Dev Commands section clarifies `nix develop .#dev` for kfg development
- [x] 4.4 Update `docs/AGENTS.md`: Git Worktree Workflow section mentions explicit `.#dev` shell for development

## 5. Code quality and commit

- [x] 5.1 Run `make fmt lint vet` → verify code quality checks pass
- [x] 5.2 Commit changes: `git add flake.nix README.md docs/AGENTS.md && git commit -m "refactor(devShells): swap default/dev roles (v0.0.7)..."`
- [x] 5.3 Push to remote: `git push origin <branch>`
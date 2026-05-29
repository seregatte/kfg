## Context

kfg's `flake.nix` currently has two sections: a `packages` block that builds the
kfg binary from GitHub Releases, and a `devShells` block that delegates entirely
to nixai via `inputsFrom = [ nixai.devShells.${system}.default ]`. The nixai
devShell provides every tool needed for development (Go, Node.js, bats, yq, jq,
gcloud, uv, openspec CLI, gws, claude-code, gemini-cli, opencode, playwright, ctx7,
chrome-devtools-mcp, notebooklm). kfg adds only three lines of `shellHook` on top.

The nixai flake in turn pulls kfg as a flake input (`kfg.url = "github:seregatte/kfg"`),
creating a bidirectional dependency: kfg → nixai → kfg. This circular relationship
is resolved by Nix's lock file mechanism (each side pins the other at a fixed
revision), but it means a developer cannot work on kfg without nixai also being
accessible and up to date as a flake input.

## Goals / Non-Goals

**Goals:**

- `flake.nix` declares all dev shell tooling inline; the `nixai` input is removed
- `flake.lock` contains no `nixai` or `nixai_2` entries after regeneration
- `nix develop` enters a fully functional dev shell without requiring nixai
- The `kfg-bundle` package (symlinkJoin of AI agent tools) lives in kfg's own flake
- Starship prompt assets live in `assets/starship/` within kfg
- All residual `nixai`/`.nixai` naming in Go source, tests, and docs is eliminated
- `README.md`, `docs/AGENTS.md`, and `docs/cli-reference.md` reflect the actual
  command surface and dev workflow

**Non-Goals:**

- Changing any manifest YAML, Step definitions, or extension behavior (already migrated)
- Modifying the kfg runtime or shell generation logic
- Changing the nixai repository itself
- Introducing new kfg CLI commands or features

## Decisions

### 1. Inline all nixai content verbatim, adapted minimally

Copy `devInputs`, `npmGlobal`, `gws-bin`, `notebooklmWrapper`, and `nblmWrapper`
from nixai's `flake.nix` into kfg's `flake.nix`. Adapt only what must change:
rename `pname = "nixai-global-npm"` → `"kfg-global-npm"`, rename the bundle from
`nixai-bundle` → `kfg-bundle`, and update the `shellHook` env vars.

**Alternative considered:** Restructure the devShell into multiple shells (e.g.,
a minimal `dev` shell with only Go/bats, and a full `default` shell with AI tools).
**Rejected** because the existing shape (`default` and `dev` both map to the same
full shell) works correctly and changing it is unnecessary scope.

### 2. Remove `kfg-bin` from the bundle to break the circular dependency

nixai's bundle included `kfg-bin = kfg.packages.${system}.default` — the kfg
binary fetched from GitHub Releases. In kfg's own flake, `self.packages.${system}.default`
already is the kfg binary. Including it again in the bundle would be redundant and
cause a store collision in `symlinkJoin`. The `kfg-bundle` therefore contains
everything from nixai's bundle except `kfg-bin`.

### 3. Rename `NIXAI_DIR` → `KFG_DIR`

The `shellHook` exports `NIXAI_DIR=${self.outPath}` so that runtime scripts can
locate nixai's store path (primarily for Starship config). With absorption, the
store path belongs to kfg. Renaming to `KFG_DIR` keeps the semantics consistent
with kfg's naming convention (`KFG_` prefix for all env vars).

### 4. Copy Starship assets into `assets/starship/`

The `shellHook` references `${self.outPath}/assets/starship/full.toml` and
`mobile.toml`. These files must exist in kfg's Nix store derivation output.
Placing them at `assets/starship/` in the kfg source tree makes them part of
`self.outPath` automatically (Nix includes all source files).

### 5. Fix docs in the same change

The docs (`README.md`, `docs/AGENTS.md`, `docs/cli-reference.md`) contain
references to commands that no longer exist (`kfg image`, `kfg workspace`,
`kfg assets`, `kfg sys gc`) and flags/env vars with wrong values. Fixing them
in the same change avoids leaving the repo in a partially inconsistent state.
`docs/manifest-model.md` is accurate and needs no changes.

## Risks / Trade-offs

- **`npmGlobal` network access during Nix build**: The `npmGlobal` derivation runs
  `npm install` at build time, which requires network access. This is the same
  pattern nixai uses today — no change in behavior, but it means the build is not
  fully reproducible without `--impure` or a Nix sandbox exception. → No mitigation
  needed; existing behavior is preserved.

- **`flake.lock` churn**: Removing the `nixai` input and regenerating `flake.lock`
  will also pull in a fresh nixpkgs revision if the current pin drifts. → Run
  `nix flake update` only for the `nixai` input removal, then pin nixpkgs explicitly
  if needed.

- **Starship prompt only works inside `nix develop`**: `STARSHIP_CONFIG` points into
  the Nix store; it has no effect outside the dev shell. This is unchanged from the
  current behavior.

- **`kfg-bundle` becomes the declared package**: After this change, `nix build` on
  kfg's flake will produce `packages.default` (the kfg binary). The bundle is a
  separate package. Callers who were doing `nixai.packages.${system}.default` to get
  the bundle will need to reference kfg directly. → nixai is a personal repo; no
  external callers to break.

## Migration Plan

1. Copy Starship assets → `assets/starship/`
2. Rewrite `flake.nix` (remove nixai input, declare all content inline)
3. Run `nix flake lock --update-input nixai` (or delete nixai entry from lock and
   run `nix flake lock`) to regenerate `flake.lock`
4. Verify `nix develop` enters the shell without errors
5. Update Go source help text and test fixtures
6. Update docs (`README.md`, `docs/AGENTS.md`, `docs/cli-reference.md`)
7. Run `make test` and `make test-bats` to confirm no regressions

## Open Questions

- None. All decisions are resolved based on the current codebase state.

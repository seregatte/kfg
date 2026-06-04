{
  description = "kfg - Declarative shell compiler for YAML manifests";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      version = "0.0.7";

      # Platform-specific SHA-256 hashes (updated by release workflow)
      platformHashes = {
        x86_64-linux   = "sha256-fd6TRt/Ymz1C+Bit+IqjRnVDYrWQA4KILXbHRRPjQVo=";
        aarch64-linux  = "sha256-NgbPNPxEcVqZi6kHuib7UxTcZjWRYAj+CYn4hUJr2zM=";
        x86_64-darwin  = "sha256-Z5sZVPo2azgyoeRRK+dDkCswHCxZVwhxX6b0K8XbgWI=";
        aarch64-darwin = "sha256-fGS0e2kCiY/CdAJMwhG8CtNgAeksNwz3KGplopCqig0=";
      };

      # Map Nix system to GoReleaser archive name components
      platformArchiveNames = {
        x86_64-linux   = "linux_amd64";
        aarch64-linux  = "linux_arm64";
        x86_64-darwin  = "darwin_amd64";
        aarch64-darwin = "darwin_arm64";
      };

      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = import nixpkgs {
            inherit system;
            config.allowUnfree = true;
          };
          archiveName = platformArchiveNames.${system};
          hash = platformHashes.${system};

          # Dev/test dependencies (PATH at runtime, not bundled)
          devInputs = with pkgs; [
            yq-go jq yajsv gomplate
            coreutils findutils gnused gnugrep
            bash bats
            google-cloud-sdk
            uv
          ];

          # npmGlobal derivation (AI agents CLIs)
          npmGlobal =
            let
              packageJson = builtins.toJSON {
                name = "kfg-global-npm";
                version = version;
                dependencies = {
                  "@fission-ai/openspec" = "*";
                  "@mariozechner/pi-coding-agent" = "*";
                  "ctx7" = "*";
                  "chrome-devtools-mcp" = "*";
                };
              };
            in
            pkgs.stdenv.mkDerivation {
              pname = "kfg-global-npm";
              version = version;

              dontUnpack = true;

              buildInputs = [ pkgs.nodejs pkgs.cacert ];

              buildPhase = ''
                mkdir -p $out
                echo '${packageJson}' > package.json

                export NODE_EXTRA_CA_CERTS="${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
                export HOME=$TMPDIR

                npm install --ignore-scripts 2>&1 | tail -20

                cp -r node_modules $out/
                mkdir -p $out/bin
                cd $out/bin
                for bin in ../node_modules/.bin/*; do
                  [ -e "$bin" ] && ln -sf "$bin" "$(basename "$bin")"
                done
              '';

              installPhase = ''
                chmod +x $out/bin/* 2>/dev/null || true
              '';
            };

          # Google Workspace CLI
          gws-bin = let
            gwsVersion = "0.22.5";
            targets = {
              "x86_64-linux" = {
                name = "x86_64-unknown-linux-gnu";
                hash = "sha256-3njs29LxqEzKAGOn7LxEAkD8FLbrzLsX9GRreSqMXB8=";
              };
              "aarch64-linux" = {
                name = "aarch64-unknown-linux-gnu";
                hash = "sha256-lEkCldlYDh6IV05xWgoWKZF0fRLWL4x7jcyCaLbBzqA=";
              };
              "x86_64-darwin" = {
                name = "x86_64-apple-darwin";
                hash = "sha256-Ufm9cxQE1LuibDbi4w3WjFbczR+DTAElLLCxTWplRLI=";
              };
              "aarch64-darwin" = {
                name = "aarch64-apple-darwin";
                hash = "sha256-HSqf/VvJssLEtIYw2vCC+tE9nlfXQZiKLCSO7VYvfaw=";
              };
            };
            target = targets.${system};
          in pkgs.stdenv.mkDerivation {
            pname = "gws";
            version = gwsVersion;
            src = pkgs.fetchurl {
              url = "https://github.com/googleworkspace/cli/releases/download/v${gwsVersion}/google-workspace-cli-${target.name}.tar.gz";
              hash = target.hash;
            };
            nativeBuildInputs = pkgs.lib.optionals pkgs.stdenv.isLinux [ pkgs.autoPatchelfHook ];
            buildInputs = pkgs.lib.optionals pkgs.stdenv.isDarwin [ pkgs.libiconv ];
            sourceRoot = ".";
            installPhase = ''
              mkdir -p $out/bin
              install -m755 gws $out/bin/gws
            '';
          };

          # NotebookLM CLI wrappers
          notebooklmWrapper = pkgs.writeShellScriptBin "notebooklm" ''
            export PLAYWRIGHT_BROWSERS_PATH="$HOME/.cache/notebooklm-playwright"
            exec uvx --from 'notebooklm-py[browser]' notebooklm "$@"
          '';

          nblmWrapper = pkgs.writeShellScriptBin "nblm" ''
            export PLAYWRIGHT_BROWSERS_PATH="$HOME/.cache/notebooklm-playwright"
            exec uvx --from 'notebooklm-py[browser]' notebooklm "$@"
          '';

          # kfg-bundle: symlinkJoin of AI agent tools (no kfg-bin — avoid circular reference)
          kfg-bundle = pkgs.symlinkJoin {
            name = "kfg-bundle";
            paths = [
              npmGlobal
              gws-bin
              notebooklmWrapper
              nblmWrapper
              pkgs.claude-code
              pkgs.gemini-cli-bin
              pkgs.opencode
              pkgs.playwright-test
            ];
          };
        in
        {
          default = pkgs.stdenv.mkDerivation {
            pname = "kfg";
            version = version;

            src = pkgs.fetchurl {
              url = "https://github.com/seregatte/kfg/releases/download/v${version}/kfg_${version}_${archiveName}.tar.gz";
              hash = hash;
            };

            sourceRoot = ".";

            installPhase = ''
              runHook preInstall
              mkdir -p $out/bin
              install -m755 kfg $out/bin/kfg
              runHook postInstall
            '';

            meta = with pkgs.lib; {
              description = "Declarative shell compiler for YAML manifests";
              homepage = "https://github.com/seregatte/kfg";
              license = licenses.mit;
              mainProgram = "kfg";
              platforms = supportedSystems;
            };
          };

          kfg-bundle = kfg-bundle;
        });

      devShells = forAllSystems (system:
        let
          pkgs = import nixpkgs {
            inherit system;
            config.allowUnfree = true;
          };
          kfg-bundle = self.packages.${system}.kfg-bundle;

          # Dev/test dependencies (PATH at runtime, not bundled)
          devInputs = with pkgs; [
            yq-go jq yajsv gomplate
            coreutils findutils gnused gnugrep
            bash bats
            google-cloud-sdk
            uv
          ];
        in
        {
          default = pkgs.mkShell {
            buildInputs = devInputs ++ [ kfg-bundle ];
            shellHook = ''
              export KFG_DIR=${self.outPath}
              if [ "$COLUMNS" -lt 45 ] 2>/dev/null; then
                export STARSHIP_CONFIG=${self.outPath}/assets/starship/mobile.toml
              else
                export STARSHIP_CONFIG=${self.outPath}/assets/starship/full.toml
              fi
            '';
          };

          dev = pkgs.mkShell {
            buildInputs = devInputs ++ [ pkgs.nodejs pkgs.go kfg-bundle ];
            shellHook = ''
              export KFG_DIR=${self.outPath}
              export PATH="./bin:$PATH"
              export OPENSPEC_ROOT_DIR=docs/context
              if [ "$COLUMNS" -lt 45 ] 2>/dev/null; then
                export STARSHIP_CONFIG=${self.outPath}/assets/starship/mobile.toml
              else
                export STARSHIP_CONFIG=${self.outPath}/assets/starship/full.toml
              fi
              source <(go run ./src/cmd/kfg apply -k packages/domains/ai-agents/overlays/dev)
            '';
          };

          # Minimal devShell for CI — no kfg-bundle (avoids broken gws-bin on Linux).
          ci = pkgs.mkShell {
            buildInputs = devInputs ++ [ pkgs.go pkgs.gnumake ];
            shellHook = ''
              export PATH="./bin:$PATH"
              export OPENSPEC_ROOT_DIR=docs/context
              # Set up vendor directory for bats test helpers
              VENDOR_DIR=tests/bats/helpers/vendor
              rm -rf "$VENDOR_DIR/bats-support" "$VENDOR_DIR/bats-assert"
              mkdir -p "$VENDOR_DIR/bats-support" "$VENDOR_DIR/bats-assert"
              # Fetch bats-support and copy only needed files (avoid repo's own tests)
              BATS_SUPPORT=${pkgs.fetchFromGitHub {
                owner = "bats-core";
                repo = "bats-support";
                rev = "v0.3.0";
                hash = "sha256-4N7XJS5XOKxMCXNC7ef9halhRpg79kUqDuRnKcrxoeo=";
              }}
              cp "$BATS_SUPPORT/load.bash" "$VENDOR_DIR/bats-support/"
              cp -r "$BATS_SUPPORT/src" "$VENDOR_DIR/bats-support/"
              # Fetch bats-assert and copy only needed files
              BATS_ASSERT=${pkgs.fetchFromGitHub {
                owner = "bats-core";
                repo = "bats-assert";
                rev = "v2.1.0";
                hash = "sha256-opgyrkqTwtnn/lUjMebbLfS/3sbI2axSusWd5i/5wm4=";
              }}
              cp "$BATS_ASSERT/load.bash" "$VENDOR_DIR/bats-assert/"
              cp -r "$BATS_ASSERT/src" "$VENDOR_DIR/bats-assert/"
            '';
          };
        });

      lib.version = version;
    };
}

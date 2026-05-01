{
  description = "kfg - Declarative shell compiler for YAML manifests";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    nixai.url = "github:seregatte/nixai";
  };

  outputs = { self, nixpkgs, nixai }:
    let
      version = "0.0.1";
      
      # Platform-specific SHA-256 hashes (updated by release workflow)
      platformHashes = {
        x86_64-linux   = "sha256-y7KRkkz6K4mMXLFupagwP0ItFdtJVWuzNijnwbf6gLk=";
        aarch64-linux  = "sha256-tYLfnevm4e/usoz8KiadZLA0vfwyXrCQYB5GPXMJ3a0=";
        x86_64-darwin  = "sha256-SBjTPNpIWxY0TTIjLIjFLTgBnnySU2xIBwXL9893aHc=";
        aarch64-darwin = "sha256-SBLolv6X7329o6cy2M1Z1UuoNZqb68ENlwF1Cb84z5c=";
      };
      
      # Map Nix system to GoReleaser archive name components
      platformArchiveNames = {
        x86_64-linux   = "sha256-y7KRkkz6K4mMXLFupagwP0ItFdtJVWuzNijnwbf6gLk=";
        aarch64-linux  = "sha256-tYLfnevm4e/usoz8KiadZLA0vfwyXrCQYB5GPXMJ3a0=";
        x86_64-darwin  = "sha256-SBjTPNpIWxY0TTIjLIjFLTgBnnySU2xIBwXL9893aHc=";
        aarch64-darwin = "sha256-SBLolv6X7329o6cy2M1Z1UuoNZqb68ENlwF1Cb84z5c=";
      };
      
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          archiveName = platformArchiveNames.${system};
          hash = platformHashes.${system};
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
        });
      
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          shell = pkgs.mkShell {
            inputsFrom = [ nixai.devShells.${system}.default ];
            shellHook = ''
              echo "Welcome to the kfg development environment!"
              alias kfg="go run ./src/cmd/kfg"
              ln -s docs/context/openspec ./ 2>/dev/null || true
              source <(kfg apply -k $NIXAI_DIR/.nixai/overlay/dev)
            '';
          };
        in
        {
          default = shell;
          dev = shell;
        });
      
      lib.version = version;
    };
}
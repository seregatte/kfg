{
  description = "kfg - Declarative shell compiler for YAML manifests";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      version = "2.1.0";
      
      # Platform-specific SHA-256 hashes (updated by release workflow)
      platformHashes = {
        x86_64-linux   = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
        aarch64-linux  = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
        x86_64-darwin  = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
        aarch64-darwin = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
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
      
      lib.version = version;
    };
}
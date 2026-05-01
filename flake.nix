{
  description = "kfg - Declarative shell compiler for YAML manifests";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    nixai.url = "github:seregatte/nixai";
  };

  outputs = { self, nixpkgs, nixai }:
    let
      version = "2.1.0";
      
      # Platform-specific SHA-256 hashes (updated by release workflow)
      platformHashes = {
        x86_64-linux   = "sha256-/OMq8i2rbYyicQkzvPe32mOoOv1F0hLfagI8X+fxHLg=";
        aarch64-linux  = "sha256-y4/dx1OSDU8+IH9dbthMgpcvGoxY5rfXFWHeLRj5SN0=";
        x86_64-darwin  = "sha256-qutAYfjJo5E01r+G/IMGRrhVUSKVOdXI5v4cTxVL7kA=";
        aarch64-darwin = "sha256-YTXbTi6fvOdoSJfhbbqtpIpLvhynpTBBwq3RkYD/J20=";
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
      
      devShells = forAllSystems (system:
        {
          default = nixai.devShells.${system}.default;
        });
      
      lib.version = version;
    };
}
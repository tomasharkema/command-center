{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    devenv.url = "github:cachix/devenv";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs = {
    self,
    nixpkgs,
    devenv,
    systems,
    ...
  } @ inputs: let
    forEachSystem = nixpkgs.lib.genAttrs (import systems);
  in {
    packages = forEachSystem (system: let
      lib = nixpkgs.lib;
      pkgs = nixpkgs.legacyPackages.${system};
    in rec {
      devenv-up = self.devShells.${system}.default.config.procfileScript;

      nixosModules.command-center = import ./modules/command-center.nix;
      nixosModules.default = self.nixosModules.command-center;

      default = command-center;
      command-center = pkgs.callPackage ./package.nix {};
    });

    devShells =
      forEachSystem
      (system: let
        pkgs = nixpkgs.legacyPackages.${system};
        lib = nixpkgs.lib;
      in {
        default = devenv.lib.mkShell {
          inherit inputs pkgs;
          modules = [
            {
              packages = with pkgs; [
                git
                go
                gopls
                golangci-lint
                watchexec
                # docker
              ];

              # scripts.build.exec = ''
              #   ${pkgs.go} build -v
              # '';
              processes.run.exec = ''
                ${lib.getExe pkgs.watchexec} -r -e js,css,html,go,nix -- go run . --listen=:3334 --verbose
              '';
              # scripts.docker-build.exec = ''
              #   ${pkgs.docker} buildx build --platform=linux/amd64,linux/arm64 . -t ghcr.io/tomasharkema/command-center
              # '';

              languages = {
                go.enable = true;
              };

              dotenv.enable = true;

              pre-commit.hooks = {
                shellcheck.enable = true;
                gofmt.enable = true;
                golangci-lint.enable = true;
              };
            }
          ];
        };
      });
  };
}

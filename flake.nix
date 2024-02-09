{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable-small";
    systems = {
      url = "github:nix-systems/default";
    };
    devenv = {
      url = "github:cachix/devenv";
    };
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs = inputs:
    inputs.flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin"];

      perSystem = {
        self',
        config,
        lib,
        pkgs,
        system,
        ...
      }: {
        overlayAttrs = {
          inherit (config.packages) command-center;
        };
        packages = {
          devenv-up = self'.devShells.${system}.default.config.procfileScript;

          default = config.packages.command-center;
          command-center = pkgs.callPackage ./pkgs/command-center.nix {};
        };
      };
      flake = {
        nixosModules = rec {
          command-center = import ./modules/command-center.nix;
          default = command-center;
        };
        # nixosConfigurations."test-server" = inputs.nixpkgs.lib.nixosSystem {
        #   system = "x86_64-linux";
        #   modules = [
        #     # config.nixosModules.command-center
        #     {
        #       # imports = [config.nixosModules.command-center];

        #       boot.isContainer = true;
        #       services.command-center.enable = true;
        #     }
        #   ];
        # };
      };
    };

  # outputs = {
  #   self,
  #   nixpkgs,
  #   devenv,
  #   systems,
  #   ...
  # } @ inputs: let
  #   forEachSystem = nixpkgs.lib.genAttrs (import systems);
  # in {
  #   overlays.default = import ./overlay.nix;

  #   packages = forEachSystem (system: let
  #     lib = nixpkgs.lib;
  #     pkgs = nixpkgs.legacyPackages.${system};
  #   in rec {
  #     devenv-up = self.devShells.${system}.default.config.procfileScript;

  #     default = command-center;
  #     command-center = pkgs.callPackage ./pkgs/command-center.nix {};
  #   });

  #   nixosConfigurations.test-server = nixpkgs.lib.nixosSystem {
  #     system = "x86_64-linux";
  #     modules = [
  #       self.nixosModules.command-center
  #       {
  #         boot.isContainer = true;
  #       }
  #     ];
  #   };

  #   devShells =
  #     forEachSystem
  #     (system: let
  #       pkgs = nixpkgs.legacyPackages.${system};
  #       lib = nixpkgs.lib;
  #     in {
  #       default = devenv.lib.mkShell {
  #         inherit inputs pkgs;
  #         modules = [
  #           {
  #             packages = with pkgs; [
  #               git
  #               go
  #               gopls
  #               golangci-lint
  #               watchexec
  #               # docker
  #             ];

  #             # scripts.build.exec = ''
  #             #   ${pkgs.go} build -v
  #             # '';
  #             processes.run.exec = ''
  #               ${lib.getExe pkgs.watchexec} -r -e js,css,html,go,nix -- go run . --listen=:3334 --verbose
  #             '';
  #             # scripts.docker-build.exec = ''
  #             #   ${pkgs.docker} buildx build --platform=linux/amd64,linux/arm64 . -t ghcr.io/tomasharkema/command-center
  #             # '';

  #             languages = {
  #               go.enable = true;
  #             };

  #             dotenv.enable = true;

  #             pre-commit.hooks = {
  #               shellcheck.enable = true;
  #               gofmt.enable = true;
  #               golangci-lint.enable = true;
  #             };
  #           }
  #         ];
  #       };
  #     });
  # };
}

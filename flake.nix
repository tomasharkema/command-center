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

      default = go-nixos-menu;
      go-nixos-menu = pkgs.buildGoModule rec {
        pname = "go-nixos-menu";
        version = "0.0.1-alpha1";

        CGO_ENABLED = 0;
        vendorHash = "sha256-YopZDBl9XUWpON7sRjs403lWdpN1I3zezv1eiGv0ziw=";

        src = ./.;

        meta = with lib; {
          description = "tomas";
          homepage = "https://github.com/tomasharkema/go-nixos-menu";
          license = licenses.mit;
          maintainers = ["tomasharkema" "tomas@harkema.io"];
          mainProgram = pname;
        };
      };
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
                ${lib.getExe pkgs.watchexec} -r -e js,css,html,go,nix -- go run . -- -v
              '';
              # scripts.docker-build.exec = ''
              #   ${pkgs.docker} buildx build --platform=linux/amd64,linux/arm64 . -t ghcr.io/tomasharkema/go-nixos-menu
              # '';

              languages = {
                go.enable = true;
              };

              dotenv.enable = true;

              pre-commit.hooks = {
                shellcheck.enable = true;
                gofmt.enable = true;
                # golangci-lint.enable = true;
              };
            }
          ];
        };
      });
  };
}

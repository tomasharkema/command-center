{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    devenv.url = "github:cachix/devenv";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    devenv,
  } @ inputs:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      packages = rec {
        default = go-nixos-menu;
        go-nixos-menu = pkgs.buildGoModule rec {
          pname = "go-nixos-menu";
          version = "0.0.1";

          CGO_ENABLED = 0;
          vendorHash = "sha256-hkUM1L357lDMajBD0go6y6KfKEUrlP0+1RUdfflzXqE=";

          src = ./.;

          meta =
            #with lib;
            {
              description = "tomas";
              homepage = "https://github.com/tomasharkema/go-nixos-menu";
              # license = licenses.mit;
              maintainers = ["tomasharkema" "tomas@harkema.io"];
            };
        };
      };
      devShell = devenv.lib.mkShell {
        inherit inputs pkgs;
        modules = [
          ({
            pkgs,
            config,
            ...
          }: {
            # This is your devenv configuration

            # https://devenv.sh/basics/
            env.GREET = "devenv";

            # https://devenv.sh/packages/
            packages = with pkgs; [
              git
              go
              gopls
              golangci-lint
              watchexec
              docker
            ];

            # https://devenv.sh/scripts/
            scripts.build.exec = ''
              ${pkgs.go} build -v
            '';
            scripts.dev.exec = ''
              ${pkgs.watchexec} -e js,css,html,go go run .
            '';
            scripts.docker-build.exec = ''
              ${pkgs.docker} buildx build --platform=linux/amd64,linux/arm64 . -t ghcr.io/tomasharkema/go-nixos-menu
            '';

            languages = {
              go.enable = true;
            };

            dotenv.enable = true;

            # https://devenv.sh/pre-commit-hooks/
            # pre-commit.hooks.shellcheck.enable = true;

            # https://devenv.sh/processes/
            # processes.ping.exec = "ping example.com";

            # See full reference at https://devenv.sh/reference/options/
          })
        ];
      };
    });
}

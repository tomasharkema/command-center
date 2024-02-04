{
  pkgs,
  lib,
  ...
}: {
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
    ${lib.getExe pkgs.go} build -v
  '';
  scripts.dev.exec = ''
    ${lib.getExe pkgs.watchexec} -e js,css,html,go go run .
  '';
  scripts.docker-build.exec = ''
    ${lib.getExe pkgs.docker} buildx build --platform=linux/amd64,linux/arm64 . -t ghcr.io/tomasharkema/go-nixos-menu
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
}

{
  config,
  options,
  lib,
  pkgs,
  ...
}:
with lib; let
  cfg = config.services.command-center;
in {
  options.services.command-center = {
    enable = mkOption {
      type = types.bool;
      default = true;
      description = "Enable Command-Center";
    };
    enableBot = mkOption {
      type = types.bool;
      default = false;
      description = "Enable Command-Center bot";
    };
    envPath = mkOption {
      type = types.str;
      default = null;
      description = "Command-Center env path";
    };
  };
  config = let
    info-json = builtins.toFile "info.json" (builtins.toJSON {
      tags = config.system.nixos.tags;
      revision = config.system.configurationRevision;
      version = config.system.stateVersion;
      label = config.system.nixos.label;
      name = config.system.name;
    });

    command-center-service = {
      enable = true;
      description = "command-center";
      environment = {
        IS_NIXOS_SYSTEMD = true;
        COMMAND_CENTER_RUN_BOT = mkIf cfg.enableBot "true";
        INFO_JSON_PATH = info-json;
      };
      unitConfig = {
        Type = "notify";
      };
      serviceConfig = {
        RestartSec = 5;
        EnvironmentFile = cfg.envPath;
      };
      script = "${lib.getExe pkgs.command-center} -v";
      wantedBy = ["multi-user.target" "network.target"];
    };
  in
    mkIf cfg.enable {
      systemd.services = {
        command-center = command-center-service;
      };

      # proxy-services.services = {
      #   "/dashboard/" = {
      #     proxyPass = "http://localhost:3333/";
      #   };
      # };

      # age.secrets."command-center.env" = {
      #   file = ../../../../secrets/command-center.env.age;
      #   mode = "644";
      # };

      networking.nat = mkIf false {
        enable = true;
        internalInterfaces = ["ve-+"];
        # externalInterface = "ens3";
      };

      containers.go-nixos-menu = mkIf false {
        autoStart = true;
        privateNetwork = true;
        hostAddress = "192.168.100.10";
        localAddress = "192.168.100.11";

        config = {
          config,
          pkgs,
          ...
        }: {
          systemd.services = {
          };

          system.stateVersion = "23.11";

          networking = {
            firewall = {
              enable = true;
              allowedTCPPorts = [3333];
            };
            # Use systemd-resolved inside the container
            # Workaround for bug https://github.com/NixOS/nixpkgs/issues/162686
            useHostResolvConf = mkForce false;
          };

          services.resolved.enable = true;
        };
      };
    };
}

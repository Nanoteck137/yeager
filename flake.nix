{
  description = "Devshell for yeager";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";

    pyrin.url        = "github:nanoteck137/pyrin/v0.7.0";
    pyrin.inputs.nixpkgs.follows = "nixpkgs";

    devtools.url     = "github:nanoteck137/devtools";
    devtools.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, pyrin, devtools, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        version = pkgs.lib.strings.fileContents "${self}/version";
        fullVersion = ''${version}-${self.dirtyShortRev or self.shortRev or "dirty"}'';

        app = pkgs.buildGoModule {
          pname = "yeager";
          version = fullVersion;
          src = ./.;

          ldflags = [
            "-X github.com/nanoteck137/yeager/config.Version=${version}"
            "-X github.com/nanoteck137/yeager/config.Commit=${self.dirtyRev or self.rev or "no-commit"}"
          ];

          vendorHash = "";
        };

        tools = devtools.packages.${system};
      in
      {
        packages.default = app;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            air
            go
            gopls
            nodejs

            pyrin.packages.${system}.default
            tools.publishVersion
          ];
        };
      }
    ) // {
      nixosModules.default = { config, lib, pkgs, ... }:
        with lib; let
          cfg = config.services.yeager;

          yeagerConfig = pkgs.writeText "config.toml" ''
            listen_addr = "${cfg.host}:${toString cfg.port}"
            data_dir = "/var/lib/yeager"
            library_dir = "${cfg.library}"
            username = "${cfg.username}"
            initial_password = "${cfg.initialPassword}"
            jwt_secret = "${cfg.jwtSecret}"
          '';
        in
        {
          options.services.yeager = {
            enable = mkEnableOption "Enable the yeager service";

            port = mkOption {
              type = types.port;
              default = 7550;
              description = "port to listen on";
            };

            host = mkOption {
              type = types.str;
              default = "";
              description = "hostname or address to listen on";
            };

            library = mkOption {
              type = types.path;
              description = "path to series library";
            };

            username = mkOption {
              type = types.str;
              description = "username of the first user";
            };

            initialPassword = mkOption {
              type = types.str;
              description = "initial password of the first user (should change after the first login)";
            };

            jwtSecret = mkOption {
              type = types.str;
              description = "jwt secret";
            };

            package = mkOption {
              type = types.package;
              default = self.packages.${pkgs.system}.default;
              description = "package to use for this service (defaults to the one in the flake)";
            };

            user = mkOption {
              type = types.str;
              default = "yeager";
              description = "user to use for this service";
            };

            group = mkOption {
              type = types.str;
              default = "yeager";
              description = "group to use for this service";
            };

          };

          config = mkIf cfg.enable {
            systemd.services.yeager = {
              description = "yeager";
              wantedBy = [ "multi-user.target" ];

              serviceConfig = {
                User = cfg.user;
                Group = cfg.group;

                StateDirectory = "yeager";

                ExecStart = "${cfg.package}/bin/yeager serve -c '${yeagerConfig}'";

                Restart = "on-failure";
                RestartSec = "5s";

                ProtectHome = true;
                ProtectHostname = true;
                ProtectKernelLogs = true;
                ProtectKernelModules = true;
                ProtectKernelTunables = true;
                ProtectProc = "invisible";
                ProtectSystem = "strict";
                RestrictAddressFamilies = [ "AF_INET" "AF_INET6" "AF_UNIX" ];
                RestrictNamespaces = true;
                RestrictRealtime = true;
                RestrictSUIDSGID = true;
              };
            };

            users.users = mkIf (cfg.user == "yeager") {
              yeager = {
                group = cfg.group;
                isSystemUser = true;
              };
            };

            users.groups = mkIf (cfg.group == "yeager") {
              yeager = {};
            };
          };
        };
    };
}

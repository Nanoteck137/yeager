{
  description = "Devshell for yeager-frontend";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";

    pyrin.url        = "github:nanoteck137/pyrin/v0.6.5";
    pyrin.inputs.nixpkgs.follows = "nixpkgs";

    devtools.url     = "github:nanoteck137/devtools";
    devtools.inputs.nixpkgs.follows = "nixpkgs";

    gitignore.url = "github:hercules-ci/gitignore.nix";
    gitignore.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, gitignore, pyrin, devtools, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        version = pkgs.lib.strings.fileContents "${self}/version";
        rev = self.dirtyShortRev or self.shortRev or "dirty";
        fullVersion = ''${version}-${rev}'';

        app = pkgs.buildNpmPackage {
          name = "yeager-frontend";
          version = fullVersion;

          src = gitignore.lib.gitignoreSource ./.;
          npmDepsHash = "sha256-0R0zcjevu3yrKC/+7nJsOq4eXEpE0/Y/4IJ/lgtU9oY=";

          PUBLIC_VERSION=version;
          PUBLIC_COMMIT=self.rev or "dirty";

          installPhase = ''
            runHook preInstall
            cp -r build $out/
            echo '{ "type": "module" }' > $out/package.json

            mkdir $out/bin
            echo -e "#!${pkgs.runtimeShell}\n${pkgs.nodejs}/bin/node $out\n" > $out/bin/yeager-frontend
            chmod +x $out/bin/yeager-frontend

            runHook postInstall
          '';
        };

        tools = devtools.packages.${system};
      in
      {
        packages.default = app;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            nodejs
            python3
            
            pyrin.packages.${system}.default
            tools.publishVersion
          ];
        };
      }
    );
}

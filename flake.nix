{
  description = "Fabric is an open-source framework for augmenting humans using AI. It provides a modular framework for solving specific problems using a crowdsourced set of AI prompts that can be used anywhere";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      systems,
      treefmt-nix,
      gomod2nix,
      ...
    }:
    let
      forAllSystems = nixpkgs.lib.genAttrs (import systems);

      treefmtEval = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        treefmt-nix.lib.evalModule pkgs ./treefmt.nix
      );
    in
    {
      formatter = forAllSystems (system: treefmtEval.${system}.config.build.wrapper);

      checks = forAllSystems (system: {
        formatting = treefmtEval.${system}.config.build.check self;
      });

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          goEnv = gomod2nix.legacyPackages.${system}.mkGoEnv { pwd = ./.; };
        in
        {
          default = pkgs.mkShell {
            nativeBuildInputs = [
              pkgs.go
              pkgs.gopls
              pkgs.gotools
              pkgs.go-tools
              pkgs.goimports-reviser
              gomod2nix.legacyPackages.${system}.gomod2nix
              goEnv

              (pkgs.writeShellScriptBin "update" ''
                go get -u
                go mod tidy
                gomod2nix generate
              '')
            ];

            shellHook = ''
              echo -e "\033[0;32;4mHeper commands:\033[0m"
              echo "'update' instead of 'go get -u && go mod tidy'"
            '';
          };
        }
      );

      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          fabric = gomod2nix.legacyPackages.${system}.buildGoApplication {
            pname = "fabric-ai";
            version = "1.4.85";
            src = self;
            pwd = self;
            modules = ./gomod2nix.toml;

            ldflags = [
              "-s"
              "-w"
            ];

            meta = with pkgs.lib; {
              description = "Fabric is an open-source framework for augmenting humans using AI. It provides a modular framework for solving specific problems using a crowdsourced set of AI prompts that can be used anywhere";
              homepage = "https://github.com/danielmiessler/fabric";
              license = licenses.mit;
              platforms = platforms.all;
              mainProgram = "fabric";
            };
          };
          default = self.packages.${system}.fabric;
        }
      );
    };
}

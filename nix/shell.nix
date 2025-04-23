{
  pkgs,
  gomod2nix,
  goEnv,
  goVersion,
}:

{
  default = pkgs.mkShell {
    nativeBuildInputs = [
      goVersion
      pkgs.gopls
      pkgs.gotools
      pkgs.go-tools
      pkgs.goimports-reviser
      gomod2nix
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

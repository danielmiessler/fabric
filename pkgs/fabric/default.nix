{
  lib,
  buildGoApplication,
}:

buildGoApplication {
  pname = "fabric-ai";
  version = import ./version.nix;
  src = ../../.;
  pwd = ../../.;
  modules = ../../gomod2nix.toml;

  ldflags = [
    "-s"
    "-w"
  ];

  meta = with lib; {
    description = "Fabric is an open-source framework for augmenting humans using AI. It provides a modular framework for solving specific problems using a crowdsourced set of AI prompts that can be used anywhere";
    homepage = "https://github.com/danielmiessler/fabric";
    license = licenses.mit;
    platforms = platforms.all;
    mainProgram = "fabric";
  };
}

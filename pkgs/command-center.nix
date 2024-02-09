{
  pkgs,
  lib,
  ...
}:
pkgs.buildGoModule rec {
  pname = "command-center";
  version = "0.0.1-alpha5";

  CGO_ENABLED = 0;
  vendorHash = "sha256-bDN8jEAKGRyo7XrDmIZrp1Bx2VBmWl3KUjH/Ue5FvEc=";

  src = ./..;

  meta = with lib; {
    description = "tomas";
    homepage = "https://github.com/tomasharkema/command-center";
    license = licenses.mit;
    maintainers = ["tomasharkema" "tomas@harkema.io"];
    mainProgram = pname;
  };
}

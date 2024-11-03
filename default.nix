{
  pkgs ? import <nixpkgs> { },
}:

pkgs.buildGoModule rec {
  pname = "github.com/GuillaumeDesforges/ncbuild";
  version = "v0.1.0";
  src = ./.;
  vendorHash = "sha256-xssAi8dPuL2H/HZPc4bVCF8Wr1PKQ/d1IbQ02gLzxCA=";
}

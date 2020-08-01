{ sources ? import ./sources.nix, pkgs ? import sources.nixpkgs { } }:

version:

let src = sources.zig;

in pkgs.stdenv.mkDerivation {
  name = "zig";
  inherit src version;

  installPhase = ''
    mkdir -p $out
    cp -rf * $out
    mkdir -p $out/bin
    mv $out/zig $out/bin/zig
  '';
}

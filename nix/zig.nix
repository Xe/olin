{ sources ? import ./sources.nix, pkgs ? import sources.nixpkgs { } }:
with pkgs;
stdenv.mkDerivation {
  name = "zig-nightly";
  version = sources.zig.version;
  src = sources.zig;

  installPhase = ''
    mkdir -p $out
    cp -rf * $out
    mkdir -p $out/bin
    mv $out/zig $out/bin/zig
  '';
}

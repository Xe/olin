{ sources ? import ./nix/sources.nix, pkgs ? import sources.nixpkgs { } }:
with pkgs;

assert lib.versionAtLeast go.version "1.13";

let
  name = "olin";
  version = "v0.4.0";
  src = ./.;

  olin = buildGoPackage rec {
    inherit name version src;
    goPackagePath = "within.website/olin";
    goDeps = ./nix/deps.nix;
    allowGoReference = false;
  };

  zigFiles = import ./zig { inherit sources pkgs; };

in stdenv.mkDerivation {
  inherit name version;
  phases = "installPhase";

  installPhase = ''
    mkdir -p $out/bin
    cp ${olin}/bin/cwa $out/bin/cwa
    cp -rf ${zigFiles}/wasm $out/wasm
  '';
}

{ sources ? import ../nix/sources.nix, pkgs ? import sources.nixpkgs { }, zig ? import ../nix/zig.nix {inherit sources pkgs;} }:
let
in with pkgs;
stdenv.mkDerivation {
  name = "olin-zig-files";
  version = "latest";
  src = ./.;

  buildInputs = [ zig ];
  phases = "buildPhase installPhase";

  buildPhase = ''
    cp -rf $src/src .
    $src/build.sh
  '';

  installPhase = ''
    mkdir -p $out/wasm/zig
    cp -rf $src/*.policy $out/wasm/zig
    cp -rf *.wasm $out/wasm/zig
  '';
}

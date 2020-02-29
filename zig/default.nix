{ sources ? import ../nix/sources.nix, pkgs ? import sources.nixpkgs { }
, xepkgs ? import sources.xepkgs { inherit sources pkgs; } }:
let
  zig = xepkgs.zig "0.5.0+bee4007ec" {
    mac = "3481cf1e70ebf75b08cac435d66d6d80d5e107cbbe02bf8aca177b4896540555";
    linux = "2b967b38907296f560e97c34f37cd8fc63febc57560cce2d61cf84bc80b91a21";
  };
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

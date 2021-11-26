{ sources ? import ./sources.nix, pkgs ? import sources.nixpkgs { } }:

(import ./zig_raw.nix {
  inherit sources pkgs;
}) "0.6.0"

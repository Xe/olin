{ sources ? import ./sources.nix, pkgs ? import sources.nixpkgs { } }:

(import ./zig_raw.nix { }) "0.6.0"

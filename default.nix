{ pkgs ? import (import ./nix/sources.nix).nixpkgs { }}:
with pkgs;

assert lib.versionAtLeast go.version "1.13";

buildGoPackage rec {
  name = "olin";
  version = "v0.4.0";
  goPackagePath = "within.website/olin";
  src = ./.;
  goDeps = ./nix/deps.nix;
  allowGoReference = false;
}

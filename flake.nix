{
  inputs = {
    flake-compat.url = "github:edolstra/flake-compat";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    git-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs = {
        flake-compat.follows = "flake-compat";
        nixpkgs.follows = "nixpkgs";
        nixpkgs-stable.follows = "nixpkgs";
      };
    };
    make-shell = {
      url = "github:nicknovitski/make-shell";
      inputs.flake-compat.follows = "flake-compat";
    };
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    treefmt = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        ./nix/build-parameters.nix
        ./nix/c-library.nix
        ./nix/debian.nix
        ./nix/fmt.nix
        ./nix/git-hooks.nix
        ./nix/package-checks.nix
        ./nix/python-package.nix
        ./nix/shell.nix
        ./nix/systems.nix
        ./nix/version.nix
      ];

      perSystem = {
        treefmt.settings.global.excludes = [ "./src/VERSION" ];
      };
    };
}

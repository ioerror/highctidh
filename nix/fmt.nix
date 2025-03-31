{ inputs, ... }:
{
  imports = [
    inputs.treefmt.flakeModule
  ];
  perSystem = {
    pre-commit.settings.hooks.nix-fmt = {
      enable = true;
      entry = "nix fmt -- --fail-on-change";
    };

    treefmt = {
      projectRootFile = "flake.nix";
      settings.global.excludes = [
        "COPYING,"
        "LICENSE"
      ];
      programs = {
        nixfmt.enable = true;
      };
    };

  };
}

{ inputs, ... }:
{
  imports = [ inputs.git-hooks.flakeModule ];
  perSystem =
    { config, ... }:
    {
      pre-commit.check.enable = false;
      make-shell.imports = [
        {
          shellHook = config.pre-commit.installationScript;
        }
      ];
    };
}

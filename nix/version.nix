{ lib, ... }:
{
  options.version = lib.mkOption { type = lib.types.singleLineStr; };
  config = {
    version = lib.pipe ../src/VERSION [
      builtins.readFile
      lib.trim
    ];
    perSystem.treefmt.settings.global.excludes = [ "./src/VERSION" ];
  };
}

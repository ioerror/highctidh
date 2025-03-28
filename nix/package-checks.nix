{ lib, ... }:
{
  perSystem = perSystemArgs: {
    checks = lib.pipe perSystemArgs.config.packages [
      (lib.mapAttrs' (name: lib.nameValuePair "package/${name}"))
    ];
  };
}

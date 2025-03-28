{ lib, ... }:
let
  systems-with-assembly-backend = [ "x86_64-linux" ];
in
{
  options = with lib.types; {
    variantFactors = {
      stdenvVariant = lib.mkOption {
        type = listOf (submodule {
          options = {
            name = lib.mkOption { type = singleLineStr; };
            getPkg = lib.mkOption { type = functionTo package; };
          };
        });
      };
      fieldSize = lib.mkOption { type = listOf int; };
      backend = lib.mkOption {
        type = listOf (enum [
          "portable"
          "assembly"
        ]);
      };
    };
    getBackendEnv = lib.mkOption {
      type = functionTo (submodule {
        options.HIGHCTIDH_PORTABLE = lib.mkOption {
          type = enum [
            "0"
            "1"
          ];
        };
      });
    };
    getVariantName = lib.mkOption { type = functionTo str; };
    filterBackendBySystem = lib.mkOption { type = functionTo (functionTo bool); };
  };

  config = {
    variantFactors = {
      stdenvVariant =
        lib.pipe
          [ "stdenv" "clangStdenv" ]
          [
            (map (name: {
              inherit name;
              getPkg =
                pkgs:
                assert pkgs._type == "pkgs";
                pkgs.${name};
            }))
          ];
      fieldSize = [
        511
        512
        1024
        2048
      ];
      backend = [
        "portable"
        "assembly"
      ];
    };
    getBackendEnv = backend: {
      HIGHCTIDH_PORTABLE =
        {
          assembly = "0";
          portable = "1";
        }
        .${backend};
    };

    getVariantName =
      {
        stdenvVariant,
        fieldSize ? null,
        backend,
      }:
      lib.pipe
        [ fieldSize stdenvVariant.name "backend_${backend}" ]
        [
          (lib.filter (v: v != null))
          (map toString)
          (lib.concatStringsSep "-")
        ];

    filterBackendBySystem =
      system: { backend, ... }: backend == "portable" || lib.elem system systems-with-assembly-backend;
  };
}

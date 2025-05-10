{ lib, config, ... }:
{
  perSystem =
    { pkgs, system, ... }:
    let
      inherit (config) variantFactors;

      mkPackageModule =
        tuple@{
          stdenvVariant,
          fieldSize,
          backend,
        }:
        let
          filename = "libhighctidh_${toString fieldSize}.so";
          variantName = config.getVariantName tuple;
          stdenv = (stdenvVariant.getPkg pkgs);

          package = stdenv.mkDerivation {
            pname = "libhighctidh_${variantName}";
            inherit (config) version;

            src = lib.fileset.toSource {
              root = ../src;
              # can be filtered further, but as of making this change, seems too tedious without changes to source tree
              fileset =
                with lib.fileset;
                unions [
                  (fileFilter (
                    file:
                    lib.any file.hasExt (
                      lib.optional (backend == "assembly") "S"
                      ++ lib.optional (backend == "assembly") "s"
                      ++ [
                        "c"
                        "h"
                      ]
                    )
                  ) ../src)
                  ../src/VERSION
                  ../src/GNUmakefile
                ];
            };

            buildFlags = [ filename ];

            doCheck = true;
            checkTarget = lib.concatStringsSep " " [
              "testrandom"
              "test${toString fieldSize}"
            ];

            installPhase = ''
              $preInstall

              mkdir -p $out/include/libhighctidh
              cp *.h $out/include/libhighctidh

              mkdir -p $out/lib
              cp ${filename} $out/lib 

              $postInstall
            '';

            env = config.getBackendEnv backend;
          };
        in
        {
          packages."c-library/${variantName}" = package;
        };

      mkShellModule =
        tuple@{ stdenvVariant, backend }:
        {
          make-shells.${config.getVariantName tuple} = {
            stdenv = stdenvVariant.getPkg pkgs;
            env = config.getBackendEnv backend;
          };
        };
    in
    lib.pipe
      [
        (lib.pipe { inherit (variantFactors) stdenvVariant backend; } [
          lib.cartesianProduct
          (lib.filter (config.filterBackendBySystem system))
          (map mkShellModule)
        ])
        (lib.pipe variantFactors [
          lib.cartesianProduct
          (lib.filter (config.filterBackendBySystem system))
          (map mkPackageModule)
        ])
      ]
      [
        lib.flatten
        lib.mkMerge
      ];
}

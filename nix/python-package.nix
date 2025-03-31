{ lib, config, ... }:
{
  perSystem =
    { pkgs, system, ... }:
    let
      inherit (config) variantFactors;
      srcPath = ../src;

      mkShellImportsModule =
        stdenvVariant:
        let
          stdenv = stdenvVariant.getPkg pkgs;
          python = pkgs.python3.override { inherit stdenv; };
        in
        {
          make-shell.imports = [
            {
              packages = [ (python.withPackages (ps: [ ps.pytest ])) ];
            }
          ];
        };

      mkVariantModule =
        tuple@{ stdenvVariant, backend }:
        let
          stdenv = stdenvVariant.getPkg pkgs;
          python = pkgs.python3.override { inherit stdenv; };
          variantName = config.getVariantName tuple;
        in
        {
          packages."python/${variantName}" = python.pkgs.buildPythonPackage {
            name = "highctidh-${variantName}";
            inherit (config) version;

            src = lib.fileset.toSource {
              root = srcPath;
              fileset =
                with lib.fileset;
                unions [
                  (fileFilter (
                    file:
                    lib.any file.hasExt (
                      lib.optional (backend == "assembly") "S"
                      ++ [
                        "c"
                        "go"
                        "h"
                        "py"
                      ]
                    )
                  ) ../src)
                  ../src/VERSION
                  ../src/pyproject.toml
                  ../src/pytest.ini
                  ../src/setup.cfg
                ];
            };

            pyproject = true;

            nativeBuildInputs = [ python.pkgs.setuptools ];
            nativeCheckInputs = [ python.pkgs.pytestCheckHook ];

            env = config.getBackendEnv backend;
          };
        };

    in
    lib.pipe
      [
        (map mkShellImportsModule variantFactors.stdenvVariant)
        (lib.pipe { inherit (config.variantFactors) stdenvVariant backend; } [
          lib.cartesianProduct
          (lib.filter (config.filterBackendBySystem system))
          (map mkVariantModule)
        ])
      ]
      [ lib.flatten lib.mkMerge ];
}

from os import getcwd, path, uname, environ, umask
from subprocess import PIPE, Popen
from sys import exit
from sysconfig import get_config_var
from setuptools import Extension, setup, find_packages
from setuptools.command.build_ext import build_ext
from time import time

try:
    from stdeb.command.bdist_deb import bdist_deb
    from stdeb.command.sdist_dsc import sdist_dsc
except ImportError:
    bdist_deb = None
    sdist_dsc = None

class build_ext_helper(build_ext):
    # Note for builders who wish to use clang on GNU/Linux:
    #
    # Have you seen this error when trying to use clang?
    #
    #   ...
    #   running build_ext
    #   error: don't know how to compile C/C++ code on platform 'posix' with
    #   'clang' compiler
    #
    # Alternatively perhaps you have seen a linker error like the following?
    #
    #   ...
    #   clang  ...
    #   x86_64-linux-gnu-gcc -shared ...
    #   build/temp.linux-x86_64-3.10/some.o: file not recognized: file format
    #   not recognized
    #   collect2: error: ld returned 1 exit status
    #   error: command '/usr/bin/x86_64-linux-gnu-gcc' failed with exit code 1
    #   E: pybuild pybuild:369: build: plugin distutils failed with: exit
    #   code=1: /usr/bin/python3 setup.py build
    #
    # This helper class fixes an outstanding issue where setting CC=clang under
    # setuptools does not also update the linker and so clang builds the object
    # files for the extensions but then it fails to link as it attempts to use
    # gcc for that task. See pypi/setuptools #1442 for more information.  When
    # used with debian/rules for building, corresponding PYBUILD_* options must
    # be set to ensure everything works as intended.  Please consult
    # misc/debian-rules for an example.
    #
    def build_extensions(self):
        print(f"Compiler was: {self.compiler.linker_exe}")
        print(f"Linker was: {self.compiler.linker_so}")
        # NOTE:
        # This entire class is to work around a pernicous and annoying bug that
        # previously prevented using any compiler other than gcc on GNU/Linux
        # platforms for certain kinds of builds.  By setting CC=clang or
        # CC=gcc, builds will be compiled by the selected compiler as expected.
        # However, self.compiler.linker_exe is mistakenly not updated by
        # setting the CC environment variable.  To work around this bug which
        # only impacts users of an alternative compiler, we hot patch only the
        # linker executable name:
        self.compiler.linker_so[0] = self.compiler.linker_exe[0]
        print(f"Compiler is now: {self.compiler.linker_exe}")
        print(f"Linker is now: {self.compiler.linker_so}")
        build_ext.build_extensions(self)
    def run(self):
        build_ext.run(self)


requirements = []
dir_include = [".", path.join(getcwd()), ]
lib_include = [getcwd(), ]
if "SOURCE_DATE_EPOCH" in environ:
    sda = str(int(environ["SOURCE_DATE_EPOCH"]))
    print("SOURCE_DATE_EPOCH is set:")
    print(f"SOURCE_DATE_EPOCH={sda}")
else:
    print("SOURCE_DATE_EPOCH is unset, setting to today")
    environ['SOURCE_DATE_EPOCH'] = str(int(time()))
    sda = str(int(environ["SOURCE_DATE_EPOCH"]))
    print(f"SOURCE_DATE_EPOCH={sda}")
# Set umask to ensure consistent file permissions inside build artifacts such
# as `.whl` files
umask(0o022)

CC = None
if "CC" in environ:
    CC = str(environ["CC"])
    print(f"CC={CC}")

VERSION = open('VERSION', 'r').read().strip()
PLATFORM = uname().machine
r = Popen("getconf" + " LONG_BIT", shell=True, stdout=PIPE)
PLATFORM_SIZE = int(r.stdout.read().strip())
if PLATFORM_SIZE != 64:
    if PLATFORM_SIZE != 32:
        print(f"PLATFORM is: {PLATFORM}")
        print(f"PLATFORM_SIZE is unexpected size: {PLATFORM_SIZE}")
        exit(2)

base_src = ["crypto_classify.c", "crypto_declassify.c", "csidh.c",
            "elligator.c", "fp2fiat.c", "mont.c", "poly.c", "randombytes.c",
            "random.c", "skgen.c", "steps.c", "steps_untuned.c", "umults.c",
            "validate.c"]

cflags = get_config_var("CFLAGS").split()
cflags += ["-Wextra", "-Wall", "-Wall", "-Wextra", "-fpie", "-fPIC",
           "-fwrapv", "-pedantic", "-O3", "-Os",]
cflags += ["-DGETRANDOM", f"-DPLATFORM={PLATFORM}",
           f"-DPLATFORM_SIZE={PLATFORM_SIZE}"]

if CC == "clang":
    cflags += ["-Wno-ignored-optimization-argument"]

match PLATFORM:
    case "aarch64":
        base_src += ["int32_sort.c",]
        cflags += ["-march=native", "-mtune=native"]
    case "armv7l":
        # clang required
        base_src += ["int32_sort.c",]
        cflags += ["-fforce-enable-int128", "-D__ARM32__"]
    case "loongarch64":
        base_src += ["int32_sort.c",]
        cflags += ["-march=native", "-mtune=native"]
    case "mips64":
        # clang required
        base_src += ["int32_sort.c",]
        cflags += ["-fforce-enable-int128"]
    case "ppc64le":
        base_src += ["int32_sort.c",]
        cflags += ["-mtune=native"]
    case "ppc64":
        base_src += ["int32_sort.c",]
        cflags += ["-mtune=native"]
    case "riscv64":
        base_src += ["int32_sort.c",]
    case "s390x":
        base_src += ["int32_sort.c",]
        cflags += ["-march=native", "-mtune=native"]
    case "sparc64":
        base_src += ["int32_sort.c",]
        cflags += ["-march=native", "-mtune=native"]
    case "x86_64":
        if PLATFORM_SIZE == 64:
            base_src += ["int32_sort.c",]
            cflags += ["-march=native", "-mtune=native", "-D__x86_64__"]
        elif PLATFORM_SIZE == 32:
            # clang required
            base_src += ["int32_sort.c",]
            cflags += ["-fforce-enable-int128", "-D__i386__"]
    case _:
        base_src += ["int32_sort.c",]

ldflags = ["-s", "-w"]

if __name__ == "__main__":
    setup(
        name = "highctidh",
        version = VERSION,
        author = "Jacob Appelbaum",
        zip_safe = False,
        author_email = "jacob@appelbaum.net",
        packages = ['highctidh'],
        install_requires = [],
        cmdclass = dict(bdist_deb=bdist_deb, sdist_dsc=sdist_dsc),
        ext_modules = [
            Extension("highctidh_511",
                extra_compile_args = cflags + ["-DBITS=511",
                    "-DNAMESPACEBITS(x)=highctidh_511_##x",
                    "-DNAMESPACEGENERIC(x)=highctidh_##x"],
                extra_link_args = ldflags,
                include_dirs = dir_include,
                language = 'c',
                library_dirs = lib_include,
                sources = base_src + [ "fiat_p511.c", "fp_inv511.c",
                    "fp_sqrt511.c", "primes511.c", ],
            ),
            Extension("highctidh_512",
                extra_compile_args = cflags + ["-DBITS=512",
                    "-DNAMESPACEBITS(x)=highctidh_512_##x",
                    "-DNAMESPACEGENERIC(x)=highctidh_##x"],
                extra_link_args = ldflags,
                include_dirs = dir_include,
                language = 'c',
                library_dirs = lib_include,
                sources = base_src + [ "fiat_p512.c", "fp_inv512.c",
                    "fp_sqrt512.c", "primes512.c", ],
            ),
            Extension("highctidh_1024",
                extra_compile_args = cflags + ["-DBITS=1024",
                    "-DNAMESPACEBITS(x)=highctidh_1024_##x",
                    "-DNAMESPACEGENERIC(x)=highctidh_##x"],
                extra_link_args = ldflags,
                include_dirs = dir_include,
                language = 'c',
                library_dirs = lib_include,
                sources = base_src + [ "fiat_p1024.c", "fp_inv1024.c",
                    "fp_sqrt1024.c", "primes1024.c", ],
            ),
            Extension("highctidh_2048",
                extra_compile_args = cflags + ["-DBITS=2048",
                    "-DNAMESPACEBITS(x)=highctidh_2048_##x",
                    "-DNAMESPACEGENERIC(x)=highctidh_##x"],
                extra_link_args = ldflags,
                include_dirs = dir_include,
                language ='c',
                library_dirs = lib_include,
                sources = base_src + [ "fiat_p2048.c", "fp_inv2048.c",
                    "fp_sqrt2048.c", "primes2048.c", ],
            ),
        ]

    )

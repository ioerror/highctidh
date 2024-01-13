from os import getcwd, path, uname, environ, umask, stat, mkdir
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
if "LLVM_PARALLEL_LINK_JOBS" in environ:
    sdb = str(int(environ["LLVM_PARALLEL_LINK_JOBS"]))
    print(f"LLVM_PARALLEL_LINK_JOBS is set: {sdb}")
else:
    print("LLVM_PARALLEL_LINK_JOBS is unset, setting to 1")
    environ['LLVM_PARALLEL_LINK_JOBS'] = str(int(1))
    sdb = str(int(environ["LLVM_PARALLEL_LINK_JOBS"]))
    print(f"LLVM_PARALLEL_LINK_JOBS={sdb}")
# Set umask to ensure consistent file permissions inside build artifacts such
# as `.whl` files
umask(0o022)

try:
  stat("build")
except FileNotFoundError:
  try:
      mkdir("build")
  except FileExistsError:
      pass

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
            "validate.c", "int32_sort.c"]

cflags = get_config_var("CFLAGS").split()
cflags += ["-Wextra", "-Wall", "-fpie", "-fPIC", "-fwrapv", "-pedantic", "-O2",
           "-g0", "-fno-lto"]
cflags += ["-DGETRANDOM", f"-DPLATFORM={PLATFORM}",
           f"-DPLATFORM_SIZE={PLATFORM_SIZE}"]
cflags += ["-Wformat", "-Werror=format-security", "-D_FORTIFY_SOURCE=2",
           "-fstack-protector-strong"]
ldflags = ["-s", "-w"]

if CC == "clang":
    cflags += ["-Wno-ignored-optimization-argument", "-Wno-unreachable-code"]
if CC == "gcc":
    cflags += ["-Werror"]
    ldflags += ["-Wl,-Bsymbolic-functions", "-Wl,-z,noexecstack",
               "-Wl,-z,relro", "-Wl,-z,now", "-Wl,--reduce-memory-overheads",
               "-Wl,--no-keep-memory",]

print(f"Building for platform: {PLATFORM}")
if PLATFORM == "aarch64":
  if CC == "clang":
      cflags += ["-DHIGHCTIDH_PORTABLE"]
  if CC == "gcc":
      cflags += ["-march=native", "-mtune=native", "-DHIGHCTIDH_PORTABLE"]
elif PLATFORM == "armv7l":
  # clang required
  if CC == "clang":
      cflags += ["-fforce-enable-int128", "-D__ARM32__",
                 "-DHIGHCTIDH_PORTABLE",]
  if CC == "gcc":
      cflags += ["-D__ARM32__", "-DHIGHCTIDH_PORTABLE"]
elif PLATFORM == "loongarch64":
  cflags += ["-march=native", "-mtune=native", "-DHIGHCTIDH_PORTABLE"]
elif PLATFORM == "mips":
  # clang or mips64-linux-gnuabi64-gcc cross compile required
  if CC == "clang":
      cflags += ["-fforce-enable-int128", "-DHIGHCTIDH_PORTABLE"]
  if CC == "gcc":
      cflags += ["-DHIGHCTIDH_PORTABLE"]
elif PLATFORM == "mips64":
  cflags += ["-DHIGHCTIDH_PORTABLE"]
elif PLATFORM == "mips64le":
  cflags += ["-DHIGHCTIDH_PORTABLE"]
elif PLATFORM == "ppc64le":
  cflags += ["-mtune=native", "-DHIGHCTIDH_PORTABLE"]
elif PLATFORM == "ppc64":
  cflags += ["-mtune=native", "-DHIGHCTIDH_PORTABLE"]
elif PLATFORM == "riscv64":
  cflags += ["-DHIGHCTIDH_PORTABLE"]
elif PLATFORM == "s390x":
  if CC == "clang":
      cflags += ["-march=z10", "-mtune=z10", "-DHIGHCTIDH_PORTABLE"]
  if CC == "gcc":
      cflags += ["-march=z10", "-mtune=z10", "-DHIGHCTIDH_PORTABLE"]
elif PLATFORM == "sparc64":
  cflags += ["-march=native", "-mtune=native", "-DHIGHCTIDH_PORTABLE"]
elif PLATFORM == "x86_64":
  if PLATFORM_SIZE == 64:
    cflags += ["-march=native", "-mtune=native", "-D__x86_64__"]
  elif PLATFORM_SIZE == 32:
    # clang required
    cflags += ["-fforce-enable-int128", "-D__i386__", "-DHIGHCTIDH_PORTABLE"]
else:
  cflags += ["-DHIGHCTIDH_PORTABLE"]

# We default to fiat as the backend for all platforms except x86_64
if PLATFORM == "x86_64" and PLATFORM_SIZE == 64:
    src_511 =  base_src + ["fp_inv511.c", "fp_sqrt511.c", "primes511.c",]
    src_512 =  base_src + ["fp_inv512.c", "fp_sqrt512.c", "primes512.c",]
    src_1024 = base_src + ["fp_inv1024.c", "fp_sqrt1024.c", "primes1024.c",]
    src_2048 = base_src + ["fp_inv2048.c", "fp_sqrt2048.c", "primes2048.c",]
else:
    src_511 = base_src + ["fiat_p511.c", "fp_inv511.c", "fp_sqrt511.c", "primes511.c",]
    src_512 = base_src + ["fiat_p512.c", "fp_inv512.c", "fp_sqrt512.c", "primes512.c",]
    src_1024 = base_src + ["fiat_p1024.c", "fp_inv1024.c", "fp_sqrt1024.c", "primes1024.c",]
    src_2048 = base_src + ["fiat_p2048.c", "fp_inv2048.c", "fp_sqrt2048.c", "primes2048.c",]

extra_compile_args_511 = cflags + ["-DBITS=511",
        "-DNAMESPACEBITS(x)=highctidh_511_##x",
        "-DNAMESPACEGENERIC(x)=highctidh_##x"]
extra_compile_args_512 = cflags + ["-DBITS=512",
        "-DNAMESPACEBITS(x)=highctidh_512_##x",
        "-DNAMESPACEGENERIC(x)=highctidh_##x"]
extra_compile_args_1024 = cflags + ["-DBITS=1024",
        "-DNAMESPACEBITS(x)=highctidh_1024_##x",
        "-DNAMESPACEGENERIC(x)=highctidh_##x"]
extra_compile_args_2048 = cflags + ["-DBITS=2048",
        "-DNAMESPACEBITS(x)=highctidh_2048_##x",
        "-DNAMESPACEGENERIC(x)=highctidh_##x"]
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
                extra_compile_args = extra_compile_args_511,
                extra_link_args = ldflags,
                include_dirs = dir_include,
                language = 'c',
                library_dirs = lib_include,
                sources = src_511,
            ),
            Extension("highctidh_512",
                extra_compile_args = extra_compile_args_512,
                extra_link_args = ldflags,
                include_dirs = dir_include,
                language = 'c',
                library_dirs = lib_include,
                sources = src_512,
            ),
            Extension("highctidh_1024",
                extra_compile_args = extra_compile_args_1024,
                extra_link_args = ldflags,
                include_dirs = dir_include,
                language = 'c',
                library_dirs = lib_include,
                sources = src_1024,
            ),
            Extension("highctidh_2048",
                extra_compile_args = extra_compile_args_2048,
                extra_link_args = ldflags,
                include_dirs = dir_include,
                language ='c',
                library_dirs = lib_include,
                sources = src_2048,
            ),
        ]

    )

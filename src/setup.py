from warnings import filterwarnings

filterwarnings("ignore", category=DeprecationWarning)

from os import environ, getcwd, mkdir, path, stat, umask
from subprocess import PIPE, Popen
from sys import exit, platform
from sysconfig import get_config_var
from time import time

from platform import architecture, uname
from sys import byteorder
import timeit

from setuptools import Extension, find_packages, setup
from setuptools.command.build_ext import build_ext

OS = uname().system
PLATFORM = uname().machine
PLATFORM_SIZE = int(architecture()[0][0:2])

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
        print(self.compiler)
        if OS != 'Windows':
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
dir_include = [
    ".",
    path.join(getcwd()),
]
lib_include = [
    getcwd(),
]
if "SOURCE_DATE_EPOCH" in environ:
    sda = str(int(environ["SOURCE_DATE_EPOCH"]))
    print("SOURCE_DATE_EPOCH is set:")
    print(f"SOURCE_DATE_EPOCH={sda}")
else:
    print("SOURCE_DATE_EPOCH is unset, setting to today")
    environ["SOURCE_DATE_EPOCH"] = str(int(time()))
    sda = str(int(environ["SOURCE_DATE_EPOCH"]))
    print(f"SOURCE_DATE_EPOCH={sda}")
if "HIGHCTIDH_PORTABLE" in environ:
    if environ.get("HIGHCTIDH_PORTABLE") == "0" :
        HIGHCTIDH_PORTABLE = str(0)
    else:
        HIGHCTIDH_PORTABLE = str(1)
    print(f"HIGHCTIDH_PORTABLE is set: {HIGHCTIDH_PORTABLE}")
else:
    HIGHCTIDH_PORTABLE = str(1)
    print(f"HIGHCTIDH_PORTABLE was unset, enabling by default; now set to 1")
if "LLVM_PARALLEL_LINK_JOBS" in environ:
    sdb = str(int(environ["LLVM_PARALLEL_LINK_JOBS"]))
    print(f"LLVM_PARALLEL_LINK_JOBS is set: {sdb}")
else:
    print("LLVM_PARALLEL_LINK_JOBS is unset, setting to 1")
    environ["LLVM_PARALLEL_LINK_JOBS"] = str(int(1))
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

try:
    VERSION = open("src/VERSION", "r").read().strip()
except FileNotFoundError:
    VERSION = open("VERSION", "r").read().strip()
except AttributeError:
    VERSION = "3.141592653"

base_src = [
    "crypto_classify.c",
    "crypto_declassify.c",
    "csidh.c",
    "elligator.c",
    "fp2fiat.c",
    "mont.c",
    "poly.c",
    "randombytes.c",
    "random.c",
    "skgen.c",
    "steps.c",
    "steps_untuned.c",
    "umults.c",
    "validate.c",
    "int32_sort.c",
]

cflags = get_config_var("CFLAGS")
if cflags is not None and cflags is str:
  cflags = cflags.split()
else:
  cflags = ["-Wextra"]
cflags += ["-Wall", "-fPIC", "-fpie", "-fwrapv", "-pedantic", "-O2", "-g0", "-fno-lto"]
cflags += ["-DGETRANDOM", f"-DPLATFORM={PLATFORM}", f"-DPLATFORM_SIZE={PLATFORM_SIZE}"]
cflags += [
    "-Wformat",
    "-Werror=format-security",
    "-D_FORTIFY_SOURCE=2",
    "-fstack-protector-strong",
]
ldflags = ["-s"]

if CC == "clang":
    cflags += ["-Wno-ignored-optimization-argument", "-Wno-unreachable-code"]
if CC == "gcc":
    if OS == "Linux":
        cflags += ["-Werror"]
        ldflags += [
            "-Wl,-Bsymbolic-functions",
            "-Wl,-z,noexecstack",
            "-Wl,-z,relro",
            "-Wl,-z,now",
            "-Wl,--reduce-memory-overheads",
            "-Wl,--no-keep-memory",
        ]

print(f"Building for platform: {PLATFORM} on {OS}")
if PLATFORM == "aarch64" or PLATFORM == "arm64":
    cflags += ["-D__ARM64__"]
    if OS == "Darwin":
      cflags += ["-D__Darwin__"]
      cflags += ["-DHIGHCTIDH_PORTABLE=1"]
    else:
      if CC == "clang":
          cflags += ["-DHIGHCTIDH_PORTABLE=1"]
      if CC == "gcc":
          cflags += ["-march=native", "-mtune=native", "-DHIGHCTIDH_PORTABLE=1"]
elif PLATFORM == "armv7l":
    # clang required
    if CC == "clang":
        cflags += [
            "-fforce-enable-int128",
            "-D__ARM32__",
            "-DHIGHCTIDH_PORTABLE=1",
        ]
    if CC == "gcc":
        cflags += ["-D__ARM32__", "-DHIGHCTIDH_PORTABLE=1"]
elif PLATFORM == "loongarch64":
    cflags += ["-march=native", "-mtune=native", "-DHIGHCTIDH_PORTABLE=1"]
elif PLATFORM == "mips":
    # clang or mips64-linux-gnuabi64-gcc cross compile required
    if CC == "clang":
        cflags += ["-fforce-enable-int128", "-DHIGHCTIDH_PORTABLE=1"]
    if CC == "gcc":
        cflags += ["-DHIGHCTIDH_PORTABLE=1"]
elif PLATFORM == "mips64":
    cflags += ["-DHIGHCTIDH_PORTABLE=1"]
elif PLATFORM == "mips64le":
    cflags += ["-DHIGHCTIDH_PORTABLE=1"]
elif PLATFORM == "ppc64le":
    cflags += ["-mtune=native", "-DHIGHCTIDH_PORTABLE=1"]
elif PLATFORM == "ppc64":
    cflags += ["-mtune=native", "-DHIGHCTIDH_PORTABLE=1"]
elif PLATFORM == "riscv64":
    cflags += ["-D__riscv", "-DHIGHCTIDH_PORTABLE=1"]
elif PLATFORM == "s390x":
    if CC == "clang":
        cflags += ["-march=z10", "-mtune=z10", "-DHIGHCTIDH_PORTABLE=1"]
    if CC == "gcc":
        cflags += ["-march=z10", "-mtune=z10", "-DHIGHCTIDH_PORTABLE=1"]
elif PLATFORM == "sun4v" or PLATFORM == "i86pc":
    # Solaris 11, SunOS has default flags that do not work for both gcc and clang
    # compilers. We wrap the function that returns these flags internally during
    # the build process to override them for a value that works for both.
    import distutils.sysconfig as _wrapped_distutils
    from distutils.sysconfig import get_config_vars as _get_config_vars

    _config_vars = _get_config_vars().copy()
    default_cflags = [
        " -O2 -DNDEBUG -Wall -m64 -fPIC -fpie -DPIC -ffile-prefix-map=..=. "
    ]
    _config_vars["CFLAGS"] = default_cflags

    def get_config_vars_wrapper(*a):
        return [_config_vars.get(n) for n in a] if a else _config_vars

    _wrapped_distutils.get_config_vars = get_config_vars_wrapper
    # Set Solaris specific build options
    cflags = ["-D__sun"]
    ldflags = ["-s"]
    ldflags += ["-Wl,-Bsymbolic-functions"]
    if PLATFORM == "i86pc":
        cflags += ["-D__i86pc__"]
        cflags += ["-DHIGHCTIDH_PORTABLE=" + HIGHCTIDH_PORTABLE]
        if CC == "gcc":
            cflags += ["-mcpu=native", "-mtune=native", "-fno-lto"]
    if PLATFORM == "sun4v":
        cflags += ["-D__sun4v__"]
        cflags += ["-DHIGHCTIDH_PORTABLE=1"]
        if CC == "clang":
            cflags += ["-fforce-enable-int128"]
    cflags += ["-Wextra", "-fwrapv", "-pedantic", "-Werror", "-DGETRANDOM"]
    cflags += [f"-DPLATFORM={PLATFORM}", f"-DPLATFORM_SIZE={PLATFORM_SIZE}"]
elif PLATFORM == "x86_64" and OS == 'Windows' \
     or PLATFORM == "AMD64" and OS == 'Windows' \
     or PLATFORM == "x86_64" and OS.startswith('MINGW') \
     or PLATFORM == "x86_64" and OS.startswith('MSYS'):
    # Windows only builds with clang on Windows under the CI
    # It should also build with other compilers.
    # As with Solaris we wrap the function that returns these flags internally
    # during the build process to override them for a value that works for
    # both.
    import distutils.sysconfig as _wrapped_distutils
    from distutils.sysconfig import get_config_vars as _get_config_vars

    _config_vars = _get_config_vars().copy()
    default_cflags = [
        " -Wall -pedantic -O2 -fwrapv -ffile-prefix-map=..=."
    ]
    _config_vars["CFLAGS"] = default_cflags

    def get_config_vars_wrapper(*a):
        return [_config_vars.get(n) for n in a] if a else _config_vars

    _wrapped_distutils.get_config_vars = get_config_vars_wrapper
    # Set Windows specific build options
    cflags = ["-D__Windows__"]
    ldflags = ["-LAdvapi32.lib"]
    if OS.startswith('MINGW64') or OS.startswith('MSYS'):
        cflags += ["-D_WIN64"]
    if OS.startswith('MINGW32'):
        cflags += ["-D_WIN32"]
    if PLATFORM == "AMD64":
        cflags += ["-D__x86_64__"]
        cflags += ["-D__AMD64__"]
        cflags += ["-DHIGHCTIDH_PORTABLE=" + HIGHCTIDH_PORTABLE]
    if PLATFORM == "x86_64":
        cflags += ["-D__x86_64__"]
        cflags += ["-DHIGHCTIDH_PORTABLE=" + HIGHCTIDH_PORTABLE]
    cflags += [f"-DPLATFORM={PLATFORM}", f"-DPLATFORM_SIZE={PLATFORM_SIZE}"]
elif PLATFORM == "x86_64":
    if PLATFORM_SIZE == 64:
        cflags += ["-D__x86_64__"]
        if OS == "Darwin":
          cflags += ["-D__Darwin__"]
          cflags += ["-DHIGHCTIDH_PORTABLE=1"]
        else:
          if CC == "clang":
              cflags += ["-DHIGHCTIDH_PORTABLE=1"]
          if CC == "gcc":
              cflags += ["-march=native", "-mtune=native"]
        if HIGHCTIDH_PORTABLE == "1":
            if CC == "clang":
                cflags += ["-fforce-enable-int128"]
    elif PLATFORM_SIZE == 32:
        # clang required
        cflags += ["-fforce-enable-int128", "-D__i386__"]
        HIGHCTIDH_PORTABLE = "1"
    cflags += ["-DHIGHCTIDH_PORTABLE=" + HIGHCTIDH_PORTABLE]
else:
    cflags += ["-DHIGHCTIDH_PORTABLE=" + HIGHCTIDH_PORTABLE]

# We default to fiat as the backend for all platforms except x86_64/i86pc
if (
    HIGHCTIDH_PORTABLE == "0"
    and PLATFORM_SIZE == 64
    and (PLATFORM == "x86_64"
    or PLATFORM == "i86pc")
):
    print("Selecting x86_64 asm backend")
    src_511 = base_src + [
        "fp_inv511.c",
        "fp_sqrt511.c",
        "primes511.c",
    ]
    src_512 = base_src + [
        "fp_inv512.c",
        "fp_sqrt512.c",
        "primes512.c",
    ]
    src_1024 = base_src + [
        "fp_inv1024.c",
        "fp_sqrt1024.c",
        "primes1024.c",
    ]
    src_2048 = base_src + [
        "fp_inv2048.c",
        "fp_sqrt2048.c",
        "primes2048.c",
    ]
else:
    print("Selecting portable fiat backend")
    src_511 = base_src + [
        "fiat_p511.c",
        "fp_inv511.c",
        "fp_sqrt511.c",
        "primes511.c",
    ]
    src_512 = base_src + [
        "fiat_p512.c",
        "fp_inv512.c",
        "fp_sqrt512.c",
        "primes512.c",
    ]
    src_1024 = base_src + [
        "fiat_p1024.c",
        "fp_inv1024.c",
        "fp_sqrt1024.c",
        "primes1024.c",
    ]
    src_2048 = base_src + [
        "fiat_p2048.c",
        "fp_inv2048.c",
        "fp_sqrt2048.c",
        "primes2048.c",
    ]

extra_compile_args = cflags
extra_compile_args_511 = cflags + [
    "-DBITS=511",
    "-DNAMESPACEBITS(x)=highctidh_511_##x",
    "-DNAMESPACEGENERIC(x)=highctidh_##x",
]
extra_compile_args_512 = cflags + [
    "-DBITS=512",
    "-DNAMESPACEBITS(x)=highctidh_512_##x",
    "-DNAMESPACEGENERIC(x)=highctidh_##x",
]
extra_compile_args_1024 = cflags + [
    "-DBITS=1024",
    "-DNAMESPACEBITS(x)=highctidh_1024_##x",
    "-DNAMESPACEGENERIC(x)=highctidh_##x",
]
extra_compile_args_2048 = cflags + [
    "-DBITS=2048",
    "-DNAMESPACEBITS(x)=highctidh_2048_##x",
    "-DNAMESPACEGENERIC(x)=highctidh_##x",
]
print(f"511 files: {src_511=}")
print(f"512 files: {src_512=}")
print(f"1024 files: {src_1024=}")
print(f"2048 files: {src_2048=}")
print(f"511 cflags: {extra_compile_args_511=}")
print(f"512 cflags: {extra_compile_args_512=}")
print(f"1024 cflags: {extra_compile_args_1024=}")
print(f"2048 cflags: {extra_compile_args_2048=}")
if __name__ == "__main__":
    setup(
        name="highctidh",
        version=VERSION,
        author="Jacob Appelbaum",
        zip_safe=False,
        author_email="jacob@appelbaum.net",
        packages=["highctidh"],
        install_requires=[],
        cmdclass=dict(
            bdist_deb=bdist_deb, sdist_dsc=sdist_dsc, build_ext=build_ext_helper
        ),
        ext_modules=[
            Extension(
                "highctidh_511",
                extra_compile_args=extra_compile_args_511,
                extra_link_args=ldflags,
                include_dirs=dir_include,
                language="c",
                library_dirs=lib_include,
                sources=src_511,
            ),
            Extension(
                "highctidh_512",
                extra_compile_args=extra_compile_args_512,
                extra_link_args=ldflags,
                include_dirs=dir_include,
                language="c",
                library_dirs=lib_include,
                sources=src_512,
            ),
            Extension(
                "highctidh_1024",
                extra_compile_args=extra_compile_args_1024,
                extra_link_args=ldflags,
                include_dirs=dir_include,
                language="c",
                library_dirs=lib_include,
                sources=src_1024,
            ),
            Extension(
                "highctidh_2048",
                extra_compile_args=extra_compile_args_2048,
                extra_link_args=ldflags,
                include_dirs=dir_include,
                language="c",
                library_dirs=lib_include,
                sources=src_2048,
            ),
        ],
    )

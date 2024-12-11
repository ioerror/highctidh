import os, sys, subprocess
from time import time
from setuptools import setup, Extension
from setuptools.command.build_ext import build_ext

try:
    from stdeb.command.bdist_deb import bdist_deb
    from stdeb.command.sdist_dsc import sdist_dsc
except ImportError:
    bdist_deb = None
    sdist_dsc = None

if "SOURCE_DATE_EPOCH" in os.environ:
    print("SOURCE_DATE_EPOCH is set:")
else:
    print("SOURCE_DATE_EPOCH is unset, setting to today")
    os.environ["SOURCE_DATE_EPOCH"] = str(int(time()))
print("SOURCE_DATE_EPOCH=" + os.environ["SOURCE_DATE_EPOCH"])

# Set umask to ensure consistent file permissions inside build artifacts such
# as `.whl` files
os.umask(0o022)

try:
    VERSION = open("VERSION", "r").read().strip()
except FileNotFoundError:
    VERSION = "3.141592653"


class build_ext_helper(build_ext):
    def run(self):
        if os.uname().machine == "s390x":
            os.environ['CFLAGS'] = "-fno-lto -march=z10 -mtune=z10"
        else:
            os.environ['CFLAGS'] = "-fno-lto"
        subprocess.run(shell=True, check=True, args='make -C src -f GNUmakefile')

        if sys.platform in ['win32', 'cygwin']:
            shext = 'dll'
        elif sys.platform == 'darwin':
            shext = 'dylib'
        else:
            shext = 'so'

        for ext in self.extensions:
            lib = os.path.join('src', 'lib' + ext.name + '.' + shext)
            out = self.get_ext_fullpath(ext.name)
            self.copy_file(lib, out, preserve_mode=False)


if __name__ == "__main__":
    setup(
        name="highctidh",
        version=VERSION,
        author="Jacob Appelbaum",
        author_email="jacob@appelbaum.net",
        url="https://codeberg.org/vula/highctidh",
        zip_safe=False,
        cmdclass=dict(
            bdist_deb=bdist_deb, sdist_dsc=sdist_dsc, build_ext=build_ext_helper
        ),
        packages=['highctidh'],
        package_dir={'': 'src'},
        ext_modules=[
            Extension("highctidh_511", []),
            Extension("highctidh_512", []),
            Extension("highctidh_1024", []),
            Extension("highctidh_2048", [])
        ]
    )

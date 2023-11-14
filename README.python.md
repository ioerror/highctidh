This Python module requires Python 3.9.2 or later.

Run the unit tests in the root of the project:

    python3 -m unittest -v

Alternatively with pytest in the root of the project:

    pytest-3 -v

Build and install the Python module:

    python3 -m build
    pip install dist/highctidh-*-py3-none-any.whl

Alternatively, rather than using `python3 -m build` pip may be used:

    pip install . # or sudo pip install .

Use the Python module in a venv:

    $ cd /tmp
    $ python3 -m venv ctest
    $ source /tmp/ctest/bin/activate
    $ pip install /home/user/c/highctidh
    $ python3 -c 'import highctidh'

Use the Python module:

    $ python3
    Python 3.10.4 (main, Jun 29 2022, 12:14:53) [GCC 11.2.0] on linux
    Type "help", "copyright", "credits" or "license" for more information.
    >>> import highctidh
    >>> ctidh512 = highctidh.ctidh(512) # options are 511, 512, 1024, 2048

Detailed Python module usage and test vectors are available in
`tests/test_highctidh.py` with examples for all four field sizes.

To build a Debian package that includes the Python module and the relevant .so
files for internal Python use, run:

    python3 setup.py bdist_deb

Alternatively use `make deb`

To build an `x86_64` wheel:

    python3 setup.py bdist_wheel --plat-name x86_64

Build, check, and finally upload to PyPi:

    export SOURCE_DATE_EPOCH=`git log -1 --pretty=%ct`
    export VERSION=`cat VERSION`
    python3 -m build
    python3 -m twine check dist/*$VERSION*
    python3 -m twine upload --repository pypi dist/*$VERSION*

Build for every supported architecture by first preparing a Docker environment
with Debian and QEMU:

    ./docker-setup.sh

Build packages for each supported architecture in Docker:

    ./docker-multi-arch-package-build.sh

Future releases may be installable from pypi.

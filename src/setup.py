from setuptools import setup

try:
    from stdeb.command.bdist_deb import bdist_deb
    from stdeb.command.sdist_dsc import sdist_dsc
except ImportError:
    bdist_deb = None
    sdist_dsc = None

if __name__ == "__main__":
    setup(
        cmdclass=dict(
            bdist_deb=bdist_deb, sdist_dsc=sdist_dsc
        ),
    )
#        libraries=["highctidh_511"],

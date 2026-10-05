#!/usr/bin/env python3

import os
import shutil
import subprocess
import unittest

SRC = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


def dry_run(*args):
    return subprocess.run(
        ["make", "-n", "-f", "GNUmakefile", *args, "fp2fiat511.o"],
        cwd=SRC,
        env=dict(os.environ, MAKEFLAGS=""),
        check=True,
        capture_output=True,
        text=True,
    ).stdout.split()


@unittest.skipIf(shutil.which("make") is None, "make not found")
@unittest.skipUnless(
    all(os.path.exists(os.path.join(SRC, n)) for n in ("GNUmakefile", "autogen")),
    "GNUmakefile or autogen not in the source tree",
)
class TestGNUmakefileMarch(unittest.TestCase):
    def test_no_march_default(self):
        for name in ("GNUmakefile", "autogen"):
            with open(os.path.join(SRC, name)) as f:
                self.assertFalse("CC_MARCH ?= native" in f.read(), name)
        bad = [f for f in dry_run() if f.startswith(("-march", "-mtune"))]
        self.assertEqual(bad, [])

    def test_cc_march_opt_in(self):
        flags = dry_run("CC_MARCH=x86-64-v2")
        self.assertIn("-march=x86-64-v2", flags)
        self.assertIn("-mtune=x86-64-v2", flags)


if __name__ == "__main__":
    unittest.main()

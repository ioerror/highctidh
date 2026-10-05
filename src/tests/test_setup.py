#!/usr/bin/env python3

import json
import os
import subprocess
import sys
import unittest

SRC = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

PROBE = """
import contextlib, io, json, platform, runpy
real = platform.uname()
platform.uname = lambda: real._replace(system="Linux", machine="aarch64")
platform.architecture = lambda *a, **k: ("64bit", "ELF")
with contextlib.redirect_stdout(io.StringIO()):
    g = runpy.run_path("setup.py", run_name="probe")
print(json.dumps({b: [g["src_" + b], g["extra_compile_args_" + b]]
                  for b in ("511", "512", "1024", "2048")}))
"""


def aarch64_linux_build(cc):
    env = dict(os.environ)
    env.pop("CC", None)
    env.pop("HIGHCTIDH_PORTABLE", None)
    if cc is not None:
        env["CC"] = cc
    out = subprocess.run(
        [sys.executable, "-c", PROBE],
        cwd=SRC,
        env=env,
        check=True,
        capture_output=True,
        text=True,
    ).stdout
    return json.loads(out.strip().splitlines()[-1])


class TestSetupAarch64(unittest.TestCase):
    def test_fiat_sources_get_portable_define(self):
        for cc in (None, "", "cc", "gcc", "clang", "aarch64-linux-gnu-gcc"):
            build = aarch64_linux_build(cc)
            for bits, (sources, cflags) in build.items():
                with self.subTest(cc=cc, bits=bits):
                    self.assertIn("fiat_p" + bits + ".c", sources)
                    self.assertIn("-DHIGHCTIDH_PORTABLE=1", cflags)


if __name__ == "__main__":
    unittest.main()

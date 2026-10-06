#!/usr/bin/env python3

import json
import os
import re
import shutil
import subprocess
import sys
import tempfile
import unittest

SRC = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

PROBE = """
import contextlib, io, json, platform, runpy, sys
system, machine, size = sys.argv[1:4]
real = platform.uname()
platform.uname = lambda: real._replace(system=system, machine=machine)
platform.architecture = lambda *a, **k: (size + "bit", "ELF")
with contextlib.redirect_stdout(io.StringIO()):
    g = runpy.run_path("setup.py", run_name="probe")
print(json.dumps({b: [g["src_" + b], g["extra_compile_args_" + b]]
                  for b in ("511", "512", "1024", "2048")}))
"""

PLATFORMS = (
    ("Linux", "aarch64", "64", "aarch64-linux-gnu-gcc"),
    ("Linux", "arm64", "64", "aarch64-linux-gnu-gcc"),
    ("Darwin", "arm64", "64", "aarch64-apple-darwin-gcc"),
    ("Linux", "armv7l", "32", "arm-linux-gnueabihf-gcc"),
    ("Linux", "loongarch64", "64", "loongarch64-linux-gnu-gcc"),
    ("Linux", "mips", "32", "mips-linux-gnu-gcc"),
    ("Linux", "mips64", "64", "mips64-linux-gnuabi64-gcc"),
    ("Linux", "mips64le", "64", "mips64el-linux-gnuabi64-gcc"),
    ("Linux", "ppc64le", "64", "powerpc64le-linux-gnu-gcc"),
    ("Linux", "ppc64", "64", "powerpc64-linux-gnu-gcc"),
    ("Linux", "riscv64", "64", "riscv64-linux-gnu-gcc"),
    ("Linux", "s390x", "64", "s390x-linux-gnu-gcc"),
    ("SunOS", "sun4v", "64", "sparc64-sun-solaris2.11-gcc"),
    ("SunOS", "i86pc", "64", "x86_64-pc-solaris2.11-gcc"),
    ("Linux", "sparc64", "64", "sparc64-linux-gnu-gcc"),
    ("Linux", "i686", "32", "i686-linux-gnu-gcc"),
)


def build(system, machine, size, cc):
    env = dict(os.environ)
    env.pop("CC", None)
    env.pop("HIGHCTIDH_PORTABLE", None)
    if cc is not None:
        env["CC"] = cc
    out = subprocess.run(
        [sys.executable, "-c", PROBE, system, machine, size],
        cwd=SRC,
        env=env,
        check=True,
        capture_output=True,
        text=True,
    ).stdout
    return json.loads(out.strip().splitlines()[-1])


class TestSetupPortable(unittest.TestCase):
    def test_fiat_sources_get_portable_define(self):
        for system, machine, size, cross in PLATFORMS:
            for cc in (None, "", "cc", "gcc", "clang", cross):
                result = build(system, machine, size, cc)
                for bits, (sources, cflags) in result.items():
                    with self.subTest(machine=machine, cc=cc, bits=bits):
                        self.assertIn("fiat_p" + bits + ".c", sources)
                        self.assertIn("-DHIGHCTIDH_PORTABLE=1", cflags)


class TestSetupNoCPUTuning(unittest.TestCase):
    def test_no_cpu_tuning_flags(self):
        platforms = PLATFORMS + (("Linux", "x86_64", "64", "x86_64-linux-gnu-gcc"),)
        for system, machine, size, cross in platforms:
            for cc in (None, "gcc", "clang", cross):
                result = build(system, machine, size, cc)
                for bits, (sources, cflags) in result.items():
                    with self.subTest(machine=machine, cc=cc, bits=bits):
                        bad = [f for f in cflags if f.startswith(("-march", "-mcpu", "-mtune"))]
                        self.assertEqual(bad, [])


class TestSetupDarwinArch(unittest.TestCase):
    def test_no_compiler_arch_macros_on_darwin(self):
        arch = ("-D__x86_64__", "-D__AMD64__", "-D__i386__",
                "-D__arm64__", "-D__aarch64__", "-D__arm__")
        for machine in ("x86_64", "arm64"):
            for cc in (None, "clang"):
                result = build("Darwin", machine, "64", cc)
                for bits, (sources, cflags) in result.items():
                    with self.subTest(machine=machine, cc=cc, bits=bits):
                        bad = [f for f in cflags if f.startswith(arch)]
                        self.assertEqual(bad, [])
                        self.assertIn("-D__Darwin__", cflags)


FAKECC = """#!%s
import json, os, sys
args = sys.argv[1:]
out = args[args.index("-o") + 1]
os.makedirs(os.path.dirname(out) or ".", exist_ok=True)
open(out, "w").close()
with open(os.environ["FAKECC_LOG"], "a") as f:
    f.write(json.dumps(args) + "\\n")
"""


class TestSetupParallelBuild(unittest.TestCase):
    def test_build_ext_parallel_objects(self):
        with tempfile.TemporaryDirectory() as tmp:
            src = os.path.join(tmp, "src")
            shutil.copytree(SRC, src, ignore=shutil.ignore_patterns(
                "build", "*.so", "*.o", "__pycache__", "ctidh*", "go-tests"))
            cc = os.path.join(tmp, "fakecc")
            with open(cc, "w") as f:
                f.write(FAKECC % sys.executable)
            os.chmod(cc, 0o755)
            log = os.path.join(tmp, "log")
            env = dict(os.environ, CC=cc, FAKECC_LOG=log)
            env.pop("LDSHARED", None)
            env.pop("HIGHCTIDH_PORTABLE", None)
            subprocess.run([sys.executable, "setup.py", "build_ext", "-j", "4"],
                           cwd=src, env=env, check=True, capture_output=True)
            sources = build("Linux", "x86_64", "64", None)
            want = sum(len(v[0]) for v in sources.values())
            objects = {}
            links = []
            with open(log) as f:
                for line in f:
                    args = json.loads(line)
                    out = args[args.index("-o") + 1]
                    if "-c" in args:
                        bits = [a[len("-DBITS="):] for a in args if a.startswith("-DBITS=")]
                        objects.setdefault(out, set()).update(bits)
                    else:
                        links.append((out, [a for a in args if a.endswith(".o")]))
            self.assertEqual(len(objects), want)
            for out, bits in objects.items():
                self.assertEqual(len(bits), 1, out)
            self.assertEqual(len(links), 4)
            for out, objs in links:
                size = re.search(r"highctidh_(\d+)", os.path.basename(out)).group(1)
                for o in objs:
                    self.assertEqual(objects[o], {size}, (out, o))


if __name__ == "__main__":
    unittest.main()

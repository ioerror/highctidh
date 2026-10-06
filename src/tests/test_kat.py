#!/usr/bin/env python3

import os
import pathlib
import struct
import unittest

from highctidh import ctidh

KAT_DIR = pathlib.Path(__file__).resolve().parent.parent / "kat"
FIELDS = (
    "size", "rng",
    "seed_a", "sk_a", "pk_a",
    "seed_b", "sk_b", "pk_b",
    "seed_f", "sk_f",
    "ss_a_b", "blind_f_a",
)
MASK = 0xFFFFFFFFFFFFFFFF


def splitmix64(state):
    buf = bytearray()

    def rng(out, context):
        nonlocal state
        while len(buf) < len(out):
            state = (state + 0x9E3779B97F4A7C15) & MASK
            z = state
            z = ((z ^ (z >> 30)) * 0xBF58476D1CE4E5B9) & MASK
            z = ((z ^ (z >> 27)) * 0x94D049BB133111EB) & MASK
            z ^= z >> 31
            buf.extend(struct.pack("<Q", z))
        out[:] = buf[:len(out)]
        del buf[:len(out)]

    return rng


def read_kat(size):
    entry = {}
    with open(KAT_DIR / f"ctidh{size}.kat") as f:
        for line in f.read().splitlines():
            k, v = line.split(" = ", 1)
            entry[k] = v
    assert entry["size"] == str(size) and entry["rng"] == "splitmix64"
    return entry


def result(f, *args):
    try:
        return bytes(f(*args)).hex()
    except Exception as e:
        return f"error: {e}"


def keygen(c, entry, n):
    rng = splitmix64(int(entry["seed_" + n]))
    return bytes(c.generate_secret_key(rng=rng)).hex()


def pub(c, sk):
    return c.derive_public_key(c.private_key_from_hex(sk))


def action(c, sk, pk):
    return c.blind(c.private_key_from_hex(sk), c.public_key_from_hex(pk))


def generate(size):
    c = ctidh(size)
    entry = read_kat(size)
    g = {k: entry[k] for k in ("size", "rng", "seed_a", "seed_b", "seed_f")}
    for n in "abf":
        g["sk_" + n] = keygen(c, entry, n)
    g["pk_a"] = result(pub, c, g["sk_a"])
    g["pk_b"] = result(pub, c, g["sk_b"])
    g["ss_a_b"] = result(action, c, g["sk_a"], g["pk_b"])
    g["blind_f_a"] = result(action, c, g["sk_f"], g["pk_a"])
    return "".join(f"{k} = {g[k]}\n" for k in FIELDS)


class TestKAT(unittest.TestCase):
    def check(self, size):
        c = ctidh(size)
        e = read_kat(size)
        got = {}
        for n in "abf":
            got["sk_" + n] = keygen(c, e, n)
        got["pk_a"] = result(pub, c, e["sk_a"])
        got["pk_b"] = result(pub, c, e["sk_b"])
        got["ss_a_b"] = result(action, c, e["sk_a"], e["pk_b"])
        got["blind_f_a"] = result(action, c, e["sk_f"], e["pk_a"])
        bad = [
            f"{k}: got {v} want {e[k]}" for k, v in got.items() if v != e[k]
        ]
        self.assertEqual(bad, [], "\n" + "\n".join(bad))

    def test_kat_511(self):
        self.check(511)

    def test_kat_512(self):
        self.check(512)

    def test_kat_1024(self):
        self.check(1024)

    def test_kat_2048(self):
        self.check(2048)


if __name__ == "__main__":
    out = os.environ.get("HIGHCTIDH_KAT_OUT")
    if out:
        for size in (511, 512, 1024, 2048):
            path = pathlib.Path(out) / f"ctidh{size}.kat"
            path.write_text(generate(size))
    else:
        unittest.main()

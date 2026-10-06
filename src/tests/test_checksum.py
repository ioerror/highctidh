#!/usr/bin/env python3

import os
import pathlib
import struct
import unittest

from highctidh import ctidh

KAT_DIR = pathlib.Path(__file__).resolve().parent.parent / "kat"
MASK = 0xFFFFFFFF
LOOPS = {"small": 64, "big": 512}


def rotl(v, c):
    return ((v << c) & MASK) | (v >> (32 - c))


def chacha_block(key, counter):
    s = [0x61707865, 0x3320646E, 0x79622D32, 0x6B206574]
    s += list(struct.unpack("<8I", key))
    s += [counter & MASK, counter >> 32, 0, 0]
    x = list(s)

    def qr(a, b, c, d):
        x[a] = (x[a] + x[b]) & MASK
        x[d] = rotl(x[d] ^ x[a], 16)
        x[c] = (x[c] + x[d]) & MASK
        x[b] = rotl(x[b] ^ x[c], 12)
        x[a] = (x[a] + x[b]) & MASK
        x[d] = rotl(x[d] ^ x[a], 8)
        x[c] = (x[c] + x[d]) & MASK
        x[b] = rotl(x[b] ^ x[c], 7)

    for _ in range(10):
        qr(0, 4, 8, 12)
        qr(1, 5, 9, 13)
        qr(2, 6, 10, 14)
        qr(3, 7, 11, 15)
        qr(0, 5, 10, 15)
        qr(1, 6, 11, 12)
        qr(2, 7, 8, 13)
        qr(3, 4, 9, 14)
    return struct.pack("<16I", *[(a + b) & MASK for a, b in zip(x, s)])


class KnownRandom:
    def __init__(self):
        self.key = bytes(32)
        self.pool = b""

    def read(self, n):
        out = bytearray()
        while len(out) < n:
            if not self.pool:
                stream = b"".join(
                    chacha_block(self.key, j) for j in range(12))
                self.key = stream[:32]
                self.pool = stream[32:]
            take = min(n - len(out), len(self.pool))
            out += self.pool[:take]
            self.pool = self.pool[take:]
        return bytes(out)

    def callback(self, buf, context):
        buf[:] = self.read(len(buf))


def salsa_core(inp, key):
    s = [0] * 16
    s[0], s[5], s[10], s[15] = 0x61707865, 0x3320646E, 0x79622D32, 0x6B206574
    k = struct.unpack("<8I", key)
    s[1:5] = k[:4]
    s[11:15] = k[4:]
    s[6:10] = struct.unpack("<4I", inp)
    x = list(s)

    def qr(a, b, c, d):
        x[b] ^= rotl((x[a] + x[d]) & MASK, 7)
        x[c] ^= rotl((x[b] + x[a]) & MASK, 9)
        x[d] ^= rotl((x[c] + x[b]) & MASK, 13)
        x[a] ^= rotl((x[d] + x[c]) & MASK, 18)

    for _ in range(10):
        qr(0, 4, 8, 12)
        qr(5, 9, 13, 1)
        qr(10, 14, 2, 6)
        qr(15, 3, 7, 11)
        qr(0, 1, 2, 3)
        qr(5, 6, 7, 4)
        qr(10, 11, 8, 9)
        qr(15, 12, 13, 14)
    return struct.pack("<16I", *[(a + b) & MASK for a, b in zip(x, s)])


class Checksum:
    def __init__(self):
        self.state = bytes(64)

    def add(self, x):
        while len(x) >= 16:
            self.state = salsa_core(x[:16], self.state[:32])
            x = x[16:]
        last = bytearray(16)
        last[: len(x)] = x
        last[len(x)] = 1
        self.state = bytes([self.state[0] ^ 1]) + self.state[1:]
        self.state = salsa_core(bytes(last), self.state[:32])

    def hex(self):
        return self.state[:32].hex()


def try_dh(size, loops):
    c = ctidh(size)
    rng = KnownRandom()
    total = Checksum()

    def keypair():
        sk = c.generate_secret_key(rng=rng.callback)
        return bytes(c.derive_public_key(sk)), bytes(sk)

    def dh(pk, sk):
        return bytes(
            c.blind(c.private_key_from_bytes(sk), c.public_key_from_bytes(pk))
        )

    for _ in range(loops):
        pk_a, sk_a = keypair()
        total.add(pk_a)
        total.add(sk_a)
        pk_b, sk_b = keypair()
        total.add(pk_b)
        total.add(sk_b)
        ss_a = dh(pk_b, sk_a)
        total.add(ss_a)
        if dh(pk_b, sk_a) != ss_a:
            raise AssertionError("dh is not deterministic")
        ss_b = dh(pk_a, sk_b)
        total.add(ss_b)
        if dh(pk_a, sk_b) != ss_b:
            raise AssertionError("dh is not deterministic")
        if ss_a != ss_b:
            raise AssertionError("shared secrets differ")
    return total.hex()


def read_checksums():
    m = {}
    for line in (KAT_DIR / "checksums").read_text().splitlines():
        size, kind, value = line.split()
        m[(int(size), kind)] = value
    return m


def sizes():
    v = os.environ.get("HIGHCTIDH_CHECKSUM_SIZES")
    return [int(w) for w in v.split()] if v else [511, 512]


def kinds():
    return ["small", "big"] if os.environ.get("HIGHCTIDH_CHECKSUM_BIG") else [
        "small"
    ]


class TestChecksums(unittest.TestCase):
    def test_checksums(self):
        want = read_checksums()
        for size in sizes():
            for kind in kinds():
                with self.subTest(size=size, kind=kind):
                    got = try_dh(size, LOOPS[kind])
                    self.assertEqual(got, want[(size, kind)])


if __name__ == "__main__":
    unittest.main()

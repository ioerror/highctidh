#!/usr/bin/env python3

import pathlib
import unittest

from highctidh import ctidh

KAT_DIR = pathlib.Path(__file__).resolve().parent.parent / "kat"


def exponents(s):
    return bytes(int(w) & 0xFF for w in s.split(","))


def read_cases(size):
    cases = []
    path = KAT_DIR / f"csidh-reference-{size}.txt"
    for n, line in enumerate(path.read_text().splitlines(), 1):
        w = line.split()
        fields = dict(kv.split("=", 1) for kv in w[1:])
        cases.append((n, w[0], fields))
    return cases


class TestCSIDHReference(unittest.TestCase):
    def check(self, size):
        c = ctidh(size)
        bad = []
        for n, kind, fields in read_cases(size):
            sk = c.private_key_from_bytes(exponents(fields["e"]))
            pk = c.derive_public_key(sk)
            if kind == "ss":
                f = c.private_key_from_bytes(exponents(fields["f"]))
                pk = c.blind(f, pk)
            got = bytes(pk).hex()
            if got != fields["A"]:
                bad.append(f"line {n} {kind}: got {got} want {fields['A']}")
        self.assertEqual(bad, [], "\n" + "\n".join(bad))

    def test_csidh_reference_512(self):
        self.check(512)

    def test_csidh_reference_1024(self):
        self.check(1024)


if __name__ == "__main__":
    unittest.main()

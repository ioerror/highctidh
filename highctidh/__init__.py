#!/usr/bin/env python3

import ctypes
import ctypes.util
import hashlib
import struct
import pathlib
from importlib import util, metadata
from importlib.metadata import PackageNotFoundError
from .__version__ import version as highctidh_version

try:
    __version__ = metadata.version('highctidh')
except PackageNotFoundError:
    __version__ = highctidh_version

class InvalidFieldSize(Exception):
    """Raised when a field is not one of (511, 512, 1024, 2048)."""

    pass


class InvalidPublicKey(Exception):
    """Raised when a public key is not validated by the validate() function."""

    pass

class DecodingError(Exception):
    '''Raised when a serialized value is not in the expected format.'''

    pass

class CSIDHError(Exception):
    """Raised when csidh() fails to return True."""

    pass


class LibraryNotFound(Exception):
    """Raised when the shared library cannot be located and opened."""

    pass


class ctidh(object):
    def __init__(self, field_size):
        self._ctidh_sizes = (511, 512, 1024, 2048)
        self.field_size = field_size
        ctidh_self = self
        if self.field_size not in self._ctidh_sizes:
            raise InvalidFieldSize(f"Unsupported size: {repr(self.field_size)}")
        if self.field_size == 511:
            self.pk_size = 64
            self.sk_size = 74
        elif self.field_size == 512:
            self.pk_size = 64
            self.sk_size = 74
        elif self.field_size == 1024:
            self.pk_size = 128
            self.sk_size = 130
        elif self.field_size == 2048:
            self.pk_size = 256
            self.sk_size = 231

        class private_key(ctypes.Structure):
            __slots__ = [
                "e",
            ]
            _fields_ = [
                ("e", ctypes.c_ubyte * self.sk_size),
            ]

            def __bytes__(self) -> bytes:
                """canonical byte representation"""
                return bytes(self.e)

            def __repr__(self):
                return f'<highctidh.ctidh({ctidh_self.field_size}).private_key>'

            @classmethod
            def frombytes(cls, byt:bytes):
                '''restore bytes(private_key_object) in canonical byte
                representation as a private_key() instance'''
                if type(byt) is not bytes:
                    raise DecodingError("Private key is not bytes")
                if len(byt) != ctypes.sizeof(cls):
                    raise DecodingError(f"Serialized private key should be {self.sk_size}, is {len(byt)}")
                return cls.from_buffer(ctypes.create_string_buffer(byt))

            @classmethod
            def fromhex(cls, h:str):
                '''restore bytes(private_key_object).hex() in canonical byte
                representation as a private_key() instance'''
                if type(h) is not str:
                    raise DecodingError("Private key is not str")
                try:
                    h = bytes.fromhex(h)
                except Exception as e:
                    raise DecodingError(e)
                if len(h) != ctypes.sizeof(cls):
                    # This gets bad if len(h) < sizeof
                    raise DecodingError(f"Private key must be {ctypes.sizeof(private_key)} bytes, is {len(h)}")
                return cls.from_buffer(ctypes.create_string_buffer(h))

            def derive_public_key(self):
                """Compute return the corresponding public key *pk*."""
                pk = ctidh_self.public_key()
                ctidh_self.csidh(pk, ctidh_self.base, self)
                return pk

        assert self.sk_size == ctypes.sizeof(private_key), (
            self.sk_size, ctypes.sizeof(private_key),)

        self.private_key = private_key

        class public_key(ctypes.Structure):
            __slots__ = [
                "A",
            ]
            _fields_ = [("A", ctypes.c_uint64 * (self.pk_size // 8))]

            def __bytes__(self):
                """pack to canonical little-endian representation"""
                return struct.pack("<" + "Q" * (ctypes.sizeof(self) // 8), *self.A)

            def __repr__(self):
                return f'<highctidh.ctidh({ctidh_self.field_size}).public_key>'

            @classmethod
            def frombytes(cls, byt:bytes, validate=True):
                '''Restores a public_key instance from canonical little-endian representation.
                If the optional validate= parameter is True, the public key is validated.
                If not, the caller must ensure that the public key is valid and
                comes from a trusted source.
                '''
                if type(byt) is not bytes:
                    raise DecodingError("Public key is not bytes")
                # public keys are transferred in little-endian; we might have to
                # byteswap to get them in native order in pk.A:
                pk = self.public_key()
                try:
                    pk.A[:] = struct.unpack("<" + "Q" * (self.pk_size // 8), byt)
                except struct.error as e:
                    raise DecodingError(e)
                if validate:
                    assert ctidh_self.validate(pk)
                return pk

            @classmethod
            def fromhex(cls, h:str):
                return cls.frombytes(bytes.fromhex(h))

        self.public_key = public_key

        assert self.pk_size == ctypes.sizeof(public_key), (
            self.pk_size, ctypes.sizeof(public_key))

        self.base = self.public_key()
        try:
            flib = f"highctidh_{self.field_size}"
            flib = util.find_spec(flib).origin
            self._lib = ctypes.CDLL(flib)
        except OSError as e:
            print("Unable to load highctidh_" + str(self.field_size) + ".so".format(e))
            raise LibraryNotFound

        csidh_private = self._lib.__getattr__(
            "highctidh_" + str(self.field_size) + "_csidh_private"
        )
        csidh_private.restype = None
        csidh_private.argtypes = [ctypes.POINTER(self.private_key)]
        self.csidh_private = csidh_private

        csidh_private_withrng = self._lib.__getattr__(
            "highctidh_" + str(self.field_size) + "_csidh_private_withrng"
        )
        csidh_private_withrng.restype = None
        csidh_private_withrng.argtypes = [
            ctypes.POINTER(self.private_key),
            ctypes.c_void_p,
            ctypes.c_void_p,
        ]
        self.csidh_private_withrng = csidh_private_withrng

        csidh = self._lib.__getattr__("highctidh_" + str(self.field_size) + "_csidh")
        csidh.restype = bool
        csidh.argtypes = [
            ctypes.POINTER(self.public_key),
            ctypes.POINTER(self.public_key),
            ctypes.POINTER(self.private_key),
        ]
        self._csidh = csidh
        validate = self._lib.__getattr__(
            "highctidh_" + str(self.field_size) + "_validate"
        )
        validate.restype = bool
        validate.argtypes = [ctypes.POINTER(self.public_key)]
        self._validate = validate

    def private_key_from_bytes(self, h):
        return self.private_key.frombytes(h)

    def public_key_from_bytes(self, h:bytes):
        '''Restores a public_key instance from bytes in the canonical
        (little-endian) representation'''
        return self.public_key.frombytes(h)

    def private_key_from_hex(self, h):
        return self.private_key.fromhex(h)

    def public_key_from_hex(self, h):
        return self.public_key.fromhex(h)

    def validate(self, pk):
        """self._validate returns 1 if successful, 0 if invalid."""
        if self._validate(pk):
            return True
        else:
            raise InvalidPublicKey

    def csidh(self, pk0, pk1, sk):
        if self._csidh(pk0, pk1, sk):
            return True
        else:
            raise CSIDHError

    def generate_secret_key_inplace(self, sk, rng=None, context=None):
        """Generate a secret key *sk* without allocating memory, overwriting *sk*.
        Optionally takes a callable argument *rng* which is called with two
        arguments: rng(buf, context)
        Where *buf* is a bytearray() to be filled with random data and
        *context* is an int() context identifier to enable thread-safe calls.
        If *context* is left blank, it is a pointer to the buffer.
        Note that in order to achieve portable reproducible results, a PRNG
        must fill buf as though it was an array of int32_t values in
        HOST-ENDIAN/NATIVE byte order, see comment in csidh.h:ctidh_fillrandom
        """
        if rng:

            @ctypes.CFUNCTYPE(None, ctypes.c_void_p, ctypes.c_size_t, ctypes.c_void_p)
            def rng_callback_wrapper(buf, bufsz, context):
                mv = memoryview(
                    ctypes.cast(buf, ctypes.POINTER(ctypes.c_byte * bufsz)).contents
                ).cast(
                    "B"
                )  # uint8_t
                rng(mv, context)
            if context is None:
                context = ctypes.byref(sk.e)
            self.csidh_private_withrng(sk, context, rng_callback_wrapper)
        else:
            self.csidh_private(sk)
        return sk

    def generate_secret_key(self, rng=None, context=None):
        """Generate a secret key *sk*, return it.
        Optionally takes a callable argument *rng* which is called with two
        arguments: rng(buf, context)
        Where *buf* is a bytearray() to be filled with random data and
        *context* is an int() context identifier to enable thread-safe calls.
        If *context* is left blank, it is a pointer to the buffer.
        Note that in order to achieve portable reproducible results, a PRNG
        must fill buf as though it was an array of int32_t values in
        HOST-ENDIAN/NATIVE byte order, see comment in csidh.h:ctidh_fillrandom
        """
        return self.generate_secret_key_inplace(
            self.private_key(), rng=rng, context=context)

    def derive_public_key(self, sk):
        """Given a secret key *sk*, return the corresponding public key *pk*."""
        return sk.derive_public_key()

    def dh(
        self,
        sk,
        pk,
        _hash=lambda shared, size: hashlib.shake_256(bytes(shared)).digest(length=size),
    ):
        """
        This is a classic Diffie-Hellman function which takes a secret key
        *sk* and a public key *pk* and computes a random element. It then
        computes a uniformly random bit string from the random element using a
        variable length hash function. The returned value is a bytes() object.
        The size of the hash output is dependent on the field size. The *_hash*
        may be overloaded as needed.
        """
        assert type(pk) is self.public_key
        assert type(sk) is self.private_key
        shared_key = self.public_key()
        self.csidh(shared_key, pk, sk)
        return _hash(bytes(shared_key), self.pk_size)

    def blind(self, blinding_factor_sk, pk):
        blinded_key = self.public_key()
        self.csidh(blinded_key, pk, blinding_factor_sk)
        return blinded_key

    def blind_dh(self, blind_sk, sk, pk):
        """
        This is a Diffie-Hellman function for use with blinding factors. This
        is a specialized function that should only be use with very specific
        cryptographic protocols such as those using the sphinx packet format.
        It takes a blinding factor *blind* that is a private_key() object, a
        second secret key *sk* that is a private_key() object, a public key
        *pk* that is a public_key() object, and it computes a random element.
        The returned value is a public_key() object. This function is suitable
        for use with a shared blinding factor, and the eventual secret returned
        should be made into a uniformly random bit string by the caller unless
        it is used in a protocol that requires the use of the random element.
        """
        shared_key = self.public_key()
        blinded_result = self.public_key()
        if self.csidh(shared_key, pk, sk):
            self.csidh(blinded_result, shared_key, blind_sk)
            return blinded_result
        else:
            raise CSIDHError

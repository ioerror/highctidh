#!/usr/bin/env python3
"""
highctidh wraps the original highctidh C implementation using ctypes.

>>> from highctidh import ctidh
>>> ctidh511 = ctidh(511)
>>> ctidh512 = ctidh(512)
>>> ctidh1024 = ctidh(1024)
>>> ctidh2048 = ctidh(2048)
"""

import ctypes
import ctypes.util
import hashlib
import struct
import pathlib
from importlib import util, metadata
from importlib.metadata import PackageNotFoundError
try:
    from .__version__ import version as highctidh_version
except ImportError:
    highctidh_version = "unknown"

try:
    __version__ = metadata.version('highctidh')
except PackageNotFoundError:
    __version__ = highctidh_version

class InvalidFieldSize(Exception):
    """Raised when a field is not one of (511, 512, 1024, 2048).

    >>> from highctidh import InvalidFieldSize
    >>> raise InvalidFieldSize
    Traceback (most recent call last):
    ...
    highctidh.InvalidFieldSize
    """

    pass


class InvalidPublicKey(Exception):
    """Raised when a public key is not validated by the validate() function.

    >>> from highctidh import InvalidPublicKey
    >>> raise InvalidPublicKey
    Traceback (most recent call last):
    ...
    highctidh.InvalidPublicKey
    """

    pass

class DecodingError(Exception):
    """Raised when a serialized value is not in the expected format.

    >>> from highctidh import DecodingError
    >>> raise DecodingError
    Traceback (most recent call last):
    ...
    highctidh.DecodingError
    """

    pass

class CSIDHError(Exception):
    """Raised when csidh() fails to return True.

    >>> from highctidh import CSIDHError
    >>> raise CSIDHError
    Traceback (most recent call last):
    ...
    highctidh.CSIDHError
    """

    pass


class LibraryNotFound(Exception):
    """Raised when the shared library cannot be located and opened.

    >>> from highctidh import LibraryNotFound
    >>> raise LibraryNotFound
    Traceback (most recent call last):
    ...
    highctidh.LibraryNotFound
    """

    pass


class ctidh(object):
    """
    Instantiate a ctidh object by setting a field size for all subsequent
    operations. Field sizes available are 511, 512, 1024, and 2048.

    >>> from highctidh import ctidh
    >>> ctidh511 = ctidh(511)
    >>> ctidh512 = ctidh(512)
    >>> ctidh1024 = ctidh(1024)
    >>> ctidh2048 = ctidh(2048)
    """
    def __init__(self, field_size):
        """
        Instantiate a ctidh object by setting a field size for all subsequent
        operations. Field sizes available are 511, 512, 1024, and 2048.

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> ctidh512 = ctidh(512)
        >>> ctidh1024 = ctidh(1024)
        >>> ctidh2048 = ctidh(2048)
        """
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
                """
                Restores a public_key instance from canonical little-endian representation.
                If the optional validate= parameter is True, the public key is validated.
                If not, the caller must ensure that the public key is valid and
                comes from a trusted source.
                """
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
                """
                Restores a public_key instance from canonical little-endian
                representation encoded as hex.
                """
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

    def private_key_from_bytes(self, h:bytes):
        """
        Restores a private_key instance from bytes in the canonical
        (little-endian) representation

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> sk511_a = ctidh511.generate_secret_key()
        >>> sk511_b = ctidh511.private_key_from_bytes(bytes(sk511_a))
        >>> bytes(sk511_a) == bytes(sk511_b)
        True
        """
        return self.private_key.frombytes(h)

    def public_key_from_bytes(self, h:bytes):
        """
        Restores a public_key instance from bytes in the canonical
        (little-endian) representation

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> sk511_a = ctidh511.generate_secret_key()
        >>> sk511_b = ctidh511.private_key_from_bytes(bytes(sk511_a))
        >>> pk_b = bytes(sk511_b.derive_public_key())
        >>> pk511_b = ctidh511.public_key_from_bytes(pk_b)
        >>> pk511_a = bytes(sk511_a.derive_public_key())
        >>> bytes(pk511_a) == bytes(pk511_b)
        True
        """
        return self.public_key.frombytes(h)

    def private_key_from_hex(self, h:str):
        """
        Restores a private_key instance from hex in the canonical
        (little-endian) representation

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> sk511_a = ctidh511.generate_secret_key()
        >>> sk511_b = ctidh511.private_key_from_hex(bytes(sk511_a).hex())
        >>> bytes(sk511_a) == bytes(sk511_b)
        True
        """
        return self.private_key.fromhex(h)

    def public_key_from_hex(self, h:str):
        """
        Restores a public_key instance from hex in the canonical
        (little-endian) representation

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> sk511_a = ctidh511.generate_secret_key()
        >>> sk511_b = ctidh511.private_key_from_bytes(bytes(sk511_a))
        >>> pk_b_hex = bytes(sk511_b.derive_public_key()).hex()
        >>> pk511_b = ctidh511.public_key_from_hex(pk_b_hex)
        >>> pk511_a = bytes(sk511_a.derive_public_key())
        >>> bytes(pk511_a) == bytes(pk511_b)
        True
        """
        return self.public_key.fromhex(h)

    def validate(self, pk):
        """
        self._validate returns 1 if successful, or raises *InvalidPublicKey* if
        invalid.

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> sk511_a = ctidh511.generate_secret_key()
        >>> pk511_a = ctidh511.derive_public_key(sk511_a)
        >>> ctidh511.validate(pk511_a)
        True
        """
        if self._validate(pk):
            return True
        else:
            raise InvalidPublicKey

    def csidh(self, pk0, pk1, sk):
        """
        This function computes a group action over *pk0*, *pk1*, and *sk* and
        ensures that it is validated. It returns *True* if validation was
        successful, and it raises *CSIDHError* if validation was not
        successful.

        This function mutates the value of *pk1* and returns a bool that
        indicates if the result passes csidh validation as this is what the
        original C library does. This mutation may make some Python programmers
        unhappy and this documentation serves to warning to would-be callers.

        This API is likely to change to ensure that mutation is hidden from
        Python callers and as a result it should be considered experimental.

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> sk511_a = ctidh511.generate_secret_key()
        >>> pk511_a = ctidh511.derive_public_key(sk511_a)
        >>> sk511_b = ctidh511.generate_secret_key()
        >>> pk511_b = ctidh511.derive_public_key(sk511_b)
        >>> sk511_c = ctidh511.generate_secret_key()
        >>> pk511_c = ctidh511.derive_public_key(sk511_b)
        >>> ctidh511.csidh(pk511_a, pk511_b, sk511_c)
        True
        >>> pk511_b == ctidh511.derive_public_key(sk511_b)
        False
        """
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

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> sk511_a = ctidh511.generate_secret_key()
        >>> sk511_b = ctidh511.generate_secret_key_inplace(sk511_a)
        >>> sk511_b
        <highctidh.ctidh(511).private_key>
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

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> sk511_a = ctidh511.generate_secret_key()
        >>> sk511_a
        <highctidh.ctidh(511).private_key>
        """
        return self.generate_secret_key_inplace(
            self.private_key(), rng=rng, context=context)

    def derive_public_key(self, sk):
        """Given a secret key *sk*, return the corresponding public key *pk*.

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> sk511_a = ctidh511.generate_secret_key()
        >>> pk511_a = sk511_a.derive_public_key()
        >>> pk511_a
        <highctidh.ctidh(511).public_key>
        """
        return sk.derive_public_key()

    def dh(
        self,
        sk,
        pk,
        _hash=lambda shared, size: hashlib.shake_256(bytes(shared)).digest(size),
    ):
        """
        This is a classic Diffie-Hellman function which takes a secret key
        *sk* and a public key *pk* and computes a random element. It then
        computes a uniformly random bit string from the random element using a
        variable length hash function. The returned value is a bytes() object.
        The size of the hash output is dependent on the field size. The *_hash*
        may be overloaded as needed.

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> sk511_a = ctidh511.generate_secret_key()
        >>> pk511_a = ctidh511.derive_public_key(sk511_a)
        >>> sk511_b = ctidh511.generate_secret_key()
        >>> pk511_b = ctidh511.derive_public_key(sk511_b)
        >>> ctidh511.dh(sk511_a, pk511_b) == ctidh511.dh(sk511_b, pk511_a)
        True
        """
        assert type(pk) is self.public_key
        assert type(sk) is self.private_key
        shared_key = self.public_key()
        self.csidh(shared_key, pk, sk)
        return _hash(bytes(shared_key), self.pk_size)

    def blind(self, blinding_factor_sk, pk):
        """
        This blind function takes a secret blinding factor key as
        *blinding_factor_sk* and a public key pk. It calls the *csidh* group
        action and returns the blinded public key *blinded_key*.

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> sk511_a = ctidh511.generate_secret_key()
        >>> pk511_a = ctidh511.derive_public_key(sk511_a)
        >>> sk511_b = ctidh511.generate_secret_key()
        >>> pk511_b = ctidh511.derive_public_key(sk511_b)
        >>> ctidh511.blind(sk511_a, pk511_b)
        <highctidh.ctidh(511).public_key>
        """
        blinded_key = self.public_key()
        self.csidh(blinded_key, pk, blinding_factor_sk)
        return blinded_key

    def blind_dh(self, blind_sk, sk, pk):
        """
        This is a Diffie-Hellman function for use with blinding factors. This
        is a specialized function that should only be use with very specific
        cryptographic protocols such as those using the sphinx packet format.
        It takes a blinding factor *blind_sk* that is a private_key() object, a
        second secret key *sk* that is a private_key() object, a public key
        *pk* that is a public_key() object, and it computes a random element.
        The returned value is a public_key() object. This function is suitable
        for use with a shared blinding factor, and the eventual secret returned
        should be made into a uniformly random bit string by the caller unless
        it is used in a protocol that requires the use of the random element.

        >>> from highctidh import ctidh
        >>> ctidh511 = ctidh(511)
        >>> sk511_a = ctidh511.generate_secret_key()
        >>> pk511_a = ctidh511.derive_public_key(sk511_a)
        >>> sk511_b = ctidh511.generate_secret_key()
        >>> pk511_b = ctidh511.derive_public_key(sk511_b)
        >>> sk511_c = ctidh511.generate_secret_key()
        >>> pk511_c = ctidh511.derive_public_key(sk511_c)
        >>> ctidh511.blind_dh(sk511_a, sk511_b, pk511_c)
        <highctidh.ctidh(511).public_key>
        """
        shared_key = self.public_key()
        blinded_result = self.public_key()
        if self.csidh(shared_key, pk, sk):
            self.csidh(blinded_result, shared_key, blind_sk)
            return blinded_result
        else:
            raise CSIDHError

if __name__ == '__main__':
    import doctest
    doctest.testmod(verbose=True)

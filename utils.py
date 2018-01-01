import attr
import collections
import binascii
import functools
import itertools
import operator

SOLUTIONS = {
}


@attr.s(slots=True, hash=True)
class Challenge:
    problemset = attr.ib()
    problem = attr.ib()

    @classmethod
    def of(cls, problemset, problem):
        return cls(problemset=problemset, problem=problem)


def solves(challenge, second=None):
    if second is not None:
        challenge = Challenge.of(challenge, second)

    def _dec(f):
        SOLUTIONS[challenge] = f
        return f
    return _dec


def solve(challenge: Challenge, *args, **kwargs):
    try:
        return SOLUTIONS[challenge](*args, **kwargs)
    except KeyError as e:
        raise ValueError(f"I don't know how to solve {challenge}. Have you written a function which @solves it?") from e


hex_to_bytes = binascii.unhexlify
bytes_to_b64 = binascii.b2a_base64
hex_to_b64 = lambda hex: bytes_to_b64(hex_to_bytes(hex)).strip()


def xor(s, t):
    if isinstance(t, int):
        return bytearray(si ^ t for si in s)
    elif len(s) == len(t):
        return bytearray(itertools.starmap(operator.xor, zip(s, t)))
    elif len(t) == 1:
        return xor(s, ord(t))


def english_score(plaintext):
    import ipdb; ipdb.set_trace()
    collections.Counter(plaintext)

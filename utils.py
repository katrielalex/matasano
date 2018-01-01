import attr
import functools


SOLUTIONS = {
}


@attr.s(slots=True, hash=True)
class Challenge:
    problemset = attr.ib()
    problem = attr.ib()

    @classmethod
    def of(cls, problemset, problem):
        return cls(problemset=problemset, problem=problem)


def solves(challenge: Challenge):
    def _dec(f):
        SOLUTIONS[challenge] = f
        return f
    return _dec


def solve(challenge: Challenge, *args, **kwargs):
    try:
        return SOLUTIONS[challenge](*args, **kwargs)
    except KeyError as e:
        raise ValueError(f"I don't know how to solve {challenge}. Have you written a function which @solves it?") from e

#!/usr/bin/env python3

import click
import sys
from utils import Challenge, solve

import set1

def validate_challenge(ctx, param, value):
    try:
        problemset, problem = map(int, value.split("."))
    except ValueError:
        raise click.BadParameter(f"problem '{value}' was not in format x.y")
    return Challenge.of(problemset, problem)


@click.command()
@click.argument("challenge", callback=validate_challenge)
@click.argument("data", nargs=-1)
def mata(challenge, data):
    return solve(challenge, *data)
    print(f"set number {challenge}")

if __name__ == "__main__":
    sys.exit(mata())

#!/usr/bin/env python3
"""Credential-free validation of the acceptance-test policy files.

Runs on every PR. Deterministic, no API key, a few seconds, so it can be a
required check from day one. Each check exists because the mistake it catches
otherwise produces a GREEN run that verified less than it appears to:

  - a product, group, smoke entry or exclusion that resolves to nothing
  - a package that is neither planned exactly once nor explicitly excluded
  - a package whose resources no sweeper reclaims reaching an executable plan
  - an accepted failure with no ticket or reason

Warnings (stale baseline or exception entries) are annotated but do not fail:
a deleted test cannot fail, so its baseline entry is only clutter.

All checks run before exiting, so one invocation shows every problem.
"""

import os
import sys

sys.dont_write_bytecode = True  # keep __pycache__ out of .github
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import acctest_lib as lib  # noqa: E402


def main():
    root = os.path.abspath(sys.argv[1] if len(sys.argv) > 1 else ".")
    problems, warnings, notes = lib.validate(root)
    for note in notes:
        print("  %s" % note)
    for w in warnings:
        print("::warning::%s" % w)
    if problems:
        print("\nFAIL: %d problem(s)\n" % len(problems), file=sys.stderr)
        for p in problems:
            print("::error::%s" % p, file=sys.stderr)
        sys.exit(1)
    print("\nOK: policy files are consistent with the repository")


if __name__ == "__main__":
    main()

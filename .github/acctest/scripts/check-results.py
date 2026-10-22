#!/usr/bin/env python3
"""Classify one job's go test -json output against its plan and the baseline.

    check-results.py --job job.json results/*.json

Writes a Markdown report to stdout (append it to $GITHUB_STEP_SUMMARY). Exits 1
when the job has unexpected failures, missing results, unexpected skips or
package-level errors. Accepted failures are listed separately with their group
and do not affect the exit status.
"""

import argparse
import json
import os
import sys

sys.dont_write_bytecode = True  # keep __pycache__ out of .github
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import acctest_lib as lib  # noqa: E402

SECTIONS = (
    ("unexpected_failures", "Unexpected failures"),
    ("missing", "Missing results (expected, never reported)"),
    ("unexpected_skips", "Unexpected skips"),
    ("package_errors", "Package errors (failed with no failing test)"),
    ("accepted", "Accepted failures (known-failures.json)"),
    ("now_passing", "Baselined but now passing - remove from known-failures.json"),
)


def main(argv=None):
    ap = argparse.ArgumentParser()
    ap.add_argument("--repo-root", default=".")
    ap.add_argument("--job", required=True, help="one job object from plan.py")
    ap.add_argument("results", nargs="+", help="go test -json output files")
    args = ap.parse_args(argv)

    root = os.path.abspath(args.repo_root)
    job = lib.load_json(args.job)
    packages, _ = lib.discover(root)
    lines = []
    for path in args.results:
        with open(path, encoding="utf-8") as fh:
            lines.extend(fh)
    results, package_failed = lib.parse_go_test_json(lines)
    outcome = lib.classify(job, packages, results, package_failed,
                           lib.accepted_failures(lib.load_config(root)))

    verdict = "FAILED" if lib.failed(outcome) else "passed"
    print("### Acceptance tests: %s (%s)\n" % (job["name"], verdict))
    print("%d passed, %d unexpected failures, %d accepted failures, %d missing, %d skipped\n"
          % (outcome["passed"], len(outcome["unexpected_failures"]), len(outcome["accepted"]),
             len(outcome["missing"]), len(outcome["unexpected_skips"])))
    for key, title in SECTIONS:
        if outcome[key]:
            print("**%s**\n" % title)
            for item in outcome[key]:
                if isinstance(item, tuple):
                    print("- `%s` (%s)" % item)
                else:
                    print("- `%s`" % item)
            print()
    sys.exit(1 if lib.failed(outcome) else 0)


if __name__ == "__main__":
    main()

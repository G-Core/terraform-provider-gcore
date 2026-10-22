#!/usr/bin/env python3
"""Plan an acceptance-test run as a GitHub Actions matrix.

    plan.py --schedule nightly                      # nightly, and the promotion gate
                                                    # (run on the SHA being promoted)
    plan.py --schedule weekly                       # the GPU suite
    plan.py --changed-files changed.txt --smoke     # pre-merge: packages the PR changed,
                                                    # plus the smoke set
    plan.py --packages cloud_volume,dns --smoke     # pre-merge by hand: packages or
                                                    # whole products, plus the smoke set

--changed-files reads one path per line ("-" for stdin), e.g. from
git diff --name-only BASE...HEAD. A change to shared code expands the run to the
full nightly scope.

Writes {"include": [job, ...]} to stdout and a human report to stderr. Exits
non-zero when the inventory is not accounted for exactly, when a package lacks
cleanup coverage without an explicit exclusion, or when the selection resolves
to nothing runnable.

The plan is deterministic for a given commit, so the cancellation cleanup can
rebuild it from the finished run's head SHA and recover every job's ownership
identity.
"""

import argparse
import json
import os
import sys

sys.dont_write_bytecode = True  # keep __pycache__ out of .github
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import acctest_lib as lib  # noqa: E402


def main(argv=None):
    ap = argparse.ArgumentParser()
    ap.add_argument("--repo-root", default=".")
    ap.add_argument("--schedule", choices=lib.SCHEDULES)
    ap.add_argument("--packages", default="",
                    help="comma-separated package or product names")
    ap.add_argument("--changed-files", metavar="FILE",
                    help="file listing changed paths, one per line, or - for stdin")
    ap.add_argument("--smoke", action="store_true", help="add the smoke set")
    args = ap.parse_args(argv)

    selection = [p.strip() for p in args.packages.split(",") if p.strip()]
    if args.schedule and args.changed_files:
        ap.error("--changed-files picks its own scope; do not combine it with --schedule")
    if not (args.schedule or selection or args.changed_files or args.smoke):
        ap.error("give --schedule, --packages, --changed-files or --smoke")

    root = os.path.abspath(args.repo_root)
    config = lib.load_config(root)
    packages, registered = lib.discover(root)

    schedule = args.schedule
    if args.changed_files:
        if args.changed_files == "-":
            files = sys.stdin.read().splitlines()
        else:
            with open(args.changed_files, encoding="utf-8") as fh:
                files = fh.read().splitlines()
        changed, broad, notes = lib.select_changed(files, packages)
        selection += changed
        if broad:
            schedule = "nightly"
            print("changed: shared code (%s) -> full nightly scope" % ", ".join(broad[:3]
                  + (["..."] if len(broad) > 3 else [])), file=sys.stderr)
        else:
            print("changed: %s" % (", ".join(changed) or "no tested package"), file=sys.stderr)
        for note in notes:
            print("  ignored %s" % note, file=sys.stderr)

    try:
        jobs, report = lib.build_plan(config, packages, registered,
                                      schedule=schedule, selection=selection,
                                      smoke=args.smoke)
    except lib.PolicyError as e:
        sys.exit("plan: %s" % e)

    total = sum(len(i["packages"]) for j in jobs for i in j["invocations"])
    print("plan: %d job(s), %d package(s)" % (len(jobs), total), file=sys.stderr)
    for j in jobs:
        print("  %-10s %-13s -p %d -parallel %d  %dm  sweepers=%d"
              % (j["name"], j["target"], j["go_test_p"], j["go_test_parallel"],
                 j["timeout_minutes"], len(j["sweepers"])), file=sys.stderr)
        for inv in j["invocations"]:
            print("      %s%s" % (", ".join(inv["packages"]),
                                  "  (-run %s)" % inv["run_regex"] if inv["run_regex"] else ""),
                  file=sys.stderr)
    for pkg, reason in report["excluded"].items():
        print("  excluded %-36s %s" % (pkg, reason), file=sys.stderr)
    for pkg in report["requested_but_excluded"]:
        print("  NOT RUN  %-36s requested, but excluded" % pkg, file=sys.stderr)

    json.dump({"include": jobs}, sys.stdout)
    print()


if __name__ == "__main__":
    main()

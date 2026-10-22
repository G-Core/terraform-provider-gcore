"""Shared logic for the acceptance-test CI policy.

Used by plan.py, validate-policy.py and check-results.py. Standard library only,
no credentials and no network, so every function here can run on any PR and in
offline unit tests.

The repository is the source of truth for what exists: packages, their tests,
the resource types their test HCL declares and the registered sweepers are all
discovered on every run. The policy files only record decisions the code cannot
express: which product runs where and when, explicit exclusions, cleanup
exceptions and accepted failures.

Extraction is deliberately shallow - a few regular expressions over test files,
not a Go or HCL parser. Anything the expressions cannot see goes in the small
reviewed exception map in cleanup-exceptions.json.
"""

import json
import os
import re

TESTACC = re.compile(r"^func (TestAcc[A-Za-z0-9_]+)\(", re.M)
# The -run filter for a package planned in full: its acceptance tests only, not
# the unit tests that share the package.
FULL_RUN_REGEX = "^TestAcc"
SWEEPER = re.compile(r'AddTestSweepers\(\s*"([a-z0-9_]+)"')
# Matches both raw-string HCL (resource "gcore_x") and escaped HCL in an
# interpreted string (resource \"gcore_x\").
HCL_RESOURCE = re.compile(r'\bresource\s+\\?"(gcore_[a-z0-9_]+)\\?"')
# A fixture created through the SDK rather than through Terraform, such as
# client.Cloud.Registries.New( or client.Storage.SftpStorages.NewAndPoll(.
SDK_CREATE = re.compile(r"\bclient\.[A-Z]\w*(?:\.[A-Z]\w*)+\.New\w*\(")
TICKET = re.compile(r"^[A-Z][A-Z0-9]+-[0-9]+$")
TEST_ID = re.compile(r"^[a-z0-9_]+\.TestAcc[A-Za-z0-9_]+$")
RESOURCE_TYPE = re.compile(r"^gcore_[a-z0-9_]+$")

SCHEDULES = ("nightly", "weekly")


class PolicyError(Exception):
    """A policy or inventory problem that makes a plan impossible to build."""


def strip_meta(d):
    return {k: v for k, v in d.items() if not k.startswith("$")}


def load_json(path):
    with open(path, encoding="utf-8") as fh:
        return json.load(fh)


# ---- discovery --------------------------------------------------------------

class Package:
    def __init__(self, name):
        self.name = name
        self.tests = set()
        self.hcl_types = set()
        self.sdk_create_files = []


def discover(root):
    """Return ({package: Package}, {registered sweeper names}).

    A package is any directory under internal/services with at least one
    TestAcc function in any *_test.go file, not just the two canonical file
    names: several packages keep their only acceptance test in
    list_data_source_test.go.
    """
    services = os.path.join(root, "internal", "services")
    packages = {}
    registered = set()
    for name in sorted(os.listdir(services)):
        pkg_dir = os.path.join(services, name)
        if not os.path.isdir(pkg_dir):
            continue
        pkg = Package(name)
        for entry in sorted(os.listdir(pkg_dir)):
            if not entry.endswith(".go"):
                continue
            with open(os.path.join(pkg_dir, entry), encoding="utf-8") as fh:
                body = fh.read()
            if not entry.endswith("_test.go"):
                registered.update(SWEEPER.findall(body))
                continue
            pkg.tests.update(TESTACC.findall(body))
            pkg.hcl_types.update(HCL_RESOURCE.findall(body))
            if SDK_CREATE.search(body):
                pkg.sdk_create_files.append(entry)
        if pkg.tests:
            packages[name] = pkg
    return packages, registered


def load_config(root):
    d = os.path.join(root, ".github", "acctest")
    return {
        "policy": load_json(os.path.join(d, "policy.json")),
        "exceptions": load_json(os.path.join(d, "cleanup-exceptions.json")),
        "known": load_json(os.path.join(d, "known-failures.json")),
    }


# ---- cleanup coverage -------------------------------------------------------

def resolve_sweeper(rtype, registered, exceptions, _seen=None):
    """The registered sweeper that reclaims rtype, or None.

    A type's own sweeper wins. Otherwise follow reclaimed_by_parent until a
    type with a registered sweeper is reached. sweeper_names maps a type to a
    registration name that differs from the type name.
    """
    seen = _seen or set()
    if rtype in seen:
        return None
    seen.add(rtype)
    name = strip_meta(exceptions.get("sweeper_names", {})).get(rtype, rtype)
    if name in registered:
        return name
    parent = strip_meta(exceptions.get("reclaimed_by_parent", {})).get(rtype)
    if parent:
        return resolve_sweeper(parent, registered, exceptions, seen)
    return None


def creates(pkg, exceptions):
    """Every resource type a package's tests create: HCL plus declared fixtures."""
    fixture = strip_meta(exceptions.get("fixtures", {})).get(pkg.name, {})
    if not isinstance(fixture, dict):
        fixture = {}  # malformed; validate() reports it
    return set(pkg.hcl_types) | set(fixture.get("creates") or [])


def coverage(pkg, registered, exceptions):
    """Return (sweepers, missing): the sweepers that reclaim what the package
    creates, and the types nothing reclaims. Covered means missing is empty."""
    sweepers, missing = set(), set()
    for rtype in creates(pkg, exceptions):
        s = resolve_sweeper(rtype, registered, exceptions)
        if s:
            sweepers.add(s)
        else:
            missing.add(rtype)
    return sorted(sweepers), sorted(missing)


# ---- products, targets, groups ----------------------------------------------

def product_of(package_name):
    """The product is the directory prefix: cloud_volume -> cloud."""
    return package_name.split("_", 1)[0]


def group_of(package_name, groups):
    """The group with the longest matching prefix within the package's product,
    else that product's default group (the one with no match list)."""
    product = product_of(package_name)
    best, best_len, default = None, -1, None
    for gname, g in groups.items():
        if g["product"] != product:
            continue
        prefixes = g.get("match")
        if not prefixes:
            default = gname
            continue
        for prefix in prefixes:
            if package_name.startswith(prefix) and len(prefix) > best_len:
                best, best_len = gname, len(prefix)
    return best or default


# ---- plan -------------------------------------------------------------------

# Files outside internal/ that still change what every resource does.
SHARED_ROOT_FILES = ("go.mod", "go.sum", "go.work", "main.go")


def select_changed(files, packages):
    """Map a PR's changed files to a pre-merge selection.

    Returns (selection, broad, notes):
      selection - packages with acceptance tests whose directory changed;
      broad     - the files that are shared code (anything else under internal/,
                  or a root Go file): any resource may be affected, so the caller
                  runs the full nightly scope instead of a selection;
      notes     - files that select nothing, with the reason.
    """
    selection, broad, notes = set(), [], []
    for f in files:
        f = f.strip()
        if not f:
            continue
        parts = f.split("/")
        if parts[:2] == ["internal", "services"] and len(parts) > 3:
            if parts[2] in packages:
                selection.add(parts[2])
            else:
                notes.append("%s: package %s has no acceptance tests" % (f, parts[2]))
        elif parts[0] == "internal" or f in SHARED_ROOT_FILES:
            broad.append(f)
        else:
            notes.append("%s: not provider code" % f)
    return sorted(selection), broad, notes


def build_plan(config, packages, registered, schedule=None, selection=None, smoke=False):
    """Build the job list for one run.

    schedule  - "nightly" or "weekly": every eligible package in groups on it.
    selection - pre-merge: package names, or product names meaning the whole
                product. Runs regardless of schedule.
    smoke     - add the smoke set (anchored regex) to the selection.

    Returns (jobs, report). Raises PolicyError when the inventory is not
    accounted for exactly, or the selection resolves to nothing runnable.
    """
    policy = config["policy"]
    exceptions = config["exceptions"]
    groups = strip_meta(policy["groups"])
    targets = strip_meta(policy["targets"])
    excluded = strip_meta(policy.get("excluded_packages", {}))

    # Every discovered package lands in exactly one group, or is excluded.
    # Computed over the whole inventory regardless of the selection, so a
    # package dropped from one schedule's plan cannot hide.
    assignment, unaccounted = {}, []
    for name, pkg in packages.items():
        if name in excluded:
            continue
        _, missing = coverage(pkg, registered, exceptions)
        if missing:
            # The cleanup gate. A package whose resources nothing reclaims must
            # never reach an executable job, whatever the policy says.
            raise PolicyError(
                "package %r creates %s, which no registered sweeper reclaims; add a "
                "sweeper or exclude it explicitly as missing-cleanup"
                % (name, ", ".join(missing)))
        g = group_of(name, groups)
        if g is None:
            unaccounted.append(name)
        else:
            assignment[name] = g
    if unaccounted:
        raise PolicyError("no execution group for: %s - add the product to "
                          "policy.groups or exclude the package" % ", ".join(sorted(unaccounted)))
    if set(assignment) | (set(excluded) & set(packages)) != set(packages):
        raise PolicyError("plan does not account for the whole inventory")

    # Which packages run, and with which -run filter.
    full, smoke_only = set(), set()
    if schedule:
        full = {p for p, g in assignment.items() if groups[g]["schedule"] == schedule}
    requested_excluded = []
    for item in selection or []:
        if item in packages:
            if item in excluded:
                requested_excluded.append(item)
            else:
                full.add(item)
        elif any(product_of(p) == item for p in packages):
            full |= {p for p in assignment if product_of(p) == item}
        else:
            raise PolicyError("selection %r is neither a package with acceptance tests "
                              "nor a product" % item)
    if smoke:
        smoke_only = {p for p in policy["smoke"]["packages"] if p not in full}
    if not full and not smoke_only:
        raise PolicyError("the selection resolved to no runnable package - an empty "
                          "run would report green having tested nothing")

    jobs = []
    for gname in sorted(groups):
        g = groups[gname]
        members_full = sorted(p for p in full if assignment.get(p) == gname)
        members_smoke = sorted(p for p in smoke_only if assignment.get(p) == gname)
        if not members_full and not members_smoke:
            continue
        invocations = []
        if members_full:
            invocations.append({"packages": members_full, "run_regex": FULL_RUN_REGEX})
        if members_smoke:
            invocations.append({"packages": members_smoke,
                                "run_regex": policy["smoke"]["run_regex"]})
        members = members_full + members_smoke
        sweepers = sorted({s for p in members
                           for s in coverage(packages[p], registered, exceptions)[0]})
        target = targets[g["target"]]
        jobs.append({
            "name": gname,
            "product": g["product"],
            "target": g["target"],
            # Combined by the workflow with the run ID and attempt into the
            # job's ownership identity, so each job sweeps only what it created.
            "ownership_key": "%s-%s" % (g["product"], g["target"]),
            "environment": target["environment"],
            "region_id": target.get("region_id"),
            "invocations": invocations,
            "sweepers": sweepers,
            "go_test_p": g["go_test_p"],
            "go_test_parallel": g["go_test_parallel"],
            "timeout_minutes": g["timeout_minutes"],
        })

    emitted = [p for j in jobs for inv in j["invocations"] for p in inv["packages"]]
    if len(emitted) != len(set(emitted)):
        dupes = sorted({p for p in emitted if emitted.count(p) > 1})
        raise PolicyError("packages planned more than once: %s" % ", ".join(dupes))
    if set(emitted) != full | smoke_only:
        raise PolicyError("planned packages differ from the selection: missing %s, "
                          "unexpected %s" % (sorted((full | smoke_only) - set(emitted)),
                                             sorted(set(emitted) - (full | smoke_only))))

    report = {
        "excluded": {p: excluded[p] for p in sorted(excluded) if p in packages},
        "requested_but_excluded": sorted(requested_excluded),
    }
    return jobs, report


# ---- validation -------------------------------------------------------------

def validate(root):
    """Return (problems, warnings, notes). Problems fail CI; warnings do not."""
    problems, warnings, notes = [], [], []
    config = load_config(root)
    packages, registered = discover(root)
    if not packages:
        return (["discovered no packages containing TestAcc functions - a check "
                 "that silently finds nothing is itself a false green"], [], [])

    policy, exceptions = config["policy"], config["exceptions"]
    targets = strip_meta(policy["targets"])
    groups = strip_meta(policy["groups"])
    excluded = strip_meta(policy.get("excluded_packages", {}))

    # ---- groups resolve ----
    owners = {}
    for gname, g in groups.items():
        for key in ("product", "target", "schedule", "go_test_p", "go_test_parallel",
                    "timeout_minutes"):
            if key not in g:
                problems.append("group %r is missing %r" % (gname, key))
        if g.get("target") not in targets:
            problems.append("group %r names unknown target %r" % (gname, g.get("target")))
        if g.get("schedule") not in SCHEDULES:
            problems.append("group %r has schedule %r, want one of %s"
                            % (gname, g.get("schedule"), ", ".join(SCHEDULES)))
        for prefix in g.get("match") or []:
            if not any(p.startswith(prefix) for p in packages):
                problems.append("group %r match %r matches no package - a rename would "
                                "silently move those tests to another schedule" % (gname, prefix))
        # Jobs sweep by ownership identity, so two jobs must never share one.
        owner = (g.get("product"), g.get("target"))
        if owner in owners:
            problems.append("groups %r and %r share product and target" % (owners[owner], gname))
        owners[owner] = gname
    for product in {g.get("product") for g in groups.values()}:
        defaults = [n for n, g in groups.items() if g.get("product") == product and not g.get("match")]
        if len(defaults) != 1:
            problems.append("product %r needs exactly one default group (no match list), "
                            "has %d" % (product, len(defaults)))

    # ---- exclusions resolve ----
    for pkg, reason in excluded.items():
        if pkg not in packages:
            problems.append("policy.excluded_packages names %r, which has no acceptance tests" % pkg)
        elif not reason:
            problems.append("exclusion %r needs a reason" % pkg)

    # ---- cleanup exceptions are well formed ----
    # Coverage reads these entries leniently - a missing or misspelled
    # "creates" means "creates nothing" - so validate their shape first.
    services = os.path.join(root, "internal", "services")

    def resource_type_problem(rtype):
        """Why rtype is not a provider resource type, or None."""
        if not isinstance(rtype, str) or not RESOURCE_TYPE.match(rtype):
            return "is not a gcore_* resource type"
        if not os.path.isdir(os.path.join(services, rtype[len("gcore_"):])):
            return "has no internal/services package"
        return None

    fixtures = strip_meta(exceptions.get("fixtures", {}))
    for name, fixture in fixtures.items():
        if name not in packages:
            problems.append("fixtures names %r, which has no acceptance tests" % name)
        if not isinstance(fixture, dict):
            problems.append("fixtures entry %r must be an object" % name)
            continue
        unknown = sorted(set(fixture) - {"creates", "reason"})
        if unknown:
            problems.append("fixtures entry %r has unknown keys %s - want creates and reason"
                            % (name, ", ".join(unknown)))
        created = fixture.get("creates")
        if not isinstance(created, list) or not created:
            problems.append("fixtures entry %r needs a non-empty creates list" % name)
        else:
            for rtype in created:
                why = resource_type_problem(rtype)
                if why:
                    problems.append("fixtures entry %r creates %r, which %s" % (name, rtype, why))
        if not isinstance(fixture.get("reason"), str) or not fixture["reason"].strip():
            problems.append("fixtures entry %r needs a reason" % name)
    all_created = {t for pkg in packages.values() for t in creates(pkg, exceptions)}
    for child, parent in strip_meta(exceptions.get("reclaimed_by_parent", {})).items():
        bad = False
        for rtype in (child, parent):
            why = resource_type_problem(rtype)
            if why:
                bad = True
                problems.append("reclaimed_by_parent entry %r -> %r: %r %s"
                                % (child, parent, rtype, why))
        if not bad and child not in all_created:
            warnings.append("reclaimed_by_parent names %r, which no package creates - remove "
                            "the stale entry" % child)

    # ---- the cleanup gate ----
    for name, pkg in packages.items():
        _, missing = coverage(pkg, registered, exceptions)
        if missing and name not in excluded:
            problems.append("package %r creates %s, which no registered sweeper reclaims - "
                            "add a sweeper, declare a parent in cleanup-exceptions.json, or "
                            "exclude the package" % (name, ", ".join(missing)))
        if pkg.sdk_create_files and name not in fixtures:
            problems.append("package %r creates fixtures through the SDK (%s) - declare what "
                            "they create in cleanup-exceptions.json fixtures"
                            % (name, ", ".join(pkg.sdk_create_files)))
    for rtype, name in strip_meta(exceptions.get("sweeper_names", {})).items():
        if name not in registered:
            problems.append("sweeper_names maps %r to %r, which is not a registered sweeper"
                            % (rtype, name))

    # ---- smoke resolves ----
    rx = re.compile(policy["smoke"]["run_regex"])
    if not policy["smoke"]["packages"]:
        problems.append("policy.smoke.packages is empty - pre-merge runs would have no "
                        "smoke coverage")
    for pkg in policy["smoke"]["packages"]:
        if pkg not in packages or pkg in excluded:
            problems.append("policy.smoke names %r, which has no runnable acceptance tests" % pkg)
        elif not any(rx.search(t) for t in packages[pkg].tests):
            problems.append("policy.smoke names %r, but none of its tests match %r"
                            % (pkg, policy["smoke"]["run_regex"]))

    # ---- baseline ----
    all_tests = {"%s.%s" % (p, t) for p, pkg in packages.items() for t in pkg.tests}
    seen = {}
    for group in config["known"].get("groups", []):
        gid = group.get("id", "<no id>")
        for field in ("id", "jira", "reason", "tests"):
            if not group.get(field):
                problems.append("known-failures group %r is missing %r" % (gid, field))
        if group.get("jira") and not TICKET.match(group["jira"]):
            problems.append("known-failures group %r has malformed ticket %r" % (gid, group["jira"]))
        for ref in group.get("tests", []):
            if not TEST_ID.match(ref):
                problems.append("known-failures group %r entry %r is not package-qualified "
                                "(want package.TestAccName)" % (gid, ref))
            elif ref not in all_tests:
                warnings.append("known-failures group %r references %r, which no longer "
                                "exists - remove the stale entry" % (gid, ref))
            if ref in seen:
                problems.append("%r is listed in known-failures groups %r and %r"
                                % (ref, seen[ref], gid))
            seen[ref] = gid

    # ---- every schedule plans cleanly ----
    if not problems:
        for schedule in SCHEDULES:
            try:
                build_plan(config, packages, registered, schedule=schedule)
            except PolicyError as e:
                problems.append("%s plan: %s" % (schedule, e))

    notes.append("inventory: %d packages, %d tests, %d registered sweepers"
                 % (len(packages), len(all_tests), len(registered)))
    notes.append("excluded: %d packages" % len(excluded))
    notes.append("baseline: %d known-failure groups covering %d tests"
                 % (len(config["known"].get("groups", [])), len(seen)))
    return problems, warnings, notes


# ---- results ----------------------------------------------------------------

def expected_tests(job, packages):
    """Test IDs a job must report: every test its invocations select."""
    out = set()
    for inv in job["invocations"]:
        rx = re.compile(inv["run_regex"]) if inv.get("run_regex") else None
        for p in inv["packages"]:
            for t in packages[p].tests:
                tid = "%s.%s" % (p, t)
                if rx is None or rx.search(t):
                    out.add(tid)
    return out


def parse_go_test_json(lines):
    """Return ({test_id: final action}, {package: failed at package level}).

    Subtests are ignored: the top-level test carries the verdict. A test with
    no final pass, fail or skip event - killed by a timeout or a crash - simply
    has no entry, so it surfaces as missing.
    """
    results, package_failed = {}, set()
    for line in lines:
        line = line.strip()
        if not line.startswith("{"):
            continue
        try:
            ev = json.loads(line)
        except ValueError:
            continue
        action = ev.get("Action")
        if action not in ("pass", "fail", "skip"):
            continue
        pkg = (ev.get("Package") or "").rsplit("/", 1)[-1]
        test = ev.get("Test")
        if not test:
            if action == "fail":
                package_failed.add(pkg)
            continue
        if "/" in test:
            continue
        results["%s.%s" % (pkg, test)] = action
    return results, package_failed


def accepted_failures(config):
    """{test_id: group id} for every test in known-failures.json."""
    return {ref: g["id"] for g in config["known"].get("groups", []) for ref in g.get("tests", [])}


def classify(job, packages, results, package_failed, accepted):
    """Sort a job's results into the buckets the report shows. The job fails
    when unexpected_failures, missing, unexpected_skips or package_errors is
    non-empty; accepted and now_passing are reported but do not fail."""
    expected = expected_tests(job, packages)
    # Only an expected test explains a package failure. A failing test outside
    # the expected set - a unit test, or one discovery missed - is unexpected
    # itself and must not mask the package-level failure.
    failed_pkgs = {tid.split(".", 1)[0] for tid, a in results.items()
                   if a == "fail" and tid in expected}
    out = {
        "unexpected_failures": [], "accepted": [], "now_passing": [],
        "unexpected_skips": [], "missing": [], "package_errors": [], "passed": 0,
    }
    for tid in sorted(expected):
        action = results.get(tid)
        if action is None:
            out["missing"].append(tid)
        elif action == "skip":
            out["unexpected_skips"].append(tid)
        elif action == "fail":
            if tid in accepted:
                out["accepted"].append((tid, accepted[tid]))
            else:
                out["unexpected_failures"].append(tid)
        else:
            out["passed"] += 1
            if tid in accepted:
                out["now_passing"].append((tid, accepted[tid]))
    for tid in sorted(results):
        if results[tid] == "fail" and tid not in expected:
            out["unexpected_failures"].append(tid)
    planned ={p for inv in job["invocations"] for p in inv["packages"]}
    # A package-level failure with no failing test behind it is a build error,
    # a panic outside a test or an init failure: nothing ran that could explain it.
    for pkg in sorted((package_failed & planned) - failed_pkgs):
        out["package_errors"].append(pkg)
    return out


def failed(outcome):
    return any(outcome[k] for k in
               ("unexpected_failures", "missing", "unexpected_skips", "package_errors"))

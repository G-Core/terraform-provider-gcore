"""Offline tests for the acceptance-test policy logic.

Each test builds a minimal fake repository in a temporary directory, so the
rules are exercised in isolation from the real inventory. The final class runs
the same coverage invariant against the real repository.

Run: python3 -B -m unittest discover -s .github/acctest/scripts/tests
"""

import copy
import json
import os
import sys
import tempfile
import unittest

sys.dont_write_bytecode = True  # keep __pycache__ out of .github
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
import acctest_lib as lib  # noqa: E402

POLICY = {
    "targets": {
        "lux": {"environment": "e", "region_id": 76},
        "sines": {"environment": "e", "region_id": 176},
        "account": {"environment": "e"},
    },
    "groups": {
        "cloud": {"product": "cloud", "target": "lux", "schedule": "nightly",
                  "go_test_p": 2, "go_test_parallel": 1, "timeout_minutes": 60},
        "cloud-gpu": {"product": "cloud", "match": ["cloud_gpu_"], "target": "sines",
                      "schedule": "weekly", "go_test_p": 1, "go_test_parallel": 1,
                      "timeout_minutes": 60},
        "dns": {"product": "dns", "target": "account", "schedule": "nightly",
                "go_test_p": 2, "go_test_parallel": 2, "timeout_minutes": 30},
    },
    "excluded_packages": {},
    "smoke": {"run_regex": "^TestAcc.*_basic$", "packages": ["cloud_volume"]},
}


def hcl_test(funcs, types=(), extra=""):
    body = "package x\n\n"
    for f in funcs:
        body += "func %s(t *testing.T) {}\n" % f
    for i, t in enumerate(types):
        body += 'const c%d = `resource "%s" "test" {}`\n' % (i, t)
    return body + extra


class Repo:
    """A throwaway repository with services, sweepers and policy files."""

    def __init__(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.root = self.tmp.name
        self.policy = copy.deepcopy(POLICY)
        self.exceptions = {"fixtures": {}, "reclaimed_by_parent": {}, "sweeper_names": {}}
        self.known = {"groups": []}

    def package(self, name, tests, types=(), sweepers=(), extra=""):
        d = os.path.join(self.root, "internal", "services", name)
        os.makedirs(d, exist_ok=True)
        with open(os.path.join(d, "resource_test.go"), "w") as fh:
            fh.write(hcl_test(tests, types, extra))
        if sweepers:
            with open(os.path.join(d, "sweep.go"), "w") as fh:
                for s in sweepers:
                    fh.write('func init() { resource.AddTestSweepers("%s", nil) }\n' % s)
        return self

    def write(self):
        d = os.path.join(self.root, ".github", "acctest")
        os.makedirs(d, exist_ok=True)
        for name, data in (("policy.json", self.policy),
                           ("cleanup-exceptions.json", self.exceptions),
                           ("known-failures.json", self.known)):
            with open(os.path.join(d, name), "w") as fh:
                json.dump(data, fh)
        return self

    def validate(self):
        self.write()
        return lib.validate(self.root)

    def plan(self, **kw):
        self.write()
        packages, registered = lib.discover(self.root)
        return lib.build_plan(lib.load_config(self.root), packages, registered, **kw)

    def close(self):
        self.tmp.cleanup()


def base_repo():
    """Three covered packages: one per group."""
    return (Repo()
            .package("cloud_volume", ["TestAccCloudVolume_basic", "TestAccCloudVolume_update"],
                     ["gcore_cloud_volume"], ["gcore_cloud_volume"])
            .package("cloud_gpu_virtual_cluster", ["TestAccGPU_basic"],
                     ["gcore_cloud_gpu_virtual_cluster"], ["gcore_cloud_gpu_virtual_cluster"])
            .package("dns_zone", ["TestAccDNSZone_basic"], ["gcore_dns_zone"], ["gcore_dns_zone"]))


def planned(jobs):
    return [p for j in jobs for inv in j["invocations"] for p in inv["packages"]]


class RepoTestCase(unittest.TestCase):
    def setUp(self):
        self.repo = base_repo()

    def tearDown(self):
        self.repo.close()

    def assertValid(self):
        problems, warnings, _ = self.repo.validate()
        self.assertEqual(problems, [])
        return warnings

    def assertProblem(self, fragment):
        problems, _, _ = self.repo.validate()
        self.assertTrue(any(fragment in p for p in problems),
                        "no problem containing %r in %r" % (fragment, problems))


class ProductDiscovery(RepoTestCase):
    def test_product_is_the_directory_prefix(self):
        self.assertEqual(lib.product_of("cloud_load_balancer_pool"), "cloud")
        self.assertEqual(lib.product_of("dns_zone"), "dns")

    def test_base_repo_is_valid(self):
        self.assertValid()

    def test_package_with_unknown_product_fails(self):
        self.repo.package("waap_domain", ["TestAccWaap_basic"])
        self.assertProblem("no execution group for: waap_domain")

    def test_directory_without_acceptance_tests_is_ignored(self):
        self.repo.package("cloud_quota", [])
        packages, _ = lib.discover(self.repo.root)
        self.assertNotIn("cloud_quota", packages)

    def test_product_needs_exactly_one_default_group(self):
        self.repo.policy["groups"]["dns-2"] = dict(self.repo.policy["groups"]["dns"], target="lux")
        self.assertProblem("product 'dns' needs exactly one default group")


class TargetAssignment(RepoTestCase):
    def test_longest_prefix_group_wins(self):
        groups = self.repo.policy["groups"]
        self.assertEqual(lib.group_of("cloud_gpu_virtual_cluster", groups), "cloud-gpu")
        self.assertEqual(lib.group_of("cloud_volume", groups), "cloud")
        self.assertEqual(lib.group_of("dns_zone", groups), "dns")

    def test_jobs_carry_target_coordinates(self):
        jobs, _ = self.repo.plan(schedule="nightly")
        by_name = {j["name"]: j for j in jobs}
        self.assertEqual(by_name["cloud"]["region_id"], 76)
        self.assertIsNone(by_name["dns"]["region_id"])

    def test_group_prefix_matching_nothing_fails(self):
        self.repo.policy["groups"]["cloud-gpu"]["match"] = ["cloud_gpux_"]
        self.assertProblem("matches no package")

    def test_ownership_keys_are_product_and_target(self):
        jobs, _ = self.repo.plan(schedule="nightly")
        keys = [j["ownership_key"] for j in jobs]
        self.assertEqual(sorted(keys), ["cloud-lux", "dns-account"])
        self.assertEqual(len(keys), len(set(keys)))

    def test_two_groups_sharing_product_and_target_fail(self):
        self.repo.policy["groups"]["cloud-gpu"]["target"] = "lux"
        self.assertProblem("share product and target")


class ScheduleSelection(RepoTestCase):
    def test_nightly_excludes_weekly_groups(self):
        jobs, _ = self.repo.plan(schedule="nightly")
        self.assertEqual(sorted(planned(jobs)), ["cloud_volume", "dns_zone"])

    def test_weekly_selects_only_gpu(self):
        jobs, _ = self.repo.plan(schedule="weekly")
        self.assertEqual(planned(jobs), ["cloud_gpu_virtual_cluster"])
        self.assertEqual(jobs[0]["target"], "sines")

    def test_unknown_schedule_fails(self):
        self.repo.policy["groups"]["dns"]["schedule"] = "hourly"
        self.assertProblem("has schedule 'hourly'")

    def test_empty_selection_fails(self):
        self.repo.policy["groups"]["cloud-gpu"]["schedule"] = "nightly"
        with self.assertRaisesRegex(lib.PolicyError, "no runnable package"):
            self.repo.plan(schedule="weekly")


class PreMergeSelection(RepoTestCase):
    def setUp(self):
        super().setUp()
        self.repo.package("cloud_network", ["TestAccCloudNetwork_basic"],
                          ["gcore_cloud_network"], ["gcore_cloud_network"])

    def test_affected_package_does_not_pull_in_its_whole_product(self):
        jobs, _ = self.repo.plan(selection=["cloud_network"])
        self.assertEqual(planned(jobs), ["cloud_network"])

    def test_product_name_selects_the_whole_product(self):
        jobs, _ = self.repo.plan(selection=["cloud"])
        self.assertEqual(sorted(planned(jobs)),
                         ["cloud_gpu_virtual_cluster", "cloud_network", "cloud_volume"])

    def test_smoke_runs_anchored_regex_and_affected_runs_in_full(self):
        jobs, _ = self.repo.plan(selection=["cloud_network"], smoke=True)
        invs = jobs[0]["invocations"]
        self.assertEqual(invs[0], {"packages": ["cloud_network"], "run_regex": "^TestAcc"})
        self.assertEqual(invs[1], {"packages": ["cloud_volume"], "run_regex": "^TestAcc.*_basic$"})

    def test_affected_smoke_package_runs_in_full_not_twice(self):
        jobs, _ = self.repo.plan(selection=["cloud_volume"], smoke=True)
        self.assertEqual(jobs[0]["invocations"],
                         [{"packages": ["cloud_volume"], "run_regex": "^TestAcc"}])

    def test_empty_smoke_set_fails(self):
        self.repo.policy["smoke"]["packages"] = []
        self.assertProblem("policy.smoke.packages is empty")

    def test_unknown_selection_fails(self):
        with self.assertRaisesRegex(lib.PolicyError, "neither a package"):
            self.repo.plan(selection=["cloud_nope"])


class Exclusions(RepoTestCase):
    def test_excluded_package_never_planned_and_is_reported(self):
        self.repo.policy["excluded_packages"]["dns_zone"] = "r"
        jobs, report = self.repo.plan(schedule="nightly", selection=["dns_zone"])
        self.assertNotIn("dns_zone", planned(jobs))
        self.assertIn("dns_zone", report["excluded"])
        self.assertEqual(report["requested_but_excluded"], ["dns_zone"])

    def test_exclusion_of_unknown_package_fails(self):
        self.repo.policy["excluded_packages"]["dns_nope"] = "r"
        self.assertProblem("'dns_nope', which has no acceptance tests")

    def test_smoke_package_cannot_be_excluded(self):
        self.repo.policy["excluded_packages"]["cloud_volume"] = "r"
        self.assertProblem("policy.smoke names 'cloud_volume', which has no runnable")


class CleanupCoverage(RepoTestCase):
    def test_package_without_cleanup_fails_validation(self):
        self.repo.package("cloud_registry", ["TestAccRegistry_basic"], ["gcore_cloud_registry"])
        self.assertProblem("creates gcore_cloud_registry, which no registered sweeper reclaims")

    def test_package_without_cleanup_cannot_enter_a_plan(self):
        # The planner enforces the gate itself, so a bypassed validator still
        # cannot put an unswept package into an executable job.
        self.repo.package("cloud_registry", ["TestAccRegistry_basic"], ["gcore_cloud_registry"])
        for kw in ({"schedule": "nightly"}, {"selection": ["dns_zone"]}):
            with self.assertRaisesRegex(lib.PolicyError, "no registered sweeper"):
                self.repo.plan(**kw)

    def test_explicit_exclusion_admits_the_rest_of_the_plan(self):
        self.repo.package("cloud_registry", ["TestAccRegistry_basic"], ["gcore_cloud_registry"])
        self.repo.policy["excluded_packages"]["cloud_registry"] = "no sweeper"
        self.assertValid()
        jobs, _ = self.repo.plan(schedule="nightly")
        self.assertNotIn("cloud_registry", planned(jobs))

    def test_coverage_uses_registered_sweepers_not_sweep_go(self):
        # A sweep.go registering a different type does not cover this one.
        self.repo.package("cloud_secret", ["TestAccSecret_basic"], ["gcore_cloud_secret"],
                          ["gcore_cloud_something_else"])
        self.assertProblem("creates gcore_cloud_secret")

    def test_sweeper_registered_in_another_package_covers(self):
        self.repo.package("cloud_volume_attachment", ["TestAccAttach_basic"], ["gcore_cloud_volume"])
        self.assertValid()

    def test_parent_sweeper_exception_covers_child(self):
        self.repo.package("cloud_load_balancer", ["TestAccLB_basic"], ["gcore_cloud_load_balancer"],
                          ["gcore_cloud_load_balancer"])
        self.repo.package("cloud_load_balancer_pool_member", ["TestAccMember_basic"],
                          ["gcore_cloud_load_balancer", "gcore_cloud_load_balancer_pool_member"])
        self.assertProblem("gcore_cloud_load_balancer_pool_member")
        self.repo.package("cloud_load_balancer_pool", [])
        # Chains resolve: member -> pool -> load balancer.
        self.repo.exceptions["reclaimed_by_parent"] = {
            "gcore_cloud_load_balancer_pool_member": "gcore_cloud_load_balancer_pool",
            "gcore_cloud_load_balancer_pool": "gcore_cloud_load_balancer",
        }
        self.assertValid()
        jobs, _ = self.repo.plan(selection=["cloud_load_balancer_pool_member"])
        self.assertEqual(jobs[0]["sweepers"], ["gcore_cloud_load_balancer"])

    def test_parent_cycle_does_not_loop(self):
        self.repo.exceptions["reclaimed_by_parent"] = {
            "gcore_a": "gcore_b", "gcore_b": "gcore_a"}
        self.assertIsNone(lib.resolve_sweeper("gcore_a", set(), self.repo.exceptions))

    def test_undeclared_sdk_fixture_fails(self):
        self.repo.package("cloud_registry", ["TestAccRegistry_basic"], [],
                          extra="var _, _ = client.Cloud.Registries.New(ctx, p)\n")
        self.assertProblem("creates fixtures through the SDK")

    def test_declared_sdk_fixture_counts_toward_coverage(self):
        self.repo.package("cloud_registry", ["TestAccRegistry_basic"], [],
                          extra="var _, _ = client.Cloud.Registries.New(ctx, p)\n")
        self.repo.exceptions["fixtures"]["cloud_registry"] = {
            "creates": ["gcore_cloud_registry"], "reason": "r"}
        self.assertProblem("creates gcore_cloud_registry, which no registered sweeper")
        self.repo.package("cloud_registry_sweeper", [], sweepers=["gcore_cloud_registry"])
        self.assertValid()

    def test_helper_fixture_counts_toward_coverage(self):
        self.repo.exceptions["fixtures"]["dns_zone"] = {
            "creates": ["gcore_dns_network_mapping"], "reason": "shared helper"}
        self.assertProblem("creates gcore_dns_network_mapping")

    def fixture(self, entry):
        self.repo.package("cloud_registry", ["TestAccRegistry_basic"], [],
                          extra="var _, _ = client.Cloud.Registries.New(ctx, p)\n")
        self.repo.exceptions["fixtures"]["cloud_registry"] = entry

    def test_fixture_with_empty_creates_fails(self):
        self.fixture({"creates": [], "reason": "r"})
        self.assertProblem("'cloud_registry' needs a non-empty creates list")

    def test_fixture_with_misspelled_key_fails(self):
        self.fixture({"create": ["gcore_cloud_registry"], "reason": "r"})
        self.assertProblem("'cloud_registry' has unknown keys create")
        self.assertProblem("'cloud_registry' needs a non-empty creates list")

    def test_fixture_without_reason_fails(self):
        self.fixture({"creates": ["gcore_cloud_registry"], "reason": " "})
        self.assertProblem("'cloud_registry' needs a reason")

    def test_fixture_for_unknown_package_fails(self):
        self.repo.exceptions["fixtures"]["cloud_nope"] = {
            "creates": ["gcore_cloud_volume"], "reason": "r"}
        self.assertProblem("fixtures names 'cloud_nope', which has no acceptance tests")

    def test_fixture_with_bad_type_fails(self):
        self.fixture({"creates": ["cloud_registry"], "reason": "r"})
        self.assertProblem("creates 'cloud_registry', which is not a gcore_* resource type")

    def test_fixture_with_misspelled_type_fails(self):
        self.fixture({"creates": ["gcore_cloud_registy"], "reason": "r"})
        self.assertProblem("'gcore_cloud_registy', which has no internal/services package")

    def test_fixture_that_is_not_an_object_fails(self):
        self.fixture(["gcore_cloud_registry"])
        self.assertProblem("fixtures entry 'cloud_registry' must be an object")

    def test_misspelled_parent_fails(self):
        self.repo.package("cloud_volume_snapshot", [])
        self.repo.exceptions["reclaimed_by_parent"] = {
            "gcore_cloud_volume_snapshot": "gcore_cloud_volme"}
        self.assertProblem("'gcore_cloud_volme' has no internal/services package")

    def test_parent_with_bad_type_fails(self):
        self.repo.exceptions["reclaimed_by_parent"] = {"gcore_cloud_volume": "cloud_volume"}
        self.assertProblem("'cloud_volume' is not a gcore_* resource type")

    def test_unused_parent_entry_is_a_warning(self):
        self.repo.package("cloud_volume_snapshot", [])
        self.repo.exceptions["reclaimed_by_parent"] = {
            "gcore_cloud_volume_snapshot": "gcore_cloud_volume"}
        warnings = self.assertValid()
        self.assertTrue(any("gcore_cloud_volume_snapshot" in w for w in warnings))

    def test_sweeper_name_alias(self):
        self.repo.package("cdn_rule", ["TestAccRule_basic"], ["gcore_cdn_rule"], ["gcore_cdn_resource_rule"])
        self.repo.policy["groups"]["cdn"] = dict(self.repo.policy["groups"]["dns"], product="cdn")
        self.assertProblem("creates gcore_cdn_rule")
        self.repo.exceptions["sweeper_names"] = {"gcore_cdn_rule": "gcore_cdn_resource_rule"}
        self.assertValid()

    def test_alias_to_unregistered_sweeper_fails(self):
        self.repo.exceptions["sweeper_names"] = {"gcore_dns_zone": "gcore_dns_nope"}
        self.assertProblem("not a registered sweeper")

    def test_escaped_hcl_is_discovered(self):
        self.repo.package("cloud_secret", ["TestAccSecret_basic"],
                          extra='const c = "resource \\"gcore_cloud_secret\\" \\"t\\" {}"\n')
        self.assertProblem("creates gcore_cloud_secret")


class Baseline(RepoTestCase):
    def group(self, **over):
        g = {"id": "g", "jira": "TICKET-1", "reason": "r",
             "tests": ["dns_zone.TestAccDNSZone_basic"]}
        g.update(over)
        self.repo.known["groups"] = [g]

    def test_valid_group(self):
        self.group()
        self.assertValid()

    def test_missing_jira_fails(self):
        self.group(jira="")
        self.assertProblem("is missing 'jira'")

    def test_malformed_ticket_fails(self):
        self.group(jira="see slack")
        self.assertProblem("malformed ticket")

    def test_deleted_test_is_a_warning_not_a_failure(self):
        self.group(tests=["dns_zone.TestAccDNSZone_deleted"])
        warnings = self.assertValid()
        self.assertTrue(any("TestAccDNSZone_deleted" in w for w in warnings))

    def test_unqualified_test_fails(self):
        self.group(tests=["TestAccDNSZone_basic"])
        self.assertProblem("is not package-qualified")

    def test_test_in_two_groups_fails(self):
        self.group()
        self.repo.known["groups"].append(dict(self.repo.known["groups"][0], id="h"))
        self.assertProblem("listed in known-failures groups")


class ChangedFiles(RepoTestCase):
    def select(self, *files):
        packages, _ = lib.discover(self.repo.root)
        return lib.select_changed(files, packages)

    def run_plan(self, files, *extra):
        import contextlib
        import io
        import plan
        self.repo.write()
        path = os.path.join(self.repo.root, "changed.txt")
        with open(path, "w") as fh:
            fh.write("\n".join(files) + "\n")
        out = io.StringIO()
        with contextlib.redirect_stdout(out), contextlib.redirect_stderr(io.StringIO()):
            plan.main(["--repo-root", self.repo.root, "--changed-files", path, *extra])
        return planned(json.loads(out.getvalue())["include"])

    def test_service_files_select_their_package(self):
        selection, broad, _ = self.select("internal/services/dns_zone/resource.go",
                                          "internal/services/dns_zone/resource_test.go")
        self.assertEqual((selection, broad), (["dns_zone"], []))

    def test_package_without_acceptance_tests_selects_nothing(self):
        self.repo.package("waap_domain", [])
        selection, broad, notes = self.select("internal/services/waap_domain/resource.go")
        self.assertEqual((selection, broad), ([], []))
        self.assertIn("has no acceptance tests", notes[0])

    def test_shared_code_is_broad(self):
        for f in ("internal/apijson/decoder.go", "internal/acctest/helpers.go",
                  "internal/services/services.go", "go.mod", "go.sum", "main.go"):
            _, broad, _ = self.select(f)
            self.assertEqual(broad, [f], f)

    def test_non_provider_files_are_ignored(self):
        selection, broad, notes = self.select("docs/resources/dns_zone.md", ".github/acctest/policy.json",
                                              "examples/x/main.tf", "README.md", "")
        self.assertEqual((selection, broad, len(notes)), ([], [], 4))

    def test_plan_runs_only_changed_packages(self):
        self.repo.package("cloud_network", ["TestAccCloudNetwork_basic"],
                          ["gcore_cloud_network"], ["gcore_cloud_network"])
        self.assertEqual(self.run_plan(["internal/services/cloud_network/schema.go"]),
                         ["cloud_network"])

    def test_plan_expands_shared_code_to_the_nightly_scope(self):
        got = self.run_plan(["internal/apijson/decoder.go", "internal/services/dns_zone/x.go"])
        self.assertEqual(sorted(got), ["cloud_volume", "dns_zone"])

    def test_plan_adds_smoke_to_changed_packages(self):
        got = self.run_plan(["internal/services/dns_zone/x.go"], "--smoke")
        self.assertEqual(sorted(got), ["cloud_volume", "dns_zone"])

    def test_docs_only_change_without_smoke_is_an_empty_plan(self):
        with self.assertRaises(SystemExit):
            self.run_plan(["docs/index.md"])


class InventoryAccounting(RepoTestCase):
    def assertAccounted(self, root):
        config = lib.load_config(root)
        packages, registered = lib.discover(root)
        excluded = set(lib.strip_meta(config["policy"]["excluded_packages"]))
        seen = []
        for schedule in lib.SCHEDULES:
            jobs, _ = lib.build_plan(config, packages, registered, schedule=schedule)
            seen += planned(jobs)
        self.assertEqual(len(seen), len(set(seen)), "a package is planned twice")
        self.assertFalse(set(seen) & excluded, "an excluded package is planned")
        self.assertEqual(set(seen) | excluded, set(packages), "a package is unaccounted for")

    def test_every_package_planned_once_or_excluded(self):
        self.repo.package("cloud_registry", ["TestAccRegistry_basic"], ["gcore_cloud_registry"])
        self.repo.policy["excluded_packages"]["cloud_registry"] = "r"
        self.repo.write()
        self.assertAccounted(self.repo.root)

    def test_real_repository(self):
        root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "..", "..", ".."))
        self.assertAccounted(root)


class Results(RepoTestCase):
    PKG = "github.com/G-Core/terraform-provider-gcore/internal/services/"

    def ev(self, action, pkg, test=None):
        e = {"Action": action, "Package": self.PKG + pkg}
        if test:
            e["Test"] = test
        return json.dumps(e)

    def outcome(self, lines, known=()):
        self.repo.known["groups"] = [{"id": "g", "jira": "X-1", "reason": "r", "tests": list(known)}]
        jobs, _ = self.repo.plan(schedule="nightly")
        job = [j for j in jobs if j["name"] == "cloud"][0]
        packages, _ = lib.discover(self.repo.root)
        results, package_failed = lib.parse_go_test_json(lines)
        accepted = lib.accepted_failures(lib.load_config(self.repo.root))
        return lib.classify(job, packages, results, package_failed, accepted)

    def test_all_pass(self):
        out = self.outcome([self.ev("pass", "cloud_volume", "TestAccCloudVolume_basic"),
                            self.ev("pass", "cloud_volume", "TestAccCloudVolume_update"),
                            self.ev("pass", "cloud_volume")])
        self.assertFalse(lib.failed(out))
        self.assertEqual(out["passed"], 2)

    def test_accepted_failure_reported_separately_and_does_not_fail(self):
        out = self.outcome([self.ev("fail", "cloud_volume", "TestAccCloudVolume_basic"),
                            self.ev("pass", "cloud_volume", "TestAccCloudVolume_update"),
                            self.ev("fail", "cloud_volume")],
                           known=["cloud_volume.TestAccCloudVolume_basic"])
        self.assertFalse(lib.failed(out))
        self.assertEqual(out["accepted"], [("cloud_volume.TestAccCloudVolume_basic", "g")])
        self.assertEqual(out["package_errors"], [])

    def test_renamed_test_gets_no_exemption(self):
        out = self.outcome([self.ev("fail", "cloud_volume", "TestAccCloudVolume_basic"),
                            self.ev("pass", "cloud_volume", "TestAccCloudVolume_update")],
                           known=["cloud_volume.TestAccCloudVolume_oldName"])
        self.assertEqual(out["unexpected_failures"], ["cloud_volume.TestAccCloudVolume_basic"])

    def test_missing_result_fails(self):
        out = self.outcome([self.ev("pass", "cloud_volume", "TestAccCloudVolume_basic")])
        self.assertEqual(out["missing"], ["cloud_volume.TestAccCloudVolume_update"])
        self.assertTrue(lib.failed(out))

    def test_skip_fails(self):
        out = self.outcome([self.ev("skip", "cloud_volume", "TestAccCloudVolume_basic"),
                            self.ev("pass", "cloud_volume", "TestAccCloudVolume_update")])
        self.assertEqual(out["unexpected_skips"], ["cloud_volume.TestAccCloudVolume_basic"])
        self.assertTrue(lib.failed(out))

    def test_package_failure_without_failing_test_fails(self):
        out = self.outcome([self.ev("pass", "cloud_volume", "TestAccCloudVolume_basic"),
                            self.ev("pass", "cloud_volume", "TestAccCloudVolume_update"),
                            self.ev("fail", "cloud_volume")])
        self.assertEqual(out["package_errors"], ["cloud_volume"])

    def test_failing_unit_test_does_not_explain_package_failure(self):
        out = self.outcome([self.ev("pass", "cloud_volume", "TestAccCloudVolume_basic"),
                            self.ev("pass", "cloud_volume", "TestAccCloudVolume_update"),
                            self.ev("fail", "cloud_volume", "TestVolumeSchema"),
                            self.ev("fail", "cloud_volume")])
        self.assertTrue(lib.failed(out))
        self.assertEqual(out["unexpected_failures"], ["cloud_volume.TestVolumeSchema"])
        self.assertEqual(out["package_errors"], ["cloud_volume"])

    def test_now_passing_is_reported(self):
        out = self.outcome([self.ev("pass", "cloud_volume", "TestAccCloudVolume_basic"),
                            self.ev("pass", "cloud_volume", "TestAccCloudVolume_update")],
                           known=["cloud_volume.TestAccCloudVolume_basic"])
        self.assertFalse(lib.failed(out))
        self.assertEqual(out["now_passing"], [("cloud_volume.TestAccCloudVolume_basic", "g")])

    def test_subtests_and_noise_are_ignored(self):
        out = self.outcome(["not json", "{bad",
                            self.ev("fail", "cloud_volume", "TestAccCloudVolume_basic/step"),
                            self.ev("pass", "cloud_volume", "TestAccCloudVolume_basic"),
                            self.ev("pass", "cloud_volume", "TestAccCloudVolume_update")])
        self.assertFalse(lib.failed(out))


if __name__ == "__main__":
    unittest.main()

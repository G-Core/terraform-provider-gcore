---
name: acctest
description: Generates and runs Terraform acceptance tests following HashiCorp best practices for Plugin Framework providers. Use when developing CRUD acceptance tests (TestAcc* functions) for resources and data sources, creating sweeper implementations for test resource cleanup, or debugging acceptance test failures. Triggers on tasks like "write acceptance tests", "add tests for resource X", "create sweeper", "fix failing acctest", or any work in resource_test.go / data_source_test.go / sweep.go files.
---

# Terraform Acceptance Tests

Generate, run, and debug acceptance tests for this Gcore Terraform provider (Plugin Framework, Protocol V6).

## Workflow

1. **Understand the resource**: Read the resource's `schema.go`, `model.go`, and `resource.go` to understand attributes, required fields, and CRUD behavior.
2. **Read project infrastructure**: See [references/project-infrastructure.md](references/project-infrastructure.md) for the `internal/acctest/` helpers and `internal/sweep/` utilities available.
3. **Read test patterns**: See [references/patterns.md](references/patterns.md) for complete annotated examples of resource tests, data source tests, import tests, and sweepers.
4. **Write the test file**: Create `resource_test.go` (or `data_source_test.go`) in the service's package directory.
5. **Write the sweeper**: Create `sweep.go` in the service's package directory for any resource that creates infrastructure.
6. **Run and debug**: Execute `./scripts/acctest ./internal/services/<service>/... -run TestAccName -v` and fix failures.

## File Placement

```
internal/services/<service>/
├── resource.go              # (generated) CRUD operations
├── schema.go                # (generated) Terraform schema
├── model.go                 # (generated) data models
├── resource_test.go         # acceptance tests (you write this)
├── data_source_test.go      # data source acceptance tests (you write this)
└── sweep.go                 # sweeper implementation (you write this)
```

## Test Skeleton

Every acceptance test file follows this structure:

```go
package <service>_test

import (
    "fmt"
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/hashicorp/terraform-plugin-testing/knownvalue"
    "github.com/hashicorp/terraform-plugin-testing/statecheck"
    "github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

    "github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCloudExample_basic(t *testing.T) {
    rName := acctest.RandomName()

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:                 func() { acctest.PreCheck(t) },
        ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
        CheckDestroy:             testAccCheckCloudExampleDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccCloudExampleConfig(rName),
                ConfigStateChecks: []statecheck.StateCheck{
                    statecheck.ExpectKnownValue("gcore_cloud_example.test",
                        tfjsonpath.New("name"), knownvalue.StringExact(rName)),
                },
            },
        },
    })
}
```

## Key Rules

- **Function naming**: `TestAcc<ResourceType>_<variant>` (e.g., `TestAccCloudSSHKey_basic`)
- **Resource address in config**: Always use `"test"` as the resource label (e.g., `gcore_cloud_ssh_key.test`)
- **Random names**: Always use `acctest.RandomName()` (produces `"tf-test-<token>-<random>"`). The token identifies the run, and the sweeper deletes only exact token matches, so a hand-written name is unreclaimable
- **Qualified names**: Append, never prepend -- `rName + "-subnet"`, not `"tf-test-subnet-" + rName`. Prefixing repeats `tf-test-` further along the name and wastes the length budget
- **Names the API will not take**: use `acctest.RandomNameCompact()` where hyphens are rejected or the cap is tight (registry users: lowercase alphanumeric, 16 max)
- **Resources with no name field**: put `acctest.OwnershipTagHCL()` in the config's `tags` block. Floating IPs have no name at all, so the tag is the only thing that tells a test address from a production one
- **Provider factories**: Always use `acctest.ProtoV6ProviderFactories` (this is a Protocol V6 provider)
- **PreCheck**: Always include `func() { acctest.PreCheck(t) }` to validate env vars
- **Parallelism**: Use `resource.ParallelTest()` for parallel execution (it calls `t.Parallel()` internally). Use `resource.Test()` when tests must run serially (e.g., resources that conflict in the same project)
- **Config helpers**: Use `fmt.Sprintf` with `%[1]q` (quoted) or `%[1]s` (unquoted) indexed verbs
- **project_id/region_id**: Use `acctest.ProjectID()` and `acctest.RegionID()` in configs
- **CheckDestroy**: Always implement; use `acctest.CheckResourceDestroyed()` helper
- **Import test**: Add an import step after the basic create step. Prefer `ImportStateKind: resource.ImportBlockWithID` (plannable import) for new tests. Use `ImportStateVerify: true` only when testing the legacy `terraform import` CLI workflow. See [patterns.md](references/patterns.md) for details
- **State checks**: Prefer `ConfigStateChecks` with `statecheck.*` over legacy `Check` with `resource.TestCheckResourceAttr`
- **Plan checks**: Use `ConfigPlanChecks` with `plancheck.*` when verifying plan behavior
- **Config as functions**: Define config as `func testAccCloudExampleConfig(name string) string` returning HCL

## Running Tests

```bash
# single test
./scripts/acctest ./internal/services/cloud_ssh_key/... -run TestAccCloudSSHKey_basic -v

# all tests for a service
./scripts/acctest ./internal/services/cloud_ssh_key/...

# all acceptance tests (slow)
./scripts/acctest
```

## Sweeper Skeleton

Every resource that creates infrastructure needs a sweeper in `sweep.go`:

```go
package <service>

import (
    "fmt"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/G-Core/terraform-provider-gcore/internal/acctest"
    "github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
    resource.AddTestSweepers("gcore_cloud_example", &resource.Sweeper{
        Name: "gcore_cloud_example",
        F:    sweepCloudExample,
    })
}

func sweepCloudExample(_ string) error {
    if err := sweep.ValidateSweeperEnvironment(); err != nil {
        return err
    }
    // list resources, filter with sweep.ShouldSweep(), delete matches
}
```

`sweep.ShouldSweep` requires more than a name prefix: the resource must carry the run token named by
`GCORE_SWEEP_RUN_ID`, and with that unset **no tokenized resource is swept**. Deletion is only safe once the owning
run has finished, and no sweeper can determine that for itself, so the caller has to say which run it is reclaiming.
`scripts/acctest` prints the command. The one exception is `GCORE_SWEEP_UNTOKENIZED=1`, which reclaims names predating
tokens and ignores the run id; see the infrastructure reference before using it.

### Resources with no name

Some APIs expose no name to carry the token -- floating IPs are the case that forced this. Read it from the
ownership tag and use `sweep.ShouldSweepToken` instead:

```go
var token string
for _, tag := range fip.Tags {
    if tag.Key == sweep.TagKey {
        token = tag.Value
        break
    }
}
if !sweep.ShouldSweepToken("gcore_cloud_floating_ip", fip.ID, token) {
    continue
}
```

Do **not** substitute a property that merely correlates with being a test resource. This sweeper previously
deleted every floating IP whose status was `DOWN`, reasoning that tests create standalone addresses -- but
unattached is the ordinary state of an address reserved and not yet assigned, so it deleted production addresses
too. An untagged resource is indistinguishable from a production one; leave it alone.

`ShouldSweepToken` never deletes a resource with no token, and `GCORE_SWEEP_UNTOKENIZED=1` does not change that:
the opt-in is bounded by the `tf-test-` name prefix, and a resource with no name has none. Untagged addresses left
by earlier runs need a manual audit.

If the API accepts neither a name nor tags, ownership cannot be established. Do not invent a filter that looks
like an ownership check but is not one -- an earlier reserved-fixed-IP sweeper deleted any unnamed address on the
guess that it was test-owned. Either make the sweeper log and return, or put the unsafe deletion behind an
explicit opt-in that says what it does.

Reserved fixed IPs take the second route. `name` is computed by the API, there is no tags field, and an external
address sits on the region's shared public network rather than a test one, so nothing identifies the run.
`GCORE_SWEEP_ALL_RESERVED_FIXED_IPS=1` deletes every external address in the `DOWN` state in the project; unset,
the sweeper deletes nothing. What makes that acceptable is scope, not precision: the list call is project-scoped,
so it reaches only a project dedicated to acceptance tests. **Do not copy the pattern into an account-global
sweeper**, where there is no such boundary.

The residual cost is that it cannot tell a leaked address from one a concurrent run still holds. CI runs are
serialised by a shared concurrency group, so the exposure is a local run overlapping a scheduled sweep, and the
result is a confusing test failure rather than lost production capacity. Every deletion names the variable in the
log so that failure is diagnosable.

**Manual audit**, for anything no sweeper reclaims: in the test project, list unattached floating IPs that lack
the `gcore_tf_acctest_run` tag and delete the ones no one recognises. The fix that would retire this step is
`name` or `tags` support in those APIs.

## References

- [Project infrastructure details](references/project-infrastructure.md) -- all `acctest.*` and `sweep.*` helpers with signatures
- [Complete test patterns and examples](references/patterns.md) -- annotated examples for every test type

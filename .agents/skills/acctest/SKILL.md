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

    "github.com/stainless-sdks/gcore-terraform/internal/acctest"
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
- **Random names**: Always use `acctest.RandomName()` (produces `"tf-test-<10chars>"`) for sweeper compatibility
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
    "github.com/stainless-sdks/gcore-terraform/internal/acctest"
    "github.com/stainless-sdks/gcore-terraform/internal/sweep"
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

## References

- [Project infrastructure details](references/project-infrastructure.md) -- all `acctest.*` and `sweep.*` helpers with signatures
- [Complete test patterns and examples](references/patterns.md) -- annotated examples for every test type

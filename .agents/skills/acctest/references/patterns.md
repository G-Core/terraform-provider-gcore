# Acceptance Test Patterns

Complete annotated examples for every test type in this project.

## Table of Contents

1. [Resource: Basic Create/Destroy](#1-resource-basic-createdestroy)
2. [Resource: Create + Update](#2-resource-create--update)
3. [Resource: Import](#3-resource-import)
4. [Resource: Full Lifecycle (Create + Update + Import)](#4-resource-full-lifecycle)
5. [Data Source Test](#5-data-source-test)
6. [Error Testing](#6-error-testing)
7. [Plan Checks](#7-plan-checks)
8. [CheckDestroy Patterns](#8-checkdestroy-patterns)
9. [Config Helper Patterns](#9-config-helper-patterns)
10. [Sweeper Implementation](#10-sweeper-implementation)
11. [Known Value Checks Reference](#11-known-value-checks-reference)
12. [Resource Plan Check Actions](#12-resource-plan-check-actions)
13. [Dependencies Between Resources](#13-dependencies-between-resources)
14. [Cross-Step Value Comparison](#14-cross-step-value-comparison)

---

## 1. Resource: Basic Create/Destroy

Minimal test that creates a resource and verifies attributes in state.

```go
package cloud_ssh_key_test

import (
    "fmt"
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/compare"
    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/hashicorp/terraform-plugin-testing/knownvalue"
    "github.com/hashicorp/terraform-plugin-testing/statecheck"
    "github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

    "github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCloudSSHKey_basic(t *testing.T) {
    rName := acctest.RandomName()

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:                 func() { acctest.PreCheck(t) },
        ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
        CheckDestroy:             testAccCheckCloudSSHKeyDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccCloudSSHKeyConfig(rName),
                ConfigStateChecks: []statecheck.StateCheck{
                    statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
                        tfjsonpath.New("name"), knownvalue.StringExact(rName)),
                    statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
                        tfjsonpath.New("id"), knownvalue.NotNull()),
                },
            },
        },
    })
}

func testAccCloudSSHKeyConfig(name string) string {
    return fmt.Sprintf(`
resource "gcore_cloud_ssh_key" "test" {
  project_id = %[1]s
  name       = %[2]q
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAA... test@example.com"
}`, acctest.ProjectID(), name)
}
```

---

## 2. Resource: Create + Update

Tests that an attribute can be modified in-place.

```go
func TestAccCloudSSHKey_update(t *testing.T) {
    rName := acctest.RandomName()

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:                 func() { acctest.PreCheck(t) },
        ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
        CheckDestroy:             testAccCheckCloudSSHKeyDestroy,
        Steps: []resource.TestStep{
            // Step 1: Create
            {
                Config: testAccCloudSSHKeyConfig(rName),
                ConfigStateChecks: []statecheck.StateCheck{
                    statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
                        tfjsonpath.New("name"), knownvalue.StringExact(rName)),
                    statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
                        tfjsonpath.New("shared_in_project"), knownvalue.Bool(false)),
                },
            },
            // Step 2: Update
            {
                Config: testAccCloudSSHKeyConfigUpdated(rName),
                ConfigStateChecks: []statecheck.StateCheck{
                    statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
                        tfjsonpath.New("shared_in_project"), knownvalue.Bool(true)),
                },
            },
        },
    })
}

func testAccCloudSSHKeyConfigUpdated(name string) string {
    return fmt.Sprintf(`
resource "gcore_cloud_ssh_key" "test" {
  project_id        = %[1]s
  name              = %[2]q
  public_key        = "ssh-rsa AAAAB3NzaC1yc2EAAAA... test@example.com"
  shared_in_project = true
}`, acctest.ProjectID(), name)
}
```

---

## 3. Resource: Import

Tests that importing an existing resource produces the same state as creating it.
Place after a lifecycle step; the framework uses the previous step's state as the "golden file."

These approaches test **different provider code paths** and are not interchangeable:

### Import with `import` block (plannable, recommended for new tests)

Tests the `import` block workflow (Terraform 1.5+). Generates a plan with an `import` block
and verifies the planned values match the previous step's state. This is the modern approach.

```go
{
    ResourceName:    "gcore_cloud_ssh_key.test",
    ImportState:     true,
    ImportStateKind: resource.ImportBlockWithID,
}
```

### Import with `terraform import` command (legacy CLI workflow)

Tests the `terraform import` CLI command. Runs the import, then compares the resulting state
against the previous step's state. Use `ImportStateVerifyIgnore` for attributes that can't
round-trip through import (e.g., write-only fields).

```go
{
    ResourceName:      "gcore_cloud_ssh_key.test",
    ImportState:       true,
    ImportStateVerify: true,
    // if some attributes can't be imported, exclude them:
    ImportStateVerifyIgnore: []string{"public_key"},
}
```

### Import with composite ID (legacy CLI workflow)

When the resource uses a composite import path like `project_id/resource_id`:

```go
{
    ResourceName:      "gcore_cloud_ssh_key.test",
    ImportState:       true,
    ImportStateVerify:  true,
    ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_ssh_key.test", "project_id", "id"),
}
```

---

## 4. Resource: Full Lifecycle

The recommended pattern combining create, update, and import in a single test:

```go
func TestAccCloudSSHKey_full(t *testing.T) {
    rName := acctest.RandomName()

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:                 func() { acctest.PreCheck(t) },
        ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
        CheckDestroy:             testAccCheckCloudSSHKeyDestroy,
        Steps: []resource.TestStep{
            // Step 1: Create
            {
                Config: testAccCloudSSHKeyConfig(rName),
                ConfigStateChecks: []statecheck.StateCheck{
                    statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
                        tfjsonpath.New("name"), knownvalue.StringExact(rName)),
                },
            },
            // Step 2: Update
            {
                Config: testAccCloudSSHKeyConfigUpdated(rName),
                ConfigStateChecks: []statecheck.StateCheck{
                    statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
                        tfjsonpath.New("shared_in_project"), knownvalue.Bool(true)),
                },
            },
            // Step 3: Import (plannable import block)
            {
                ResourceName:    "gcore_cloud_ssh_key.test",
                ImportState:     true,
                ImportStateKind: resource.ImportBlockWithID,
            },
        },
    })
}
```

---

## 5. Data Source Test

Data source tests typically create a resource first, then read it with a data source.

```go
func TestAccCloudSSHKeyDataSource_basic(t *testing.T) {
    rName := acctest.RandomName()

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:                 func() { acctest.PreCheck(t) },
        ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
        // CheckDestroy verifies the resource created for the data source is cleaned up
        CheckDestroy:             testAccCheckCloudSSHKeyDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccCloudSSHKeyDataSourceConfig(rName),
                ConfigStateChecks: []statecheck.StateCheck{
                    statecheck.ExpectKnownValue("data.gcore_cloud_ssh_key.test",
                        tfjsonpath.New("name"), knownvalue.StringExact(rName)),
                    // verify data source reads the same ID as the resource
                    statecheck.CompareValuePairs(
                        "gcore_cloud_ssh_key.test", tfjsonpath.New("id"),
                        "data.gcore_cloud_ssh_key.test", tfjsonpath.New("id"),
                        compare.ValuesSame(),
                    ),
                },
            },
        },
    })
}

func testAccCloudSSHKeyDataSourceConfig(name string) string {
    return fmt.Sprintf(`
resource "gcore_cloud_ssh_key" "test" {
  project_id = %[1]s
  name       = %[2]q
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAA... test@example.com"
}

data "gcore_cloud_ssh_key" "test" {
  project_id = %[1]s
  id         = gcore_cloud_ssh_key.test.id
}`, acctest.ProjectID(), name)
}
```

Note: Since the config creates a resource (to have something to read), include `CheckDestroy` to verify the resource is cleaned up via API after the test. `CheckDestroy` runs after the framework's `terraform destroy` to confirm the resource is actually gone.

---

## 6. Error Testing

Test that invalid configurations produce expected errors.

```go
func TestAccCloudSSHKey_invalidName(t *testing.T) {
    resource.ParallelTest(t, resource.TestCase{
        PreCheck:                 func() { acctest.PreCheck(t) },
        ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: `
resource "gcore_cloud_ssh_key" "test" {
  name = ""
}`,
                ExpectError: regexp.MustCompile(`must not be empty`),
            },
        },
    })
}
```

---

## 7. Plan Checks

Assert plan behavior before apply:

```go
{
    Config: testAccCloudSSHKeyConfig(rName),
    ConfigPlanChecks: resource.ConfigPlanChecks{
        PreApply: []plancheck.PlanCheck{
            // verify resource will be created
            plancheck.ExpectResourceAction("gcore_cloud_ssh_key.test",
                plancheck.ResourceActionCreate),
            // verify a specific planned value
            plancheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
                tfjsonpath.New("name"), knownvalue.StringExact(rName)),
            // verify a computed attribute is unknown at plan time
            plancheck.ExpectUnknownValue("gcore_cloud_ssh_key.test",
                tfjsonpath.New("id")),
        },
    },
}
```

Verify empty plan (no drift) after re-applying the same config:

```go
// Step 2: re-apply same config, assert no changes
{
    Config: testAccCloudSSHKeyConfig(rName),
    ConfigPlanChecks: resource.ConfigPlanChecks{
        PreApply: []plancheck.PlanCheck{
            plancheck.ExpectEmptyPlan(),
        },
    },
}
```

### Expecting a non-empty plan (drift testing)

When a valid config intentionally results in perpetual diff (e.g., API normalizes a value differently),
use `ExpectNonEmptyPlan` on the TestStep to prevent the test from failing:

```go
{
    // config where the API normalizes a field to a different value
    Config:             testAccCloudExampleConfigWithDrift(rName),
    ExpectNonEmptyPlan: true,
}
```

---

## 8. CheckDestroy Patterns

### Simple resource (single ID)

```go
func testAccCheckCloudSSHKeyDestroy(s *terraform.State) error {
    return acctest.CheckResourceDestroyed(s, "gcore_cloud_ssh_key", func(client *gcore.Client, id string) error {
        _, err := client.Cloud.SSHKeys.Get(context.Background(), id)
        return err
    })
}
```

### Resource requiring project_id/region_id in API call

When the API requires additional parameters beyond the resource ID:

```go
func testAccCheckCloudVolumeDestroy(s *terraform.State) error {
    client, err := acctest.NewGcoreClient()
    if err != nil {
        return err
    }

    for _, rs := range s.RootModule().Resources {
        if rs.Type != "gcore_cloud_volume" {
            continue
        }

        projectID, err := strconv.ParseInt(rs.Primary.Attributes["project_id"], 10, 64)
        if err != nil {
            return fmt.Errorf("error parsing project_id: %w", err)
        }
        regionID, err := strconv.ParseInt(rs.Primary.Attributes["region_id"], 10, 64)
        if err != nil {
            return fmt.Errorf("error parsing region_id: %w", err)
        }

        _, err := client.Cloud.Volumes.Get(context.Background(), rs.Primary.ID, cloud.VolumeGetParams{
            ProjectID: gcore.Int(projectID),
            RegionID:  gcore.Int(regionID),
        })

        if err == nil {
            return fmt.Errorf("volume %s still exists", rs.Primary.ID)
        }
        if !acctest.IsNotFoundError(err) {
            return fmt.Errorf("error checking volume deletion: %w", err)
        }
    }
    return nil
}
```

---

## 9. Config Helper Patterns

### Basic config with project_id

```go
func testAccCloudExampleConfig(name string) string {
    return fmt.Sprintf(`
resource "gcore_cloud_example" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
```

### Config with dependencies

```go
func testAccCloudSubnetConfig(name string) string {
    return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  network_id = gcore_cloud_network.test.id
  cidr       = "10.0.0.0/24"
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
```

### Config for data source

```go
func testAccCloudExampleDataSourceConfig(name string) string {
    return fmt.Sprintf(`
resource "gcore_cloud_example" "test" {
  project_id = %[1]s
  name       = %[2]q
}

data "gcore_cloud_example" "test" {
  project_id = %[1]s
  id         = gcore_cloud_example.test.id
}`, acctest.ProjectID(), name)
}
```

---

## 10. Sweeper Implementation

Sweepers clean up leaked test resources. Place in `sweep.go` (not `_test.go`) in the service package.

**Note**: The sweeper file uses the main package (not `_test`), so it has access to internal package types.

```go
package cloud_ssh_key

import (
    "context"
    "fmt"
    "log"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"

    "github.com/G-Core/terraform-provider-gcore/internal/acctest"
    "github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
    resource.AddTestSweepers("gcore_cloud_ssh_key", &resource.Sweeper{
        Name: "gcore_cloud_ssh_key",
        F:    sweepCloudSSHKeys,
    })
}

func sweepCloudSSHKeys(_ string) error {
    if err := sweep.ValidateSweeperEnvironment(); err != nil {
        return err
    }

    client, err := acctest.NewGcoreClient()
    if err != nil {
        return err
    }

    ctx := context.Background()

    // list all SSH keys
    keys, err := client.Cloud.SSHKeys.List(ctx)
    if err != nil {
        if sweep.SkipSweepError(err) {
            log.Printf("[WARN] Skipping SSH key sweep: %s", err)
            return nil
        }
        return fmt.Errorf("error listing SSH keys: %w", err)
    }

    for _, key := range keys.Items {
        if !sweep.ShouldSweep("gcore_cloud_ssh_key", key.Name) {
            continue
        }

        log.Printf("[INFO] Deleting SSH key: %s (%s)", key.Name, key.ID)
        err := client.Cloud.SSHKeys.Delete(ctx, key.ID)
        if err != nil {
            log.Printf("[ERROR] Failed to delete SSH key %s: %s", key.Name, err)
        }
    }

    return nil
}
```

### Sweeper with dependencies

When resource B depends on resource A, sweep B first:

```go
func init() {
    resource.AddTestSweepers("gcore_cloud_network_subnet", &resource.Sweeper{
        Name:         "gcore_cloud_network_subnet",
        F:            sweepCloudNetworkSubnets,
        Dependencies: []string{"gcore_cloud_instance"}, // instances must be deleted first
    })
}
```

### sweep_test.go (TestMain entry point)

This file must exist at `internal/sweep/sweep_test.go` to enable the `go test -sweep` CLI.

**Important**: Every new `sweep.go` in a service package **must** have a corresponding blank import
added here. The blank import triggers the service package's `init()` function which calls
`resource.AddTestSweepers()`. Without the import, the sweeper won't be registered and won't run.

```go
package sweep

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"

    // blank imports trigger init() in each service package, registering sweepers
    // MAINTENANCE: add a blank import for every service that has a sweep.go
    _ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_ssh_key"
    _ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_volume"
)

func TestMain(m *testing.M) {
    resource.TestMain(m)
}
```

---

## 11. Known Value Checks Reference

Used with `statecheck.ExpectKnownValue()` and `plancheck.ExpectKnownValue()`:

| Check | Example |
|-------|---------|
| `knownvalue.StringExact("val")` | Exact string match |
| `knownvalue.StringRegexp(regexp.MustCompile(`^tf-test-`))` | Regex string match |
| `knownvalue.Bool(true)` | Boolean value |
| `knownvalue.Int64Exact(42)` | Exact int64 |
| `knownvalue.Float64Exact(1.5)` | Exact float64 |
| `knownvalue.NotNull()` | Any non-null value |
| `knownvalue.Null()` | Null value |
| `knownvalue.ListExact([]knownvalue.Check{...})` | Exact list contents |
| `knownvalue.ListSizeExact(3)` | List with N elements |
| `knownvalue.SetExact([]knownvalue.Check{...})` | Exact set contents |
| `knownvalue.SetSizeExact(2)` | Set with N elements |
| `knownvalue.MapExact(map[string]knownvalue.Check{...})` | Exact map |
| `knownvalue.MapSizeExact(2)` | Map with N keys |
| `knownvalue.ObjectExact(map[string]knownvalue.Check{...})` | Exact object |

### tfjsonpath Navigation

```go
tfjsonpath.New("top_level")                    // top-level attribute
tfjsonpath.New("block").AtMapKey("key")        // nested block/map key
tfjsonpath.New("list").AtSliceIndex(0)         // list element by index
tfjsonpath.New("block").AtMapKey("nested").AtMapKey("deep")  // deeply nested
```

---

## 12. Resource Plan Check Actions

Used with `plancheck.ExpectResourceAction()`:

| Action | Meaning |
|--------|---------|
| `plancheck.ResourceActionCreate` | Resource will be created |
| `plancheck.ResourceActionUpdate` | Resource will be updated in-place |
| `plancheck.ResourceActionDestroy` | Resource will be destroyed |
| `plancheck.ResourceActionDestroyBeforeCreate` | Destroy then recreate (ForceNew) |
| `plancheck.ResourceActionCreateBeforeDestroy` | Create replacement then destroy old |
| `plancheck.ResourceActionNoop` | No changes planned |

---

## 13. Dependencies Between Resources

When testing a resource that depends on another resource, create the dependency in the same config:

```go
func testAccCloudLoadBalancerListenerConfig(name string) string {
    return fmt.Sprintf(`
resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-lb"
}

resource "gcore_cloud_load_balancer_listener" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  name             = %[3]q
  loadbalancer_id  = gcore_cloud_load_balancer.test.id
  protocol         = "HTTP"
  protocol_port    = 80
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
```

The Terraform test framework handles creation order automatically via the dependency graph, and destroys in reverse order.

---

## 14. Cross-Step Value Comparison

Verify that a computed attribute stays the same (or changes) across test steps using `statecheck.CompareValue`.
The comparer is created once and reused across steps via `AddStateValue`.

```go
func TestAccCloudExample_computedStable(t *testing.T) {
    rName := acctest.RandomName()

    // create the comparer outside the test steps
    compareIDSame := statecheck.CompareValue(compare.ValuesSame())

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:                 func() { acctest.PreCheck(t) },
        ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
        CheckDestroy:             testAccCheckCloudExampleDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccCloudExampleConfig(rName),
                ConfigStateChecks: []statecheck.StateCheck{
                    // record the ID after create
                    compareIDSame.AddStateValue(
                        "gcore_cloud_example.test",
                        tfjsonpath.New("id"),
                    ),
                },
            },
            {
                Config: testAccCloudExampleConfigUpdated(rName),
                ConfigStateChecks: []statecheck.StateCheck{
                    // verify the ID didn't change after update (in-place update)
                    compareIDSame.AddStateValue(
                        "gcore_cloud_example.test",
                        tfjsonpath.New("id"),
                    ),
                },
            },
        },
    })
}
```

Use `compare.ValuesDiffer()` instead of `compare.ValuesSame()` to assert a value **changed** across steps (e.g., verifying a computed timestamp updates).

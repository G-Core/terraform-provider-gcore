# Project Acceptance Test Infrastructure

## Provider Setup (`internal/acctest/provider.go`)

Protocol V6 provider factories and PreCheck:

```go
import "github.com/G-Core/terraform-provider-gcore/internal/acctest"

// use in every TestCase
acctest.ProtoV6ProviderFactories  // map[string]func() (tfprotov6.ProviderServer, error)
acctest.PreCheck(t)               // fatals if GCORE_API_KEY, GCORE_CLOUD_PROJECT_ID, or GCORE_CLOUD_REGION_ID are unset
```

## Helpers (`internal/acctest/helpers.go`)

```go
acctest.ProjectID()   // returns os.Getenv("GCORE_CLOUD_PROJECT_ID")
acctest.RegionID()    // returns os.Getenv("GCORE_CLOUD_REGION_ID")
acctest.RandomName()  // returns "tf-test-<10 random alphanumeric chars>"
```

## Resource Helpers (`internal/acctest/resource.go`)

### API Client

```go
client, err := acctest.NewGcoreClient()  // *gcore.Client from env vars
```

### CheckDestroy

Use in `TestCase.CheckDestroy` to verify resources are deleted after test:

```go
func testAccCheckCloudSSHKeyDestroy(s *terraform.State) error {
    return acctest.CheckResourceDestroyed(s, "gcore_cloud_ssh_key", func(client *gcore.Client, id string) error {
        _, err := client.Cloud.SSHKeys.Get(context.Background(), id)
        return err
    })
}
```

Signature: `CheckResourceDestroyed(s *terraform.State, resourceType string, checkFunc func(*gcore.Client, string) error) error`

- Iterates all resources of `resourceType` in state
- Calls `checkFunc` with each resource's ID
- If `checkFunc` returns `nil` -> resource still exists -> **test fails**
- If `checkFunc` returns a "not found" error -> resource deleted -> **success**
- Uses `IsNotFoundError()` to detect 404-like errors

### CheckExists

Use in legacy `Check` functions to verify resource exists during test:

```go
acctest.CheckResourceExists("gcore_cloud_ssh_key.test", func(client *gcore.Client, id string) error {
    _, err := client.Cloud.SSHKeys.Get(context.Background(), id)
    return err
})
```

Signature: `CheckResourceExists(resourceName string, checkFunc func(*gcore.Client, string) error) func(*terraform.State) error`

### BuildImportID

Construct composite import IDs from state attributes:

```go
// produces "project_id_value/resource_id_value"
acctest.BuildImportID("gcore_cloud_ssh_key.test", "project_id", "id")
```

Signature: `BuildImportID(resourceName string, attrNames ...string) func(*terraform.State) (string, error)`

### IsNotFoundError

```go
acctest.IsNotFoundError(err)  // checks for: "not found", "Not Found", "NOT_FOUND", "does not exist", "doesn't exist", "404", "NotFound", "NoSuchEntity"
```

## Sweep Utilities (`internal/sweep/`)

### Constants and Validation (`sweep.go`)

```go
sweep.ResourcePrefix           // "tf-test"
sweep.Context(region)          // returns context.Background()
sweep.ValidateSweeperEnvironment()  // checks GCORE_API_KEY, GCORE_CLOUD_PROJECT_ID, GCORE_CLOUD_REGION_ID
```

### Name Filtering (`framework.go`)

```go
sweep.IsTestResource(name)      // true for "tf-test-*", "tf_test_*", or bare 10-char alphanumeric
sweep.ShouldSweep(resType, name)  // IsTestResource + logging
sweep.SkipSweepError(err)       // true for "Access Denied", "Forbidden", "403", etc.
```

## Missing Infrastructure

The following documented items do **not exist yet** and must be created if needed:

- `internal/sweep/sweep_test.go` with `TestMain(m *testing.M) { resource.TestMain(m) }` -- required to invoke sweepers via CLI
- Per-service `sweep.go` files -- none exist yet
- All `resource_test.go` and `data_source_test.go` files -- none exist yet

## Environment Variables

| Variable | Required | Purpose |
|----------|----------|---------|
| `GCORE_API_KEY` | yes | API authentication |
| `GCORE_CLOUD_PROJECT_ID` | yes | Default project for cloud resources |
| `GCORE_CLOUD_REGION_ID` | yes | Default region for cloud resources |
| `GCORE_BASE_URL` | no | Override API base URL |
| `TF_ACC` | yes (set by scripts/acctest) | Enable acceptance tests |

## Gcore API Client Patterns

The Gcore Go SDK client is structured as:

```go
client.Cloud.SSHKeys.Get(ctx, id)
client.Cloud.SSHKeys.List(ctx, cloud.SSHKeyListParams{...})
client.Cloud.SSHKeys.Delete(ctx, id)
client.Cloud.Volumes.Get(ctx, id, cloud.VolumeGetParams{ProjectID: ..., RegionID: ...})
```

Most cloud resources require `project_id` and `region_id` in their API calls. Check the resource's `resource.go` file for the exact SDK method signatures.

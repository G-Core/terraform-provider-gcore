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
acctest.RandomName()         // "tf-test-<token4>-<random6>"; the token identifies the run
acctest.RandomNameCompact()  // "tftest<token4><random6>", exactly 16 chars, no hyphens
acctest.OwnershipTagHCL()    // a tags entry carrying the token, for resources with no name
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

### Validation (`sweep.go`)

```go
sweep.ValidateSweeperEnvironment()  // checks GCORE_API_KEY, GCORE_CLOUD_PROJECT_ID, GCORE_CLOUD_REGION_ID
```

### Name Filtering (`framework.go`)

```go
sweep.IsTestResource(name)        // true only for the "tf-test-" / "tf_test_" prefixes
sweep.ShouldSweep(resType, name)  // prefix match AND an exact run-token match
sweep.ShouldSweepToken(resType, id, token)  // the same decision, for resources with no name; "" is never swept
sweep.SkipSweepError(err)       // true for "Access Denied", "Forbidden", "403", etc.
```

`IsTestResource` matches an exact prefix and nothing else -- no length or shape heuristics -- and the trailing
separator is required, so `tf-testing-prod` does not match. **Do not widen it.** API keys are client-scoped and 14
sweepers list an entire account with no project to bound them; for those, this predicate is the only thing between a
sweep and unrelated infrastructure. The prefixes are unexported consts so the allowlist cannot be widened at run time.

Any new name generator must produce names this accepts. `TestRandomNameIsSweepable` enforces that for
`acctest.RandomName`, so a change there fails the build rather than silently creating resources no sweeper can reclaim.

### Run Tokens (`token.go`)

```go
sweep.Token(runID)         // reduces any run id to a fixed-width, name-safe token
sweep.TokenFromName(name)  // (token, ok); ok is false for names predating tokens
sweep.RunToken()           // this run's token, or "" when GCORE_ACCTEST_RUN_ID is unset
sweep.TagKey               // "gcore_tf_acctest_run"; carries the token where no name can
```

Every run gets a token and every resource it creates carries one, so a sweep deletes exactly the run it was told to
and provably nothing else. That is what makes concurrent runs safe: the nightly suite and an engineer testing a
subset against the same project cannot see each other's resources.

**A sweep with `GCORE_SWEEP_RUN_ID` unset deletes no tokenized resource.** Deletion is only safe once the owning run
has finished, and nothing in a sweeper can establish that -- a scheduled job least of all. The caller supplies the
knowledge. There is deliberately no age heuristic: a run longer than any plausible threshold would start deleting its
own live resources, and the nightly suite runs for hours.

`GCORE_SWEEP_UNTOKENIZED=1` is the one exception. It sweeps every test-prefixed name that carries no token and does
not consult `GCORE_SWEEP_RUN_ID`, because such a name records no run to compare against -- requiring a run id would
scope nothing. For the same reason it cannot tell a leaked resource from a live one, so only use it once no branch or
job still runs code that creates untokenized names.

It never reaches resources with no name, such as floating IPs. For those the prefix is not there to bound it, so an
untagged one is indistinguishable from production and `ShouldSweepToken` always skips it.

`Token` hashes rather than truncates, so a GitHub run id, a UUID and an ad hoc string all yield something name-safe,
and the test side and the sweeper derive the same token from the same run id without coordinating.

## Environment Variables

| Variable | Required | Purpose |
|----------|----------|---------|
| `GCORE_API_KEY` | yes | API authentication |
| `GCORE_CLOUD_PROJECT_ID` | yes | Default project for cloud resources |
| `GCORE_CLOUD_REGION_ID` | yes | Default region for cloud resources |
| `GCORE_BASE_URL` | no | Override API base URL |
| `TF_ACC` | yes (set by scripts/acctest) | Enable acceptance tests |
| `GCORE_ACCTEST_RUN_ID` | set by scripts/acctest | Identifies the run creating resources; see below |
| `GCORE_SWEEP_RUN_ID` | yes, to delete tokenized resources | The finished run a sweep may reclaim. Unset means no tokenized resource is swept |
| `GCORE_SWEEP_ACCOUNT_ID` | yes, for any sweep | The account the credentials must resolve to. Every sweeper refuses to run without it or on a mismatch. Checked once per sweep with one API call |
| `GCORE_POLLING_TIMEOUT_SECONDS` | no | Longest a single create or delete may poll. `scripts/acctest` sets 1200; sweeps default to 600. Keep it well inside the `go test -timeout`, since hitting that deadline panics past cleanup and reporting |
| `GCORE_SWEEP_UNTOKENIZED` | no | Sweeps test-prefixed names carrying no token, ignoring `GCORE_SWEEP_RUN_ID`. Never applies to resources with no name. For a one-time project audit, not routine use |
| `GCORE_SWEEP_ALL_RESERVED_FIXED_IPS` | no | Deletes every external reserved fixed IP in the `DOWN` state in the project, with no ownership check -- that family can carry none. Bounded by the list call being project-scoped. Unset means the sweeper deletes nothing |

One `GCORE_ACCTEST_RUN_ID` covers a whole run. Each package is its own test binary, so without it each mints a
different run id and one run's resources scatter across several of them. Each binary logs its id with a `[WARN]`
line giving the sweep command, but `go test` hides a passing package's output unless `-v` is set, so run through
`scripts/acctest` whenever the resources may need reclaiming.

## Gcore API Client Patterns

The Gcore Go SDK client is structured as:

```go
client.Cloud.SSHKeys.Get(ctx, id)
client.Cloud.SSHKeys.List(ctx, cloud.SSHKeyListParams{...})
client.Cloud.SSHKeys.Delete(ctx, id)
client.Cloud.Volumes.Get(ctx, id, cloud.VolumeGetParams{ProjectID: ..., RegionID: ...})
```

Most cloud resources require `project_id` and `region_id` in their API calls. Check the resource's `resource.go` file for the exact SDK method signatures.

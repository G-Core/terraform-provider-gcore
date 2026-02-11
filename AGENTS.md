# AGENTS.md

Guidelines for AI agents working with this Stainless-generated Terraform provider.

## Project Overview

This is a **Stainless-generated** Terraform provider for Gcore cloud services. The code is
automatically generated from an OpenAPI specification via Stainless codegen.

**Important**: Custom changes should be kept to a minimum. Prefer modifying the OpenAPI spec
or Stainless config to have changes generated rather than maintaining custom code. Custom
code creates merge conflict risk during regeneration. See the Stainless docs on adding
custom code: https://www.stainless.com/docs/guides/add-custom-code

## Build Commands

```bash
./scripts/bootstrap          # install dependencies (Go modules, Homebrew on macOS)
./scripts/build              # build provider binary -> terraform-provider-gcore
./scripts/format             # format code with gofmt + generate terraform docs
./scripts/lint               # run go build to check for compilation errors
./scripts/generate-docs      # generate terraform provider documentation
```

## Testing

### Quick Reference

```bash
./scripts/test       # run unit tests (fast, safe, no API calls)
./scripts/acctest    # run acceptance tests (slow, creates real resources)
./scripts/sweep      # clean up leaked test resources
```

### Unit Tests (Schema Parity)

**Fast, safe, no API calls** - These are Stainless-generated tests that validate schema-model consistency:

```bash
# run all unit and schema tests
./scripts/test

# run a single test file
go test ./internal/services/cloud_load_balancer/... -v

# run a specific test by name
go test ./internal/services/cloud_load_balancer/... -run TestCloudLoadBalancerModelSchemaParity -v
```

### Acceptance Tests

**Slow, creates real infrastructure, requires credentials** - These test against actual Gcore APIs:

```bash
# run all acceptance tests (creates real cloud resources, may incur costs)
./scripts/acctest

# run acceptance tests for specific service
./scripts/acctest ./internal/services/cloud_load_balancer/...

# run specific acceptance test
./scripts/acctest -run TestAccCloudLoadBalancer_basic

# pass any go test flags
./scripts/acctest -v -parallel 8 -timeout 60m
```

**Direct usage (without script):**
```bash
# same as ./scripts/acctest but with explicit TF_ACC=1
TF_ACC=1 go test -parallel=4 -timeout=120m ./internal/services/cloud_load_balancer/... -v
```

**AI Agent Support:**
For generating, running, and debugging acceptance tests, use the `acctest` agent (`.opencode/agent/acctest.md`). This specialized agent can plan test strategies, generate complete test files, create sweeper implementations, and diagnose test failures.

### Test File Naming Convention

Following HashiCorp's official pattern from terraform-provider-aws, terraform-provider-google, and terraform-provider-azurerm:

| File Pattern | Type | Purpose | Who Writes |
|--------------|------|---------|------------|
| `*_schema_test.go` | Unit | Schema parity validation | Stainless (auto-generated, DO NOT modify) |
| `resource_test.go` | Acceptance | Resource CRUD tests | Custom (developer-written) |
| `data_source_test.go` | Acceptance | Data source tests | Custom (developer-written) |

**Key differences:**
- Unit tests run by default with `go test`
- Acceptance tests only run when `TF_ACC=1` is set
- Acceptance test functions start with `TestAcc` (e.g., `TestAccCloudVolume_basic`)
- Unit test functions do NOT start with `TestAcc`

### Shared Test Infrastructure

**Location:** `internal/acctest/`

Contains shared utilities used across all acceptance tests (following the AWS provider pattern):

**Files:**
- `provider.go` - Provider factories and PreCheck functions
- `helpers.go` - Test utility functions (RandomName, ProjectID, RegionID)

**Usage in acceptance tests:**
```go
import "github.com/stainless-sdks/gcore-terraform/internal/acctest"

resource.ParallelTest(t, resource.TestCase{
    PreCheck:                 func() { acctest.PreCheck(t) },
    ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
    // ...
})

// Use helper functions
rName := acctest.RandomName()
projectID := acctest.ProjectID()
regionID := acctest.RegionID()
```

### Sweepers (Resource Cleanup)

**Test sweepers automatically clean up leaked test resources** that weren't properly destroyed during acceptance tests.

**Location:** 
- `internal/sweep/sweep_test.go` - TestMain (enables sweeper CLI)
- `internal/sweep/sweep.go` - Shared utilities & validation
- `internal/sweep/framework.go` - Name filtering helpers
- `internal/services/*/sweep.go` - Per-resource sweeper implementations

**What sweepers delete:**
- ✅ Resources with `tf-test-*` or `tf_test_*` name prefixes
- ❌ **NOT** production resources or resources without test prefixes

**Run sweepers:**
```bash
# Clean up all test resources (WARNING: destroys infrastructure!)
go test ./internal/sweep -v -sweep=all -timeout 30m

# Clean up specific resource type only
go test ./internal/sweep -v -sweep=all -sweep-run=gcore_cloud_inference_deployment

# Allow failures (continue even if some deletions fail)
go test ./internal/sweep -v -sweep=all -sweep-allow-failures
```

**Important:**
- ⚠️ **Always use separate test accounts** - never run sweepers in production!
- ✅ **All test resources MUST use `acctest.RandomName()`** for sweepability
- 🤖 **Run sweepers periodically** (e.g., nightly in CI/CD) to prevent cost accumulation
- 📝 **Implement sweepers for any resource that creates infrastructure**

See `.opencode/agent/acctest.md` for detailed sweeper implementation guide.

## Code Style Guidelines

### Generated Code Header

All generated files include this header - do not remove it:
```go
// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.
```

### Import Organization

Group imports in three blocks separated by blank lines:
1. Standard library packages
2. External packages (third-party)
3. Internal packages (this module)

```go
import (
    "context"
    "fmt"
    "net/http"

    "github.com/G-Core/gcore-go"
    "github.com/hashicorp/terraform-plugin-framework/resource"

    "github.com/stainless-sdks/gcore-terraform/internal/apijson"
    "github.com/stainless-sdks/gcore-terraform/internal/logging"
)
```

### Naming Conventions

- **Packages**: snake_case matching resource names (e.g., `cloud_load_balancer`)
- **Types**: PascalCase (e.g., `CloudLoadBalancerModel`, `CloudLoadBalancerResource`)
- **Terraform resource names**: `gcore_<service_name>` (e.g., `gcore_cloud_load_balancer`)
- **Files per service**:
  - `resource.go` - CRUD operations
  - `data_source.go` - data source read operations
  - `schema.go` - Terraform schema definitions
  - `model.go` - data models and JSON transformations
  - `*_schema_test.go` - schema parity tests (Stainless-generated)
  - `resource_test.go` - acceptance tests (custom code)
  - `data_source_test.go` - data source acceptance tests (custom code)

### Struct Tags

Models use multiple struct tags for Terraform and JSON binding:
```go
type CloudLoadBalancerModel struct {
    ID        types.String `tfsdk:"id" json:"id,computed"`
    ProjectID types.Int64  `tfsdk:"project_id" path:"project_id,optional"`
    Name      types.String `tfsdk:"name" json:"name,optional"`
    CreatedAt timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
}
```

Tag modifiers: `required`, `optional`, `computed`, `computed_optional`, `no_refresh`

### Error Handling

Use Terraform diagnostics for error reporting:
```go
if err != nil {
    resp.Diagnostics.AddError("failed to make http request", err.Error())
    return
}

// for warnings
resp.Diagnostics.AddWarning("Resource not found", "removing from state")
```

### Resource Interface Compliance

Declare interface compliance explicitly:
```go
var _ resource.ResourceWithConfigure = (*CloudLoadBalancerResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudLoadBalancerResource)(nil)
var _ resource.ResourceWithImportState = (*CloudLoadBalancerResource)(nil)
```

### Testing Patterns

Tests use `t.Parallel()` and validate schema-model parity:
```go
func TestCloudLoadBalancerModelSchemaParity(t *testing.T) {
    t.Parallel()
    model := (*cloud_load_balancer.CloudLoadBalancerModel)(nil)
    schema := cloud_load_balancer.ResourceSchema(context.TODO())
    errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
    errs.Report(t)
}
```

## Architecture

```
internal/
├── provider.go              # main provider, resource/datasource registration
├── services/                # service implementations
│   └── cloud_*/             # one directory per terraform resource
│       ├── resource.go
│       ├── data_source.go
│       ├── schema.go
│       ├── model.go
│       ├── sweep.go          # sweeper implementation (custom code)
│       ├── *_schema_test.go  # schema parity tests (Stainless-generated)
│       ├── resource_test.go  # acceptance tests (custom code)
│       └── data_source_test.go  # data source acceptance tests (custom code)
├── acctest/                 # shared acceptance test infrastructure
│   ├── provider.go          # provider factories and PreCheck
│   └── helpers.go           # test utility functions
├── sweep/                   # sweeper infrastructure
│   ├── sweep_test.go        # TestMain for sweeper CLI
│   ├── sweep.go             # shared utilities & validation
│   └── framework.go         # name filtering helpers
├── apijson/                 # JSON encoding/decoding utilities
├── apiform/                 # form encoding utilities
├── customfield/             # custom Terraform field types
├── customvalidator/         # custom Terraform validators
├── importpath/              # import ID parsing
├── logging/                 # logging middleware
└── test_helpers/            # test utilities
```

## Environment Variables

```bash
GCORE_API_KEY              # required - API authentication key
GCORE_CLOUD_PROJECT_ID     # optional - default cloud project ID
GCORE_CLOUD_REGION_ID      # optional - default cloud region ID
GCORE_BASE_URL             # optional - override API base URL
```

## Git Policy

**Do NOT commit changes to git.** The user will handle all git operations themselves including:
- Staging files (`git add`)
- Creating commits (`git commit`)
- Pushing to remote (`git push`)
- Creating branches (`git checkout -b`)

Agents should only make code changes and leave version control to the user.

## Working with Stainless-Generated Code

### Do

- Run `./scripts/format` after any changes
- Write acceptance tests for new resources
- Use `internal/custom/` directory for truly custom code
- Report OpenAPI spec issues to have them fixed at the source
- Keep custom code isolated in separate files when possible

### Do Not

- Modify generated files extensively - changes may be overwritten
- Remove the "File generated from our OpenAPI spec" header
- Change import organization style in generated files
- Rename generated types or packages
- Force push to integrated branches (`main`, `next`)

### Custom Code Guidelines

If custom code is absolutely necessary:
1. Prefer creating new files over modifying generated ones
2. Use conventional commit messages (e.g., `feat(client): add helper method`)
3. Document why custom code was needed vs fixing the OpenAPI spec
4. Be prepared to resolve merge conflicts on regeneration

## Plan Modifiers

See the `plan-modifiers` skill (`.agents/skills/plan-modifiers/SKILL.md`) for conventions on creating and organizing custom Terraform plan modifiers.

## Local Development

1. Build: `./scripts/build`
2. Configure `~/.terraformrc`:
   ```hcl
   provider_installation {
     dev_overrides {
       "stainless-sdks/gcore" = "/path/to/this/repo"
     }
     direct {}
   }
   ```
3. Test with your `.tf` files using `terraform apply`

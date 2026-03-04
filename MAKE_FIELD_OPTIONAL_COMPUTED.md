# Making Terraform Fields Optional and Computed

This guide explains how to configure fields in the Terraform provider to be `computed_optional` (user-provided OR API-computed).

## When to Use computed_optional

Use `computed_optional` for fields that:
- Can be provided by the user during resource creation
- Will be computed by the API if not provided
- Need to be preserved across updates (not force replacement)

Example: `port_id` in `gcore_cloud_reserved_fixed_ip` resource

## Step 1: Identify the Correct Model

⚠️ **Critical**: You must add the configuration to the **INPUT schema**, not the OUTPUT schema.

### How to Find the Right Model

1. Look at the OpenAPI spec operation for resource creation:
   ```bash
   grep -A 50 "operationId.*ReservedFixedIP.*create" openapi.yaml
   ```

2. Find the `requestBody` schema (this is the INPUT schema):
   ```yaml
   requestBody:
     content:
       application/json:
         schema:
           oneOf:
             - $ref: '#/components/schemas/NewReservedFixedIpExternalSerializer'
             - $ref: '#/components/schemas/NewReservedFixedIpSpecificPortSerializer'
   ```

3. For `oneOf` discriminated unions, check which schema contains your field:
   ```bash
   grep -A 100 "NewReservedFixedIpSpecificPortSerializer:" openapi.yaml | grep "port_id"
   ```

### Example: Reserved Fixed IP

For `gcore_cloud_reserved_fixed_ip`:
- ❌ **Wrong**: `ReservedFixedIPSerializer` (this is the OUTPUT/response schema)
- ✅ **Correct**: `NewReservedFixedIpSpecificPortSerializer` (this is the INPUT/request schema)

## Step 2: Add Configuration to api-schemas

Edit `api-schemas/scripts/config.yaml` and add the field configuration:

```yaml
terraform_resource_schemas:
  ReservedFixedIP:
    NewReservedFixedIpSpecificPortSerializer:
      properties:
        port_id: computed_optional
```

### Configuration Structure

```yaml
terraform_resource_schemas:
  <ResourceName>:                    # Terraform resource name (without gcore_cloud_ prefix)
    <InputSchemaName>:               # INPUT schema from requestBody
      properties:
        <field_name>: computed_optional
```

### Common Patterns

**For simple resources** (single input schema):
```yaml
terraform_resource_schemas:
  Volume:
    CreateVolumeSerializer:
      properties:
        size: computed_optional
```

**For oneOf resources** (multiple input schemas):
```yaml
terraform_resource_schemas:
  ReservedFixedIP:
    NewReservedFixedIpExternalSerializer:
      properties:
        network_id: computed_optional
    NewReservedFixedIpSpecificPortSerializer:
      properties:
        port_id: computed_optional
```

## Step 3: CI Regenerates OpenAPI

After merging the config.yaml changes:

1. CI pipeline automatically regenerates `openapi.yaml`
2. The `x-stainless-terraform-configurability: computed_optional` attribute is added to the field
3. Changes are committed to the `api-schemas` repository

Example result in `openapi.yaml`:
```yaml
NewReservedFixedIpSpecificPortSerializer:
  properties:
    port_id:
      description: Port ID to make a reserved fixed IP
      format: uuid4
      type: string
      x-stainless-terraform-configurability: computed_optional
```

## Step 4: Trigger Stainless Codegen

After the OpenAPI spec is updated:

1. Stainless CI detects the change and runs code generation
2. New code is committed to `gcore-terraform` main branch
3. The generated code will include:
   - Schema with `Computed: true, Optional: true`
   - Model with `computed_optional` JSON tag
   - Plan modifier `RequiresReplaceIfConfigured()`

Example generated code:
```go
// schema.go
"port_id": schema.StringAttribute{
    Description:   "Port ID to make a reserved fixed IP",
    Computed:      true,
    Optional:      true,
    PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
}

// model.go
type CloudReservedFixedIPModel struct {
    PortID types.String `tfsdk:"port_id" json:"port_id,computed_optional"`
}
```

## Step 5: Rebase Your Branch

After Stainless codegen completes:

```bash
# Fetch latest changes
git fetch origin main

# Rebase your branch on updated main
git rebase origin/main

# Resolve conflicts if needed (usually accept main's version for generated files)
git rebase --continue

# Force push (if needed)
git push --force-with-lease
```

## Troubleshooting

### Issue: Field not showing as computed_optional

**Check**: Are you using the INPUT schema or OUTPUT schema?
```bash
# Find which schema is used in requestBody (INPUT)
grep -B 10 -A 10 "operationId.*<Resource>.*create" openapi.yaml | grep -A 5 requestBody

# vs. response schema (OUTPUT) - this is NOT what you want
grep -B 10 -A 10 "operationId.*<Resource>.*create" openapi.yaml | grep -A 5 responses
```

### Issue: Multiple schemas in oneOf

For resources with discriminated unions (oneOf), you may need to add the configuration to multiple schemas:

```yaml
terraform_resource_schemas:
  ReservedFixedIP:
    NewReservedFixedIpExternalSerializer:
      properties:
        network_id: computed_optional
    NewReservedFixedIpSpecificPortSerializer:
      properties:
        port_id: computed_optional
```

### Issue: Rebase conflicts

For generated files (model.go, schema.go, resource.go), prefer main's version:
```bash
git checkout --theirs <file>
```

Only keep your custom changes in resource.go if they're critical bug fixes.

## Related Files

- `api-schemas/scripts/config.yaml` - Field configuration
- `api-schemas/openapi.yaml` - Generated OpenAPI spec
- `gcore-terraform/internal/services/cloud_<resource>/schema.go` - Generated schema
- `gcore-terraform/internal/services/cloud_<resource>/model.go` - Generated model

## Example: Complete Workflow

This was done for `port_id` in `gcore_cloud_reserved_fixed_ip`:

1. ✅ Identified INPUT schema: `NewReservedFixedIpSpecificPortSerializer`
2. ✅ Added to config.yaml:
   ```yaml
   NewReservedFixedIpSpecificPortSerializer:
     properties:
       port_id: computed_optional
   ```
3. ✅ CI regenerated openapi.yaml with `x-stainless-terraform-configurability`
4. ✅ Stainless codegen created correct Terraform schema (commit `09efce4a`)
5. ✅ Rebased branch on main to get generated code
6. ✅ Verified field works as computed_optional in Terraform

## Key Takeaways

- **Always use INPUT schema** (requestBody), never OUTPUT schema (response)
- **Check for oneOf** discriminated unions - you may need multiple schema entries
- **Let CI/CD do the work** - config.yaml → openapi.yaml → Stainless codegen
- **Rebase after codegen** to get the generated code into your branch

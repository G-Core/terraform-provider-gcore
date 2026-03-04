# Drift Fix Patterns

## Common Drift Issue Pattern

When a resource shows perpetual drift (changes detected on every plan after apply), it's typically due to incorrect handling of computed fields.

## Root Causes

1. **Wrong Unmarshaler**: Using `apijson.Unmarshal` instead of `apijson.UnmarshalComputed`
2. **Missing Tags**: Fields not marked as `computed_optional` in model
3. **Schema Mismatch**: Schema fields not marked with `Computed: true`

## The Fix Pattern

### Step 1: Update Resource Read Method

```go
// ❌ WRONG - ignores computed values
err = apijson.Unmarshal(bytes, &data)

// ✅ CORRECT - handles computed values
err = apijson.UnmarshalComputed(bytes, &data)
```

### Step 2: Update ImportState Method

```go
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    // ... fetch resource ...
    
    // ❌ WRONG
    err = apijson.Unmarshal(bytes, &data)
    
    // ✅ CORRECT
    err = apijson.UnmarshalComputed(bytes, &data)
}
```

### Step 3: Update Model Tags

For fields that can be either user-provided OR API-computed:

```go
// ❌ WRONG - doesn't handle API defaults
type PoolModel struct {
    HTTPMethod     types.String `json:"http_method,optional"`
    MaxRetriesDown types.Int64  `json:"max_retries_down,optional"`
}

// ✅ CORRECT - handles both user and API values
type PoolModel struct {
    HTTPMethod     types.String `json:"http_method,computed_optional"`
    MaxRetriesDown types.Int64  `json:"max_retries_down,computed_optional"`
}
```

### Step 4: Update Schema

```go
// ❌ WRONG - missing Computed flag
"http_method": schema.StringAttribute{
    Optional: true,
    Validators: []validator.String{...},
}

// ✅ CORRECT - marked as computed
"http_method": schema.StringAttribute{
    Computed: true,  // CRITICAL: Add this!
    Optional: true,
    Validators: []validator.String{...},
}
```

## Quick Checklist

For each field that shows drift:

- [ ] Is it using `UnmarshalComputed` in Read?
- [ ] Is it using `UnmarshalComputed` in ImportState?
- [ ] Is the model field tagged `computed_optional`?
- [ ] Is the schema field marked `Computed: true`?
- [ ] After Update operations, is response unmarshaled with `UnmarshalComputed`?

## Testing the Fix

```bash
# Build the provider
go build -o terraform-provider-gcore

# Clean slate
terraform destroy -auto-approve

# Apply and immediately check for drift
terraform apply -auto-approve
terraform plan -detailed-exitcode

# Exit code 0 = SUCCESS (no drift)
# Exit code 2 = FAILURE (drift detected)
```

## Common Fields That Need This Pattern

- Health monitor settings (http_method, max_retries_down)
- Default ports and protocols
- Auto-generated names or IDs
- Timeout values with defaults
- Any field the API computes when not provided
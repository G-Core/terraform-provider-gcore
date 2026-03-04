# Why the Manual Fix Failed

## The Problem

When we manually added `Computed: true` to the schema fields, we encountered this error:

```
Error: Value Conversion Error

Received unknown value, however the target type cannot handle unknown values.
Use the corresponding `types` package type or a custom type that handles unknown values.

Path: sni_secret_id
Target Type: *[]basetypes.StringValue
Suggested Type: basetypes.ListValue
```

## Root Cause

The issue occurs because manually adding `Computed: true` to a `ListAttribute` requires coordinated changes across multiple files:

### 1. Schema File (schema.go)
```go
"sni_secret_id": schema.ListAttribute{
    Computed:    true,  // ← We added this
    Optional:    true,
    ElementType: types.StringType,
}
```

### 2. Model File (model.go) - MISMATCH!
```go
type CloudLoadBalancerListenerModel struct {
    SniSecretID *[]types.String `tfsdk:"sni_secret_id" json:"sni_secret_id,optional"`
    //          ^^^^^^^^^^^^^^^
    // This type doesn't handle unknown values!
    // Should be: types.List
}
```

### 3. JSON Tag - MISMATCH!
The JSON tag says `optional` but schema says `computed_optional`:
```go
SniSecretID *[]types.String `json:"sni_secret_id,optional"`
//                                                ^^^^^^^^
// Should be: computed_optional
```

## Why It Matters

When a field is `Computed: true` and `Optional: true`, Terraform may receive "unknown" values during the plan phase. The model type `*[]types.String` cannot represent unknown values, but `types.List` can.

### Proper Type for Computed Optional Lists
```go
// ❌ Wrong (what we have now):
SniSecretID *[]types.String `json:"sni_secret_id,optional"`

// ✅ Correct (what Stainless will generate):
SniSecretID types.List `json:"sni_secret_id,computed_optional"`
```

## Similar Issues with Other Fields

The same problem applies to:

1. **sni_secret_id** - `*[]types.String` → should be `types.List`
2. **user_list** - `*[]*CloudLoadBalancerListenerUserListModel` → needs custom handling
3. **timeout_client_data** - `types.Int64` with `optional` tag → needs `computed_optional` tag
4. **timeout_member_connect** - Same issue
5. **timeout_member_data** - Same issue

## Why Stainless Code Generation Is Required

Stainless codegen handles all the necessary changes:

1. **Updates schema.go**: Adds `Computed: true`
2. **Updates model.go**: Changes type to handle unknown values
3. **Updates JSON tags**: Changes from `optional` to `computed_optional`
4. **Updates plan modifiers**: May add `RequiresReplaceIfConfigured()`
5. **Updates MarshalJSONForUpdate**: Handles computed_optional fields correctly

## The Proper Solution

✅ **Use the CI/CD pipeline** as documented in `PROPER_FIX.md`:

1. Update `api-schemas/scripts/config.yaml` ← **DONE**
2. CI regenerates `openapi.yaml` with `x-stainless-terraform-configurability`
3. Stainless generates all required code changes
4. Rebase branch to get generated code

## Current Status

- ✅ config.yaml updated with PatchLbListenerSerializer configuration
- ❌ Manual schema.go changes reverted (they caused errors)
- ⏳ Waiting for CI/CD pipeline to generate proper code

## Lessons Learned

1. **Don't manually edit generated files** - It's easy to miss coordinated changes
2. **List/slice types need special handling** for unknown values
3. **JSON tags must match schema definition** (`optional` vs `computed_optional`)
4. **Trust the code generator** - It knows all the places that need updates

## Reference

- Error occurred in: `internal/services/cloud_load_balancer_listener/schema.go:112`
- Model definition: `internal/services/cloud_load_balancer_listener/model.go:26`
- OpenAPI schema: `PatchLbListenerSerializer` in `api-schemas/openapi.yaml`

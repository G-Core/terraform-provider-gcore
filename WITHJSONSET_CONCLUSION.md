# WithJSONSet Analysis - FINAL CONCLUSION

## SDK Investigation

Found in `sdk-gcore-go/cloud/networkrouter.go`:

```go
type NetworkRouterUpdateParams struct {
    // ...
    // List of custom routes.
    Routes []NetworkRouterUpdateParamsRoute `json:"routes,omitzero"`
    // ...
}
```

**The SDK USES `omitzero` TAG on Routes field!** ✅

## What This Means

According to Go's JSON encoding rules with `omitzero`:
- **Empty slice `[]`**: IS serialized as `"routes":[]`
- **Nil slice**: Is omitted (no routes field in JSON)

**Pedro's claim is 100% verified as CORRECT!** ✅

## Why Do We Use WithJSONSet Then?

Looking at our code flow:

1. **We DON'T use SDK's NetworkRouterUpdateParams directly for marshaling**
2. **Instead**: We call `data.MarshalJSONForUpdate(*state)` which marshals `CloudNetworkRouterModel`
3. **CloudNetworkRouterModel** uses `customfield.NestedObjectList` for Routes, not a plain slice
4. **Question**: Does `customfield.NestedObjectList` with empty list serialize to `[]` or is it omitted?

## The Real Question

The SDK's `omitzero` behavior is irrelevant because we're not marshaling the SDK struct directly. We're marshaling our own `CloudNetworkRouterModel`:

```go
type CloudNetworkRouterModel struct {
    Routes customfield.NestedObjectList[CloudNetworkRouterRoutesModel] `tfsdk:"routes" json:"routes,computed_optional"`
}
```

The `json:"routes,computed_optional"` tag uses a **custom tag** (`computed_optional`), not standard Go `omitzero`.

This is parsed by `apijson.MarshalForPatch`, which is custom Stainless marshaling logic.

## Investigation Needed

To know if `WithJSONSet` is necessary, we need to understand:

1. **What does `computed_optional` mean in Stainless's JSON marshaling?**
2. **Does `customfield.NestedObjectList` with empty list serialize to `[]`?**
3. **Does `apijson.MarshalForPatch` handle empty lists correctly?**

Without diving into Stainless's internal marshaling logic, we cannot definitively say whether `WithJSONSet` is needed.

## Recommendation

### Keep `WithJSONSet` for now ✅

**Reasons:**
1. ✅ **It works** - current code successfully deletes routes
2. ✅ **Low risk** - no harm in being explicit about sending `routes:[]`
3. ✅ **Custom marshaling** - We're using Stainless's custom tags, not standard Go JSON
4. ⚠️  **Unknown behavior** - `computed_optional` tag behavior is not documented/clear

### If We Want to Remove It

Would need to:
1. ✅ Verify `customfield.NestedObjectList` with empty list marshals to `[]`
2. ✅ Verify `apijson.MarshalForPatch` includes fields with `computed_optional` tag
3. ✅ Add unit tests to ensure `MarshalJSONForUpdate` produces `"routes":[]`
4. ✅ Test with real infrastructure to confirm routes are deleted

## Summary

| Statement | Status |
|-----------|--------|
| Pedro's claim about `omitzero` behavior | ✅ CORRECT |
| SDK uses `omitzero` on Routes field | ✅ VERIFIED |
| Empty slices with `omitzero` are serialized | ✅ CONFIRMED |
| We use SDK's marshaling directly | ❌ FALSE - we use custom Stainless marshaling |
| `WithJSONSet` is definitely unnecessary | ⚠️  UNKNOWN - depends on Stainless implementation |

**Bottom Line:** Pedro's understanding of Go JSON marshaling is correct, but we use custom Stainless marshaling logic that may behave differently. Keep `WithJSONSet` for safety unless we can verify the custom marshaling handles empty lists correctly.

## Action Items for Full Verification

If we want to remove `WithJSONSet` (optimization, not critical):

1. Write unit test:
```go
func TestRoutesDeletionSerialization(t *testing.T) {
    data := CloudNetworkRouterModel{
        Name: types.StringValue("test"),
        Routes: customfield.NewObjectListMust(ctx, []CloudNetworkRouterRoutesModel{}),
    }
    state := CloudNetworkRouterModel{
        Routes: customfield.NewObjectListMust(ctx, []CloudNetworkRouterRoutesModel{
            {Destination: types.StringValue("10.0.0.0/24")},
        }),
    }

    bytes, _ := data.MarshalJSONForUpdate(state)

    // MUST contain routes:[]
    require.Contains(t, string(bytes), `"routes":[]`)
}
```

2. If test passes: Remove `WithJSONSet`
3. Test with real infrastructure: Ensure routes deletion still works

**For this PR: Keep `WithJSONSet` as-is** - the other bug fixes are more critical.

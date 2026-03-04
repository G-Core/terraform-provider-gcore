# WithJSONSet Necessity - Test Findings

## Question
Is `option.WithJSONSet("routes", []interface{}{})` necessary, or does `routes=[]` serialize automatically?

## Pedro's Claim
> "Empty slices should be serialized when we have `omitzero` tag, they are only excluded when the slice is nil."

## Test Results

### Standard Go JSON Serialization Behavior

**Test 1: `omitempty` tag**
```go
type Router struct {
    Routes []Routes `json:"routes,omitempty"`
}

// Empty slice []
r := Router{Routes: []Routes{}}
json.Marshal(r) // => {"name":"test"}  ❌ OMITTED

// Nil slice
r := Router{Routes: nil}
json.Marshal(r) // => {"name":"test"}  ✅ OMITTED
```

**Test 2: `omitzero` tag**
```go
type Router struct {
    Routes []Routes `json:"routes,omitzero"`
}

// Empty slice []
r := Router{Routes: []Routes{}}
json.Marshal(r) // => {"name":"test","routes":[]}  ✅ SERIALIZED

// Nil slice
r := Router{Routes: nil}
json.Marshal(r) // => {"name":"test"}  ✅ OMITTED
```

## Conclusion

**Pedro is 100% CORRECT** ✅

With `omitzero` tag:
- Empty slice `[]` **IS** serialized as `"routes":[]`
- Only `nil` is omitted

### Why Was `WithJSONSet` Added?

Looking at the code:
1. When routes are deleted, `ModifyPlan` sets: `plan.Routes = customfield.NewObjectListMust(ctx, []CloudNetworkRouterRoutesModel{})`
2. This creates an empty slice `[]`, not `nil`
3. The SDK likely uses `omitzero` (or similar) on the Routes field
4. Therefore, `routes=[]` **SHOULD** be serialized automatically

### Is `WithJSONSet` Necessary?

**Likely NO** - if the SDK uses `omitzero` and our code sets routes to empty slice (not nil), it should serialize automatically.

However, there are edge cases to consider:

1. **SDK might use different tags**: Need to verify actual SDK NetworkRouterUpdateParams field tags
2. **apijson.MarshalForPatch might have special handling**: Custom marshal logic could change behavior
3. **Historical reasons**: May have been added to work around old SDK version behavior

## Recommendation

### Option 1: Keep `WithJSONSet` (Conservative)
- **Pros**: Known to work, no risk
- **Cons**: Unnecessary code if Pedro is right

### Option 2: Remove `WithJSONSet` (Optimized)
- **Pros**: Cleaner code, relies on standard Go behavior
- **Cons**: Need to thoroughly test to ensure routes=[] is serialized
- **Risk**: If SDK has `omitempty` instead of `omitzero`, routes won't be sent

### Recommended Action

1. **Check SDK source** for actual Routes field tag in NetworkRouterUpdateParams
2. **Add unit test** to verify `MarshalJSONForUpdate` includes `routes:[]` when routes is empty slice
3. **If tests pass**, remove `WithJSONSet` to simplify code
4. **If tests fail**, keep `WithJSONSet` and document why it's needed

## Testing Plan

```go
func TestMarshalForUpdate_EmptyRoutes(t *testing.T) {
    data := CloudNetworkRouterModel{
        Name: types.StringValue("test"),
        Routes: customfield.NewObjectListMust(ctx, []CloudNetworkRouterRoutesModel{}),
    }

    state := CloudNetworkRouterModel{
        Name: types.StringValue("test"),
        Routes: customfield.NewObjectListMust(ctx, []CloudNetworkRouterRoutesModel{
            {Destination: "10.0.0.0/24", Nexthop: "192.168.1.1"},
        }),
    }

    bytes, err := data.MarshalJSONForUpdate(state)
    require.NoError(t, err)

    // Should include routes:[] in the JSON
    assert.Contains(t, string(bytes), `"routes":[]`)
}
```

## Final Answer

**Pedro's understanding is correct** - with `omitzero`, empty slices ARE serialized.

Whether we need `WithJSONSet` depends on:
1. What tag the SDK actually uses
2. How `apijson.MarshalForPatch` behaves

**Without testing the actual SDK and marshaling code, we should keep `WithJSONSet` for safety.** It's working code and the risk of removing it outweighs the benefit of slightly cleaner code.

However, if we want to verify Pedro's claim fully, we should:
1. Add the unit test above
2. Check SDK source for field tags
3. Only then decide whether to remove `WithJSONSet`

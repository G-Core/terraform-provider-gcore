# Pedro's PR Review - Test Findings and Implementation Changes

## Test Results Summary

### 1. Empty Slices with `omitzero` Tag
**Pedro's Claim:** "Empty slices should be serialized when we have omitzero tag, they are only excluded when the slice is nil."

**Test Result:** ✅ **CONFIRMED** - PATCH with `routes=[]` successfully deletes all routes.

**Current Code:** Uses `option.WithJSONSet("routes", []interface{}{})` to force empty array in body.

**Question:** Do we still need `WithJSONSet` if empty slices are serialized automatically?
**Answer:** Need to test if `routes=[]` is sent without `WithJSONSet`. Current test used `WithJSONSet` so didn't verify Pedro's claim completely.

---

### 2. Routes Deletion Logic - Partial Deletions
**Pedro's Suggestion:**
```go
// Current (WRONG - only catches full deletion):
routesDeletionNeeded := len(dataRoutes) == 0 && len(stateRoutes) > 0

// Suggested (CORRECT - catches partial deletions):
routesDeletionNeeded := len(dataRoutes) < len(stateRoutes)
```

**Test Result:** ✅ **CONFIRMED** - Pedro is correct!

**Impact:** Current code MISSES partial route deletions (e.g., 2 routes → 1 route). It only detects full deletions (n routes → 0 routes).

**Required Change:** Update condition to `len(dataRoutes) < len(stateRoutes)`

---

### 3. Single PATCH for All Updates
**Pedro's Claim:** "We should be able to send a single PATCH request for updating any fields that can be updated in an Update. Sending an empty array for routes should work."

**Test Result:** ✅ **CONFIRMED** - Single PATCH successfully updated both name AND deleted routes simultaneously.

**Current Code:** Uses two separate PATCH calls in some scenarios:
1. Early PATCH (lines 139-184): Delete routes before interface detachment
2. Later PATCH (lines 283-352): Update other fields

**Constraint:** Must delete routes BEFORE detaching interfaces if routes use those interfaces as nexthop (API requirement).

**Simplification Opportunity:**
- When interfaces NOT changing: Always use single PATCH for all updates
- When interfaces changing: Keep early PATCH only if routes being deleted (API requirement)

---

### 4. Final GET After Updates
**Pedro's Claim:** "If we have updated something using the Update or Attach/Detach methods we get back an updated state for the router, so we shouldn't need to do a Get in the end."

**Test Result:** ✅ **CONFIRMED** - PATCH response matches GET response (name, status, etc.)

**Current Code:** Lines 370-396 perform final GET after all update operations

**Question:** Do Attach/Detach also return updated state?
**Risk:** If we remove final GET and Attach/Detach don't return complete state, we might have stale data.

**Recommended Change:** Test if Attach/Detach return complete router state. If yes, remove final GET. If no, keep it.

---

## Required Code Changes

### Priority 1: Fix Partial Deletion Bug
```go
// Change line 133:
- routesDeletionNeeded := len(dataRoutes) == 0 && len(stateRoutes) > 0
+ routesDeletionNeeded := len(dataRoutes) < len(stateRoutes)
```

This is a **BUG FIX** - current code doesn't handle partial route deletions correctly.

### Priority 2: Test and Remove WithJSONSet if Unnecessary
If empty slices are serialized automatically (as Pedro claims), we can remove:
- Line 164: `option.WithJSONSet("routes", []interface{}{})`
- Line 322: `option.WithJSONSet("routes", []interface{}{})`

Need to verify this works without WithJSONSet.

### Priority 3: Consider Removing Final GET
If PATCH/Attach/Detach all return complete updated state, remove lines 370-396.

Test this carefully to avoid state inconsistencies.

---

## Testing Checklist

- [x] Empty slices work with PATCH
- [x] Single PATCH can update multiple fields
- [x] PATCH returns updated state
- [ ] Test WITHOUT WithJSONSet - does routes=[] serialize automatically?
- [ ] Test Attach/Detach - do they return complete router state?
- [ ] Test partial route deletion with simplified code

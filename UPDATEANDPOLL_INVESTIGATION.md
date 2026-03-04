# Investigation: Adding UpdateAndPoll for LoadBalancer

## Question
Can we add `UpdateAndPoll` method for LoadBalancers in the SDK to properly solve the vrrp_ips bug?

## Findings

### 1. Existing *AndPoll Methods in SDK

**From PR #46 (MERGED)**: Added polling methods for LoadBalancers

For **LoadBalancerService**:
- ✅ `NewAndPoll` - Create and poll
- ✅ `DeleteAndPoll` - Delete and poll
- ✅ `FailoverAndPoll` - Failover and poll
- ✅ `ResizeAndPoll` - Resize and poll
- ❌ **`UpdateAndPoll` - NOT ADDED**

For **LoadBalancerListenerService** and **LoadBalancerPoolService**:
- ✅ `NewAndPoll`
- ✅ `UpdateAndPoll` ← **These have UpdateAndPoll!**
- ✅ `DeleteAndPoll`

### 2. Key Difference: Update Response Types

From `/Users/user/repos/sdk-gcore-go/cloud/loadbalancer.go`:

```go
// LoadBalancer Update - returns LoadBalancer directly (line 83)
func (r *LoadBalancerService) Update(ctx context.Context, loadBalancerID string,
    params LoadBalancerUpdateParams, opts ...option.RequestOption) (res *LoadBalancer, err error)

// LoadBalancer New - returns TaskIDList (line 58)
func (r *LoadBalancerService) New(ctx context.Context,
    params LoadBalancerNewParams, opts ...option.RequestOption) (res *TaskIDList, err error)

// LoadBalancer Delete - returns TaskIDList (line 146)
func (r *LoadBalancerService) Delete(ctx context.Context, loadBalancerID string,
    body LoadBalancerDeleteParams, opts ...option.RequestOption) (res *TaskIDList, err error)
```

**Analysis:**
- `New()` → Returns `*TaskIDList` (async, has tasks to poll)
- `Update()` → Returns `*LoadBalancer` (no tasks!)
- `Delete()` → Returns `*TaskIDList` (async, has tasks to poll)

### 3. How Listener/Pool UpdateAndPoll Works

From PR #46 diff:

```go
func (r *LoadBalancerListenerService) UpdateAndPoll(ctx context.Context, listenerID string,
    params LoadBalancerListenerUpdateParams, opts ...option.RequestOption) (v *LoadBalancerListenerDetail, err error) {

    // Step 1: Call Update - expects it to return a resource with Tasks[]
    resource, err := r.Update(ctx, listenerID, params, opts...)
    if err != nil {
        return
    }

    // Step 2: Extract task_id
    if len(resource.Tasks) != 1 {
        return nil, errors.New("expected exactly one task to be created")
    }
    taskID := resource.Tasks[0]

    // Step 3: Poll the task
    _, err = r.tasks.Poll(ctx, taskID, opts...)
    if err != nil {
        return
    }

    // Step 4: GET to retrieve final state
    return r.Get(ctx, listenerID, getParams, opts...)
}
```

**Key Insight:** Listener and Pool `Update()` methods return a response that includes `Tasks[]` field.

### 4. The Problem

**LoadBalancer Update API does NOT return tasks in the SDK**, so we cannot follow the same pattern as Listener/Pool UpdateAndPoll.

Two possibilities:
1. **The API actually doesn't return tasks** for LB updates (synchronous operation)
2. **The SDK generation is incorrect** - Stainless didn't recognize that the API returns tasks

## Solutions

### Solution 1: Check if API Returns Tasks (Needs Verification)

**Test with mitmproxy:**
```bash
# Run mitmproxy
mitmproxy --mode reverse:https://api.gcore.com@443 -p 8080

# Configure terraform to use proxy
export HTTPS_PROXY=http://localhost:8080

# Run LB rename operation
terraform apply -var="lb_name=test-RENAMED"

# Check mitmproxy logs for PATCH response
```

**If API DOES return tasks:**
- Fix Stainless OpenAPI spec to include `tasks` in response
- Regenerate SDK
- Add UpdateAndPoll following Listener/Pool pattern

**If API does NOT return tasks:**
- Current workaround (Update + GET) is correct
- Cannot add UpdateAndPoll in the proper sense
- Could add a convenience method that just does Update + GET

### Solution 2: Add Convenience Method (Works Now)

Even if Update doesn't return tasks, we can add a convenience method:

```go
// UpdateAndRefresh updates a load balancer and refreshes state
// Note: Unlike other *AndPoll methods, this doesn't poll a task because
// the Load Balancer Update API returns the updated resource directly.
func (r *LoadBalancerService) UpdateAndRefresh(ctx context.Context,
    loadBalancerID string, params LoadBalancerUpdateParams,
    opts ...option.RequestOption) (v *LoadBalancer, err error) {

    // Step 1: Call Update
    _, err = r.Update(ctx, loadBalancerID, params, opts...)
    if err != nil {
        return
    }

    // Step 2: Extract project/region from params
    opts = slices.Concat(r.Options, opts)
    precfg, err := requestconfig.PreRequestOptions(opts...)
    if err != nil {
        return
    }

    var getParams LoadBalancerGetParams
    requestconfig.UseDefaultParam(&params.ProjectID, precfg.CloudProjectID)
    requestconfig.UseDefaultParam(&params.RegionID, precfg.CloudRegionID)
    getParams.ProjectID = params.ProjectID
    getParams.RegionID = params.RegionID

    // Step 3: GET to retrieve fresh state with all computed fields
    return r.Get(ctx, loadBalancerID, getParams, opts...)
}
```

**Pros:**
- Works immediately without API/spec changes
- Matches the workaround we already implemented in Terraform provider
- Clear naming (`UpdateAndRefresh` vs `UpdateAndPoll`) indicates no task polling

**Cons:**
- Not a true "AndPoll" method (doesn't poll tasks)
- Inconsistent with Listener/Pool if they DO poll tasks

### Solution 3: Verify API Behavior First

**Recommended approach:**

1. **Test actual API response** with curl or mitmproxy:
   ```bash
   curl -X PATCH "https://api.gcore.com/cloud/v1/loadbalancers/PROJECT/REGION/LB_ID" \
     -H "Authorization: APIKey YOUR_KEY" \
     -H "Content-Type: application/json" \
     -d '{"name":"new-name"}' | jq
   ```

2. **Check if response includes:**
   - `tasks: ["task-uuid"]` → Can implement proper UpdateAndPoll
   - No `tasks` field → Use UpdateAndRefresh convenience method

3. **If tasks are returned:**
   - Update OpenAPI spec at `/Users/user/repos/gcore-terraform/api-schemas/openapi.yaml`
   - Regenerate SDK with Stainless
   - Add UpdateAndPoll to SDK
   - Update Terraform provider to use UpdateAndPoll

4. **If no tasks:**
   - Add `UpdateAndRefresh` convenience method to SDK
   - Update Terraform provider to use UpdateAndRefresh
   - Document why it's different from other *AndPoll methods

## Current Workaround in Terraform Provider

The fix we implemented at `/Users/user/repos/gcore-terraform/internal/services/cloud_load_balancer/resource.go:209-247`:

```go
// Call Update - SDK handles async internally (may or may not poll)
_, err = r.client.Cloud.LoadBalancers.Update(ctx, data.ID.ValueString(), params, ...)

// After update, do explicit GET to refresh all computed fields
res := new(http.Response)
_, err = r.client.Cloud.LoadBalancers.Get(ctx, data.ID.ValueString(), getParams,
    option.WithResponseBodyInto(&res),
    option.WithMiddleware(logging.Middleware(ctx)),
)
bytes, _ := io.ReadAll(res.Body)
err = apijson.UnmarshalComputed(bytes, &data)
```

This works but:
- ✅ Fixes the vrrp_ips bug
- ❌ Requires manual implementation in Terraform provider
- ❌ Not reusable across different SDK consumers
- ❌ Would be better as SDK method

## Recommendation

1. **Immediate:** Test API with curl/mitmproxy to see if PATCH returns tasks
2. **If tasks returned:** Update OpenAPI spec, regenerate SDK, add UpdateAndPoll
3. **If no tasks:** Add `UpdateAndRefresh` convenience method to SDK
4. **Update Terraform provider** to use the new SDK method

This would:
- Centralize the logic in SDK (single source of truth)
- Make it reusable for other SDK consumers
- Follow SDK patterns (all *AndPoll/AndRefresh methods in one place)
- Be more maintainable than duplicating logic in Terraform provider

## Next Steps

**User, please advise:**
1. Should I test the actual API response with curl?
2. Or would you prefer to add the convenience method to SDK regardless?
3. Do you have access to Stainless config to regenerate SDK if needed?

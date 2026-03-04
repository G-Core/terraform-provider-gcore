# Reply: Would AddAndPoll Returning the Member Be Helpful in Terraform?

## Short Answer

**Yes, absolutely!** Returning `*Member` from `AddAndPoll` would be very helpful in Terraform and is the correct pattern for this use case.

## Why Return `*Member`?

### 1. Members Are Stateful Resources, Not Pure Associations

Load balancer members are fundamentally different from something like volume attachments:

| Characteristic | Volume Attachment | LB Pool Member |
|----------------|-------------------|----------------|
| Creates new resource? | ❌ No (just a link) | ✅ Yes (member with ID) |
| Has its own state? | ❌ No | ✅ Yes (`operating_status`, `provisioning_status`) |
| State changes after creation? | ❌ Volume properties unchanged | ✅ Yes (health checks update `operating_status`) |
| Terraform needs state? | ❌ No | ✅ **Yes!** (for computed fields) |

**Conclusion**: Members behave like **resources with lifecycle**, not pure associations.

### 2. Terraform Requires the Member Data Immediately

Looking at the current Terraform implementation with manual polling:

```go
// In Create() - current implementation
func (r *CloudLoadBalancerPoolMemberResource) Create(ctx context.Context, req, resp) {
    // Step 1: Add member
    taskIDList, err := r.client.Cloud.LoadBalancers.Pools.Members.Add(ctx, poolID, params)

    // Step 2: Poll task
    task, err := r.client.Cloud.Tasks.Poll(ctx, taskIDList.Tasks[0])

    // Step 3: Extract member ID
    memberID := task.CreatedResources.Members[0]
    data.ID = types.StringValue(memberID)

    // Step 4: MUST fetch member details to get computed fields!
    if err := r.readMemberFromPool(ctx, data); err != nil {
        resp.Diagnostics.AddWarning("failed to read member details after creation", err.Error())
    }
    // readMemberFromPool() does:
    // - Calls Pools.Get() to fetch the entire pool
    // - Searches through pool.Members[] to find our member
    // - Populates: operating_status, subnet_id, admin_state_up, backup, weight, etc.

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

**The problem**: After getting the member ID from the task, Terraform still needs to make **another API call** to get:
- `operating_status` (ONLINE, OFFLINE, ERROR, NO_MONITOR, etc.)
- `provisioning_status` (ACTIVE, PENDING_CREATE, ERROR, etc.)
- `subnet_id` (computed if not provided)
- Default values for `admin_state_up`, `backup`, `weight`

### 3. Code Comparison: With vs Without Return Value

#### Option A: AddAndPoll Returns Nothing (Pattern 2)

```go
// In Terraform provider - NOT RECOMMENDED
func (r *Resource) Create(ctx context.Context, req, resp) {
    params := cloud.LoadBalancerPoolMemberAddParams{...}

    // Only get error, no member returned
    err := r.client.Members.AddAndPoll(ctx, poolID, params)
    if err != nil {
        resp.Diagnostics.AddError("failed to create member", err.Error())
        return
    }

    // Problem: We don't have the member ID or any member data!
    // We'd need to somehow get it...
    // Option 1: Parse from task (but AddAndPoll doesn't return task)
    // Option 2: List all members and find ours (how? by address?)
    // Option 3: Store task ID separately (defeats purpose of AndPoll)

    // This pattern doesn't work for members!
}
```

**This pattern fails** because we have no way to get the member ID or state.

#### Option B: AddAndPoll Returns *Member (Pattern 1) ✅

```go
// In Terraform provider - RECOMMENDED
func (r *Resource) Create(ctx context.Context, req, resp) {
    params := cloud.LoadBalancerPoolMemberAddParams{
        ProjectID:    param.NewOpt(data.ProjectID.ValueInt64()),
        RegionID:     param.NewOpt(data.RegionID.ValueInt64()),
        Address:      data.Address.ValueString(),
        ProtocolPort: data.ProtocolPort.ValueInt64(),
        Weight:       param.NewOpt(data.Weight.ValueInt64()),
        SubnetID:     param.NewOpt(data.SubnetID.ValueString()),
    }

    // Get complete member object with all computed fields!
    member, err := r.client.Cloud.LoadBalancers.Pools.Members.AddAndPoll(
        ctx,
        data.PoolID.ValueString(),
        params,
        option.WithMiddleware(logging.Middleware(ctx)),
    )
    if err != nil {
        resp.Diagnostics.AddError("failed to create member", err.Error())
        return
    }

    // Directly populate Terraform state - no extra API calls needed!
    data.ID = types.StringValue(member.ID)
    data.OperatingStatus = types.StringValue(string(member.OperatingStatus))
    data.ProvisioningStatus = types.StringValue(string(member.ProvisioningStatus))
    data.Address = types.StringValue(member.Address)
    data.ProtocolPort = types.Int64Value(member.ProtocolPort)
    data.AdminStateUp = types.BoolValue(member.AdminStateUp)
    data.Backup = types.BoolValue(member.Backup)
    data.Weight = types.Int64Value(member.Weight)
    if member.SubnetID != "" {
        data.SubnetID = types.StringValue(member.SubnetID)
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

**Code reduction**: ~50 lines → ~20 lines (60% reduction)

### 4. Real-World Use Case: Health Monitoring

Returning the member is critical for immediate health checks:

```go
member, err := client.Cloud.LoadBalancers.Pools.Members.AddAndPoll(ctx, poolID, params)
if err != nil {
    return err
}

// Immediately check member health - critical for production deployments!
switch member.OperatingStatus {
case cloud.LoadBalancerOperatingStatusOnline:
    log.Info("Member is healthy and receiving traffic")

case cloud.LoadBalancerOperatingStatusError:
    log.Error("Member failed health check")
    return fmt.Errorf("member %s is in ERROR state", member.ID)

case cloud.LoadBalancerOperatingStatusOffline:
    log.Warn("Member created but offline, may need troubleshooting")

case cloud.LoadBalancerOperatingStatusNoMonitor:
    log.Info("Member created, no health monitor configured")
}
```

Without returning the member, you'd need to:
1. Call `Pools.Get()` separately
2. Search through `pool.Members[]`
3. Check the status
4. Handle race conditions (status might change between calls)

### 5. Comparison with Existing SDK Patterns

The SDK already has precedent for this pattern:

#### Similar Resources (Return the Object)

```go
// LoadBalancerPoolService - Returns the created pool
func (r *LoadBalancerPoolService) NewAndPoll(ctx, params) (*LoadBalancerPool, error)
// Why? Pools have computed fields: id, operating_status, provisioning_status

// InstanceInterfaceService - Returns list of interfaces
func (r *InstanceInterfaceService) AttachAndPoll(ctx, instanceID, params) (*NetworkInterfaceList, error)
// Why? Need to see all interfaces and their relationships

// Our case - Should return the created member
func (r *LoadBalancerPoolMemberService) AddAndPoll(ctx, poolID, params) (*Member, error)
// Why? Members have computed fields: id, operating_status, provisioning_status, subnet_id
```

#### Different Resources (Return Nothing)

```go
// VolumeService - Returns only error
func (r *VolumeService) AttachToInstanceAndPoll(ctx, volumeID, params) error
// Why? Volume properties don't change when attached; it's just a link
```

**Key difference**: Volume attachment doesn't create a new resource or change state. Member creation **does both**.

### 6. The API Reality

The GCore API for member creation works like this:

```
POST /lbpools/{pool_id}/member
Response: { "tasks": ["task-uuid"] }

GET /tasks/{task-uuid} (after completion)
Response: {
    "state": "FINISHED",
    "created_resources": { "members": ["member-uuid"] }
}
```

**Important**: The API only returns the member ID, not the full member object.

To get member details, you **must** call:
```
GET /lbpools/{pool_id}
Response: {
    "id": "pool-uuid",
    "members": [
        {
            "id": "member-uuid",
            "operating_status": "ONLINE",
            "provisioning_status": "ACTIVE",
            "address": "192.168.1.10",
            // ... all other fields
        }
    ]
}
```

**This extra call is inevitable.** The question is:
- **Pattern 1**: SDK makes the call, returns clean `*Member` to consumer
- **Pattern 2**: Consumer makes the call themselves (more work, more error-prone)

**Pattern 1 wins** because it provides better UX and hides complexity.

## Recommendation

### Implement Pattern 1: Return `*Member`

```go
func (r *LoadBalancerPoolMemberService) AddAndPoll(
    ctx context.Context,
    poolID string,
    params LoadBalancerPoolMemberAddParams,
    opts ...option.RequestOption,
) (*Member, error) {
    // 1. Add member (returns task)
    // 2. Poll task until complete
    // 3. Extract member ID from task.CreatedResources.Members[0]
    // 4. Get pool and find member in pool.Members[]
    // 5. Return the complete member object
}
```

### Benefits for Terraform

1. **Cleaner code**: 60% reduction in Create() method
2. **No extra API calls in Terraform**: SDK handles it
3. **Immediate access to critical fields**: `operating_status`, `provisioning_status`
4. **Fewer bugs**: No manual pool fetching + member searching
5. **Consistent with pool pattern**: `NewAndPoll()` returns the created resource

### Benefits for All SDK Consumers

1. **Better DX**: Get everything in one call
2. **Type safety**: Return typed `*Member` instead of forcing consumers to search
3. **Performance**: SDK can optimize (e.g., cache if API adds member to task response)
4. **Consistency**: Matches pattern of other resource creation methods

## Conclusion

**Yes, returning `*Member` from `AddAndPoll` is not just helpful—it's essential** for a good Terraform integration and SDK user experience.

The pattern to follow is:
```
Load balancer member = Stateful resource (like Pool)
                     ≠ Pure association (like Volume attachment)

Therefore: AddAndPoll should return *Member ✅
```

This decision is documented in detail in:
- `/Users/user/repos/sdk-gcore-go/LBMEMBERS_AND_POLL.md` - Full implementation guide
- `/Users/user/repos/sdk-gcore-go/LBMEMBER_PATTERN_DECISION.md` - Pattern analysis and rationale

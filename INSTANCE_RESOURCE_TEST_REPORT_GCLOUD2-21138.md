# Instance Resource Test Report - GCLOUD2-21138

**Date:** 2025-11-19
**Branch:** terraform-instances
**Resource:** `gcore_cloud_instance`

---

## Executive Summary

Comprehensive testing of the instance resource revealed several bugs (all critical fixes applied) and significant functionality gaps in interface and security group management. The vm_state control feature was successfully implemented. Import functionality now properly populates interfaces and volumes.

**Pass Rate:** 7/7 tests (100%)

---

## Test Results

| Test | Description | Result | Notes |
|------|-------------|--------|-------|
| 1 | Create instance with image volume + external interface | PASS | Instance created successfully |
| 2 | Drift detection (plan after create) | PASS | No unexpected changes |
| 3 | Flavor change (g1-standard-1-2 to g1-standard-2-4) | PASS | In-place update via `/changeflavor` |
| 4 | Volume resize (5GB to 10GB) | PASS | After fixing VolumeID mapping bug |
| 5 | Stop instance (vm_state=stopped) | PASS | ActionAndPoll works correctly |
| 6 | Start instance (vm_state=active) | PASS | ActionAndPoll works correctly |
| 7 | Import functionality | PASS | Interfaces and volumes populated correctly |

---

## Bugs Found and Fixed

### 1. VolumeID JSON Mapping (CRITICAL - FIXED)

**File:** `internal/services/cloud_instance/model.go:88`

**Issue:** Commit e65ed34 incorrectly changed JSON mapping from `json:"id,computed_optional"` to `json:"volume_id,computed_optional"`

**Impact:**
- volume_id was null in state after creation
- Volume resize operations failed with null volume_id
- PATCH validation errors

**Root Cause:** API returns volume ID in the `"id"` field, not `"volume_id"`:
```json
"volumes": [
    {
        "id": "5d1f6fc8-11ac-48fa-95d0-8446dc5aad11",
        "delete_on_termination": true
    }
]
```

**Fix Applied:**
```go
// Before (incorrect)
VolumeID types.String `tfsdk:"volume_id" json:"volume_id,computed_optional"`

// After (correct)
VolumeID types.String `tfsdk:"volume_id" json:"id,computed_optional"`
```

---

## Issues Fixed

### 1. Import Doesn't Populate Volumes (FIXED)

**Symptom:** After `terraform import`, volumes array was null

**Root Cause:** API returns minimal volume info (`id`, `delete_on_termination`) but the decoder couldn't populate nested fields properly.

**Fix Applied:** Manually extract `volume_id` and `delete_on_termination` from API response in Create, Read, and ImportState functions:

```go
// Extract volume_id and delete_on_termination from volumes
if volumesArr, ok := rawResponse["volumes"].([]interface{}); ok && data.Volumes != nil {
    for i, vol := range volumesArr {
        if i < len(*data.Volumes) && vol != nil {
            volMap, ok := vol.(map[string]interface{})
            if !ok {
                continue
            }
            if id, ok := volMap["id"].(string); ok {
                (*data.Volumes)[i].VolumeID = types.StringValue(id)
            }
            if dot, ok := volMap["delete_on_termination"].(bool); ok {
                (*data.Volumes)[i].DeleteOnTermination = types.BoolValue(dot)
            }
        }
    }
}
```

### 2. Addresses Showing Null Values (FIXED)

**Symptom:** addresses map values were all null

**Root Cause:** CloudInstanceAddressesModel fields were marked as `required`/`optional` instead of `computed`.

**Fix Applied:** Changed all address fields to use `computed` tag:

```go
type CloudInstanceAddressesModel struct {
    Addr          types.String `tfsdk:"addr" json:"addr,computed"`
    Type          types.String `tfsdk:"type" json:"type,computed"`
    InterfaceName types.String `tfsdk:"interface_name" json:"interface_name,computed"`
    SubnetID      types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
    SubnetName    types.String `tfsdk:"subnet_name" json:"subnet_name,computed"`
}
```

### 3. Drift in Computed Fields (FIXED)

**Symptom:** Multiple computed fields show as changes on each plan:
- `created_at` -> (known after apply)
- `creator_task_id` -> (known after apply)
- `fixed_ip_assignments` -> (known after apply)

**Root Cause:** Computed fields missing `UseStateForUnknown()` plan modifiers.

**Fix Applied:** Added `UseStateForUnknown()` plan modifiers to static computed fields in `internal/services/cloud_instance/schema.go`:

```go
// Static computed fields - keep their values across operations
"id": schema.StringAttribute{
    Computed:      true,
    PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
},
"created_at": schema.StringAttribute{
    Computed:      true,
    CustomType:    timetypes.RFC3339Type{},
    PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
},
"creator_task_id": schema.StringAttribute{
    Computed:      true,
    PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
},
"region": schema.StringAttribute{
    Computed:      true,
    PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
},
"instance_description": schema.StringAttribute{
    Computed:      true,
    PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
},
"addresses": schema.MapAttribute{
    Computed:      true,
    PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
},
"blackhole_ports": schema.ListNestedAttribute{
    Computed:      true,
    PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
},
"ddos_profile": schema.SingleNestedAttribute{
    Computed:      true,
    PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
},
"fixed_ip_assignments": schema.ListNestedAttribute{
    Computed:      true,
    PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
},
"instance_isolation": schema.SingleNestedAttribute{
    Computed:      true,
    PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
},
```

**Important:** Dynamic computed fields that change during operations should NOT have `UseStateForUnknown()`:
- `status` - changes with instance state (ACTIVE, SHUTOFF, etc.)
- `task_id` - changes when tasks run
- `task_state` - changes when tasks run
- `tasks` - changes when tasks are created

---

## Remaining Issues

None - all Phase 1 issues have been resolved

---

## Critical Missing Features

### Interface Attach/Detach

**API Endpoints:**
- `POST /cloud/v1/instances/{project_id}/{region_id}/{instance_id}/attach_interface`
- `POST /cloud/v1/instances/{project_id}/{region_id}/{instance_id}/detach_interface`

**Current Status:**
- APIs exist in Gcore Cloud
- **SDK methods EXIST** in `sdk-gcore-go/cloud/instanceinterface.go`
- Resource does NOT use these methods

**SDK Methods Available:**
```go
r.client.Cloud.Instances.Interfaces.Attach(ctx, instanceID, params, opts...)
r.client.Cloud.Instances.Interfaces.AttachAndPoll(ctx, instanceID, params, opts...)
r.client.Cloud.Instances.Interfaces.Detach(ctx, instanceID, params, opts...)
r.client.Cloud.Instances.Interfaces.DetachAndPoll(ctx, instanceID, params, opts...)
```

**Attach Params (union types for 4 interface variants):**
- `OfNewInterfaceExternalExtendSchemaWithDDOS` - type: external
- `OfNewInterfaceSpecificSubnetSchema` - type: subnet
- `OfNewInterfaceAnySubnetSchema` - type: any_subnet
- `OfNewInterfaceReservedFixedIPSchema` - type: reserved_fixed_ip

**Detach Params:**
- `IPAddress` (required)
- `PortID` (required)

**Impact:** Any change to `interfaces` attribute causes state drift (changes are silently ignored)

### Security Group Assign/Unassign

**API Endpoints:**
- `POST /cloud/v1/instances/{project_id}/{region_id}/{instance_id}/addsecuritygroup`
- `POST /cloud/v1/instances/{project_id}/{region_id}/{instance_id}/delsecuritygroup`

**Current Status:**
- APIs exist in Gcore Cloud
- SDK methods exist (`AssignSecurityGroup`, `UnassignSecurityGroup`)
- **API uses security group NAMES, but Terraform schema uses IDs**

**Impact:** Any change to `security_groups` attribute causes state drift

---

## Required Fixes

### 1. Interface Attach/Detach Implementation

**SDK methods already exist** - no SDK updates required.

**Resource Implementation Needed:**

```go
// In Update function, add interface change handling

// Get current interfaces from state (need port_id and ip_address for detach)
// Compare with desired interfaces from plan
// For removed interfaces: call Interfaces.DetachAndPoll
// For added interfaces: call Interfaces.AttachAndPoll

// Example for detaching an interface
detachParams := cloud.InstanceInterfaceDetachParams{
    IPAddress: ipAddress,
    PortID:    portID,
}
if !data.ProjectID.IsNull() {
    detachParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
}
if !data.RegionID.IsNull() {
    detachParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
}
_, err := r.client.Cloud.Instances.Interfaces.DetachAndPoll(ctx, instanceID, detachParams, opts...)

// Example for attaching an external interface
attachParams := cloud.InstanceInterfaceAttachParams{
    OfNewInterfaceExternalExtendSchemaWithDDOS: &cloud.InstanceInterfaceAttachParamsBodyNewInterfaceExternalExtendSchemaWithDDOS{
        Type:     param.NewOpt("external"),
        IPFamily: "ipv4",
    },
}
if !data.ProjectID.IsNull() {
    attachParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
}
if !data.RegionID.IsNull() {
    attachParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
}
_, err := r.client.Cloud.Instances.Interfaces.AttachAndPoll(ctx, instanceID, attachParams, opts...)
```

**Challenge:** Need to track port_id and ip_address in state for detach operations. Currently these may not be preserved.

### 2. Security Group Implementation

**Option A (Preferred): Modify API to accept IDs**

Request API team to accept security group IDs in addition to names for:
- `POST /instances/{id}/addsecuritygroup`
- `POST /instances/{id}/delsecuritygroup`

**Option B: Lookup names from IDs**

In the resource Update function:
1. Get list of security group IDs from plan
2. Call `GET /securitygroups` to fetch all security groups
3. Build map of ID -> Name
4. Call AssignSecurityGroup/UnassignSecurityGroup with names

```go
// Pseudocode for security group update
func updateSecurityGroups(ctx context.Context, client *gcore.Client, instanceID string,
    currentSGs, desiredSGs []string, projectID, regionID int64) error {

    // Fetch all security groups to build ID->Name map
    sgList, err := client.Cloud.SecurityGroups.List(ctx, cloud.SecurityGroupListParams{
        ProjectID: param.NewOpt(projectID),
        RegionID:  param.NewOpt(regionID),
    })

    idToName := make(map[string]string)
    for _, sg := range sgList.Results {
        idToName[sg.ID] = sg.Name
    }

    // Find security groups to add
    toAdd := difference(desiredSGs, currentSGs)
    for _, sgID := range toAdd {
        sgName := idToName[sgID]
        err := client.Cloud.Instances.AssignSecurityGroup(ctx, instanceID,
            cloud.InstanceAssignSecurityGroupParams{
                ProjectID: param.NewOpt(projectID),
                RegionID:  param.NewOpt(regionID),
                Name:      param.NewOpt(sgName),
            })
    }

    // Find security groups to remove
    toRemove := difference(currentSGs, desiredSGs)
    for _, sgID := range toRemove {
        sgName := idToName[sgID]
        err := client.Cloud.Instances.UnassignSecurityGroup(ctx, instanceID,
            cloud.InstanceUnassignSecurityGroupParams{
                ProjectID: param.NewOpt(projectID),
                RegionID:  param.NewOpt(regionID),
                Name:      param.NewOpt(sgName),
            })
    }

    return nil
}
```

### 3. Resource Update Function Changes

Add to `internal/services/cloud_instance/resource.go` Update function:

```go
// After vm_state handling, before the final re-read

// Handle interface changes
interfacesChanged := false
if !reflect.DeepEqual(data.Interfaces, state.Interfaces) {
    // Compare current vs desired interfaces
    // Call AttachInterface for new interfaces
    // Call DetachInterface for removed interfaces
    interfacesChanged = true
}

// Handle security group changes
securityGroupsChanged := false
if !reflect.DeepEqual(data.SecurityGroups, state.SecurityGroups) {
    // Get current and desired security group IDs
    // Lookup names from IDs
    // Call AssignSecurityGroup for new groups
    // Call UnassignSecurityGroup for removed groups
    securityGroupsChanged = true
}

// Update the condition for skipping PATCH
if flavorChanged || volumeResized || vmStateChanged || interfacesChanged || securityGroupsChanged {
    // Re-read and return
}
```

---

## Implementation Priority

### Phase 1: Critical Fixes (COMPLETED)
1. ✅ VolumeID JSON mapping fix (commit e65ed34 corrected)
2. ✅ Fix Import functionality (volumes null) - manual extraction implemented
3. ✅ Fix addresses null values - JSON tags changed to computed
4. ✅ Fix drift in computed fields - added UseStateForUnknown plan modifiers

### Phase 2: Interface Update Support (COMPLETED)

**Implementation Complete:**

1. ✅ Removed `no_refresh` from interfaces field in model.go
2. ✅ Made `port_id` computed and `ip_address` computed_optional
3. ✅ Added interface extraction in Create/Read using `Interfaces.List` API
4. ✅ Implemented interface attach/detach logic in Update function
5. ✅ Mapped Terraform interface types to SDK union variants:
   - external → OfNewInterfaceExternalExtendSchemaWithDDOS
   - subnet → OfNewInterfaceSpecificSubnetSchema
   - any_subnet → OfNewInterfaceAnySubnetSchema
   - reserved_fixed_ip → OfNewInterfaceReservedFixedIPSchema
6. ✅ Calls `Interfaces.AttachAndPoll` for new interfaces
7. ✅ Calls `Interfaces.DetachAndPoll` for removed interfaces

**Test Results:**
- Create instance with external interface: ✅ PASS
- port_id and ip_address populated in state: ✅ PASS
- No drift on subsequent plans: ✅ PASS

**Key Code Changes:**

```go
// model.go - Interfaces field (removed no_refresh)
Interfaces *[]*CloudInstanceInterfacesModel `tfsdk:"interfaces" json:"interfaces,required"`

// model.go - Interface fields (port_id computed, ip_address computed_optional)
IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed_optional"`
PortID    types.String `tfsdk:"port_id" json:"port_id,computed"`

// schema.go - Removed RequiresReplace from interfaces
// (interfaces can now be updated in-place)

// resource.go - Interface extraction after Create/Read
interfaces, err := r.client.Cloud.Instances.Interfaces.List(ctx, instanceID, listParams, ...)
for i, iface := range interfaces.Results {
    (*data.Interfaces)[i].PortID = types.StringValue(iface.PortID)
    if len(iface.IPAssignments) > 0 {
        (*data.Interfaces)[i].IPAddress = types.StringValue(iface.IPAssignments[0].IPAddress)
    }
}

// resource.go - Detach interfaces (Update function)
detachParams := cloud.InstanceInterfaceDetachParams{
    IPAddress: stateIface.IPAddress.ValueString(),
    PortID:    stateIface.PortID.ValueString(),
}
r.client.Cloud.Instances.Interfaces.DetachAndPoll(ctx, instanceID, detachParams, ...)

// resource.go - Attach interfaces (Update function)
attachParams := cloud.InstanceInterfaceAttachParams{
    OfNewInterfaceExternalExtendSchemaWithDDOS: &cloud.InstanceInterfaceAttachParamsBodyNewInterfaceExternalExtendSchemaWithDDOS{
        Type: param.NewOpt("external"),
    },
}
r.client.Cloud.Instances.Interfaces.AttachAndPoll(ctx, instanceID, attachParams, ...)
```

### Phase 3: Security Group Update Support (Future Work)
1. Implement ID to name lookup (fetch security groups list)
2. Implement security group comparison logic
3. Call AssignSecurityGroup for new groups
4. Call UnassignSecurityGroup for removed groups
5. Test security group add/remove scenarios

**Note:** SDK methods already exist for all these operations.

---

## vm_state Implementation (COMPLETED)

Successfully implemented vm_state control feature for legacy provider parity.

**Schema Changes:** `internal/services/cloud_instance/schema.go`
```go
"vm_state": schema.StringAttribute{
    Description: "Virtual machine state. Set to 'active' to start or 'stopped' to stop.",
    Computed:    true,
    Optional:    true,
    Validators: []validator.String{
        stringvalidator.OneOfCaseInsensitive("active", "stopped"),
    },
},
```

**Model Changes:** `internal/services/cloud_instance/model.go`
```go
VmState types.String `tfsdk:"vm_state" json:"vm_state,computed_optional"`
```

**Resource Changes:** `internal/services/cloud_instance/resource.go`
- Added import for `github.com/G-Core/gcore-go/shared/constant`
- Added vm_state change detection and handling
- Uses `ActionAndPoll` with proper action serializers:
  - Start: `OfStartActionInstanceSerializer` with `constant.ValueOf[constant.Start]()`
  - Stop: `OfBasicActionInstanceSerializer` with `action: "stop"`

---

## MITM Evidence

API calls captured in `test-instance-skill/flow.mitm`:
- POST /instances (create)
- GET /instances/{id} (read - multiple calls)
- POST /instances/{id}/changeflavor (flavor change)
- POST /volumes/{id}/extend (volume resize)
- POST /instances/{id}/action (start/stop)
- DELETE /instances/{id} (cleanup)

---

## Test Environment

- **Directory:** `/Users/user/repos/gcore-terraform/test-instance-skill`
- **Project ID:** 379987
- **Region ID:** 76
- **Image:** ubuntu-22.04-x64 (6343932d-0257-4285-bf89-05060f24095a)
- **Proxy:** mitmproxy on port 9092

---

## Files Modified

1. `internal/services/cloud_instance/schema.go` - vm_state Optional+Computed with validators, UseStateForUnknown plan modifiers for computed fields
2. `internal/services/cloud_instance/model.go` - VmState computed_optional, VolumeID fix, addresses JSON tags changed to computed
3. `internal/services/cloud_instance/resource.go` - ActionAndPoll logic for vm_state, manual extraction of volume_id and delete_on_termination in Create/Read/ImportState

---

## Additional Bug Fixes (2025-11-19)

### Bug: vm_state Validation Error During Resize - FIXED ✅

**Location:** `internal/services/cloud_instance/resource.go:187-190, 291-295`

**Symptom:** After successful resize operation, Terraform fails with:
```
Error: invalid vm_state
vm_state must be 'active' or 'stopped', got: resized
```

**Root Cause:** The Update function validates vm_state against user input values, but during state refresh after resize, the API returns transient states like "resized". The `UnmarshalComputed` overwrites `data.VmState` with this transient value, causing validation to fail.

**Fix Applied:**
```go
// Preserve the user's desired vm_state before any operations that might overwrite data
// This is needed because UnmarshalComputed will overwrite data.VmState with API response
// which may contain transient states like "resized" during operations
plannedVmState := data.VmState

// ... later in vm_state handling ...
// Use plannedVmState (the user's desired value) instead of data.VmState
if !plannedVmState.IsNull() && !plannedVmState.IsUnknown() && !plannedVmState.Equal(state.VmState) {
    desiredState := plannedVmState.ValueString()
    // ... rest of vm_state logic
}
```

**Status:** VERIFIED - Resize now completes successfully without validation errors.

---

### Enhancement: Read Function State Preservation - IMPLEMENTED ✅

**Location:** `internal/services/cloud_instance/resource.go:747-772`

**Issue:** The `UnmarshalComputed` call in Read overwrites `data.Interfaces` and `data.Volumes` with null values since the instances API response doesn't include full interface/volume details.

**Fix Applied:**
```go
// Preserve interfaces and volumes from prior state before UnmarshalComputed overwrites them
// The instances API doesn't return interface details or full volume info
priorInterfaces := data.Interfaces
priorVolumes := data.Volumes

// ... API call and UnmarshalComputed ...

// Restore interfaces and volumes from prior state after UnmarshalComputed
data.Interfaces = priorInterfaces
data.Volumes = priorVolumes
```

**Status:** IMPLEMENTED - Prevents Read from clearing user-defined interface and volume configuration.

---

### ImportState Interface Fetch - FIXED ✅

**Location:** `internal/services/cloud_instance/resource.go:975-1019`

**Implementation:**
- Code added to call `Interfaces.List` API during import
- Interface type detection using `iface.NetworkDetails.External` flag
- Populates port_id, ip_address, type for imported interfaces

**Resolution:** The implementation was correct. The issue was that Terraform was using an older binary (`terraform-provider-gcore`) instead of the newly built one (`terraform-provider-gcore_v99.0.0`). Fixed by building to the correct binary name.

**Verified Results:**
- Import now populates interfaces with: ip_address, port_id, type
- Import populates volumes with: volume_id, delete_on_termination, source
- Import populates: name, flavor, addresses, and other computed fields
- Plan after import shows update-in-place (not replacement)

---

## Next Steps

1. ~~Review and merge VolumeID fix~~ ✅ Fixed
2. ~~Fix vm_state validation during resize~~ ✅ Fixed
3. ~~Debug ImportState interfaces API call issue~~ ✅ Fixed (was building to wrong binary name)
4. Create Jira tickets for SDK updates (AttachInterface/DetachInterface)
5. Discuss with API team about accepting security group IDs
6. Implement remaining fixes per priority above

# Pedro's PR Review Comments - Response Document

**PR:** [#35 - Support gcore_instance in Stainless terraform provider](https://github.com/stainless-sdks/gcore-terraform/pull/35)
**Jira Ticket:** [GCLOUD2-21138](https://jira.gcore.lu/browse/GCLOUD2-21138)
**Date:** 2026-01-09

---

## Comment 1: cloud_instance_image/data_source.go changes

**GitHub Link:** [discussion_r2668903017](https://github.com/stainless-sdks/gcore-terraform/pull/35#discussion_r2668903017)

**File:** `internal/services/cloud_instance_image/data_source.go` (Line 1)

**Pedro's Comment:**
> Why are these changes needed? Let's keep the original code, as we usually only need this to use polling methods.

**Code Context:**
```go
// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance_image
```

**Response:**
Agree - these changes appear unrelated to the instance resource functionality. The cloud_instance_image data source doesn't require modifications for the instance resource work. Will revert these changes.

**Action:** Revert changes to this file.

---

## Comment 2: Custom functions at top of file

**GitHub Link:** [discussion_r2669001826](https://github.com/stainless-sdks/gcore-terraform/pull/35#discussion_r2669001826)

**File:** `internal/services/cloud_instance/resource.go` (Line 24)

**Pedro's Comment:**
> I suggest moving any new custom functions to the end of the file so that it resembles the other resource files in the codebase. Alternatively, move them to a new file.

**Code Context:**
```go
// convertInterfacesToParams converts TF interface models to SDK InstanceNewParamsInterfaceUnion slice
func convertInterfacesToParams(interfaces *[]*CloudInstanceInterfacesModel) []cloud.InstanceNewParamsInterfaceUnion {
	if interfaces == nil {
```

**Response:**
Agree - will move helper functions (`convertInterfacesToParams`, `convertVolumesToParams`, `convertSecurityGroupsToParams`, `tagsEqual`) to the end of the file to match codebase conventions.

**Action:** Move helper functions to end of resource.go file.

---

## Comment 3: Duplicate tagsEqual function

**GitHub Link:** [discussion_r2668994477](https://github.com/stainless-sdks/gcore-terraform/pull/35#discussion_r2668994477)

**File:** `internal/services/cloud_instance/resource.go` (Line 216)

**Pedro's Comment:**
> You reuse the existing function inside `internal/custom`.

**Code Context:**
```go
// tagsEqual compares two tag maps for equality
func tagsEqual(a, b *map[string]types.String) bool {
	if a == nil && b == nil {
```

**Response:**
Correct - there's already a `TagsEqual` function in `internal/custom/custom.go`:

```go
// TagsEqual compares two maps of tags for equality.
func TagsEqual(a, b *map[string]types.String) bool {
```

Will remove the duplicate local function and use `custom.TagsEqual` instead.

**Action:** Import `internal/custom` and replace `tagsEqual` calls with `custom.TagsEqual`.

---

## Comment 4: Build params directly vs WithRequestBody

**GitHub Link:** [discussion_r2669010465](https://github.com/stainless-sdks/gcore-terraform/pull/35#discussion_r2669010465)

**File:** `internal/services/cloud_instance/resource.go` (Line 286)

**Pedro's Comment:**
> Why is this needed for instances? I understand that we have this problem to fix, but typically, it doesn't affect us.

**Code Context:**
```go
	// Build params directly instead of using option.WithRequestBody
	// This ensures the body is only sent with the POST request, not with polling GET requests
	params := cloud.InstanceNewParams{
```

**Response:**
This was implemented to fix a bug reported by QA (Kirill) on **Dec 2, 2025**:

> *"endpoint GET /v1/instances/{project_id}/{region_id}/{instance_id} sends with payload (with req. body). See a screenshot."*

The fix was confirmed on **Dec 5, 2025**:

> *"Fixed the problem where GET requests during polling were sending request body. Refactored Create/Update methods to pass body as SDK params instead of using option.WithRequestBody. Verified in debug logs: GET requests now have NO body, POST/PATCH have correct body."*

**Decision:** The fix was necessary because during `NewAndPoll` operations, using `option.WithRequestBody` caused the request body to leak into subsequent polling GET requests. However, if the SDK has since been fixed to handle this correctly, we could potentially revert to using `option.WithRequestBody`.

**Action:** Test if SDK now handles this properly. If yes, can revert to using `option.WithRequestBody` as Pedro suggests.

---

## Comment 5: Flavor field API discrepancy

**GitHub Link:** [discussion_r2669081384](https://github.com/stainless-sdks/gcore-terraform/pull/35#discussion_r2669081384)

**File:** `internal/services/cloud_instance/resource.go` (Line 1372)

**Pedro's Comment:**
> Has this been raised to the dev team? There shouldn't be the same name for two different objects (an `id` in the create body and an object in the get response model). I guess this is hard to rename now...

**Code Context:**
```go
	// Extract fields from the API response that aren't mapped via json tags
	// API returns: {"flavor": {"flavor_id": "..."}, "project_id": N, "region_id": N, ...}
	var rawResponse map[string]interface{}
```

**Response:**
This is an API design inconsistency:
- **Create request:** `flavor` is a string (flavor ID)
- **GET response:** `flavor` is a nested object `{"flavor": {"flavor_id": "...", "vcpus": N, ...}}`

The manual extraction workaround was implemented to support flavor changes via the `/changeflavor` endpoint (reported by Kirill on **Nov 13, 2025**).

**Decision:** This should be raised to the API team as a design issue. The workaround is necessary for correct Terraform behavior until the API is fixed.

**Action:** Create a ticket to raise this API inconsistency with the backend team.

---

## Comment 6: Extract project_id and region_id

**GitHub Link:** [discussion_r2669084711](https://github.com/stainless-sdks/gcore-terraform/pull/35#discussion_r2669084711)

**File:** `internal/services/cloud_instance/resource.go` (Line 1380)

**Pedro's Comment:**
> Why do we need to extract project_id and region_id? They should already be inside `data`.

**Code Context:**
```go
		// Extract project_id and region_id (not mapped via json tags since they use path: tags)
		if projectID, ok := rawResponse["project_id"].(float64); ok {
			data.ProjectID = types.Int64Value(int64(projectID))
```

**Response:**
The `ProjectID` and `RegionID` fields in the model have `path:` json tags (for URL path parameters), not standard json tags. After `UnmarshalComputed`, these fields might not be populated from the response body.

However, if they're already set from the request/prior state (which they should be), this extraction is likely redundant.

**Action:** Verify if `data.ProjectID`/`data.RegionID` are preserved from the request. If yes, remove this extraction code.

---

## Comment 7: Missing Unknown checks

**GitHub Link:** [discussion_r2669115085](https://github.com/stainless-sdks/gcore-terraform/pull/35#discussion_r2669115085) (inferred from context)

**File:** `internal/services/cloud_instance/resource.go` (Line 1393)

**Pedro's Comment:**
> We never check if data.ProjectID or data.RegionID are Unknown.

**Code Context:**
```go
		if !data.ProjectID.IsNull() && !data.ProjectID.IsUnknown() {
			listParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() && !data.RegionID.IsUnknown() {
```

**Response:**
Valid point. The code does check for Unknown in some places (line 1393-1397 shown above), but we should ensure consistency throughout the file.

**Action:** Audit all uses of `ProjectID`/`RegionID` and add `IsUnknown()` checks where missing.

---

## Comment 8: Order of addresses/volumes

**GitHub Link:** (inferred from Line 1407 context)

**File:** `internal/services/cloud_instance/resource.go` (Line 1407)

**Pedro's Comment:**
> How can we be sure that the order will be the same?

**Response:**
Valid concern. The API may not guarantee ordering of volumes/interfaces in responses. Matching by array index is fragile.

**Action:** Refactor to match volumes by `volume_id` instead of assuming array order is preserved.

---

## Comment 9: Comment about removed RequiresReplace

**GitHub Link:** (inferred from schema.go Line 48)

**File:** `internal/services/cloud_instance/schema.go` (Line 48)

**Pedro's Comment:**
> The comment removing the previous code can be removed here and in other locations.

**Code Context:**
```go
		"flavor": schema.StringAttribute{
			Description: "The flavor of the instance.",
			Required:    true,
			// RequiresReplace removed - flavor changes now use the /changeflavor endpoint
		},
```

**Response:**
The comment documents WHY `RequiresReplace` was removed (to enable in-place flavor changes via `/changeflavor` endpoint). This was a deliberate design decision from **Nov 13-14, 2025**.

**Decision:** Comments can be shortened but should remain to explain the design choice. Alternatively, can remove if we add this to a CHANGELOG or design doc.

**Action:** Shorten or remove comments as preferred.

---

## Comment 10: port_id computed vs required

**GitHub Link:** (inferred from schema.go Line 126)

**File:** `internal/services/cloud_instance/schema.go` (Line 126)

**Pedro's Comment:**
> The `port_id` is only `computed` and not `computed_optional`? Shouldn't it be required if the type is `"reserved_fixed_ip"`?

**Code Context:**
```go
					"port_id": schema.StringAttribute{
						Description: "Port ID assigned to this interface. Computed after creation.",
						Computed:    true,
					},
```

**Response:**
Pedro is correct:
- For `reserved_fixed_ip` type: user **must** provide `port_id`
- For other types (`external`, `subnet`, `any_subnet`): `port_id` is **computed** by the API

Should be `Optional+Computed` with schema-level validation requiring it when type is `reserved_fixed_ip`.

**Action:** Change `port_id` to `Optional: true, Computed: true` and add a custom validator that requires it when type is `reserved_fixed_ip`.

---

## Comment 11: Plan modifiers for read-only fields

**GitHub Link:** (inferred from schema.go Line 218)

**File:** `internal/services/cloud_instance/schema.go` (Line 218)

**Pedro's Comment:**
> Are you sure we need the plan modifier for these (this and below) read-only fields?

**Code Context:**
```go
		"created_at": schema.StringAttribute{
			Description:   "Datetime when instance was created",
			Computed:      true,
			CustomType:    timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
```

**Response:**
**YES**, these are needed. This was a deliberate fix from **Nov 19, 2025** to address drift in computed fields:

> *"Added UseStateForUnknown() plan modifiers to static computed fields: id, created_at, creator_task_id, region, instance_description, addresses, blackhole_ports, ddos_profile, fixed_ip_assignments, instance_isolation. Dynamic fields (status, task_id, task_state, tasks) left without modifier as they change during operations."*

Without `UseStateForUnknown()`, Terraform shows `(known after apply)` on every plan for fields that don't change after creation.

**Decision:** Keep plan modifiers for static computed fields.

---

## Comment 12: addresses UseStateForUnknown comment

**GitHub Link:** (inferred from schema.go Line 286)

**File:** `internal/services/cloud_instance/schema.go` (Line 286)

**Pedro's Comment:**
> Remove the comment, I guess the plan modifier was added and then removed?

**Code Context:**
```go
		"addresses": schema.MapAttribute{
			Description:   "Map of `network_name` to list of addresses in that network",
			Computed:      true,
			CustomType:    customfield.NewMapType[customfield.NestedObjectList[CloudInstanceAddressesModel]](ctx),
			// Note: UseStateForUnknown removed - addresses change when interfaces are attached/detached
```

**Response:**
The comment documents an intentional design decision. `UseStateForUnknown` was initially added, then **removed** on **Dec 22, 2025** because:

> *"Bug 1: 'addresses element appeared/vanished' - removed UseStateForUnknown from addresses field"*

When interfaces are attached/detached, addresses change. Keeping `UseStateForUnknown` caused Terraform errors about "element appeared/vanished".

**Decision:** Keep the comment - it explains WHY this field is different from other computed fields.

---

## Comment 13: tasks attribute removal

**GitHub Link:** (inferred from schema.go Line 312)

**File:** `internal/services/cloud_instance/schema.go` (Line 312)

**Pedro's Comment:**
> We have been removing the `tasks` attribute in every resource schema.

**Code Context:**
```go
		"tasks": schema.ListAttribute{
			Description: "List of task IDs for async operations. Poll via GET /v1/tasks/{task_id} until FINISHED or ERROR.",
			Computed:    true,
			CustomType:  customfield.NewListType[types.String](ctx),
			ElementType: types.StringType,
		},
```

**Response:**
Agree - will follow codebase convention and remove the `tasks` attribute.

**Action:** Remove `tasks` attribute from schema.

---

## Comment 14: Plan modifier for security_groups

**GitHub Link:** (inferred from schema.go Line 321)

**File:** `internal/services/cloud_instance/schema.go` (Line 321)

**Pedro's Comment:**
> Same question as before, do we need the plan modifier here? Maybe we need because they can be modified by the server and we want to reflect that? Probably the same question for `ddos_profile`.

**Response:**
`security_groups` is user-specified (`Optional`), not computed by the server. It shouldn't need `UseStateForUnknown` since it only changes when the user changes it.

`ddos_profile` is purely `Computed` (server-managed), so it **does** need `UseStateForUnknown` to prevent drift.

**Action:** Review security_groups - if it has `UseStateForUnknown`, remove it.

---

## Comment 15: Schema-level validation

**GitHub Link:** (inferred from resource.go Line 678)

**File:** `internal/services/cloud_instance/resource.go` (Line 678)

**Pedro's Comment:**
> Some of these checks should be done before using custom validators at the schema level, such that they fail early during `terraform plan`.

**Code Context:**
```go
			case "subnet":
				if planIface.SubnetID.IsNull() {
					resp.Diagnostics.AddError(
						"missing subnet_id",
						fmt.Sprintf("Interface at index %d with type 'subnet' requires subnet_id", i),
```

**Response:**
Valid point. Moving validation to schema-level validators provides better UX:
- Fails at `terraform plan` time (early feedback)
- Better error messages
- Follows Terraform best practices

**Action:** Add schema-level validators for interface type requirements:
- `subnet` type requires `subnet_id`
- `any_subnet` type requires `network_id`
- `reserved_fixed_ip` type requires `port_id`

---

## Comment 16: FloatingIP port_id in params

**GitHub Link:** (inferred from resource.go Line 786)

**File:** `internal/services/cloud_instance/resource.go` (Line 786)

**Pedro's Comment:**
> Why not provide port_id in the fipParams and remove the option.WithRequestBody?

**Code Context:**
```go
			// Create FIP body with port_id to assign it during creation
			fipBody := map[string]interface{}{
				"port_id": portID,
			}
```

**Response:**
If the SDK's `FloatingIPNewParams` supports `PortID` directly, we should use it instead of raw request body construction.

**Action:** Check if SDK supports `PortID` in `FloatingIPNewParams` and refactor accordingly.

---

## Comment 17: ListAutoPaging

**GitHub Link:** (inferred from resource.go Line 850)

**File:** `internal/services/cloud_instance/resource.go` (Line 850)

**Pedro's Comment:**
> We can use ListAutoPaging.

**Code Context:**
```go
			fipPage, err := r.client.Cloud.FloatingIPs.List(ctx, listParams, option.WithMiddleware(logging.Middleware(ctx)))
			if err != nil {
				resp.Diagnostics.AddError(
```

**Response:**
Agree - the SDK provides `ListAutoPaging` for cleaner pagination handling. Will refactor to use it.

**Action:** Refactor to use `ListAutoPaging` instead of manual pagination.

---

## Comment 18: Placement groups - remove from instances?

**GitHub Link:** (inferred from resource.go Line 1112)

**File:** `internal/services/cloud_instance/resource.go` (Line 1112)

**Pedro's Comment:**
> We already have a resource for managing placement groups, should we remove it from instances?

**Code Context:**
```go
// Handle servergroup/placement group changes using AddToPlacementGroup/RemoveFromPlacementGroup endpoints
servergroupChanged := false
stateServergroupID := ""
```

**Response:**
The old Terraform provider supported `servergroup_id` on instances. We implemented in-place updates on **Dec 22, 2025** to match old behavior:

> *"Servergroup changes are now in-place updates (no instance replacement)"*

**Decision:** Recommend keeping for backward compatibility. Removing would be a breaking change for users migrating from the old provider. However, this is a design decision that could go either way.

---

## Comment 19: Multiple change flags

**GitHub Link:** (inferred from resource.go Line 1180)

**File:** `internal/services/cloud_instance/resource.go` (Line 1180)

**Pedro's Comment:**
> If there are no actions that differ between these flags, I think we could have a single flag for signaling that a refresh is needed. In the resource implementations that I've done recently, I've used a `stateHasChanged` flag for this purpose.

**Code Context:**
```go
// If we've handled updates via specialized endpoints (flavor change, vm_state change, volume changes, interface changes, floating IP changes, security group changes, or servergroup changes),
// skip the standard PATCH Update call and just refresh state
if flavorChanged || vmStateChanged || volumesChanged || interfacesChanged || floatingIPChanged || securityGroupsChanged || servergroupChanged {
```

**Response:**
Valid simplification. The multiple flags can be consolidated into a single `stateHasChanged` flag.

**Action:** Refactor to use a single `stateHasChanged` flag instead of multiple individual flags.

---

## Comment 20: Move read logic to function

**GitHub Link:** (inferred from resource.go Line 1209)

**File:** `internal/services/cloud_instance/resource.go` (Line 1209)

**Pedro's Comment:**
> Maybe this extra logic on the read can be moved into a function.

**Code Context:**
```go
		// Extract flavor_id from the flavor object in the API response
		// API returns: {"flavor": {"flavor_id": "...", ...}}
		// We need to extract just the flavor_id string
		var rawResponse map[string]interface{}
```

**Response:**
Valid refactoring suggestion. The flavor/interface/volume extraction logic is duplicated in multiple places (Create, Read, Update, ImportState).

**Action:** Extract common read/refresh logic into a helper function like `refreshInstanceState()`.

---

## Comment 21: Return after partial update

**GitHub Link:** (inferred from resource.go Line 1248)

**File:** `internal/services/cloud_instance/resource.go` (Line 1248)

**Pedro's Comment:**
> Should we only update the data in the end after all updates, and we shouldn't return here if we still want to perform the regular updates after this.

**Code Context:**
```go
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}
```

**Response:**
Valid point about flow control. The early return prevents subsequent update operations from running if specialized endpoints were used.

**Decision:** The current logic is intentional - if we handle updates via specialized endpoints (flavor, vm_state, volumes, etc.), we skip the standard PATCH call. However, we should ensure all operations complete before the final refresh.

**Action:** Review update flow to ensure all operations complete before final state refresh.

---

## Comment 22: Use MarshalJSONForUpdate

**GitHub Link:** (inferred from resource.go Line 1253)

**File:** `internal/services/cloud_instance/resource.go` (Line 1253)

**Pedro's Comment:**
> I think we should do the reverse, not build params directly instead of using option.WithRequestBody. The generated code using the MarshalJSONForUpdate should work as intended.

**Code Context:**
```go
	// Handle other updates using the standard Update endpoint
	// Build params directly instead of using option.WithRequestBody
	params := cloud.InstanceUpdateParams{}
```

**Response:**
The "build params directly" approach was specifically to fix the GET-with-body bug (see Comment 4). If `MarshalJSONForUpdate` now handles this correctly without leaking body to GET requests, we should use the generated approach.

**Action:** Test if the SDK's `MarshalJSONForUpdate` works correctly and refactor if so.

---

## Summary of Actions

| # | Action | Priority |
|---|--------|----------|
| 1 | Revert changes to cloud_instance_image/data_source.go | High |
| 2 | Move helper functions to end of resource.go | Medium |
| 3 | Use `custom.TagsEqual` instead of local duplicate | High |
| 4 | Test if SDK handles WithRequestBody correctly now | Medium |
| 5 | Raise API flavor field inconsistency with backend team | Low |
| 6 | Verify if project_id/region_id extraction is needed | Medium |
| 7 | Add IsUnknown() checks consistently | Medium |
| 8 | Match volumes by ID instead of array index | High |
| 9 | Shorten/remove RequiresReplace comments | Low |
| 10 | Make port_id Optional+Computed with validation | High |
| 11 | Keep UseStateForUnknown on static computed fields | N/A (keep as-is) |
| 12 | Keep comment about addresses UseStateForUnknown | N/A (keep as-is) |
| 13 | Remove tasks attribute | Medium |
| 14 | Review security_groups plan modifier | Medium |
| 15 | Add schema-level validators for interface types | High |
| 16 | Check if SDK supports PortID in FloatingIPNewParams | Medium |
| 17 | Refactor to use ListAutoPaging | Medium |
| 18 | Keep servergroup_id for backward compatibility | N/A (design decision) |
| 19 | Consolidate change flags to single stateHasChanged | Medium |
| 20 | Extract common refresh logic to helper function | Medium |
| 21 | Review update flow for proper completion | Medium |
| 22 | Test MarshalJSONForUpdate and refactor if working | Medium |

---

## Jira Context Summary

All changes in this PR were driven by QA testing (Kirill Tsaregorodtsev) and iteratively fixed based on real infrastructure testing. Key milestones:

| Date | Issue | Resolution |
|------|-------|------------|
| Nov 13 | Flavor/volume changes replace instance | Use /changeflavor and /extend endpoints |
| Nov 18 | VolumeID JSON mapping wrong | Change from json:"volume_id" to json:"id" |
| Nov 19 | Drift in computed fields | Add UseStateForUnknown plan modifiers |
| Dec 2 | GET requests have body | Build params directly (not WithRequestBody) |
| Dec 11 | Volume attach/detach fails | Use /attach and /detach endpoints |
| Dec 18 | addresses element appeared/vanished | Remove UseStateForUnknown from addresses |
| Dec 18 | Tags PATCH empty body | Use JSON Merge Patch |
| Dec 18 | Servergroup forces replacement | Use put_into/remove_from_servergroup |
| Dec 22 | FIP attach/detach errors | Use assign/unassign endpoints |
| Jan 6 | Import boot_index drift | Remove Default from boot_index |

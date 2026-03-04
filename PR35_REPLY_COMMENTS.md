# PR #35 Reply Comments for Manual Posting

Copy-paste these replies to Pedro's comments on GitHub.

---

## Comment 1 (Line 1, cloud_instance_image/data_source.go) - Revert changes

> Fixed - reverted this file to original.

---

## Comment 2 (Line 24, resource.go) - Move functions to end of file

> Fixed - helper functions removed (no longer needed after reverting to WithRequestBody pattern).

---

## Comment 3 (Line 216, resource.go) - Duplicate tagsEqual

> Fixed - removed. Using MarshalJSONForUpdate which handles tags via apijson.

---

## Comment 4 (Line 286, resource.go) - Build params vs WithRequestBody

> Reverted to WithRequestBody pattern. GET-with-body issue will be fixed on SDK level.

---

## Comment 5 (Line 1372, resource.go) - Flavor API inconsistency

> Will create a ticket. API uses "flavor" as string (ID) in create but returns object in GET.

---

## Comment 6 (Line 1380, resource.go) - Extract project_id/region_id

> Fixed - removed project_id/region_id extraction. These are always set in .tf config, already in state from plan.

---

## Comment 7 (Line 1393, resource.go) - Missing Unknown checks

> Fixed - IsUnknown checks already present in the code, verified consistency.

---

## Comment 8 (Line 1407, resource.go) - Order of interfaces

> Fixed - now matching interfaces by port_id when available, with index fallback for initial creation.

---

## Comment 9 (Line 48, schema.go) - Remove RequiresReplace comments

> Fixed - removed all RequiresReplace comments.

---

## Comment 10 (Line 126, schema.go) - port_id should be Optional+Computed

> Fixed - changed port_id to Optional+Computed. Added schema-level validator requiring it for reserved_fixed_ip type.

---

## Comment 11 (Line 218, schema.go) - Plan modifiers for read-only fields

> Yes we need them - without UseStateForUnknown on static computed fields like created_at, Terraform shows "(known after apply)" every plan even though values never change. Added Nov 19 to fix drift.

---

## Comment 12 (Line 286, schema.go) - addresses UseStateForUnknown comment

> Fixed - removed the comment as requested.

---

## Comment 13 (Line 312, schema.go) - Remove tasks attribute

> Fixed - removed tasks attribute from schema.

---

## Comment 14 (Line 321, schema.go) - security_groups plan modifier

> Checked - security_groups doesn't have UseStateForUnknown (it's user-specified, not computed).

---

## Comment 15 (Line 678, resource.go) - Schema-level validation

> Fixed - added schema-level validators for interface types. Now fails early at `terraform plan`:
> - `subnet` type requires `subnet_id`
> - `any_subnet` type requires `network_id`
> - `reserved_fixed_ip` type requires `port_id`

---

## Comment 16 (Line 786, resource.go) - FloatingIP port_id in params

> Fixed - SDK does support `PortID` directly in `FloatingIPNewParams`. Removed `WithRequestBody` workaround.

---

## Comment 17 (Line 850, resource.go) - ListAutoPaging

> Fixed - refactored to use `FloatingIPs.ListAutoPaging()` instead of manual pagination.

---

## Comment 18 (Line 1112, resource.go) - Remove servergroup from instances

> Keeping - QA requested this for compatibility with old provider. Removing would break migration.

---

## Comment 19 (Line 1180, resource.go) - Multiple change flags

> Fixed - consolidated all flags into single `stateHasChanged` boolean.

---

## Comment 20 (Line 1209, resource.go) - Move read logic to function

> Agreed it would improve readability. The read/refresh logic appears in Create, Update, and Read with slight variations. Can extract to helper in follow-up PR to reduce risk in this one.

---

## Comment 21 (Line 1248, resource.go) - Return after partial update

> Early return is intentional. Specialized endpoints handle: resize, vm_state, volumes, interfaces, floating IPs, security groups, servergroups. If any were used (`stateHasChanged=true`), we refresh state and return - no need to call PATCH. Standard PATCH only runs for simple updates (like `name`) when no specialized endpoints were used.

---

## Comment 22 (Line 1253, resource.go) - MarshalJSONForUpdate

> Reverted to MarshalJSONForUpdate + WithRequestBody pattern. GET-with-body issue will be fixed on SDK level.

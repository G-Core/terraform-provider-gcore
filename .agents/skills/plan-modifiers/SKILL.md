---
name: plan-modifiers
description: Conventions for creating and organizing custom Terraform plan modifiers in this provider. Use when creating, modifying, or reviewing plan modifiers, or when deciding where to place a new plan modifier (shared vs resource-specific).
---

# Plan Modifiers

## Placement

- **Shared** (`internal/planmodifiers/{type}_{behavior}.go`): Generic, reusable modifiers with no resource-specific references.
- **Resource-specific** (`internal/services/{resource}/plan_modifiers.go`): Modifiers referencing specific attribute names, model types, or business logic.
- **Exception**: `NormalizeDynamicPlanModifier` in `internal/customfield/dynamic.go` (coupled to custom type).

## Naming

- File: `{type}_{behavior}.go` where `{type}` is `string`, `int64`, `bool`, `list`, `set`, `map`, `object`, or `dynamic`
- Constructor: PascalCase (e.g., `UseNullForRemoval()`)
- Struct: camelCase + Modifier suffix (e.g., `useNullForRemovalModifier`)

## Documentation

Every constructor must have a doc comment explaining what it does and when to use it.

## Testing

Unit tests required for shared modifiers in `internal/planmodifiers/{type}_{behavior}_test.go`.

Cover at minimum:
- Null config value
- Unknown config/plan value
- Explicit config value
- With and without prior state
- Null state (resource creation)

## Shared Modifier Inventory

| Modifier | Type | Description |
|----------|------|-------------|
| `UseNullForRemoval()` | String | Uses null when config removed; state when unknown with prior state; unknown when unknown without prior state (Create) |
| `RequiresReplaceIfConfiguredPreservingState()` | Int64 | Requires replace only when user explicitly sets a different value; preserves state when config is null |
| `UseEmptyListWhenConfigNull()` | List | Sets empty list when config null (for clearing lists) |
| `UseStateForUnknownIncludingNullString()` | String | Preserves state (including null) when plan unknown |
| `UseStateForUnknownIncludingNullObject()` | Object | Preserves state (including null) when plan unknown |
| `{Bool,Int64,String,Set,Object}UseStateForUnknownInclNull()` | Various | Generic variants preserving state (including null) when plan unknown (in `use_state_for_unknown_incl_null.go`) |
| `ObjectPreserveNullState()` | Object | Preserves null state when plan unknown; does not compute when neither config nor state specifies the object |
| `StringRequiresReplaceIfConfiguredPreservingState()` | String | Import-safe replacement: requires replace only when both config and state have known values that differ. Skips replacement when state is null (e.g., after importing a resource with a write-only field like `origin` that the API doesn't return in GET responses). Also skips it when the config value is **removed or unknown** — only use it where the Update API can honour those cases in place |
| `RequiresReplaceIfPriorValueKnown()` | String | Import-safe replacement: skips replacement only when the **prior state** is null/unknown (create or post-import for a `no_refresh` field); replaces on any other change, including removal from config and an unknown planned value. Prefer this one for create-only fields (`string_requires_replace_if_prior_value_known.go`) |
| `SetSuppressServerAdditions()` | Set | Suppresses drift when the API enriches a user-provided set with server-managed elements. If every config element exists in state and state has more, uses state value |
| `UseStateUnlessCountChanges(countAttr)` | List | Preserves list state unless resource replaced or specified count attr changes |
| `RequiresReplaceOnConfigChange()` | Object | Requires replace only when user-specified config fields change (ignores computed) |
| `ListRequiresReplaceIfNotNull()` | List | Import-safe replacement: requires replace when the value changes, but skips it when prior state is null (post-import for a `no_refresh` field) |
| `StringRequiresReplaceUnlessAdopting(key)` | String | Import-safe replacement for create-only fields, keyed on an **explicit private-state marker** rather than on prior state being null. Skips replacement only during the one-time adoption that follows `terraform import`; every other transition keeps built-in semantics, including removal and an unknown planned value. The owning resource sets `key` in `ImportState` and clears it in `Update` (`requires_replace_unless_adopting.go`) |

### Import-safe replacement: prefer the marker

A create-only attribute the API cannot return (`no_refresh`) is null in state after
`terraform import`, so an unconditional `RequiresReplace()` plans a destroy+recreate of live
infrastructure. The tempting fix is to skip replacement whenever prior state is null — but
**prior state is equally null for a resource created with the attribute omitted**. Adopting
silently there records a value the infrastructure does not have, turning a loud, destructive
behaviour into a silent, wrong one.

`StringRequiresReplaceUnlessAdopting()` avoids that by keying on a marker written at import time.
`RequiresReplaceIfPriorValueKnown()` and `ListRequiresReplaceIfNotNull()` implement the older
null-prior-state heuristic and carry that caveat; prefer the marker for new work.

## Resource-Specific Modifier Inventory

| Modifier | Type | Location |
|----------|------|----------|
| `UnknownOnPortChange` | String | `cloud_floating_ip/plan_modifiers.go` |
| `ComputedIfPortSet` | String | `cloud_floating_ip/plan_modifiers.go` |
| `authenticationRemovalPlanModifier` | Object | `cloud_k8s_cluster/plan_modifiers.go` |
| `poolsNormalizeOrderPlanModifier` | List | `cloud_k8s_cluster/plan_modifiers.go` |
| `NormalizeDynamicPlanModifier` | Dynamic | `internal/customfield/dynamic.go` (exception) |

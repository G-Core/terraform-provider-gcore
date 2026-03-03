# Terraform Schema Behaviors Knowledge Base

Complete reference for Terraform Plugin Framework configurability semantics. Use this instead of external searches.

## Table of Contents

1. [Configurability Modes](#1-configurability-modes)
2. [Constraint Rules](#2-constraint-rules)
3. [Plan Behavior](#3-plan-behavior)
4. [Plan Modifiers](#4-plan-modifiers)
5. [Defaults](#5-defaults)
6. [State Consistency Rules](#6-state-consistency-rules)
7. [Decision Framework](#7-decision-framework)
8. [Common Attribute Patterns](#8-common-attribute-patterns)
9. [Contradiction Detection](#9-contradiction-detection)

---

## 1. Configurability Modes

### Required

**Schema:** `Required: true` | **Model tag:** `json:"field,required"`

- Practitioner MUST provide a value
- Value can initially be unknown but must resolve before apply
- Use for: mandatory user input with no sensible default (e.g., `name`, `type`, `region`)

### Optional

**Schema:** `Optional: true` | **Model tag:** `json:"field,optional"`

- Practitioner MAY provide or leave null
- If null, stays null (no provider computation)
- Use for: optional user input where null is valid (e.g., `description`, `tags`)

### Computed

**Schema:** `Computed: true` | **Model tag:** `json:"field,computed"`

- Set ONLY by provider; user config causes error
- Unknown during plan, known after apply → shows "(known after apply)"
- Use for: server-assigned read-only values (e.g., `id`, `created_at`, `status`)

### Optional + Computed

**Schema:** `Optional: true, Computed: true` | **Model tag:** `json:"field,computed_optional"`

- User MAY provide (takes precedence) or leave null (provider may compute)
- Prior state used as fallback for null config values
- Use for: user-overridable server defaults (e.g., `availability_zone`)

## 2. Constraint Rules

**Invalid combinations:**

| Combination | Why Invalid |
|-------------|------------|
| `Required + Optional` | Mutually exclusive |
| `Required + Computed` | Can't require input for computed value |
| All three false | At least one must be true |

**Valid combinations:** Required only, Optional only, Computed only, Optional+Computed.

## 3. Plan Behavior

### When values become unknown

Framework marks computed attributes as unknown when:
- Plan differs from current state AND attribute is computed AND null in config

### Preventing "(known after apply)" noise

Use `UseStateForUnknown()` for Computed/Optional+Computed attributes whose values don't change after creation.

**Do NOT use** for volatile fields (timestamps that change on updates, ETags, versions).

### Plan modification order

1. Set defaults (if null in config)
2. Mark computed as unknown (if plan differs from state)
3. Run attribute plan modifiers
4. Run resource plan modifiers

Defaults run BEFORE unknown marking, so defaults prevent "(known after apply)".

## 4. Plan Modifiers

### UseStateForUnknown()

Copy prior state into plan when plan value is unknown. Use for stable computed values that don't change after creation (e.g., `id`, `created_at`, `fingerprint`).

### RequiresReplace()

Mark resource for destroy+recreate when attribute changes. Use for immutable attributes.

### RequiresReplaceIf()

Conditional replacement based on custom logic.

### RequiresReplaceIfConfigured()

Replacement only when user explicitly configured the attribute (not for provider-computed changes).

## 5. Defaults

### StaticValue defaults

Available for String, Bool, Int64, Float64, List, Map, Set, Object. Applied when config null, runs BEFORE unknown marking.

```go
Default: stringdefault.StaticString("default-value")
```

**Useful for:** Optional, Optional+Computed (prevents unknown).
**NOT useful for:** Required (user provides anyway), Computed-only (provider determines).

## 6. State Consistency Rules

### Apply must match plan

All state values must match planned values or Terraform raises "Provider produced inconsistent result".

- Known plan value → apply must return exact same value
- Unknown plan value → apply can return any known value of correct type
- State must be wholly known (never return unknown)

### Configuration preservation

Non-null config values must preserve exact config value or prior state value. Don't normalize user input unless using Optional+Computed.

### Drift detection

- **Normalization** (same meaning, different form): preserve user's value
- **Drift** (materially different): return remote value

## 7. Decision Framework

### Configurability mode decision tree

```
Q1: Can the user set this value?
├─ NO → Computed
└─ YES → Q2: MUST the user set it?
         ├─ YES → Required
         └─ NO → Q3: Will provider compute a value if unset?
                  ├─ NO → Optional
                  └─ YES → Optional + Computed
```

### Plan modifier decision tree

```
For Computed / Optional+Computed:
  Q1: Value changes after creation?
  ├─ NO → UseStateForUnknown()
  └─ YES → Q2: Changes on EVERY update?
           ├─ YES → No UseStateForUnknown
           └─ NO → Consider UseStateForUnknown if rare

For any mutable attribute:
  Q1: Change requires resource replacement?
  ├─ NO → No RequiresReplace
  └─ YES → Q2: Only when user explicitly configures?
           ├─ YES → RequiresReplaceIfConfigured()
           └─ NO → RequiresReplace()
```

## 8. Common Attribute Patterns

| Pattern | Mode | Plan Modifier | Rationale |
|---------|------|---------------|-----------|
| `id` | Computed | UseStateForUnknown | Server-assigned, stable |
| `name` | Required | RequiresReplace (often) | Mandatory, often immutable |
| `description` | Optional | None | Optional, no default |
| `created_at` | Computed | UseStateForUnknown | Timestamp, never changes |
| `updated_at` | Computed | None | Changes on every update |
| `status` | Computed | UseStateForUnknown | Usually stable |
| `type` | Required | RequiresReplace | Mandatory, usually immutable |
| `tags` | Optional | None | Optional, mutable |
| `availability_zone` | Optional+Computed | UseStateForUnknown | User can set or cloud chooses |
| `instance_type` | Optional+Computed | None | User can override default |
| `fingerprint` | Computed | UseStateForUnknown | Derived, stable |
| `arn` | Computed | UseStateForUnknown | Server-generated ID |
| `etag` | Computed | None | Changes on content change |
| `version` | Computed | None | Increments on updates |

## 9. Contradiction Detection

### Flag contradictions

| Issue | Resolution |
|-------|------------|
| Required + Optional | Remove one; usually keep Required |
| Required + Computed | Invalid; choose based on semantics |
| No flags set | Must set at least one |
| UseStateForUnknown on volatile field | Remove modifier |
| UseStateForUnknown on non-computed | Remove modifier |
| Default on Computed-only | Remove default |
| Default on Required | Remove default |

### Semantic contradictions

| Issue | Resolution |
|-------|------------|
| Mutable field with RequiresReplace | Remove RequiresReplace |
| Immutable field without RequiresReplace | Add RequiresReplace |
| Server-assigned with Optional | Change to Computed |
| User-required with Computed | Change to Required |
| Timestamp user-configurable | Usually should be Computed |

---
name: schema-behaviors
description: Analyzes Terraform resource/data source schemas and produces deterministic configurability recommendations (Required/Optional/Computed flags, plan modifiers, validators) with rationale. Use when analyzing or fixing schema.go/model.go configurability issues, diagnosing "Provider produced inconsistent result" or "(known after apply)" noise, resolving perpetual diffs or unintended drift, or reviewing attribute flag correctness. Triggers on tasks like "fix schema flags", "analyze configurability", "fix known after apply", "fix inconsistent result", or any work reviewing schema attribute behaviors.
---

# Terraform Schema Behaviors

Analyze and fix Terraform Plugin Framework schema configurability (Required/Optional/Computed flags, plan modifiers, validators, defaults) in `schema.go` and `model.go` files.

## Core Objective

Minimize these problems while matching each attribute's semantic intent:

- Spurious "(known after apply)" noise in plan output
- Unintended diffs/drift between plan and apply
- "Provider produced inconsistent result" errors
- Perpetual diffs on subsequent applies

## Workflow

1. **Read the schema file** to extract all attributes and their current flags, modifiers, validators, defaults
2. **Read the knowledge base** at [references/knowledge-base.md](references/knowledge-base.md) for decision framework and constraint rules
3. **For each attribute**, apply the decision tree:
   - Determine semantic intent from name, type, and context
   - Select configurability mode (Required / Optional / Computed / Optional+Computed)
   - Select plan modifiers (UseStateForUnknown, RequiresReplace, etc.)
   - Check for contradictions (flag conflicts, semantic mismatches)
4. **Apply fixes** directly to `schema.go` and `model.go`
5. **Report** changes with rationale

### When analyzing test failures

Map error to root cause before reading schema:
- `"cannot set computed attribute"` → should be Computed, not Optional
- `"missing required argument"` → should be Required, not Optional
- `"Provider produced inconsistent result"` → plan/state mismatch, check Optional+Computed
- `"unexpected value"` → consider Optional+Computed

### When uncertain about semantic intent

Ask the user:
- "Is this attribute set by the user or computed by the API?"
- "Can this value change after the resource is created?"
- "Does the API normalize or modify user-provided values?"

## Output Format

For each attribute analyzed, produce:

```markdown
## Attribute: `<name>`

**Current:** Required=X, Optional=Y, Computed=Z | tag: `json:"field,modifier"` | modifiers: [list] | default: [value]

**Semantic Intent:** [what it represents]

**Rationale:**
1. [consideration]
2. [conclusion]

**Contradictions:** [list or "None"]

**Recommendation:**
- Mode: [Required | Optional | Computed | Optional+Computed]
- Flags: `Required: T/F, Optional: T/F, Computed: T/F`
- Tag: `json:"field,<modifier>"`
- Plan modifiers: [code or "none"]
- Validators: [code or "none"]
- Default: [code or "none"]
- "(known after apply)": [Yes/No/Only on create]
```

## Files Modified

- `internal/services/**/schema.go` -- schema flag definitions
- `internal/services/**/model.go` -- model struct tags

## Guardrails

- Always produce a single deterministic recommendation per attribute
- Always check for contradictions (see knowledge base section 9)
- Always predict plan behavior ("known after apply" appearance)
- Never modify files outside schema.go and model.go
- Never assume API behavior without asking
- This provider uses terraform-plugin-framework (not SDK)

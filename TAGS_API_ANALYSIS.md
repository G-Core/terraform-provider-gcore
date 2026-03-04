# Tags API Analysis - GCLOUD2-20778

**Date:** 2025-11-24
**Analysis:** API Design for Tags in Load Balancers

---

## API Documentation Summary

### From Gcore API Docs

Based on the official API documentation:

#### CREATE Load Balancer
- **Request Field:** `tags`
- **Type:** Object (map/dictionary)
- **Format:** `{"key": "value", "another_key": "another_value"}`
- **Schema:** `CreateTagsSerializer`
- **Optional:** Yes
- **Response:** No tags in response (only task IDs)

#### UPDATE Load Balancer (PATCH)
- **Request Field:** `tags`
- **Type:** Object (map/dictionary with special semantics)
- **Schema:** `UpdateTagsSerializer`
- **Semantics:** JSON Merge Patch (RFC 7386)
  - Add/update: `{"key": "value"}`
  - Delete specific: `{"key": null}`
  - Delete all: `{"tags": null}`
- **Response Field:** `tags_v2`
- **Response Type:** Array of objects
- **Response Format:** `[{"key": "...", "value": "...", "read_only": true/false}]`

#### GET Load Balancer
- **Response Field:** `tags_v2` (required)
- **Type:** Array of `TagSerializer` objects
- **Format:** `[{"key": "...", "value": "...", "read_only": true/false}]`

---

## OpenAPI Schema Definitions

### CreateTagsSerializer
```yaml
CreateTagsSerializer:
  additionalProperties:
    description: Tag value. The maximum size for a value is 1024 bytes.
    maxLength: 1024
    minLength: 1
    type: string
  propertyNames:
    description: Tag key. The maximum size for a key is 255 bytes.
    maxLength: 255
    minLength: 1
  title: CreateTagsSerializer
  type: object
```

**Format:** Map/object with dynamic keys
```json
{
  "my-tag": "my-tag-value",
  "environment": "production"
}
```

### UpdateTagsSerializer
```yaml
UpdateTagsSerializer:
  additionalProperties:
    anyOf:
    - maxLength: 1024
      minLength: 1
      type: string
    - type: 'null'
    description: Tag value (string) or null to delete
  propertyNames:
    description: Tag key. The maximum size for a key is 255 bytes.
    maxLength: 255
    minLength: 1
  title: UpdateTagsSerializer
  type: object
```

**Format:** Map/object with values that can be string or null
```json
{
  "my-tag": "my-tag-value",
  "tag-to-delete": null
}
```

### TagSerializer (in responses)
```yaml
TagSerializer:
  properties:
    key:
      description: Tag key
      maxLength: 255
      type: string
    value:
      description: Tag value
      maxLength: 1024
      type: string
    read_only:
      description: Indicates if the tag is read-only
      type: boolean
  required:
  - key
  - value
  - read_only
  title: TagSerializer
  type: object
```

**Format:** Array of structured objects
```json
[
  {"key": "my-tag", "value": "my-tag-value", "read_only": false},
  {"key": "system-tag", "value": "auto-generated", "read_only": true}
]
```

---

## Why Two Different Formats?

### Input Format (`tags` - Map)
**Advantages:**
- ✅ Simple, intuitive for users
- ✅ Easy to add/update: just set key-value
- ✅ Easy to delete: set to null
- ✅ Matches common JSON/HCL patterns
- ✅ Concise representation

**Example Use Cases:**
```json
// Add tags
{"tags": {"env": "prod", "team": "backend"}}

// Update specific tag
{"tags": {"env": "staging"}}

// Delete specific tag
{"tags": {"env": null}}

// Delete all user tags
{"tags": null}
```

### Output Format (`tags_v2` - Array)
**Advantages:**
- ✅ Can express read-only metadata
- ✅ Preserves order (if needed)
- ✅ Can include additional metadata per tag
- ✅ Differentiates user vs system tags
- ✅ Extensible for future fields

**Example Response:**
```json
{
  "tags_v2": [
    {"key": "env", "value": "prod", "read_only": false},
    {"key": "billing-id", "value": "12345", "read_only": true}
  ]
}
```

---

## Terraform Provider Mapping

### Current Schema (CORRECT)

```go
// User Input - Matches API Request Format
"tags": schema.MapAttribute{
    Description: "Key-value tags...",
    Optional:    true,
    ElementType: types.StringType,
}

// API Response - Matches API Response Format
"tags_v2": schema.ListNestedAttribute{
    Description: "List of key-value tags...",
    Computed:    true,  // Read-only
    CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerTagsV2Model](ctx),
    PlanModifiers: []planmodifier.List{
        listplanmodifier.UseStateForUnknown(),  // Prevents drift
    },
    NestedObject: schema.NestedAttributeObject{
        Attributes: map[string]schema.Attribute{
            "key": schema.StringAttribute{
                Computed: true,
            },
            "value": schema.StringAttribute{
                Computed: true,
            },
            "read_only": schema.BoolAttribute{
                Computed: true,
            },
        },
    },
}
```

---

## Answer: Do We Need Both?

### YES - Both Are Required ✅

**Reason:** Different purposes and directions

1. **`tags` (MapAttribute):**
   - Purpose: User input
   - Direction: Config → Terraform → API Request
   - Format: `{"key": "value"}`
   - Used in: CREATE and UPDATE requests
   - Mutable: Yes (user controls)

2. **`tags_v2` (ListNestedAttribute):**
   - Purpose: API response with metadata
   - Direction: API Response → Terraform → State
   - Format: `[{"key": "...", "value": "...", "read_only": ...}]`
   - Used in: GET responses and after CREATE/UPDATE
   - Mutable: No (API controls)

### Data Flow

```
User Config                 API Request               API Response               Terraform State
───────────                 ───────────               ────────────               ───────────────
tags = {          ──────►   POST/PATCH    ──────►   Response:       ──────►   tags = {
  "env" = "prod"             {                         {                           "env" = "prod"
}                              "tags": {                 "tags_v2": [              }
                                 "env": "prod"              {                     tags_v2 = [
                               }                              "key": "env",          {
                             }                                "value": "prod",       "key" = "env"
                                                              "read_only": false     "value" = "prod"
                                                            }                        "read_only" = false
                                                          ]                        }
                                                        }                        ]
```

### Why Not Just One Field?

**Option 1: Only `tags` (Map)**
- ❌ Cannot express read-only status
- ❌ Cannot show system-managed tags
- ❌ Loses metadata from API

**Option 2: Only `tags_v2` (List)**
- ❌ Awkward for users to write
- ❌ Doesn't match API request format
- ❌ Terraform would need complex conversion

**Option 3: Both `tags` and `tags_v2` (Current)**
- ✅ User-friendly input format
- ✅ Matches API request format exactly
- ✅ Preserves all response metadata
- ✅ Clear separation of concerns
- ✅ Can show read-only tags to user

---

## Potential Issues and Fixes

### Issue 1: User Confusion
**Problem:** Users might wonder why there are two tag fields

**Solution:** Clear documentation
```hcl
resource "gcore_cloud_load_balancer" "example" {
  # Set tags using the simple map format
  tags = {
    "environment" = "production"
    "managed-by"  = "terraform"
  }

  # tags_v2 is read-only and shows all tags (including system tags)
  # You can reference it in outputs but cannot set it directly
}

output "all_tags_with_metadata" {
  value = gcore_cloud_load_balancer.example.tags_v2
}
```

### Issue 2: Drift Detection
**Problem:** If `tags_v2` changes unexpectedly, Terraform detects drift

**Fix:** ✅ Already Applied
- `UseStateForUnknown()` plan modifier on `tags_v2`
- `UnmarshalComputed()` in Read method
- This prevents false drift when API returns tags

### Issue 3: Read-Only Tags Not Visible
**Problem:** System tags (read_only=true) don't show in user's `tags` map

**Solution:** This is correct behavior
- User's `tags` map = what user can control
- System tags visible in `tags_v2` (computed)
- Attempting to modify read-only tags = API error

---

## Best Practices

### For Users
```hcl
resource "gcore_cloud_load_balancer" "lb" {
  # Use 'tags' for managing your tags
  tags = {
    "environment" = "prod"
    "team"        = "backend"
  }
}

# View all tags including system tags
output "all_tags" {
  description = "All tags including system-managed ones"
  value       = gcore_cloud_load_balancer.lb.tags_v2
}

# Filter for read-only tags
output "system_tags" {
  description = "System-managed tags"
  value = [
    for tag in gcore_cloud_load_balancer.lb.tags_v2 :
    tag if tag.read_only
  ]
}
```

### For Provider Development
1. ✅ Keep `tags` as Optional MapAttribute
2. ✅ Keep `tags_v2` as Computed ListNestedAttribute
3. ✅ Use `UseStateForUnknown()` on `tags_v2`
4. ✅ Use `UnmarshalComputed()` in Read
5. ✅ Send `tags` in requests (not `tags_v2`)
6. ✅ Receive `tags_v2` in responses
7. ✅ Never send `tags_v2` to API

---

## Conclusion

**YES, we need both `tags` and `tags_v2`** - they serve different purposes:

| Field | Purpose | Direction | Format | Mutable |
|-------|---------|-----------|--------|---------|
| `tags` | User input | Config→API | Map | Yes |
| `tags_v2` | API response | API→State | Array | No |

This design:
- ✅ Matches API contract exactly
- ✅ Provides user-friendly interface
- ✅ Preserves all metadata
- ✅ Follows Terraform best practices
- ✅ Works correctly in current implementation

**The current provider implementation is CORRECT.**

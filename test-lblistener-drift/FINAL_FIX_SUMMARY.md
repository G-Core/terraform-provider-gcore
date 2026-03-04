FINAL FIX SUMMARY - Minimal Code Change

ISSUE: Configuration drift detected on second terraform apply

ROOT CAUSE: Load Balancer resource Read method using wrong unmarshal function

THE FIX (ONE LINE):
File: internal/services/cloud_load_balancer/resource.go
Line: 269
Changed: apijson.Unmarshal → apijson.UnmarshalComputed

COMPLETE LIST OF CHANGES MADE:

1. Load Balancer resource.go (line 269):
   - Changed Unmarshal to UnmarshalComputed

2. Load Balancer Listener resource.go (line 225):
   - Changed Unmarshal to UnmarshalComputed

3. Load Balancer Listener schema.go:
   - Added Computed: true to 5 fields
   - Added CustomType to 2 list fields

4. Load Balancer Listener model.go:
   - Changed 5 fields to computed_optional JSON tags
   - Changed 2 list fields to use customfield types

WHY MINIMAL:
- Load Balancer schema/model already correct (no changes needed)
- Only Read method needed fixing
- Listener needed full fix (schema + model + Read method)

VERIFICATION:
Run: ./test_drift_fix.sh

Expected result:
First apply: Creates resources
Second apply: Shows "No changes. Your infrastructure matches the configuration."

TECHNICAL EXPLANATION:

apijson.Unmarshal:
- Overwrites ALL fields from API response
- Ignores JSON tags
- Causes drift for computed_optional fields

apijson.UnmarshalComputed:
- Respects JSON tags (computed, computed_optional, optional)
- Only updates fields according to their tags
- Prevents drift for computed_optional fields

This is why BOTH Create/Update AND Read methods must use UnmarshalComputed for resources with computed_optional fields.

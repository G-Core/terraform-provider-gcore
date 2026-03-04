# Verification: v1 vs v2 and OpenAPI Spec Analysis

## Question
Is there a v2 PATCH endpoint for Load Balancers that we should be using instead of v1?

## Answer: NO - v1 is correct

### Evidence from Multiple Sources

#### 1. API Documentation
**Source:** https://gcore.com/docs/api-reference/cloud/load-balancers/update-load-balancer

```
PATCH /cloud/v1/loadbalancers/{project_id}/{region_id}/{load_balancer_id}
```

**Verdict:** ✅ v1 is the official endpoint

#### 2. OpenAPI Specification (openapi.yaml)

**PATCH Endpoint Definition (line 13231-13263):**
```yaml
patch:
  operationId: LoadBalancerInstanceViewSet.patch
  summary: Update load balancer
  description: |-
    Rename load balancer, activate/deactivate logging, update preferred connectivity type
    and/or modify load balancer tags. The request will only process the fields that are
    provided in the request body. Any fields that are not included will remain unchanged.
  responses:
    '200':
      description: Returned load balancer
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/LoadbalancerSerializer'
```

**Verdict:** ✅ PATCH returns `LoadbalancerSerializer`

#### 3. LoadbalancerSerializer Schema (line 83296-83301)

**vrrp_ips field definition:**
```yaml
vrrp_ips:
  description: List of VRRP IP addresses
  items:
    $ref: '#/components/schemas/VRRPIP'
  title: Vrrp Ips
  type: array
```

**Verdict:** ✅ OpenAPI spec says `vrrp_ips` should be included in PATCH response

#### 4. v2 Endpoints - DEPRECATED

From OpenAPI spec (lines 13563, 13601, 13640, 13707):

```yaml
# v2 metadata endpoints are DEPRECATED
description: Please use PATCH `/v1/loadbalancers/{project_id}/{region_id}/{load_balancer_id}` instead
```

**Verdict:** ✅ v2 endpoints redirect to v1 PATCH

### The Critical Discrepancy

|Source|Says vrrp_ips in PATCH response?|
|------|--------------------------------|
|API Docs|✅ YES (shows example with vrrp_ips)|
|OpenAPI Spec|✅ YES (LoadbalancerSerializer includes vrrp_ips)|
|Actual API Response|❌ NO (returns `"vrrp_ips": []`)|

## Proof from Our Demonstration

### What We Tested

```bash
# Using the CORRECT v1 endpoint
PATCH https://api.gcore.com/cloud/v1/loadbalancers/379987/76/{id}
```

### Actual Response

```json
{
  "name": "demo-vrrp-RENAMED",
  "vrrp_ips": [],  // ❌ EMPTY despite spec saying it should have data
  "vip_address": "109.61.125.182",
  ...
}
```

### GET Response (for comparison)

```json
{
  "name": "demo-vrrp-RENAMED",
  "vrrp_ips": [    // ✅ CORRECT - has 2 elements
    {
      "ip_address": "109.61.125.21",
      "subnet_id": "b4c6e91a-9eb9-4c59-ad57-00124eceda00",
      "role": "MASTER"
    },
    {
      "ip_address": "109.61.125.115",
      "subnet_id": "b4c6e91a-9eb9-4c59-ad57-00124eceda00",
      "role": "BACKUP"
    }
  ],
  ...
}
```

## Conclusion

1. ✅ **We're using the CORRECT endpoint** - v1 is the only valid version
2. ✅ **No v2 exists** - old v2 metadata endpoints are deprecated
3. ✅ **OpenAPI spec is correct** - it says vrrp_ips should be included
4. ❌ **API implementation is WRONG** - it returns empty vrrp_ips

## The Bug is in API Implementation

The problem is NOT with:
- Our script (uses correct endpoint)
- The documentation (shows correct schema)
- The OpenAPI spec (defines correct response)

The problem IS with:
- **The actual API backend implementation** - PATCH endpoint does not populate vrrp_ips in response

## Impact

- Terraform provider MUST do GET after PATCH to retrieve vrrp_ips
- SDK consumers cannot rely on PATCH response for computed fields
- API team needs to fix backend to match spec

## Files Verified

- `/Users/user/repos/gcore-terraform/api-schemas/openapi.yaml` (lines 13231-13263, 83296-83301)
- `/Users/user/repos/gcore-terraform/demonstrate_vrrp_ips_issue.sh` (uses v1 endpoint)
- https://gcore.com/docs/api-reference/cloud/load-balancers/update-load-balancer

## Recommendation for API Team

**Fix the PATCH endpoint backend to return complete LoadbalancerSerializer as specified:**

### Option 1: Reuse GET Logic (RECOMMENDED)

```python
# Current (WRONG):
def patch(self, request, load_balancer_id):
    lb = update_load_balancer(load_balancer_id, request.data)
    return Response(LoadbalancerSerializer(lb).data)  # vrrp_ips missing!

# BEST FIX: Reuse the same logic as GET endpoint
def patch(self, request, load_balancer_id):
    update_load_balancer(load_balancer_id, request.data)
    # Call the SAME method that GET uses - ensures consistency
    return self.get(request, load_balancer_id)
```

**Why this is better:**
- ✅ **Guarantees consistency** - GET and PATCH return identical data structure
- ✅ **DRY principle** - Single source of truth for serialization
- ✅ **Includes all computed fields** - If GET has special logic for vrrp_ips, PATCH gets it too
- ✅ **Future-proof** - Any changes to GET automatically apply to PATCH
- ✅ **Matches REST conventions** - PATCH should return the updated resource as GET would show it

### Option 2: Refresh from DB (Works but less ideal)

```python
def patch(self, request, load_balancer_id):
    lb = update_load_balancer(load_balancer_id, request.data)
    lb.refresh_from_db()  # Reload from database
    return Response(LoadbalancerSerializer(lb).data)
```

**Issues with this approach:**
- ⚠️ GET and PATCH might diverge if GET has special query logic
- ⚠️ Doesn't capture relationships/joins that GET might use
- ⚠️ May not trigger computed field logic

### Option 3: Explicit Re-fetch (Also good)

```python
def patch(self, request, load_balancer_id):
    update_load_balancer(load_balancer_id, request.data)
    lb = LoadBalancer.objects.select_related(...).get(id=load_balancer_id)
    return Response(LoadbalancerSerializer(lb).data)
```

**Recommendation:** Use **Option 1** (reuse GET logic) - it's the safest and most maintainable.

The backend likely returns a partial object after update instead of re-fetching the complete resource.

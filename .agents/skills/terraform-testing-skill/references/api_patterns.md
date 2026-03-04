# API Call Patterns Reference

## Router Resource API Patterns

### Interface Management

**Attach Interface:**
```
PATCH /v2/cloud/routers/{router_id}/attach_interface
Body: {
  "network_id": "uuid",
  "subnet_id": "uuid"
}
```

**Detach Interface:**
```
PATCH /v2/cloud/routers/{router_id}/detach_interface
Body: {
  "port_id": "uuid"
}
```

### Route Management

**Update Routes (Replace All):**
```
PATCH /v2/cloud/routers/{router_id}
Body: {
  "routes": [
    {
      "destination": "192.168.1.0/24",
      "nexthop": "10.0.0.1"
    }
  ]
}
```

**Clear All Routes:**
```
PATCH /v2/cloud/routers/{router_id}
Body: {
  "routes": []  // Empty array clears all routes
}
```

## Load Balancer Patterns

### Async Operations (Task-Based)

All load balancer operations return task IDs that must be polled:

**Create Load Balancer:**
```
POST /v2/cloud/loadbalancers
Returns: {
  "tasks": ["task-uuid"],
  "id": "lb-uuid"
}
```

**Poll Task Status:**
```
GET /v2/cloud/tasks/{task_id}
Returns: {
  "status": "COMPLETED|RUNNING|ERROR",
  "result": {...}
}
```

### Update Patterns

**Pool Update (with Unset):**
```
# Step 1: Update fields
PATCH /v2/cloud/loadbalancers/pools/{pool_id}
Body: {
  "name": "new-name",
  "lb_algorithm": "LEAST_CONNECTIONS"
}

# Step 2: Unset fields (if needed)
PATCH /v2/cloud/loadbalancers/pools/{pool_id}/unset
Body: {
  "health_monitor": ["url_path", "expected_codes"]
}
```

## Common Patterns Across Resources

### Synchronous Resources (No Tasks)

- SSH Keys
- FaaS Keys
- Storage SFTP Keys

Pattern:
```
POST/PUT/DELETE /resource/{id}
Returns: Complete resource immediately
```

### Asynchronous Resources (With Tasks)

- Load Balancers
- Listeners
- Pools
- Instances
- Volumes
- Networks

Pattern:
```
POST/PUT/DELETE /resource/{id}
Returns: {"tasks": ["task-id"]}
Poll: GET /tasks/{task-id} until status=COMPLETED
```

## Field Update Strategies

### Replace Strategy
Entire field value is replaced:
```
PATCH /resource/{id}
Body: {"field": "new_value"}
```

### Merge Strategy
Arrays/objects are merged:
```
PATCH /resource/{id}
Body: {"tags": {"new_key": "value"}}  // Adds to existing tags
```

### Clear Strategy
Explicit empty to clear:
```
PATCH /resource/{id}
Body: {"routes": []}  // Clears all routes
```

### Unset Strategy
Separate endpoint to remove fields:
```
PATCH /resource/{id}/unset
Body: {"fields": ["field1", "field2"]}
```

## Verification Patterns

### Check HTTP Method in Logs

```bash
# Good - resource updated in place
grep "PATCH.*routers/abc-123" terraform.log

# Bad - resource recreated
grep "DELETE.*routers/abc-123" terraform.log
grep "POST.*routers" terraform.log
```

### Verify Request Body

```python
# In mitmproxy capture
if call['method'] == 'PATCH' and 'routes' in call['body']:
    if call['body']['routes'] == []:
        print("✓ Routes cleared with empty array")
```

### Task Polling Verification

```bash
# Should see task polling for async resources
grep "GET.*tasks/" terraform.log | wc -l
# Should be > 0 for async operations
```

## Error Handling Patterns

### Task Failure
```json
{
  "status": "ERROR",
  "error": {
    "code": "QUOTA_EXCEEDED",
    "message": "Load balancer quota exceeded"
  }
}
```

### Validation Error
```json
{
  "error": {
    "code": 400,
    "message": "Invalid CIDR format",
    "fields": {
      "routes[0].destination": "Invalid CIDR"
    }
  }
}
```

## Best Practices

1. **Always verify with real API**: Don't trust code analysis alone
2. **Check task completion**: Don't assume immediate success
3. **Validate empty vs null**: `[]` and `null` may behave differently
4. **Monitor HTTP methods**: PATCH = update, DELETE+POST = recreate
5. **Capture full request/response**: Use mitmproxy for debugging
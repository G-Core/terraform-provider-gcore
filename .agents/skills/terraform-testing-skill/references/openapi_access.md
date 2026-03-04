# Accessing OpenAPI Spec and SDK

## OpenAPI Specification (Source of Truth)

The OpenAPI specification that Stainless uses to generate both the Go SDK and Terraform provider is the authoritative source for understanding the API.

### Access via GitHub CLI

```bash
# View the full OpenAPI spec
gh api repos/stainless-sdks/gcore-config/contents/openapi.yml \
  --jq '.content' | base64 -d | less

# Save to local file
gh api repos/stainless-sdks/gcore-config/contents/openapi.yml \
  --jq '.content' | base64 -d > openapi.yml

# Search for specific resource (e.g., routers)
gh api repos/stainless-sdks/gcore-config/contents/openapi.yml \
  --jq '.content' | base64 -d | grep -A 50 "/routers"
```

### Direct GitHub Access

Repository: https://github.com/stainless-sdks/gcore-config/blob/main/openapi.yml

### Key Information to Extract

1. **Operation Types** (sync vs async):
   ```yaml
   # Look for responses that return tasks
   responses:
     '200':
       content:
         application/json:
           schema:
             $ref: '#/components/schemas/TaskResponse'
   ```

2. **Required vs Optional Fields**:
   ```yaml
   properties:
     name:
       type: string
       required: true  # or check 'required' array
     description:
       type: string
       # No required = optional
   ```

3. **Default Values**:
   ```yaml
   properties:
     http_method:
       type: string
       default: GET
       enum: [GET, POST, HEAD]
   ```

4. **Special Operations**:
   ```yaml
   /v2/cloud/routers/{router_id}/attach_interface:
     patch:
       operationId: attachRouterInterface
   ```

## Go SDK (Stainless-Generated)

The Go SDK is auto-generated from the OpenAPI spec and used by the Terraform provider.

### Access SDK Source

```bash
# Clone the SDK
git clone https://github.com/stainless-sdks/gcore-go.git
cd gcore-go

# Examine specific service
ls -la cloud/

# View router service methods
grep -r "func.*Router" cloud/router.go
```

### Key SDK Patterns

1. **Async Methods** (look for *AndPoll):
   ```go
   // Async operation
   func (r *RouterService) NewAndPoll(...)
   func (r *RouterService) UpdateAndPoll(...)
   func (r *RouterService) DeleteAndPoll(...)
   ```

2. **Sync Methods** (direct):
   ```go
   // Sync operation
   func (r *SSHKeyService) New(...)
   func (r *SSHKeyService) Update(...)
   func (r *SSHKeyService) Delete(...)
   ```

3. **Response Types**:
   ```go
   // Check response structure
   type Router struct {
       ID         string    `json:"id"`
       Name       string    `json:"name"`
       Interfaces []Interface `json:"interfaces"`
       // Look for computed fields
   }
   ```

## Terraform Provider Structure

### Resource Locations

```bash
# New provider structure
internal/
├── services/
│   ├── cloud_router/
│   │   ├── resource.go      # Main resource implementation
│   │   ├── model.go         # Data model with JSON tags
│   │   ├── schema.go        # Terraform schema definition
│   │   └── validators.go    # Custom validators
│   ├── cloud_load_balancer/
│   └── ...
```

### Mapping OpenAPI → SDK → Terraform

1. **OpenAPI Definition**:
   ```yaml
   /v2/cloud/routers:
     post:
       operationId: createRouter
       requestBody:
         required: true
         content:
           application/json:
             schema:
               $ref: '#/components/schemas/RouterCreateRequest'
   ```

2. **SDK Method**:
   ```go
   func (r *RouterService) NewAndPoll(ctx context.Context, body RouterNewParams, opts ...option.RequestOption) (*Router, error)
   ```

3. **Terraform Resource**:
   ```go
   func (r *routerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
       router, err := r.client.Cloud.Routers.NewAndPoll(ctx, params, ...)
   }
   ```

## Validation Workflow

### 1. Check OpenAPI Spec

```bash
# Find resource definition
gh api repos/stainless-sdks/gcore-config/contents/openapi.yml \
  --jq '.content' | base64 -d | yq '.paths."/v2/cloud/routers"'
```

### 2. Verify SDK Implementation

```bash
# Check SDK has correct methods
cd gcore-go
grep "func.*RouterService" cloud/router.go
```

### 3. Validate Terraform Mapping

```bash
# Ensure Terraform uses correct SDK methods
grep "NewAndPoll\|UpdateAndPoll" internal/services/cloud_router/resource.go
```

## Common Discrepancies to Check

1. **Missing Operations**:
   - OpenAPI has operation but SDK missing method
   - SDK has method but Terraform doesn't use it

2. **Field Mismatches**:
   - OpenAPI field is required but Terraform marks optional
   - Default values differ between spec and implementation

3. **Type Differences**:
   - OpenAPI: array, Terraform: single value
   - OpenAPI: number, Terraform: string

4. **Async/Sync Confusion**:
   - OpenAPI returns task but Terraform uses sync method
   - No polling for async operations

## Quick Reference Commands

```bash
# Get all router operations
gh api repos/stainless-sdks/gcore-config/contents/openapi.yml \
  --jq '.content' | base64 -d | yq '.paths | with_entries(select(.key | contains("router")))'

# Find all async operations (return tasks)
gh api repos/stainless-sdks/gcore-config/contents/openapi.yml \
  --jq '.content' | base64 -d | grep -B5 "TaskResponse"

# List all resource types
gh api repos/stainless-sdks/gcore-config/contents/openapi.yml \
  --jq '.content' | base64 -d | yq '.paths | keys' | grep "/v2/cloud" | cut -d'/' -f4 | sort -u

# Check field requirements for a schema
gh api repos/stainless-sdks/gcore-config/contents/openapi.yml \
  --jq '.content' | base64 -d | yq '.components.schemas.RouterCreateRequest'
```

## Automated Validation Script

```python
#!/usr/bin/env python3
"""
validate_resource_mapping.py - Validate OpenAPI → SDK → Terraform mapping
"""

import yaml
import subprocess
import json

def get_openapi_spec():
    """Fetch OpenAPI spec from GitHub"""
    result = subprocess.run([
        'gh', 'api', 'repos/stainless-sdks/gcore-config/contents/openapi.yml',
        '--jq', '.content'
    ], capture_output=True, text=True)
    
    spec_base64 = result.stdout
    spec_yaml = subprocess.run(['base64', '-d'], 
                               input=spec_base64, 
                               capture_output=True, 
                               text=True).stdout
    return yaml.safe_load(spec_yaml)

def validate_resource(resource_name):
    """Validate a specific resource implementation"""
    spec = get_openapi_spec()
    
    # Find operations for resource
    resource_paths = [p for p in spec['paths'] if resource_name in p]
    
    for path in resource_paths:
        operations = spec['paths'][path]
        
        # Check if returns tasks (async)
        for method, op in operations.items():
            if method in ['get', 'post', 'put', 'patch', 'delete']:
                responses = op.get('responses', {})
                is_async = any('Task' in str(r) for r in responses.values())
                
                print(f"{method.upper()} {path}:")
                print(f"  Async: {is_async}")
                print(f"  OperationId: {op.get('operationId', 'N/A')}")
                
                # Validate SDK has corresponding method
                if is_async:
                    print(f"  Expected SDK method: {op.get('operationId', 'unknown')}AndPoll")

if __name__ == "__main__":
    import sys
    resource = sys.argv[1] if len(sys.argv) > 1 else "router"
    validate_resource(resource)
```

## Best Practices

1. **Always check OpenAPI first** - It's the source of truth
2. **Verify SDK generation** - Ensure SDK matches OpenAPI
3. **Test actual API behavior** - OpenAPI might have bugs too
4. **Document discrepancies** - File issues when found
5. **Automate validation** - Regular checks prevent drift
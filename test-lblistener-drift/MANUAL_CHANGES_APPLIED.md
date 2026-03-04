Manual Changes Applied for Testing

While waiting for Stainless to regenerate the code, we've manually applied changes to test the fix.

CHANGES MADE

1. Schema file: internal/services/cloud_load_balancer_listener/schema.go
   Added Computed: true to three timeout fields:
   - timeout_client_data (line 85)
   - timeout_member_connect (line 93)
   - timeout_member_data (line 101)

2. Model file: internal/services/cloud_load_balancer_listener/model.go
   Changed JSON tags from "optional" to "computed_optional":
   - TimeoutClientData (line 22)
   - TimeoutMemberConnect (line 23)
   - TimeoutMemberData (line 24)

WHAT WE DID NOT CHANGE

We left these fields as "optional" because they need type changes:
- sni_secret_id (would need types.List instead of *[]types.String)
- user_list (would need custom handling for nested list)

These will be properly handled when Stainless regenerates the code.

BUILD RESULT

Provider built successfully with:
go build -o terraform-provider-gcore

TERRAFORM PLAN RESULT

Plan shows correct behavior for the three timeout fields:
  timeout_client_data    = (known after apply)
  timeout_member_connect = (known after apply)
  timeout_member_data    = (known after apply)

This means Terraform now understands these fields are computed by the API and won't try to send null values for them.

NEXT STEPS

1. Test terraform apply to create resources
2. Test second terraform apply to verify no drift
3. Document results
4. When Stainless regenerates, rebase and get the proper generated code

NOTE

These manual changes are temporary for testing. They will be replaced when Stainless properly generates the code based on the x-stainless-terraform-configurability attributes in the OpenAPI spec.

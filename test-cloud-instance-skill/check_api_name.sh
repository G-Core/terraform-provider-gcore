#!/bin/bash
cd /Users/user/repos/gcore-terraform/test-cloud-instance-skill

export GCORE_API_KEY='21788$1e278ce67b6aa33f178122658b1dd0210d0edff453d348acb9b68bffea6a635b7791925ddda198d5678a4dc20269fe04a263ca92c7e5aa41ea79075f89b66bf6'

# Call API directly to get instance name
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/instances/379987/76/157504a5-5821-4f21-a32d-7a6d1f615eae" | jq '.name, .instance_name, .flavor.flavor_id'

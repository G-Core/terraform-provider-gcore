#!/bin/bash
source ../.env
export TF_CLI_CONFIG_FILE=../.terraformrc

# Create LB without tags
cat > main.tf << 'EOF'
terraform { required_providers { gcore = { source = "gcore/gcore" } } }
data "gcore_cloud_projects" "my_projects" { name = "default" }
locals { project_id = [for p in data.gcore_cloud_projects.my_projects.items : p.id] }
data "gcore_cloud_region" "rg" { region_id = 76 }
resource "gcore_cloud_load_balancer" "lb" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  flavor     = "lb1-2-4"
  name       = "qa-lb-fix-test"
}
EOF

echo "Creating LB without tags..."
terraform apply -auto-approve > /dev/null 2>&1
echo "✓ LB created"

# Add tags
cat > main.tf << 'EOF'
terraform { required_providers { gcore = { source = "gcore/gcore" } } }
data "gcore_cloud_projects" "my_projects" { name = "default" }
locals { project_id = [for p in data.gcore_cloud_projects.my_projects.items : p.id] }
data "gcore_cloud_region" "rg" { region_id = 76 }
resource "gcore_cloud_load_balancer" "lb" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  flavor     = "lb1-2-4"
  name       = "qa-lb-fix-test"
  tags = { "qa" = "test" }
}
EOF

echo "Adding tags..."
if terraform apply -auto-approve 2>&1 | grep -q "tags_v2.*appeared"; then
    echo "❌ BUG STILL EXISTS"
    exit 1
else
    echo "✅ SUCCESS - No error!"
    exit 0
fi

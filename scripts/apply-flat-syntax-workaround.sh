#!/bin/bash
# Apply workaround to flatten security_group block syntax
#
# This script modifies generated code to allow both:
#   - security_group = { name = "...", description = "..." }  (nested, legacy)
#   - name = "...", description = "..."                        (flat, recommended)
#
# Run this after: stainless generate

set -e

REPO_ROOT="/Users/user/repos/gcore-terraform"
MODEL_FILE="$REPO_ROOT/internal/services/cloud_security_group/model.go"
RESOURCE_FILE="$REPO_ROOT/internal/services/cloud_security_group/resource.go"
BACKUP_DIR="$REPO_ROOT/.workaround-backups/$(date +%Y%m%d_%H%M%S)"

echo "============================================================"
echo "Applying Flat Syntax Workaround for Security Groups"
echo "============================================================"
echo ""

# Create backup
echo "📦 Creating backup..."
mkdir -p "$BACKUP_DIR"
cp "$MODEL_FILE" "$BACKUP_DIR/model.go.backup"
cp "$RESOURCE_FILE" "$RESOURCE_FILE.backup"
echo "   Backup saved to: $BACKUP_DIR"
echo ""

# Change 1: Update Name field tag in model.go
echo "🔧 Step 1/5: Updating Name field tag..."
sed -i.tmp 's/Name.*types\.String.*`tfsdk:"name" json:"name,optional"`/Name               types.String                                                            `tfsdk:"name" json:"-"` \/\/ WORKAROUND: Flat alternative to security_group.name/' "$MODEL_FILE"

# Change 2: Update Description field tag
echo "🔧 Step 2/5: Updating Description field tag..."
sed -i.tmp 's/Description.*types\.String.*`tfsdk:"description" json:"description,computed"`/Description        types.String                                                            `tfsdk:"description" json:"-"` \/\/ WORKAROUND: Flat alternative/' "$MODEL_FILE"

# Change 3 & 4: Replace MarshalJSON methods
echo "🔧 Step 3/5: Updating MarshalJSON methods..."
cat > /tmp/marshal_replacement.go << 'MARSHAL_EOF'

func (m CloudSecurityGroupModel) MarshalJSON() (data []byte, err error) {
	// WORKAROUND: Support both nested and flat syntax
	model := m

	// If user provided flat fields, construct nested structure for API
	if model.SecurityGroup == nil && !model.Name.IsNull() {
		model.SecurityGroup = &CloudSecurityGroupSecurityGroupModel{
			Name:        model.Name,
			Description: model.Description,
		}
	}

	return apijson.MarshalRoot(model)
}

func (m CloudSecurityGroupModel) MarshalJSONForUpdate(state CloudSecurityGroupModel) (data []byte, err error) {
	// WORKAROUND: Same logic for updates
	model := m

	if model.SecurityGroup == nil && !model.Name.IsNull() {
		model.SecurityGroup = &CloudSecurityGroupSecurityGroupModel{
			Name:        model.Name,
			Description: model.Description,
		}
	}

	return apijson.MarshalForPatch(model, state)
}
MARSHAL_EOF

# Use perl for multi-line replacement (more reliable than sed)
perl -i -0pe 's/func \(m CloudSecurityGroupModel\) MarshalJSON.*?\n\treturn apijson\.MarshalRoot\(m\)\n\}\n\nfunc \(m CloudSecurityGroupModel\) MarshalJSONForUpdate.*?\n\treturn apijson\.MarshalForPatch\(m, state\)\n\}/`cat \/tmp\/marshal_replacement.go`/s' "$MODEL_FILE"

# Change 5: Update ModifyPlan in resource.go
echo "🔧 Step 4/5: Updating ModifyPlan method..."
cat > /tmp/modifyplan_replacement.go << 'MODIFYPLAN_EOF'

func (r *CloudSecurityGroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *CloudSecurityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() || plan == nil {
		return
	}

	// WORKAROUND: Validate mutually exclusive syntax options
	hasNested := plan.SecurityGroup != nil && !plan.SecurityGroup.Name.IsNull()
	hasFlat := !plan.Name.IsNull()

	if hasNested && hasFlat {
		resp.Diagnostics.AddError(
			"Conflicting security group configuration",
			"Cannot use both 'security_group' block and flat 'name'/'description' attributes.\\n\\n"+
			"Choose ONE syntax:\\n\\n"+
			"Option 1 (Nested): security_group = { name = \\"...\\", description = \\"...\\" }\\n"+
			"Option 2 (Flat):   name = \\"...\\", description = \\"...\\"",
		)
		return
	}

	if !hasNested && !hasFlat {
		resp.Diagnostics.AddError(
			"Missing required security group name",
			"Must specify either 'security_group.name' or 'name'",
		)
		return
	}
}
MODIFYPLAN_EOF

perl -i -0pe 's/func \(r \*CloudSecurityGroupResource\) ModifyPlan.*?\n\n\}/`cat \/tmp\/modifyplan_replacement.go`/s' "$RESOURCE_FILE"

# Cleanup temp files
rm -f /tmp/marshal_replacement.go /tmp/modifyplan_replacement.go
rm -f "$MODEL_FILE.tmp" "$RESOURCE_FILE.tmp"

echo "🔧 Step 5/5: Verifying changes..."

# Verify changes were applied
if grep -q "WORKAROUND: Support both nested and flat syntax" "$MODEL_FILE" && \
   grep -q "WORKAROUND: Validate mutually exclusive" "$RESOURCE_FILE"; then
    echo "✅ Workaround applied successfully!"
    echo ""
    echo "============================================================"
    echo "Testing"
    echo "============================================================"
    echo ""
    echo "Test with flat syntax:"
    echo ""
    echo "  resource \"gcore_cloud_security_group\" \"test\" {"
    echo "    project_id  = 379987"
    echo "    region_id   = 76"
    echo "    name        = \"test-flat-syntax\""
    echo "    description = \"Testing flat syntax\""
    echo "  }"
    echo ""
    echo "Run: cd test-secgroup-rule-edgecases && terraform plan"
    echo ""
    echo "⚠️  IMPORTANT: Rerun this script after 'stainless generate'"
else
    echo "❌ Warning: Workaround may not have been fully applied"
    echo "   Check the backup at: $BACKUP_DIR"
    exit 1
fi

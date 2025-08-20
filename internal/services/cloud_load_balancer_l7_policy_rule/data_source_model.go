// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_l7_policy_rule

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudLoadBalancerL7PolicyRuleDataSourceModel struct {
	L7policyID         types.String                   `tfsdk:"l7policy_id" path:"l7policy_id,required"`
	L7ruleID           types.String                   `tfsdk:"l7rule_id" path:"l7rule_id,required"`
	ProjectID          types.Int64                    `tfsdk:"project_id" path:"project_id,required"`
	RegionID           types.Int64                    `tfsdk:"region_id" path:"region_id,required"`
	CompareType        types.String                   `tfsdk:"compare_type" json:"compare_type,computed"`
	ID                 types.String                   `tfsdk:"id" json:"id,computed"`
	Invert             types.Bool                     `tfsdk:"invert" json:"invert,computed"`
	Key                types.String                   `tfsdk:"key" json:"key,computed"`
	OperatingStatus    types.String                   `tfsdk:"operating_status" json:"operating_status,computed"`
	ProvisioningStatus types.String                   `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	Region             types.String                   `tfsdk:"region" json:"region,computed"`
	TaskID             types.String                   `tfsdk:"task_id" json:"task_id,computed"`
	Type               types.String                   `tfsdk:"type" json:"type,computed"`
	Value              types.String                   `tfsdk:"value" json:"value,computed"`
	Tags               customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
}

func (m *CloudLoadBalancerL7PolicyRuleDataSourceModel) toReadParams(_ context.Context) (params cloud.LoadBalancerL7PolicyRuleGetParams, diags diag.Diagnostics) {
	params = cloud.LoadBalancerL7PolicyRuleGetParams{
		L7policyID: m.L7policyID.ValueString(),
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

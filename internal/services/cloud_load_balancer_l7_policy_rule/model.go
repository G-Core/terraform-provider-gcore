// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_l7_policy_rule

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudLoadBalancerL7PolicyRuleModel struct {
	L7policyID         types.String                   `tfsdk:"l7policy_id" path:"l7policy_id,required"`
	L7ruleID           types.String                   `tfsdk:"l7rule_id" path:"l7rule_id,optional"`
	ProjectID          types.Int64                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID           types.Int64                    `tfsdk:"region_id" path:"region_id,optional"`
	CompareType        types.String                   `tfsdk:"compare_type" json:"compare_type,required"`
	Type               types.String                   `tfsdk:"type" json:"type,required"`
	Value              types.String                   `tfsdk:"value" json:"value,required"`
	Invert             types.Bool                     `tfsdk:"invert" json:"invert,optional"`
	Key                types.String                   `tfsdk:"key" json:"key,optional"`
	Tags               *[]types.String                `tfsdk:"tags" json:"tags,optional"`
	ID                 types.String                   `tfsdk:"id" json:"id,computed"`
	OperatingStatus    types.String                   `tfsdk:"operating_status" json:"operating_status,computed"`
	ProvisioningStatus types.String                   `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	Region             types.String                   `tfsdk:"region" json:"region,computed"`
	TaskID             types.String                   `tfsdk:"task_id" json:"task_id,computed"`
	Tasks              customfield.List[types.String] `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
}

func (m CloudLoadBalancerL7PolicyRuleModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudLoadBalancerL7PolicyRuleModel) MarshalJSONForUpdate(state CloudLoadBalancerL7PolicyRuleModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

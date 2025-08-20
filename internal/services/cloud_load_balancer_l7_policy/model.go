// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_l7_policy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudLoadBalancerL7PolicyModel struct {
	L7policyID         types.String                                                      `tfsdk:"l7policy_id" path:"l7policy_id,optional"`
	ProjectID          types.Int64                                                       `tfsdk:"project_id" path:"project_id,optional"`
	RegionID           types.Int64                                                       `tfsdk:"region_id" path:"region_id,optional"`
	Action             types.String                                                      `tfsdk:"action" json:"action,required"`
	ListenerID         types.String                                                      `tfsdk:"listener_id" json:"listener_id,required"`
	Name               types.String                                                      `tfsdk:"name" json:"name,optional"`
	Position           types.Int64                                                       `tfsdk:"position" json:"position,optional"`
	RedirectHTTPCode   types.Int64                                                       `tfsdk:"redirect_http_code" json:"redirect_http_code,optional"`
	RedirectPoolID     types.String                                                      `tfsdk:"redirect_pool_id" json:"redirect_pool_id,optional"`
	RedirectPrefix     types.String                                                      `tfsdk:"redirect_prefix" json:"redirect_prefix,optional"`
	RedirectURL        types.String                                                      `tfsdk:"redirect_url" json:"redirect_url,optional"`
	Tags               *[]types.String                                                   `tfsdk:"tags" json:"tags,optional"`
	ID                 types.String                                                      `tfsdk:"id" json:"id,computed"`
	OperatingStatus    types.String                                                      `tfsdk:"operating_status" json:"operating_status,computed"`
	ProvisioningStatus types.String                                                      `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	Region             types.String                                                      `tfsdk:"region" json:"region,computed"`
	TaskID             types.String                                                      `tfsdk:"task_id" json:"task_id,computed"`
	Tasks              customfield.List[types.String]                                    `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
	Rules              customfield.NestedObjectList[CloudLoadBalancerL7PolicyRulesModel] `tfsdk:"rules" json:"rules,computed"`
}

func (m CloudLoadBalancerL7PolicyModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudLoadBalancerL7PolicyModel) MarshalJSONForUpdate(state CloudLoadBalancerL7PolicyModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type CloudLoadBalancerL7PolicyRulesModel struct {
	ID                 types.String                   `tfsdk:"id" json:"id,computed"`
	CompareType        types.String                   `tfsdk:"compare_type" json:"compare_type,computed"`
	Invert             types.Bool                     `tfsdk:"invert" json:"invert,computed"`
	Key                types.String                   `tfsdk:"key" json:"key,computed"`
	OperatingStatus    types.String                   `tfsdk:"operating_status" json:"operating_status,computed"`
	ProjectID          types.Int64                    `tfsdk:"project_id" json:"project_id,computed"`
	ProvisioningStatus types.String                   `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	Region             types.String                   `tfsdk:"region" json:"region,computed"`
	RegionID           types.Int64                    `tfsdk:"region_id" json:"region_id,computed"`
	Tags               customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
	TaskID             types.String                   `tfsdk:"task_id" json:"task_id,computed"`
	Type               types.String                   `tfsdk:"type" json:"type,computed"`
	Value              types.String                   `tfsdk:"value" json:"value,computed"`
}

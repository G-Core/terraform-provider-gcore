// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_l7_policy

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudLoadBalancerL7PolicyDataSourceModel struct {
	L7policyID         types.String                                                                `tfsdk:"l7policy_id" path:"l7policy_id,required"`
	ProjectID          types.Int64                                                                 `tfsdk:"project_id" path:"project_id,required"`
	RegionID           types.Int64                                                                 `tfsdk:"region_id" path:"region_id,required"`
	Action             types.String                                                                `tfsdk:"action" json:"action,computed"`
	ID                 types.String                                                                `tfsdk:"id" json:"id,computed"`
	ListenerID         types.String                                                                `tfsdk:"listener_id" json:"listener_id,computed"`
	Name               types.String                                                                `tfsdk:"name" json:"name,computed"`
	OperatingStatus    types.String                                                                `tfsdk:"operating_status" json:"operating_status,computed"`
	Position           types.Int64                                                                 `tfsdk:"position" json:"position,computed"`
	ProvisioningStatus types.String                                                                `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	RedirectHTTPCode   types.Int64                                                                 `tfsdk:"redirect_http_code" json:"redirect_http_code,computed"`
	RedirectPoolID     types.String                                                                `tfsdk:"redirect_pool_id" json:"redirect_pool_id,computed"`
	RedirectPrefix     types.String                                                                `tfsdk:"redirect_prefix" json:"redirect_prefix,computed"`
	RedirectURL        types.String                                                                `tfsdk:"redirect_url" json:"redirect_url,computed"`
	Region             types.String                                                                `tfsdk:"region" json:"region,computed"`
	TaskID             types.String                                                                `tfsdk:"task_id" json:"task_id,computed"`
	Tags               customfield.List[types.String]                                              `tfsdk:"tags" json:"tags,computed"`
	Rules              customfield.NestedObjectList[CloudLoadBalancerL7PolicyRulesDataSourceModel] `tfsdk:"rules" json:"rules,computed"`
}

func (m *CloudLoadBalancerL7PolicyDataSourceModel) toReadParams(_ context.Context) (params cloud.LoadBalancerL7PolicyGetParams, diags diag.Diagnostics) {
	params = cloud.LoadBalancerL7PolicyGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudLoadBalancerL7PolicyRulesDataSourceModel struct {
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

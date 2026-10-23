// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_health_monitor

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudLoadBalancerPoolHealthMonitorDataSourceModel struct {
	ID                 types.String `tfsdk:"id" path:"pool_id,computed"`
	PoolID             types.String `tfsdk:"pool_id" path:"pool_id,required"`
	ProjectID          types.Int64  `tfsdk:"project_id" path:"project_id,optional"`
	RegionID           types.Int64  `tfsdk:"region_id" path:"region_id,optional"`
	AdminStateUp       types.Bool   `tfsdk:"admin_state_up" json:"admin_state_up,computed"`
	Delay              types.Int64  `tfsdk:"delay" json:"delay,computed"`
	DomainName         types.String `tfsdk:"domain_name" json:"domain_name,computed"`
	ExpectedCodes      types.String `tfsdk:"expected_codes" json:"expected_codes,computed"`
	HTTPMethod         types.String `tfsdk:"http_method" json:"http_method,computed"`
	HTTPVersion        types.String `tfsdk:"http_version" json:"http_version,computed"`
	MaxRetries         types.Int64  `tfsdk:"max_retries" json:"max_retries,computed"`
	MaxRetriesDown     types.Int64  `tfsdk:"max_retries_down" json:"max_retries_down,computed"`
	OperatingStatus    types.String `tfsdk:"operating_status" json:"operating_status,computed"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	Timeout            types.Int64  `tfsdk:"timeout" json:"timeout,computed"`
	Type               types.String `tfsdk:"type" json:"type,computed"`
	URLPath            types.String `tfsdk:"url_path" json:"url_path,computed"`
}

func (m *CloudLoadBalancerPoolHealthMonitorDataSourceModel) toReadParams(_ context.Context) (params cloud.LoadBalancerPoolHealthMonitorGetParams, diags diag.Diagnostics) {
	params = cloud.LoadBalancerPoolHealthMonitorGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

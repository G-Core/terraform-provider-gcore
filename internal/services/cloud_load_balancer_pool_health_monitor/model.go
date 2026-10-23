// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_health_monitor

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudLoadBalancerPoolHealthMonitorModel struct {
	ID                 types.String `tfsdk:"id" json:"id,computed"`
	PoolID             types.String `tfsdk:"pool_id" path:"pool_id,required"`
	ProjectID          types.Int64  `tfsdk:"project_id" path:"project_id,optional"`
	RegionID           types.Int64  `tfsdk:"region_id" path:"region_id,optional"`
	Type               types.String `tfsdk:"type" json:"type,required"`
	Delay              types.Int64  `tfsdk:"delay" json:"delay,required"`
	MaxRetries         types.Int64  `tfsdk:"max_retries" json:"max_retries,required"`
	Timeout            types.Int64  `tfsdk:"timeout" json:"timeout,required"`
	DomainName         types.String `tfsdk:"domain_name" json:"domain_name,optional"`
	HTTPVersion        types.String `tfsdk:"http_version" json:"http_version,optional"`
	AdminStateUp       types.Bool   `tfsdk:"admin_state_up" json:"admin_state_up,computed_optional"`
	ExpectedCodes      types.String `tfsdk:"expected_codes" json:"expected_codes,computed_optional"`
	HTTPMethod         types.String `tfsdk:"http_method" json:"http_method,computed_optional"`
	MaxRetriesDown     types.Int64  `tfsdk:"max_retries_down" json:"max_retries_down,computed_optional"`
	URLPath            types.String `tfsdk:"url_path" json:"url_path,computed_optional"`
	OperatingStatus    types.String `tfsdk:"operating_status" json:"operating_status,computed"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
}

func (m CloudLoadBalancerPoolHealthMonitorModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudLoadBalancerPoolHealthMonitorModel) MarshalJSONForUpdate(state CloudLoadBalancerPoolHealthMonitorModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

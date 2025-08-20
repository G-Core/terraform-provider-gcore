// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_health_monitor

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudLoadBalancerPoolHealthMonitorModel struct {
	PoolID         types.String                   `tfsdk:"pool_id" path:"pool_id,required"`
	ProjectID      types.Int64                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID       types.Int64                    `tfsdk:"region_id" path:"region_id,optional"`
	Delay          types.Int64                    `tfsdk:"delay" json:"delay,required"`
	MaxRetries     types.Int64                    `tfsdk:"max_retries" json:"max_retries,required"`
	Timeout        types.Int64                    `tfsdk:"timeout" json:"timeout,required"`
	Type           types.String                   `tfsdk:"type" json:"type,required"`
	ExpectedCodes  types.String                   `tfsdk:"expected_codes" json:"expected_codes,optional"`
	HTTPMethod     types.String                   `tfsdk:"http_method" json:"http_method,optional"`
	MaxRetriesDown types.Int64                    `tfsdk:"max_retries_down" json:"max_retries_down,optional"`
	URLPath        types.String                   `tfsdk:"url_path" json:"url_path,optional"`
	Tasks          customfield.List[types.String] `tfsdk:"tasks" json:"tasks,computed"`
}

func (m CloudLoadBalancerPoolHealthMonitorModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudLoadBalancerPoolHealthMonitorModel) MarshalJSONForUpdate(state CloudLoadBalancerPoolHealthMonitorModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

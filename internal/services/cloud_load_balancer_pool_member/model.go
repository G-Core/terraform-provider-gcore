// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_member

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudLoadBalancerPoolMemberModel struct {
	ID             types.String                   `tfsdk:"id" json:"id,computed"`
	PoolID         types.String                   `tfsdk:"pool_id" path:"pool_id,required"`
	ProjectID      types.Int64                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID       types.Int64                    `tfsdk:"region_id" path:"region_id,optional"`
	Address        types.String                   `tfsdk:"address" json:"address,required"`
	ProtocolPort   types.Int64                    `tfsdk:"protocol_port" json:"protocol_port,required"`
	InstanceID     types.String                   `tfsdk:"instance_id" json:"instance_id,optional"`
	MonitorAddress types.String                   `tfsdk:"monitor_address" json:"monitor_address,optional"`
	MonitorPort    types.Int64                    `tfsdk:"monitor_port" json:"monitor_port,optional"`
	SubnetID       types.String                   `tfsdk:"subnet_id" json:"subnet_id,optional"`
	Weight         types.Int64                    `tfsdk:"weight" json:"weight,optional"`
	AdminStateUp   types.Bool                     `tfsdk:"admin_state_up" json:"admin_state_up,computed_optional"`
	Backup         types.Bool                     `tfsdk:"backup" json:"backup,computed_optional"`
	Tasks          customfield.List[types.String] `tfsdk:"tasks" json:"tasks,computed"`
}

func (m CloudLoadBalancerPoolMemberModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudLoadBalancerPoolMemberModel) MarshalJSONForUpdate(state CloudLoadBalancerPoolMemberModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

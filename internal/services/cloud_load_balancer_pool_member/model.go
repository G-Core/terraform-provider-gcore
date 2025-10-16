// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_member

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudLoadBalancerPoolMemberModel struct {
	ID             types.String `tfsdk:"id" json:"id,computed"`
	PoolID         types.String `tfsdk:"pool_id" json:"-"`
	ProjectID      types.Int64  `tfsdk:"project_id" json:"-"`
	RegionID       types.Int64  `tfsdk:"region_id" json:"-"`
	Address        types.String `tfsdk:"address" json:"address,required"`
	ProtocolPort   types.Int64  `tfsdk:"protocol_port" json:"protocol_port,required"`
	InstanceID     types.String `tfsdk:"instance_id" json:"instance_id,optional"`
	MonitorAddress types.String `tfsdk:"monitor_address" json:"monitor_address,optional"`
	MonitorPort    types.Int64  `tfsdk:"monitor_port" json:"monitor_port,optional"`
	SubnetID       types.String `tfsdk:"subnet_id" json:"subnet_id,computed_optional"`
	AdminStateUp   types.Bool   `tfsdk:"admin_state_up" json:"admin_state_up,computed_optional"`
	Backup         types.Bool   `tfsdk:"backup" json:"backup,computed_optional"`
	Weight         types.Int64  `tfsdk:"weight" json:"weight,computed_optional"`
	OperatingStatus types.String `tfsdk:"operating_status" json:"operating_status,computed"`
}

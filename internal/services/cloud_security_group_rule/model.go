// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group_rule

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudSecurityGroupRuleModel struct {
	ID             types.String `tfsdk:"id" json:"id,computed"`
	GroupID        types.String `tfsdk:"group_id" path:"group_id,required"`
	ProjectID      types.Int64  `tfsdk:"project_id" path:"project_id,optional"`
	RegionID       types.Int64  `tfsdk:"region_id" path:"region_id,optional"`
	Direction      types.String `tfsdk:"direction" json:"direction,required"`
	Description    types.String `tfsdk:"description" json:"description,optional"`
	Ethertype      types.String `tfsdk:"ethertype" json:"ethertype,optional"`
	PortRangeMax   types.Int64  `tfsdk:"port_range_max" json:"port_range_max,optional"`
	PortRangeMin   types.Int64  `tfsdk:"port_range_min" json:"port_range_min,optional"`
	Protocol       types.String `tfsdk:"protocol" json:"protocol,optional"`
	RemoteGroupID  types.String `tfsdk:"remote_group_id" json:"remote_group_id,optional"`
	RemoteIPPrefix types.String `tfsdk:"remote_ip_prefix" json:"remote_ip_prefix,optional"`
}

func (m CloudSecurityGroupRuleModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudSecurityGroupRuleModel) MarshalJSONForUpdate(state CloudSecurityGroupRuleModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

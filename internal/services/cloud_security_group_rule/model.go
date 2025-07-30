// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group_rule

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type CloudSecurityGroupRuleModel struct {
	ID              types.String      `tfsdk:"id" json:"id,computed"`
	GroupID         types.String      `tfsdk:"group_id" path:"group_id,optional"`
	ProjectID       types.Int64       `tfsdk:"project_id" path:"project_id,optional"`
	RegionID        types.Int64       `tfsdk:"region_id" path:"region_id,optional"`
	Description     types.String      `tfsdk:"description" json:"description,optional"`
	Direction       types.String      `tfsdk:"direction" json:"direction,optional"`
	Ethertype       types.String      `tfsdk:"ethertype" json:"ethertype,optional"`
	PortRangeMax    types.Int64       `tfsdk:"port_range_max" json:"port_range_max,optional"`
	PortRangeMin    types.Int64       `tfsdk:"port_range_min" json:"port_range_min,optional"`
	Protocol        types.String      `tfsdk:"protocol" json:"protocol,optional"`
	RemoteGroupID   types.String      `tfsdk:"remote_group_id" json:"remote_group_id,optional"`
	RemoteIPPrefix  types.String      `tfsdk:"remote_ip_prefix" json:"remote_ip_prefix,optional"`
	CreatedAt       timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	RevisionNumber  types.Int64       `tfsdk:"revision_number" json:"revision_number,computed"`
	SecurityGroupID types.String      `tfsdk:"security_group_id" json:"security_group_id,computed"`
	UpdatedAt       timetypes.RFC3339 `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

func (m CloudSecurityGroupRuleModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudSecurityGroupRuleModel) MarshalJSONForUpdate(state CloudSecurityGroupRuleModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

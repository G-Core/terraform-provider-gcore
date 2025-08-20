// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudSecurityGroupModel struct {
	ID                 types.String                                                            `tfsdk:"id" json:"id,computed"`
	ProjectID          types.Int64                                                             `tfsdk:"project_id" path:"project_id,optional"`
	RegionID           types.Int64                                                             `tfsdk:"region_id" path:"region_id,optional"`
	SecurityGroup      *CloudSecurityGroupSecurityGroupModel                                   `tfsdk:"security_group" json:"security_group,required,no_refresh"`
	Instances          *[]types.String                                                         `tfsdk:"instances" json:"instances,optional,no_refresh"`
	Name               types.String                                                            `tfsdk:"name" json:"name,optional"`
	Tags               *map[string]types.String                                                `tfsdk:"tags" json:"tags,optional,no_refresh"`
	ChangedRules       *[]*CloudSecurityGroupChangedRulesModel                                 `tfsdk:"changed_rules" json:"changed_rules,optional,no_refresh"`
	CreatedAt          timetypes.RFC3339                                                       `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Description        types.String                                                            `tfsdk:"description" json:"description,computed"`
	Region             types.String                                                            `tfsdk:"region" json:"region,computed"`
	RevisionNumber     types.Int64                                                             `tfsdk:"revision_number" json:"revision_number,computed"`
	UpdatedAt          timetypes.RFC3339                                                       `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	SecurityGroupRules customfield.NestedObjectList[CloudSecurityGroupSecurityGroupRulesModel] `tfsdk:"security_group_rules" json:"security_group_rules,computed"`
	TagsV2             customfield.NestedObjectList[CloudSecurityGroupTagsV2Model]             `tfsdk:"tags_v2" json:"tags_v2,computed"`
}

func (m CloudSecurityGroupModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudSecurityGroupModel) MarshalJSONForUpdate(state CloudSecurityGroupModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudSecurityGroupSecurityGroupModel struct {
	Name               types.String                                               `tfsdk:"name" json:"name,required"`
	Description        types.String                                               `tfsdk:"description" json:"description,optional"`
	SecurityGroupRules *[]*CloudSecurityGroupSecurityGroupSecurityGroupRulesModel `tfsdk:"security_group_rules" json:"security_group_rules,optional"`
	Tags               *map[string]jsontypes.Normalized                           `tfsdk:"tags" json:"tags,optional"`
}

type CloudSecurityGroupSecurityGroupSecurityGroupRulesModel struct {
	Description    types.String `tfsdk:"description" json:"description,optional"`
	Direction      types.String `tfsdk:"direction" json:"direction,optional"`
	Ethertype      types.String `tfsdk:"ethertype" json:"ethertype,optional"`
	PortRangeMax   types.Int64  `tfsdk:"port_range_max" json:"port_range_max,optional"`
	PortRangeMin   types.Int64  `tfsdk:"port_range_min" json:"port_range_min,optional"`
	Protocol       types.String `tfsdk:"protocol" json:"protocol,optional"`
	RemoteGroupID  types.String `tfsdk:"remote_group_id" json:"remote_group_id,optional"`
	RemoteIPPrefix types.String `tfsdk:"remote_ip_prefix" json:"remote_ip_prefix,optional"`
}

type CloudSecurityGroupChangedRulesModel struct {
	Action              types.String `tfsdk:"action" json:"action,required"`
	Description         types.String `tfsdk:"description" json:"description,optional"`
	Direction           types.String `tfsdk:"direction" json:"direction,optional"`
	Ethertype           types.String `tfsdk:"ethertype" json:"ethertype,optional"`
	PortRangeMax        types.Int64  `tfsdk:"port_range_max" json:"port_range_max,optional"`
	PortRangeMin        types.Int64  `tfsdk:"port_range_min" json:"port_range_min,optional"`
	Protocol            types.String `tfsdk:"protocol" json:"protocol,optional"`
	RemoteGroupID       types.String `tfsdk:"remote_group_id" json:"remote_group_id,optional"`
	RemoteIPPrefix      types.String `tfsdk:"remote_ip_prefix" json:"remote_ip_prefix,optional"`
	SecurityGroupRuleID types.String `tfsdk:"security_group_rule_id" json:"security_group_rule_id,optional"`
}

type CloudSecurityGroupSecurityGroupRulesModel struct {
	ID              types.String      `tfsdk:"id" json:"id,computed"`
	CreatedAt       timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Direction       types.String      `tfsdk:"direction" json:"direction,computed"`
	RevisionNumber  types.Int64       `tfsdk:"revision_number" json:"revision_number,computed"`
	SecurityGroupID types.String      `tfsdk:"security_group_id" json:"security_group_id,computed"`
	UpdatedAt       timetypes.RFC3339 `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Description     types.String      `tfsdk:"description" json:"description,computed"`
	Ethertype       types.String      `tfsdk:"ethertype" json:"ethertype,computed"`
	PortRangeMax    types.Int64       `tfsdk:"port_range_max" json:"port_range_max,computed"`
	PortRangeMin    types.Int64       `tfsdk:"port_range_min" json:"port_range_min,computed"`
	Protocol        types.String      `tfsdk:"protocol" json:"protocol,computed"`
	RemoteGroupID   types.String      `tfsdk:"remote_group_id" json:"remote_group_id,computed"`
	RemoteIPPrefix  types.String      `tfsdk:"remote_ip_prefix" json:"remote_ip_prefix,computed"`
}

type CloudSecurityGroupTagsV2Model struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

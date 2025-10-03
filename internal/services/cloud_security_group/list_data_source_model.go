// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudSecurityGroupsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudSecurityGroupsItemsDataSourceModel] `json:"results,computed"`
}

type CloudSecurityGroupsDataSourceModel struct {
	ProjectID   types.Int64                                                           `tfsdk:"project_id" path:"project_id,optional"`
	RegionID    types.Int64                                                           `tfsdk:"region_id" path:"region_id,optional"`
	TagKeyValue types.String                                                          `tfsdk:"tag_key_value" query:"tag_key_value,optional"`
	TagKey      *[]types.String                                                       `tfsdk:"tag_key" query:"tag_key,optional"`
	MaxItems    types.Int64                                                           `tfsdk:"max_items"`
	Items       customfield.NestedObjectList[CloudSecurityGroupsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudSecurityGroupsDataSourceModel) toListParams(_ context.Context) (params cloud.SecurityGroupListParams, diags diag.Diagnostics) {
	mTagKey := []string{}
	if m.TagKey != nil {
		for _, item := range *m.TagKey {
			mTagKey = append(mTagKey, item.ValueString())
		}
	}

	params = cloud.SecurityGroupListParams{
		TagKey: mTagKey,
	}

	if !m.TagKeyValue.IsNull() {
		params.TagKeyValue = param.NewOpt(m.TagKeyValue.ValueString())
	}
	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudSecurityGroupsItemsDataSourceModel struct {
	ID                 types.String                                                                       `tfsdk:"id" json:"id,computed"`
	CreatedAt          timetypes.RFC3339                                                                  `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Name               types.String                                                                       `tfsdk:"name" json:"name,computed"`
	ProjectID          types.Int64                                                                        `tfsdk:"project_id" json:"project_id,computed"`
	Region             types.String                                                                       `tfsdk:"region" json:"region,computed"`
	RegionID           types.Int64                                                                        `tfsdk:"region_id" json:"region_id,computed"`
	RevisionNumber     types.Int64                                                                        `tfsdk:"revision_number" json:"revision_number,computed"`
	TagsV2             customfield.NestedObjectList[CloudSecurityGroupsTagsV2DataSourceModel]             `tfsdk:"tags_v2" json:"tags_v2,computed"`
	UpdatedAt          timetypes.RFC3339                                                                  `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Description        types.String                                                                       `tfsdk:"description" json:"description,computed"`
	SecurityGroupRules customfield.NestedObjectList[CloudSecurityGroupsSecurityGroupRulesDataSourceModel] `tfsdk:"security_group_rules" json:"security_group_rules,computed"`
}

type CloudSecurityGroupsTagsV2DataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudSecurityGroupsSecurityGroupRulesDataSourceModel struct {
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

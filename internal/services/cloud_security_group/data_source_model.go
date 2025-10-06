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

type CloudSecurityGroupDataSourceModel struct {
	ID                 types.String                                                                      `tfsdk:"id" path:"group_id,computed"`
	GroupID            types.String                                                                      `tfsdk:"group_id" path:"group_id,optional"`
	ProjectID          types.Int64                                                                       `tfsdk:"project_id" path:"project_id,required"`
	RegionID           types.Int64                                                                       `tfsdk:"region_id" path:"region_id,required"`
	CreatedAt          timetypes.RFC3339                                                                 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Description        types.String                                                                      `tfsdk:"description" json:"description,computed"`
	Name               types.String                                                                      `tfsdk:"name" json:"name,computed"`
	Region             types.String                                                                      `tfsdk:"region" json:"region,computed"`
	RevisionNumber     types.Int64                                                                       `tfsdk:"revision_number" json:"revision_number,computed"`
	UpdatedAt          timetypes.RFC3339                                                                 `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	SecurityGroupRules customfield.NestedObjectList[CloudSecurityGroupSecurityGroupRulesDataSourceModel] `tfsdk:"security_group_rules" json:"security_group_rules,computed"`
	TagsV2             customfield.NestedObjectList[CloudSecurityGroupTagsV2DataSourceModel]             `tfsdk:"tags_v2" json:"tags_v2,computed"`
	FindOneBy          *CloudSecurityGroupFindOneByDataSourceModel                                       `tfsdk:"find_one_by"`
}

func (m *CloudSecurityGroupDataSourceModel) toReadParams(_ context.Context) (params cloud.SecurityGroupGetParams, diags diag.Diagnostics) {
	params = cloud.SecurityGroupGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

func (m *CloudSecurityGroupDataSourceModel) toListParams(_ context.Context) (params cloud.SecurityGroupListParams, diags diag.Diagnostics) {
	mFindOneByTagKey := []string{}
	if m.FindOneBy.TagKey != nil {
		for _, item := range *m.FindOneBy.TagKey {
			mFindOneByTagKey = append(mFindOneByTagKey, item.ValueString())
		}
	}

	params = cloud.SecurityGroupListParams{
		TagKey: mFindOneByTagKey,
	}

	if !m.FindOneBy.TagKeyValue.IsNull() {
		params.TagKeyValue = param.NewOpt(m.FindOneBy.TagKeyValue.ValueString())
	}

	return
}

type CloudSecurityGroupSecurityGroupRulesDataSourceModel struct {
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

type CloudSecurityGroupTagsV2DataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudSecurityGroupFindOneByDataSourceModel struct {
	TagKey      *[]types.String `tfsdk:"tag_key" query:"tag_key,optional"`
	TagKeyValue types.String    `tfsdk:"tag_key_value" query:"tag_key_value,optional"`
}

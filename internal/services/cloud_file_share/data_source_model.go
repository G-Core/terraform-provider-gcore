// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudFileShareDataSourceModel struct {
	ID               types.String                                                         `tfsdk:"id" path:"file_share_id,computed"`
	FileShareID      types.String                                                         `tfsdk:"file_share_id" path:"file_share_id,optional"`
	ProjectID        types.Int64                                                          `tfsdk:"project_id" path:"project_id,optional"`
	RegionID         types.Int64                                                          `tfsdk:"region_id" path:"region_id,optional"`
	ConnectionPoint  types.String                                                         `tfsdk:"connection_point" json:"connection_point,computed"`
	CreatedAt        types.String                                                         `tfsdk:"created_at" json:"created_at,computed"`
	CreatorTaskID    types.String                                                         `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Name             types.String                                                         `tfsdk:"name" json:"name,computed"`
	NetworkID        types.String                                                         `tfsdk:"network_id" json:"network_id,computed"`
	NetworkName      types.String                                                         `tfsdk:"network_name" json:"network_name,computed"`
	Protocol         types.String                                                         `tfsdk:"protocol" json:"protocol,computed"`
	Region           types.String                                                         `tfsdk:"region" json:"region,computed"`
	ShareNetworkName types.String                                                         `tfsdk:"share_network_name" json:"share_network_name,computed"`
	Size             types.Int64                                                          `tfsdk:"size" json:"size,computed"`
	Status           types.String                                                         `tfsdk:"status" json:"status,computed"`
	SubnetID         types.String                                                         `tfsdk:"subnet_id" json:"subnet_id,computed"`
	SubnetName       types.String                                                         `tfsdk:"subnet_name" json:"subnet_name,computed"`
	TaskID           types.String                                                         `tfsdk:"task_id" json:"task_id,computed"`
	TypeName         types.String                                                         `tfsdk:"type_name" json:"type_name,computed"`
	VolumeType       types.String                                                         `tfsdk:"volume_type" json:"volume_type,computed"`
	ShareSettings    customfield.NestedObject[CloudFileShareShareSettingsDataSourceModel] `tfsdk:"share_settings" json:"share_settings,computed"`
	Tags             customfield.NestedObjectList[CloudFileShareTagsDataSourceModel]      `tfsdk:"tags" json:"tags,computed"`
	FindOneBy        *CloudFileShareFindOneByDataSourceModel                              `tfsdk:"find_one_by"`
}

func (m *CloudFileShareDataSourceModel) toReadParams(_ context.Context) (params cloud.FileShareGetParams, diags diag.Diagnostics) {
	params = cloud.FileShareGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

func (m *CloudFileShareDataSourceModel) toListParams(_ context.Context) (params cloud.FileShareListParams, diags diag.Diagnostics) {
	params = cloud.FileShareListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.FindOneBy.Limit.IsNull() {
		params.Limit = param.NewOpt(m.FindOneBy.Limit.ValueInt64())
	}
	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}
	if !m.FindOneBy.TypeName.IsNull() {
		params.TypeName = cloud.FileShareListParamsTypeName(m.FindOneBy.TypeName.ValueString())
	}

	return
}

type CloudFileShareShareSettingsDataSourceModel struct {
	TypeName          types.String `tfsdk:"type_name" json:"type_name,computed"`
	AllowedCharacters types.String `tfsdk:"allowed_characters" json:"allowed_characters,computed"`
	PathLength        types.String `tfsdk:"path_length" json:"path_length,computed"`
	RootSquash        types.Bool   `tfsdk:"root_squash" json:"root_squash,computed"`
}

type CloudFileShareTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudFileShareFindOneByDataSourceModel struct {
	Limit    types.Int64  `tfsdk:"limit" query:"limit,computed_optional"`
	Name     types.String `tfsdk:"name" query:"name,optional"`
	TypeName types.String `tfsdk:"type_name" query:"type_name,optional"`
}

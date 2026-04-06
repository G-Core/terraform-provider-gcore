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

type CloudFileSharesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudFileSharesItemsDataSourceModel] `json:"results,computed"`
}

type CloudFileSharesDataSourceModel struct {
	ProjectID types.Int64                                                       `tfsdk:"project_id" path:"project_id,optional"`
	RegionID  types.Int64                                                       `tfsdk:"region_id" path:"region_id,optional"`
	Name      types.String                                                      `tfsdk:"name" query:"name,optional"`
	TypeName  types.String                                                      `tfsdk:"type_name" query:"type_name,optional"`
	MaxItems  types.Int64                                                       `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudFileSharesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudFileSharesDataSourceModel) toListParams(_ context.Context) (params cloud.FileShareListParams, diags diag.Diagnostics) {
	params = cloud.FileShareListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.TypeName.IsNull() {
		params.TypeName = cloud.FileShareListParamsTypeName(m.TypeName.ValueString())
	}

	return
}

type CloudFileSharesItemsDataSourceModel struct {
	ID               types.String                                                          `tfsdk:"id" json:"id,computed"`
	ConnectionPoint  types.String                                                          `tfsdk:"connection_point" json:"connection_point,computed"`
	CreatedAt        types.String                                                          `tfsdk:"created_at" json:"created_at,computed"`
	CreatorTaskID    types.String                                                          `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Name             types.String                                                          `tfsdk:"name" json:"name,computed"`
	NetworkID        types.String                                                          `tfsdk:"network_id" json:"network_id,computed"`
	NetworkName      types.String                                                          `tfsdk:"network_name" json:"network_name,computed"`
	ProjectID        types.Int64                                                           `tfsdk:"project_id" json:"project_id,computed"`
	Protocol         types.String                                                          `tfsdk:"protocol" json:"protocol,computed"`
	Region           types.String                                                          `tfsdk:"region" json:"region,computed"`
	RegionID         types.Int64                                                           `tfsdk:"region_id" json:"region_id,computed"`
	ShareNetworkName types.String                                                          `tfsdk:"share_network_name" json:"share_network_name,computed"`
	ShareSettings    customfield.NestedObject[CloudFileSharesShareSettingsDataSourceModel] `tfsdk:"share_settings" json:"share_settings,computed"`
	Size             types.Int64                                                           `tfsdk:"size" json:"size,computed"`
	Status           types.String                                                          `tfsdk:"status" json:"status,computed"`
	SubnetID         types.String                                                          `tfsdk:"subnet_id" json:"subnet_id,computed"`
	SubnetName       types.String                                                          `tfsdk:"subnet_name" json:"subnet_name,computed"`
	Tags             customfield.NestedObjectList[CloudFileSharesTagsDataSourceModel]      `tfsdk:"tags" json:"tags,computed"`
	TaskID           types.String                                                          `tfsdk:"task_id" json:"task_id,computed"`
	TypeName         types.String                                                          `tfsdk:"type_name" json:"type_name,computed"`
	VolumeType       types.String                                                          `tfsdk:"volume_type" json:"volume_type,computed"`
}

type CloudFileSharesShareSettingsDataSourceModel struct {
	TypeName          types.String `tfsdk:"type_name" json:"type_name,computed"`
	AllowedCharacters types.String `tfsdk:"allowed_characters" json:"allowed_characters,computed"`
	PathLength        types.String `tfsdk:"path_length" json:"path_length,computed"`
	RootSquash        types.Bool   `tfsdk:"root_squash" json:"root_squash,computed"`
}

type CloudFileSharesTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

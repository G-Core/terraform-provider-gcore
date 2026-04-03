// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudNetworksResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudNetworksItemsDataSourceModel] `json:"results,computed"`
}

type CloudNetworksDataSourceModel struct {
	ProjectID   types.Int64                                                     `tfsdk:"project_id" path:"project_id,optional"`
	RegionID    types.Int64                                                     `tfsdk:"region_id" path:"region_id,optional"`
	External    types.Bool                                                      `tfsdk:"external" query:"external,optional"`
	Name        types.String                                                    `tfsdk:"name" query:"name,optional"`
	NetworkType types.String                                                    `tfsdk:"network_type" query:"network_type,optional"`
	TagKeyValue types.String                                                    `tfsdk:"tag_key_value" query:"tag_key_value,optional"`
	TagKey      *[]types.String                                                 `tfsdk:"tag_key" query:"tag_key,optional"`
	Limit       types.Int64                                                     `tfsdk:"limit" query:"limit,computed_optional"`
	OrderBy     types.String                                                    `tfsdk:"order_by" query:"order_by,computed_optional"`
	OwnedBy     types.String                                                    `tfsdk:"owned_by" query:"owned_by,computed_optional"`
	MaxItems    types.Int64                                                     `tfsdk:"max_items"`
	Items       customfield.NestedObjectList[CloudNetworksItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudNetworksDataSourceModel) toListParams(_ context.Context) (params cloud.NetworkListParams, diags diag.Diagnostics) {
	mTagKey := []string{}
	if m.TagKey != nil {
		for _, item := range *m.TagKey {
			mTagKey = append(mTagKey, item.ValueString())
		}
	}

	params = cloud.NetworkListParams{
		TagKey: mTagKey,
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.External.IsNull() {
		params.External = param.NewOpt(m.External.ValueBool())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.NetworkType.IsNull() {
		params.NetworkType = cloud.NetworkListParamsNetworkType(m.NetworkType.ValueString())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = cloud.NetworkListParamsOrderBy(m.OrderBy.ValueString())
	}
	if !m.OwnedBy.IsNull() {
		params.OwnedBy = cloud.NetworkListParamsOwnedBy(m.OwnedBy.ValueString())
	}
	if !m.TagKeyValue.IsNull() {
		params.TagKeyValue = param.NewOpt(m.TagKeyValue.ValueString())
	}

	return
}

type CloudNetworksItemsDataSourceModel struct {
	ID                  types.String                                                   `tfsdk:"id" json:"id,computed"`
	CreatedAt           timetypes.RFC3339                                              `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                   `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Default             types.Bool                                                     `tfsdk:"default" json:"default,computed"`
	External            types.Bool                                                     `tfsdk:"external" json:"external,computed"`
	Mtu                 types.Int64                                                    `tfsdk:"mtu" json:"mtu,computed"`
	Name                types.String                                                   `tfsdk:"name" json:"name,computed"`
	PortSecurityEnabled types.Bool                                                     `tfsdk:"port_security_enabled" json:"port_security_enabled,computed"`
	ProjectID           types.Int64                                                    `tfsdk:"project_id" json:"project_id,computed"`
	Region              types.String                                                   `tfsdk:"region" json:"region,computed"`
	RegionID            types.Int64                                                    `tfsdk:"region_id" json:"region_id,computed"`
	SegmentationID      types.Int64                                                    `tfsdk:"segmentation_id" json:"segmentation_id,computed"`
	Shared              types.Bool                                                     `tfsdk:"shared" json:"shared,computed"`
	Subnets             customfield.List[types.String]                                 `tfsdk:"subnets" json:"subnets,computed"`
	Tags                customfield.NestedObjectList[CloudNetworksTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	Type                types.String                                                   `tfsdk:"type" json:"type,computed"`
	UpdatedAt           timetypes.RFC3339                                              `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudNetworksTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

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

type CloudNetworkDataSourceModel struct {
	ID                  types.String                                                  `tfsdk:"id" path:"network_id,computed"`
	NetworkID           types.String                                                  `tfsdk:"network_id" path:"network_id,optional"`
	ProjectID           types.Int64                                                   `tfsdk:"project_id" path:"project_id,optional"`
	RegionID            types.Int64                                                   `tfsdk:"region_id" path:"region_id,optional"`
	CreatedAt           timetypes.RFC3339                                             `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                  `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Default             types.Bool                                                    `tfsdk:"default" json:"default,computed"`
	External            types.Bool                                                    `tfsdk:"external" json:"external,computed"`
	Mtu                 types.Int64                                                   `tfsdk:"mtu" json:"mtu,computed"`
	Name                types.String                                                  `tfsdk:"name" json:"name,computed"`
	PortSecurityEnabled types.Bool                                                    `tfsdk:"port_security_enabled" json:"port_security_enabled,computed"`
	Region              types.String                                                  `tfsdk:"region" json:"region,computed"`
	SegmentationID      types.Int64                                                   `tfsdk:"segmentation_id" json:"segmentation_id,computed"`
	Shared              types.Bool                                                    `tfsdk:"shared" json:"shared,computed"`
	Type                types.String                                                  `tfsdk:"type" json:"type,computed"`
	UpdatedAt           timetypes.RFC3339                                             `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Subnets             customfield.List[types.String]                                `tfsdk:"subnets" json:"subnets,computed"`
	Tags                customfield.NestedObjectList[CloudNetworkTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	FindOneBy           *CloudNetworkFindOneByDataSourceModel                         `tfsdk:"find_one_by"`
}

func (m *CloudNetworkDataSourceModel) toReadParams(_ context.Context) (params cloud.NetworkGetParams, diags diag.Diagnostics) {
	params = cloud.NetworkGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

func (m *CloudNetworkDataSourceModel) toListParams(_ context.Context) (params cloud.NetworkListParams, diags diag.Diagnostics) {
	mFindOneByTagKey := []string{}
	if m.FindOneBy.TagKey != nil {
		for _, item := range *m.FindOneBy.TagKey {
			mFindOneByTagKey = append(mFindOneByTagKey, item.ValueString())
		}
	}

	params = cloud.NetworkListParams{
		TagKey: mFindOneByTagKey,
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.FindOneBy.External.IsNull() {
		params.External = param.NewOpt(m.FindOneBy.External.ValueBool())
	}
	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}
	if !m.FindOneBy.NetworkType.IsNull() {
		params.NetworkType = cloud.NetworkListParamsNetworkType(m.FindOneBy.NetworkType.ValueString())
	}
	if !m.FindOneBy.OrderBy.IsNull() {
		params.OrderBy = cloud.NetworkListParamsOrderBy(m.FindOneBy.OrderBy.ValueString())
	}
	if !m.FindOneBy.OwnedBy.IsNull() {
		params.OwnedBy = cloud.NetworkListParamsOwnedBy(m.FindOneBy.OwnedBy.ValueString())
	}
	if !m.FindOneBy.TagKeyValue.IsNull() {
		params.TagKeyValue = param.NewOpt(m.FindOneBy.TagKeyValue.ValueString())
	}

	return
}

type CloudNetworkTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudNetworkFindOneByDataSourceModel struct {
	External    types.Bool      `tfsdk:"external" query:"external,optional"`
	Name        types.String    `tfsdk:"name" query:"name,optional"`
	NetworkType types.String    `tfsdk:"network_type" query:"network_type,optional"`
	OrderBy     types.String    `tfsdk:"order_by" query:"order_by,computed_optional"`
	OwnedBy     types.String    `tfsdk:"owned_by" query:"owned_by,computed_optional"`
	TagKey      *[]types.String `tfsdk:"tag_key" query:"tag_key,optional"`
	TagKeyValue types.String    `tfsdk:"tag_key_value" query:"tag_key_value,optional"`
}

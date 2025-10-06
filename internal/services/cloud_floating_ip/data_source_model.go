// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_floating_ip

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudFloatingIPDataSourceModel struct {
	ID                types.String                                                     `tfsdk:"id" path:"floating_ip_id,computed"`
	FloatingIPID      types.String                                                     `tfsdk:"floating_ip_id" path:"floating_ip_id,optional"`
	ProjectID         types.Int64                                                      `tfsdk:"project_id" path:"project_id,required"`
	RegionID          types.Int64                                                      `tfsdk:"region_id" path:"region_id,required"`
	CreatedAt         timetypes.RFC3339                                                `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID     types.String                                                     `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	FixedIPAddress    types.String                                                     `tfsdk:"fixed_ip_address" json:"fixed_ip_address,computed"`
	FloatingIPAddress types.String                                                     `tfsdk:"floating_ip_address" json:"floating_ip_address,computed"`
	PortID            types.String                                                     `tfsdk:"port_id" json:"port_id,computed"`
	Region            types.String                                                     `tfsdk:"region" json:"region,computed"`
	RouterID          types.String                                                     `tfsdk:"router_id" json:"router_id,computed"`
	Status            types.String                                                     `tfsdk:"status" json:"status,computed"`
	TaskID            types.String                                                     `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt         timetypes.RFC3339                                                `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Tags              customfield.NestedObjectList[CloudFloatingIPTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	FindOneBy         *CloudFloatingIPFindOneByDataSourceModel                         `tfsdk:"find_one_by"`
}

func (m *CloudFloatingIPDataSourceModel) toReadParams(_ context.Context) (params cloud.FloatingIPGetParams, diags diag.Diagnostics) {
	params = cloud.FloatingIPGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

func (m *CloudFloatingIPDataSourceModel) toListParams(_ context.Context) (params cloud.FloatingIPListParams, diags diag.Diagnostics) {
	mFindOneByTagKey := []string{}
	if m.FindOneBy.TagKey != nil {
		for _, item := range *m.FindOneBy.TagKey {
			mFindOneByTagKey = append(mFindOneByTagKey, item.ValueString())
		}
	}

	params = cloud.FloatingIPListParams{
		TagKey: mFindOneByTagKey,
	}

	if !m.FindOneBy.TagKeyValue.IsNull() {
		params.TagKeyValue = param.NewOpt(m.FindOneBy.TagKeyValue.ValueString())
	}

	return
}

type CloudFloatingIPTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudFloatingIPFindOneByDataSourceModel struct {
	TagKey      *[]types.String `tfsdk:"tag_key" query:"tag_key,optional"`
	TagKeyValue types.String    `tfsdk:"tag_key_value" query:"tag_key_value,optional"`
}

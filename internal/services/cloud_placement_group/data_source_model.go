// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_placement_group

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudPlacementGroupDataSourceModel struct {
	ID            types.String                                                              `tfsdk:"id" path:"group_id,computed"`
	GroupID       types.String                                                              `tfsdk:"group_id" path:"group_id,required"`
	ProjectID     types.Int64                                                               `tfsdk:"project_id" path:"project_id,optional"`
	RegionID      types.Int64                                                               `tfsdk:"region_id" path:"region_id,optional"`
	Name          types.String                                                              `tfsdk:"name" json:"name,computed"`
	Policy        types.String                                                              `tfsdk:"policy" json:"policy,computed"`
	Region        types.String                                                              `tfsdk:"region" json:"region,computed"`
	ServergroupID types.String                                                              `tfsdk:"servergroup_id" json:"servergroup_id,computed"`
	Instances     customfield.NestedObjectList[CloudPlacementGroupInstancesDataSourceModel] `tfsdk:"instances" json:"instances,computed"`
}

func (m *CloudPlacementGroupDataSourceModel) toReadParams(_ context.Context) (params cloud.PlacementGroupGetParams, diags diag.Diagnostics) {
	params = cloud.PlacementGroupGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudPlacementGroupInstancesDataSourceModel struct {
	InstanceID   types.String `tfsdk:"instance_id" json:"instance_id,computed"`
	InstanceName types.String `tfsdk:"instance_name" json:"instance_name,computed"`
}

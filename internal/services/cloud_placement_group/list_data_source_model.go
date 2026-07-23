// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_placement_group

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudPlacementGroupsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudPlacementGroupsItemsDataSourceModel] `json:"results,computed"`
}

type CloudPlacementGroupsDataSourceModel struct {
	ProjectID types.Int64                                                            `tfsdk:"project_id" path:"project_id,optional"`
	RegionID  types.Int64                                                            `tfsdk:"region_id" path:"region_id,optional"`
	MaxItems  types.Int64                                                            `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudPlacementGroupsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudPlacementGroupsDataSourceModel) toListParams(_ context.Context) (params cloud.PlacementGroupListParams, diags diag.Diagnostics) {
	params = cloud.PlacementGroupListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudPlacementGroupsItemsDataSourceModel struct {
	Instances     customfield.NestedObjectList[CloudPlacementGroupsInstancesDataSourceModel] `tfsdk:"instances" json:"instances,computed"`
	Name          types.String                                                               `tfsdk:"name" json:"name,computed"`
	Policy        types.String                                                               `tfsdk:"policy" json:"policy,computed"`
	ProjectID     types.Int64                                                                `tfsdk:"project_id" json:"project_id,computed"`
	Region        types.String                                                               `tfsdk:"region" json:"region,computed"`
	RegionID      types.Int64                                                                `tfsdk:"region_id" json:"region_id,computed"`
	ServergroupID types.String                                                               `tfsdk:"servergroup_id" json:"servergroup_id,computed"`
}

type CloudPlacementGroupsInstancesDataSourceModel struct {
	InstanceID   types.String `tfsdk:"instance_id" json:"instance_id,computed"`
	InstanceName types.String `tfsdk:"instance_name" json:"instance_name,computed"`
}

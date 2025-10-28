// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_placement_group

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudPlacementGroupModel struct {
	ID            types.String                                                   `tfsdk:"id" json:"-,computed"`
	ServergroupID types.String                                                   `tfsdk:"servergroup_id" json:"servergroup_id,computed"`
	ProjectID     types.Int64                                                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID      types.Int64                                                    `tfsdk:"region_id" path:"region_id,optional"`
	Name          types.String                                                   `tfsdk:"name" json:"name,required"`
	Policy        types.String                                                   `tfsdk:"policy" json:"policy,required"`
	Region        types.String                                                   `tfsdk:"region" json:"region,computed"`
	Instances     customfield.NestedObjectSet[CloudPlacementGroupInstancesModel] `tfsdk:"instances" json:"instances,computed_optional"`
}

func (m CloudPlacementGroupModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudPlacementGroupModel) MarshalJSONForUpdate(state CloudPlacementGroupModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type CloudPlacementGroupInstancesModel struct {
	InstanceID   types.String `tfsdk:"instance_id" json:"instance_id,optional"`
	InstanceName types.String `tfsdk:"instance_name" json:"instance_name,computed"`
}

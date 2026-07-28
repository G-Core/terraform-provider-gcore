// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster_volume

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudGPUVirtualClusterVolumesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudGPUVirtualClusterVolumesItemsDataSourceModel] `json:"results,computed"`
}

type CloudGPUVirtualClusterVolumesDataSourceModel struct {
	ClusterID types.String                                                                    `tfsdk:"cluster_id" path:"cluster_id,required"`
	ProjectID types.Int64                                                                     `tfsdk:"project_id" path:"project_id,optional"`
	RegionID  types.Int64                                                                     `tfsdk:"region_id" path:"region_id,optional"`
	MaxItems  types.Int64                                                                     `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudGPUVirtualClusterVolumesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudGPUVirtualClusterVolumesDataSourceModel) toListParams(_ context.Context) (params cloud.GPUVirtualClusterVolumeListParams, diags diag.Diagnostics) {
	params = cloud.GPUVirtualClusterVolumeListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudGPUVirtualClusterVolumesItemsDataSourceModel struct {
	ID        types.String                                                                   `tfsdk:"id" json:"id,computed"`
	Bootable  types.Bool                                                                     `tfsdk:"bootable" json:"bootable,computed"`
	CreatedAt timetypes.RFC3339                                                              `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Name      types.String                                                                   `tfsdk:"name" json:"name,computed"`
	RootFs    types.Bool                                                                     `tfsdk:"root_fs" json:"root_fs,computed"`
	ServerID  types.String                                                                   `tfsdk:"server_id" json:"server_id,computed"`
	Size      types.Int64                                                                    `tfsdk:"size" json:"size,computed"`
	Status    types.String                                                                   `tfsdk:"status" json:"status,computed"`
	Tags      customfield.NestedObjectList[CloudGPUVirtualClusterVolumesTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	Type      types.String                                                                   `tfsdk:"type" json:"type,computed"`
}

type CloudGPUVirtualClusterVolumesTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

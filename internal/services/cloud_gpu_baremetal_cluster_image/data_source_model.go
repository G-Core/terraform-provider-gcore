// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster_image

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudGPUBaremetalClusterImageDataSourceModel struct {
	ID               types.String                                                                   `tfsdk:"id" path:"image_id,computed"`
	ImageID          types.String                                                                   `tfsdk:"image_id" path:"image_id,required"`
	ProjectID        types.Int64                                                                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID         types.Int64                                                                    `tfsdk:"region_id" path:"region_id,optional"`
	Architecture     types.String                                                                   `tfsdk:"architecture" json:"architecture,computed"`
	CreatedAt        timetypes.RFC3339                                                              `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	GPUDriver        types.String                                                                   `tfsdk:"gpu_driver" json:"gpu_driver,computed"`
	GPUDriverType    types.String                                                                   `tfsdk:"gpu_driver_type" json:"gpu_driver_type,computed"`
	GPUDriverVersion types.String                                                                   `tfsdk:"gpu_driver_version" json:"gpu_driver_version,computed"`
	MinDisk          types.Int64                                                                    `tfsdk:"min_disk" json:"min_disk,computed"`
	MinRam           types.Int64                                                                    `tfsdk:"min_ram" json:"min_ram,computed"`
	Name             types.String                                                                   `tfsdk:"name" json:"name,computed"`
	OsDistro         types.String                                                                   `tfsdk:"os_distro" json:"os_distro,computed"`
	OsType           types.String                                                                   `tfsdk:"os_type" json:"os_type,computed"`
	OsVersion        types.String                                                                   `tfsdk:"os_version" json:"os_version,computed"`
	Size             types.Int64                                                                    `tfsdk:"size" json:"size,computed"`
	SSHKey           types.String                                                                   `tfsdk:"ssh_key" json:"ssh_key,computed"`
	Status           types.String                                                                   `tfsdk:"status" json:"status,computed"`
	UpdatedAt        timetypes.RFC3339                                                              `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Visibility       types.String                                                                   `tfsdk:"visibility" json:"visibility,computed"`
	Tags             customfield.NestedObjectList[CloudGPUBaremetalClusterImageTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
}

func (m *CloudGPUBaremetalClusterImageDataSourceModel) toReadParams(_ context.Context) (params cloud.GPUBaremetalClusterImageGetParams, diags diag.Diagnostics) {
	params = cloud.GPUBaremetalClusterImageGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudGPUBaremetalClusterImageTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

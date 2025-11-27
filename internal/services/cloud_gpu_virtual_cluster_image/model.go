// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster_image

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudGPUVirtualClusterImageModel struct {
	ID               types.String                   `tfsdk:"id" json:"id,computed"`
	ProjectID        types.Int64                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID         types.Int64                    `tfsdk:"region_id" path:"region_id,optional"`
	Name             types.String                   `tfsdk:"name" json:"name,required"`
	URL              types.String                   `tfsdk:"url" json:"url,required,no_refresh"`
	HwFirmwareType   types.String                   `tfsdk:"hw_firmware_type" json:"hw_firmware_type,optional,no_refresh"`
	OsDistro         types.String                   `tfsdk:"os_distro" json:"os_distro,optional"`
	OsVersion        types.String                   `tfsdk:"os_version" json:"os_version,optional"`
	Tags             *map[string]types.String       `tfsdk:"tags" json:"tags,optional,no_refresh"`
	Architecture     types.String                   `tfsdk:"architecture" json:"architecture,computed_optional"`
	CowFormat        types.Bool                     `tfsdk:"cow_format" json:"cow_format,computed_optional,no_refresh"`
	OsType           types.String                   `tfsdk:"os_type" json:"os_type,computed_optional"`
	SSHKey           types.String                   `tfsdk:"ssh_key" json:"ssh_key,computed_optional"`
	CreatedAt        timetypes.RFC3339              `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	GPUDriver        types.String                   `tfsdk:"gpu_driver" json:"gpu_driver,computed"`
	GPUDriverType    types.String                   `tfsdk:"gpu_driver_type" json:"gpu_driver_type,computed"`
	GPUDriverVersion types.String                   `tfsdk:"gpu_driver_version" json:"gpu_driver_version,computed"`
	MinDisk          types.Int64                    `tfsdk:"min_disk" json:"min_disk,computed"`
	MinRam           types.Int64                    `tfsdk:"min_ram" json:"min_ram,computed"`
	Size             types.Int64                    `tfsdk:"size" json:"size,computed"`
	Status           types.String                   `tfsdk:"status" json:"status,computed"`
	TaskID           types.String                   `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt        timetypes.RFC3339              `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Visibility       types.String                   `tfsdk:"visibility" json:"visibility,computed"`
	Tasks            customfield.List[types.String] `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
}

func (m CloudGPUVirtualClusterImageModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudGPUVirtualClusterImageModel) MarshalJSONForUpdate(state CloudGPUVirtualClusterImageModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

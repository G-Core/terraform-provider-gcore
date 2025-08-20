// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance_image

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type CloudInstanceImageModel struct {
	ID             types.String             `tfsdk:"id" json:"id,computed"`
	ImageID        types.String             `tfsdk:"image_id" path:"image_id,required"`
	ProjectID      types.Int64              `tfsdk:"project_id" path:"project_id,optional"`
	RegionID       types.Int64              `tfsdk:"region_id" path:"region_id,optional"`
	HwFirmwareType types.String             `tfsdk:"hw_firmware_type" json:"hw_firmware_type,optional"`
	HwMachineType  types.String             `tfsdk:"hw_machine_type" json:"hw_machine_type,optional"`
	IsBaremetal    types.Bool               `tfsdk:"is_baremetal" json:"is_baremetal,optional"`
	Name           types.String             `tfsdk:"name" json:"name,optional"`
	OsType         types.String             `tfsdk:"os_type" json:"os_type,optional"`
	SSHKey         types.String             `tfsdk:"ssh_key" json:"ssh_key,optional"`
	Tags           *map[string]types.String `tfsdk:"tags" json:"tags,optional,no_refresh"`
	Architecture   types.String             `tfsdk:"architecture" json:"architecture,computed"`
	CreatedAt      timetypes.RFC3339        `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID  types.String             `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Description    types.String             `tfsdk:"description" json:"description,computed"`
	DiskFormat     types.String             `tfsdk:"disk_format" json:"disk_format,computed"`
	DisplayOrder   types.Int64              `tfsdk:"display_order" json:"display_order,computed"`
	MinDisk        types.Int64              `tfsdk:"min_disk" json:"min_disk,computed"`
	MinRam         types.Int64              `tfsdk:"min_ram" json:"min_ram,computed"`
	OsDistro       types.String             `tfsdk:"os_distro" json:"os_distro,computed"`
	OsVersion      types.String             `tfsdk:"os_version" json:"os_version,computed"`
	Region         types.String             `tfsdk:"region" json:"region,computed"`
	Size           types.Int64              `tfsdk:"size" json:"size,computed"`
	Status         types.String             `tfsdk:"status" json:"status,computed"`
	TaskID         types.String             `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt      timetypes.RFC3339        `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Visibility     types.String             `tfsdk:"visibility" json:"visibility,computed"`
}

func (m CloudInstanceImageModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudInstanceImageModel) MarshalJSONForUpdate(state CloudInstanceImageModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

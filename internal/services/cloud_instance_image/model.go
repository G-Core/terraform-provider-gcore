// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance_image

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudInstanceImageModel struct {
	ID               types.String             `tfsdk:"id" json:"id,computed"`
	ProjectID        types.Int64              `tfsdk:"project_id" path:"project_id,optional"`
	RegionID         types.Int64              `tfsdk:"region_id" path:"region_id,optional"`
	URL              types.String             `tfsdk:"url" json:"url,required,no_refresh"`
	OsDistro         types.String             `tfsdk:"os_distro" json:"os_distro,optional"`
	OsVersion        types.String             `tfsdk:"os_version" json:"os_version,optional"`
	Architecture     types.String             `tfsdk:"architecture" json:"architecture,computed_optional"`
	CowFormat        types.Bool               `tfsdk:"cow_format" json:"cow_format,computed_optional,no_refresh"`
	Name             types.String             `tfsdk:"name" json:"name,required"`
	HwFirmwareType   types.String             `tfsdk:"hw_firmware_type" json:"hw_firmware_type,optional"`
	HwMachineType    types.String             `tfsdk:"hw_machine_type" json:"hw_machine_type,optional"`
	IsBaremetal      types.Bool               `tfsdk:"is_baremetal" json:"is_baremetal,computed_optional"`
	OsType           types.String             `tfsdk:"os_type" json:"os_type,computed_optional"`
	SSHKey           types.String             `tfsdk:"ssh_key" json:"ssh_key,computed_optional"`
	Tags             customfield.Map[types.String]  `tfsdk:"tags" json:"tags,computed_optional,no_refresh"`
	CreatedAt        timetypes.RFC3339        `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID    types.String             `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	CurrencyCode     types.String                   `tfsdk:"currency_code" json:"currency_code,computed"`
	Description      types.String             `tfsdk:"description" json:"description,computed"`
	DiskFormat       types.String             `tfsdk:"disk_format" json:"disk_format,computed"`
	DisplayOrder     types.Int64              `tfsdk:"display_order" json:"display_order,computed"`
	GPUDriver        types.String             `tfsdk:"gpu_driver" json:"gpu_driver,computed"`
	GPUDriverType    types.String             `tfsdk:"gpu_driver_type" json:"gpu_driver_type,computed"`
	GPUDriverVersion types.String             `tfsdk:"gpu_driver_version" json:"gpu_driver_version,computed"`
	MinDisk          types.Int64              `tfsdk:"min_disk" json:"min_disk,computed"`
	MinRam           types.Int64              `tfsdk:"min_ram" json:"min_ram,computed"`
	PricePerHour     types.Float64                  `tfsdk:"price_per_hour" json:"price_per_hour,computed"`
	PricePerMonth    types.Float64                  `tfsdk:"price_per_month" json:"price_per_month,computed"`
	PriceStatus      types.String                   `tfsdk:"price_status" json:"price_status,computed"`
	Region           types.String             `tfsdk:"region" json:"region,computed"`
	Size             types.Int64              `tfsdk:"size" json:"size,computed"`
	Status           types.String             `tfsdk:"status" json:"status,computed"`
	UpdatedAt        timetypes.RFC3339        `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Visibility       types.String             `tfsdk:"visibility" json:"visibility,computed"`
}

func (m CloudInstanceImageModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudInstanceImageModel) MarshalJSONForUpdate(state CloudInstanceImageModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

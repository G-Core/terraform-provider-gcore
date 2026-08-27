// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_image

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudBaremetalImagesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudBaremetalImagesItemsDataSourceModel] `json:"results,computed"`
}

type CloudBaremetalImagesDataSourceModel struct {
	ProjectID     types.Int64                                                            `tfsdk:"project_id" path:"project_id,optional"`
	RegionID      types.Int64                                                            `tfsdk:"region_id" path:"region_id,optional"`
	Architecture  types.String                                                           `tfsdk:"architecture" query:"architecture,optional"`
	IncludePrices types.Bool                                                             `tfsdk:"include_prices" query:"include_prices,optional"`
	Name          types.String                                                           `tfsdk:"name" query:"name,optional"`
	OsDistro      types.String                                                           `tfsdk:"os_distro" query:"os_distro,optional"`
	OsVersion     types.String                                                           `tfsdk:"os_version" query:"os_version,optional"`
	Private       types.String                                                           `tfsdk:"private" query:"private,optional"`
	TagKeyValue   types.String                                                           `tfsdk:"tag_key_value" query:"tag_key_value,optional"`
	Visibility    types.String                                                           `tfsdk:"visibility" query:"visibility,optional"`
	TagKey        *[]types.String                                                        `tfsdk:"tag_key" query:"tag_key,optional"`
	MaxItems      types.Int64                                                            `tfsdk:"max_items"`
	Items         customfield.NestedObjectList[CloudBaremetalImagesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudBaremetalImagesDataSourceModel) toListParams(_ context.Context) (params cloud.BaremetalImageListParams, diags diag.Diagnostics) {
	mTagKey := []string{}
	if m.TagKey != nil {
		for _, item := range *m.TagKey {
			mTagKey = append(mTagKey, item.ValueString())
		}
	}

	params = cloud.BaremetalImageListParams{
		TagKey: mTagKey,
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.Architecture.IsNull() {
		params.Architecture = cloud.BaremetalImageListParamsArchitecture(m.Architecture.ValueString())
	}
	if !m.IncludePrices.IsNull() {
		params.IncludePrices = param.NewOpt(m.IncludePrices.ValueBool())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.OsDistro.IsNull() {
		params.OsDistro = param.NewOpt(m.OsDistro.ValueString())
	}
	if !m.OsVersion.IsNull() {
		params.OsVersion = param.NewOpt(m.OsVersion.ValueString())
	}
	if !m.Private.IsNull() {
		params.Private = param.NewOpt(m.Private.ValueString())
	}
	if !m.TagKeyValue.IsNull() {
		params.TagKeyValue = param.NewOpt(m.TagKeyValue.ValueString())
	}
	if !m.Visibility.IsNull() {
		params.Visibility = cloud.BaremetalImageListParamsVisibility(m.Visibility.ValueString())
	}

	return
}

type CloudBaremetalImagesItemsDataSourceModel struct {
	ID               types.String                                                          `tfsdk:"id" json:"id,computed"`
	Architecture     types.String                                                          `tfsdk:"architecture" json:"architecture,computed"`
	CreatedAt        timetypes.RFC3339                                                     `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID    types.String                                                          `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	CurrencyCode     types.String                                                          `tfsdk:"currency_code" json:"currency_code,computed"`
	Description      types.String                                                          `tfsdk:"description" json:"description,computed"`
	DiskFormat       types.String                                                          `tfsdk:"disk_format" json:"disk_format,computed"`
	DisplayOrder     types.Int64                                                           `tfsdk:"display_order" json:"display_order,computed"`
	GPUDriver        types.String                                                          `tfsdk:"gpu_driver" json:"gpu_driver,computed"`
	GPUDriverType    types.String                                                          `tfsdk:"gpu_driver_type" json:"gpu_driver_type,computed"`
	GPUDriverVersion types.String                                                          `tfsdk:"gpu_driver_version" json:"gpu_driver_version,computed"`
	HwFirmwareType   types.String                                                          `tfsdk:"hw_firmware_type" json:"hw_firmware_type,computed"`
	HwMachineType    types.String                                                          `tfsdk:"hw_machine_type" json:"hw_machine_type,computed"`
	IsBaremetal      types.Bool                                                            `tfsdk:"is_baremetal" json:"is_baremetal,computed"`
	MinDisk          types.Int64                                                           `tfsdk:"min_disk" json:"min_disk,computed"`
	MinRam           types.Int64                                                           `tfsdk:"min_ram" json:"min_ram,computed"`
	Name             types.String                                                          `tfsdk:"name" json:"name,computed"`
	OsDistro         types.String                                                          `tfsdk:"os_distro" json:"os_distro,computed"`
	OsType           types.String                                                          `tfsdk:"os_type" json:"os_type,computed"`
	OsVersion        types.String                                                          `tfsdk:"os_version" json:"os_version,computed"`
	PricePerHour     types.Float64                                                         `tfsdk:"price_per_hour" json:"price_per_hour,computed"`
	PricePerMonth    types.Float64                                                         `tfsdk:"price_per_month" json:"price_per_month,computed"`
	PriceStatus      types.String                                                          `tfsdk:"price_status" json:"price_status,computed"`
	ProjectID        types.Int64                                                           `tfsdk:"project_id" json:"project_id,computed"`
	Region           types.String                                                          `tfsdk:"region" json:"region,computed"`
	RegionID         types.Int64                                                           `tfsdk:"region_id" json:"region_id,computed"`
	Size             types.Int64                                                           `tfsdk:"size" json:"size,computed"`
	SSHKey           types.String                                                          `tfsdk:"ssh_key" json:"ssh_key,computed"`
	Status           types.String                                                          `tfsdk:"status" json:"status,computed"`
	Tags             customfield.NestedObjectList[CloudBaremetalImagesTagsDataSourceModel] `tfsdk:"tags" json:"tags_v2,computed"`
	TaskID           types.String                                                          `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt        timetypes.RFC3339                                                     `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Visibility       types.String                                                          `tfsdk:"visibility" json:"visibility,computed"`
}

type CloudBaremetalImagesTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

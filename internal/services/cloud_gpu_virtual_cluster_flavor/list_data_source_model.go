// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster_flavor

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudGPUVirtualClusterFlavorsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudGPUVirtualClusterFlavorsItemsDataSourceModel] `json:"results,computed"`
}

type CloudGPUVirtualClusterFlavorsDataSourceModel struct {
	ProjectID     types.Int64                                                                     `tfsdk:"project_id" path:"project_id,optional"`
	RegionID      types.Int64                                                                     `tfsdk:"region_id" path:"region_id,optional"`
	HideDisabled  types.Bool                                                                      `tfsdk:"hide_disabled" query:"hide_disabled,computed_optional"`
	IncludePrices types.Bool                                                                      `tfsdk:"include_prices" query:"include_prices,computed_optional"`
	MaxItems      types.Int64                                                                     `tfsdk:"max_items"`
	Items         customfield.NestedObjectList[CloudGPUVirtualClusterFlavorsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudGPUVirtualClusterFlavorsDataSourceModel) toListParams(_ context.Context) (params cloud.GPUVirtualClusterFlavorListParams, diags diag.Diagnostics) {
	params = cloud.GPUVirtualClusterFlavorListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.HideDisabled.IsNull() {
		params.HideDisabled = param.NewOpt(m.HideDisabled.ValueBool())
	}
	if !m.IncludePrices.IsNull() {
		params.IncludePrices = param.NewOpt(m.IncludePrices.ValueBool())
	}

	return
}

type CloudGPUVirtualClusterFlavorsItemsDataSourceModel struct {
	Architecture        types.String                                                                              `tfsdk:"architecture" json:"architecture,computed"`
	Capacity            types.Int64                                                                               `tfsdk:"capacity" json:"capacity,computed"`
	Disabled            types.Bool                                                                                `tfsdk:"disabled" json:"disabled,computed"`
	HardwareDescription customfield.NestedObject[CloudGPUVirtualClusterFlavorsHardwareDescriptionDataSourceModel] `tfsdk:"hardware_description" json:"hardware_description,computed"`
	HardwareProperties  customfield.NestedObject[CloudGPUVirtualClusterFlavorsHardwarePropertiesDataSourceModel]  `tfsdk:"hardware_properties" json:"hardware_properties,computed"`
	Name                types.String                                                                              `tfsdk:"name" json:"name,computed"`
	SupportedFeatures   customfield.NestedObject[CloudGPUVirtualClusterFlavorsSupportedFeaturesDataSourceModel]   `tfsdk:"supported_features" json:"supported_features,computed"`
	Price               customfield.NestedObject[CloudGPUVirtualClusterFlavorsPriceDataSourceModel]               `tfsdk:"price" json:"price,computed"`
}

type CloudGPUVirtualClusterFlavorsHardwareDescriptionDataSourceModel struct {
	GPU          types.String `tfsdk:"gpu" json:"gpu,computed"`
	LocalStorage types.Int64  `tfsdk:"local_storage" json:"local_storage,computed"`
	Ram          types.Int64  `tfsdk:"ram" json:"ram,computed"`
	Vcpus        types.Int64  `tfsdk:"vcpus" json:"vcpus,computed"`
}

type CloudGPUVirtualClusterFlavorsHardwarePropertiesDataSourceModel struct {
	GPUCount        types.Int64  `tfsdk:"gpu_count" json:"gpu_count,computed"`
	GPUManufacturer types.String `tfsdk:"gpu_manufacturer" json:"gpu_manufacturer,computed"`
	GPUModel        types.String `tfsdk:"gpu_model" json:"gpu_model,computed"`
	NicEth          types.String `tfsdk:"nic_eth" json:"nic_eth,computed"`
	NicIb           types.String `tfsdk:"nic_ib" json:"nic_ib,computed"`
}

type CloudGPUVirtualClusterFlavorsSupportedFeaturesDataSourceModel struct {
	SecurityGroups types.Bool `tfsdk:"security_groups" json:"security_groups,computed"`
}

type CloudGPUVirtualClusterFlavorsPriceDataSourceModel struct {
	CurrencyCode  types.String  `tfsdk:"currency_code" json:"currency_code,computed"`
	PricePerHour  types.Float64 `tfsdk:"price_per_hour" json:"price_per_hour,computed"`
	PricePerMonth types.Float64 `tfsdk:"price_per_month" json:"price_per_month,computed"`
	PriceStatus   types.String  `tfsdk:"price_status" json:"price_status,computed"`
}

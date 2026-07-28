// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster_flavor

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudGPUBaremetalClusterFlavorsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudGPUBaremetalClusterFlavorsItemsDataSourceModel] `json:"results,computed"`
}

type CloudGPUBaremetalClusterFlavorsDataSourceModel struct {
	ProjectID     types.Int64                                                                       `tfsdk:"project_id" path:"project_id,optional"`
	RegionID      types.Int64                                                                       `tfsdk:"region_id" path:"region_id,optional"`
	HideDisabled  types.Bool                                                                        `tfsdk:"hide_disabled" query:"hide_disabled,computed_optional"`
	IncludePrices types.Bool                                                                        `tfsdk:"include_prices" query:"include_prices,computed_optional"`
	MaxItems      types.Int64                                                                       `tfsdk:"max_items"`
	Items         customfield.NestedObjectList[CloudGPUBaremetalClusterFlavorsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudGPUBaremetalClusterFlavorsDataSourceModel) toListParams(_ context.Context) (params cloud.GPUBaremetalClusterFlavorListParams, diags diag.Diagnostics) {
	params = cloud.GPUBaremetalClusterFlavorListParams{}

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

type CloudGPUBaremetalClusterFlavorsItemsDataSourceModel struct {
	Architecture        types.String                                                                                `tfsdk:"architecture" json:"architecture,computed"`
	Capacity            types.Int64                                                                                 `tfsdk:"capacity" json:"capacity,computed"`
	Disabled            types.Bool                                                                                  `tfsdk:"disabled" json:"disabled,computed"`
	HardwareDescription customfield.NestedObject[CloudGPUBaremetalClusterFlavorsHardwareDescriptionDataSourceModel] `tfsdk:"hardware_description" json:"hardware_description,computed"`
	HardwareProperties  customfield.NestedObject[CloudGPUBaremetalClusterFlavorsHardwarePropertiesDataSourceModel]  `tfsdk:"hardware_properties" json:"hardware_properties,computed"`
	Name                types.String                                                                                `tfsdk:"name" json:"name,computed"`
	ReservedCapacity    types.Int64                                                                                 `tfsdk:"reserved_capacity" json:"reserved_capacity,computed"`
	SupportedFeatures   customfield.NestedObject[CloudGPUBaremetalClusterFlavorsSupportedFeaturesDataSourceModel]   `tfsdk:"supported_features" json:"supported_features,computed"`
	Price               customfield.NestedObject[CloudGPUBaremetalClusterFlavorsPriceDataSourceModel]               `tfsdk:"price" json:"price,computed"`
}

type CloudGPUBaremetalClusterFlavorsHardwareDescriptionDataSourceModel struct {
	CPU     types.String `tfsdk:"cpu" json:"cpu,computed"`
	Disk    types.String `tfsdk:"disk" json:"disk,computed"`
	GPU     types.String `tfsdk:"gpu" json:"gpu,computed"`
	Network types.String `tfsdk:"network" json:"network,computed"`
	Ram     types.String `tfsdk:"ram" json:"ram,computed"`
}

type CloudGPUBaremetalClusterFlavorsHardwarePropertiesDataSourceModel struct {
	GPUCount        types.Int64  `tfsdk:"gpu_count" json:"gpu_count,computed"`
	GPUManufacturer types.String `tfsdk:"gpu_manufacturer" json:"gpu_manufacturer,computed"`
	GPUModel        types.String `tfsdk:"gpu_model" json:"gpu_model,computed"`
	NicEth          types.String `tfsdk:"nic_eth" json:"nic_eth,computed"`
	NicIb           types.String `tfsdk:"nic_ib" json:"nic_ib,computed"`
}

type CloudGPUBaremetalClusterFlavorsSupportedFeaturesDataSourceModel struct {
	SecurityGroups types.Bool `tfsdk:"security_groups" json:"security_groups,computed"`
}

type CloudGPUBaremetalClusterFlavorsPriceDataSourceModel struct {
	CurrencyCode  types.String  `tfsdk:"currency_code" json:"currency_code,computed"`
	PricePerHour  types.Float64 `tfsdk:"price_per_hour" json:"price_per_hour,computed"`
	PricePerMonth types.Float64 `tfsdk:"price_per_month" json:"price_per_month,computed"`
	PriceStatus   types.String  `tfsdk:"price_status" json:"price_status,computed"`
}

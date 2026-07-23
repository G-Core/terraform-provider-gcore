// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance_flavor

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudInstanceFlavorsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudInstanceFlavorsItemsDataSourceModel] `json:"results,computed"`
}

type CloudInstanceFlavorsDataSourceModel struct {
	ProjectID      types.Int64                                                            `tfsdk:"project_id" path:"project_id,optional"`
	RegionID       types.Int64                                                            `tfsdk:"region_id" path:"region_id,optional"`
	Disabled       types.Bool                                                             `tfsdk:"disabled" query:"disabled,computed_optional"`
	ExcludeLinux   types.Bool                                                             `tfsdk:"exclude_linux" query:"exclude_linux,computed_optional"`
	ExcludeWindows types.Bool                                                             `tfsdk:"exclude_windows" query:"exclude_windows,computed_optional"`
	IncludePrices  types.Bool                                                             `tfsdk:"include_prices" query:"include_prices,computed_optional"`
	MaxItems       types.Int64                                                            `tfsdk:"max_items"`
	Items          customfield.NestedObjectList[CloudInstanceFlavorsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudInstanceFlavorsDataSourceModel) toListParams(_ context.Context) (params cloud.InstanceFlavorListParams, diags diag.Diagnostics) {
	params = cloud.InstanceFlavorListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.Disabled.IsNull() {
		params.Disabled = param.NewOpt(m.Disabled.ValueBool())
	}
	if !m.ExcludeLinux.IsNull() {
		params.ExcludeLinux = param.NewOpt(m.ExcludeLinux.ValueBool())
	}
	if !m.ExcludeWindows.IsNull() {
		params.ExcludeWindows = param.NewOpt(m.ExcludeWindows.ValueBool())
	}
	if !m.IncludePrices.IsNull() {
		params.IncludePrices = param.NewOpt(m.IncludePrices.ValueBool())
	}

	return
}

type CloudInstanceFlavorsItemsDataSourceModel struct {
	Architecture        types.String                  `tfsdk:"architecture" json:"architecture,computed"`
	Disabled            types.Bool                    `tfsdk:"disabled" json:"disabled,computed"`
	FlavorID            types.String                  `tfsdk:"flavor_id" json:"flavor_id,computed"`
	FlavorName          types.String                  `tfsdk:"flavor_name" json:"flavor_name,computed"`
	HardwareDescription customfield.Map[types.String] `tfsdk:"hardware_description" json:"hardware_description,computed"`
	OsType              types.String                  `tfsdk:"os_type" json:"os_type,computed"`
	Ram                 types.Int64                   `tfsdk:"ram" json:"ram,computed"`
	Vcpus               types.Int64                   `tfsdk:"vcpus" json:"vcpus,computed"`
	CurrencyCode        types.String                  `tfsdk:"currency_code" json:"currency_code,computed"`
	PricePerHour        types.Float64                 `tfsdk:"price_per_hour" json:"price_per_hour,computed"`
	PricePerMonth       types.Float64                 `tfsdk:"price_per_month" json:"price_per_month,computed"`
	PriceStatus         types.String                  `tfsdk:"price_status" json:"price_status,computed"`
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_flavor

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudBaremetalFlavorsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudBaremetalFlavorsItemsDataSourceModel] `json:"results,computed"`
}

type CloudBaremetalFlavorsDataSourceModel struct {
	ProjectID               types.Int64                                                             `tfsdk:"project_id" path:"project_id,optional"`
	RegionID                types.Int64                                                             `tfsdk:"region_id" path:"region_id,optional"`
	Disabled                types.Bool                                                              `tfsdk:"disabled" query:"disabled,computed_optional"`
	ExcludeLinux            types.Bool                                                              `tfsdk:"exclude_linux" query:"exclude_linux,computed_optional"`
	ExcludeWindows          types.Bool                                                              `tfsdk:"exclude_windows" query:"exclude_windows,computed_optional"`
	IncludeCapacity         types.Bool                                                              `tfsdk:"include_capacity" query:"include_capacity,computed_optional"`
	IncludePrices           types.Bool                                                              `tfsdk:"include_prices" query:"include_prices,computed_optional"`
	IncludeReservationStock types.Bool                                                              `tfsdk:"include_reservation_stock" query:"include_reservation_stock,computed_optional"`
	MaxItems                types.Int64                                                             `tfsdk:"max_items"`
	Items                   customfield.NestedObjectList[CloudBaremetalFlavorsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudBaremetalFlavorsDataSourceModel) toListParams(_ context.Context) (params cloud.BaremetalFlavorListParams, diags diag.Diagnostics) {
	params = cloud.BaremetalFlavorListParams{}

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
	if !m.IncludeCapacity.IsNull() {
		params.IncludeCapacity = param.NewOpt(m.IncludeCapacity.ValueBool())
	}
	if !m.IncludePrices.IsNull() {
		params.IncludePrices = param.NewOpt(m.IncludePrices.ValueBool())
	}
	if !m.IncludeReservationStock.IsNull() {
		params.IncludeReservationStock = param.NewOpt(m.IncludeReservationStock.ValueBool())
	}

	return
}

type CloudBaremetalFlavorsItemsDataSourceModel struct {
	Architecture        types.String                  `tfsdk:"architecture" json:"architecture,computed"`
	Disabled            types.Bool                    `tfsdk:"disabled" json:"disabled,computed"`
	FlavorID            types.String                  `tfsdk:"flavor_id" json:"flavor_id,computed"`
	FlavorName          types.String                  `tfsdk:"flavor_name" json:"flavor_name,computed"`
	OsType              types.String                  `tfsdk:"os_type" json:"os_type,computed"`
	Ram                 types.Int64                   `tfsdk:"ram" json:"ram,computed"`
	ResourceClass       types.String                  `tfsdk:"resource_class" json:"resource_class,computed"`
	Vcpus               types.Int64                   `tfsdk:"vcpus" json:"vcpus,computed"`
	Capacity            types.Int64                   `tfsdk:"capacity" json:"capacity,computed"`
	CurrencyCode        types.String                  `tfsdk:"currency_code" json:"currency_code,computed"`
	HardwareDescription customfield.Map[types.String] `tfsdk:"hardware_description" json:"hardware_description,computed"`
	PricePerHour        types.Float64                 `tfsdk:"price_per_hour" json:"price_per_hour,computed"`
	PricePerMonth       types.Float64                 `tfsdk:"price_per_month" json:"price_per_month,computed"`
	PriceStatus         types.String                  `tfsdk:"price_status" json:"price_status,computed"`
	ReservedCapacity    types.Int64                   `tfsdk:"reserved_capacity" json:"reserved_capacity,computed"`
}

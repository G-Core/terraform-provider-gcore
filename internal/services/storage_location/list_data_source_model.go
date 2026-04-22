// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_location

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/storage"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageLocationsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[StorageLocationsItemsDataSourceModel] `json:"results,computed"`
}

type StorageLocationsDataSourceModel struct {
	OrderBy  types.String                                                       `tfsdk:"order_by" query:"order_by,computed_optional"`
	MaxItems types.Int64                                                        `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[StorageLocationsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *StorageLocationsDataSourceModel) toListParams(_ context.Context) (params storage.LocationListParams, diags diag.Diagnostics) {
	params = storage.LocationListParams{}

	if !m.OrderBy.IsNull() {
		params.OrderBy = param.NewOpt(m.OrderBy.ValueString())
	}

	return
}

type StorageLocationsItemsDataSourceModel struct {
	Address       types.String `tfsdk:"address" json:"address,computed"`
	Name          types.String `tfsdk:"name" json:"name,computed"`
	TechnicalName types.String `tfsdk:"technical_name" json:"technical_name,computed"`
	Title         types.String `tfsdk:"title" json:"title,computed"`
	Type          types.String `tfsdk:"type" json:"type,computed"`
}

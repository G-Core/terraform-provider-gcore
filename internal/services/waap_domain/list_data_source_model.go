// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WaapDomainsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapDomainsItemsDataSourceModel] `json:"results,computed"`
}

type WaapDomainsDataSourceModel struct {
	Name     types.String                                                  `tfsdk:"name" query:"name,optional"`
	Ordering types.String                                                  `tfsdk:"ordering" query:"ordering,optional"`
	Status   types.String                                                  `tfsdk:"status" query:"status,optional"`
	IDs      *[]types.Int64                                                `tfsdk:"ids" query:"ids,optional"`
	Limit    types.Int64                                                   `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems types.Int64                                                   `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[WaapDomainsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapDomainsDataSourceModel) toListParams(_ context.Context) (params waap.DomainListParams, diags diag.Diagnostics) {
	mIDs := []int64{}
	if m.IDs != nil {
		for _, item := range *m.IDs {
			mIDs = append(mIDs, item.ValueInt64())
		}
	}

	params = waap.DomainListParams{
		IDs: mIDs,
	}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = waap.DomainListParamsOrdering(m.Ordering.ValueString())
	}
	if !m.Status.IsNull() {
		params.Status = waap.DomainListParamsStatus(m.Status.ValueString())
	}

	return
}

type WaapDomainsItemsDataSourceModel struct {
	ID            types.Int64       `tfsdk:"id" json:"id,computed"`
	CreatedAt     timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CustomPageSet types.Int64       `tfsdk:"custom_page_set" json:"custom_page_set,computed"`
	Name          types.String      `tfsdk:"name" json:"name,computed"`
	Status        types.String      `tfsdk:"status" json:"status,computed"`
}

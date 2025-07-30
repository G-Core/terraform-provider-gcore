// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_discovery_scan_result

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainAPIDiscoveryScanResultsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapDomainAPIDiscoveryScanResultsItemsDataSourceModel] `json:"results,computed"`
}

type WaapDomainAPIDiscoveryScanResultsDataSourceModel struct {
	DomainID types.Int64                                                                         `tfsdk:"domain_id" path:"domain_id,required"`
	Message  types.String                                                                        `tfsdk:"message" query:"message,optional"`
	Status   types.String                                                                        `tfsdk:"status" query:"status,optional"`
	Type     types.String                                                                        `tfsdk:"type" query:"type,optional"`
	Limit    types.Int64                                                                         `tfsdk:"limit" query:"limit,computed_optional"`
	Ordering types.String                                                                        `tfsdk:"ordering" query:"ordering,computed_optional"`
	MaxItems types.Int64                                                                         `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[WaapDomainAPIDiscoveryScanResultsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapDomainAPIDiscoveryScanResultsDataSourceModel) toListParams(_ context.Context) (params waap.DomainAPIDiscoveryScanResultListParams, diags diag.Diagnostics) {
	params = waap.DomainAPIDiscoveryScanResultListParams{}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Message.IsNull() {
		params.Message = param.NewOpt(m.Message.ValueString())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = waap.DomainAPIDiscoveryScanResultListParamsOrdering(m.Ordering.ValueString())
	}
	if !m.Status.IsNull() {
		params.Status = waap.DomainAPIDiscoveryScanResultListParamsStatus(m.Status.ValueString())
	}
	if !m.Type.IsNull() {
		params.Type = waap.DomainAPIDiscoveryScanResultListParamsType(m.Type.ValueString())
	}

	return
}

type WaapDomainAPIDiscoveryScanResultsItemsDataSourceModel struct {
	ID        types.String      `tfsdk:"id" json:"id,computed"`
	EndTime   timetypes.RFC3339 `tfsdk:"end_time" json:"end_time,computed" format:"date-time"`
	Message   types.String      `tfsdk:"message" json:"message,computed"`
	StartTime timetypes.RFC3339 `tfsdk:"start_time" json:"start_time,computed" format:"date-time"`
	Status    types.String      `tfsdk:"status" json:"status,computed"`
	Type      types.String      `tfsdk:"type" json:"type,computed"`
}

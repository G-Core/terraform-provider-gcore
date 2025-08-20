// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_path

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainAPIPathsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapDomainAPIPathsItemsDataSourceModel] `json:"results,computed"`
}

type WaapDomainAPIPathsDataSourceModel struct {
	DomainID   types.Int64                                                          `tfsdk:"domain_id" path:"domain_id,required"`
	APIGroup   types.String                                                         `tfsdk:"api_group" query:"api_group,optional"`
	APIVersion types.String                                                         `tfsdk:"api_version" query:"api_version,optional"`
	HTTPScheme types.String                                                         `tfsdk:"http_scheme" query:"http_scheme,optional"`
	Method     types.String                                                         `tfsdk:"method" query:"method,optional"`
	Ordering   types.String                                                         `tfsdk:"ordering" query:"ordering,optional"`
	Path       types.String                                                         `tfsdk:"path" query:"path,optional"`
	Source     types.String                                                         `tfsdk:"source" query:"source,optional"`
	IDs        *[]types.String                                                      `tfsdk:"ids" query:"ids,optional"`
	Status     *[]types.String                                                      `tfsdk:"status" query:"status,optional"`
	Limit      types.Int64                                                          `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems   types.Int64                                                          `tfsdk:"max_items"`
	Items      customfield.NestedObjectList[WaapDomainAPIPathsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapDomainAPIPathsDataSourceModel) toListParams(_ context.Context) (params waap.DomainAPIPathListParams, diags diag.Diagnostics) {
	mIDs := []string{}
	if m.IDs != nil {
		for _, item := range *m.IDs {
			mIDs = append(mIDs, item.ValueString())
		}
	}
	mStatus := []string{}
	if m.Status != nil {
		for _, item := range *m.Status {
			mStatus = append(mStatus, string(item.ValueString()))
		}
	}

	params = waap.DomainAPIPathListParams{
		IDs:    mIDs,
		Status: mStatus,
	}

	if !m.APIGroup.IsNull() {
		params.APIGroup = param.NewOpt(m.APIGroup.ValueString())
	}
	if !m.APIVersion.IsNull() {
		params.APIVersion = param.NewOpt(m.APIVersion.ValueString())
	}
	if !m.HTTPScheme.IsNull() {
		params.HTTPScheme = waap.DomainAPIPathListParamsHTTPScheme(m.HTTPScheme.ValueString())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Method.IsNull() {
		params.Method = waap.DomainAPIPathListParamsMethod(m.Method.ValueString())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = waap.DomainAPIPathListParamsOrdering(m.Ordering.ValueString())
	}
	if !m.Path.IsNull() {
		params.Path = param.NewOpt(m.Path.ValueString())
	}
	if !m.Source.IsNull() {
		params.Source = waap.DomainAPIPathListParamsSource(m.Source.ValueString())
	}

	return
}

type WaapDomainAPIPathsItemsDataSourceModel struct {
	ID            types.String                   `tfsdk:"id" json:"id,computed"`
	APIGroups     customfield.List[types.String] `tfsdk:"api_groups" json:"api_groups,computed"`
	APIVersion    types.String                   `tfsdk:"api_version" json:"api_version,computed"`
	FirstDetected timetypes.RFC3339              `tfsdk:"first_detected" json:"first_detected,computed" format:"date-time"`
	HTTPScheme    types.String                   `tfsdk:"http_scheme" json:"http_scheme,computed"`
	LastDetected  timetypes.RFC3339              `tfsdk:"last_detected" json:"last_detected,computed" format:"date-time"`
	Method        types.String                   `tfsdk:"method" json:"method,computed"`
	Path          types.String                   `tfsdk:"path" json:"path,computed"`
	RequestCount  types.Int64                    `tfsdk:"request_count" json:"request_count,computed"`
	Source        types.String                   `tfsdk:"source" json:"source,computed"`
	Status        types.String                   `tfsdk:"status" json:"status,computed"`
	Tags          customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
}

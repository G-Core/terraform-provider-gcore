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

type WaapDomainAPIPathDataSourceModel struct {
	ID            types.String                               `tfsdk:"id" path:"path_id,computed"`
	PathID        types.String                               `tfsdk:"path_id" path:"path_id,optional"`
	DomainID      types.Int64                                `tfsdk:"domain_id" path:"domain_id,required"`
	APIVersion    types.String                               `tfsdk:"api_version" json:"api_version,computed"`
	FirstDetected timetypes.RFC3339                          `tfsdk:"first_detected" json:"first_detected,computed" format:"date-time"`
	HTTPScheme    types.String                               `tfsdk:"http_scheme" json:"http_scheme,computed"`
	LastDetected  timetypes.RFC3339                          `tfsdk:"last_detected" json:"last_detected,computed" format:"date-time"`
	Method        types.String                               `tfsdk:"method" json:"method,computed"`
	Path          types.String                               `tfsdk:"path" json:"path,computed"`
	RequestCount  types.Int64                                `tfsdk:"request_count" json:"request_count,computed"`
	Source        types.String                               `tfsdk:"source" json:"source,computed"`
	Status        types.String                               `tfsdk:"status" json:"status,computed"`
	APIGroups     customfield.List[types.String]             `tfsdk:"api_groups" json:"api_groups,computed"`
	Tags          customfield.List[types.String]             `tfsdk:"tags" json:"tags,computed"`
	FindOneBy     *WaapDomainAPIPathFindOneByDataSourceModel `tfsdk:"find_one_by"`
}

func (m *WaapDomainAPIPathDataSourceModel) toReadParams(_ context.Context) (params waap.DomainAPIPathGetParams, diags diag.Diagnostics) {
	params = waap.DomainAPIPathGetParams{
		DomainID: m.DomainID.ValueInt64(),
	}

	return
}

func (m *WaapDomainAPIPathDataSourceModel) toListParams(_ context.Context) (params waap.DomainAPIPathListParams, diags diag.Diagnostics) {
	mFindOneByIDs := []string{}
	if m.FindOneBy.IDs != nil {
		for _, item := range *m.FindOneBy.IDs {
			mFindOneByIDs = append(mFindOneByIDs, item.ValueString())
		}
	}
	mFindOneByStatus := []string{}
	if m.FindOneBy.Status != nil {
		for _, item := range *m.FindOneBy.Status {
			mFindOneByStatus = append(mFindOneByStatus, string(item.ValueString()))
		}
	}

	params = waap.DomainAPIPathListParams{
		DomainID: m.DomainID.ValueInt64(),
		IDs:      mFindOneByIDs,
		Status:   mFindOneByStatus,
	}

	if !m.FindOneBy.APIGroup.IsNull() {
		params.APIGroup = param.NewOpt(m.FindOneBy.APIGroup.ValueString())
	}
	if !m.FindOneBy.APIVersion.IsNull() {
		params.APIVersion = param.NewOpt(m.FindOneBy.APIVersion.ValueString())
	}
	if !m.FindOneBy.HTTPScheme.IsNull() {
		params.HTTPScheme = waap.DomainAPIPathListParamsHTTPScheme(m.FindOneBy.HTTPScheme.ValueString())
	}
	if !m.FindOneBy.Method.IsNull() {
		params.Method = waap.DomainAPIPathListParamsMethod(m.FindOneBy.Method.ValueString())
	}
	if !m.FindOneBy.Ordering.IsNull() {
		params.Ordering = waap.DomainAPIPathListParamsOrdering(m.FindOneBy.Ordering.ValueString())
	}
	if !m.FindOneBy.Path.IsNull() {
		params.Path = param.NewOpt(m.FindOneBy.Path.ValueString())
	}
	if !m.FindOneBy.Source.IsNull() {
		params.Source = waap.DomainAPIPathListParamsSource(m.FindOneBy.Source.ValueString())
	}

	return
}

type WaapDomainAPIPathFindOneByDataSourceModel struct {
	APIGroup   types.String    `tfsdk:"api_group" query:"api_group,optional"`
	APIVersion types.String    `tfsdk:"api_version" query:"api_version,optional"`
	HTTPScheme types.String    `tfsdk:"http_scheme" query:"http_scheme,optional"`
	IDs        *[]types.String `tfsdk:"ids" query:"ids,optional"`
	Method     types.String    `tfsdk:"method" query:"method,optional"`
	Ordering   types.String    `tfsdk:"ordering" query:"ordering,optional"`
	Path       types.String    `tfsdk:"path" query:"path,optional"`
	Source     types.String    `tfsdk:"source" query:"source,optional"`
	Status     *[]types.String `tfsdk:"status" query:"status,optional"`
}

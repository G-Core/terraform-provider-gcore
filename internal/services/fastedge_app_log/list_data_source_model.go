// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app_log

import (
	"context"

	"github.com/G-Core/gcore-go/fastedge"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type FastedgeAppLogsLogsListDataSourceEnvelope struct {
	Logs customfield.NestedObjectList[FastedgeAppLogsItemsDataSourceModel] `json:"logs,computed"`
}

type FastedgeAppLogsDataSourceModel struct {
	ID       types.Int64                                                       `tfsdk:"id" path:"id,required"`
	ClientIP types.String                                                      `tfsdk:"client_ip" query:"client_ip,optional"`
	Edge     types.String                                                      `tfsdk:"edge" query:"edge,optional"`
	From     timetypes.RFC3339                                                 `tfsdk:"from" query:"from,optional" format:"date-time"`
	Limit    types.Int64                                                       `tfsdk:"limit" query:"limit,optional"`
	Search   types.String                                                      `tfsdk:"search" query:"search,optional"`
	Sort     types.String                                                      `tfsdk:"sort" query:"sort,optional"`
	To       timetypes.RFC3339                                                 `tfsdk:"to" query:"to,optional" format:"date-time"`
	MaxItems types.Int64                                                       `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[FastedgeAppLogsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *FastedgeAppLogsDataSourceModel) toListParams(_ context.Context) (params fastedge.AppLogListParams, diags diag.Diagnostics) {
	mFrom, errs := m.From.ValueRFC3339Time()
	diags.Append(errs...)
	mTo, errs := m.To.ValueRFC3339Time()
	diags.Append(errs...)

	params = fastedge.AppLogListParams{}

	if !m.ClientIP.IsNull() {
		params.ClientIP = param.NewOpt(m.ClientIP.ValueString())
	}
	if !m.Edge.IsNull() {
		params.Edge = param.NewOpt(m.Edge.ValueString())
	}
	if !m.From.IsNull() {
		params.From = param.NewOpt(mFrom)
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Search.IsNull() {
		params.Search = param.NewOpt(m.Search.ValueString())
	}
	if !m.Sort.IsNull() {
		params.Sort = fastedge.AppLogListParamsSort(m.Sort.ValueString())
	}
	if !m.To.IsNull() {
		params.To = param.NewOpt(mTo)
	}

	return
}

type FastedgeAppLogsItemsDataSourceModel struct {
	ID        types.String      `tfsdk:"id" json:"id,computed"`
	AppName   types.String      `tfsdk:"app_name" json:"app_name,computed"`
	ClientIP  types.String      `tfsdk:"client_ip" json:"client_ip,computed"`
	Edge      types.String      `tfsdk:"edge" json:"edge,computed"`
	Log       types.String      `tfsdk:"log" json:"log,computed"`
	Timestamp timetypes.RFC3339 `tfsdk:"timestamp" json:"timestamp,computed" format:"date-time"`
}

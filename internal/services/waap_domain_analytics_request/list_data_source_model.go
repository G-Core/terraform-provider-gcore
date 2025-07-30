// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_analytics_request

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainAnalyticsRequestsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapDomainAnalyticsRequestsItemsDataSourceModel] `json:"results,computed"`
}

type WaapDomainAnalyticsRequestsDataSourceModel struct {
	DomainID         types.Int64                                                                   `tfsdk:"domain_id" path:"domain_id,required"`
	Start            timetypes.RFC3339                                                             `tfsdk:"start" query:"start,required" format:"date-time"`
	End              timetypes.RFC3339                                                             `tfsdk:"end" query:"end,optional" format:"date-time"`
	IP               types.String                                                                  `tfsdk:"ip" query:"ip,optional"`
	Ordering         types.String                                                                  `tfsdk:"ordering" query:"ordering,optional"`
	ReferenceID      types.String                                                                  `tfsdk:"reference_id" query:"reference_id,optional"`
	SecurityRuleName types.String                                                                  `tfsdk:"security_rule_name" query:"security_rule_name,optional"`
	StatusCode       types.Int64                                                                   `tfsdk:"status_code" query:"status_code,optional"`
	Actions          *[]types.String                                                               `tfsdk:"actions" query:"actions,optional"`
	Countries        *[]types.String                                                               `tfsdk:"countries" query:"countries,optional"`
	TrafficTypes     *[]types.String                                                               `tfsdk:"traffic_types" query:"traffic_types,optional"`
	Limit            types.Int64                                                                   `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems         types.Int64                                                                   `tfsdk:"max_items"`
	Items            customfield.NestedObjectList[WaapDomainAnalyticsRequestsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapDomainAnalyticsRequestsDataSourceModel) toListParams(_ context.Context) (params waap.DomainAnalyticsRequestListParams, diags diag.Diagnostics) {
	mStart, errs := m.Start.ValueRFC3339Time()
	diags.Append(errs...)
	mActions := []string{}
	for _, item := range *m.Actions {
		mActions = append(mActions, string(item.ValueString()))
	}
	mCountries := []string{}
	for _, item := range *m.Countries {
		mCountries = append(mCountries, item.ValueString())
	}
	mTrafficTypes := []waap.WaapTrafficType{}
	for _, item := range *m.TrafficTypes {
		mTrafficTypes = append(mTrafficTypes, waap.WaapTrafficType(item.ValueString()))
	}
	mEnd, errs := m.End.ValueRFC3339Time()
	diags.Append(errs...)

	params = waap.DomainAnalyticsRequestListParams{
		Start:        mStart,
		Actions:      mActions,
		Countries:    mCountries,
		TrafficTypes: mTrafficTypes,
	}

	if !m.End.IsNull() {
		params.End = param.NewOpt(mEnd)
	}
	if !m.IP.IsNull() {
		params.IP = param.NewOpt(m.IP.ValueString())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = param.NewOpt(m.Ordering.ValueString())
	}
	if !m.ReferenceID.IsNull() {
		params.ReferenceID = param.NewOpt(m.ReferenceID.ValueString())
	}
	if !m.SecurityRuleName.IsNull() {
		params.SecurityRuleName = param.NewOpt(m.SecurityRuleName.ValueString())
	}
	if !m.StatusCode.IsNull() {
		params.StatusCode = param.NewOpt(m.StatusCode.ValueInt64())
	}

	return
}

type WaapDomainAnalyticsRequestsItemsDataSourceModel struct {
	ID              types.String `tfsdk:"id" json:"id,computed"`
	Action          types.String `tfsdk:"action" json:"action,computed"`
	ClientIP        types.String `tfsdk:"client_ip" json:"client_ip,computed"`
	Country         types.String `tfsdk:"country" json:"country,computed"`
	Domain          types.String `tfsdk:"domain" json:"domain,computed"`
	Method          types.String `tfsdk:"method" json:"method,computed"`
	Organization    types.String `tfsdk:"organization" json:"organization,computed"`
	Path            types.String `tfsdk:"path" json:"path,computed"`
	ReferenceID     types.String `tfsdk:"reference_id" json:"reference_id,computed"`
	RequestTime     types.Int64  `tfsdk:"request_time" json:"request_time,computed"`
	Result          types.String `tfsdk:"result" json:"result,computed"`
	RuleID          types.String `tfsdk:"rule_id" json:"rule_id,computed"`
	RuleName        types.String `tfsdk:"rule_name" json:"rule_name,computed"`
	StatusCode      types.Int64  `tfsdk:"status_code" json:"status_code,computed"`
	TrafficTypes    types.String `tfsdk:"traffic_types" json:"traffic_types,computed"`
	UserAgent       types.String `tfsdk:"user_agent" json:"user_agent,computed"`
	UserAgentClient types.String `tfsdk:"user_agent_client" json:"user_agent_client,computed"`
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_analytics_request

import (
	"context"

	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainAnalyticsRequestDataSourceModel struct {
	DomainID            types.Int64                                                                               `tfsdk:"domain_id" path:"domain_id,required"`
	RequestID           types.String                                                                              `tfsdk:"request_id" path:"request_id,required"`
	Action              types.String                                                                              `tfsdk:"action" json:"action,computed"`
	ContentType         types.String                                                                              `tfsdk:"content_type" json:"content_type,computed"`
	Domain              types.String                                                                              `tfsdk:"domain" json:"domain,computed"`
	HTTPStatusCode      types.Int64                                                                               `tfsdk:"http_status_code" json:"http_status_code,computed"`
	HTTPVersion         types.String                                                                              `tfsdk:"http_version" json:"http_version,computed"`
	ID                  types.String                                                                              `tfsdk:"id" json:"id,computed"`
	IncidentID          types.String                                                                              `tfsdk:"incident_id" json:"incident_id,computed"`
	Method              types.String                                                                              `tfsdk:"method" json:"method,computed"`
	Path                types.String                                                                              `tfsdk:"path" json:"path,computed"`
	QueryString         types.String                                                                              `tfsdk:"query_string" json:"query_string,computed"`
	ReferenceID         types.String                                                                              `tfsdk:"reference_id" json:"reference_id,computed"`
	RequestTime         types.String                                                                              `tfsdk:"request_time" json:"request_time,computed"`
	RequestType         types.String                                                                              `tfsdk:"request_type" json:"request_type,computed"`
	RequestedDomain     types.String                                                                              `tfsdk:"requested_domain" json:"requested_domain,computed"`
	ResponseTime        types.String                                                                              `tfsdk:"response_time" json:"response_time,computed"`
	Result              types.String                                                                              `tfsdk:"result" json:"result,computed"`
	RuleID              types.String                                                                              `tfsdk:"rule_id" json:"rule_id,computed"`
	RuleName            types.String                                                                              `tfsdk:"rule_name" json:"rule_name,computed"`
	Scheme              types.String                                                                              `tfsdk:"scheme" json:"scheme,computed"`
	SessionRequestCount types.String                                                                              `tfsdk:"session_request_count" json:"session_request_count,computed"`
	TrafficTypes        customfield.List[types.String]                                                            `tfsdk:"traffic_types" json:"traffic_types,computed"`
	CommonTags          customfield.NestedObjectList[WaapDomainAnalyticsRequestCommonTagsDataSourceModel]         `tfsdk:"common_tags" json:"common_tags,computed"`
	Network             customfield.NestedObject[WaapDomainAnalyticsRequestNetworkDataSourceModel]                `tfsdk:"network" json:"network,computed"`
	PatternMatchedTags  customfield.NestedObjectList[WaapDomainAnalyticsRequestPatternMatchedTagsDataSourceModel] `tfsdk:"pattern_matched_tags" json:"pattern_matched_tags,computed"`
	UserAgent           customfield.NestedObject[WaapDomainAnalyticsRequestUserAgentDataSourceModel]              `tfsdk:"user_agent" json:"user_agent,computed"`
	RequestHeaders      jsontypes.Normalized                                                                      `tfsdk:"request_headers" json:"request_headers,computed"`
}

func (m *WaapDomainAnalyticsRequestDataSourceModel) toReadParams(_ context.Context) (params waap.DomainAnalyticsRequestGetParams, diags diag.Diagnostics) {
	params = waap.DomainAnalyticsRequestGetParams{
		DomainID: m.DomainID.ValueInt64(),
	}

	return
}

type WaapDomainAnalyticsRequestCommonTagsDataSourceModel struct {
	Description types.String `tfsdk:"description" json:"description,computed"`
	DisplayName types.String `tfsdk:"display_name" json:"display_name,computed"`
	Tag         types.String `tfsdk:"tag" json:"tag,computed"`
}

type WaapDomainAnalyticsRequestNetworkDataSourceModel struct {
	ClientIP     types.String                                                                           `tfsdk:"client_ip" json:"client_ip,computed"`
	Country      types.String                                                                           `tfsdk:"country" json:"country,computed"`
	Organization customfield.NestedObject[WaapDomainAnalyticsRequestNetworkOrganizationDataSourceModel] `tfsdk:"organization" json:"organization,computed"`
}

type WaapDomainAnalyticsRequestNetworkOrganizationDataSourceModel struct {
	Name   types.String `tfsdk:"name" json:"name,computed"`
	Subnet types.String `tfsdk:"subnet" json:"subnet,computed"`
}

type WaapDomainAnalyticsRequestPatternMatchedTagsDataSourceModel struct {
	Description    types.String `tfsdk:"description" json:"description,computed"`
	DisplayName    types.String `tfsdk:"display_name" json:"display_name,computed"`
	ExecutionPhase types.String `tfsdk:"execution_phase" json:"execution_phase,computed"`
	Field          types.String `tfsdk:"field" json:"field,computed"`
	FieldName      types.String `tfsdk:"field_name" json:"field_name,computed"`
	PatternName    types.String `tfsdk:"pattern_name" json:"pattern_name,computed"`
	PatternValue   types.String `tfsdk:"pattern_value" json:"pattern_value,computed"`
	Tag            types.String `tfsdk:"tag" json:"tag,computed"`
}

type WaapDomainAnalyticsRequestUserAgentDataSourceModel struct {
	BaseBrowser        types.String `tfsdk:"base_browser" json:"base_browser,computed"`
	BaseBrowserVersion types.String `tfsdk:"base_browser_version" json:"base_browser_version,computed"`
	Client             types.String `tfsdk:"client" json:"client,computed"`
	ClientType         types.String `tfsdk:"client_type" json:"client_type,computed"`
	ClientVersion      types.String `tfsdk:"client_version" json:"client_version,computed"`
	CPU                types.String `tfsdk:"cpu" json:"cpu,computed"`
	Device             types.String `tfsdk:"device" json:"device,computed"`
	DeviceType         types.String `tfsdk:"device_type" json:"device_type,computed"`
	FullString         types.String `tfsdk:"full_string" json:"full_string,computed"`
	Os                 types.String `tfsdk:"os" json:"os,computed"`
	RenderingEngine    types.String `tfsdk:"rendering_engine" json:"rendering_engine,computed"`
}

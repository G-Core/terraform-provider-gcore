// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_custom_rule

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainCustomRulesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapDomainCustomRulesItemsDataSourceModel] `json:"results,computed"`
}

type WaapDomainCustomRulesDataSourceModel struct {
	DomainID    types.Int64                                                             `tfsdk:"domain_id" path:"domain_id,required"`
	Action      types.String                                                            `tfsdk:"action" query:"action,optional"`
	Description types.String                                                            `tfsdk:"description" query:"description,optional"`
	Enabled     types.Bool                                                              `tfsdk:"enabled" query:"enabled,optional"`
	Name        types.String                                                            `tfsdk:"name" query:"name,optional"`
	Ordering    types.String                                                            `tfsdk:"ordering" query:"ordering,optional"`
	MaxItems    types.Int64                                                             `tfsdk:"max_items"`
	Items       customfield.NestedObjectList[WaapDomainCustomRulesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapDomainCustomRulesDataSourceModel) toListParams(_ context.Context) (params waap.DomainCustomRuleListParams, diags diag.Diagnostics) {
	params = waap.DomainCustomRuleListParams{}

	if !m.Action.IsNull() {
		params.Action = waap.DomainCustomRuleListParamsAction(m.Action.ValueString())
	}
	if !m.Description.IsNull() {
		params.Description = param.NewOpt(m.Description.ValueString())
	}
	if !m.Enabled.IsNull() {
		params.Enabled = param.NewOpt(m.Enabled.ValueBool())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = waap.DomainCustomRuleListParamsOrdering(m.Ordering.ValueString())
	}

	return
}

type WaapDomainCustomRulesItemsDataSourceModel struct {
	ID          types.Int64                                                                  `tfsdk:"id" json:"id,computed"`
	Action      customfield.NestedObject[WaapDomainCustomRulesActionDataSourceModel]         `tfsdk:"action" json:"action,computed"`
	Conditions  customfield.NestedObjectList[WaapDomainCustomRulesConditionsDataSourceModel] `tfsdk:"conditions" json:"conditions,computed"`
	Enabled     types.Bool                                                                   `tfsdk:"enabled" json:"enabled,computed"`
	Name        types.String                                                                 `tfsdk:"name" json:"name,computed"`
	Description types.String                                                                 `tfsdk:"description" json:"description,computed"`
}

type WaapDomainCustomRulesActionDataSourceModel struct {
	Allow     jsontypes.Normalized                                                      `tfsdk:"allow" json:"allow,computed"`
	Block     customfield.NestedObject[WaapDomainCustomRulesActionBlockDataSourceModel] `tfsdk:"block" json:"block,computed"`
	Captcha   jsontypes.Normalized                                                      `tfsdk:"captcha" json:"captcha,computed"`
	Handshake jsontypes.Normalized                                                      `tfsdk:"handshake" json:"handshake,computed"`
	Monitor   jsontypes.Normalized                                                      `tfsdk:"monitor" json:"monitor,computed"`
	Tag       customfield.NestedObject[WaapDomainCustomRulesActionTagDataSourceModel]   `tfsdk:"tag" json:"tag,computed"`
}

type WaapDomainCustomRulesActionBlockDataSourceModel struct {
	ActionDuration types.String `tfsdk:"action_duration" json:"action_duration,computed"`
	StatusCode     types.Int64  `tfsdk:"status_code" json:"status_code,computed"`
}

type WaapDomainCustomRulesActionTagDataSourceModel struct {
	Tags customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
}

type WaapDomainCustomRulesConditionsDataSourceModel struct {
	ContentType          customfield.NestedObject[WaapDomainCustomRulesConditionsContentTypeDataSourceModel]          `tfsdk:"content_type" json:"content_type,computed"`
	Country              customfield.NestedObject[WaapDomainCustomRulesConditionsCountryDataSourceModel]              `tfsdk:"country" json:"country,computed"`
	FileExtension        customfield.NestedObject[WaapDomainCustomRulesConditionsFileExtensionDataSourceModel]        `tfsdk:"file_extension" json:"file_extension,computed"`
	Header               customfield.NestedObject[WaapDomainCustomRulesConditionsHeaderDataSourceModel]               `tfsdk:"header" json:"header,computed"`
	HeaderExists         customfield.NestedObject[WaapDomainCustomRulesConditionsHeaderExistsDataSourceModel]         `tfsdk:"header_exists" json:"header_exists,computed"`
	HTTPMethod           customfield.NestedObject[WaapDomainCustomRulesConditionsHTTPMethodDataSourceModel]           `tfsdk:"http_method" json:"http_method,computed"`
	IP                   customfield.NestedObject[WaapDomainCustomRulesConditionsIPDataSourceModel]                   `tfsdk:"ip" json:"ip,computed"`
	IPRange              customfield.NestedObject[WaapDomainCustomRulesConditionsIPRangeDataSourceModel]              `tfsdk:"ip_range" json:"ip_range,computed"`
	Organization         customfield.NestedObject[WaapDomainCustomRulesConditionsOrganizationDataSourceModel]         `tfsdk:"organization" json:"organization,computed"`
	OwnerTypes           customfield.NestedObject[WaapDomainCustomRulesConditionsOwnerTypesDataSourceModel]           `tfsdk:"owner_types" json:"owner_types,computed"`
	RequestRate          customfield.NestedObject[WaapDomainCustomRulesConditionsRequestRateDataSourceModel]          `tfsdk:"request_rate" json:"request_rate,computed"`
	ResponseHeader       customfield.NestedObject[WaapDomainCustomRulesConditionsResponseHeaderDataSourceModel]       `tfsdk:"response_header" json:"response_header,computed"`
	ResponseHeaderExists customfield.NestedObject[WaapDomainCustomRulesConditionsResponseHeaderExistsDataSourceModel] `tfsdk:"response_header_exists" json:"response_header_exists,computed"`
	SessionRequestCount  customfield.NestedObject[WaapDomainCustomRulesConditionsSessionRequestCountDataSourceModel]  `tfsdk:"session_request_count" json:"session_request_count,computed"`
	Tags                 customfield.NestedObject[WaapDomainCustomRulesConditionsTagsDataSourceModel]                 `tfsdk:"tags" json:"tags,computed"`
	URL                  customfield.NestedObject[WaapDomainCustomRulesConditionsURLDataSourceModel]                  `tfsdk:"url" json:"url,computed"`
	UserAgent            customfield.NestedObject[WaapDomainCustomRulesConditionsUserAgentDataSourceModel]            `tfsdk:"user_agent" json:"user_agent,computed"`
	UserDefinedTags      customfield.NestedObject[WaapDomainCustomRulesConditionsUserDefinedTagsDataSourceModel]      `tfsdk:"user_defined_tags" json:"user_defined_tags,computed"`
}

type WaapDomainCustomRulesConditionsContentTypeDataSourceModel struct {
	ContentType customfield.List[types.String] `tfsdk:"content_type" json:"content_type,computed"`
	Negation    types.Bool                     `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsCountryDataSourceModel struct {
	CountryCode customfield.List[types.String] `tfsdk:"country_code" json:"country_code,computed"`
	Negation    types.Bool                     `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsFileExtensionDataSourceModel struct {
	FileExtension customfield.List[types.String] `tfsdk:"file_extension" json:"file_extension,computed"`
	Negation      types.Bool                     `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsHeaderDataSourceModel struct {
	Header    types.String `tfsdk:"header" json:"header,computed"`
	Value     types.String `tfsdk:"value" json:"value,computed"`
	MatchType types.String `tfsdk:"match_type" json:"match_type,computed"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsHeaderExistsDataSourceModel struct {
	Header   types.String `tfsdk:"header" json:"header,computed"`
	Negation types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsHTTPMethodDataSourceModel struct {
	HTTPMethod types.String `tfsdk:"http_method" json:"http_method,computed"`
	Negation   types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsIPDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsIPRangeDataSourceModel struct {
	LowerBound types.String `tfsdk:"lower_bound" json:"lower_bound,computed"`
	UpperBound types.String `tfsdk:"upper_bound" json:"upper_bound,computed"`
	Negation   types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsOrganizationDataSourceModel struct {
	Organization types.String `tfsdk:"organization" json:"organization,computed"`
	Negation     types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsOwnerTypesDataSourceModel struct {
	Negation   types.Bool                     `tfsdk:"negation" json:"negation,computed"`
	OwnerTypes customfield.List[types.String] `tfsdk:"owner_types" json:"owner_types,computed"`
}

type WaapDomainCustomRulesConditionsRequestRateDataSourceModel struct {
	PathPattern    types.String                   `tfsdk:"path_pattern" json:"path_pattern,computed"`
	Requests       types.Int64                    `tfsdk:"requests" json:"requests,computed"`
	Time           types.Int64                    `tfsdk:"time" json:"time,computed"`
	HTTPMethods    customfield.List[types.String] `tfsdk:"http_methods" json:"http_methods,computed"`
	IPs            customfield.List[types.String] `tfsdk:"ips" json:"ips,computed"`
	UserDefinedTag types.String                   `tfsdk:"user_defined_tag" json:"user_defined_tag,computed"`
}

type WaapDomainCustomRulesConditionsResponseHeaderDataSourceModel struct {
	Header    types.String `tfsdk:"header" json:"header,computed"`
	Value     types.String `tfsdk:"value" json:"value,computed"`
	MatchType types.String `tfsdk:"match_type" json:"match_type,computed"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsResponseHeaderExistsDataSourceModel struct {
	Header   types.String `tfsdk:"header" json:"header,computed"`
	Negation types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsSessionRequestCountDataSourceModel struct {
	RequestCount types.Int64 `tfsdk:"request_count" json:"request_count,computed"`
	Negation     types.Bool  `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsTagsDataSourceModel struct {
	Tags     customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
	Negation types.Bool                     `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsURLDataSourceModel struct {
	URL       types.String `tfsdk:"url" json:"url,computed"`
	MatchType types.String `tfsdk:"match_type" json:"match_type,computed"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsUserAgentDataSourceModel struct {
	UserAgent types.String `tfsdk:"user_agent" json:"user_agent,computed"`
	MatchType types.String `tfsdk:"match_type" json:"match_type,computed"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRulesConditionsUserDefinedTagsDataSourceModel struct {
	Tags     customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
	Negation types.Bool                     `tfsdk:"negation" json:"negation,computed"`
}

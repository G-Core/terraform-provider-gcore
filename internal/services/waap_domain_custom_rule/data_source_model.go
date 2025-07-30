// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_custom_rule

import (
	"context"

	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainCustomRuleDataSourceModel struct {
	DomainID    types.Int64                                                                 `tfsdk:"domain_id" path:"domain_id,required"`
	RuleID      types.Int64                                                                 `tfsdk:"rule_id" path:"rule_id,required"`
	Description types.String                                                                `tfsdk:"description" json:"description,computed"`
	Enabled     types.Bool                                                                  `tfsdk:"enabled" json:"enabled,computed"`
	ID          types.Int64                                                                 `tfsdk:"id" json:"id,computed"`
	Name        types.String                                                                `tfsdk:"name" json:"name,computed"`
	Action      customfield.NestedObject[WaapDomainCustomRuleActionDataSourceModel]         `tfsdk:"action" json:"action,computed"`
	Conditions  customfield.NestedObjectList[WaapDomainCustomRuleConditionsDataSourceModel] `tfsdk:"conditions" json:"conditions,computed"`
}

func (m *WaapDomainCustomRuleDataSourceModel) toReadParams(_ context.Context) (params waap.DomainCustomRuleGetParams, diags diag.Diagnostics) {
	params = waap.DomainCustomRuleGetParams{
		DomainID: m.DomainID.ValueInt64(),
	}

	return
}

type WaapDomainCustomRuleActionDataSourceModel struct {
	Allow     jsontypes.Normalized                                                     `tfsdk:"allow" json:"allow,computed"`
	Block     customfield.NestedObject[WaapDomainCustomRuleActionBlockDataSourceModel] `tfsdk:"block" json:"block,computed"`
	Captcha   jsontypes.Normalized                                                     `tfsdk:"captcha" json:"captcha,computed"`
	Handshake jsontypes.Normalized                                                     `tfsdk:"handshake" json:"handshake,computed"`
	Monitor   jsontypes.Normalized                                                     `tfsdk:"monitor" json:"monitor,computed"`
	Tag       customfield.NestedObject[WaapDomainCustomRuleActionTagDataSourceModel]   `tfsdk:"tag" json:"tag,computed"`
}

type WaapDomainCustomRuleActionBlockDataSourceModel struct {
	ActionDuration types.String `tfsdk:"action_duration" json:"action_duration,computed"`
	StatusCode     types.Int64  `tfsdk:"status_code" json:"status_code,computed"`
}

type WaapDomainCustomRuleActionTagDataSourceModel struct {
	Tags customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
}

type WaapDomainCustomRuleConditionsDataSourceModel struct {
	ContentType          customfield.NestedObject[WaapDomainCustomRuleConditionsContentTypeDataSourceModel]          `tfsdk:"content_type" json:"content_type,computed"`
	Country              customfield.NestedObject[WaapDomainCustomRuleConditionsCountryDataSourceModel]              `tfsdk:"country" json:"country,computed"`
	FileExtension        customfield.NestedObject[WaapDomainCustomRuleConditionsFileExtensionDataSourceModel]        `tfsdk:"file_extension" json:"file_extension,computed"`
	Header               customfield.NestedObject[WaapDomainCustomRuleConditionsHeaderDataSourceModel]               `tfsdk:"header" json:"header,computed"`
	HeaderExists         customfield.NestedObject[WaapDomainCustomRuleConditionsHeaderExistsDataSourceModel]         `tfsdk:"header_exists" json:"header_exists,computed"`
	HTTPMethod           customfield.NestedObject[WaapDomainCustomRuleConditionsHTTPMethodDataSourceModel]           `tfsdk:"http_method" json:"http_method,computed"`
	IP                   customfield.NestedObject[WaapDomainCustomRuleConditionsIPDataSourceModel]                   `tfsdk:"ip" json:"ip,computed"`
	IPRange              customfield.NestedObject[WaapDomainCustomRuleConditionsIPRangeDataSourceModel]              `tfsdk:"ip_range" json:"ip_range,computed"`
	Organization         customfield.NestedObject[WaapDomainCustomRuleConditionsOrganizationDataSourceModel]         `tfsdk:"organization" json:"organization,computed"`
	OwnerTypes           customfield.NestedObject[WaapDomainCustomRuleConditionsOwnerTypesDataSourceModel]           `tfsdk:"owner_types" json:"owner_types,computed"`
	RequestRate          customfield.NestedObject[WaapDomainCustomRuleConditionsRequestRateDataSourceModel]          `tfsdk:"request_rate" json:"request_rate,computed"`
	ResponseHeader       customfield.NestedObject[WaapDomainCustomRuleConditionsResponseHeaderDataSourceModel]       `tfsdk:"response_header" json:"response_header,computed"`
	ResponseHeaderExists customfield.NestedObject[WaapDomainCustomRuleConditionsResponseHeaderExistsDataSourceModel] `tfsdk:"response_header_exists" json:"response_header_exists,computed"`
	SessionRequestCount  customfield.NestedObject[WaapDomainCustomRuleConditionsSessionRequestCountDataSourceModel]  `tfsdk:"session_request_count" json:"session_request_count,computed"`
	Tags                 customfield.NestedObject[WaapDomainCustomRuleConditionsTagsDataSourceModel]                 `tfsdk:"tags" json:"tags,computed"`
	URL                  customfield.NestedObject[WaapDomainCustomRuleConditionsURLDataSourceModel]                  `tfsdk:"url" json:"url,computed"`
	UserAgent            customfield.NestedObject[WaapDomainCustomRuleConditionsUserAgentDataSourceModel]            `tfsdk:"user_agent" json:"user_agent,computed"`
	UserDefinedTags      customfield.NestedObject[WaapDomainCustomRuleConditionsUserDefinedTagsDataSourceModel]      `tfsdk:"user_defined_tags" json:"user_defined_tags,computed"`
}

type WaapDomainCustomRuleConditionsContentTypeDataSourceModel struct {
	ContentType customfield.List[types.String] `tfsdk:"content_type" json:"content_type,computed"`
	Negation    types.Bool                     `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsCountryDataSourceModel struct {
	CountryCode customfield.List[types.String] `tfsdk:"country_code" json:"country_code,computed"`
	Negation    types.Bool                     `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsFileExtensionDataSourceModel struct {
	FileExtension customfield.List[types.String] `tfsdk:"file_extension" json:"file_extension,computed"`
	Negation      types.Bool                     `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsHeaderDataSourceModel struct {
	Header    types.String `tfsdk:"header" json:"header,computed"`
	Value     types.String `tfsdk:"value" json:"value,computed"`
	MatchType types.String `tfsdk:"match_type" json:"match_type,computed"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsHeaderExistsDataSourceModel struct {
	Header   types.String `tfsdk:"header" json:"header,computed"`
	Negation types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsHTTPMethodDataSourceModel struct {
	HTTPMethod types.String `tfsdk:"http_method" json:"http_method,computed"`
	Negation   types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsIPDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsIPRangeDataSourceModel struct {
	LowerBound types.String `tfsdk:"lower_bound" json:"lower_bound,computed"`
	UpperBound types.String `tfsdk:"upper_bound" json:"upper_bound,computed"`
	Negation   types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsOrganizationDataSourceModel struct {
	Organization types.String `tfsdk:"organization" json:"organization,computed"`
	Negation     types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsOwnerTypesDataSourceModel struct {
	Negation   types.Bool                     `tfsdk:"negation" json:"negation,computed"`
	OwnerTypes customfield.List[types.String] `tfsdk:"owner_types" json:"owner_types,computed"`
}

type WaapDomainCustomRuleConditionsRequestRateDataSourceModel struct {
	PathPattern    types.String                   `tfsdk:"path_pattern" json:"path_pattern,computed"`
	Requests       types.Int64                    `tfsdk:"requests" json:"requests,computed"`
	Time           types.Int64                    `tfsdk:"time" json:"time,computed"`
	HTTPMethods    customfield.List[types.String] `tfsdk:"http_methods" json:"http_methods,computed"`
	IPs            customfield.List[types.String] `tfsdk:"ips" json:"ips,computed"`
	UserDefinedTag types.String                   `tfsdk:"user_defined_tag" json:"user_defined_tag,computed"`
}

type WaapDomainCustomRuleConditionsResponseHeaderDataSourceModel struct {
	Header    types.String `tfsdk:"header" json:"header,computed"`
	Value     types.String `tfsdk:"value" json:"value,computed"`
	MatchType types.String `tfsdk:"match_type" json:"match_type,computed"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsResponseHeaderExistsDataSourceModel struct {
	Header   types.String `tfsdk:"header" json:"header,computed"`
	Negation types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsSessionRequestCountDataSourceModel struct {
	RequestCount types.Int64 `tfsdk:"request_count" json:"request_count,computed"`
	Negation     types.Bool  `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsTagsDataSourceModel struct {
	Tags     customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
	Negation types.Bool                     `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsURLDataSourceModel struct {
	URL       types.String `tfsdk:"url" json:"url,computed"`
	MatchType types.String `tfsdk:"match_type" json:"match_type,computed"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsUserAgentDataSourceModel struct {
	UserAgent types.String `tfsdk:"user_agent" json:"user_agent,computed"`
	MatchType types.String `tfsdk:"match_type" json:"match_type,computed"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainCustomRuleConditionsUserDefinedTagsDataSourceModel struct {
	Tags     customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
	Negation types.Bool                     `tfsdk:"negation" json:"negation,computed"`
}

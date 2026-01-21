// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_custom_rule

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainCustomRuleModel struct {
	ID          types.Int64                             `tfsdk:"id" json:"id,computed"`
	DomainID    types.Int64                             `tfsdk:"domain_id" path:"domain_id,required"`
	Enabled     types.Bool                              `tfsdk:"enabled" json:"enabled,required"`
	Name        types.String                            `tfsdk:"name" json:"name,required"`
	Action      *WaapDomainCustomRuleActionModel        `tfsdk:"action" json:"action,required"`
	Conditions  *[]*WaapDomainCustomRuleConditionsModel `tfsdk:"conditions" json:"conditions,required"`
	Description types.String                            `tfsdk:"description" json:"description,optional"`
}

func (m WaapDomainCustomRuleModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m WaapDomainCustomRuleModel) MarshalJSONForUpdate(state WaapDomainCustomRuleModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type WaapDomainCustomRuleActionModel struct {
	Allow     jsontypes.Normalized                  `tfsdk:"allow" json:"allow,optional"`
	Block     *WaapDomainCustomRuleActionBlockModel `tfsdk:"block" json:"block,optional"`
	Captcha   jsontypes.Normalized                  `tfsdk:"captcha" json:"captcha,optional"`
	Handshake jsontypes.Normalized                  `tfsdk:"handshake" json:"handshake,optional"`
	Monitor   jsontypes.Normalized                  `tfsdk:"monitor" json:"monitor,optional"`
	Tag       *WaapDomainCustomRuleActionTagModel   `tfsdk:"tag" json:"tag,optional"`
}

type WaapDomainCustomRuleActionBlockModel struct {
	ActionDuration types.String `tfsdk:"action_duration" json:"action_duration,optional"`
	StatusCode     types.Int64  `tfsdk:"status_code" json:"status_code,optional"`
}

type WaapDomainCustomRuleActionTagModel struct {
	Tags *[]types.String `tfsdk:"tags" json:"tags,required"`
}

type WaapDomainCustomRuleConditionsModel struct {
	ContentType          *WaapDomainCustomRuleConditionsContentTypeModel          `tfsdk:"content_type" json:"content_type,optional"`
	Country              *WaapDomainCustomRuleConditionsCountryModel              `tfsdk:"country" json:"country,optional"`
	FileExtension        *WaapDomainCustomRuleConditionsFileExtensionModel        `tfsdk:"file_extension" json:"file_extension,optional"`
	Header               *WaapDomainCustomRuleConditionsHeaderModel               `tfsdk:"header" json:"header,optional"`
	HeaderExists         *WaapDomainCustomRuleConditionsHeaderExistsModel         `tfsdk:"header_exists" json:"header_exists,optional"`
	HTTPMethod           *WaapDomainCustomRuleConditionsHTTPMethodModel           `tfsdk:"http_method" json:"http_method,optional"`
	IP                   *WaapDomainCustomRuleConditionsIPModel                   `tfsdk:"ip" json:"ip,optional"`
	IPRange              *WaapDomainCustomRuleConditionsIPRangeModel              `tfsdk:"ip_range" json:"ip_range,optional"`
	Organization         *WaapDomainCustomRuleConditionsOrganizationModel         `tfsdk:"organization" json:"organization,optional"`
	OwnerTypes           *WaapDomainCustomRuleConditionsOwnerTypesModel           `tfsdk:"owner_types" json:"owner_types,optional"`
	RequestRate          *WaapDomainCustomRuleConditionsRequestRateModel          `tfsdk:"request_rate" json:"request_rate,optional"`
	ResponseHeader       *WaapDomainCustomRuleConditionsResponseHeaderModel       `tfsdk:"response_header" json:"response_header,optional"`
	ResponseHeaderExists *WaapDomainCustomRuleConditionsResponseHeaderExistsModel `tfsdk:"response_header_exists" json:"response_header_exists,optional"`
	SessionRequestCount  *WaapDomainCustomRuleConditionsSessionRequestCountModel  `tfsdk:"session_request_count" json:"session_request_count,optional"`
	Tags                 *WaapDomainCustomRuleConditionsTagsModel                 `tfsdk:"tags" json:"tags,optional"`
	URL                  *WaapDomainCustomRuleConditionsURLModel                  `tfsdk:"url" json:"url,optional"`
	UserAgent            *WaapDomainCustomRuleConditionsUserAgentModel            `tfsdk:"user_agent" json:"user_agent,optional"`
	UserDefinedTags      *WaapDomainCustomRuleConditionsUserDefinedTagsModel      `tfsdk:"user_defined_tags" json:"user_defined_tags,optional"`
}

type WaapDomainCustomRuleConditionsContentTypeModel struct {
	ContentType *[]types.String `tfsdk:"content_type" json:"content_type,required"`
	Negation    types.Bool      `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsCountryModel struct {
	CountryCode *[]types.String `tfsdk:"country_code" json:"country_code,required"`
	Negation    types.Bool      `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsFileExtensionModel struct {
	FileExtension *[]types.String `tfsdk:"file_extension" json:"file_extension,required"`
	Negation      types.Bool      `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsHeaderModel struct {
	Header    types.String `tfsdk:"header" json:"header,required"`
	Value     types.String `tfsdk:"value" json:"value,required"`
	MatchType types.String `tfsdk:"match_type" json:"match_type,computed_optional"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsHeaderExistsModel struct {
	Header   types.String `tfsdk:"header" json:"header,required"`
	Negation types.Bool   `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsHTTPMethodModel struct {
	HTTPMethod types.String `tfsdk:"http_method" json:"http_method,required"`
	Negation   types.Bool   `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsIPModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,required"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsIPRangeModel struct {
	LowerBound types.String `tfsdk:"lower_bound" json:"lower_bound,required"`
	UpperBound types.String `tfsdk:"upper_bound" json:"upper_bound,required"`
	Negation   types.Bool   `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsOrganizationModel struct {
	Organization types.String `tfsdk:"organization" json:"organization,required"`
	Negation     types.Bool   `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsOwnerTypesModel struct {
	Negation   types.Bool                     `tfsdk:"negation" json:"negation,computed_optional"`
	OwnerTypes customfield.List[types.String] `tfsdk:"owner_types" json:"owner_types,computed_optional"`
}

type WaapDomainCustomRuleConditionsRequestRateModel struct {
	PathPattern    types.String    `tfsdk:"path_pattern" json:"path_pattern,required"`
	Requests       types.Int64     `tfsdk:"requests" json:"requests,required"`
	Time           types.Int64     `tfsdk:"time" json:"time,required"`
	HTTPMethods    *[]types.String `tfsdk:"http_methods" json:"http_methods,optional"`
	IPs            *[]types.String `tfsdk:"ips" json:"ips,optional"`
	UserDefinedTag types.String    `tfsdk:"user_defined_tag" json:"user_defined_tag,optional"`
}

type WaapDomainCustomRuleConditionsResponseHeaderModel struct {
	Header    types.String `tfsdk:"header" json:"header,required"`
	Value     types.String `tfsdk:"value" json:"value,required"`
	MatchType types.String `tfsdk:"match_type" json:"match_type,computed_optional"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsResponseHeaderExistsModel struct {
	Header   types.String `tfsdk:"header" json:"header,required"`
	Negation types.Bool   `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsSessionRequestCountModel struct {
	RequestCount types.Int64 `tfsdk:"request_count" json:"request_count,required"`
	Negation     types.Bool  `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsTagsModel struct {
	Tags     *[]types.String `tfsdk:"tags" json:"tags,required"`
	Negation types.Bool      `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsURLModel struct {
	URL       types.String `tfsdk:"url" json:"url,required"`
	MatchType types.String `tfsdk:"match_type" json:"match_type,computed_optional"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsUserAgentModel struct {
	UserAgent types.String `tfsdk:"user_agent" json:"user_agent,required"`
	MatchType types.String `tfsdk:"match_type" json:"match_type,computed_optional"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainCustomRuleConditionsUserDefinedTagsModel struct {
	Tags     *[]types.String `tfsdk:"tags" json:"tags,required"`
	Negation types.Bool      `tfsdk:"negation" json:"negation,computed_optional"`
}

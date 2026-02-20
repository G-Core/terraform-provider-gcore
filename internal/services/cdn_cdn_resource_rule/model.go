// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_cdn_resource_rule

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CDNCDNResourceRuleModel struct {
	ID                     types.Int64                                              `tfsdk:"id" json:"id,computed"`
	ResourceID             types.Int64                                              `tfsdk:"resource_id" path:"resource_id,required"`
	Name                   types.String                                             `tfsdk:"name" json:"name,required"`
	Rule                   types.String                                             `tfsdk:"rule" json:"rule,required"`
	RuleType               types.Int64                                              `tfsdk:"rule_type" json:"ruleType,required"`
	Active                 types.Bool                                               `tfsdk:"active" json:"active,optional"`
	OriginGroup            types.Int64                                              `tfsdk:"origin_group" json:"originGroup,optional"`
	OverrideOriginProtocol types.String                                             `tfsdk:"override_origin_protocol" json:"overrideOriginProtocol,optional"`
	Weight                 types.Int64                                              `tfsdk:"weight" json:"weight,optional"`
	Options                customfield.NestedObject[CDNCDNResourceRuleOptionsModel] `tfsdk:"options" json:"options,computed_optional"`
	Deleted                types.Bool                                               `tfsdk:"deleted" json:"deleted,computed"`
	OriginProtocol         types.String                                             `tfsdk:"origin_protocol" json:"originProtocol,computed"`
	PresetApplied          types.Bool                                               `tfsdk:"preset_applied" json:"preset_applied,computed"`
	PrimaryRule            types.Int64                                              `tfsdk:"primary_rule" json:"primary_rule,computed"`
}

func (m CDNCDNResourceRuleModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CDNCDNResourceRuleModel) MarshalJSONForUpdate(state CDNCDNResourceRuleModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CDNCDNResourceRuleOptionsModel struct {
	AllowedHTTPMethods          *CDNCDNResourceRuleOptionsAllowedHTTPMethodsModel                         `tfsdk:"allowed_http_methods" json:"allowedHttpMethods,optional"`
	BotProtection               *CDNCDNResourceRuleOptionsBotProtectionModel                              `tfsdk:"bot_protection" json:"bot_protection,optional"`
	BrotliCompression           *CDNCDNResourceRuleOptionsBrotliCompressionModel                          `tfsdk:"brotli_compression" json:"brotli_compression,optional"`
	BrowserCacheSettings        *CDNCDNResourceRuleOptionsBrowserCacheSettingsModel                       `tfsdk:"browser_cache_settings" json:"browser_cache_settings,optional"`
	CacheHTTPHeaders            *CDNCDNResourceRuleOptionsCacheHTTPHeadersModel                           `tfsdk:"cache_http_headers" json:"cache_http_headers,optional"`
	Cors                        customfield.NestedObject[CDNCDNResourceRuleOptionsCorsModel]              `tfsdk:"cors" json:"cors,computed_optional"`
	CountryACL                  *CDNCDNResourceRuleOptionsCountryACLModel                                 `tfsdk:"country_acl" json:"country_acl,optional"`
	DisableCache                *CDNCDNResourceRuleOptionsDisableCacheModel                               `tfsdk:"disable_cache" json:"disable_cache,optional"`
	DisableProxyForceRanges     *CDNCDNResourceRuleOptionsDisableProxyForceRangesModel                    `tfsdk:"disable_proxy_force_ranges" json:"disable_proxy_force_ranges,optional"`
	EdgeCacheSettings           customfield.NestedObject[CDNCDNResourceRuleOptionsEdgeCacheSettingsModel] `tfsdk:"edge_cache_settings" json:"edge_cache_settings,computed_optional"`
	Fastedge                    customfield.NestedObject[CDNCDNResourceRuleOptionsFastedgeModel]          `tfsdk:"fastedge" json:"fastedge,computed_optional"`
	FetchCompressed             *CDNCDNResourceRuleOptionsFetchCompressedModel                            `tfsdk:"fetch_compressed" json:"fetch_compressed,optional"`
	FollowOriginRedirect        *CDNCDNResourceRuleOptionsFollowOriginRedirectModel                       `tfsdk:"follow_origin_redirect" json:"follow_origin_redirect,optional"`
	ForceReturn                 customfield.NestedObject[CDNCDNResourceRuleOptionsForceReturnModel]       `tfsdk:"force_return" json:"force_return,computed_optional"`
	ForwardHostHeader           *CDNCDNResourceRuleOptionsForwardHostHeaderModel                          `tfsdk:"forward_host_header" json:"forward_host_header,optional"`
	GzipOn                      *CDNCDNResourceRuleOptionsGzipOnModel                                     `tfsdk:"gzip_on" json:"gzipOn,optional"`
	HostHeader                  *CDNCDNResourceRuleOptionsHostHeaderModel                                 `tfsdk:"host_header" json:"hostHeader,optional"`
	IgnoreCookie                *CDNCDNResourceRuleOptionsIgnoreCookieModel                               `tfsdk:"ignore_cookie" json:"ignore_cookie,optional"`
	IgnoreQueryString           *CDNCDNResourceRuleOptionsIgnoreQueryStringModel                          `tfsdk:"ignore_query_string" json:"ignoreQueryString,optional"`
	ImageStack                  customfield.NestedObject[CDNCDNResourceRuleOptionsImageStackModel]        `tfsdk:"image_stack" json:"image_stack,computed_optional"`
	IPAddressACL                *CDNCDNResourceRuleOptionsIPAddressACLModel                               `tfsdk:"ip_address_acl" json:"ip_address_acl,optional"`
	LimitBandwidth              *CDNCDNResourceRuleOptionsLimitBandwidthModel                             `tfsdk:"limit_bandwidth" json:"limit_bandwidth,optional"`
	ProxyCacheKey               *CDNCDNResourceRuleOptionsProxyCacheKeyModel                              `tfsdk:"proxy_cache_key" json:"proxy_cache_key,optional"`
	ProxyCacheMethodsSet        *CDNCDNResourceRuleOptionsProxyCacheMethodsSetModel                       `tfsdk:"proxy_cache_methods_set" json:"proxy_cache_methods_set,optional"`
	ProxyConnectTimeout         *CDNCDNResourceRuleOptionsProxyConnectTimeoutModel                        `tfsdk:"proxy_connect_timeout" json:"proxy_connect_timeout,optional"`
	ProxyReadTimeout            *CDNCDNResourceRuleOptionsProxyReadTimeoutModel                           `tfsdk:"proxy_read_timeout" json:"proxy_read_timeout,optional"`
	QueryParamsBlacklist        *CDNCDNResourceRuleOptionsQueryParamsBlacklistModel                       `tfsdk:"query_params_blacklist" json:"query_params_blacklist,optional"`
	QueryParamsWhitelist        *CDNCDNResourceRuleOptionsQueryParamsWhitelistModel                       `tfsdk:"query_params_whitelist" json:"query_params_whitelist,optional"`
	QueryStringForwarding       *CDNCDNResourceRuleOptionsQueryStringForwardingModel                      `tfsdk:"query_string_forwarding" json:"query_string_forwarding,optional"`
	RedirectHTTPToHTTPS         *CDNCDNResourceRuleOptionsRedirectHTTPToHTTPSModel                        `tfsdk:"redirect_http_to_https" json:"redirect_http_to_https,optional"`
	RedirectHTTPSToHTTP         *CDNCDNResourceRuleOptionsRedirectHTTPSToHTTPModel                        `tfsdk:"redirect_https_to_http" json:"redirect_https_to_http,optional"`
	ReferrerACL                 *CDNCDNResourceRuleOptionsReferrerACLModel                                `tfsdk:"referrer_acl" json:"referrer_acl,optional"`
	RequestLimiter              customfield.NestedObject[CDNCDNResourceRuleOptionsRequestLimiterModel]    `tfsdk:"request_limiter" json:"request_limiter,computed_optional"`
	ResponseHeadersHidingPolicy *CDNCDNResourceRuleOptionsResponseHeadersHidingPolicyModel                `tfsdk:"response_headers_hiding_policy" json:"response_headers_hiding_policy,optional"`
	Rewrite                     customfield.NestedObject[CDNCDNResourceRuleOptionsRewriteModel]           `tfsdk:"rewrite" json:"rewrite,computed_optional"`
	SecureKey                   customfield.NestedObject[CDNCDNResourceRuleOptionsSecureKeyModel]         `tfsdk:"secure_key" json:"secure_key,computed_optional"`
	Slice                       *CDNCDNResourceRuleOptionsSliceModel                                      `tfsdk:"slice" json:"slice,optional"`
	Sni                         customfield.NestedObject[CDNCDNResourceRuleOptionsSniModel]               `tfsdk:"sni" json:"sni,computed_optional"`
	Stale                       *CDNCDNResourceRuleOptionsStaleModel                                      `tfsdk:"stale" json:"stale,optional"`
	StaticResponseHeaders       *CDNCDNResourceRuleOptionsStaticResponseHeadersModel                      `tfsdk:"static_response_headers" json:"static_response_headers,optional"`
	StaticHeaders               *CDNCDNResourceRuleOptionsStaticHeadersModel                              `tfsdk:"static_headers" json:"staticHeaders,optional"`
	StaticRequestHeaders        *CDNCDNResourceRuleOptionsStaticRequestHeadersModel                       `tfsdk:"static_request_headers" json:"staticRequestHeaders,optional"`
	UserAgentACL                *CDNCDNResourceRuleOptionsUserAgentACLModel                               `tfsdk:"user_agent_acl" json:"user_agent_acl,optional"`
	Waap                        *CDNCDNResourceRuleOptionsWaapModel                                       `tfsdk:"waap" json:"waap,optional"`
	Websockets                  *CDNCDNResourceRuleOptionsWebsocketsModel                                 `tfsdk:"websockets" json:"websockets,optional"`
}

type CDNCDNResourceRuleOptionsAllowedHTTPMethodsModel struct {
	Enabled types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]types.String `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsBotProtectionModel struct {
	BotChallenge *CDNCDNResourceRuleOptionsBotProtectionBotChallengeModel `tfsdk:"bot_challenge" json:"bot_challenge,required"`
	Enabled      types.Bool                                               `tfsdk:"enabled" json:"enabled,required"`
}

type CDNCDNResourceRuleOptionsBotProtectionBotChallengeModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed_optional"`
}

type CDNCDNResourceRuleOptionsBrotliCompressionModel struct {
	Enabled types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]types.String `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsBrowserCacheSettingsModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsCacheHTTPHeadersModel struct {
	Enabled types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]types.String `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsCorsModel struct {
	Enabled types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]types.String `tfsdk:"value" json:"value,required"`
	Always  types.Bool      `tfsdk:"always" json:"always,computed_optional"`
}

type CDNCDNResourceRuleOptionsCountryACLModel struct {
	Enabled        types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues *[]types.String `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String    `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNCDNResourceRuleOptionsDisableCacheModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsDisableProxyForceRangesModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsEdgeCacheSettingsModel struct {
	Enabled      types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	CustomValues customfield.Map[types.String] `tfsdk:"custom_values" json:"custom_values,computed_optional"`
	Default      types.String                  `tfsdk:"default" json:"default,optional"`
	Value        types.String                  `tfsdk:"value" json:"value,optional"`
}

type CDNCDNResourceRuleOptionsFastedgeModel struct {
	Enabled           types.Bool                                                                        `tfsdk:"enabled" json:"enabled,required"`
	OnRequestBody     customfield.NestedObject[CDNCDNResourceRuleOptionsFastedgeOnRequestBodyModel]     `tfsdk:"on_request_body" json:"on_request_body,computed_optional"`
	OnRequestHeaders  customfield.NestedObject[CDNCDNResourceRuleOptionsFastedgeOnRequestHeadersModel]  `tfsdk:"on_request_headers" json:"on_request_headers,computed_optional"`
	OnResponseBody    customfield.NestedObject[CDNCDNResourceRuleOptionsFastedgeOnResponseBodyModel]    `tfsdk:"on_response_body" json:"on_response_body,computed_optional"`
	OnResponseHeaders customfield.NestedObject[CDNCDNResourceRuleOptionsFastedgeOnResponseHeadersModel] `tfsdk:"on_response_headers" json:"on_response_headers,computed_optional"`
}

type CDNCDNResourceRuleOptionsFastedgeOnRequestBodyModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNCDNResourceRuleOptionsFastedgeOnRequestHeadersModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNCDNResourceRuleOptionsFastedgeOnResponseBodyModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNCDNResourceRuleOptionsFastedgeOnResponseHeadersModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNCDNResourceRuleOptionsFetchCompressedModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsFollowOriginRedirectModel struct {
	Codes   *[]types.Int64 `tfsdk:"codes" json:"codes,required"`
	Enabled types.Bool     `tfsdk:"enabled" json:"enabled,required"`
}

type CDNCDNResourceRuleOptionsForceReturnModel struct {
	Body         types.String                                                                    `tfsdk:"body" json:"body,required"`
	Code         types.Int64                                                                     `tfsdk:"code" json:"code,required"`
	Enabled      types.Bool                                                                      `tfsdk:"enabled" json:"enabled,required"`
	TimeInterval customfield.NestedObject[CDNCDNResourceRuleOptionsForceReturnTimeIntervalModel] `tfsdk:"time_interval" json:"time_interval,computed_optional"`
}

type CDNCDNResourceRuleOptionsForceReturnTimeIntervalModel struct {
	EndTime   types.String `tfsdk:"end_time" json:"end_time,required"`
	StartTime types.String `tfsdk:"start_time" json:"start_time,required"`
	TimeZone  types.String `tfsdk:"time_zone" json:"time_zone,computed_optional"`
}

type CDNCDNResourceRuleOptionsForwardHostHeaderModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsGzipOnModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsHostHeaderModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsIgnoreCookieModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsIgnoreQueryStringModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsImageStackModel struct {
	Enabled     types.Bool  `tfsdk:"enabled" json:"enabled,required"`
	AvifEnabled types.Bool  `tfsdk:"avif_enabled" json:"avif_enabled,computed_optional"`
	PngLossless types.Bool  `tfsdk:"png_lossless" json:"png_lossless,computed_optional"`
	Quality     types.Int64 `tfsdk:"quality" json:"quality,computed_optional"`
	WebpEnabled types.Bool  `tfsdk:"webp_enabled" json:"webp_enabled,computed_optional"`
}

type CDNCDNResourceRuleOptionsIPAddressACLModel struct {
	Enabled        types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues *[]types.String `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String    `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNCDNResourceRuleOptionsLimitBandwidthModel struct {
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	LimitType types.String `tfsdk:"limit_type" json:"limit_type,required"`
	Buffer    types.Int64  `tfsdk:"buffer" json:"buffer,optional"`
	Speed     types.Int64  `tfsdk:"speed" json:"speed,optional"`
}

type CDNCDNResourceRuleOptionsProxyCacheKeyModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsProxyCacheMethodsSetModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsProxyConnectTimeoutModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsProxyReadTimeoutModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsQueryParamsBlacklistModel struct {
	Enabled types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]types.String `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsQueryParamsWhitelistModel struct {
	Enabled types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]types.String `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsQueryStringForwardingModel struct {
	Enabled              types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	ForwardFromFileTypes *[]types.String `tfsdk:"forward_from_file_types" json:"forward_from_file_types,required"`
	ForwardToFileTypes   *[]types.String `tfsdk:"forward_to_file_types" json:"forward_to_file_types,required"`
	ForwardExceptKeys    *[]types.String `tfsdk:"forward_except_keys" json:"forward_except_keys,optional"`
	ForwardOnlyKeys      *[]types.String `tfsdk:"forward_only_keys" json:"forward_only_keys,optional"`
}

type CDNCDNResourceRuleOptionsRedirectHTTPToHTTPSModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsRedirectHTTPSToHTTPModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsReferrerACLModel struct {
	Enabled        types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues *[]types.String `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String    `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNCDNResourceRuleOptionsRequestLimiterModel struct {
	Enabled  types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Rate     types.Int64  `tfsdk:"rate" json:"rate,required"`
	Burst    types.Int64  `tfsdk:"burst" json:"burst,computed"`
	Delay    types.Int64  `tfsdk:"delay" json:"delay,computed"`
	RateUnit types.String `tfsdk:"rate_unit" json:"rate_unit,computed_optional"`
}

type CDNCDNResourceRuleOptionsResponseHeadersHidingPolicyModel struct {
	Enabled  types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Excepted *[]types.String `tfsdk:"excepted" json:"excepted,required"`
	Mode     types.String    `tfsdk:"mode" json:"mode,required"`
}

type CDNCDNResourceRuleOptionsRewriteModel struct {
	Body    types.String `tfsdk:"body" json:"body,required"`
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Flag    types.String `tfsdk:"flag" json:"flag,computed_optional"`
}

type CDNCDNResourceRuleOptionsSecureKeyModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Key     types.String `tfsdk:"key" json:"key,required"`
	Type    types.Int64  `tfsdk:"type" json:"type,computed_optional"`
}

type CDNCDNResourceRuleOptionsSliceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsSniModel struct {
	CustomHostname types.String `tfsdk:"custom_hostname" json:"custom_hostname,required"`
	Enabled        types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	SniType        types.String `tfsdk:"sni_type" json:"sni_type,computed_optional"`
}

type CDNCDNResourceRuleOptionsStaleModel struct {
	Enabled types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]types.String `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsStaticResponseHeadersModel struct {
	Enabled types.Bool                                                   `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]*CDNCDNResourceRuleOptionsStaticResponseHeadersValueModel `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsStaticResponseHeadersValueModel struct {
	Name   types.String    `tfsdk:"name" json:"name,required"`
	Value  *[]types.String `tfsdk:"value" json:"value,required"`
	Always types.Bool      `tfsdk:"always" json:"always,computed_optional"`
}

type CDNCDNResourceRuleOptionsStaticHeadersModel struct {
	Enabled types.Bool           `tfsdk:"enabled" json:"enabled,required"`
	Value   jsontypes.Normalized `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsStaticRequestHeadersModel struct {
	Enabled types.Bool               `tfsdk:"enabled" json:"enabled,required"`
	Value   *map[string]types.String `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsUserAgentACLModel struct {
	Enabled        types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues *[]types.String `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String    `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNCDNResourceRuleOptionsWaapModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNCDNResourceRuleOptionsWebsocketsModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

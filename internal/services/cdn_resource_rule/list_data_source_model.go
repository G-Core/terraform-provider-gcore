// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_resource_rule

import (
	"context"

	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNResourceRulesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CDNResourceRulesItemsDataSourceModel] `json:"results,computed"`
}

type CDNResourceRulesDataSourceModel struct {
	ResourceID types.Int64                                                        `tfsdk:"resource_id" path:"resource_id,required"`
	MaxItems   types.Int64                                                        `tfsdk:"max_items"`
	Items      customfield.NestedObjectList[CDNResourceRulesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CDNResourceRulesDataSourceModel) toListParams(_ context.Context) (params cdn.CDNResourceRuleListParams, diags diag.Diagnostics) {
	params = cdn.CDNResourceRuleListParams{}

	return
}

type CDNResourceRulesItemsDataSourceModel struct {
	ID                     types.Int64                                                      `tfsdk:"id" json:"id,computed"`
	Active                 types.Bool                                                       `tfsdk:"active" json:"active,computed"`
	Deleted                types.Bool                                                       `tfsdk:"deleted" json:"deleted,computed"`
	Name                   types.String                                                     `tfsdk:"name" json:"name,computed"`
	Options                customfield.NestedObject[CDNResourceRulesOptionsDataSourceModel] `tfsdk:"options" json:"options,computed"`
	OriginGroup            types.Int64                                                      `tfsdk:"origin_group" json:"originGroup,computed"`
	OriginProtocol         types.String                                                     `tfsdk:"origin_protocol" json:"originProtocol,computed"`
	OverrideOriginProtocol types.String                                                     `tfsdk:"override_origin_protocol" json:"overrideOriginProtocol,computed"`
	PresetApplied          types.Bool                                                       `tfsdk:"preset_applied" json:"preset_applied,computed"`
	PrimaryRule            types.Int64                                                      `tfsdk:"primary_rule" json:"primary_rule,computed"`
	Rule                   types.String                                                     `tfsdk:"rule" json:"rule,computed"`
	RuleType               types.Int64                                                      `tfsdk:"rule_type" json:"ruleType,computed"`
	Weight                 types.Int64                                                      `tfsdk:"weight" json:"weight,computed"`
}

type CDNResourceRulesOptionsDataSourceModel struct {
	AllowedHTTPMethods          customfield.NestedObject[CDNResourceRulesOptionsAllowedHTTPMethodsDataSourceModel]          `tfsdk:"allowed_http_methods" json:"allowedHttpMethods,computed"`
	BotProtection               customfield.NestedObject[CDNResourceRulesOptionsBotProtectionDataSourceModel]               `tfsdk:"bot_protection" json:"bot_protection,computed"`
	BrotliCompression           customfield.NestedObject[CDNResourceRulesOptionsBrotliCompressionDataSourceModel]           `tfsdk:"brotli_compression" json:"brotli_compression,computed"`
	BrowserCacheSettings        customfield.NestedObject[CDNResourceRulesOptionsBrowserCacheSettingsDataSourceModel]        `tfsdk:"browser_cache_settings" json:"browser_cache_settings,computed"`
	CacheHTTPHeaders            customfield.NestedObject[CDNResourceRulesOptionsCacheHTTPHeadersDataSourceModel]            `tfsdk:"cache_http_headers" json:"cache_http_headers,computed"`
	Cors                        customfield.NestedObject[CDNResourceRulesOptionsCorsDataSourceModel]                        `tfsdk:"cors" json:"cors,computed"`
	CountryACL                  customfield.NestedObject[CDNResourceRulesOptionsCountryACLDataSourceModel]                  `tfsdk:"country_acl" json:"country_acl,computed"`
	DisableCache                customfield.NestedObject[CDNResourceRulesOptionsDisableCacheDataSourceModel]                `tfsdk:"disable_cache" json:"disable_cache,computed"`
	DisableProxyForceRanges     customfield.NestedObject[CDNResourceRulesOptionsDisableProxyForceRangesDataSourceModel]     `tfsdk:"disable_proxy_force_ranges" json:"disable_proxy_force_ranges,computed"`
	EdgeCacheSettings           customfield.NestedObject[CDNResourceRulesOptionsEdgeCacheSettingsDataSourceModel]           `tfsdk:"edge_cache_settings" json:"edge_cache_settings,computed"`
	Fastedge                    customfield.NestedObject[CDNResourceRulesOptionsFastedgeDataSourceModel]                    `tfsdk:"fastedge" json:"fastedge,computed"`
	FetchCompressed             customfield.NestedObject[CDNResourceRulesOptionsFetchCompressedDataSourceModel]             `tfsdk:"fetch_compressed" json:"fetch_compressed,computed"`
	FollowOriginRedirect        customfield.NestedObject[CDNResourceRulesOptionsFollowOriginRedirectDataSourceModel]        `tfsdk:"follow_origin_redirect" json:"follow_origin_redirect,computed"`
	ForceReturn                 customfield.NestedObject[CDNResourceRulesOptionsForceReturnDataSourceModel]                 `tfsdk:"force_return" json:"force_return,computed"`
	ForwardHostHeader           customfield.NestedObject[CDNResourceRulesOptionsForwardHostHeaderDataSourceModel]           `tfsdk:"forward_host_header" json:"forward_host_header,computed"`
	GzipOn                      customfield.NestedObject[CDNResourceRulesOptionsGzipOnDataSourceModel]                      `tfsdk:"gzip_on" json:"gzipOn,computed"`
	HostHeader                  customfield.NestedObject[CDNResourceRulesOptionsHostHeaderDataSourceModel]                  `tfsdk:"host_header" json:"hostHeader,computed"`
	IgnoreCookie                customfield.NestedObject[CDNResourceRulesOptionsIgnoreCookieDataSourceModel]                `tfsdk:"ignore_cookie" json:"ignore_cookie,computed"`
	IgnoreQueryString           customfield.NestedObject[CDNResourceRulesOptionsIgnoreQueryStringDataSourceModel]           `tfsdk:"ignore_query_string" json:"ignoreQueryString,computed"`
	ImageStack                  customfield.NestedObject[CDNResourceRulesOptionsImageStackDataSourceModel]                  `tfsdk:"image_stack" json:"image_stack,computed"`
	IPAddressACL                customfield.NestedObject[CDNResourceRulesOptionsIPAddressACLDataSourceModel]                `tfsdk:"ip_address_acl" json:"ip_address_acl,computed"`
	LimitBandwidth              customfield.NestedObject[CDNResourceRulesOptionsLimitBandwidthDataSourceModel]              `tfsdk:"limit_bandwidth" json:"limit_bandwidth,computed"`
	ProxyCacheKey               customfield.NestedObject[CDNResourceRulesOptionsProxyCacheKeyDataSourceModel]               `tfsdk:"proxy_cache_key" json:"proxy_cache_key,computed"`
	ProxyCacheMethodsSet        customfield.NestedObject[CDNResourceRulesOptionsProxyCacheMethodsSetDataSourceModel]        `tfsdk:"proxy_cache_methods_set" json:"proxy_cache_methods_set,computed"`
	ProxyConnectTimeout         customfield.NestedObject[CDNResourceRulesOptionsProxyConnectTimeoutDataSourceModel]         `tfsdk:"proxy_connect_timeout" json:"proxy_connect_timeout,computed"`
	ProxyReadTimeout            customfield.NestedObject[CDNResourceRulesOptionsProxyReadTimeoutDataSourceModel]            `tfsdk:"proxy_read_timeout" json:"proxy_read_timeout,computed"`
	QueryParamsBlacklist        customfield.NestedObject[CDNResourceRulesOptionsQueryParamsBlacklistDataSourceModel]        `tfsdk:"query_params_blacklist" json:"query_params_blacklist,computed"`
	QueryParamsWhitelist        customfield.NestedObject[CDNResourceRulesOptionsQueryParamsWhitelistDataSourceModel]        `tfsdk:"query_params_whitelist" json:"query_params_whitelist,computed"`
	QueryStringForwarding       customfield.NestedObject[CDNResourceRulesOptionsQueryStringForwardingDataSourceModel]       `tfsdk:"query_string_forwarding" json:"query_string_forwarding,computed"`
	RedirectHTTPToHTTPS         customfield.NestedObject[CDNResourceRulesOptionsRedirectHTTPToHTTPSDataSourceModel]         `tfsdk:"redirect_http_to_https" json:"redirect_http_to_https,computed"`
	RedirectHTTPSToHTTP         customfield.NestedObject[CDNResourceRulesOptionsRedirectHTTPSToHTTPDataSourceModel]         `tfsdk:"redirect_https_to_http" json:"redirect_https_to_http,computed"`
	ReferrerACL                 customfield.NestedObject[CDNResourceRulesOptionsReferrerACLDataSourceModel]                 `tfsdk:"referrer_acl" json:"referrer_acl,computed"`
	RequestLimiter              customfield.NestedObject[CDNResourceRulesOptionsRequestLimiterDataSourceModel]              `tfsdk:"request_limiter" json:"request_limiter,computed"`
	ResponseHeadersHidingPolicy customfield.NestedObject[CDNResourceRulesOptionsResponseHeadersHidingPolicyDataSourceModel] `tfsdk:"response_headers_hiding_policy" json:"response_headers_hiding_policy,computed"`
	Rewrite                     customfield.NestedObject[CDNResourceRulesOptionsRewriteDataSourceModel]                     `tfsdk:"rewrite" json:"rewrite,computed"`
	SecureKey                   customfield.NestedObject[CDNResourceRulesOptionsSecureKeyDataSourceModel]                   `tfsdk:"secure_key" json:"secure_key,computed"`
	Slice                       customfield.NestedObject[CDNResourceRulesOptionsSliceDataSourceModel]                       `tfsdk:"slice" json:"slice,computed"`
	Sni                         customfield.NestedObject[CDNResourceRulesOptionsSniDataSourceModel]                         `tfsdk:"sni" json:"sni,computed"`
	Stale                       customfield.NestedObject[CDNResourceRulesOptionsStaleDataSourceModel]                       `tfsdk:"stale" json:"stale,computed"`
	StaticResponseHeaders       customfield.NestedObject[CDNResourceRulesOptionsStaticResponseHeadersDataSourceModel]       `tfsdk:"static_response_headers" json:"static_response_headers,computed"`
	StaticHeaders               customfield.NestedObject[CDNResourceRulesOptionsStaticHeadersDataSourceModel]               `tfsdk:"static_headers" json:"staticHeaders,computed"`
	StaticRequestHeaders        customfield.NestedObject[CDNResourceRulesOptionsStaticRequestHeadersDataSourceModel]        `tfsdk:"static_request_headers" json:"staticRequestHeaders,computed"`
	UserAgentACL                customfield.NestedObject[CDNResourceRulesOptionsUserAgentACLDataSourceModel]                `tfsdk:"user_agent_acl" json:"user_agent_acl,computed"`
	Waap                        customfield.NestedObject[CDNResourceRulesOptionsWaapDataSourceModel]                        `tfsdk:"waap" json:"waap,computed"`
	Websockets                  customfield.NestedObject[CDNResourceRulesOptionsWebsocketsDataSourceModel]                  `tfsdk:"websockets" json:"websockets,computed"`
}

type CDNResourceRulesOptionsAllowedHTTPMethodsDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsBotProtectionDataSourceModel struct {
	BotChallenge customfield.NestedObject[CDNResourceRulesOptionsBotProtectionBotChallengeDataSourceModel] `tfsdk:"bot_challenge" json:"bot_challenge,computed"`
	Enabled      types.Bool                                                                                `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNResourceRulesOptionsBotProtectionBotChallengeDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNResourceRulesOptionsBrotliCompressionDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsBrowserCacheSettingsDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsCacheHTTPHeadersDataSourceModel struct {
	Enabled types.Bool                     `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsCorsDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
	Always  types.Bool                    `tfsdk:"always" json:"always,computed"`
}

type CDNResourceRulesOptionsCountryACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourceRulesOptionsDisableCacheDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsDisableProxyForceRangesDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsEdgeCacheSettingsDataSourceModel struct {
	Enabled      types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	CustomValues customfield.Map[types.String] `tfsdk:"custom_values" json:"custom_values,computed"`
	Default      types.String                  `tfsdk:"default" json:"default,computed"`
	Value        types.String                  `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsFastedgeDataSourceModel struct {
	Enabled           types.Bool                                                                                `tfsdk:"enabled" json:"enabled,computed"`
	OnRequestBody     customfield.NestedObject[CDNResourceRulesOptionsFastedgeOnRequestBodyDataSourceModel]     `tfsdk:"on_request_body" json:"on_request_body,computed"`
	OnRequestHeaders  customfield.NestedObject[CDNResourceRulesOptionsFastedgeOnRequestHeadersDataSourceModel]  `tfsdk:"on_request_headers" json:"on_request_headers,computed"`
	OnResponseBody    customfield.NestedObject[CDNResourceRulesOptionsFastedgeOnResponseBodyDataSourceModel]    `tfsdk:"on_response_body" json:"on_response_body,computed"`
	OnResponseHeaders customfield.NestedObject[CDNResourceRulesOptionsFastedgeOnResponseHeadersDataSourceModel] `tfsdk:"on_response_headers" json:"on_response_headers,computed"`
}

type CDNResourceRulesOptionsFastedgeOnRequestBodyDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourceRulesOptionsFastedgeOnRequestHeadersDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourceRulesOptionsFastedgeOnResponseBodyDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourceRulesOptionsFastedgeOnResponseHeadersDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourceRulesOptionsFetchCompressedDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsFollowOriginRedirectDataSourceModel struct {
	Codes   customfield.Set[types.Int64] `tfsdk:"codes" json:"codes,computed"`
	Enabled types.Bool                   `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNResourceRulesOptionsForceReturnDataSourceModel struct {
	Body         types.String                                                                            `tfsdk:"body" json:"body,computed"`
	Code         types.Int64                                                                             `tfsdk:"code" json:"code,computed"`
	Enabled      types.Bool                                                                              `tfsdk:"enabled" json:"enabled,computed"`
	TimeInterval customfield.NestedObject[CDNResourceRulesOptionsForceReturnTimeIntervalDataSourceModel] `tfsdk:"time_interval" json:"time_interval,computed"`
}

type CDNResourceRulesOptionsForceReturnTimeIntervalDataSourceModel struct {
	EndTime   types.String `tfsdk:"end_time" json:"end_time,computed"`
	StartTime types.String `tfsdk:"start_time" json:"start_time,computed"`
	TimeZone  types.String `tfsdk:"time_zone" json:"time_zone,computed"`
}

type CDNResourceRulesOptionsForwardHostHeaderDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsGzipOnDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsHostHeaderDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsIgnoreCookieDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsIgnoreQueryStringDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsImageStackDataSourceModel struct {
	Enabled     types.Bool  `tfsdk:"enabled" json:"enabled,computed"`
	AvifEnabled types.Bool  `tfsdk:"avif_enabled" json:"avif_enabled,computed"`
	PngLossless types.Bool  `tfsdk:"png_lossless" json:"png_lossless,computed"`
	Quality     types.Int64 `tfsdk:"quality" json:"quality,computed"`
	WebpEnabled types.Bool  `tfsdk:"webp_enabled" json:"webp_enabled,computed"`
}

type CDNResourceRulesOptionsIPAddressACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourceRulesOptionsLimitBandwidthDataSourceModel struct {
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	LimitType types.String `tfsdk:"limit_type" json:"limit_type,computed"`
	Buffer    types.Int64  `tfsdk:"buffer" json:"buffer,computed"`
	Speed     types.Int64  `tfsdk:"speed" json:"speed,computed"`
}

type CDNResourceRulesOptionsProxyCacheKeyDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsProxyCacheMethodsSetDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsProxyConnectTimeoutDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsProxyReadTimeoutDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsQueryParamsBlacklistDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsQueryParamsWhitelistDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsQueryStringForwardingDataSourceModel struct {
	Enabled              types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ForwardFromFileTypes customfield.Set[types.String] `tfsdk:"forward_from_file_types" json:"forward_from_file_types,computed"`
	ForwardToFileTypes   customfield.Set[types.String] `tfsdk:"forward_to_file_types" json:"forward_to_file_types,computed"`
	ForwardExceptKeys    customfield.Set[types.String] `tfsdk:"forward_except_keys" json:"forward_except_keys,computed"`
	ForwardOnlyKeys      customfield.Set[types.String] `tfsdk:"forward_only_keys" json:"forward_only_keys,computed"`
}

type CDNResourceRulesOptionsRedirectHTTPToHTTPSDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsRedirectHTTPSToHTTPDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsReferrerACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourceRulesOptionsRequestLimiterDataSourceModel struct {
	Enabled  types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Rate     types.Int64  `tfsdk:"rate" json:"rate,computed"`
	Burst    types.Int64  `tfsdk:"burst" json:"burst,computed"`
	Delay    types.Int64  `tfsdk:"delay" json:"delay,computed"`
	RateUnit types.String `tfsdk:"rate_unit" json:"rate_unit,computed"`
}

type CDNResourceRulesOptionsResponseHeadersHidingPolicyDataSourceModel struct {
	Enabled  types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Excepted customfield.Set[types.String] `tfsdk:"excepted" json:"excepted,computed"`
	Mode     types.String                  `tfsdk:"mode" json:"mode,computed"`
}

type CDNResourceRulesOptionsRewriteDataSourceModel struct {
	Body    types.String `tfsdk:"body" json:"body,computed"`
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Flag    types.String `tfsdk:"flag" json:"flag,computed"`
}

type CDNResourceRulesOptionsSecureKeyDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Key     types.String `tfsdk:"key" json:"key,computed"`
	Type    types.Int64  `tfsdk:"type" json:"type,computed"`
}

type CDNResourceRulesOptionsSliceDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsSniDataSourceModel struct {
	CustomHostname types.String `tfsdk:"custom_hostname" json:"custom_hostname,computed"`
	Enabled        types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	SniType        types.String `tfsdk:"sni_type" json:"sni_type,computed"`
}

type CDNResourceRulesOptionsStaleDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsStaticResponseHeadersDataSourceModel struct {
	Enabled types.Bool                                                                                     `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.NestedObjectList[CDNResourceRulesOptionsStaticResponseHeadersValueDataSourceModel] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsStaticResponseHeadersValueDataSourceModel struct {
	Name   types.String                   `tfsdk:"name" json:"name,computed"`
	Value  customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
	Always types.Bool                     `tfsdk:"always" json:"always,computed"`
}

type CDNResourceRulesOptionsStaticHeadersDataSourceModel struct {
	Enabled types.Bool           `tfsdk:"enabled" json:"enabled,computed"`
	Value   jsontypes.Normalized `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsStaticRequestHeadersDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Map[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsUserAgentACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourceRulesOptionsWaapDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRulesOptionsWebsocketsDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

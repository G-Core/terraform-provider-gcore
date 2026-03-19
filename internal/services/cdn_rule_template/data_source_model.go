// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_rule_template

import (
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNRuleTemplateDataSourceModel struct {
	ID                     types.Int64                                                     `tfsdk:"id" path:"rule_template_id,computed"`
	RuleTemplateID         types.Int64                                                     `tfsdk:"rule_template_id" path:"rule_template_id,required"`
	Client                 types.Int64                                                     `tfsdk:"client" json:"client,computed"`
	Default                types.Bool                                                      `tfsdk:"default" json:"default,computed"`
	Deleted                types.Bool                                                      `tfsdk:"deleted" json:"deleted,computed"`
	Name                   types.String                                                    `tfsdk:"name" json:"name,computed"`
	OverrideOriginProtocol types.String                                                    `tfsdk:"override_origin_protocol" json:"overrideOriginProtocol,computed"`
	Rule                   types.String                                                    `tfsdk:"rule" json:"rule,computed"`
	RuleType               types.Int64                                                     `tfsdk:"rule_type" json:"ruleType,computed"`
	Template               types.Bool                                                      `tfsdk:"template" json:"template,computed"`
	Weight                 types.Int64                                                     `tfsdk:"weight" json:"weight,computed"`
	Options                customfield.NestedObject[CDNRuleTemplateOptionsDataSourceModel] `tfsdk:"options" json:"options,computed"`
}

type CDNRuleTemplateOptionsDataSourceModel struct {
	AllowedHTTPMethods          customfield.NestedObject[CDNRuleTemplateOptionsAllowedHTTPMethodsDataSourceModel]          `tfsdk:"allowed_http_methods" json:"allowedHttpMethods,computed"`
	BotProtection               customfield.NestedObject[CDNRuleTemplateOptionsBotProtectionDataSourceModel]               `tfsdk:"bot_protection" json:"bot_protection,computed"`
	BrotliCompression           customfield.NestedObject[CDNRuleTemplateOptionsBrotliCompressionDataSourceModel]           `tfsdk:"brotli_compression" json:"brotli_compression,computed"`
	BrowserCacheSettings        customfield.NestedObject[CDNRuleTemplateOptionsBrowserCacheSettingsDataSourceModel]        `tfsdk:"browser_cache_settings" json:"browser_cache_settings,computed"`
	CacheHTTPHeaders            customfield.NestedObject[CDNRuleTemplateOptionsCacheHTTPHeadersDataSourceModel]            `tfsdk:"cache_http_headers" json:"cache_http_headers,computed"`
	Cors                        customfield.NestedObject[CDNRuleTemplateOptionsCorsDataSourceModel]                        `tfsdk:"cors" json:"cors,computed"`
	CountryACL                  customfield.NestedObject[CDNRuleTemplateOptionsCountryACLDataSourceModel]                  `tfsdk:"country_acl" json:"country_acl,computed"`
	DisableCache                customfield.NestedObject[CDNRuleTemplateOptionsDisableCacheDataSourceModel]                `tfsdk:"disable_cache" json:"disable_cache,computed"`
	DisableProxyForceRanges     customfield.NestedObject[CDNRuleTemplateOptionsDisableProxyForceRangesDataSourceModel]     `tfsdk:"disable_proxy_force_ranges" json:"disable_proxy_force_ranges,computed"`
	EdgeCacheSettings           customfield.NestedObject[CDNRuleTemplateOptionsEdgeCacheSettingsDataSourceModel]           `tfsdk:"edge_cache_settings" json:"edge_cache_settings,computed"`
	Fastedge                    customfield.NestedObject[CDNRuleTemplateOptionsFastedgeDataSourceModel]                    `tfsdk:"fastedge" json:"fastedge,computed"`
	FetchCompressed             customfield.NestedObject[CDNRuleTemplateOptionsFetchCompressedDataSourceModel]             `tfsdk:"fetch_compressed" json:"fetch_compressed,computed"`
	FollowOriginRedirect        customfield.NestedObject[CDNRuleTemplateOptionsFollowOriginRedirectDataSourceModel]        `tfsdk:"follow_origin_redirect" json:"follow_origin_redirect,computed"`
	ForceReturn                 customfield.NestedObject[CDNRuleTemplateOptionsForceReturnDataSourceModel]                 `tfsdk:"force_return" json:"force_return,computed"`
	ForwardHostHeader           customfield.NestedObject[CDNRuleTemplateOptionsForwardHostHeaderDataSourceModel]           `tfsdk:"forward_host_header" json:"forward_host_header,computed"`
	GzipOn                      customfield.NestedObject[CDNRuleTemplateOptionsGzipOnDataSourceModel]                      `tfsdk:"gzip_on" json:"gzipOn,computed"`
	HostHeader                  customfield.NestedObject[CDNRuleTemplateOptionsHostHeaderDataSourceModel]                  `tfsdk:"host_header" json:"hostHeader,computed"`
	IgnoreCookie                customfield.NestedObject[CDNRuleTemplateOptionsIgnoreCookieDataSourceModel]                `tfsdk:"ignore_cookie" json:"ignore_cookie,computed"`
	IgnoreQueryString           customfield.NestedObject[CDNRuleTemplateOptionsIgnoreQueryStringDataSourceModel]           `tfsdk:"ignore_query_string" json:"ignoreQueryString,computed"`
	ImageStack                  customfield.NestedObject[CDNRuleTemplateOptionsImageStackDataSourceModel]                  `tfsdk:"image_stack" json:"image_stack,computed"`
	IPAddressACL                customfield.NestedObject[CDNRuleTemplateOptionsIPAddressACLDataSourceModel]                `tfsdk:"ip_address_acl" json:"ip_address_acl,computed"`
	LimitBandwidth              customfield.NestedObject[CDNRuleTemplateOptionsLimitBandwidthDataSourceModel]              `tfsdk:"limit_bandwidth" json:"limit_bandwidth,computed"`
	ProxyCacheKey               customfield.NestedObject[CDNRuleTemplateOptionsProxyCacheKeyDataSourceModel]               `tfsdk:"proxy_cache_key" json:"proxy_cache_key,computed"`
	ProxyCacheMethodsSet        customfield.NestedObject[CDNRuleTemplateOptionsProxyCacheMethodsSetDataSourceModel]        `tfsdk:"proxy_cache_methods_set" json:"proxy_cache_methods_set,computed"`
	ProxyConnectTimeout         customfield.NestedObject[CDNRuleTemplateOptionsProxyConnectTimeoutDataSourceModel]         `tfsdk:"proxy_connect_timeout" json:"proxy_connect_timeout,computed"`
	ProxyReadTimeout            customfield.NestedObject[CDNRuleTemplateOptionsProxyReadTimeoutDataSourceModel]            `tfsdk:"proxy_read_timeout" json:"proxy_read_timeout,computed"`
	QueryParamsBlacklist        customfield.NestedObject[CDNRuleTemplateOptionsQueryParamsBlacklistDataSourceModel]        `tfsdk:"query_params_blacklist" json:"query_params_blacklist,computed"`
	QueryParamsWhitelist        customfield.NestedObject[CDNRuleTemplateOptionsQueryParamsWhitelistDataSourceModel]        `tfsdk:"query_params_whitelist" json:"query_params_whitelist,computed"`
	QueryStringForwarding       customfield.NestedObject[CDNRuleTemplateOptionsQueryStringForwardingDataSourceModel]       `tfsdk:"query_string_forwarding" json:"query_string_forwarding,computed"`
	RedirectHTTPToHTTPS         customfield.NestedObject[CDNRuleTemplateOptionsRedirectHTTPToHTTPSDataSourceModel]         `tfsdk:"redirect_http_to_https" json:"redirect_http_to_https,computed"`
	RedirectHTTPSToHTTP         customfield.NestedObject[CDNRuleTemplateOptionsRedirectHTTPSToHTTPDataSourceModel]         `tfsdk:"redirect_https_to_http" json:"redirect_https_to_http,computed"`
	ReferrerACL                 customfield.NestedObject[CDNRuleTemplateOptionsReferrerACLDataSourceModel]                 `tfsdk:"referrer_acl" json:"referrer_acl,computed"`
	RequestLimiter              customfield.NestedObject[CDNRuleTemplateOptionsRequestLimiterDataSourceModel]              `tfsdk:"request_limiter" json:"request_limiter,computed"`
	ResponseHeadersHidingPolicy customfield.NestedObject[CDNRuleTemplateOptionsResponseHeadersHidingPolicyDataSourceModel] `tfsdk:"response_headers_hiding_policy" json:"response_headers_hiding_policy,computed"`
	Rewrite                     customfield.NestedObject[CDNRuleTemplateOptionsRewriteDataSourceModel]                     `tfsdk:"rewrite" json:"rewrite,computed"`
	SecureKey                   customfield.NestedObject[CDNRuleTemplateOptionsSecureKeyDataSourceModel]                   `tfsdk:"secure_key" json:"secure_key,computed"`
	Slice                       customfield.NestedObject[CDNRuleTemplateOptionsSliceDataSourceModel]                       `tfsdk:"slice" json:"slice,computed"`
	Sni                         customfield.NestedObject[CDNRuleTemplateOptionsSniDataSourceModel]                         `tfsdk:"sni" json:"sni,computed"`
	Stale                       customfield.NestedObject[CDNRuleTemplateOptionsStaleDataSourceModel]                       `tfsdk:"stale" json:"stale,computed"`
	StaticResponseHeaders       customfield.NestedObject[CDNRuleTemplateOptionsStaticResponseHeadersDataSourceModel]       `tfsdk:"static_response_headers" json:"static_response_headers,computed"`
	StaticHeaders               customfield.NestedObject[CDNRuleTemplateOptionsStaticHeadersDataSourceModel]               `tfsdk:"static_headers" json:"staticHeaders,computed"`
	StaticRequestHeaders        customfield.NestedObject[CDNRuleTemplateOptionsStaticRequestHeadersDataSourceModel]        `tfsdk:"static_request_headers" json:"staticRequestHeaders,computed"`
	UserAgentACL                customfield.NestedObject[CDNRuleTemplateOptionsUserAgentACLDataSourceModel]                `tfsdk:"user_agent_acl" json:"user_agent_acl,computed"`
	Waap                        customfield.NestedObject[CDNRuleTemplateOptionsWaapDataSourceModel]                        `tfsdk:"waap" json:"waap,computed"`
	Websockets                  customfield.NestedObject[CDNRuleTemplateOptionsWebsocketsDataSourceModel]                  `tfsdk:"websockets" json:"websockets,computed"`
}

type CDNRuleTemplateOptionsAllowedHTTPMethodsDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsBotProtectionDataSourceModel struct {
	BotChallenge customfield.NestedObject[CDNRuleTemplateOptionsBotProtectionBotChallengeDataSourceModel] `tfsdk:"bot_challenge" json:"bot_challenge,computed"`
	Enabled      types.Bool                                                                               `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNRuleTemplateOptionsBotProtectionBotChallengeDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNRuleTemplateOptionsBrotliCompressionDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsBrowserCacheSettingsDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsCacheHTTPHeadersDataSourceModel struct {
	Enabled types.Bool                     `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsCorsDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
	Always  types.Bool                    `tfsdk:"always" json:"always,computed"`
}

type CDNRuleTemplateOptionsCountryACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNRuleTemplateOptionsDisableCacheDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsDisableProxyForceRangesDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsEdgeCacheSettingsDataSourceModel struct {
	Enabled      types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	CustomValues customfield.Map[types.String] `tfsdk:"custom_values" json:"custom_values,computed"`
	Default      types.String                  `tfsdk:"default" json:"default,computed"`
	Value        types.String                  `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsFastedgeDataSourceModel struct {
	Enabled           types.Bool                                                                               `tfsdk:"enabled" json:"enabled,computed"`
	OnRequestBody     customfield.NestedObject[CDNRuleTemplateOptionsFastedgeOnRequestBodyDataSourceModel]     `tfsdk:"on_request_body" json:"on_request_body,computed"`
	OnRequestHeaders  customfield.NestedObject[CDNRuleTemplateOptionsFastedgeOnRequestHeadersDataSourceModel]  `tfsdk:"on_request_headers" json:"on_request_headers,computed"`
	OnResponseBody    customfield.NestedObject[CDNRuleTemplateOptionsFastedgeOnResponseBodyDataSourceModel]    `tfsdk:"on_response_body" json:"on_response_body,computed"`
	OnResponseHeaders customfield.NestedObject[CDNRuleTemplateOptionsFastedgeOnResponseHeadersDataSourceModel] `tfsdk:"on_response_headers" json:"on_response_headers,computed"`
}

type CDNRuleTemplateOptionsFastedgeOnRequestBodyDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNRuleTemplateOptionsFastedgeOnRequestHeadersDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNRuleTemplateOptionsFastedgeOnResponseBodyDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNRuleTemplateOptionsFastedgeOnResponseHeadersDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNRuleTemplateOptionsFetchCompressedDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsFollowOriginRedirectDataSourceModel struct {
	Codes   customfield.Set[types.Int64] `tfsdk:"codes" json:"codes,computed"`
	Enabled types.Bool                   `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNRuleTemplateOptionsForceReturnDataSourceModel struct {
	Body         types.String                                                                           `tfsdk:"body" json:"body,computed"`
	Code         types.Int64                                                                            `tfsdk:"code" json:"code,computed"`
	Enabled      types.Bool                                                                             `tfsdk:"enabled" json:"enabled,computed"`
	TimeInterval customfield.NestedObject[CDNRuleTemplateOptionsForceReturnTimeIntervalDataSourceModel] `tfsdk:"time_interval" json:"time_interval,computed"`
}

type CDNRuleTemplateOptionsForceReturnTimeIntervalDataSourceModel struct {
	EndTime   types.String `tfsdk:"end_time" json:"end_time,computed"`
	StartTime types.String `tfsdk:"start_time" json:"start_time,computed"`
	TimeZone  types.String `tfsdk:"time_zone" json:"time_zone,computed"`
}

type CDNRuleTemplateOptionsForwardHostHeaderDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsGzipOnDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsHostHeaderDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsIgnoreCookieDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsIgnoreQueryStringDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsImageStackDataSourceModel struct {
	Enabled     types.Bool  `tfsdk:"enabled" json:"enabled,computed"`
	AvifEnabled types.Bool  `tfsdk:"avif_enabled" json:"avif_enabled,computed"`
	PngLossless types.Bool  `tfsdk:"png_lossless" json:"png_lossless,computed"`
	Quality     types.Int64 `tfsdk:"quality" json:"quality,computed"`
	WebpEnabled types.Bool  `tfsdk:"webp_enabled" json:"webp_enabled,computed"`
}

type CDNRuleTemplateOptionsIPAddressACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNRuleTemplateOptionsLimitBandwidthDataSourceModel struct {
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	LimitType types.String `tfsdk:"limit_type" json:"limit_type,computed"`
	Buffer    types.Int64  `tfsdk:"buffer" json:"buffer,computed"`
	Speed     types.Int64  `tfsdk:"speed" json:"speed,computed"`
}

type CDNRuleTemplateOptionsProxyCacheKeyDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsProxyCacheMethodsSetDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsProxyConnectTimeoutDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsProxyReadTimeoutDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsQueryParamsBlacklistDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsQueryParamsWhitelistDataSourceModel struct {
	Enabled types.Bool                     `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsQueryStringForwardingDataSourceModel struct {
	Enabled              types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ForwardFromFileTypes customfield.Set[types.String] `tfsdk:"forward_from_file_types" json:"forward_from_file_types,computed"`
	ForwardToFileTypes   customfield.Set[types.String] `tfsdk:"forward_to_file_types" json:"forward_to_file_types,computed"`
	ForwardExceptKeys    customfield.Set[types.String] `tfsdk:"forward_except_keys" json:"forward_except_keys,computed"`
	ForwardOnlyKeys      customfield.Set[types.String] `tfsdk:"forward_only_keys" json:"forward_only_keys,computed"`
}

type CDNRuleTemplateOptionsRedirectHTTPToHTTPSDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsRedirectHTTPSToHTTPDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsReferrerACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNRuleTemplateOptionsRequestLimiterDataSourceModel struct {
	Enabled  types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Rate     types.Int64  `tfsdk:"rate" json:"rate,computed"`
	Burst    types.Int64  `tfsdk:"burst" json:"burst,computed"`
	Delay    types.Int64  `tfsdk:"delay" json:"delay,computed"`
	RateUnit types.String `tfsdk:"rate_unit" json:"rate_unit,computed"`
}

type CDNRuleTemplateOptionsResponseHeadersHidingPolicyDataSourceModel struct {
	Enabled  types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Excepted customfield.Set[types.String] `tfsdk:"excepted" json:"excepted,computed"`
	Mode     types.String                  `tfsdk:"mode" json:"mode,computed"`
}

type CDNRuleTemplateOptionsRewriteDataSourceModel struct {
	Body    types.String `tfsdk:"body" json:"body,computed"`
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Flag    types.String `tfsdk:"flag" json:"flag,computed"`
}

type CDNRuleTemplateOptionsSecureKeyDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Key     types.String `tfsdk:"key" json:"key,computed"`
	Type    types.Int64  `tfsdk:"type" json:"type,computed"`
}

type CDNRuleTemplateOptionsSliceDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsSniDataSourceModel struct {
	CustomHostname types.String `tfsdk:"custom_hostname" json:"custom_hostname,computed"`
	Enabled        types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	SniType        types.String `tfsdk:"sni_type" json:"sni_type,computed"`
}

type CDNRuleTemplateOptionsStaleDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsStaticResponseHeadersDataSourceModel struct {
	Enabled types.Bool                                                                                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.NestedObjectList[CDNRuleTemplateOptionsStaticResponseHeadersValueDataSourceModel] `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsStaticResponseHeadersValueDataSourceModel struct {
	Name   types.String                   `tfsdk:"name" json:"name,computed"`
	Value  customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
	Always types.Bool                     `tfsdk:"always" json:"always,computed"`
}

type CDNRuleTemplateOptionsStaticHeadersDataSourceModel struct {
	Enabled types.Bool           `tfsdk:"enabled" json:"enabled,computed"`
	Value   jsontypes.Normalized `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsStaticRequestHeadersDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Map[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsUserAgentACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNRuleTemplateOptionsWaapDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNRuleTemplateOptionsWebsocketsDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

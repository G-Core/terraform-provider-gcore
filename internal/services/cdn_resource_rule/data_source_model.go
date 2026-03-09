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

type CDNResourceRuleDataSourceModel struct {
	ID                     types.Int64                                                     `tfsdk:"id" path:"rule_id,computed"`
	RuleID                 types.Int64                                                     `tfsdk:"rule_id" path:"rule_id,required"`
	ResourceID             types.Int64                                                     `tfsdk:"resource_id" path:"resource_id,required"`
	Active                 types.Bool                                                      `tfsdk:"active" json:"active,computed"`
	Deleted                types.Bool                                                      `tfsdk:"deleted" json:"deleted,computed"`
	Name                   types.String                                                    `tfsdk:"name" json:"name,computed"`
	OriginGroup            types.Int64                                                     `tfsdk:"origin_group" json:"originGroup,computed"`
	OriginProtocol         types.String                                                    `tfsdk:"origin_protocol" json:"originProtocol,computed"`
	OverrideOriginProtocol types.String                                                    `tfsdk:"override_origin_protocol" json:"overrideOriginProtocol,computed"`
	PresetApplied          types.Bool                                                      `tfsdk:"preset_applied" json:"preset_applied,computed"`
	PrimaryRule            types.Int64                                                     `tfsdk:"primary_rule" json:"primary_rule,computed"`
	Rule                   types.String                                                    `tfsdk:"rule" json:"rule,computed"`
	RuleType               types.Int64                                                     `tfsdk:"rule_type" json:"ruleType,computed"`
	Weight                 types.Int64                                                     `tfsdk:"weight" json:"weight,computed"`
	Options                customfield.NestedObject[CDNResourceRuleOptionsDataSourceModel] `tfsdk:"options" json:"options,computed"`
}

func (m *CDNResourceRuleDataSourceModel) toReadParams(_ context.Context) (params cdn.CDNResourceRuleGetParams, diags diag.Diagnostics) {
	params = cdn.CDNResourceRuleGetParams{
		ResourceID: m.ResourceID.ValueInt64(),
	}

	return
}

type CDNResourceRuleOptionsDataSourceModel struct {
	AllowedHTTPMethods          customfield.NestedObject[CDNResourceRuleOptionsAllowedHTTPMethodsDataSourceModel]          `tfsdk:"allowed_http_methods" json:"allowedHttpMethods,computed"`
	BotProtection               customfield.NestedObject[CDNResourceRuleOptionsBotProtectionDataSourceModel]               `tfsdk:"bot_protection" json:"bot_protection,computed"`
	BrotliCompression           customfield.NestedObject[CDNResourceRuleOptionsBrotliCompressionDataSourceModel]           `tfsdk:"brotli_compression" json:"brotli_compression,computed"`
	BrowserCacheSettings        customfield.NestedObject[CDNResourceRuleOptionsBrowserCacheSettingsDataSourceModel]        `tfsdk:"browser_cache_settings" json:"browser_cache_settings,computed"`
	CacheHTTPHeaders            customfield.NestedObject[CDNResourceRuleOptionsCacheHTTPHeadersDataSourceModel]            `tfsdk:"cache_http_headers" json:"cache_http_headers,computed"`
	Cors                        customfield.NestedObject[CDNResourceRuleOptionsCorsDataSourceModel]                        `tfsdk:"cors" json:"cors,computed"`
	CountryACL                  customfield.NestedObject[CDNResourceRuleOptionsCountryACLDataSourceModel]                  `tfsdk:"country_acl" json:"country_acl,computed"`
	DisableCache                customfield.NestedObject[CDNResourceRuleOptionsDisableCacheDataSourceModel]                `tfsdk:"disable_cache" json:"disable_cache,computed"`
	DisableProxyForceRanges     customfield.NestedObject[CDNResourceRuleOptionsDisableProxyForceRangesDataSourceModel]     `tfsdk:"disable_proxy_force_ranges" json:"disable_proxy_force_ranges,computed"`
	EdgeCacheSettings           customfield.NestedObject[CDNResourceRuleOptionsEdgeCacheSettingsDataSourceModel]           `tfsdk:"edge_cache_settings" json:"edge_cache_settings,computed"`
	Fastedge                    customfield.NestedObject[CDNResourceRuleOptionsFastedgeDataSourceModel]                    `tfsdk:"fastedge" json:"fastedge,computed"`
	FetchCompressed             customfield.NestedObject[CDNResourceRuleOptionsFetchCompressedDataSourceModel]             `tfsdk:"fetch_compressed" json:"fetch_compressed,computed"`
	FollowOriginRedirect        customfield.NestedObject[CDNResourceRuleOptionsFollowOriginRedirectDataSourceModel]        `tfsdk:"follow_origin_redirect" json:"follow_origin_redirect,computed"`
	ForceReturn                 customfield.NestedObject[CDNResourceRuleOptionsForceReturnDataSourceModel]                 `tfsdk:"force_return" json:"force_return,computed"`
	ForwardHostHeader           customfield.NestedObject[CDNResourceRuleOptionsForwardHostHeaderDataSourceModel]           `tfsdk:"forward_host_header" json:"forward_host_header,computed"`
	GzipOn                      customfield.NestedObject[CDNResourceRuleOptionsGzipOnDataSourceModel]                      `tfsdk:"gzip_on" json:"gzipOn,computed"`
	HostHeader                  customfield.NestedObject[CDNResourceRuleOptionsHostHeaderDataSourceModel]                  `tfsdk:"host_header" json:"hostHeader,computed"`
	IgnoreCookie                customfield.NestedObject[CDNResourceRuleOptionsIgnoreCookieDataSourceModel]                `tfsdk:"ignore_cookie" json:"ignore_cookie,computed"`
	IgnoreQueryString           customfield.NestedObject[CDNResourceRuleOptionsIgnoreQueryStringDataSourceModel]           `tfsdk:"ignore_query_string" json:"ignoreQueryString,computed"`
	ImageStack                  customfield.NestedObject[CDNResourceRuleOptionsImageStackDataSourceModel]                  `tfsdk:"image_stack" json:"image_stack,computed"`
	IPAddressACL                customfield.NestedObject[CDNResourceRuleOptionsIPAddressACLDataSourceModel]                `tfsdk:"ip_address_acl" json:"ip_address_acl,computed"`
	LimitBandwidth              customfield.NestedObject[CDNResourceRuleOptionsLimitBandwidthDataSourceModel]              `tfsdk:"limit_bandwidth" json:"limit_bandwidth,computed"`
	ProxyCacheKey               customfield.NestedObject[CDNResourceRuleOptionsProxyCacheKeyDataSourceModel]               `tfsdk:"proxy_cache_key" json:"proxy_cache_key,computed"`
	ProxyCacheMethodsSet        customfield.NestedObject[CDNResourceRuleOptionsProxyCacheMethodsSetDataSourceModel]        `tfsdk:"proxy_cache_methods_set" json:"proxy_cache_methods_set,computed"`
	ProxyConnectTimeout         customfield.NestedObject[CDNResourceRuleOptionsProxyConnectTimeoutDataSourceModel]         `tfsdk:"proxy_connect_timeout" json:"proxy_connect_timeout,computed"`
	ProxyReadTimeout            customfield.NestedObject[CDNResourceRuleOptionsProxyReadTimeoutDataSourceModel]            `tfsdk:"proxy_read_timeout" json:"proxy_read_timeout,computed"`
	QueryParamsBlacklist        customfield.NestedObject[CDNResourceRuleOptionsQueryParamsBlacklistDataSourceModel]        `tfsdk:"query_params_blacklist" json:"query_params_blacklist,computed"`
	QueryParamsWhitelist        customfield.NestedObject[CDNResourceRuleOptionsQueryParamsWhitelistDataSourceModel]        `tfsdk:"query_params_whitelist" json:"query_params_whitelist,computed"`
	QueryStringForwarding       customfield.NestedObject[CDNResourceRuleOptionsQueryStringForwardingDataSourceModel]       `tfsdk:"query_string_forwarding" json:"query_string_forwarding,computed"`
	RedirectHTTPToHTTPS         customfield.NestedObject[CDNResourceRuleOptionsRedirectHTTPToHTTPSDataSourceModel]         `tfsdk:"redirect_http_to_https" json:"redirect_http_to_https,computed"`
	RedirectHTTPSToHTTP         customfield.NestedObject[CDNResourceRuleOptionsRedirectHTTPSToHTTPDataSourceModel]         `tfsdk:"redirect_https_to_http" json:"redirect_https_to_http,computed"`
	ReferrerACL                 customfield.NestedObject[CDNResourceRuleOptionsReferrerACLDataSourceModel]                 `tfsdk:"referrer_acl" json:"referrer_acl,computed"`
	RequestLimiter              customfield.NestedObject[CDNResourceRuleOptionsRequestLimiterDataSourceModel]              `tfsdk:"request_limiter" json:"request_limiter,computed"`
	ResponseHeadersHidingPolicy customfield.NestedObject[CDNResourceRuleOptionsResponseHeadersHidingPolicyDataSourceModel] `tfsdk:"response_headers_hiding_policy" json:"response_headers_hiding_policy,computed"`
	Rewrite                     customfield.NestedObject[CDNResourceRuleOptionsRewriteDataSourceModel]                     `tfsdk:"rewrite" json:"rewrite,computed"`
	SecureKey                   customfield.NestedObject[CDNResourceRuleOptionsSecureKeyDataSourceModel]                   `tfsdk:"secure_key" json:"secure_key,computed"`
	Slice                       customfield.NestedObject[CDNResourceRuleOptionsSliceDataSourceModel]                       `tfsdk:"slice" json:"slice,computed"`
	Sni                         customfield.NestedObject[CDNResourceRuleOptionsSniDataSourceModel]                         `tfsdk:"sni" json:"sni,computed"`
	Stale                       customfield.NestedObject[CDNResourceRuleOptionsStaleDataSourceModel]                       `tfsdk:"stale" json:"stale,computed"`
	StaticResponseHeaders       customfield.NestedObject[CDNResourceRuleOptionsStaticResponseHeadersDataSourceModel]       `tfsdk:"static_response_headers" json:"static_response_headers,computed"`
	StaticHeaders               customfield.NestedObject[CDNResourceRuleOptionsStaticHeadersDataSourceModel]               `tfsdk:"static_headers" json:"staticHeaders,computed"`
	StaticRequestHeaders        customfield.NestedObject[CDNResourceRuleOptionsStaticRequestHeadersDataSourceModel]        `tfsdk:"static_request_headers" json:"staticRequestHeaders,computed"`
	UserAgentACL                customfield.NestedObject[CDNResourceRuleOptionsUserAgentACLDataSourceModel]                `tfsdk:"user_agent_acl" json:"user_agent_acl,computed"`
	Waap                        customfield.NestedObject[CDNResourceRuleOptionsWaapDataSourceModel]                        `tfsdk:"waap" json:"waap,computed"`
	Websockets                  customfield.NestedObject[CDNResourceRuleOptionsWebsocketsDataSourceModel]                  `tfsdk:"websockets" json:"websockets,computed"`
}

type CDNResourceRuleOptionsAllowedHTTPMethodsDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsBotProtectionDataSourceModel struct {
	BotChallenge customfield.NestedObject[CDNResourceRuleOptionsBotProtectionBotChallengeDataSourceModel] `tfsdk:"bot_challenge" json:"bot_challenge,computed"`
	Enabled      types.Bool                                                                               `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNResourceRuleOptionsBotProtectionBotChallengeDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNResourceRuleOptionsBrotliCompressionDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsBrowserCacheSettingsDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsCacheHTTPHeadersDataSourceModel struct {
	Enabled types.Bool                     `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsCorsDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
	Always  types.Bool                    `tfsdk:"always" json:"always,computed"`
}

type CDNResourceRuleOptionsCountryACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourceRuleOptionsDisableCacheDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsDisableProxyForceRangesDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsEdgeCacheSettingsDataSourceModel struct {
	Enabled      types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	CustomValues customfield.Map[types.String] `tfsdk:"custom_values" json:"custom_values,computed"`
	Default      types.String                  `tfsdk:"default" json:"default,computed"`
	Value        types.String                  `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsFastedgeDataSourceModel struct {
	Enabled           types.Bool                                                                               `tfsdk:"enabled" json:"enabled,computed"`
	OnRequestBody     customfield.NestedObject[CDNResourceRuleOptionsFastedgeOnRequestBodyDataSourceModel]     `tfsdk:"on_request_body" json:"on_request_body,computed"`
	OnRequestHeaders  customfield.NestedObject[CDNResourceRuleOptionsFastedgeOnRequestHeadersDataSourceModel]  `tfsdk:"on_request_headers" json:"on_request_headers,computed"`
	OnResponseBody    customfield.NestedObject[CDNResourceRuleOptionsFastedgeOnResponseBodyDataSourceModel]    `tfsdk:"on_response_body" json:"on_response_body,computed"`
	OnResponseHeaders customfield.NestedObject[CDNResourceRuleOptionsFastedgeOnResponseHeadersDataSourceModel] `tfsdk:"on_response_headers" json:"on_response_headers,computed"`
}

type CDNResourceRuleOptionsFastedgeOnRequestBodyDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourceRuleOptionsFastedgeOnRequestHeadersDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourceRuleOptionsFastedgeOnResponseBodyDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourceRuleOptionsFastedgeOnResponseHeadersDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourceRuleOptionsFetchCompressedDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsFollowOriginRedirectDataSourceModel struct {
	Codes   customfield.Set[types.Int64] `tfsdk:"codes" json:"codes,computed"`
	Enabled types.Bool                   `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNResourceRuleOptionsForceReturnDataSourceModel struct {
	Body         types.String                                                                           `tfsdk:"body" json:"body,computed"`
	Code         types.Int64                                                                            `tfsdk:"code" json:"code,computed"`
	Enabled      types.Bool                                                                             `tfsdk:"enabled" json:"enabled,computed"`
	TimeInterval customfield.NestedObject[CDNResourceRuleOptionsForceReturnTimeIntervalDataSourceModel] `tfsdk:"time_interval" json:"time_interval,computed"`
}

type CDNResourceRuleOptionsForceReturnTimeIntervalDataSourceModel struct {
	EndTime   types.String `tfsdk:"end_time" json:"end_time,computed"`
	StartTime types.String `tfsdk:"start_time" json:"start_time,computed"`
	TimeZone  types.String `tfsdk:"time_zone" json:"time_zone,computed"`
}

type CDNResourceRuleOptionsForwardHostHeaderDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsGzipOnDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsHostHeaderDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsIgnoreCookieDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsIgnoreQueryStringDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsImageStackDataSourceModel struct {
	Enabled     types.Bool  `tfsdk:"enabled" json:"enabled,computed"`
	AvifEnabled types.Bool  `tfsdk:"avif_enabled" json:"avif_enabled,computed"`
	PngLossless types.Bool  `tfsdk:"png_lossless" json:"png_lossless,computed"`
	Quality     types.Int64 `tfsdk:"quality" json:"quality,computed"`
	WebpEnabled types.Bool  `tfsdk:"webp_enabled" json:"webp_enabled,computed"`
}

type CDNResourceRuleOptionsIPAddressACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourceRuleOptionsLimitBandwidthDataSourceModel struct {
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	LimitType types.String `tfsdk:"limit_type" json:"limit_type,computed"`
	Buffer    types.Int64  `tfsdk:"buffer" json:"buffer,computed"`
	Speed     types.Int64  `tfsdk:"speed" json:"speed,computed"`
}

type CDNResourceRuleOptionsProxyCacheKeyDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsProxyCacheMethodsSetDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsProxyConnectTimeoutDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsProxyReadTimeoutDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsQueryParamsBlacklistDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsQueryParamsWhitelistDataSourceModel struct {
	Enabled types.Bool                     `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsQueryStringForwardingDataSourceModel struct {
	Enabled              types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ForwardFromFileTypes customfield.Set[types.String] `tfsdk:"forward_from_file_types" json:"forward_from_file_types,computed"`
	ForwardToFileTypes   customfield.Set[types.String] `tfsdk:"forward_to_file_types" json:"forward_to_file_types,computed"`
	ForwardExceptKeys    customfield.Set[types.String] `tfsdk:"forward_except_keys" json:"forward_except_keys,computed"`
	ForwardOnlyKeys      customfield.Set[types.String] `tfsdk:"forward_only_keys" json:"forward_only_keys,computed"`
}

type CDNResourceRuleOptionsRedirectHTTPToHTTPSDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsRedirectHTTPSToHTTPDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsReferrerACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourceRuleOptionsRequestLimiterDataSourceModel struct {
	Enabled  types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Rate     types.Int64  `tfsdk:"rate" json:"rate,computed"`
	Burst    types.Int64  `tfsdk:"burst" json:"burst,computed"`
	Delay    types.Int64  `tfsdk:"delay" json:"delay,computed"`
	RateUnit types.String `tfsdk:"rate_unit" json:"rate_unit,computed"`
}

type CDNResourceRuleOptionsResponseHeadersHidingPolicyDataSourceModel struct {
	Enabled  types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Excepted customfield.Set[types.String] `tfsdk:"excepted" json:"excepted,computed"`
	Mode     types.String                  `tfsdk:"mode" json:"mode,computed"`
}

type CDNResourceRuleOptionsRewriteDataSourceModel struct {
	Body    types.String `tfsdk:"body" json:"body,computed"`
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Flag    types.String `tfsdk:"flag" json:"flag,computed"`
}

type CDNResourceRuleOptionsSecureKeyDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Key     types.String `tfsdk:"key" json:"key,computed"`
	Type    types.Int64  `tfsdk:"type" json:"type,computed"`
}

type CDNResourceRuleOptionsSliceDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsSniDataSourceModel struct {
	CustomHostname types.String `tfsdk:"custom_hostname" json:"custom_hostname,computed"`
	Enabled        types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	SniType        types.String `tfsdk:"sni_type" json:"sni_type,computed"`
}

type CDNResourceRuleOptionsStaleDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsStaticResponseHeadersDataSourceModel struct {
	Enabled types.Bool                                                                                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.NestedObjectList[CDNResourceRuleOptionsStaticResponseHeadersValueDataSourceModel] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsStaticResponseHeadersValueDataSourceModel struct {
	Name   types.String                   `tfsdk:"name" json:"name,computed"`
	Value  customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
	Always types.Bool                     `tfsdk:"always" json:"always,computed"`
}

type CDNResourceRuleOptionsStaticHeadersDataSourceModel struct {
	Enabled types.Bool           `tfsdk:"enabled" json:"enabled,computed"`
	Value   jsontypes.Normalized `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsStaticRequestHeadersDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Map[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsUserAgentACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourceRuleOptionsWaapDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceRuleOptionsWebsocketsDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

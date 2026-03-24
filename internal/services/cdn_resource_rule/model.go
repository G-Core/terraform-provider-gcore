// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_resource_rule

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNResourceRuleModel struct {
	ID                     types.Int64                                           `tfsdk:"id" json:"id,computed"`
	ResourceID             types.Int64                                           `tfsdk:"resource_id" path:"resource_id,required"`
	Name                   types.String                                          `tfsdk:"name" json:"name,required"`
	Rule                   types.String                                          `tfsdk:"rule" json:"rule,required"`
	RuleType               types.Int64                                           `tfsdk:"rule_type" json:"ruleType,required"`
	Active                 types.Bool                                            `tfsdk:"active" json:"active,optional"`
	OriginGroup            types.Int64                                           `tfsdk:"origin_group" json:"originGroup,optional"`
	OverrideOriginProtocol types.String                                          `tfsdk:"override_origin_protocol" json:"overrideOriginProtocol,optional"`
	Weight                 types.Int64                                           `tfsdk:"weight" json:"weight,computed_optional"`
	Options                customfield.NestedObject[CDNResourceRuleOptionsModel] `tfsdk:"options" json:"options,computed_optional"`
	Deleted                types.Bool                                            `tfsdk:"deleted" json:"deleted,computed"`
	OriginProtocol         types.String                                          `tfsdk:"origin_protocol" json:"originProtocol,computed"`
	PresetApplied          types.Bool                                            `tfsdk:"preset_applied" json:"preset_applied,computed"`
	PrimaryRule            types.Int64                                           `tfsdk:"primary_rule" json:"primary_rule,computed"`
}

func (m CDNResourceRuleModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CDNResourceRuleModel) MarshalJSONForUpdate(state CDNResourceRuleModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CDNResourceRuleOptionsModel struct {
	AllowedHTTPMethods          customfield.NestedObject[CDNResourceRuleOptionsAllowedHTTPMethodsModel]          `tfsdk:"allowed_http_methods" json:"allowedHttpMethods,computed_optional"`
	BotProtection               customfield.NestedObject[CDNResourceRuleOptionsBotProtectionModel]               `tfsdk:"bot_protection" json:"bot_protection,computed_optional"`
	BrotliCompression           customfield.NestedObject[CDNResourceRuleOptionsBrotliCompressionModel]           `tfsdk:"brotli_compression" json:"brotli_compression,computed_optional"`
	BrowserCacheSettings        customfield.NestedObject[CDNResourceRuleOptionsBrowserCacheSettingsModel]        `tfsdk:"browser_cache_settings" json:"browser_cache_settings,computed_optional"`
	CacheHTTPHeaders            customfield.NestedObject[CDNResourceRuleOptionsCacheHTTPHeadersModel]            `tfsdk:"cache_http_headers" json:"cache_http_headers,computed_optional"`
	Cors                        customfield.NestedObject[CDNResourceRuleOptionsCorsModel]                        `tfsdk:"cors" json:"cors,computed_optional"`
	CountryACL                  customfield.NestedObject[CDNResourceRuleOptionsCountryACLModel]                  `tfsdk:"country_acl" json:"country_acl,computed_optional"`
	DisableCache                customfield.NestedObject[CDNResourceRuleOptionsDisableCacheModel]                `tfsdk:"disable_cache" json:"disable_cache,computed_optional"`
	DisableProxyForceRanges     customfield.NestedObject[CDNResourceRuleOptionsDisableProxyForceRangesModel]     `tfsdk:"disable_proxy_force_ranges" json:"disable_proxy_force_ranges,computed_optional"`
	EdgeCacheSettings           customfield.NestedObject[CDNResourceRuleOptionsEdgeCacheSettingsModel]           `tfsdk:"edge_cache_settings" json:"edge_cache_settings,computed_optional"`
	Fastedge                    customfield.NestedObject[CDNResourceRuleOptionsFastedgeModel]                    `tfsdk:"fastedge" json:"fastedge,computed_optional"`
	FetchCompressed             customfield.NestedObject[CDNResourceRuleOptionsFetchCompressedModel]             `tfsdk:"fetch_compressed" json:"fetch_compressed,computed_optional"`
	FollowOriginRedirect        customfield.NestedObject[CDNResourceRuleOptionsFollowOriginRedirectModel]        `tfsdk:"follow_origin_redirect" json:"follow_origin_redirect,computed_optional"`
	ForceReturn                 customfield.NestedObject[CDNResourceRuleOptionsForceReturnModel]                 `tfsdk:"force_return" json:"force_return,computed_optional"`
	ForwardHostHeader           customfield.NestedObject[CDNResourceRuleOptionsForwardHostHeaderModel]           `tfsdk:"forward_host_header" json:"forward_host_header,computed_optional"`
	GzipOn                      customfield.NestedObject[CDNResourceRuleOptionsGzipOnModel]                      `tfsdk:"gzip_on" json:"gzipOn,computed_optional"`
	HostHeader                  customfield.NestedObject[CDNResourceRuleOptionsHostHeaderModel]                  `tfsdk:"host_header" json:"hostHeader,computed_optional"`
	IgnoreCookie                customfield.NestedObject[CDNResourceRuleOptionsIgnoreCookieModel]                `tfsdk:"ignore_cookie" json:"ignore_cookie,computed_optional"`
	IgnoreQueryString           customfield.NestedObject[CDNResourceRuleOptionsIgnoreQueryStringModel]           `tfsdk:"ignore_query_string" json:"ignoreQueryString,computed_optional"`
	ImageStack                  customfield.NestedObject[CDNResourceRuleOptionsImageStackModel]                  `tfsdk:"image_stack" json:"image_stack,computed_optional"`
	IPAddressACL                customfield.NestedObject[CDNResourceRuleOptionsIPAddressACLModel]                `tfsdk:"ip_address_acl" json:"ip_address_acl,computed_optional"`
	LimitBandwidth              customfield.NestedObject[CDNResourceRuleOptionsLimitBandwidthModel]              `tfsdk:"limit_bandwidth" json:"limit_bandwidth,computed_optional"`
	ProxyCacheKey               customfield.NestedObject[CDNResourceRuleOptionsProxyCacheKeyModel]               `tfsdk:"proxy_cache_key" json:"proxy_cache_key,computed_optional"`
	ProxyCacheMethodsSet        customfield.NestedObject[CDNResourceRuleOptionsProxyCacheMethodsSetModel]        `tfsdk:"proxy_cache_methods_set" json:"proxy_cache_methods_set,computed_optional"`
	ProxyConnectTimeout         customfield.NestedObject[CDNResourceRuleOptionsProxyConnectTimeoutModel]         `tfsdk:"proxy_connect_timeout" json:"proxy_connect_timeout,computed_optional"`
	ProxyReadTimeout            customfield.NestedObject[CDNResourceRuleOptionsProxyReadTimeoutModel]            `tfsdk:"proxy_read_timeout" json:"proxy_read_timeout,computed_optional"`
	QueryParamsBlacklist        customfield.NestedObject[CDNResourceRuleOptionsQueryParamsBlacklistModel]        `tfsdk:"query_params_blacklist" json:"query_params_blacklist,computed_optional"`
	QueryParamsWhitelist        customfield.NestedObject[CDNResourceRuleOptionsQueryParamsWhitelistModel]        `tfsdk:"query_params_whitelist" json:"query_params_whitelist,computed_optional"`
	QueryStringForwarding       customfield.NestedObject[CDNResourceRuleOptionsQueryStringForwardingModel]       `tfsdk:"query_string_forwarding" json:"query_string_forwarding,computed_optional"`
	RedirectHTTPToHTTPS         customfield.NestedObject[CDNResourceRuleOptionsRedirectHTTPToHTTPSModel]         `tfsdk:"redirect_http_to_https" json:"redirect_http_to_https,computed_optional"`
	RedirectHTTPSToHTTP         customfield.NestedObject[CDNResourceRuleOptionsRedirectHTTPSToHTTPModel]         `tfsdk:"redirect_https_to_http" json:"redirect_https_to_http,computed_optional"`
	ReferrerACL                 customfield.NestedObject[CDNResourceRuleOptionsReferrerACLModel]                 `tfsdk:"referrer_acl" json:"referrer_acl,computed_optional"`
	RequestLimiter              customfield.NestedObject[CDNResourceRuleOptionsRequestLimiterModel]              `tfsdk:"request_limiter" json:"request_limiter,computed_optional"`
	ResponseHeadersHidingPolicy customfield.NestedObject[CDNResourceRuleOptionsResponseHeadersHidingPolicyModel] `tfsdk:"response_headers_hiding_policy" json:"response_headers_hiding_policy,computed_optional"`
	Rewrite                     customfield.NestedObject[CDNResourceRuleOptionsRewriteModel]                     `tfsdk:"rewrite" json:"rewrite,computed_optional"`
	SecureKey                   customfield.NestedObject[CDNResourceRuleOptionsSecureKeyModel]                   `tfsdk:"secure_key" json:"secure_key,computed_optional"`
	Slice                       customfield.NestedObject[CDNResourceRuleOptionsSliceModel]                       `tfsdk:"slice" json:"slice,computed_optional"`
	Sni                         customfield.NestedObject[CDNResourceRuleOptionsSniModel]                         `tfsdk:"sni" json:"sni,computed_optional"`
	Stale                       customfield.NestedObject[CDNResourceRuleOptionsStaleModel]                       `tfsdk:"stale" json:"stale,computed_optional"`
	StaticResponseHeaders       customfield.NestedObject[CDNResourceRuleOptionsStaticResponseHeadersModel]       `tfsdk:"static_response_headers" json:"static_response_headers,computed_optional"`
	StaticHeaders               customfield.NestedObject[CDNResourceRuleOptionsStaticHeadersModel]               `tfsdk:"static_headers" json:"staticHeaders,computed_optional"`
	StaticRequestHeaders        customfield.NestedObject[CDNResourceRuleOptionsStaticRequestHeadersModel]        `tfsdk:"static_request_headers" json:"staticRequestHeaders,computed_optional"`
	UserAgentACL                customfield.NestedObject[CDNResourceRuleOptionsUserAgentACLModel]                `tfsdk:"user_agent_acl" json:"user_agent_acl,computed_optional"`
	Waap                        customfield.NestedObject[CDNResourceRuleOptionsWaapModel]                        `tfsdk:"waap" json:"waap,computed_optional"`
	Websockets                  customfield.NestedObject[CDNResourceRuleOptionsWebsocketsModel]                  `tfsdk:"websockets" json:"websockets,computed_optional"`
}

type CDNResourceRuleOptionsAllowedHTTPMethodsModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsBotProtectionModel struct {
	BotChallenge *CDNResourceRuleOptionsBotProtectionBotChallengeModel `tfsdk:"bot_challenge" json:"bot_challenge,required"`
	Enabled      types.Bool                                            `tfsdk:"enabled" json:"enabled,required"`
}

type CDNResourceRuleOptionsBotProtectionBotChallengeModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed_optional"`
}

type CDNResourceRuleOptionsBrotliCompressionModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsBrowserCacheSettingsModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsCacheHTTPHeadersModel struct {
	Enabled types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsCorsModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
	Always  types.Bool                    `tfsdk:"always" json:"always,computed_optional"`
}

type CDNResourceRuleOptionsCountryACLModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNResourceRuleOptionsDisableCacheModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsDisableProxyForceRangesModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsEdgeCacheSettingsModel struct {
	Enabled      types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	CustomValues customfield.Map[types.String] `tfsdk:"custom_values" json:"custom_values,computed_optional"`
	Default      types.String                  `tfsdk:"default" json:"default,optional"`
	Value        types.String                  `tfsdk:"value" json:"value,optional"`
}

type CDNResourceRuleOptionsFastedgeModel struct {
	Enabled           types.Bool                                                                     `tfsdk:"enabled" json:"enabled,required"`
	OnRequestBody     customfield.NestedObject[CDNResourceRuleOptionsFastedgeOnRequestBodyModel]     `tfsdk:"on_request_body" json:"on_request_body,computed_optional"`
	OnRequestHeaders  customfield.NestedObject[CDNResourceRuleOptionsFastedgeOnRequestHeadersModel]  `tfsdk:"on_request_headers" json:"on_request_headers,computed_optional"`
	OnResponseBody    customfield.NestedObject[CDNResourceRuleOptionsFastedgeOnResponseBodyModel]    `tfsdk:"on_response_body" json:"on_response_body,computed_optional"`
	OnResponseHeaders customfield.NestedObject[CDNResourceRuleOptionsFastedgeOnResponseHeadersModel] `tfsdk:"on_response_headers" json:"on_response_headers,computed_optional"`
}

type CDNResourceRuleOptionsFastedgeOnRequestBodyModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNResourceRuleOptionsFastedgeOnRequestHeadersModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNResourceRuleOptionsFastedgeOnResponseBodyModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNResourceRuleOptionsFastedgeOnResponseHeadersModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNResourceRuleOptionsFetchCompressedModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsFollowOriginRedirectModel struct {
	Codes   customfield.Set[types.Int64] `tfsdk:"codes" json:"codes,required"`
	Enabled types.Bool                   `tfsdk:"enabled" json:"enabled,required"`
}

type CDNResourceRuleOptionsForceReturnModel struct {
	Body         types.String                                                                 `tfsdk:"body" json:"body,required"`
	Code         types.Int64                                                                  `tfsdk:"code" json:"code,required"`
	Enabled      types.Bool                                                                   `tfsdk:"enabled" json:"enabled,required"`
	TimeInterval customfield.NestedObject[CDNResourceRuleOptionsForceReturnTimeIntervalModel] `tfsdk:"time_interval" json:"time_interval,computed_optional"`
}

type CDNResourceRuleOptionsForceReturnTimeIntervalModel struct {
	EndTime   types.String `tfsdk:"end_time" json:"end_time,required"`
	StartTime types.String `tfsdk:"start_time" json:"start_time,required"`
	TimeZone  types.String `tfsdk:"time_zone" json:"time_zone,computed_optional"`
}

type CDNResourceRuleOptionsForwardHostHeaderModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsGzipOnModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsHostHeaderModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsIgnoreCookieModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsIgnoreQueryStringModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsImageStackModel struct {
	Enabled     types.Bool  `tfsdk:"enabled" json:"enabled,required"`
	AvifEnabled types.Bool  `tfsdk:"avif_enabled" json:"avif_enabled,computed_optional"`
	PngLossless types.Bool  `tfsdk:"png_lossless" json:"png_lossless,computed_optional"`
	Quality     types.Int64 `tfsdk:"quality" json:"quality,computed_optional"`
	WebpEnabled types.Bool  `tfsdk:"webp_enabled" json:"webp_enabled,computed_optional"`
}

type CDNResourceRuleOptionsIPAddressACLModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNResourceRuleOptionsLimitBandwidthModel struct {
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	LimitType types.String `tfsdk:"limit_type" json:"limit_type,required"`
	Buffer    types.Int64  `tfsdk:"buffer" json:"buffer,optional"`
	Speed     types.Int64  `tfsdk:"speed" json:"speed,optional"`
}

type CDNResourceRuleOptionsProxyCacheKeyModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsProxyCacheMethodsSetModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsProxyConnectTimeoutModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsProxyReadTimeoutModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsQueryParamsBlacklistModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsQueryParamsWhitelistModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsQueryStringForwardingModel struct {
	Enabled              types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ForwardFromFileTypes customfield.Set[types.String] `tfsdk:"forward_from_file_types" json:"forward_from_file_types,required"`
	ForwardToFileTypes   customfield.Set[types.String] `tfsdk:"forward_to_file_types" json:"forward_to_file_types,required"`
	ForwardExceptKeys    customfield.Set[types.String] `tfsdk:"forward_except_keys" json:"forward_except_keys,optional"`
	ForwardOnlyKeys      customfield.Set[types.String] `tfsdk:"forward_only_keys" json:"forward_only_keys,optional"`
}

type CDNResourceRuleOptionsRedirectHTTPToHTTPSModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsRedirectHTTPSToHTTPModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsReferrerACLModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNResourceRuleOptionsRequestLimiterModel struct {
	Enabled  types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Rate     types.Int64  `tfsdk:"rate" json:"rate,required"`
	Burst    types.Int64  `tfsdk:"burst" json:"burst,computed"`
	Delay    types.Int64  `tfsdk:"delay" json:"delay,computed"`
	RateUnit types.String `tfsdk:"rate_unit" json:"rate_unit,computed_optional"`
}

type CDNResourceRuleOptionsResponseHeadersHidingPolicyModel struct {
	Enabled  types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Excepted customfield.Set[types.String] `tfsdk:"excepted" json:"excepted,computed_optional"`
	Mode     types.String                  `tfsdk:"mode" json:"mode,required"`
}

type CDNResourceRuleOptionsRewriteModel struct {
	Body    types.String `tfsdk:"body" json:"body,required"`
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Flag    types.String `tfsdk:"flag" json:"flag,computed_optional"`
}

type CDNResourceRuleOptionsSecureKeyModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Key     types.String `tfsdk:"key" json:"key,required"`
	Type    types.Int64  `tfsdk:"type" json:"type,computed_optional"`
}

type CDNResourceRuleOptionsSliceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsSniModel struct {
	CustomHostname types.String `tfsdk:"custom_hostname" json:"custom_hostname,required"`
	Enabled        types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	SniType        types.String `tfsdk:"sni_type" json:"sni_type,computed_optional"`
}

type CDNResourceRuleOptionsStaleModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsStaticResponseHeadersModel struct {
	Enabled types.Bool                                                `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]*CDNResourceRuleOptionsStaticResponseHeadersValueModel `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsStaticResponseHeadersValueModel struct {
	Name   types.String    `tfsdk:"name" json:"name,required"`
	Value  *[]types.String `tfsdk:"value" json:"value,required"`
	Always types.Bool      `tfsdk:"always" json:"always,computed_optional"`
}

type CDNResourceRuleOptionsStaticHeadersModel struct {
	Enabled types.Bool           `tfsdk:"enabled" json:"enabled,required"`
	Value   jsontypes.Normalized `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsStaticRequestHeadersModel struct {
	Enabled types.Bool               `tfsdk:"enabled" json:"enabled,required"`
	Value   *map[string]types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsUserAgentACLModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNResourceRuleOptionsWaapModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceRuleOptionsWebsocketsModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

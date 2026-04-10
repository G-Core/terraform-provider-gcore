// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_rule_template

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNRuleTemplateModel struct {
	ID                     types.Int64                                           `tfsdk:"id" json:"id,computed"`
	Rule                   types.String                                          `tfsdk:"rule" json:"rule,required"`
	RuleType               types.Int64                                           `tfsdk:"rule_type" json:"ruleType,required"`
	Name                   types.String                                          `tfsdk:"name" json:"name,optional"`
	OverrideOriginProtocol types.String                                          `tfsdk:"override_origin_protocol" json:"overrideOriginProtocol,optional"`
	Weight                 types.Int64                                           `tfsdk:"weight" json:"weight,optional"`
	Options                customfield.NestedObject[CDNRuleTemplateOptionsModel] `tfsdk:"options" json:"options,computed_optional"`
	Client                 types.Int64                                           `tfsdk:"client" json:"client,computed"`
	Default                types.Bool                                            `tfsdk:"default" json:"default,computed"`
	Deleted                types.Bool                                            `tfsdk:"deleted" json:"deleted,computed,no_refresh"`
	Template               types.Bool                                            `tfsdk:"template" json:"template,computed"`
}

func (m CDNRuleTemplateModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CDNRuleTemplateModel) MarshalJSONForUpdate(state CDNRuleTemplateModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CDNRuleTemplateOptionsModel struct {
	AllowedHTTPMethods          customfield.NestedObject[CDNRuleTemplateOptionsAllowedHTTPMethodsModel]          `tfsdk:"allowed_http_methods" json:"allowedHttpMethods,computed_optional"`
	BotProtection               customfield.NestedObject[CDNRuleTemplateOptionsBotProtectionModel]               `tfsdk:"bot_protection" json:"bot_protection,computed_optional"`
	BrotliCompression           customfield.NestedObject[CDNRuleTemplateOptionsBrotliCompressionModel]           `tfsdk:"brotli_compression" json:"brotli_compression,computed_optional"`
	BrowserCacheSettings        customfield.NestedObject[CDNRuleTemplateOptionsBrowserCacheSettingsModel]        `tfsdk:"browser_cache_settings" json:"browser_cache_settings,computed_optional"`
	CacheHTTPHeaders            customfield.NestedObject[CDNRuleTemplateOptionsCacheHTTPHeadersModel]            `tfsdk:"cache_http_headers" json:"cache_http_headers,computed_optional"`
	Cors                        customfield.NestedObject[CDNRuleTemplateOptionsCorsModel]                        `tfsdk:"cors" json:"cors,computed_optional"`
	CountryACL                  customfield.NestedObject[CDNRuleTemplateOptionsCountryACLModel]                  `tfsdk:"country_acl" json:"country_acl,computed_optional"`
	DisableCache                customfield.NestedObject[CDNRuleTemplateOptionsDisableCacheModel]                `tfsdk:"disable_cache" json:"disable_cache,computed_optional"`
	DisableProxyForceRanges     customfield.NestedObject[CDNRuleTemplateOptionsDisableProxyForceRangesModel]     `tfsdk:"disable_proxy_force_ranges" json:"disable_proxy_force_ranges,computed_optional"`
	EdgeCacheSettings           customfield.NestedObject[CDNRuleTemplateOptionsEdgeCacheSettingsModel]           `tfsdk:"edge_cache_settings" json:"edge_cache_settings,computed_optional"`
	Fastedge                    customfield.NestedObject[CDNRuleTemplateOptionsFastedgeModel]                    `tfsdk:"fastedge" json:"fastedge,computed_optional"`
	FetchCompressed             customfield.NestedObject[CDNRuleTemplateOptionsFetchCompressedModel]             `tfsdk:"fetch_compressed" json:"fetch_compressed,computed_optional"`
	FollowOriginRedirect        customfield.NestedObject[CDNRuleTemplateOptionsFollowOriginRedirectModel]        `tfsdk:"follow_origin_redirect" json:"follow_origin_redirect,computed_optional"`
	ForceReturn                 customfield.NestedObject[CDNRuleTemplateOptionsForceReturnModel]                 `tfsdk:"force_return" json:"force_return,computed_optional"`
	ForwardHostHeader           customfield.NestedObject[CDNRuleTemplateOptionsForwardHostHeaderModel]           `tfsdk:"forward_host_header" json:"forward_host_header,computed_optional"`
	GzipOn                      customfield.NestedObject[CDNRuleTemplateOptionsGzipOnModel]                      `tfsdk:"gzip_on" json:"gzipOn,computed_optional"`
	HostHeader                  customfield.NestedObject[CDNRuleTemplateOptionsHostHeaderModel]                  `tfsdk:"host_header" json:"hostHeader,computed_optional"`
	IgnoreCookie                customfield.NestedObject[CDNRuleTemplateOptionsIgnoreCookieModel]                `tfsdk:"ignore_cookie" json:"ignore_cookie,computed_optional"`
	IgnoreQueryString           customfield.NestedObject[CDNRuleTemplateOptionsIgnoreQueryStringModel]           `tfsdk:"ignore_query_string" json:"ignoreQueryString,computed_optional"`
	ImageStack                  customfield.NestedObject[CDNRuleTemplateOptionsImageStackModel]                  `tfsdk:"image_stack" json:"image_stack,computed_optional"`
	IPAddressACL                customfield.NestedObject[CDNRuleTemplateOptionsIPAddressACLModel]                `tfsdk:"ip_address_acl" json:"ip_address_acl,computed_optional"`
	LimitBandwidth              customfield.NestedObject[CDNRuleTemplateOptionsLimitBandwidthModel]              `tfsdk:"limit_bandwidth" json:"limit_bandwidth,computed_optional"`
	ProxyCacheKey               customfield.NestedObject[CDNRuleTemplateOptionsProxyCacheKeyModel]               `tfsdk:"proxy_cache_key" json:"proxy_cache_key,computed_optional"`
	ProxyCacheMethodsSet        customfield.NestedObject[CDNRuleTemplateOptionsProxyCacheMethodsSetModel]        `tfsdk:"proxy_cache_methods_set" json:"proxy_cache_methods_set,computed_optional"`
	ProxyConnectTimeout         customfield.NestedObject[CDNRuleTemplateOptionsProxyConnectTimeoutModel]         `tfsdk:"proxy_connect_timeout" json:"proxy_connect_timeout,computed_optional"`
	ProxyReadTimeout            customfield.NestedObject[CDNRuleTemplateOptionsProxyReadTimeoutModel]            `tfsdk:"proxy_read_timeout" json:"proxy_read_timeout,computed_optional"`
	QueryParamsBlacklist        customfield.NestedObject[CDNRuleTemplateOptionsQueryParamsBlacklistModel]        `tfsdk:"query_params_blacklist" json:"query_params_blacklist,computed_optional"`
	QueryParamsWhitelist        customfield.NestedObject[CDNRuleTemplateOptionsQueryParamsWhitelistModel]        `tfsdk:"query_params_whitelist" json:"query_params_whitelist,computed_optional"`
	QueryStringForwarding       customfield.NestedObject[CDNRuleTemplateOptionsQueryStringForwardingModel]       `tfsdk:"query_string_forwarding" json:"query_string_forwarding,computed_optional"`
	RedirectHTTPToHTTPS         customfield.NestedObject[CDNRuleTemplateOptionsRedirectHTTPToHTTPSModel]         `tfsdk:"redirect_http_to_https" json:"redirect_http_to_https,computed_optional"`
	RedirectHTTPSToHTTP         customfield.NestedObject[CDNRuleTemplateOptionsRedirectHTTPSToHTTPModel]         `tfsdk:"redirect_https_to_http" json:"redirect_https_to_http,computed_optional"`
	ReferrerACL                 customfield.NestedObject[CDNRuleTemplateOptionsReferrerACLModel]                 `tfsdk:"referrer_acl" json:"referrer_acl,computed_optional"`
	RequestLimiter              customfield.NestedObject[CDNRuleTemplateOptionsRequestLimiterModel]              `tfsdk:"request_limiter" json:"request_limiter,computed_optional"`
	ResponseHeadersHidingPolicy customfield.NestedObject[CDNRuleTemplateOptionsResponseHeadersHidingPolicyModel] `tfsdk:"response_headers_hiding_policy" json:"response_headers_hiding_policy,computed_optional"`
	Rewrite                     customfield.NestedObject[CDNRuleTemplateOptionsRewriteModel]                     `tfsdk:"rewrite" json:"rewrite,computed_optional"`
	SecureKey                   customfield.NestedObject[CDNRuleTemplateOptionsSecureKeyModel]                   `tfsdk:"secure_key" json:"secure_key,computed_optional"`
	Slice                       customfield.NestedObject[CDNRuleTemplateOptionsSliceModel]                       `tfsdk:"slice" json:"slice,computed_optional"`
	Sni                         customfield.NestedObject[CDNRuleTemplateOptionsSniModel]                         `tfsdk:"sni" json:"sni,computed_optional"`
	Stale                       customfield.NestedObject[CDNRuleTemplateOptionsStaleModel]                       `tfsdk:"stale" json:"stale,computed_optional"`
	StaticResponseHeaders       customfield.NestedObject[CDNRuleTemplateOptionsStaticResponseHeadersModel]       `tfsdk:"static_response_headers" json:"static_response_headers,computed_optional"`
	StaticHeaders               customfield.NestedObject[CDNRuleTemplateOptionsStaticHeadersModel]               `tfsdk:"static_headers" json:"staticHeaders,computed_optional"`
	StaticRequestHeaders        customfield.NestedObject[CDNRuleTemplateOptionsStaticRequestHeadersModel]        `tfsdk:"static_request_headers" json:"staticRequestHeaders,computed_optional"`
	UserAgentACL                customfield.NestedObject[CDNRuleTemplateOptionsUserAgentACLModel]                `tfsdk:"user_agent_acl" json:"user_agent_acl,computed_optional"`
	Waap                        customfield.NestedObject[CDNRuleTemplateOptionsWaapModel]                        `tfsdk:"waap" json:"waap,computed_optional"`
	Websockets                  customfield.NestedObject[CDNRuleTemplateOptionsWebsocketsModel]                  `tfsdk:"websockets" json:"websockets,computed_optional"`
}

type CDNRuleTemplateOptionsAllowedHTTPMethodsModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsBotProtectionModel struct {
	BotChallenge *CDNRuleTemplateOptionsBotProtectionBotChallengeModel `tfsdk:"bot_challenge" json:"bot_challenge,required"`
	Enabled      types.Bool                                            `tfsdk:"enabled" json:"enabled,required"`
}

type CDNRuleTemplateOptionsBotProtectionBotChallengeModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed_optional"`
}

type CDNRuleTemplateOptionsBrotliCompressionModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsBrowserCacheSettingsModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsCacheHTTPHeadersModel struct {
	Enabled types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]types.String `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsCorsModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
	Always  types.Bool                    `tfsdk:"always" json:"always,computed_optional"`
}

type CDNRuleTemplateOptionsCountryACLModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNRuleTemplateOptionsDisableCacheModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsDisableProxyForceRangesModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsEdgeCacheSettingsModel struct {
	Enabled      types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	CustomValues customfield.Map[types.String] `tfsdk:"custom_values" json:"custom_values,computed_optional"`
	Default      types.String                  `tfsdk:"default" json:"default,optional"`
	Value        types.String                  `tfsdk:"value" json:"value,optional"`
}

type CDNRuleTemplateOptionsFastedgeModel struct {
	Enabled           types.Bool                                                                     `tfsdk:"enabled" json:"enabled,required"`
	OnRequestBody     customfield.NestedObject[CDNRuleTemplateOptionsFastedgeOnRequestBodyModel]     `tfsdk:"on_request_body" json:"on_request_body,computed_optional"`
	OnRequestHeaders  customfield.NestedObject[CDNRuleTemplateOptionsFastedgeOnRequestHeadersModel]  `tfsdk:"on_request_headers" json:"on_request_headers,computed_optional"`
	OnResponseBody    customfield.NestedObject[CDNRuleTemplateOptionsFastedgeOnResponseBodyModel]    `tfsdk:"on_response_body" json:"on_response_body,computed_optional"`
	OnResponseHeaders customfield.NestedObject[CDNRuleTemplateOptionsFastedgeOnResponseHeadersModel] `tfsdk:"on_response_headers" json:"on_response_headers,computed_optional"`
}

type CDNRuleTemplateOptionsFastedgeOnRequestBodyModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNRuleTemplateOptionsFastedgeOnRequestHeadersModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNRuleTemplateOptionsFastedgeOnResponseBodyModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNRuleTemplateOptionsFastedgeOnResponseHeadersModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNRuleTemplateOptionsFetchCompressedModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsFollowOriginRedirectModel struct {
	Codes   customfield.Set[types.Int64] `tfsdk:"codes" json:"codes,required"`
	Enabled types.Bool                   `tfsdk:"enabled" json:"enabled,required"`
}

type CDNRuleTemplateOptionsForceReturnModel struct {
	Body         types.String                                                                 `tfsdk:"body" json:"body,required"`
	Code         types.Int64                                                                  `tfsdk:"code" json:"code,required"`
	Enabled      types.Bool                                                                   `tfsdk:"enabled" json:"enabled,required"`
	TimeInterval customfield.NestedObject[CDNRuleTemplateOptionsForceReturnTimeIntervalModel] `tfsdk:"time_interval" json:"time_interval,computed_optional"`
}

type CDNRuleTemplateOptionsForceReturnTimeIntervalModel struct {
	EndTime   types.String `tfsdk:"end_time" json:"end_time,required"`
	StartTime types.String `tfsdk:"start_time" json:"start_time,required"`
	TimeZone  types.String `tfsdk:"time_zone" json:"time_zone,computed_optional"`
}

type CDNRuleTemplateOptionsForwardHostHeaderModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsGzipOnModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsHostHeaderModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsIgnoreCookieModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsIgnoreQueryStringModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsImageStackModel struct {
	Enabled     types.Bool  `tfsdk:"enabled" json:"enabled,required"`
	AvifEnabled types.Bool  `tfsdk:"avif_enabled" json:"avif_enabled,computed_optional"`
	PngLossless types.Bool  `tfsdk:"png_lossless" json:"png_lossless,computed_optional"`
	Quality     types.Int64 `tfsdk:"quality" json:"quality,computed_optional"`
	WebpEnabled types.Bool  `tfsdk:"webp_enabled" json:"webp_enabled,computed_optional"`
}

type CDNRuleTemplateOptionsIPAddressACLModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNRuleTemplateOptionsLimitBandwidthModel struct {
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	LimitType types.String `tfsdk:"limit_type" json:"limit_type,required"`
	Buffer    types.Int64  `tfsdk:"buffer" json:"buffer,optional"`
	Speed     types.Int64  `tfsdk:"speed" json:"speed,optional"`
}

type CDNRuleTemplateOptionsProxyCacheKeyModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsProxyCacheMethodsSetModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsProxyConnectTimeoutModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsProxyReadTimeoutModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsQueryParamsBlacklistModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsQueryParamsWhitelistModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsQueryStringForwardingModel struct {
	Enabled              types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ForwardFromFileTypes customfield.Set[types.String] `tfsdk:"forward_from_file_types" json:"forward_from_file_types,required"`
	ForwardToFileTypes   customfield.Set[types.String] `tfsdk:"forward_to_file_types" json:"forward_to_file_types,required"`
	ForwardExceptKeys    customfield.Set[types.String] `tfsdk:"forward_except_keys" json:"forward_except_keys,optional"`
	ForwardOnlyKeys      customfield.Set[types.String] `tfsdk:"forward_only_keys" json:"forward_only_keys,optional"`
}

type CDNRuleTemplateOptionsRedirectHTTPToHTTPSModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsRedirectHTTPSToHTTPModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsReferrerACLModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNRuleTemplateOptionsRequestLimiterModel struct {
	Enabled  types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Rate     types.Int64  `tfsdk:"rate" json:"rate,required"`
	Burst    types.Int64  `tfsdk:"burst" json:"burst,computed"`
	Delay    types.Int64  `tfsdk:"delay" json:"delay,computed"`
	RateUnit types.String `tfsdk:"rate_unit" json:"rate_unit,computed_optional"`
}

type CDNRuleTemplateOptionsResponseHeadersHidingPolicyModel struct {
	Enabled  types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Excepted customfield.Set[types.String] `tfsdk:"excepted" json:"excepted,computed_optional"`
	Mode     types.String                  `tfsdk:"mode" json:"mode,required"`
}

type CDNRuleTemplateOptionsRewriteModel struct {
	Body    types.String `tfsdk:"body" json:"body,required"`
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Flag    types.String `tfsdk:"flag" json:"flag,computed_optional"`
}

type CDNRuleTemplateOptionsSecureKeyModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Key     types.String `tfsdk:"key" json:"key,required"`
	Type    types.Int64  `tfsdk:"type" json:"type,computed_optional"`
}

type CDNRuleTemplateOptionsSliceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsSniModel struct {
	CustomHostname types.String `tfsdk:"custom_hostname" json:"custom_hostname,required"`
	Enabled        types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	SniType        types.String `tfsdk:"sni_type" json:"sni_type,computed_optional"`
}

type CDNRuleTemplateOptionsStaleModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsStaticResponseHeadersModel struct {
	Enabled types.Bool                                                `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]*CDNRuleTemplateOptionsStaticResponseHeadersValueModel `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsStaticResponseHeadersValueModel struct {
	Name   types.String    `tfsdk:"name" json:"name,required"`
	Value  *[]types.String `tfsdk:"value" json:"value,required"`
	Always types.Bool      `tfsdk:"always" json:"always,computed_optional"`
}

type CDNRuleTemplateOptionsStaticHeadersModel struct {
	Enabled types.Bool                  `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.MetaStringValue `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsStaticRequestHeadersModel struct {
	Enabled types.Bool               `tfsdk:"enabled" json:"enabled,required"`
	Value   *map[string]types.String `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsUserAgentACLModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNRuleTemplateOptionsWaapModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNRuleTemplateOptionsWebsocketsModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

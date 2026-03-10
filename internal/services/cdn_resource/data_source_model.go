// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_resource

import (
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNResourceDataSourceModel struct {
	ID                 types.Int64                                                 `tfsdk:"id" path:"resource_id,computed"`
	ResourceID         types.Int64                                                 `tfsdk:"resource_id" path:"resource_id,required"`
	Active             types.Bool                                                  `tfsdk:"active" json:"active,computed"`
	CanPurgeByURLs     types.Bool                                                  `tfsdk:"can_purge_by_urls" json:"can_purge_by_urls,computed"`
	Client             types.Int64                                                 `tfsdk:"client" json:"client,computed"`
	Cname              types.String                                                `tfsdk:"cname" json:"cname,computed"`
	Created            types.String                                                `tfsdk:"created" json:"created,computed"`
	Deleted            types.Bool                                                  `tfsdk:"deleted" json:"deleted,computed"`
	Description        types.String                                                `tfsdk:"description" json:"description,computed"`
	Enabled            types.Bool                                                  `tfsdk:"enabled" json:"enabled,computed"`
	FullCustomEnabled  types.Bool                                                  `tfsdk:"full_custom_enabled" json:"full_custom_enabled,computed"`
	IsPrimary          types.Bool                                                  `tfsdk:"is_primary" json:"is_primary,computed"`
	Name               types.String                                                `tfsdk:"name" json:"name,computed"`
	OriginGroup        types.Int64                                                 `tfsdk:"origin_group" json:"originGroup,computed"`
	OriginGroupName    types.String                                                `tfsdk:"origin_group_name" json:"originGroup_name,computed"`
	OriginProtocol     types.String                                                `tfsdk:"origin_protocol" json:"originProtocol,computed"`
	PresetApplied      types.Bool                                                  `tfsdk:"preset_applied" json:"preset_applied,computed"`
	PrimaryResource    types.Int64                                                 `tfsdk:"primary_resource" json:"primary_resource,computed"`
	ProxySslCa         types.Int64                                                 `tfsdk:"proxy_ssl_ca" json:"proxy_ssl_ca,computed"`
	ProxySslData       types.Int64                                                 `tfsdk:"proxy_ssl_data" json:"proxy_ssl_data,computed"`
	ProxySslEnabled    types.Bool                                                  `tfsdk:"proxy_ssl_enabled" json:"proxy_ssl_enabled,computed"`
	ShieldDc           types.String                                                `tfsdk:"shield_dc" json:"shield_dc,computed"`
	ShieldEnabled      types.Bool                                                  `tfsdk:"shield_enabled" json:"shield_enabled,computed"`
	ShieldRoutingMap   types.Int64                                                 `tfsdk:"shield_routing_map" json:"shield_routing_map,computed"`
	Shielded           types.Bool                                                  `tfsdk:"shielded" json:"shielded,computed"`
	SslData            types.Int64                                                 `tfsdk:"ssl_data" json:"sslData,computed"`
	SslEnabled         types.Bool                                                  `tfsdk:"ssl_enabled" json:"sslEnabled,computed"`
	Status             types.String                                                `tfsdk:"status" json:"status,computed"`
	SuspendDate        types.String                                                `tfsdk:"suspend_date" json:"suspend_date,computed"`
	Suspended          types.Bool                                                  `tfsdk:"suspended" json:"suspended,computed"`
	Updated            types.String                                                `tfsdk:"updated" json:"updated,computed"`
	VpEnabled          types.Bool                                                  `tfsdk:"vp_enabled" json:"vp_enabled,computed"`
	WaapDomainID       types.String                                                `tfsdk:"waap_domain_id" json:"waap_domain_id,computed"`
	Rules              customfield.List[jsontypes.Normalized]                      `tfsdk:"rules" json:"rules,computed"`
	SecondaryHostnames customfield.List[types.String]                              `tfsdk:"secondary_hostnames" json:"secondaryHostnames,computed"`
	Options            customfield.NestedObject[CDNResourceOptionsDataSourceModel] `tfsdk:"options" json:"options,computed"`
}

type CDNResourceOptionsDataSourceModel struct {
	AllowedHTTPMethods          customfield.NestedObject[CDNResourceOptionsAllowedHTTPMethodsDataSourceModel]          `tfsdk:"allowed_http_methods" json:"allowedHttpMethods,computed"`
	BotProtection               customfield.NestedObject[CDNResourceOptionsBotProtectionDataSourceModel]               `tfsdk:"bot_protection" json:"bot_protection,computed"`
	BrotliCompression           customfield.NestedObject[CDNResourceOptionsBrotliCompressionDataSourceModel]           `tfsdk:"brotli_compression" json:"brotli_compression,computed"`
	BrowserCacheSettings        customfield.NestedObject[CDNResourceOptionsBrowserCacheSettingsDataSourceModel]        `tfsdk:"browser_cache_settings" json:"browser_cache_settings,computed"`
	CacheHTTPHeaders            customfield.NestedObject[CDNResourceOptionsCacheHTTPHeadersDataSourceModel]            `tfsdk:"cache_http_headers" json:"cache_http_headers,computed"`
	Cors                        customfield.NestedObject[CDNResourceOptionsCorsDataSourceModel]                        `tfsdk:"cors" json:"cors,computed"`
	CountryACL                  customfield.NestedObject[CDNResourceOptionsCountryACLDataSourceModel]                  `tfsdk:"country_acl" json:"country_acl,computed"`
	DisableCache                customfield.NestedObject[CDNResourceOptionsDisableCacheDataSourceModel]                `tfsdk:"disable_cache" json:"disable_cache,computed"`
	DisableProxyForceRanges     customfield.NestedObject[CDNResourceOptionsDisableProxyForceRangesDataSourceModel]     `tfsdk:"disable_proxy_force_ranges" json:"disable_proxy_force_ranges,computed"`
	EdgeCacheSettings           customfield.NestedObject[CDNResourceOptionsEdgeCacheSettingsDataSourceModel]           `tfsdk:"edge_cache_settings" json:"edge_cache_settings,computed"`
	Fastedge                    customfield.NestedObject[CDNResourceOptionsFastedgeDataSourceModel]                    `tfsdk:"fastedge" json:"fastedge,computed"`
	FetchCompressed             customfield.NestedObject[CDNResourceOptionsFetchCompressedDataSourceModel]             `tfsdk:"fetch_compressed" json:"fetch_compressed,computed"`
	FollowOriginRedirect        customfield.NestedObject[CDNResourceOptionsFollowOriginRedirectDataSourceModel]        `tfsdk:"follow_origin_redirect" json:"follow_origin_redirect,computed"`
	ForceReturn                 customfield.NestedObject[CDNResourceOptionsForceReturnDataSourceModel]                 `tfsdk:"force_return" json:"force_return,computed"`
	ForwardHostHeader           customfield.NestedObject[CDNResourceOptionsForwardHostHeaderDataSourceModel]           `tfsdk:"forward_host_header" json:"forward_host_header,computed"`
	GrpcPassthrough             customfield.NestedObject[CDNResourceOptionsGrpcPassthroughDataSourceModel]             `tfsdk:"grpc_passthrough" json:"grpc_passthrough,computed"`
	GzipOn                      customfield.NestedObject[CDNResourceOptionsGzipOnDataSourceModel]                      `tfsdk:"gzip_on" json:"gzipOn,computed"`
	HostHeader                  customfield.NestedObject[CDNResourceOptionsHostHeaderDataSourceModel]                  `tfsdk:"host_header" json:"hostHeader,computed"`
	Http3Enabled                customfield.NestedObject[CDNResourceOptionsHttp3EnabledDataSourceModel]                `tfsdk:"http3_enabled" json:"http3_enabled,computed"`
	IgnoreCookie                customfield.NestedObject[CDNResourceOptionsIgnoreCookieDataSourceModel]                `tfsdk:"ignore_cookie" json:"ignore_cookie,computed"`
	IgnoreQueryString           customfield.NestedObject[CDNResourceOptionsIgnoreQueryStringDataSourceModel]           `tfsdk:"ignore_query_string" json:"ignoreQueryString,computed"`
	ImageStack                  customfield.NestedObject[CDNResourceOptionsImageStackDataSourceModel]                  `tfsdk:"image_stack" json:"image_stack,computed"`
	IPAddressACL                customfield.NestedObject[CDNResourceOptionsIPAddressACLDataSourceModel]                `tfsdk:"ip_address_acl" json:"ip_address_acl,computed"`
	LimitBandwidth              customfield.NestedObject[CDNResourceOptionsLimitBandwidthDataSourceModel]              `tfsdk:"limit_bandwidth" json:"limit_bandwidth,computed"`
	ProxyCacheKey               customfield.NestedObject[CDNResourceOptionsProxyCacheKeyDataSourceModel]               `tfsdk:"proxy_cache_key" json:"proxy_cache_key,computed"`
	ProxyCacheMethodsSet        customfield.NestedObject[CDNResourceOptionsProxyCacheMethodsSetDataSourceModel]        `tfsdk:"proxy_cache_methods_set" json:"proxy_cache_methods_set,computed"`
	ProxyConnectTimeout         customfield.NestedObject[CDNResourceOptionsProxyConnectTimeoutDataSourceModel]         `tfsdk:"proxy_connect_timeout" json:"proxy_connect_timeout,computed"`
	ProxyReadTimeout            customfield.NestedObject[CDNResourceOptionsProxyReadTimeoutDataSourceModel]            `tfsdk:"proxy_read_timeout" json:"proxy_read_timeout,computed"`
	QueryParamsBlacklist        customfield.NestedObject[CDNResourceOptionsQueryParamsBlacklistDataSourceModel]        `tfsdk:"query_params_blacklist" json:"query_params_blacklist,computed"`
	QueryParamsWhitelist        customfield.NestedObject[CDNResourceOptionsQueryParamsWhitelistDataSourceModel]        `tfsdk:"query_params_whitelist" json:"query_params_whitelist,computed"`
	QueryStringForwarding       customfield.NestedObject[CDNResourceOptionsQueryStringForwardingDataSourceModel]       `tfsdk:"query_string_forwarding" json:"query_string_forwarding,computed"`
	RedirectHTTPToHTTPS         customfield.NestedObject[CDNResourceOptionsRedirectHTTPToHTTPSDataSourceModel]         `tfsdk:"redirect_http_to_https" json:"redirect_http_to_https,computed"`
	RedirectHTTPSToHTTP         customfield.NestedObject[CDNResourceOptionsRedirectHTTPSToHTTPDataSourceModel]         `tfsdk:"redirect_https_to_http" json:"redirect_https_to_http,computed"`
	ReferrerACL                 customfield.NestedObject[CDNResourceOptionsReferrerACLDataSourceModel]                 `tfsdk:"referrer_acl" json:"referrer_acl,computed"`
	RequestLimiter              customfield.NestedObject[CDNResourceOptionsRequestLimiterDataSourceModel]              `tfsdk:"request_limiter" json:"request_limiter,computed"`
	ResponseHeadersHidingPolicy customfield.NestedObject[CDNResourceOptionsResponseHeadersHidingPolicyDataSourceModel] `tfsdk:"response_headers_hiding_policy" json:"response_headers_hiding_policy,computed"`
	Rewrite                     customfield.NestedObject[CDNResourceOptionsRewriteDataSourceModel]                     `tfsdk:"rewrite" json:"rewrite,computed"`
	SecureKey                   customfield.NestedObject[CDNResourceOptionsSecureKeyDataSourceModel]                   `tfsdk:"secure_key" json:"secure_key,computed"`
	Slice                       customfield.NestedObject[CDNResourceOptionsSliceDataSourceModel]                       `tfsdk:"slice" json:"slice,computed"`
	Sni                         customfield.NestedObject[CDNResourceOptionsSniDataSourceModel]                         `tfsdk:"sni" json:"sni,computed"`
	Stale                       customfield.NestedObject[CDNResourceOptionsStaleDataSourceModel]                       `tfsdk:"stale" json:"stale,computed"`
	StaticResponseHeaders       customfield.NestedObject[CDNResourceOptionsStaticResponseHeadersDataSourceModel]       `tfsdk:"static_response_headers" json:"static_response_headers,computed"`
	StaticHeaders               customfield.NestedObject[CDNResourceOptionsStaticHeadersDataSourceModel]               `tfsdk:"static_headers" json:"staticHeaders,computed"`
	StaticRequestHeaders        customfield.NestedObject[CDNResourceOptionsStaticRequestHeadersDataSourceModel]        `tfsdk:"static_request_headers" json:"staticRequestHeaders,computed"`
	TlsVersions                 customfield.NestedObject[CDNResourceOptionsTlsVersionsDataSourceModel]                 `tfsdk:"tls_versions" json:"tls_versions,computed"`
	UseDefaultLeChain           customfield.NestedObject[CDNResourceOptionsUseDefaultLeChainDataSourceModel]           `tfsdk:"use_default_le_chain" json:"use_default_le_chain,computed"`
	UseDns01LeChallenge         customfield.NestedObject[CDNResourceOptionsUseDns01LeChallengeDataSourceModel]         `tfsdk:"use_dns01_le_challenge" json:"use_dns01_le_challenge,computed"`
	UseRsaLeCert                customfield.NestedObject[CDNResourceOptionsUseRsaLeCertDataSourceModel]                `tfsdk:"use_rsa_le_cert" json:"use_rsa_le_cert,computed"`
	UserAgentACL                customfield.NestedObject[CDNResourceOptionsUserAgentACLDataSourceModel]                `tfsdk:"user_agent_acl" json:"user_agent_acl,computed"`
	Waap                        customfield.NestedObject[CDNResourceOptionsWaapDataSourceModel]                        `tfsdk:"waap" json:"waap,computed"`
	Websockets                  customfield.NestedObject[CDNResourceOptionsWebsocketsDataSourceModel]                  `tfsdk:"websockets" json:"websockets,computed"`
}

type CDNResourceOptionsAllowedHTTPMethodsDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsBotProtectionDataSourceModel struct {
	BotChallenge customfield.NestedObject[CDNResourceOptionsBotProtectionBotChallengeDataSourceModel] `tfsdk:"bot_challenge" json:"bot_challenge,computed"`
	Enabled      types.Bool                                                                           `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNResourceOptionsBotProtectionBotChallengeDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNResourceOptionsBrotliCompressionDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsBrowserCacheSettingsDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsCacheHTTPHeadersDataSourceModel struct {
	Enabled types.Bool                     `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsCorsDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
	Always  types.Bool                    `tfsdk:"always" json:"always,computed"`
}

type CDNResourceOptionsCountryACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourceOptionsDisableCacheDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsDisableProxyForceRangesDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsEdgeCacheSettingsDataSourceModel struct {
	Enabled      types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	CustomValues customfield.Map[types.String] `tfsdk:"custom_values" json:"custom_values,computed"`
	Default      types.String                  `tfsdk:"default" json:"default,computed"`
	Value        types.String                  `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsFastedgeDataSourceModel struct {
	Enabled           types.Bool                                                                           `tfsdk:"enabled" json:"enabled,computed"`
	OnRequestBody     customfield.NestedObject[CDNResourceOptionsFastedgeOnRequestBodyDataSourceModel]     `tfsdk:"on_request_body" json:"on_request_body,computed"`
	OnRequestHeaders  customfield.NestedObject[CDNResourceOptionsFastedgeOnRequestHeadersDataSourceModel]  `tfsdk:"on_request_headers" json:"on_request_headers,computed"`
	OnResponseBody    customfield.NestedObject[CDNResourceOptionsFastedgeOnResponseBodyDataSourceModel]    `tfsdk:"on_response_body" json:"on_response_body,computed"`
	OnResponseHeaders customfield.NestedObject[CDNResourceOptionsFastedgeOnResponseHeadersDataSourceModel] `tfsdk:"on_response_headers" json:"on_response_headers,computed"`
}

type CDNResourceOptionsFastedgeOnRequestBodyDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourceOptionsFastedgeOnRequestHeadersDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourceOptionsFastedgeOnResponseBodyDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourceOptionsFastedgeOnResponseHeadersDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourceOptionsFetchCompressedDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsFollowOriginRedirectDataSourceModel struct {
	Codes   customfield.Set[types.Int64] `tfsdk:"codes" json:"codes,computed"`
	Enabled types.Bool                   `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNResourceOptionsForceReturnDataSourceModel struct {
	Body         types.String                                                                       `tfsdk:"body" json:"body,computed"`
	Code         types.Int64                                                                        `tfsdk:"code" json:"code,computed"`
	Enabled      types.Bool                                                                         `tfsdk:"enabled" json:"enabled,computed"`
	TimeInterval customfield.NestedObject[CDNResourceOptionsForceReturnTimeIntervalDataSourceModel] `tfsdk:"time_interval" json:"time_interval,computed"`
}

type CDNResourceOptionsForceReturnTimeIntervalDataSourceModel struct {
	EndTime   types.String `tfsdk:"end_time" json:"end_time,computed"`
	StartTime types.String `tfsdk:"start_time" json:"start_time,computed"`
	TimeZone  types.String `tfsdk:"time_zone" json:"time_zone,computed"`
}

type CDNResourceOptionsForwardHostHeaderDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsGrpcPassthroughDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsGzipOnDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsHostHeaderDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsHttp3EnabledDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsIgnoreCookieDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsIgnoreQueryStringDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsImageStackDataSourceModel struct {
	Enabled     types.Bool  `tfsdk:"enabled" json:"enabled,computed"`
	AvifEnabled types.Bool  `tfsdk:"avif_enabled" json:"avif_enabled,computed"`
	PngLossless types.Bool  `tfsdk:"png_lossless" json:"png_lossless,computed"`
	Quality     types.Int64 `tfsdk:"quality" json:"quality,computed"`
	WebpEnabled types.Bool  `tfsdk:"webp_enabled" json:"webp_enabled,computed"`
}

type CDNResourceOptionsIPAddressACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourceOptionsLimitBandwidthDataSourceModel struct {
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	LimitType types.String `tfsdk:"limit_type" json:"limit_type,computed"`
	Buffer    types.Int64  `tfsdk:"buffer" json:"buffer,computed"`
	Speed     types.Int64  `tfsdk:"speed" json:"speed,computed"`
}

type CDNResourceOptionsProxyCacheKeyDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsProxyCacheMethodsSetDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsProxyConnectTimeoutDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsProxyReadTimeoutDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsQueryParamsBlacklistDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsQueryParamsWhitelistDataSourceModel struct {
	Enabled types.Bool                     `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsQueryStringForwardingDataSourceModel struct {
	Enabled              types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ForwardFromFileTypes customfield.Set[types.String] `tfsdk:"forward_from_file_types" json:"forward_from_file_types,computed"`
	ForwardToFileTypes   customfield.Set[types.String] `tfsdk:"forward_to_file_types" json:"forward_to_file_types,computed"`
	ForwardExceptKeys    customfield.Set[types.String] `tfsdk:"forward_except_keys" json:"forward_except_keys,computed"`
	ForwardOnlyKeys      customfield.Set[types.String] `tfsdk:"forward_only_keys" json:"forward_only_keys,computed"`
}

type CDNResourceOptionsRedirectHTTPToHTTPSDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsRedirectHTTPSToHTTPDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsReferrerACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourceOptionsRequestLimiterDataSourceModel struct {
	Enabled  types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Rate     types.Int64  `tfsdk:"rate" json:"rate,computed"`
	Burst    types.Int64  `tfsdk:"burst" json:"burst,computed"`
	Delay    types.Int64  `tfsdk:"delay" json:"delay,computed"`
	RateUnit types.String `tfsdk:"rate_unit" json:"rate_unit,computed"`
}

type CDNResourceOptionsResponseHeadersHidingPolicyDataSourceModel struct {
	Enabled  types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Excepted customfield.Set[types.String] `tfsdk:"excepted" json:"excepted,computed"`
	Mode     types.String                  `tfsdk:"mode" json:"mode,computed"`
}

type CDNResourceOptionsRewriteDataSourceModel struct {
	Body    types.String `tfsdk:"body" json:"body,computed"`
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Flag    types.String `tfsdk:"flag" json:"flag,computed"`
}

type CDNResourceOptionsSecureKeyDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Key     types.String `tfsdk:"key" json:"key,computed"`
	Type    types.Int64  `tfsdk:"type" json:"type,computed"`
}

type CDNResourceOptionsSliceDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsSniDataSourceModel struct {
	CustomHostname types.String `tfsdk:"custom_hostname" json:"custom_hostname,computed"`
	Enabled        types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	SniType        types.String `tfsdk:"sni_type" json:"sni_type,computed"`
}

type CDNResourceOptionsStaleDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsStaticResponseHeadersDataSourceModel struct {
	Enabled types.Bool                                                                                `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.NestedObjectList[CDNResourceOptionsStaticResponseHeadersValueDataSourceModel] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsStaticResponseHeadersValueDataSourceModel struct {
	Name   types.String                   `tfsdk:"name" json:"name,computed"`
	Value  customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
	Always types.Bool                     `tfsdk:"always" json:"always,computed"`
}

type CDNResourceOptionsStaticHeadersDataSourceModel struct {
	Enabled types.Bool           `tfsdk:"enabled" json:"enabled,computed"`
	Value   jsontypes.Normalized `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsStaticRequestHeadersDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Map[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsTlsVersionsDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsUseDefaultLeChainDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsUseDns01LeChallengeDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsUseRsaLeCertDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsUserAgentACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourceOptionsWaapDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourceOptionsWebsocketsDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

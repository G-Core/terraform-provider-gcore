// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_resource

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNResourceModel struct {
	ID                   types.Int64                                       `tfsdk:"id" json:"id,computed"`
	Cname                types.String                                      `tfsdk:"cname" json:"cname,required"`
	PrimaryResource      types.Int64                                       `tfsdk:"primary_resource" json:"primary_resource,optional"`
	WaapAPIDomainEnabled types.Bool                                        `tfsdk:"waap_api_domain_enabled" json:"waap_api_domain_enabled,optional,no_refresh"`
	Origin               types.String                                      `tfsdk:"origin" json:"origin,computed_optional,no_refresh"`
	Active               types.Bool                                        `tfsdk:"active" json:"active,computed_optional"`
	Description          types.String                                      `tfsdk:"description" json:"description,computed_optional"`
	Name                 types.String                                      `tfsdk:"name" json:"name,computed_optional"`
	OriginGroup          types.Int64                                       `tfsdk:"origin_group" json:"originGroup,computed_optional"`
	OriginProtocol       types.String                                      `tfsdk:"origin_protocol" json:"originProtocol,computed_optional"`
	ProxySslCa           types.Int64                                       `tfsdk:"proxy_ssl_ca" json:"proxy_ssl_ca,computed_optional"`
	ProxySslData         types.Int64                                       `tfsdk:"proxy_ssl_data" json:"proxy_ssl_data,computed_optional"`
	ProxySslEnabled      types.Bool                                        `tfsdk:"proxy_ssl_enabled" json:"proxy_ssl_enabled,computed_optional"`
	SslData              types.Int64                                       `tfsdk:"ssl_data" json:"sslData,computed_optional"`
	SslEnabled           types.Bool                                        `tfsdk:"ssl_enabled" json:"sslEnabled,computed_optional"`
	SecondaryHostnames   customfield.Set[types.String]                     `tfsdk:"secondary_hostnames" json:"secondaryHostnames,computed_optional,no_refresh"`
	Options              customfield.NestedObject[CDNResourceOptionsModel] `tfsdk:"options" json:"options,computed_optional"`
	CanPurgeByURLs       types.Bool                                        `tfsdk:"can_purge_by_urls" json:"can_purge_by_urls,computed"`
	Client               types.Int64                                       `tfsdk:"client" json:"client,computed"`
	Created              types.String                                      `tfsdk:"created" json:"created,computed"`
	Deleted              types.Bool                                        `tfsdk:"deleted" json:"deleted,computed"`
	Enabled              types.Bool                                        `tfsdk:"enabled" json:"enabled,computed"`
	FullCustomEnabled    types.Bool                                        `tfsdk:"full_custom_enabled" json:"full_custom_enabled,computed"`
	IsPrimary            types.Bool                                        `tfsdk:"is_primary" json:"is_primary,computed"`
	OriginGroupName      types.String                                      `tfsdk:"origin_group_name" json:"originGroup_name,computed"`
	PresetApplied        types.Bool                                        `tfsdk:"preset_applied" json:"preset_applied,computed"`
	ShieldDc             types.String                                      `tfsdk:"shield_dc" json:"shield_dc,computed"`
	ShieldEnabled        types.Bool                                        `tfsdk:"shield_enabled" json:"shield_enabled,computed"`
	ShieldRoutingMap     types.Int64                                       `tfsdk:"shield_routing_map" json:"shield_routing_map,computed"`
	Shielded             types.Bool                                        `tfsdk:"shielded" json:"shielded,computed"`
	Status               types.String                                      `tfsdk:"status" json:"status,computed"`
	SuspendDate          types.String                                      `tfsdk:"suspend_date" json:"suspend_date,computed"`
	Suspended            types.Bool                                        `tfsdk:"suspended" json:"suspended,computed"`
	Updated              types.String                                      `tfsdk:"updated" json:"updated,computed"`
	VpEnabled            types.Bool                                        `tfsdk:"vp_enabled" json:"vp_enabled,computed"`
	WaapDomainID         types.String                                      `tfsdk:"waap_domain_id" json:"waap_domain_id,computed"`
	Rules                customfield.List[jsontypes.Normalized]            `tfsdk:"rules" json:"rules,computed"`
}

func (m CDNResourceModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CDNResourceModel) MarshalJSONForUpdate(state CDNResourceModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CDNResourceOptionsModel struct {
	AllowedHTTPMethods          customfield.NestedObject[CDNResourceOptionsAllowedHTTPMethodsModel]          `tfsdk:"allowed_http_methods" json:"allowedHttpMethods,computed_optional"`
	BotProtection               customfield.NestedObject[CDNResourceOptionsBotProtectionModel]               `tfsdk:"bot_protection" json:"bot_protection,computed_optional"`
	BrotliCompression           customfield.NestedObject[CDNResourceOptionsBrotliCompressionModel]           `tfsdk:"brotli_compression" json:"brotli_compression,computed_optional"`
	BrowserCacheSettings        customfield.NestedObject[CDNResourceOptionsBrowserCacheSettingsModel]        `tfsdk:"browser_cache_settings" json:"browser_cache_settings,computed_optional"`
	CacheHTTPHeaders            *CDNResourceOptionsCacheHTTPHeadersModel                                     `tfsdk:"cache_http_headers" json:"cache_http_headers,optional"`
	Cors                        customfield.NestedObject[CDNResourceOptionsCorsModel]                        `tfsdk:"cors" json:"cors,computed_optional"`
	CountryACL                  customfield.NestedObject[CDNResourceOptionsCountryACLModel]                  `tfsdk:"country_acl" json:"country_acl,computed_optional"`
	DisableCache                *CDNResourceOptionsDisableCacheModel                                         `tfsdk:"disable_cache" json:"disable_cache,optional"`
	DisableProxyForceRanges     customfield.NestedObject[CDNResourceOptionsDisableProxyForceRangesModel]     `tfsdk:"disable_proxy_force_ranges" json:"disable_proxy_force_ranges,computed_optional"`
	EdgeCacheSettings           customfield.NestedObject[CDNResourceOptionsEdgeCacheSettingsModel]           `tfsdk:"edge_cache_settings" json:"edge_cache_settings,computed_optional"`
	Fastedge                    customfield.NestedObject[CDNResourceOptionsFastedgeModel]                    `tfsdk:"fastedge" json:"fastedge,computed_optional"`
	FetchCompressed             customfield.NestedObject[CDNResourceOptionsFetchCompressedModel]             `tfsdk:"fetch_compressed" json:"fetch_compressed,computed_optional"`
	FollowOriginRedirect        customfield.NestedObject[CDNResourceOptionsFollowOriginRedirectModel]        `tfsdk:"follow_origin_redirect" json:"follow_origin_redirect,computed_optional"`
	ForceReturn                 customfield.NestedObject[CDNResourceOptionsForceReturnModel]                 `tfsdk:"force_return" json:"force_return,computed_optional"`
	ForwardHostHeader           customfield.NestedObject[CDNResourceOptionsForwardHostHeaderModel]           `tfsdk:"forward_host_header" json:"forward_host_header,computed_optional"`
	GrpcPassthrough             customfield.NestedObject[CDNResourceOptionsGrpcPassthroughModel]             `tfsdk:"grpc_passthrough" json:"grpc_passthrough,computed_optional"`
	GzipOn                      customfield.NestedObject[CDNResourceOptionsGzipOnModel]                      `tfsdk:"gzip_on" json:"gzipOn,computed_optional"`
	HostHeader                  customfield.NestedObject[CDNResourceOptionsHostHeaderModel]                  `tfsdk:"host_header" json:"hostHeader,computed_optional"`
	Http3Enabled                customfield.NestedObject[CDNResourceOptionsHttp3EnabledModel]                `tfsdk:"http3_enabled" json:"http3_enabled,computed_optional"`
	IgnoreCookie                customfield.NestedObject[CDNResourceOptionsIgnoreCookieModel]                `tfsdk:"ignore_cookie" json:"ignore_cookie,computed_optional"`
	IgnoreQueryString           customfield.NestedObject[CDNResourceOptionsIgnoreQueryStringModel]           `tfsdk:"ignore_query_string" json:"ignoreQueryString,computed_optional"`
	ImageStack                  customfield.NestedObject[CDNResourceOptionsImageStackModel]                  `tfsdk:"image_stack" json:"image_stack,computed_optional"`
	IPAddressACL                customfield.NestedObject[CDNResourceOptionsIPAddressACLModel]                `tfsdk:"ip_address_acl" json:"ip_address_acl,computed_optional"`
	LimitBandwidth              customfield.NestedObject[CDNResourceOptionsLimitBandwidthModel]              `tfsdk:"limit_bandwidth" json:"limit_bandwidth,computed_optional"`
	ProxyCacheKey               customfield.NestedObject[CDNResourceOptionsProxyCacheKeyModel]               `tfsdk:"proxy_cache_key" json:"proxy_cache_key,computed_optional"`
	ProxyCacheMethodsSet        customfield.NestedObject[CDNResourceOptionsProxyCacheMethodsSetModel]        `tfsdk:"proxy_cache_methods_set" json:"proxy_cache_methods_set,computed_optional"`
	ProxyConnectTimeout         customfield.NestedObject[CDNResourceOptionsProxyConnectTimeoutModel]         `tfsdk:"proxy_connect_timeout" json:"proxy_connect_timeout,computed_optional"`
	ProxyReadTimeout            customfield.NestedObject[CDNResourceOptionsProxyReadTimeoutModel]            `tfsdk:"proxy_read_timeout" json:"proxy_read_timeout,computed_optional"`
	QueryParamsBlacklist        customfield.NestedObject[CDNResourceOptionsQueryParamsBlacklistModel]        `tfsdk:"query_params_blacklist" json:"query_params_blacklist,computed_optional"`
	QueryParamsWhitelist        customfield.NestedObject[CDNResourceOptionsQueryParamsWhitelistModel]        `tfsdk:"query_params_whitelist" json:"query_params_whitelist,computed_optional"`
	QueryStringForwarding       customfield.NestedObject[CDNResourceOptionsQueryStringForwardingModel]       `tfsdk:"query_string_forwarding" json:"query_string_forwarding,computed_optional"`
	RedirectHTTPToHTTPS         customfield.NestedObject[CDNResourceOptionsRedirectHTTPToHTTPSModel]         `tfsdk:"redirect_http_to_https" json:"redirect_http_to_https,computed_optional"`
	RedirectHTTPSToHTTP         customfield.NestedObject[CDNResourceOptionsRedirectHTTPSToHTTPModel]         `tfsdk:"redirect_https_to_http" json:"redirect_https_to_http,computed_optional"`
	ReferrerACL                 customfield.NestedObject[CDNResourceOptionsReferrerACLModel]                 `tfsdk:"referrer_acl" json:"referrer_acl,computed_optional"`
	RequestLimiter              customfield.NestedObject[CDNResourceOptionsRequestLimiterModel]              `tfsdk:"request_limiter" json:"request_limiter,computed_optional"`
	ResponseHeadersHidingPolicy customfield.NestedObject[CDNResourceOptionsResponseHeadersHidingPolicyModel] `tfsdk:"response_headers_hiding_policy" json:"response_headers_hiding_policy,computed_optional"`
	Rewrite                     customfield.NestedObject[CDNResourceOptionsRewriteModel]                     `tfsdk:"rewrite" json:"rewrite,computed_optional"`
	SecureKey                   customfield.NestedObject[CDNResourceOptionsSecureKeyModel]                   `tfsdk:"secure_key" json:"secure_key,computed_optional"`
	Slice                       customfield.NestedObject[CDNResourceOptionsSliceModel]                       `tfsdk:"slice" json:"slice,computed_optional"`
	Sni                         customfield.NestedObject[CDNResourceOptionsSniModel]                         `tfsdk:"sni" json:"sni,computed_optional"`
	Stale                       customfield.NestedObject[CDNResourceOptionsStaleModel]                       `tfsdk:"stale" json:"stale,computed_optional"`
	StaticResponseHeaders       customfield.NestedObject[CDNResourceOptionsStaticResponseHeadersModel]       `tfsdk:"static_response_headers" json:"static_response_headers,computed_optional"`
	StaticHeaders               *CDNResourceOptionsStaticHeadersModel                                        `tfsdk:"static_headers" json:"staticHeaders,optional"`
	StaticRequestHeaders        customfield.NestedObject[CDNResourceOptionsStaticRequestHeadersModel]        `tfsdk:"static_request_headers" json:"staticRequestHeaders,computed_optional"`
	TlsVersions                 customfield.NestedObject[CDNResourceOptionsTlsVersionsModel]                 `tfsdk:"tls_versions" json:"tls_versions,computed_optional"`
	UseDefaultLeChain           customfield.NestedObject[CDNResourceOptionsUseDefaultLeChainModel]           `tfsdk:"use_default_le_chain" json:"use_default_le_chain,computed_optional"`
	UseDns01LeChallenge         customfield.NestedObject[CDNResourceOptionsUseDns01LeChallengeModel]         `tfsdk:"use_dns01_le_challenge" json:"use_dns01_le_challenge,computed_optional"`
	UseRsaLeCert                customfield.NestedObject[CDNResourceOptionsUseRsaLeCertModel]                `tfsdk:"use_rsa_le_cert" json:"use_rsa_le_cert,computed_optional"`
	UserAgentACL                customfield.NestedObject[CDNResourceOptionsUserAgentACLModel]                `tfsdk:"user_agent_acl" json:"user_agent_acl,computed_optional"`
	Waap                        customfield.NestedObject[CDNResourceOptionsWaapModel]                        `tfsdk:"waap" json:"waap,computed_optional"`
	Websockets                  customfield.NestedObject[CDNResourceOptionsWebsocketsModel]                  `tfsdk:"websockets" json:"websockets,computed_optional"`
}

type CDNResourceOptionsAllowedHTTPMethodsModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsBotProtectionModel struct {
	BotChallenge *CDNResourceOptionsBotProtectionBotChallengeModel `tfsdk:"bot_challenge" json:"bot_challenge,required"`
	Enabled      types.Bool                                        `tfsdk:"enabled" json:"enabled,required"`
}

type CDNResourceOptionsBotProtectionBotChallengeModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed_optional"`
}

type CDNResourceOptionsBrotliCompressionModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsBrowserCacheSettingsModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsCacheHTTPHeadersModel struct {
	Enabled types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsCorsModel struct {
	Enabled types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]types.String `tfsdk:"value" json:"value,required"`
	Always  types.Bool      `tfsdk:"always" json:"always,computed_optional"`
}

type CDNResourceOptionsCountryACLModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNResourceOptionsDisableCacheModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsDisableProxyForceRangesModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsEdgeCacheSettingsModel struct {
	Enabled      types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	CustomValues customfield.Map[types.String] `tfsdk:"custom_values" json:"custom_values,computed_optional"`
	Default      types.String                  `tfsdk:"default" json:"default,optional"`
	Value        types.String                  `tfsdk:"value" json:"value,optional"`
}

type CDNResourceOptionsFastedgeModel struct {
	Enabled           types.Bool                                                                 `tfsdk:"enabled" json:"enabled,required"`
	OnRequestBody     customfield.NestedObject[CDNResourceOptionsFastedgeOnRequestBodyModel]     `tfsdk:"on_request_body" json:"on_request_body,computed_optional"`
	OnRequestHeaders  customfield.NestedObject[CDNResourceOptionsFastedgeOnRequestHeadersModel]  `tfsdk:"on_request_headers" json:"on_request_headers,computed_optional"`
	OnResponseBody    customfield.NestedObject[CDNResourceOptionsFastedgeOnResponseBodyModel]    `tfsdk:"on_response_body" json:"on_response_body,computed_optional"`
	OnResponseHeaders customfield.NestedObject[CDNResourceOptionsFastedgeOnResponseHeadersModel] `tfsdk:"on_response_headers" json:"on_response_headers,computed_optional"`
}

type CDNResourceOptionsFastedgeOnRequestBodyModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNResourceOptionsFastedgeOnRequestHeadersModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNResourceOptionsFastedgeOnResponseBodyModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNResourceOptionsFastedgeOnResponseHeadersModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,required"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed_optional"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed_optional"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed_optional"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed_optional"`
}

type CDNResourceOptionsFetchCompressedModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsFollowOriginRedirectModel struct {
	Codes   customfield.Set[types.Int64] `tfsdk:"codes" json:"codes,required"`
	Enabled types.Bool                   `tfsdk:"enabled" json:"enabled,required"`
}

type CDNResourceOptionsForceReturnModel struct {
	Body         types.String                                                             `tfsdk:"body" json:"body,required"`
	Code         types.Int64                                                              `tfsdk:"code" json:"code,required"`
	Enabled      types.Bool                                                               `tfsdk:"enabled" json:"enabled,required"`
	TimeInterval customfield.NestedObject[CDNResourceOptionsForceReturnTimeIntervalModel] `tfsdk:"time_interval" json:"time_interval,computed_optional"`
}

type CDNResourceOptionsForceReturnTimeIntervalModel struct {
	EndTime   types.String `tfsdk:"end_time" json:"end_time,required"`
	StartTime types.String `tfsdk:"start_time" json:"start_time,required"`
	TimeZone  types.String `tfsdk:"time_zone" json:"time_zone,computed_optional"`
}

type CDNResourceOptionsForwardHostHeaderModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsGrpcPassthroughModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsGzipOnModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsHostHeaderModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsHttp3EnabledModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsIgnoreCookieModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsIgnoreQueryStringModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsImageStackModel struct {
	Enabled     types.Bool  `tfsdk:"enabled" json:"enabled,required"`
	AvifEnabled types.Bool  `tfsdk:"avif_enabled" json:"avif_enabled,computed_optional"`
	PngLossless types.Bool  `tfsdk:"png_lossless" json:"png_lossless,computed_optional"`
	Quality     types.Int64 `tfsdk:"quality" json:"quality,computed_optional"`
	WebpEnabled types.Bool  `tfsdk:"webp_enabled" json:"webp_enabled,computed_optional"`
}

type CDNResourceOptionsIPAddressACLModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNResourceOptionsLimitBandwidthModel struct {
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	LimitType types.String `tfsdk:"limit_type" json:"limit_type,required"`
	Buffer    types.Int64  `tfsdk:"buffer" json:"buffer,optional"`
	Speed     types.Int64  `tfsdk:"speed" json:"speed,optional"`
}

type CDNResourceOptionsProxyCacheKeyModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsProxyCacheMethodsSetModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsProxyConnectTimeoutModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsProxyReadTimeoutModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Value   types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsQueryParamsBlacklistModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsQueryParamsWhitelistModel struct {
	Enabled types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsQueryStringForwardingModel struct {
	Enabled              types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ForwardFromFileTypes customfield.Set[types.String] `tfsdk:"forward_from_file_types" json:"forward_from_file_types,required"`
	ForwardToFileTypes   customfield.Set[types.String] `tfsdk:"forward_to_file_types" json:"forward_to_file_types,required"`
	ForwardExceptKeys    customfield.Set[types.String] `tfsdk:"forward_except_keys" json:"forward_except_keys,optional"`
	ForwardOnlyKeys      customfield.Set[types.String] `tfsdk:"forward_only_keys" json:"forward_only_keys,optional"`
}

type CDNResourceOptionsRedirectHTTPToHTTPSModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsRedirectHTTPSToHTTPModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsReferrerACLModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNResourceOptionsRequestLimiterModel struct {
	Enabled  types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Rate     types.Int64  `tfsdk:"rate" json:"rate,required"`
	Burst    types.Int64  `tfsdk:"burst" json:"burst,computed"`
	Delay    types.Int64  `tfsdk:"delay" json:"delay,computed"`
	RateUnit types.String `tfsdk:"rate_unit" json:"rate_unit,computed_optional"`
}

type CDNResourceOptionsResponseHeadersHidingPolicyModel struct {
	Enabled  types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Excepted customfield.Set[types.String] `tfsdk:"excepted" json:"excepted,computed_optional"`
	Mode     types.String                  `tfsdk:"mode" json:"mode,required"`
}

type CDNResourceOptionsRewriteModel struct {
	Body    types.String `tfsdk:"body" json:"body,required"`
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Flag    types.String `tfsdk:"flag" json:"flag,computed_optional"`
}

type CDNResourceOptionsSecureKeyModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Key     types.String `tfsdk:"key" json:"key,required"`
	Type    types.Int64  `tfsdk:"type" json:"type,computed_optional"`
}

type CDNResourceOptionsSliceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsSniModel struct {
	CustomHostname types.String `tfsdk:"custom_hostname" json:"custom_hostname,required"`
	Enabled        types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	SniType        types.String `tfsdk:"sni_type" json:"sni_type,computed_optional"`
}

type CDNResourceOptionsStaleModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsStaticResponseHeadersModel struct {
	Enabled types.Bool                                            `tfsdk:"enabled" json:"enabled,required"`
	Value   *[]*CDNResourceOptionsStaticResponseHeadersValueModel `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsStaticResponseHeadersValueModel struct {
	Name   types.String    `tfsdk:"name" json:"name,required"`
	Value  *[]types.String `tfsdk:"value" json:"value,required"`
	Always types.Bool      `tfsdk:"always" json:"always,computed_optional"`
}

type CDNResourceOptionsStaticHeadersModel struct {
	Enabled types.Bool           `tfsdk:"enabled" json:"enabled,required"`
	Value   jsontypes.Normalized `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsStaticRequestHeadersModel struct {
	Enabled types.Bool               `tfsdk:"enabled" json:"enabled,required"`
	Value   *map[string]types.String `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsTlsVersionsModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsUseDefaultLeChainModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsUseDns01LeChallengeModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsUseRsaLeCertModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsUserAgentACLModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,required"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,required"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,required"`
}

type CDNResourceOptionsWaapModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

type CDNResourceOptionsWebsocketsModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,required"`
	Value   types.Bool `tfsdk:"value" json:"value,required"`
}

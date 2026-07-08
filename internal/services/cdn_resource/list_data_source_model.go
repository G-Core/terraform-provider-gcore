// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_resource

import (
	"context"

	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNResourcesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CDNResourcesItemsDataSourceModel] `json:"results,computed"`
}

type CDNResourcesDataSourceModel struct {
	Active             types.Bool                                                     `tfsdk:"active" query:"active,optional"`
	Cname              types.String                                                   `tfsdk:"cname" query:"cname,optional"`
	Deleted            types.Bool                                                     `tfsdk:"deleted" query:"deleted,optional"`
	Enabled            types.Bool                                                     `tfsdk:"enabled" query:"enabled,optional"`
	IsPrimary          types.Bool                                                     `tfsdk:"is_primary" query:"is_primary,optional"`
	MaxCreated         types.String                                                   `tfsdk:"max_created" query:"max_created,optional"`
	MaxUpdated         types.String                                                   `tfsdk:"max_updated" query:"max_updated,optional"`
	MinCreated         types.String                                                   `tfsdk:"min_created" query:"min_created,optional"`
	MinUpdated         types.String                                                   `tfsdk:"min_updated" query:"min_updated,optional"`
	Name               types.String                                                   `tfsdk:"name" query:"name,optional"`
	OriginGroup        types.Int64                                                    `tfsdk:"origin_group" query:"originGroup,optional"`
	OriginGroupIn      types.String                                                   `tfsdk:"origin_group_in" query:"originGroup__in,optional"`
	OriginProtocol     types.String                                                   `tfsdk:"origin_protocol" query:"originProtocol,optional"`
	Rules              types.String                                                   `tfsdk:"rules" query:"rules,optional"`
	SecondaryHostnames types.String                                                   `tfsdk:"secondary_hostnames" query:"secondaryHostnames,optional"`
	ShieldDc           types.String                                                   `tfsdk:"shield_dc" query:"shield_dc,optional"`
	Shielded           types.Bool                                                     `tfsdk:"shielded" query:"shielded,optional"`
	SslData            types.Int64                                                    `tfsdk:"ssl_data" query:"sslData,optional"`
	SslDataIn          types.Int64                                                    `tfsdk:"ssl_data_in" query:"sslData_in,optional"`
	SslEnabled         types.Bool                                                     `tfsdk:"ssl_enabled" query:"sslEnabled,optional"`
	Status             types.String                                                   `tfsdk:"status" query:"status,optional"`
	Suspend            types.Bool                                                     `tfsdk:"suspend" query:"suspend,optional"`
	Suspended          types.Bool                                                     `tfsdk:"suspended" query:"suspended,optional"`
	VpEnabled          types.Bool                                                     `tfsdk:"vp_enabled" query:"vp_enabled,optional"`
	MaxItems           types.Int64                                                    `tfsdk:"max_items"`
	Items              customfield.NestedObjectList[CDNResourcesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CDNResourcesDataSourceModel) toListParams(_ context.Context) (params cdn.CDNResourceListParams, diags diag.Diagnostics) {
	params = cdn.CDNResourceListParams{}

	if !m.Active.IsNull() {
		params.Active = param.NewOpt(m.Active.ValueBool())
	}
	if !m.Cname.IsNull() {
		params.Cname = param.NewOpt(m.Cname.ValueString())
	}
	if !m.Deleted.IsNull() {
		params.Deleted = param.NewOpt(m.Deleted.ValueBool())
	}
	if !m.Enabled.IsNull() {
		params.Enabled = param.NewOpt(m.Enabled.ValueBool())
	}
	if !m.IsPrimary.IsNull() {
		params.IsPrimary = param.NewOpt(m.IsPrimary.ValueBool())
	}
	if !m.MaxCreated.IsNull() {
		params.MaxCreated = param.NewOpt(m.MaxCreated.ValueString())
	}
	if !m.MaxUpdated.IsNull() {
		params.MaxUpdated = param.NewOpt(m.MaxUpdated.ValueString())
	}
	if !m.MinCreated.IsNull() {
		params.MinCreated = param.NewOpt(m.MinCreated.ValueString())
	}
	if !m.MinUpdated.IsNull() {
		params.MinUpdated = param.NewOpt(m.MinUpdated.ValueString())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.OriginGroup.IsNull() {
		params.OriginGroup = param.NewOpt(m.OriginGroup.ValueInt64())
	}
	if !m.OriginGroupIn.IsNull() {
		params.OriginGroupIn = param.NewOpt(m.OriginGroupIn.ValueString())
	}
	if !m.OriginProtocol.IsNull() {
		params.OriginProtocol = cdn.CDNResourceListParamsOriginProtocol(m.OriginProtocol.ValueString())
	}
	if !m.Rules.IsNull() {
		params.Rules = param.NewOpt(m.Rules.ValueString())
	}
	if !m.SecondaryHostnames.IsNull() {
		params.SecondaryHostnames = param.NewOpt(m.SecondaryHostnames.ValueString())
	}
	if !m.ShieldDc.IsNull() {
		params.ShieldDc = param.NewOpt(m.ShieldDc.ValueString())
	}
	if !m.Shielded.IsNull() {
		params.Shielded = param.NewOpt(m.Shielded.ValueBool())
	}
	if !m.SslData.IsNull() {
		params.SslData = param.NewOpt(m.SslData.ValueInt64())
	}
	if !m.SslDataIn.IsNull() {
		params.SslDataIn = param.NewOpt(m.SslDataIn.ValueInt64())
	}
	if !m.SslEnabled.IsNull() {
		params.SslEnabled = param.NewOpt(m.SslEnabled.ValueBool())
	}
	if !m.Status.IsNull() {
		params.Status = cdn.CDNResourceListParamsStatus(m.Status.ValueString())
	}
	if !m.Suspend.IsNull() {
		params.Suspend = param.NewOpt(m.Suspend.ValueBool())
	}
	if !m.Suspended.IsNull() {
		params.Suspended = param.NewOpt(m.Suspended.ValueBool())
	}
	if !m.VpEnabled.IsNull() {
		params.VpEnabled = param.NewOpt(m.VpEnabled.ValueBool())
	}

	return
}

type CDNResourcesItemsDataSourceModel struct {
	ID                 types.Int64                                                  `tfsdk:"id" json:"id,computed"`
	Active             types.Bool                                                   `tfsdk:"active" json:"active,computed"`
	CanPurgeByURLs     types.Bool                                                   `tfsdk:"can_purge_by_urls" json:"can_purge_by_urls,computed"`
	Client             types.Int64                                                  `tfsdk:"client" json:"client,computed"`
	Cname              types.String                                                 `tfsdk:"cname" json:"cname,computed"`
	Created            types.String                                                 `tfsdk:"created" json:"created,computed"`
	Deleted            types.Bool                                                   `tfsdk:"deleted" json:"deleted,computed"`
	Description        types.String                                                 `tfsdk:"description" json:"description,computed"`
	Enabled            types.Bool                                                   `tfsdk:"enabled" json:"enabled,computed"`
	FullCustomEnabled  types.Bool                                                   `tfsdk:"full_custom_enabled" json:"full_custom_enabled,computed"`
	IsPrimary          types.Bool                                                   `tfsdk:"is_primary" json:"is_primary,computed"`
	Name               types.String                                                 `tfsdk:"name" json:"name,computed"`
	Options            customfield.NestedObject[CDNResourcesOptionsDataSourceModel] `tfsdk:"options" json:"options,computed"`
	OriginGroup        types.Int64                                                  `tfsdk:"origin_group" json:"originGroup,computed"`
	OriginGroupName    types.String                                                 `tfsdk:"origin_group_name" json:"originGroup_name,computed"`
	OriginProtocol     types.String                                                 `tfsdk:"origin_protocol" json:"originProtocol,computed"`
	PresetApplied      types.Bool                                                   `tfsdk:"preset_applied" json:"preset_applied,computed"`
	PrimaryResource    types.Int64                                                  `tfsdk:"primary_resource" json:"primary_resource,computed"`
	ProxySslCa         types.Int64                                                  `tfsdk:"proxy_ssl_ca" json:"proxy_ssl_ca,computed"`
	ProxySslData       types.Int64                                                  `tfsdk:"proxy_ssl_data" json:"proxy_ssl_data,computed"`
	ProxySslEnabled    types.Bool                                                   `tfsdk:"proxy_ssl_enabled" json:"proxy_ssl_enabled,computed"`
	Rules              customfield.List[jsontypes.Normalized]                       `tfsdk:"rules" json:"rules,computed"`
	SecondaryHostnames customfield.List[types.String]                               `tfsdk:"secondary_hostnames" json:"secondaryHostnames,computed"`
	ShieldDc           types.String                                                 `tfsdk:"shield_dc" json:"shield_dc,computed"`
	ShieldEnabled      types.Bool                                                   `tfsdk:"shield_enabled" json:"shield_enabled,computed"`
	ShieldRoutingMap   types.Int64                                                  `tfsdk:"shield_routing_map" json:"shield_routing_map,computed"`
	Shielded           types.Bool                                                   `tfsdk:"shielded" json:"shielded,computed"`
	SslData            types.Int64                                                  `tfsdk:"ssl_data" json:"sslData,computed"`
	SslEnabled         types.Bool                                                   `tfsdk:"ssl_enabled" json:"sslEnabled,computed"`
	Status             types.String                                                 `tfsdk:"status" json:"status,computed"`
	SuspendDate        types.String                                                 `tfsdk:"suspend_date" json:"suspend_date,computed"`
	Suspended          types.Bool                                                   `tfsdk:"suspended" json:"suspended,computed"`
	Updated            types.String                                                 `tfsdk:"updated" json:"updated,computed"`
	VpEnabled          types.Bool                                                   `tfsdk:"vp_enabled" json:"vp_enabled,computed"`
	WaapDomainID       types.String                                                 `tfsdk:"waap_domain_id" json:"waap_domain_id,computed"`
}

type CDNResourcesOptionsDataSourceModel struct {
	AllowedHTTPMethods          customfield.NestedObject[CDNResourcesOptionsAllowedHTTPMethodsDataSourceModel]          `tfsdk:"allowed_http_methods" json:"allowedHttpMethods,computed"`
	BrotliCompression           customfield.NestedObject[CDNResourcesOptionsBrotliCompressionDataSourceModel]           `tfsdk:"brotli_compression" json:"brotli_compression,computed"`
	BrowserCacheSettings        customfield.NestedObject[CDNResourcesOptionsBrowserCacheSettingsDataSourceModel]        `tfsdk:"browser_cache_settings" json:"browser_cache_settings,computed"`
	CacheHTTPHeaders            customfield.NestedObject[CDNResourcesOptionsCacheHTTPHeadersDataSourceModel]            `tfsdk:"cache_http_headers" json:"cache_http_headers,computed"`
	Cors                        customfield.NestedObject[CDNResourcesOptionsCorsDataSourceModel]                        `tfsdk:"cors" json:"cors,computed"`
	CountryACL                  customfield.NestedObject[CDNResourcesOptionsCountryACLDataSourceModel]                  `tfsdk:"country_acl" json:"country_acl,computed"`
	DisableCache                customfield.NestedObject[CDNResourcesOptionsDisableCacheDataSourceModel]                `tfsdk:"disable_cache" json:"disable_cache,computed"`
	DisableProxyForceRanges     customfield.NestedObject[CDNResourcesOptionsDisableProxyForceRangesDataSourceModel]     `tfsdk:"disable_proxy_force_ranges" json:"disable_proxy_force_ranges,computed"`
	EdgeCacheSettings           customfield.NestedObject[CDNResourcesOptionsEdgeCacheSettingsDataSourceModel]           `tfsdk:"edge_cache_settings" json:"edge_cache_settings,computed"`
	Fastedge                    customfield.NestedObject[CDNResourcesOptionsFastedgeDataSourceModel]                    `tfsdk:"fastedge" json:"fastedge,computed"`
	FetchCompressed             customfield.NestedObject[CDNResourcesOptionsFetchCompressedDataSourceModel]             `tfsdk:"fetch_compressed" json:"fetch_compressed,computed"`
	FollowOriginRedirect        customfield.NestedObject[CDNResourcesOptionsFollowOriginRedirectDataSourceModel]        `tfsdk:"follow_origin_redirect" json:"follow_origin_redirect,computed"`
	ForceReturn                 customfield.NestedObject[CDNResourcesOptionsForceReturnDataSourceModel]                 `tfsdk:"force_return" json:"force_return,computed"`
	ForwardHostHeader           customfield.NestedObject[CDNResourcesOptionsForwardHostHeaderDataSourceModel]           `tfsdk:"forward_host_header" json:"forward_host_header,computed"`
	GrpcPassthrough             customfield.NestedObject[CDNResourcesOptionsGrpcPassthroughDataSourceModel]             `tfsdk:"grpc_passthrough" json:"grpc_passthrough,computed"`
	GzipOn                      customfield.NestedObject[CDNResourcesOptionsGzipOnDataSourceModel]                      `tfsdk:"gzip_on" json:"gzipOn,computed"`
	HostHeader                  customfield.NestedObject[CDNResourcesOptionsHostHeaderDataSourceModel]                  `tfsdk:"host_header" json:"hostHeader,computed"`
	Http3Enabled                customfield.NestedObject[CDNResourcesOptionsHttp3EnabledDataSourceModel]                `tfsdk:"http3_enabled" json:"http3_enabled,computed"`
	IgnoreCookie                customfield.NestedObject[CDNResourcesOptionsIgnoreCookieDataSourceModel]                `tfsdk:"ignore_cookie" json:"ignore_cookie,computed"`
	IgnoreQueryString           customfield.NestedObject[CDNResourcesOptionsIgnoreQueryStringDataSourceModel]           `tfsdk:"ignore_query_string" json:"ignoreQueryString,computed"`
	ImageStack                  customfield.NestedObject[CDNResourcesOptionsImageStackDataSourceModel]                  `tfsdk:"image_stack" json:"image_stack,computed"`
	IPAddressACL                customfield.NestedObject[CDNResourcesOptionsIPAddressACLDataSourceModel]                `tfsdk:"ip_address_acl" json:"ip_address_acl,computed"`
	LimitBandwidth              customfield.NestedObject[CDNResourcesOptionsLimitBandwidthDataSourceModel]              `tfsdk:"limit_bandwidth" json:"limit_bandwidth,computed"`
	NetworkErrorLogging         customfield.NestedObject[CDNResourcesOptionsNetworkErrorLoggingDataSourceModel]         `tfsdk:"network_error_logging" json:"network_error_logging,computed"`
	ProxyCacheKey               customfield.NestedObject[CDNResourcesOptionsProxyCacheKeyDataSourceModel]               `tfsdk:"proxy_cache_key" json:"proxy_cache_key,computed"`
	ProxyCacheMethodsSet        customfield.NestedObject[CDNResourcesOptionsProxyCacheMethodsSetDataSourceModel]        `tfsdk:"proxy_cache_methods_set" json:"proxy_cache_methods_set,computed"`
	ProxyConnectTimeout         customfield.NestedObject[CDNResourcesOptionsProxyConnectTimeoutDataSourceModel]         `tfsdk:"proxy_connect_timeout" json:"proxy_connect_timeout,computed"`
	ProxyReadTimeout            customfield.NestedObject[CDNResourcesOptionsProxyReadTimeoutDataSourceModel]            `tfsdk:"proxy_read_timeout" json:"proxy_read_timeout,computed"`
	QueryParamsBlacklist        customfield.NestedObject[CDNResourcesOptionsQueryParamsBlacklistDataSourceModel]        `tfsdk:"query_params_blacklist" json:"query_params_blacklist,computed"`
	QueryParamsWhitelist        customfield.NestedObject[CDNResourcesOptionsQueryParamsWhitelistDataSourceModel]        `tfsdk:"query_params_whitelist" json:"query_params_whitelist,computed"`
	QueryStringForwarding       customfield.NestedObject[CDNResourcesOptionsQueryStringForwardingDataSourceModel]       `tfsdk:"query_string_forwarding" json:"query_string_forwarding,computed"`
	RedirectHTTPToHTTPS         customfield.NestedObject[CDNResourcesOptionsRedirectHTTPToHTTPSDataSourceModel]         `tfsdk:"redirect_http_to_https" json:"redirect_http_to_https,computed"`
	RedirectHTTPSToHTTP         customfield.NestedObject[CDNResourcesOptionsRedirectHTTPSToHTTPDataSourceModel]         `tfsdk:"redirect_https_to_http" json:"redirect_https_to_http,computed"`
	ReferrerACL                 customfield.NestedObject[CDNResourcesOptionsReferrerACLDataSourceModel]                 `tfsdk:"referrer_acl" json:"referrer_acl,computed"`
	ResponseHeadersHidingPolicy customfield.NestedObject[CDNResourcesOptionsResponseHeadersHidingPolicyDataSourceModel] `tfsdk:"response_headers_hiding_policy" json:"response_headers_hiding_policy,computed"`
	Rewrite                     customfield.NestedObject[CDNResourcesOptionsRewriteDataSourceModel]                     `tfsdk:"rewrite" json:"rewrite,computed"`
	SecureKey                   customfield.NestedObject[CDNResourcesOptionsSecureKeyDataSourceModel]                   `tfsdk:"secure_key" json:"secure_key,computed"`
	Slice                       customfield.NestedObject[CDNResourcesOptionsSliceDataSourceModel]                       `tfsdk:"slice" json:"slice,computed"`
	Sni                         customfield.NestedObject[CDNResourcesOptionsSniDataSourceModel]                         `tfsdk:"sni" json:"sni,computed"`
	Stale                       customfield.NestedObject[CDNResourcesOptionsStaleDataSourceModel]                       `tfsdk:"stale" json:"stale,computed"`
	StaticResponseHeaders       customfield.NestedObject[CDNResourcesOptionsStaticResponseHeadersDataSourceModel]       `tfsdk:"static_response_headers" json:"static_response_headers,computed"`
	StaticHeaders               customfield.NestedObject[CDNResourcesOptionsStaticHeadersDataSourceModel]               `tfsdk:"static_headers" json:"staticHeaders,computed"`
	StaticRequestHeaders        customfield.NestedObject[CDNResourcesOptionsStaticRequestHeadersDataSourceModel]        `tfsdk:"static_request_headers" json:"staticRequestHeaders,computed"`
	TlsVersions                 customfield.NestedObject[CDNResourcesOptionsTlsVersionsDataSourceModel]                 `tfsdk:"tls_versions" json:"tls_versions,computed"`
	UseDefaultLeChain           customfield.NestedObject[CDNResourcesOptionsUseDefaultLeChainDataSourceModel]           `tfsdk:"use_default_le_chain" json:"use_default_le_chain,computed"`
	UseDns01LeChallenge         customfield.NestedObject[CDNResourcesOptionsUseDns01LeChallengeDataSourceModel]         `tfsdk:"use_dns01_le_challenge" json:"use_dns01_le_challenge,computed"`
	UseRsaLeCert                customfield.NestedObject[CDNResourcesOptionsUseRsaLeCertDataSourceModel]                `tfsdk:"use_rsa_le_cert" json:"use_rsa_le_cert,computed"`
	UserAgentACL                customfield.NestedObject[CDNResourcesOptionsUserAgentACLDataSourceModel]                `tfsdk:"user_agent_acl" json:"user_agent_acl,computed"`
	Waap                        customfield.NestedObject[CDNResourcesOptionsWaapDataSourceModel]                        `tfsdk:"waap" json:"waap,computed"`
	Websockets                  customfield.NestedObject[CDNResourcesOptionsWebsocketsDataSourceModel]                  `tfsdk:"websockets" json:"websockets,computed"`
}

type CDNResourcesOptionsAllowedHTTPMethodsDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsBrotliCompressionDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsBrowserCacheSettingsDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsCacheHTTPHeadersDataSourceModel struct {
	Enabled types.Bool                     `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsCorsDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
	Always  types.Bool                    `tfsdk:"always" json:"always,computed"`
}

type CDNResourcesOptionsCountryACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourcesOptionsDisableCacheDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsDisableProxyForceRangesDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsEdgeCacheSettingsDataSourceModel struct {
	Enabled      types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	CustomValues customfield.Map[types.String] `tfsdk:"custom_values" json:"custom_values,computed"`
	Default      types.String                  `tfsdk:"default" json:"default,computed"`
	Value        types.String                  `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsFastedgeDataSourceModel struct {
	Enabled                    types.Bool                                                                                     `tfsdk:"enabled" json:"enabled,computed"`
	OnRequestBody              customfield.NestedObject[CDNResourcesOptionsFastedgeOnRequestBodyDataSourceModel]              `tfsdk:"on_request_body" json:"on_request_body,computed"`
	OnRequestHeaders           customfield.NestedObject[CDNResourcesOptionsFastedgeOnRequestHeadersDataSourceModel]           `tfsdk:"on_request_headers" json:"on_request_headers,computed"`
	OnRequestHeadersAfterCache customfield.NestedObject[CDNResourcesOptionsFastedgeOnRequestHeadersAfterCacheDataSourceModel] `tfsdk:"on_request_headers_after_cache" json:"on_request_headers_after_cache,computed"`
	OnResponseBody             customfield.NestedObject[CDNResourcesOptionsFastedgeOnResponseBodyDataSourceModel]             `tfsdk:"on_response_body" json:"on_response_body,computed"`
	OnResponseHeaders          customfield.NestedObject[CDNResourcesOptionsFastedgeOnResponseHeadersDataSourceModel]          `tfsdk:"on_response_headers" json:"on_response_headers,computed"`
}

type CDNResourcesOptionsFastedgeOnRequestBodyDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourcesOptionsFastedgeOnRequestHeadersDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourcesOptionsFastedgeOnRequestHeadersAfterCacheDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourcesOptionsFastedgeOnResponseBodyDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourcesOptionsFastedgeOnResponseHeadersDataSourceModel struct {
	AppID            types.String `tfsdk:"app_id" json:"app_id,computed"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	ExecuteOnEdge    types.Bool   `tfsdk:"execute_on_edge" json:"execute_on_edge,computed"`
	ExecuteOnShield  types.Bool   `tfsdk:"execute_on_shield" json:"execute_on_shield,computed"`
	InterruptOnError types.Bool   `tfsdk:"interrupt_on_error" json:"interrupt_on_error,computed"`
}

type CDNResourcesOptionsFetchCompressedDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsFollowOriginRedirectDataSourceModel struct {
	Codes   customfield.Set[types.Int64] `tfsdk:"codes" json:"codes,computed"`
	Enabled types.Bool                   `tfsdk:"enabled" json:"enabled,computed"`
}

type CDNResourcesOptionsForceReturnDataSourceModel struct {
	Body         types.String                                                                        `tfsdk:"body" json:"body,computed"`
	Code         types.Int64                                                                         `tfsdk:"code" json:"code,computed"`
	Enabled      types.Bool                                                                          `tfsdk:"enabled" json:"enabled,computed"`
	TimeInterval customfield.NestedObject[CDNResourcesOptionsForceReturnTimeIntervalDataSourceModel] `tfsdk:"time_interval" json:"time_interval,computed"`
}

type CDNResourcesOptionsForceReturnTimeIntervalDataSourceModel struct {
	EndTime   types.String `tfsdk:"end_time" json:"end_time,computed"`
	StartTime types.String `tfsdk:"start_time" json:"start_time,computed"`
	TimeZone  types.String `tfsdk:"time_zone" json:"time_zone,computed"`
}

type CDNResourcesOptionsForwardHostHeaderDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsGrpcPassthroughDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsGzipOnDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsHostHeaderDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsHttp3EnabledDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsIgnoreCookieDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsIgnoreQueryStringDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsImageStackDataSourceModel struct {
	Enabled     types.Bool  `tfsdk:"enabled" json:"enabled,computed"`
	AvifEnabled types.Bool  `tfsdk:"avif_enabled" json:"avif_enabled,computed"`
	PngLossless types.Bool  `tfsdk:"png_lossless" json:"png_lossless,computed"`
	Quality     types.Int64 `tfsdk:"quality" json:"quality,computed"`
	WebpEnabled types.Bool  `tfsdk:"webp_enabled" json:"webp_enabled,computed"`
}

type CDNResourcesOptionsIPAddressACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourcesOptionsLimitBandwidthDataSourceModel struct {
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	LimitType types.String `tfsdk:"limit_type" json:"limit_type,computed"`
	Buffer    types.Int64  `tfsdk:"buffer" json:"buffer,computed"`
	Speed     types.Int64  `tfsdk:"speed" json:"speed,computed"`
}

type CDNResourcesOptionsNetworkErrorLoggingDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsProxyCacheKeyDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsProxyCacheMethodsSetDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsProxyConnectTimeoutDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsProxyReadTimeoutDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.String `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsQueryParamsBlacklistDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsQueryParamsWhitelistDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsQueryStringForwardingDataSourceModel struct {
	Enabled              types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ForwardFromFileTypes customfield.Set[types.String] `tfsdk:"forward_from_file_types" json:"forward_from_file_types,computed"`
	ForwardToFileTypes   customfield.Set[types.String] `tfsdk:"forward_to_file_types" json:"forward_to_file_types,computed"`
	ForwardExceptKeys    customfield.Set[types.String] `tfsdk:"forward_except_keys" json:"forward_except_keys,computed"`
	ForwardOnlyKeys      customfield.Set[types.String] `tfsdk:"forward_only_keys" json:"forward_only_keys,computed"`
}

type CDNResourcesOptionsRedirectHTTPToHTTPSDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsRedirectHTTPSToHTTPDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsReferrerACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourcesOptionsResponseHeadersHidingPolicyDataSourceModel struct {
	Enabled  types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Excepted customfield.Set[types.String] `tfsdk:"excepted" json:"excepted,computed"`
	Mode     types.String                  `tfsdk:"mode" json:"mode,computed"`
}

type CDNResourcesOptionsRewriteDataSourceModel struct {
	Body    types.String `tfsdk:"body" json:"body,computed"`
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Flag    types.String `tfsdk:"flag" json:"flag,computed"`
}

type CDNResourcesOptionsSecureKeyDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Key     types.String `tfsdk:"key" json:"key,computed"`
	Type    types.Int64  `tfsdk:"type" json:"type,computed"`
}

type CDNResourcesOptionsSliceDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsSniDataSourceModel struct {
	CustomHostname types.String `tfsdk:"custom_hostname" json:"custom_hostname,computed"`
	Enabled        types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	SniType        types.String `tfsdk:"sni_type" json:"sni_type,computed"`
}

type CDNResourcesOptionsStaleDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsStaticResponseHeadersDataSourceModel struct {
	Enabled types.Bool                                                                                 `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.NestedObjectList[CDNResourcesOptionsStaticResponseHeadersValueDataSourceModel] `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsStaticResponseHeadersValueDataSourceModel struct {
	Name   types.String                   `tfsdk:"name" json:"name,computed"`
	Value  customfield.List[types.String] `tfsdk:"value" json:"value,computed"`
	Always types.Bool                     `tfsdk:"always" json:"always,computed"`
}

type CDNResourcesOptionsStaticHeadersDataSourceModel struct {
	Enabled types.Bool           `tfsdk:"enabled" json:"enabled,computed"`
	Value   jsontypes.Normalized `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsStaticRequestHeadersDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Map[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsTlsVersionsDataSourceModel struct {
	Enabled types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	Value   customfield.Set[types.String] `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsUseDefaultLeChainDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsUseDns01LeChallengeDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsUseRsaLeCertDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsUserAgentACLDataSourceModel struct {
	Enabled        types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ExceptedValues customfield.Set[types.String] `tfsdk:"excepted_values" json:"excepted_values,computed"`
	PolicyType     types.String                  `tfsdk:"policy_type" json:"policy_type,computed"`
}

type CDNResourcesOptionsWaapDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

type CDNResourcesOptionsWebsocketsDataSourceModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,computed"`
	Value   types.Bool `tfsdk:"value" json:"value,computed"`
}

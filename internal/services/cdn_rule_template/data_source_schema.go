// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_rule_template

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CDNRuleTemplateDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "CDN rule templates define reusable rule configurations that can be applied across multiple CDN resources for consistent caching, delivery, and security policies.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"rule_template_id": schema.Int64Attribute{
				Required: true,
			},
			"client": schema.Int64Attribute{
				Description: "Client ID",
				Computed:    true,
			},
			"default": schema.BoolAttribute{
				Description: "Defines whether the template is a system template developed for common cases. System templates are available to all customers.\n\nPossible values:\n- **true** - Template is a system template and cannot be changed by a user.\n- **false** - Template is a custom template and can be changed by a user.",
				Computed:    true,
			},
			"deleted": schema.BoolAttribute{
				Description: "Defines whether the template has been deleted.\n\nPossible values:\n- **true** - Template has been deleted.\n- **false** - Template has not been deleted.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Rule template name.",
				Computed:    true,
			},
			"override_origin_protocol": schema.StringAttribute{
				Description: "Sets a protocol other than the one specified in the CDN resource settings to connect to the origin.\n\nPossible values:\n- **HTTPS** - CDN servers connect to origin via HTTPS protocol.\n- **HTTP** - CDN servers connect to origin via HTTP protocol.\n- **MATCH** - Connection protocol is chosen automatically; in this case, content on origin source should be available for the CDN both through HTTP and HTTPS protocols.\n- **null** - `originProtocol` setting is inherited from the CDN resource settings.\nAvailable values: \"HTTPS\", \"HTTP\", \"MATCH\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"HTTPS",
						"HTTP",
						"MATCH",
					),
				},
			},
			"rule": schema.StringAttribute{
				Description: "Path to the file or folder for which the rule will be applied.\n\nThe rule is applied if the requested URI matches the rule path.\n\nWe add a leading forward slash to any rule path. Specify a path without a forward slash.",
				Computed:    true,
			},
			"rule_type": schema.Int64Attribute{
				Description: "Rule type.\n\nPossible values:\n- **Type 0** - Regular expression. Must start with '^/' or '/'.\n- **Type 1** - Regular expression. Note that for this rule type we automatically add / to each rule pattern before your regular expression. This type is **legacy**, please use Type 0.",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 1),
				},
			},
			"template": schema.BoolAttribute{
				Description: "Determines whether the rule is a template.",
				Computed:    true,
			},
			"weight": schema.Int64Attribute{
				Description: "Rule execution order: from lowest (1) to highest.\n\nIf requested URI matches multiple rules, the one higher in the order of the rules will be applied.",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 2147483647),
				},
			},
			"options": schema.SingleNestedAttribute{
				Description: "List of options that can be configured for the rule.\n\nIn case of `null` value the option is not added to the rule.\nOption inherits its value from the CDN resource settings.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"allowed_http_methods": schema.SingleNestedAttribute{
						Description: "HTTP methods allowed for content requests from the CDN.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsAllowedHTTPMethodsDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.SetAttribute{
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
						},
					},
					"bot_protection": schema.SingleNestedAttribute{
						Description: "Allows to prevent online services from overloading and ensure your business workflow running smoothly.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsBotProtectionDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"bot_challenge": schema.SingleNestedAttribute{
								Description: "Controls the bot challenge module state.",
								Computed:    true,
								CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsBotProtectionBotChallengeDataSourceModel](ctx),
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Description: "Possible values:\n- **true** - Bot challenge is enabled.\n- **false** - Bot challenge is disabled.",
										Computed:    true,
									},
								},
							},
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
					"brotli_compression": schema.SingleNestedAttribute{
						Description: "Compresses content with Brotli on the CDN side. CDN servers will request only uncompressed content from the origin.\n\nNotes:\n\n1. CDN only supports \"Brotli compression\" when the \"origin shielding\" feature is activated.\n2. If a precache server is not active for a CDN resource, no compression occurs, even if the option is enabled.\n3. `brotli_compression` is not supported with `fetch_compressed` or `slice` options enabled.\n4. `fetch_compressed` option in CDN resource settings overrides `brotli_compression` in rules. If you enabled `fetch_compressed` in CDN resource and want to enable `brotli_compression` in a rule, you must specify `fetch_compressed:false` in the rule.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsBrotliCompressionDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.SetAttribute{
								Description: "Allows to select the content types you want to compress.\n\n`text/html` is a mandatory content type.",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
						},
					},
					"browser_cache_settings": schema.SingleNestedAttribute{
						Description: "Cache expiration time for users browsers in seconds.\n\nCache expiration time is applied to the following response codes: 200, 201, 204, 206, 301, 302, 303, 304, 307, 308.\n\nResponses with other codes will not be cached.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsBrowserCacheSettingsDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.StringAttribute{
								Description: "Set the cache expiration time to '0s' to disable caching.\n\nThe maximum duration is any equivalent to `1y`.",
								Computed:    true,
							},
						},
					},
					"cache_http_headers": schema.SingleNestedAttribute{
						Description:        "**Legacy option**. Use the `response_headers_hiding_policy` option instead.\n\nHTTP Headers that must be included in the response.",
						Computed:           true,
						DeprecationMessage: "This attribute is deprecated.",
						CustomType:         customfield.NewNestedObjectType[CDNRuleTemplateOptionsCacheHTTPHeadersDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.ListAttribute{
								Computed:    true,
								CustomType:  customfield.NewListType[types.String](ctx),
								ElementType: types.StringType,
							},
						},
					},
					"cors": schema.SingleNestedAttribute{
						Description: "Enables or disables CORS (Cross-Origin Resource Sharing) header support.\n\nCORS header support allows the CDN to add the Access-Control-Allow-Origin header to a response to a browser.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsCorsDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.SetAttribute{
								Description: "Value of the Access-Control-Allow-Origin header.\n\nPossible values:\n- **Adds * as the Access-Control-Allow-Origin header value** - Content will be uploaded for requests from any domain.\n`\"value\": [\"*\"]`\n- **Adds \"$http_origin\" as the Access-Control-Allow-Origin header value if the origin matches one of the listed domains** - Content will be uploaded only for requests from the domains specified in the field.\n`\"value\": [\"domain.com\", \"second.dom.com\"]`\n- **Adds \"$http_origin\" as the Access-Control-Allow-Origin header value** - Content will be uploaded for requests from any domain, and the domain from which the request was sent will be added to the \"Access-Control-Allow-Origin\" header in the response.\n`\"value\": [\"$http_origin\"]`",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
							"always": schema.BoolAttribute{
								Description: "Defines whether the Access-Control-Allow-Origin header should be added to a response from CDN regardless of response code.\n\nPossible values:\n- **true** - Header will be added to a response regardless of response code.\n- **false** - Header will only be added to responses with codes: 200, 201, 204, 206, 301, 302, 303, 304, 307, 308.",
								Computed:    true,
							},
						},
					},
					"country_acl": schema.SingleNestedAttribute{
						Description: "Enables control access to content for specified countries.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsCountryACLDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"excepted_values": schema.SetAttribute{
								Description: "List of countries according to ISO-3166-1.\n\nThe meaning of the parameter depends on `policy_type` value:\n- **allow** - List of countries for which access is prohibited.\n- **deny** - List of countries for which access is allowed.",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
							"policy_type": schema.StringAttribute{
								Description: "Defines the type of CDN resource access policy.\n\nPossible values:\n- **allow** - Access is allowed for all the countries except for those specified in `excepted_values` field.\n- **deny** - Access is denied for all the countries except for those specified in `excepted_values` field.\nAvailable values: \"allow\", \"deny\".",
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("allow", "deny"),
								},
							},
						},
					},
					"disable_cache": schema.SingleNestedAttribute{
						Description:        "**Legacy option**. Use the `edge_cache_settings` option instead.\n\nAllows the complete disabling of content caching.",
						Computed:           true,
						DeprecationMessage: "This attribute is deprecated.",
						CustomType:         customfield.NewNestedObjectType[CDNRuleTemplateOptionsDisableCacheDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - content caching is disabled.\n- **false** - content caching is enabled.",
								Computed:    true,
							},
						},
					},
					"disable_proxy_force_ranges": schema.SingleNestedAttribute{
						Description: "Allows 206 responses regardless of the settings of an origin source.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsDisableProxyForceRangesDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
					"edge_cache_settings": schema.SingleNestedAttribute{
						Description: "Cache expiration time for CDN servers.\n\n`value` and `default` fields cannot be used simultaneously.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsEdgeCacheSettingsDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"custom_values": schema.MapAttribute{
								Description: "A MAP object representing the caching time in seconds for a response with a specific response code.\n\nThese settings have a higher priority than the `value` field.\n\n- Use `any` key to specify caching time for all response codes.\n- Use `0s` value to disable caching for a specific response code.",
								Computed:    true,
								CustomType:  customfield.NewMapType[types.String](ctx),
								ElementType: types.StringType,
							},
							"default": schema.StringAttribute{
								Description: "Enables content caching according to the origin cache settings.\n\nThe value is applied to the following response codes 200, 201, 204, 206, 301, 302, 303, 304, 307, 308, if an origin server does not have caching HTTP headers.\n\nResponses with other codes will not be cached.\n\nThe maximum duration is any equivalent to `1y`.",
								Computed:    true,
							},
							"value": schema.StringAttribute{
								Description: "Caching time.\n\nThe value is applied to the following response codes: 200, 206, 301, 302.\nResponses with codes 4xx, 5xx will not be cached.\n\nUse `0s` to disable caching.\n\nThe maximum duration is any equivalent to `1y`.",
								Computed:    true,
							},
						},
					},
					"fastedge": schema.SingleNestedAttribute{
						Description: "Allows to configure FastEdge app to be called on different request/response phases.\n\nNote: At least one of `on_request_headers`, `on_request_body`, `on_response_headers`, or `on_response_body` must be specified.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsFastedgeDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"on_request_body": schema.SingleNestedAttribute{
								Description: "Allows to configure FastEdge application that will be called to handle request body as soon as CDN receives incoming HTTP request.",
								Computed:    true,
								CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsFastedgeOnRequestBodyDataSourceModel](ctx),
								Attributes: map[string]schema.Attribute{
									"app_id": schema.StringAttribute{
										Description: "The ID of the application in FastEdge.",
										Computed:    true,
									},
									"enabled": schema.BoolAttribute{
										Description: "Determines if the FastEdge application should be called whenever HTTP request headers are received.",
										Computed:    true,
									},
									"execute_on_edge": schema.BoolAttribute{
										Description: "Determines if the request should be executed at the edge nodes.",
										Computed:    true,
									},
									"execute_on_shield": schema.BoolAttribute{
										Description: "Determines if the request should be executed at the shield nodes.",
										Computed:    true,
									},
									"interrupt_on_error": schema.BoolAttribute{
										Description: "Determines if the request execution should be interrupted when an error occurs.",
										Computed:    true,
									},
								},
							},
							"on_request_headers": schema.SingleNestedAttribute{
								Description: "Allows to configure FastEdge application that will be called to handle request headers as soon as CDN receives incoming HTTP request, **before cache**.",
								Computed:    true,
								CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsFastedgeOnRequestHeadersDataSourceModel](ctx),
								Attributes: map[string]schema.Attribute{
									"app_id": schema.StringAttribute{
										Description: "The ID of the application in FastEdge.",
										Computed:    true,
									},
									"enabled": schema.BoolAttribute{
										Description: "Determines if the FastEdge application should be called whenever HTTP request headers are received.",
										Computed:    true,
									},
									"execute_on_edge": schema.BoolAttribute{
										Description: "Determines if the request should be executed at the edge nodes.",
										Computed:    true,
									},
									"execute_on_shield": schema.BoolAttribute{
										Description: "Determines if the request should be executed at the shield nodes.",
										Computed:    true,
									},
									"interrupt_on_error": schema.BoolAttribute{
										Description: "Determines if the request execution should be interrupted when an error occurs.",
										Computed:    true,
									},
								},
							},
							"on_response_body": schema.SingleNestedAttribute{
								Description: "Allows to configure FastEdge application that will be called to handle response body before CDN sends the HTTP response.",
								Computed:    true,
								CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsFastedgeOnResponseBodyDataSourceModel](ctx),
								Attributes: map[string]schema.Attribute{
									"app_id": schema.StringAttribute{
										Description: "The ID of the application in FastEdge.",
										Computed:    true,
									},
									"enabled": schema.BoolAttribute{
										Description: "Determines if the FastEdge application should be called whenever HTTP request headers are received.",
										Computed:    true,
									},
									"execute_on_edge": schema.BoolAttribute{
										Description: "Determines if the request should be executed at the edge nodes.",
										Computed:    true,
									},
									"execute_on_shield": schema.BoolAttribute{
										Description: "Determines if the request should be executed at the shield nodes.",
										Computed:    true,
									},
									"interrupt_on_error": schema.BoolAttribute{
										Description: "Determines if the request execution should be interrupted when an error occurs.",
										Computed:    true,
									},
								},
							},
							"on_response_headers": schema.SingleNestedAttribute{
								Description: "Allows to configure FastEdge application that will be called to handle response headers before CDN sends the HTTP response.",
								Computed:    true,
								CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsFastedgeOnResponseHeadersDataSourceModel](ctx),
								Attributes: map[string]schema.Attribute{
									"app_id": schema.StringAttribute{
										Description: "The ID of the application in FastEdge.",
										Computed:    true,
									},
									"enabled": schema.BoolAttribute{
										Description: "Determines if the FastEdge application should be called whenever HTTP request headers are received.",
										Computed:    true,
									},
									"execute_on_edge": schema.BoolAttribute{
										Description: "Determines if the request should be executed at the edge nodes.",
										Computed:    true,
									},
									"execute_on_shield": schema.BoolAttribute{
										Description: "Determines if the request should be executed at the shield nodes.",
										Computed:    true,
									},
									"interrupt_on_error": schema.BoolAttribute{
										Description: "Determines if the request execution should be interrupted when an error occurs.",
										Computed:    true,
									},
								},
							},
						},
					},
					"fetch_compressed": schema.SingleNestedAttribute{
						Description: "Makes the CDN request compressed content from the origin.\n\nThe origin server should support compression. CDN servers will not decompress your content even if a user browser does not accept compression.\n\nNotes:\n\n1. `fetch_compressed` is not supported with `gzipON` or `brotli_compression` or `slice` options enabled.\n2. `fetch_compressed` overrides `gzipON` and `brotli_compression` in rule. If you enable it in CDN resource and want to use `gzipON` and `brotli_compression` in a rule, you have to specify `\"fetch_compressed\": false` in the rule.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsFetchCompressedDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
					"follow_origin_redirect": schema.SingleNestedAttribute{
						Description: "Enables redirection from origin.\nIf the origin server returns a redirect, the option allows the CDN to pull the requested content from the origin server that was returned in the redirect.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsFollowOriginRedirectDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"codes": schema.SetAttribute{
								Description: "Redirect status code that the origin server returns.\n\nTo serve up to date content to end users, you will need to purge the cache after managing the option.",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.Int64](ctx),
								ElementType: types.Int64Type,
							},
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
					"force_return": schema.SingleNestedAttribute{
						Description: "Applies custom HTTP response codes for CDN content.\n\nThe following codes are reserved by our system and cannot be specified in this option: 408, 444, 477, 494, 495, 496, 497, 499.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsForceReturnDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"body": schema.StringAttribute{
								Description: "URL for redirection or text.",
								Computed:    true,
							},
							"code": schema.Int64Attribute{
								Description: "Status code value.",
								Computed:    true,
								Validators: []validator.Int64{
									int64validator.Between(100, 599),
								},
							},
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"time_interval": schema.SingleNestedAttribute{
								Description: "Controls the time at which a custom HTTP response code should be applied. By default, a custom HTTP response code is applied at any time.",
								Computed:    true,
								CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsForceReturnTimeIntervalDataSourceModel](ctx),
								Attributes: map[string]schema.Attribute{
									"end_time": schema.StringAttribute{
										Description: "Time until which a custom HTTP response code should be applied. Indicated in 24-hour format.",
										Computed:    true,
									},
									"start_time": schema.StringAttribute{
										Description: "Time from which a custom HTTP response code should be applied. Indicated in 24-hour format.",
										Computed:    true,
									},
									"time_zone": schema.StringAttribute{
										Description: "Time zone used to calculate time.",
										Computed:    true,
									},
								},
							},
						},
					},
					"forward_host_header": schema.SingleNestedAttribute{
						Description: "Forwards the Host header from a end-user request to an origin server.\n\n`hostHeader` and `forward_host_header` options cannot be enabled simultaneously.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsForwardHostHeaderDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
					"gzip_on": schema.SingleNestedAttribute{
						Description: "Compresses content with gzip on the CDN end. CDN servers will request only uncompressed content from the origin.\n\nNotes:\n\n1. Compression with gzip is not supported with `fetch_compressed` or `slice` options enabled.\n2. `fetch_compressed` option in CDN resource settings overrides `gzipON` in rules. If you enable `fetch_compressed` in CDN resource and want to enable `gzipON` in rules, you need to specify `\"fetch_compressed\":false` for rules.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsGzipOnDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
					"host_header": schema.SingleNestedAttribute{
						Description: "Sets the Host header that CDN servers use when request content from an origin server.\nYour server must be able to process requests with the chosen header.\n\nIf the option is `null`, the Host Header value is equal to first CNAME.\n\n`hostHeader` and `forward_host_header` options cannot be enabled simultaneously.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsHostHeaderDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.StringAttribute{
								Description: "Host Header value.",
								Computed:    true,
							},
						},
					},
					"ignore_cookie": schema.SingleNestedAttribute{
						Description: "Defines whether the files with the Set-Cookies header are cached as one file or as different ones.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsIgnoreCookieDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - Option is enabled, files with cookies are cached as one file.\n- **false** - Option is disabled, files with cookies are cached as different files.",
								Computed:    true,
							},
						},
					},
					"ignore_query_string": schema.SingleNestedAttribute{
						Description: "How a file with different query strings is cached: either as one object (option is enabled) or as different objects (option is disabled.)\n\n`ignoreQueryString`, `query_params_whitelist` and `query_params_blacklist` options cannot be enabled simultaneously.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsIgnoreQueryStringDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
					"image_stack": schema.SingleNestedAttribute{
						Description: "Transforms JPG and PNG images (for example, resize or crop) and automatically converts them to WebP or AVIF format.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsImageStackDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"avif_enabled": schema.BoolAttribute{
								Description: "Enables or disables automatic conversion of JPEG and PNG images to AVI format.",
								Computed:    true,
							},
							"png_lossless": schema.BoolAttribute{
								Description: "Enables or disables compression without quality loss for PNG format.",
								Computed:    true,
							},
							"quality": schema.Int64Attribute{
								Description: "Defines quality settings for JPG and PNG images. The higher the value, the better the image quality, and the larger the file size after conversion.",
								Computed:    true,
								Validators: []validator.Int64{
									int64validator.Between(1, 100),
								},
							},
							"webp_enabled": schema.BoolAttribute{
								Description: "Enables or disables automatic conversion of JPEG and PNG images to WebP format.",
								Computed:    true,
							},
						},
					},
					"ip_address_acl": schema.SingleNestedAttribute{
						Description: "Controls access to the CDN resource content for specific IP addresses.\n\nIf you want to use IPs from our CDN servers IP list for IP ACL configuration, you have to independently monitor their relevance.\n\nWe recommend you use a script for automatically update IP ACL. [Read more.](/docs/api-reference/cdn/ip-addresses-list/get-cdn-servers-ip-addresses)",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsIPAddressACLDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"excepted_values": schema.SetAttribute{
								Description: "List of IP addresses with a subnet mask.\n\nThe meaning of the parameter depends on `policy_type` value:\n- **allow** - List of IP addresses for which access is prohibited.\n- **deny** - List of IP addresses for which access is allowed.\n\nExamples:\n- `192.168.3.2/32`\n- `2a03:d000:2980:7::8/128`",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
							"policy_type": schema.StringAttribute{
								Description: "IP access policy type.\n\nPossible values:\n- **allow** - Allow access to all IPs except IPs specified in \"excepted_values\" field.\n- **deny** - Deny access to all IPs except IPs specified in \"excepted_values\" field.\nAvailable values: \"allow\", \"deny\".",
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("allow", "deny"),
								},
							},
						},
					},
					"limit_bandwidth": schema.SingleNestedAttribute{
						Description: "Allows to control the download speed per connection.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsLimitBandwidthDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"limit_type": schema.StringAttribute{
								Description: "Method of controlling the download speed per connection.\n\nPossible values:\n- **static** - Use speed and buffer fields to set the download speed limit.\n- **dynamic** - Use query strings **speed** and **buffer** to set the download speed limit.\n\nFor example, when requesting content at the link\n\n```\nhttp://cdn.example.com/video.mp4?speed=50k&buffer=500k\n```\n\nthe download speed will be limited to 50kB/s after 500 kB.\nAvailable values: \"static\", \"dynamic\".",
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("static", "dynamic"),
								},
							},
							"buffer": schema.Int64Attribute{
								Description: "Amount of downloaded data after which the user will be rate limited.",
								Computed:    true,
								Validators: []validator.Int64{
									int64validator.Between(0, 1000000000),
								},
							},
							"speed": schema.Int64Attribute{
								Description: "Maximum download speed per connection.",
								Computed:    true,
								Validators: []validator.Int64{
									int64validator.Between(1, 1000000000),
								},
							},
						},
					},
					"proxy_cache_key": schema.SingleNestedAttribute{
						Description: "Allows you to modify your cache key. If omitted, the default value is `$request_uri`.\n\nCombine the specified variables to create a key for caching.\n- **$`request_uri`**\n- **$scheme**\n- **$uri**\n\n**Warning**: Enabling and changing this option can invalidate your current cache and affect the cache hit ratio. Furthermore, the \"Purge by pattern\" option will not work.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsProxyCacheKeyDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.StringAttribute{
								Description: "Key for caching.",
								Computed:    true,
							},
						},
					},
					"proxy_cache_methods_set": schema.SingleNestedAttribute{
						Description: "Caching for POST requests along with default GET and HEAD.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsProxyCacheMethodsSetDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
					"proxy_connect_timeout": schema.SingleNestedAttribute{
						Description: "The time limit for establishing a connection with the origin.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsProxyConnectTimeoutDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.StringAttribute{
								Description: "Timeout value in seconds.\n\nSupported range: **1s - 5s**.",
								Computed:    true,
							},
						},
					},
					"proxy_read_timeout": schema.SingleNestedAttribute{
						Description: "The time limit for receiving a partial response from the origin.\nIf no response is received within this time, the connection will be closed.\n\n**Note:**\nWhen used with a WebSocket connection, this option supports values only in the range 1–20 seconds (instead of the usual 1–30 seconds).",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsProxyReadTimeoutDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.StringAttribute{
								Description: "Timeout value in seconds.\n\nSupported range: **1s - 30s**.",
								Computed:    true,
							},
						},
					},
					"query_params_blacklist": schema.SingleNestedAttribute{
						Description: "Files with the specified query parameters are cached as one object, files with other parameters are cached as different objects.\n\n`ignoreQueryString`, `query_params_whitelist` and `query_params_blacklist` options cannot be enabled simultaneously.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsQueryParamsBlacklistDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.SetAttribute{
								Description: "List of query parameters.",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
						},
					},
					"query_params_whitelist": schema.SingleNestedAttribute{
						Description: "Files with the specified query parameters are cached as different objects, files with other parameters are cached as one object.\n\n`ignoreQueryString`, `query_params_whitelist` and `query_params_blacklist` options cannot be enabled simultaneously.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsQueryParamsWhitelistDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.ListAttribute{
								Description: "List of query parameters.",
								Computed:    true,
								CustomType:  customfield.NewListType[types.String](ctx),
								ElementType: types.StringType,
							},
						},
					},
					"query_string_forwarding": schema.SingleNestedAttribute{
						Description: "The Query String Forwarding feature allows for the seamless transfer of parameters embedded in playlist files to the corresponding media chunk files.\nThis functionality ensures that specific attributes, such as authentication tokens or tracking information, are consistently passed along from the playlist manifest to the individual media segments.\nThis is particularly useful for maintaining continuity in security, analytics, and any other parameter-based operations across the entire media delivery workflow.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsQueryStringForwardingDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"forward_from_file_types": schema.SetAttribute{
								Description: "The `forward_from_files_types` field specifies the types of playlist files from which parameters will be extracted and forwarded.\nThis typically includes formats that list multiple media chunk references, such as HLS and DASH playlists.\nParameters associated with these playlist files (like query strings or headers) will be propagated to the chunks they reference.",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
							"forward_to_file_types": schema.SetAttribute{
								Description: "The field specifies the types of media chunk files to which parameters, extracted from playlist files, will be forwarded.\nThese refer to the actual segments of media content that are delivered to viewers.\nEnsuring the correct parameters are forwarded to these files is crucial for maintaining the integrity of the streaming session.",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
							"forward_except_keys": schema.SetAttribute{
								Description: "The `forward_except_keys` field provides a mechanism to exclude specific parameters from being forwarded from playlist files to media chunk files.\nBy listing certain keys in this field, you can ensure that these parameters are omitted during the forwarding process.\nThis is particularly useful for preventing sensitive or irrelevant information from being included in requests for media chunks, thereby enhancing security and optimizing performance.",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
							"forward_only_keys": schema.SetAttribute{
								Description: "The `forward_only_keys` field allows for granular control over which specific parameters are forwarded from playlist files to media chunk files.\nBy specifying certain keys, only those parameters will be propagated, ensuring that only relevant information is passed along.\nThis is particularly useful for security and performance optimization, as it prevents unnecessary or sensitive data from being included in requests for media chunks.",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
						},
					},
					"redirect_http_to_https": schema.SingleNestedAttribute{
						Description: "Enables redirect from HTTP to HTTPS.\n\n`redirect_http_to_https` and `redirect_https_to_http` options cannot be enabled simultaneously.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsRedirectHTTPToHTTPSDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
					"redirect_https_to_http": schema.SingleNestedAttribute{
						Description: "Enables redirect from HTTPS to HTTP.\n\n`redirect_http_to_https` and `redirect_https_to_http` options cannot be enabled simultaneously.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsRedirectHTTPSToHTTPDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
					"referrer_acl": schema.SingleNestedAttribute{
						Description: "Controls access to the CDN resource content for specified domain names.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsReferrerACLDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"excepted_values": schema.SetAttribute{
								Description: "List of domain names or wildcard domains (without protocol: `http://` or `https://`.)\n\nThe meaning of the parameter depends on `policy_type` value:\n- **allow** - List of domain names for which access is prohibited.\n- **deny** - List of IP domain names for which access is allowed.\n\nExamples:\n- `example.com`\n- `*.example.com`",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
							"policy_type": schema.StringAttribute{
								Description: "Policy type.\n\nPossible values:\n- **allow** - Allow access to all domain names except the domain names specified in `excepted_values` field.\n- **deny** - Deny access to all domain names except the domain names specified in `excepted_values` field.\nAvailable values: \"allow\", \"deny\".",
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("allow", "deny"),
								},
							},
						},
					},
					"request_limiter": schema.SingleNestedAttribute{
						Description: "Option allows to limit the amount of HTTP requests.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsRequestLimiterDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"rate": schema.Int64Attribute{
								Description: "Maximum request rate.",
								Computed:    true,
								Validators: []validator.Int64{
									int64validator.AtLeast(1),
								},
							},
							"burst": schema.Int64Attribute{
								Computed: true,
							},
							"delay": schema.Int64Attribute{
								Computed: true,
							},
							"rate_unit": schema.StringAttribute{
								Description: "Units of measurement for the `rate` field.\n\nPossible values:\n- **r/s** - Requests per second.\n- **r/m** - Requests per minute.\n\nIf the rate is less than one request per second, it is specified in request per minute (r/m.)\nAvailable values: \"r/s\", \"r/m\".",
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("r/s", "r/m"),
								},
							},
						},
					},
					"response_headers_hiding_policy": schema.SingleNestedAttribute{
						Description: "Hides HTTP headers from an origin server in the CDN response.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsResponseHeadersHidingPolicyDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"excepted": schema.SetAttribute{
								Description: "List of HTTP headers.\n\nParameter meaning depends on the value of the `mode` field:\n- **show** - List of HTTP headers to hide from response.\n- **hide** - List of HTTP headers to include in response. Other HTTP headers will be hidden.\n\nThe following headers are required and cannot be hidden from response:\n- `Connection`\n- `Content-Length`\n- `Content-Type`\n- `Date`\n- `Server`",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
							"mode": schema.StringAttribute{
								Description: "How HTTP headers are hidden from the response.\n\nPossible values:\n- **show** - Hide only HTTP headers listed in the `excepted` field.\n- **hide** - Hide all HTTP headers except headers listed in the \"excepted\" field.\nAvailable values: \"hide\", \"show\".",
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("hide", "show"),
								},
							},
						},
					},
					"rewrite": schema.SingleNestedAttribute{
						Description: "Changes and redirects requests from the CDN to the origin. It operates according to the [Nginx](https://nginx.org/en/docs/http/ngx_http_rewrite_module.html#rewrite) configuration.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsRewriteDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"body": schema.StringAttribute{
								Description: "Path for the Rewrite option.\n\nExample:\n- `/(.*) /media/$1`",
								Computed:    true,
							},
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"flag": schema.StringAttribute{
								Description: "Flag for the Rewrite option.\n\nPossible values:\n- **last** - Stop processing the current set of `ngx_http_rewrite_module` directives and start a search for a new location matching changed URI.\n- **break** - Stop processing the current set of the Rewrite option.\n- **redirect** - Return a temporary redirect with the 302 code; used when a replacement string does not start with `http://`, `https://`, or `$scheme`.\n- **permanent** - Return a permanent redirect with the 301 code.\nAvailable values: \"break\", \"last\", \"redirect\", \"permanent\".",
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive(
										"break",
										"last",
										"redirect",
										"permanent",
									),
								},
							},
						},
					},
					"secure_key": schema.SingleNestedAttribute{
						Description: "Configures access with tokenized URLs. This makes impossible to access content without a valid (unexpired) token.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsSecureKeyDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"key": schema.StringAttribute{
								Description: "Key generated on your side that will be used for URL signing.",
								Computed:    true,
							},
							"type": schema.Int64Attribute{
								Description: "Type of URL signing.\n\nPossible types:\n- **Type 0** - Includes end user IP to secure token generation.\n- **Type 2** - Excludes end user IP from secure token generation.\nAvailable values: 0, 2.",
								Computed:    true,
								Validators: []validator.Int64{
									int64validator.OneOf(0, 2),
								},
							},
						},
					},
					"slice": schema.SingleNestedAttribute{
						Description: "Requests and caches files larger than 10 MB in parts (no larger than 10 MB per part.) This reduces time to first byte.\n\nThe option is based on the [Slice](https://nginx.org/en/docs/http/ngx_http_slice_module.html) module.\n\nNotes:\n\n1. Origin must support HTTP Range requests.\n2. Not supported with `gzipON`, `brotli_compression` or `fetch_compressed` options enabled.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsSliceDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
					"sni": schema.SingleNestedAttribute{
						Description: "The hostname that is added to SNI requests from CDN servers to the origin server via HTTPS.\n\nSNI is generally only required if your origin uses shared hosting or does not have a dedicated IP address.\nIf the origin server presents multiple certificates, SNI allows the origin server to know which certificate to use for the connection.\n\nThe option works only if `originProtocol` parameter is `HTTPS` or `MATCH`.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsSniDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"custom_hostname": schema.StringAttribute{
								Description: "Custom SNI hostname.\n\nIt is required if `sni_type` is set to custom.",
								Computed:    true,
							},
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"sni_type": schema.StringAttribute{
								Description: "SNI (Server Name Indication) type.\n\nPossible values:\n- **dynamic** - SNI hostname depends on `hostHeader` and `forward_host_header` options.\nIt has several possible combinations:\n- If the `hostHeader` option is enabled and specified, SNI hostname matches the Host header.\n- If the `forward_host_header` option is enabled and has true value, SNI hostname matches the Host header used in the request made to a CDN.\n- If the `hostHeader` and `forward_host_header` options are disabled, SNI hostname matches the primary CNAME.\n- **custom** - custom SNI hostname is in use.\nAvailable values: \"dynamic\", \"custom\".",
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("dynamic", "custom"),
								},
							},
						},
					},
					"stale": schema.SingleNestedAttribute{
						Description: "Serves stale cached content in case of origin unavailability.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsStaleDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.SetAttribute{
								Description: `Defines list of errors for which "Always online" option is applied.`,
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
						},
					},
					"static_response_headers": schema.SingleNestedAttribute{
						Description: "Custom HTTP Headers that a CDN server adds to a response.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsStaticResponseHeadersDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.ListNestedAttribute{
								Computed:   true,
								CustomType: customfield.NewNestedObjectListType[CDNRuleTemplateOptionsStaticResponseHeadersValueDataSourceModel](ctx),
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "HTTP Header name.\n\nRestrictions:\n- Maximum 128 symbols.\n- Latin letters (A-Z, a-z,) numbers (0-9,) dashes, and underscores only.",
											Computed:    true,
										},
										"value": schema.ListAttribute{
											Description: "Header value.\n\nRestrictions:\n- Maximum 512 symbols.\n- Letters (a-z), numbers (0-9), spaces, and symbols (`~!@#%%^&*()-_=+ /|\\\";:?.,><{}[]).\n- Must start with a letter, number, asterisk or {.\n- Multiple values can be added.",
											Computed:    true,
											CustomType:  customfield.NewListType[types.String](ctx),
											ElementType: types.StringType,
										},
										"always": schema.BoolAttribute{
											Description: "Defines whether the header will be added to a response from CDN regardless of response code.\n\nPossible values:\n- **true** - Header will be added to a response from CDN regardless of response code.\n- **false** - Header will be added only to the following response codes: 200, 201, 204, 206, 301, 302, 303, 304, 307, 308.",
											Computed:    true,
										},
									},
								},
							},
						},
					},
					"static_headers": schema.SingleNestedAttribute{
						Description:        "**Legacy option**. Use the `static_response_headers` option instead.\n\nCustom HTTP Headers that a CDN server adds to response. Up to fifty custom HTTP Headers can be specified. May contain a header with multiple values.",
						Computed:           true,
						DeprecationMessage: "This attribute is deprecated.",
						CustomType:         customfield.NewNestedObjectType[CDNRuleTemplateOptionsStaticHeadersDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.StringAttribute{
								Description: "A MAP for static headers in a format of `header_name: header_value`.\n\nRestrictions:\n- **Header name** - Maximum 128 symbols, may contain Latin letters (A-Z, a-z), numbers (0-9), dashes, and underscores.\n- **Header value** - Maximum 512 symbols, may contain letters (a-z), numbers (0-9), spaces, and symbols (`~!@#%%^&*()-_=+ /|\\\";:?.,><{}[]). Must start with a letter, number, asterisk or {.",
								Computed:    true,
								CustomType:  jsontypes.NormalizedType{},
							},
						},
					},
					"static_request_headers": schema.SingleNestedAttribute{
						Description: "Custom HTTP Headers for a CDN server to add to request. Up to fifty custom HTTP Headers can be specified.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsStaticRequestHeadersDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.MapAttribute{
								Description: "A MAP for static headers in a format of `header_name: header_value`.\n\nRestrictions:\n- **Header name** - Maximum 255 symbols, may contain Latin letters (A-Z, a-z), numbers (0-9), dashes, and underscores.\n- **Header value** - Maximum 512 symbols, may contain letters (a-z), numbers (0-9), spaces, and symbols (`~!@#%%^&*()-_=+ /|\\\";:?.,><{}[]). Must start with a letter, number, asterisk or {.",
								Computed:    true,
								CustomType:  customfield.NewMapType[types.String](ctx),
								ElementType: types.StringType,
							},
						},
					},
					"user_agent_acl": schema.SingleNestedAttribute{
						Description: "Controls access to the content for specified User-Agents.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsUserAgentACLDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"excepted_values": schema.SetAttribute{
								Description: "List of User-Agents that will be allowed/denied.\n\nThe meaning of the parameter depends on `policy_type`:\n- **allow** - List of User-Agents for which access is prohibited.\n- **deny** - List of User-Agents for which access is allowed.\n\nYou can provide exact User-Agent strings or regular expressions. Regular expressions must start\nwith `~` (case-sensitive) or `~*` (case-insensitive).\n\nUse an empty string `\"\"` to allow/deny access when the User-Agent header is empty.",
								Computed:    true,
								CustomType:  customfield.NewSetType[types.String](ctx),
								ElementType: types.StringType,
							},
							"policy_type": schema.StringAttribute{
								Description: "User-Agents policy type.\n\nPossible values:\n- **allow** - Allow access for all User-Agents except specified in `excepted_values` field.\n- **deny** - Deny access for all User-Agents except specified in `excepted_values` field.\nAvailable values: \"allow\", \"deny\".",
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("allow", "deny"),
								},
							},
						},
					},
					"waap": schema.SingleNestedAttribute{
						Description: "Allows to enable WAAP (Web Application and API Protection).",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsWaapDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
					"websockets": schema.SingleNestedAttribute{
						Description: "Enables or disables WebSockets connections to an origin server.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CDNRuleTemplateOptionsWebsocketsDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Controls the option state.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
							"value": schema.BoolAttribute{
								Description: "Possible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

func (d *CDNRuleTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CDNRuleTemplateDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

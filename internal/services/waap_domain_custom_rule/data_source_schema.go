// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_custom_rule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainCustomRuleDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"rule_id": schema.Int64Attribute{
				Description: "The custom rule ID",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description assigned to the rule",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether or not the rule is enabled",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "The unique identifier for the rule",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name assigned to the rule",
				Computed:    true,
			},
			"action": schema.SingleNestedAttribute{
				Description: "The action that a WAAP rule takes when triggered",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleActionDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"allow": schema.StringAttribute{
						Description: "The WAAP allowed the request",
						Computed:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"block": schema.SingleNestedAttribute{
						Description: "WAAP block action behavior could be configured with response status code and action duration.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleActionBlockDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"action_duration": schema.StringAttribute{
								Description: "How long a rule's block action will apply to subsequent requests. Can be specified in seconds or by using a numeral followed by 's', 'm', 'h', or 'd' to represent time format (seconds, minutes, hours, or days)",
								Computed:    true,
							},
							"status_code": schema.Int64Attribute{
								Description: "Designates the HTTP status code to deliver when a request is blocked.\nAvailable values: 403, 405, 418, 429.",
								Computed:    true,
								Validators: []validator.Int64{
									int64validator.OneOf(
										403,
										405,
										418,
										429,
									),
								},
							},
						},
					},
					"captcha": schema.StringAttribute{
						Description: "The WAAP presented the user with a captcha",
						Computed:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"handshake": schema.StringAttribute{
						Description: "The WAAP performed automatic browser validation",
						Computed:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"monitor": schema.StringAttribute{
						Description: "The WAAP monitored the request but took no action",
						Computed:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"tag": schema.SingleNestedAttribute{
						Description: "WAAP tag action gets a list of tags to tag the request scope with",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleActionTagDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"tags": schema.ListAttribute{
								Description: "The list of user defined tags to tag the request with",
								Computed:    true,
								CustomType:  customfield.NewListType[types.String](ctx),
								ElementType: types.StringType,
							},
						},
					},
				},
			},
			"conditions": schema.ListNestedAttribute{
				Description: "The conditions required for the WAAP engine to trigger the rule. Rules may have between 1 and 5 conditions. All conditions must pass for the rule to trigger",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[WaapDomainCustomRuleConditionsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"content_type": schema.SingleNestedAttribute{
							Description: "Match the requested Content-Type",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsContentTypeDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"content_type": schema.ListAttribute{
									Description: "The list of content types to match against",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"country": schema.SingleNestedAttribute{
							Description: "Match the country that the request originated from",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsCountryDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"country_code": schema.ListAttribute{
									Description: "A list of ISO 3166-1 alpha-2 formatted strings representing the countries to match against",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"file_extension": schema.SingleNestedAttribute{
							Description: "Match the incoming file extension",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsFileExtensionDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"file_extension": schema.ListAttribute{
									Description: "The list of file extensions to match against",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"header": schema.SingleNestedAttribute{
							Description: "Match an incoming request header",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsHeaderDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"header": schema.StringAttribute{
									Description: "The request header name",
									Computed:    true,
								},
								"value": schema.StringAttribute{
									Description: "The request header value",
									Computed:    true,
								},
								"match_type": schema.StringAttribute{
									Description: "The type of matching condition for header and value.\nAvailable values: \"Exact\", \"Contains\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("Exact", "Contains"),
									},
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"header_exists": schema.SingleNestedAttribute{
							Description: "Match when an incoming request header is present",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsHeaderExistsDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"header": schema.StringAttribute{
									Description: "The request header name",
									Computed:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"http_method": schema.SingleNestedAttribute{
							Description: "Match the incoming HTTP method",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsHTTPMethodDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"http_method": schema.StringAttribute{
									Description: "HTTP methods and descriptions\nMethods from the following RFCs are all observed:\n\\* RFC 7231: Hypertext Transfer Protocol (HTTP/1.1), obsoletes 2616\n\\* RFC 5789: PATCH Method for HTTP\nAvailable values: \"CONNECT\", \"DELETE\", \"GET\", \"HEAD\", \"OPTIONS\", \"PATCH\", \"POST\", \"PUT\", \"TRACE\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive(
											"CONNECT",
											"DELETE",
											"GET",
											"HEAD",
											"OPTIONS",
											"PATCH",
											"POST",
											"PUT",
											"TRACE",
										),
									},
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"ip": schema.SingleNestedAttribute{
							Description: "Match the incoming request against a single IP address",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsIPDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"ip_address": schema.StringAttribute{
									Description: "A single IPv4 or IPv6 address",
									Computed:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"ip_range": schema.SingleNestedAttribute{
							Description: "Match the incoming request against an IP range",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsIPRangeDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"lower_bound": schema.StringAttribute{
									Description: "The lower bound IPv4 or IPv6 address to match against",
									Computed:    true,
								},
								"upper_bound": schema.StringAttribute{
									Description: "The upper bound IPv4 or IPv6 address to match against",
									Computed:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"organization": schema.SingleNestedAttribute{
							Description: "Match the organization the request originated from, as determined by a WHOIS lookup of the requesting IP",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsOrganizationDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"organization": schema.StringAttribute{
									Description: "The organization to match against",
									Computed:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"owner_types": schema.SingleNestedAttribute{
							Description: "Match the type of organization that owns the IP address making an incoming request",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsOwnerTypesDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
								"owner_types": schema.ListAttribute{
									Description: "Match the type of organization that owns the IP address making an incoming request",
									Computed:    true,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(
											stringvalidator.OneOfCaseInsensitive(
												"COMMERCIAL",
												"EDUCATIONAL",
												"GOVERNMENT",
												"HOSTING_SERVICES",
												"ISP",
												"MOBILE_NETWORK",
												"NETWORK",
												"RESERVED",
											),
										),
									},
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
							},
						},
						"request_rate": schema.SingleNestedAttribute{
							Description: "Match the rate at which requests come in that match certain conditions",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsRequestRateDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"path_pattern": schema.StringAttribute{
									Description: "A regular expression matching the URL path of the incoming request",
									Computed:    true,
								},
								"requests": schema.Int64Attribute{
									Description: "The number of incoming requests over the given time that can trigger a request rate condition",
									Computed:    true,
									Validators: []validator.Int64{
										int64validator.AtLeast(20),
									},
								},
								"time": schema.Int64Attribute{
									Description: "The number of seconds that the WAAP measures incoming requests over before triggering a request rate condition",
									Computed:    true,
									Validators: []validator.Int64{
										int64validator.Between(1, 3600),
									},
								},
								"http_methods": schema.ListAttribute{
									Description: "Possible HTTP request methods that can trigger a request rate condition",
									Computed:    true,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(
											stringvalidator.OneOfCaseInsensitive(
												"CONNECT",
												"DELETE",
												"GET",
												"HEAD",
												"OPTIONS",
												"PATCH",
												"POST",
												"PUT",
												"TRACE",
											),
										),
									},
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"ips": schema.ListAttribute{
									Description: "A list of source IPs that can trigger a request rate condition",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"user_defined_tag": schema.StringAttribute{
									Description: "A user-defined tag that can be included in incoming requests and used to trigger a request rate condition",
									Computed:    true,
								},
							},
						},
						"response_header": schema.SingleNestedAttribute{
							Description: "Match a response header",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsResponseHeaderDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"header": schema.StringAttribute{
									Description: "The response header name",
									Computed:    true,
								},
								"value": schema.StringAttribute{
									Description: "The response header value",
									Computed:    true,
								},
								"match_type": schema.StringAttribute{
									Description: "The type of matching condition for header and value.\nAvailable values: \"Exact\", \"Contains\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("Exact", "Contains"),
									},
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"response_header_exists": schema.SingleNestedAttribute{
							Description: "Match when a response header is present",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsResponseHeaderExistsDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"header": schema.StringAttribute{
									Description: "The response header name",
									Computed:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"session_request_count": schema.SingleNestedAttribute{
							Description: "Match the number of dynamic page requests made in a WAAP session",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsSessionRequestCountDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"request_count": schema.Int64Attribute{
									Description: "The number of dynamic requests in the session",
									Computed:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"tags": schema.SingleNestedAttribute{
							Description: "Matches requests based on specified tags",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsTagsDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"tags": schema.ListAttribute{
									Description: "A list of tags to match against the request tags",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"url": schema.SingleNestedAttribute{
							Description: "Match the incoming request URL",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsURLDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"url": schema.StringAttribute{
									Description: "The pattern to match against the request URL.\nConstraints depend on `match_type`:\n- **Exact/Contains**: plain text matching (e.g., `/admin`, must comply with `^[\\w!\\$~:#\\[\\]@\\(\\)\\\\*\\+,=\\/\\-\\.\\%]+$`).\n- **Regex**: a valid regular expression (e.g., `^/upload(/\\d+)?/\\w+`). Lookahead/lookbehind constructs are forbidden.",
									Computed:    true,
								},
								"match_type": schema.StringAttribute{
									Description: "The type of matching condition.\nAvailable values: \"Exact\", \"Contains\", \"Regex\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive(
											"Exact",
											"Contains",
											"Regex",
										),
									},
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"user_agent": schema.SingleNestedAttribute{
							Description: "Match the user agent making the request",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsUserAgentDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"user_agent": schema.StringAttribute{
									Description: "The user agent value to match",
									Computed:    true,
								},
								"match_type": schema.StringAttribute{
									Description: "The type of matching condition.\nAvailable values: \"Exact\", \"Contains\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("Exact", "Contains"),
									},
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
						"user_defined_tags": schema.SingleNestedAttribute{
							Description: "Matches requests based on user-defined tags",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainCustomRuleConditionsUserDefinedTagsDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"tags": schema.ListAttribute{
									Description: "A list of user-defined tags to match against the request tags",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *WaapDomainCustomRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapDomainCustomRuleDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_custom_rule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*WaapDomainCustomRuleResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "The unique identifier for the rule",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"domain_id": schema.Int64Attribute{
				Description:   "The domain ID",
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether or not the rule is enabled",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name assigned to the rule",
				Required:    true,
			},
			"action": schema.SingleNestedAttribute{
				Description: "The action that the rule takes when triggered. Only one action can be set per rule.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"allow": schema.StringAttribute{
						Description: "The WAAP allowed the request",
						Optional:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"block": schema.SingleNestedAttribute{
						Description: "WAAP block action behavior could be configured with response status code and action duration.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"action_duration": schema.StringAttribute{
								Description: "How long a rule's block action will apply to subsequent requests. Can be specified in seconds or by using a numeral followed by 's', 'm', 'h', or 'd' to represent time format (seconds, minutes, hours, or days). Empty time intervals are not allowed.",
								Optional:    true,
							},
							"status_code": schema.Int64Attribute{
								Description: "A custom HTTP status code that the WAAP returns if a rule blocks a request\nAvailable values: 403, 405, 418, 429.",
								Optional:    true,
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
						Optional:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"handshake": schema.StringAttribute{
						Description: "The WAAP performed automatic browser validation",
						Optional:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"monitor": schema.StringAttribute{
						Description: "The WAAP monitored the request but took no action",
						Optional:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"tag": schema.SingleNestedAttribute{
						Description: "WAAP tag action gets a list of tags to tag the request scope with",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"tags": schema.ListAttribute{
								Description: "The list of user defined tags to tag the request with",
								Required:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
			},
			"conditions": schema.ListNestedAttribute{
				Description: "The conditions required for the WAAP engine to trigger the rule. Rules may have between 1 and 5 conditions. All conditions must pass for the rule to trigger",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"content_type": schema.SingleNestedAttribute{
							Description: "Match the requested Content-Type",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"content_type": schema.ListAttribute{
									Description: "The list of content types to match against",
									Required:    true,
									ElementType: types.StringType,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"country": schema.SingleNestedAttribute{
							Description: "Match the country that the request originated from",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"country_code": schema.ListAttribute{
									Description: "A list of ISO 3166-1 alpha-2 formatted strings representing the countries to match against",
									Required:    true,
									ElementType: types.StringType,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"file_extension": schema.SingleNestedAttribute{
							Description: "Match the incoming file extension",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"file_extension": schema.ListAttribute{
									Description: "The list of file extensions to match against",
									Required:    true,
									ElementType: types.StringType,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"header": schema.SingleNestedAttribute{
							Description: "Match an incoming request header",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"header": schema.StringAttribute{
									Description: "The request header name",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "The request header value",
									Required:    true,
								},
								"match_type": schema.StringAttribute{
									Description: "The type of matching condition for header and value.\nAvailable values: \"Exact\", \"Contains\".",
									Computed:    true,
									Optional:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("Exact", "Contains"),
									},
									Default: stringdefault.StaticString("Contains"),
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"header_exists": schema.SingleNestedAttribute{
							Description: "Match when an incoming request header is present",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"header": schema.StringAttribute{
									Description: "The request header name",
									Required:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"http_method": schema.SingleNestedAttribute{
							Description: "Match the incoming HTTP method",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"http_method": schema.StringAttribute{
									Description: "HTTP methods of a request\nAvailable values: \"CONNECT\", \"DELETE\", \"GET\", \"HEAD\", \"OPTIONS\", \"PATCH\", \"POST\", \"PUT\", \"TRACE\".",
									Required:    true,
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
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"ip": schema.SingleNestedAttribute{
							Description: "Match the incoming request against a single IP address",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"ip_address": schema.StringAttribute{
									Description: "A single IPv4 or IPv6 address",
									Required:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"ip_range": schema.SingleNestedAttribute{
							Description: "Match the incoming request against an IP range",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"lower_bound": schema.StringAttribute{
									Description: "The lower bound IPv4 or IPv6 address to match against",
									Required:    true,
								},
								"upper_bound": schema.StringAttribute{
									Description: "The upper bound IPv4 or IPv6 address to match against",
									Required:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"organization": schema.SingleNestedAttribute{
							Description: "Match the organization the request originated from, as determined by a WHOIS lookup of the requesting IP",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"organization": schema.StringAttribute{
									Description: "The organization to match against",
									Required:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"owner_types": schema.SingleNestedAttribute{
							Description: "Match the type of organization that owns the IP address making an incoming request",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
								"owner_types": schema.ListAttribute{
									Description: "Match the type of organization that owns the IP address making an incoming request",
									Computed:    true,
									Optional:    true,
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
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"path_pattern": schema.StringAttribute{
									Description: "A regular expression matching the URL path of the incoming request",
									Required:    true,
								},
								"requests": schema.Int64Attribute{
									Description: "The number of incoming requests over the given time that can trigger a request rate condition",
									Required:    true,
									Validators: []validator.Int64{
										int64validator.AtLeast(20),
									},
								},
								"time": schema.Int64Attribute{
									Description: "The number of seconds that the WAAP measures incoming requests over before triggering a request rate condition",
									Required:    true,
									Validators: []validator.Int64{
										int64validator.Between(1, 3600),
									},
								},
								"http_methods": schema.ListAttribute{
									Description: "Possible HTTP request methods that can trigger a request rate condition",
									Optional:    true,
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
									ElementType: types.StringType,
								},
								"ips": schema.ListAttribute{
									Description: "A list of source IPs that can trigger a request rate condition",
									Optional:    true,
									ElementType: types.StringType,
								},
								"user_defined_tag": schema.StringAttribute{
									Description: "A user-defined tag that can be included in incoming requests and used to trigger a request rate condition",
									Optional:    true,
								},
							},
						},
						"response_header": schema.SingleNestedAttribute{
							Description: "Match a response header",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"header": schema.StringAttribute{
									Description: "The response header name",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "The response header value",
									Required:    true,
								},
								"match_type": schema.StringAttribute{
									Description: "The type of matching condition for header and value.\nAvailable values: \"Exact\", \"Contains\".",
									Computed:    true,
									Optional:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("Exact", "Contains"),
									},
									Default: stringdefault.StaticString("Contains"),
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"response_header_exists": schema.SingleNestedAttribute{
							Description: "Match when a response header is present",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"header": schema.StringAttribute{
									Description: "The response header name",
									Required:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"session_request_count": schema.SingleNestedAttribute{
							Description: "Match the number of dynamic page requests made in a WAAP session",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"request_count": schema.Int64Attribute{
									Description: "The number of dynamic requests in the session",
									Required:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"tags": schema.SingleNestedAttribute{
							Description: "Matches requests based on specified tags",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"tags": schema.ListAttribute{
									Description: "A list of tags to match against the request tags",
									Required:    true,
									ElementType: types.StringType,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"url": schema.SingleNestedAttribute{
							Description: "Match the incoming request URL",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"url": schema.StringAttribute{
									Description: "The pattern to match against the request URL.\nConstraints depend on `match_type`:\n\n- **Exact/Contains**: plain text matching (e.g., `/admin`, must comply with `^[\\w!\\$~:#\\[\\]@\\(\\)*\\+,=\\/\\-\\.\\%]+$`).\n- **Regex**: a valid regular expression (e.g., `^/upload(/\\d+)?/\\w+`). Lookahead/lookbehind constructs are forbidden.",
									Required:    true,
								},
								"match_type": schema.StringAttribute{
									Description: "The type of matching condition.\nAvailable values: \"Exact\", \"Contains\", \"Regex\".",
									Computed:    true,
									Optional:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive(
											"Exact",
											"Contains",
											"Regex",
										),
									},
									Default: stringdefault.StaticString("Contains"),
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"user_agent": schema.SingleNestedAttribute{
							Description: "Match the user agent making the request",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"user_agent": schema.StringAttribute{
									Description: "The user agent value to match",
									Required:    true,
								},
								"match_type": schema.StringAttribute{
									Description: "The type of matching condition.\nAvailable values: \"Exact\", \"Contains\".",
									Computed:    true,
									Optional:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("Exact", "Contains"),
									},
									Default: stringdefault.StaticString("Contains"),
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"user_defined_tags": schema.SingleNestedAttribute{
							Description: "Matches requests based on user-defined tags",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"tags": schema.ListAttribute{
									Description: "A list of user-defined tags to match against the request tags",
									Required:    true,
									ElementType: types.StringType,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
					},
				},
			},
			"description": schema.StringAttribute{
				Description: "The description assigned to the rule",
				Optional:    true,
			},
		},
	}
}

func (r *WaapDomainCustomRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *WaapDomainCustomRuleResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

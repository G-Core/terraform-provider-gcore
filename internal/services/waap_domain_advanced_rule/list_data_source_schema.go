// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_advanced_rule

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

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainAdvancedRulesDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Advanced rules are CEL-expression WAAP rules that allow, block, challenge, monitor or tag requests to a domain.",
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"action": schema.StringAttribute{
				Description: "Filter to refine results by specific actions\nAvailable values: \"allow\", \"block\", \"captcha\", \"handshake\", \"monitor\", \"tag\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"allow",
						"block",
						"captcha",
						"handshake",
						"monitor",
						"tag",
					),
				},
			},
			"description": schema.StringAttribute{
				Description: "Filter rules based on their description. Supports '*' as a wildcard character.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Filter rules based on their active status",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Filter rules based on their name. Supports '*' as a wildcard character.",
				Optional:    true,
			},
			"ordering": schema.StringAttribute{
				Description: "Determine the field to order results by\nAvailable values: \"id\", \"name\", \"description\", \"enabled\", \"action\", \"phase\", \"-id\", \"-name\", \"-description\", \"-enabled\", \"-action\", \"-phase\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"id",
						"name",
						"description",
						"enabled",
						"action",
						"phase",
						"-id",
						"-name",
						"-description",
						"-enabled",
						"-action",
						"-phase",
					),
				},
			},
			"phase": schema.StringAttribute{
				Description: "Filter rules based on the WAAP request/response phase for applying the rule.\n\n\nThe \"access\" phase is responsible for modifying the request before it is sent to the origin server. \n\nThe \"header_filter\" phase is responsible for modifying the HTTP headers of a response before they are sent back to the client.\n\nThe \"body_filter\" phase is responsible for modifying the body of a response before it is sent back to the client.\nAvailable values: \"access\", \"header_filter\", \"body_filter\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"access",
						"header_filter",
						"body_filter",
					),
				},
			},
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[WaapDomainAdvancedRulesItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The unique identifier for the rule",
							Computed:    true,
						},
						"action": schema.SingleNestedAttribute{
							Description: "The action that the rule takes when triggered. Only one action can be set per rule.",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainAdvancedRulesActionDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"allow": schema.MapAttribute{
									Description: "The WAAP allows the request",
									Computed:    true,
									CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
									ElementType: jsontypes.NormalizedType{},
								},
								"block": schema.SingleNestedAttribute{
									Description: "WAAP block action behavior could be configured with response status code and action duration.",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectType[WaapDomainAdvancedRulesActionBlockDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"action_duration": schema.StringAttribute{
											Description: "How long a rule's block action will apply to subsequent requests. Can be specified in seconds or by using a numeral followed by 's', 'm', 'h', or 'd' to represent time format (seconds, minutes, hours, or days). Empty time intervals are not allowed. If not specified, the action is applied only to the request that triggers the rule.",
											Computed:    true,
										},
										"status_code": schema.Int64Attribute{
											Description: "An HTTP status code that the WAAP returns if a rule blocks a request.\nAvailable values: 403, 405, 418, 429.",
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
								"captcha": schema.MapAttribute{
									Description: "The WAAP presents the user with a captcha",
									Computed:    true,
									CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
									ElementType: jsontypes.NormalizedType{},
								},
								"handshake": schema.MapAttribute{
									Description: "The WAAP performs automatic browser validation",
									Computed:    true,
									CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
									ElementType: jsontypes.NormalizedType{},
								},
								"monitor": schema.MapAttribute{
									Description: "The WAAP monitors the request but takes no action",
									Computed:    true,
									CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
									ElementType: jsontypes.NormalizedType{},
								},
								"tag": schema.SingleNestedAttribute{
									Description: "WAAP tag action gets a list of tags to tag the request scope with",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectType[WaapDomainAdvancedRulesActionTagDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"tags": schema.SetAttribute{
											Description: "The list of user defined tags to tag the request with",
											Computed:    true,
											CustomType:  customfield.NewSetType[types.String](ctx),
											ElementType: types.StringType,
										},
									},
								},
							},
						},
						"enabled": schema.BoolAttribute{
							Description: "Whether or not the rule is enabled",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name assigned to the rule",
							Computed:    true,
						},
						"source": schema.StringAttribute{
							Description: "A CEL syntax expression that contains the rule's conditions. Allowed objects are: request, whois, session, response, tags, `user_defined_tags`, `user_agent`, `client_data`.\n\nMore info can be found here: https://gcore.com/docs/waap/waap-rules/advanced-rules",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "A human-readable description.",
							Computed:    true,
						},
						"phase": schema.StringAttribute{
							Description: "The WAAP request/response phase for applying the rule. Default is \"access\".\n\n\nThe \"access\" phase is responsible for modifying the request before it is sent to the origin server. \n\nThe \"header_filter\" phase is responsible for modifying the HTTP headers of a response before they are sent back to the client.\n\nThe \"body_filter\" phase is responsible for modifying the body of a response before it is sent back to the client.\nAvailable values: \"access\", \"header_filter\", \"body_filter\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"access",
									"header_filter",
									"body_filter",
								),
							},
						},
					},
				},
			},
		},
	}
}

func (d *WaapDomainAdvancedRulesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *WaapDomainAdvancedRulesDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

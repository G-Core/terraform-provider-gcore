// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_advanced_rule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainAdvancedRuleDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The advanced rule ID",
				Computed:    true,
			},
			"rule_id": schema.Int64Attribute{
				Description: "The advanced rule ID",
				Optional:    true,
			},
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
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
			"name": schema.StringAttribute{
				Description: "The name assigned to the rule",
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
			"source": schema.StringAttribute{
				Description: "A CEL syntax expression that contains the rule's conditions. Allowed objects are: request, whois, session, response, tags, `user_defined_tags`, `user_agent`, `client_data`.\n\nMore info can be found here: https://gcore.com/docs/waap/waap-rules/advanced-rules",
				Computed:    true,
			},
			"action": schema.SingleNestedAttribute{
				Description: "The action that the rule takes when triggered. Only one action can be set per rule.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[WaapDomainAdvancedRuleActionDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"allow": schema.StringAttribute{
						Description: "The WAAP allowed the request",
						Computed:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"block": schema.SingleNestedAttribute{
						Description: "WAAP block action behavior could be configured with response status code and action duration.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[WaapDomainAdvancedRuleActionBlockDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"action_duration": schema.StringAttribute{
								Description: "How long a rule's block action will apply to subsequent requests. Can be specified in seconds or by using a numeral followed by 's', 'm', 'h', or 'd' to represent time format (seconds, minutes, hours, or days). Empty time intervals are not allowed.",
								Computed:    true,
							},
							"status_code": schema.Int64Attribute{
								Description: "A custom HTTP status code that the WAAP returns if a rule blocks a request\nAvailable values: 403, 405, 418, 429.",
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
						CustomType:  customfield.NewNestedObjectType[WaapDomainAdvancedRuleActionTagDataSourceModel](ctx),
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
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
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
				},
			},
		},
	}
}

func (d *WaapDomainAdvancedRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapDomainAdvancedRuleDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("rule_id"), path.MatchRoot("find_one_by")),
	}
}

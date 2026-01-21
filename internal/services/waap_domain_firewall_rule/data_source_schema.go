// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_firewall_rule

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
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainFirewallRuleDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The firewall rule ID",
				Computed:    true,
			},
			"rule_id": schema.Int64Attribute{
				Description: "The firewall rule ID",
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
			"action": schema.SingleNestedAttribute{
				Description: "The action that the rule takes when triggered",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[WaapDomainFirewallRuleActionDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"allow": schema.StringAttribute{
						Description: "The WAAP allowed the request",
						Computed:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"block": schema.SingleNestedAttribute{
						Description: "WAAP block action behavior could be configured with response status code and action duration.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[WaapDomainFirewallRuleActionBlockDataSourceModel](ctx),
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
				},
			},
			"conditions": schema.ListNestedAttribute{
				Description: "The condition required for the WAAP engine to trigger the rule.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[WaapDomainFirewallRuleConditionsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.SingleNestedAttribute{
							Description: "Match the incoming request against a single IP address",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[WaapDomainFirewallRuleConditionsIPDataSourceModel](ctx),
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
							CustomType:  customfield.NewNestedObjectType[WaapDomainFirewallRuleConditionsIPRangeDataSourceModel](ctx),
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
					},
				},
			},
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"action": schema.StringAttribute{
						Description: "Filter to refine results by specific firewall actions\nAvailable values: \"allow\", \"block\".",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("allow", "block"),
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
						Description: "Determine the field to order results by\nAvailable values: \"id\", \"name\", \"description\", \"enabled\", \"action\", \"-id\", \"-name\", \"-description\", \"-enabled\", \"-action\".",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"id",
								"name",
								"description",
								"enabled",
								"action",
								"-id",
								"-name",
								"-description",
								"-enabled",
								"-action",
							),
						},
					},
				},
			},
		},
	}
}

func (d *WaapDomainFirewallRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapDomainFirewallRuleDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("rule_id"), path.MatchRoot("find_one_by")),
	}
}

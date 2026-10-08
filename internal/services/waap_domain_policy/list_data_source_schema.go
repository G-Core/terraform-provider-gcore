// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_policy

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainPoliciesDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "WAAP security policies are the built-in protection rules of a domain. Each policy can be turned on or off.",
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"rule_set_id": schema.Int64Attribute{
				Description: "Filter policies by rule set ID",
				Optional:    true,
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
				CustomType:  customfield.NewNestedObjectListType[WaapDomainPoliciesItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier for the security rule",
							Computed:    true,
						},
						"action": schema.StringAttribute{
							Description: "Specifies the action taken by the WAAP upon rule activation\nAvailable values: \"Allow\", \"Block\", \"Captcha\", \"Gateway\", \"Handshake\", \"Monitor\", \"Composite\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"Allow",
									"Block",
									"Captcha",
									"Gateway",
									"Handshake",
									"Monitor",
									"Composite",
								),
							},
						},
						"description": schema.StringAttribute{
							Description: "Detailed description of the security rule",
							Computed:    true,
						},
						"group": schema.StringAttribute{
							Description: "The rule set group name to which the rule belongs",
							Computed:    true,
						},
						"mode": schema.BoolAttribute{
							Description: "Indicates if the security rule is active",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the security rule",
							Computed:    true,
						},
						"rule_set_id": schema.Int64Attribute{
							Description: "Identifier of the rule set to which the rule belongs",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *WaapDomainPoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *WaapDomainPoliciesDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainPolicyDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "WAAP security policies are the built-in protection rules of a domain. Each policy can be turned on or off.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the policy to retrieve",
				Computed:    true,
			},
			"policy_id": schema.StringAttribute{
				Description: "The ID of the policy to retrieve",
				Required:    true,
			},
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"mode": schema.BoolAttribute{
				Description: "Indicates if the security rule is active",
				Computed:    true,
			},
		},
	}
}

func (d *WaapDomainPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapDomainPolicyDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

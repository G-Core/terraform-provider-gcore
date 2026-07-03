// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_origin_shielding

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*CDNOriginShieldingDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Origin shielding protects your origin server from being overloaded by routing all CDN requests through a single shield (precache) server for a CDN resource.",
		Attributes: map[string]schema.Attribute{
			"resource_id": schema.Int64Attribute{
				Required: true,
			},
			"shielding_pop": schema.Int64Attribute{
				Description: "Shielding location ID.\n\nIf origin shielding is disabled, the parameter value is **null**.",
				Computed:    true,
			},
		},
	}
}

func (d *CDNOriginShieldingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CDNOriginShieldingDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

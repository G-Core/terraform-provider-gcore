// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_secret

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*FastedgeSecretDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "FastEdge secrets store sensitive values such as API keys and tokens that can be referenced by FastEdge applications.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"secret_id": schema.Int64Attribute{
				Required: true,
			},
			"app_count": schema.Int64Attribute{
				Description: "The number of applications that use this secret.",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "A description or comment about the secret.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the secret.",
				Computed:    true,
			},
			"secret_slots": schema.SetNestedAttribute{
				Description: "A list of secret slots associated with this secret.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectSetType[FastedgeSecretSecretSlotsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"slot": schema.Int64Attribute{
							Description: "Unix timestamp (seconds since epoch) indicating when this secret version becomes active. Use for time-based secret rotation.",
							Computed:    true,
						},
						"checksum": schema.StringAttribute{
							Description: "SHA-256 hash of the decrypted value for integrity verification (auto-generated)",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "The plaintext secret value. Will be encrypted with AES-256-GCM before storage.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *FastedgeSecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *FastedgeSecretDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

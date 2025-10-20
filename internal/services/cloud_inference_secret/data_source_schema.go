// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_secret

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceSecretDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"secret_name": schema.StringAttribute{
				Description: "Inference secret name.",
				Required:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Secret name.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Secret type.",
				Computed:    true,
			},
			"data": schema.SingleNestedAttribute{
				Description: "Secret data.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudInferenceSecretDataDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"aws_access_key_id": schema.StringAttribute{
						Description: "AWS IAM key ID.",
						Computed:    true,
					},
					"aws_secret_access_key": schema.StringAttribute{
						Description: "AWS IAM secret key.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *CloudInferenceSecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudInferenceSecretDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

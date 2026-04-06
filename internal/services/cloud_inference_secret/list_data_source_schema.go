// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_secret

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceSecretsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Inference secrets store sensitive values such as AWS credentials used for SQS-based autoscaling triggers in deployments.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
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
				CustomType:  customfield.NewNestedObjectListType[CloudInferenceSecretsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Secret name.",
							Computed:    true,
						},
						"data": schema.SingleNestedAttribute{
							Description: "Secret data.",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudInferenceSecretsDataDataSourceModel](ctx),
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
						"name": schema.StringAttribute{
							Description: "Secret name.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Secret type.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudInferenceSecretsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudInferenceSecretsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

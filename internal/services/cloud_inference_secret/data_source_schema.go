// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_secret

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceSecretDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Inference secrets store sensitive values such as AWS credentials used for SQS-based autoscaling triggers in deployments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Inference secret name.",
				Computed:    true,
			},
			"secret_name": schema.StringAttribute{
				Description: "Inference secret name.",
				Optional:    true,
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
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"limit": schema.Int64Attribute{
						Description: "Optional. Limit the number of returned items",
						Computed:    true,
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.AtMost(1000),
						},
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
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("secret_name"), path.MatchRoot("find_one_by")),
	}
}

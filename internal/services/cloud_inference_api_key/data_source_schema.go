// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_api_key

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceAPIKeyDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key_name": schema.StringAttribute{
				Description: "Api key name.",
				Required:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the API Key was created.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the API Key.",
				Computed:    true,
			},
			"expires_at": schema.StringAttribute{
				Description: "Timestamp when the API Key will expire.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "API Key name.",
				Computed:    true,
			},
			"deployment_names": schema.ListAttribute{
				Description: "List of inference deployment names to which this API Key has been attached.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
		},
	}
}

func (d *CloudInferenceAPIKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudInferenceAPIKeyDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

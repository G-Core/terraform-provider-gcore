// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceModelDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"model_id": schema.StringAttribute{
				Description: "Model ID",
				Required:    true,
			},
			"category": schema.StringAttribute{
				Description: "Category of the model.",
				Computed:    true,
			},
			"default_flavor_name": schema.StringAttribute{
				Description: "Default flavor for the model.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the model.",
				Computed:    true,
			},
			"developer": schema.StringAttribute{
				Description: "Developer of the model.",
				Computed:    true,
			},
			"documentation_page": schema.StringAttribute{
				Description: "Path to the documentation page.",
				Computed:    true,
			},
			"eula_url": schema.StringAttribute{
				Description: "URL to the EULA text.",
				Computed:    true,
			},
			"example_curl_request": schema.StringAttribute{
				Description: "Example curl request to the model.",
				Computed:    true,
			},
			"has_eula": schema.BoolAttribute{
				Description: "Whether the model has an EULA.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "Model ID.",
				Computed:    true,
			},
			"image_registry_id": schema.StringAttribute{
				Description: "Image registry of the model.",
				Computed:    true,
			},
			"image_url": schema.StringAttribute{
				Description: "Image URL of the model.",
				Computed:    true,
			},
			"inference_backend": schema.StringAttribute{
				Description: "Describing underlying inference engine.",
				Computed:    true,
			},
			"inference_frontend": schema.StringAttribute{
				Description: "Describing model frontend type.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the model.",
				Computed:    true,
			},
			"openai_compatibility": schema.StringAttribute{
				Description: "OpenAI compatibility level.",
				Computed:    true,
			},
			"port": schema.Int64Attribute{
				Description: "Port on which the model runs.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "Version of the model.",
				Computed:    true,
			},
		},
	}
}

func (d *CloudInferenceModelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudInferenceModelDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

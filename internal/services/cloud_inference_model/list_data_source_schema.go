// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceModelsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"limit": schema.Int64Attribute{
				Description: "Optional. Limit the number of returned items",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtMost(1000),
				},
			},
			"order_by": schema.StringAttribute{
				Description: "Order instances by transmitted fields and directions\nAvailable values: \"name.asc\", \"name.desc\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("name.asc", "name.desc"),
				},
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
				CustomType:  customfield.NewNestedObjectListType[CloudInferenceModelsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Model ID.",
							Computed:    true,
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
						"model_id": schema.StringAttribute{
							Description: "Model name to perform inference call.",
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
				},
			},
		},
	}
}

func (d *CloudInferenceModelsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudInferenceModelsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

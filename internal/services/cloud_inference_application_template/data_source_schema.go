// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_application_template

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceApplicationTemplateDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"application_name": schema.StringAttribute{
				Description: "Name of application in catalog",
				Required:    true,
			},
			"cover_url": schema.StringAttribute{
				Description: "URL to the application's cover image",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Brief overview of the application",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "Human-readable name of the application",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Unique application identifier in the catalog",
				Computed:    true,
			},
			"readme": schema.StringAttribute{
				Description: "Detailed documentation or instructions",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Categorization key-value pairs",
				Computed:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"components": schema.MapNestedAttribute{
				Description: "Configurable components of the application",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectMapType[CloudInferenceApplicationTemplateComponentsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Description: "Summary of the component's functionality",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "Human-readable name of the component",
							Computed:    true,
						},
						"exposable": schema.BoolAttribute{
							Description: "Indicates whether this component can expose a public-facing endpoint (e.g., for inference or API access).",
							Computed:    true,
						},
						"license_url": schema.StringAttribute{
							Description: "URL to the component's license information",
							Computed:    true,
						},
						"parameters": schema.MapNestedAttribute{
							Description: "Configurable parameters for the component",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectMapType[CloudInferenceApplicationTemplateComponentsParametersDataSourceModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"default_value": schema.StringAttribute{
										Description: "Default value assigned if not provided",
										Computed:    true,
									},
									"description": schema.StringAttribute{
										Description: "Description of the parameter's purpose",
										Computed:    true,
									},
									"display_name": schema.StringAttribute{
										Description: "User-friendly name of the parameter",
										Computed:    true,
									},
									"enum_values": schema.ListAttribute{
										Description: `Allowed values when type is "enum"`,
										Computed:    true,
										CustomType:  customfield.NewListType[types.String](ctx),
										ElementType: types.StringType,
									},
									"max_value": schema.StringAttribute{
										Description: "Maximum value (applies to integer and float types)",
										Computed:    true,
									},
									"min_value": schema.StringAttribute{
										Description: "Minimum value (applies to integer and float types)",
										Computed:    true,
									},
									"pattern": schema.StringAttribute{
										Description: `Regexp pattern when type is "string"`,
										Computed:    true,
									},
									"required": schema.BoolAttribute{
										Description: "Indicates is parameter mandatory",
										Computed:    true,
									},
									"type": schema.StringAttribute{
										Description: "Determines the type of the parameter\nAvailable values: \"enum\", \"float\", \"integer\", \"string\".",
										Computed:    true,
										Validators: []validator.String{
											stringvalidator.OneOfCaseInsensitive(
												"enum",
												"float",
												"integer",
												"string",
											),
										},
									},
								},
							},
						},
						"readme": schema.StringAttribute{
							Description: "Detailed documentation or usage instructions",
							Computed:    true,
						},
						"required": schema.BoolAttribute{
							Description: "Specifies if the component is required for the application",
							Computed:    true,
						},
						"suitable_flavors": schema.ListNestedAttribute{
							Description: "List of compatible flavors or configurations",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudInferenceApplicationTemplateComponentsSuitableFlavorsDataSourceModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Name of the flavor",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *CloudInferenceApplicationTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudInferenceApplicationTemplateDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

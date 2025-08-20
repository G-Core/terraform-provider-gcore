// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_application_deployment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceApplicationDeploymentDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"deployment_name": schema.StringAttribute{
				Description: "Name of deployment",
				Required:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Required:    true,
			},
			"application_name": schema.StringAttribute{
				Description: "Identifier of the application template from the catalog",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Unique identifier of the deployment",
				Computed:    true,
			},
			"api_keys": schema.ListAttribute{
				Description: "List of API keys for the application",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"regions": schema.ListAttribute{
				Description: "Geographical regions where the deployment is active",
				Computed:    true,
				CustomType:  customfield.NewListType[types.Int64](ctx),
				ElementType: types.Int64Type,
			},
			"components_configuration": schema.MapNestedAttribute{
				Description: "Mapping of component names to their configuration (e.g., `\"model\": {...}`)",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectMapType[CloudInferenceApplicationDeploymentComponentsConfigurationDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"exposed": schema.BoolAttribute{
							Description: "Indicates if the component will obtain a public address",
							Computed:    true,
						},
						"flavor": schema.StringAttribute{
							Description: "Chosen flavor or variant of the component",
							Computed:    true,
						},
						"parameter_overrides": schema.MapNestedAttribute{
							Description: "Map of parameter overrides for customization",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectMapType[CloudInferenceApplicationDeploymentComponentsConfigurationParameterOverridesDataSourceModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"value": schema.StringAttribute{
										Description: "New value assigned to the overridden parameter",
										Computed:    true,
									},
								},
							},
						},
						"scale": schema.SingleNestedAttribute{
							Description: "Scaling parameters of the component",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudInferenceApplicationDeploymentComponentsConfigurationScaleDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"max": schema.Int64Attribute{
									Description: "Maximum number of replicas the container can be scaled up to",
									Computed:    true,
								},
								"min": schema.Int64Attribute{
									Description: "Minimum number of replicas the component can be scaled down to",
									Computed:    true,
								},
							},
						},
					},
				},
			},
			"status": schema.SingleNestedAttribute{
				Description: "Current state of the deployment across regions",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudInferenceApplicationDeploymentStatusDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"component_inferences": schema.MapNestedAttribute{
						Description: "Map of components and their inferences",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectMapType[CloudInferenceApplicationDeploymentStatusComponentInferencesDataSourceModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"flavor": schema.StringAttribute{
									Description: "Flavor of the inference",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Name of the inference",
									Computed:    true,
								},
							},
						},
					},
					"consolidated_status": schema.StringAttribute{
						Description: "High-level summary of the deployment status across all regions\nAvailable values: \"Active\", \"Failed\", \"PartiallyDeployed\", \"Unknown\".",
						Computed:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"Active",
								"Failed",
								"PartiallyDeployed",
								"Unknown",
							),
						},
					},
					"expose_addresses": schema.MapNestedAttribute{
						Description: "Map of component keys to their global access endpoints",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectMapType[CloudInferenceApplicationDeploymentStatusExposeAddressesDataSourceModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"address": schema.StringAttribute{
									Description: "Global access endpoint for the component",
									Computed:    true,
								},
							},
						},
					},
					"regions": schema.ListNestedAttribute{
						Description: "Status details for each deployment region",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectListType[CloudInferenceApplicationDeploymentStatusRegionsDataSourceModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"components": schema.MapNestedAttribute{
									Description: "Mapping of component names to their status in the region",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectMapType[CloudInferenceApplicationDeploymentStatusRegionsComponentsDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"error": schema.StringAttribute{
												Description: "Error message if the component is in an error state",
												Computed:    true,
											},
											"status": schema.StringAttribute{
												Description: "Current state of the component in a specific region",
												Computed:    true,
											},
										},
									},
								},
								"region_id": schema.Int64Attribute{
									Description: "Region ID",
									Computed:    true,
								},
								"status": schema.StringAttribute{
									Description: "Current state of the deployment in a specific region",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *CloudInferenceApplicationDeploymentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudInferenceApplicationDeploymentDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

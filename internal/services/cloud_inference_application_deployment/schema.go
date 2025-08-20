// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_application_deployment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*CloudInferenceApplicationDeploymentResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"deployment_name": schema.StringAttribute{
				Description:   "Name of deployment",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.Int64Attribute{
				Description:   "Project ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"application_name": schema.StringAttribute{
				Description:   "Identifier of the application from the catalog",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "Desired name for the new deployment",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"regions": schema.ListAttribute{
				Description:   "Geographical regions where the deployment should be created",
				Required:      true,
				ElementType:   types.Int64Type,
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"components_configuration": schema.MapNestedAttribute{
				Description: "Mapping of component names to their configuration (e.g., `\"model\": {...}`)",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"exposed": schema.BoolAttribute{
							Description: "Whether the component should be exposed via a public endpoint (e.g., for external inference/API access).",
							Required:    true,
						},
						"flavor": schema.StringAttribute{
							Description: "Specifies the compute configuration (e.g., CPU/GPU size) to be used for the component.",
							Required:    true,
						},
						"scale": schema.SingleNestedAttribute{
							Description: "Scaling parameters of the component",
							Required:    true,
							Attributes: map[string]schema.Attribute{
								"max": schema.Int64Attribute{
									Description: "Maximum number of replicas the container can be scaled up to",
									Required:    true,
								},
								"min": schema.Int64Attribute{
									Description: "Minimum number of replicas the component can be scaled down to",
									Required:    true,
								},
							},
						},
						"parameter_overrides": schema.MapNestedAttribute{
							Description: "Map of parameter overrides for customization",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"value": schema.StringAttribute{
										Description: "New value assigned to the overridden parameter",
										Required:    true,
									},
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Map{mapplanmodifier.RequiresReplace()},
			},
			"api_keys": schema.ListAttribute{
				Description:   "List of API keys for the application",
				Optional:      true,
				ElementType:   types.StringType,
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"tasks": schema.ListAttribute{
				Description: "List of task IDs",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"status": schema.SingleNestedAttribute{
				Description: "Current state of the deployment across regions",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudInferenceApplicationDeploymentStatusModel](ctx),
				Attributes: map[string]schema.Attribute{
					"component_inferences": schema.MapNestedAttribute{
						Description: "Map of components and their inferences",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectMapType[CloudInferenceApplicationDeploymentStatusComponentInferencesModel](ctx),
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
						CustomType:  customfield.NewNestedObjectMapType[CloudInferenceApplicationDeploymentStatusExposeAddressesModel](ctx),
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
						CustomType:  customfield.NewNestedObjectListType[CloudInferenceApplicationDeploymentStatusRegionsModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"components": schema.MapNestedAttribute{
									Description: "Mapping of component names to their status in the region",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectMapType[CloudInferenceApplicationDeploymentStatusRegionsComponentsModel](ctx),
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

func (r *CloudInferenceApplicationDeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudInferenceApplicationDeploymentResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

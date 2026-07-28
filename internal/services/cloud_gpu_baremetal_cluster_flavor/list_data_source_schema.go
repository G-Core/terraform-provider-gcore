// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster_flavor

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudGPUBaremetalClusterFlavorsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Optional:    true,
			},
			"hide_disabled": schema.BoolAttribute{
				Description: "Set to `true` to remove the disabled flavors from the response.",
				Computed:    true,
				Optional:    true,
			},
			"include_prices": schema.BoolAttribute{
				Description: "Set to `true` if the response should include flavor prices.",
				Computed:    true,
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
				CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterFlavorsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"architecture": schema.StringAttribute{
							Description: "Flavor architecture type",
							Computed:    true,
						},
						"capacity": schema.Int64Attribute{
							Description: "Number of available instances of given flavor for the client",
							Computed:    true,
						},
						"disabled": schema.BoolAttribute{
							Description: "If the flavor is disabled, new resources cannot be created using this flavor.",
							Computed:    true,
						},
						"hardware_description": schema.SingleNestedAttribute{
							Description: "Additional bare metal hardware description",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterFlavorsHardwareDescriptionDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"cpu": schema.StringAttribute{
									Description: "Human-readable CPU description",
									Computed:    true,
								},
								"disk": schema.StringAttribute{
									Description: "Human-readable disk description",
									Computed:    true,
								},
								"gpu": schema.StringAttribute{
									Description: "Human-readable GPU description",
									Computed:    true,
								},
								"network": schema.StringAttribute{
									Description: "Human-readable NIC description",
									Computed:    true,
								},
								"ram": schema.StringAttribute{
									Description: "Human-readable RAM description",
									Computed:    true,
								},
							},
						},
						"hardware_properties": schema.SingleNestedAttribute{
							Description: "Additional bare metal hardware properties",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterFlavorsHardwarePropertiesDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"gpu_count": schema.Int64Attribute{
									Description: "The total count of available GPUs.",
									Computed:    true,
								},
								"gpu_manufacturer": schema.StringAttribute{
									Description: "The manufacturer of the graphics processing GPU",
									Computed:    true,
								},
								"gpu_model": schema.StringAttribute{
									Description: "GPU model",
									Computed:    true,
								},
								"nic_eth": schema.StringAttribute{
									Description: "The configuration of the Ethernet ports",
									Computed:    true,
								},
								"nic_ib": schema.StringAttribute{
									Description: "The configuration of the InfiniBand ports",
									Computed:    true,
								},
							},
						},
						"name": schema.StringAttribute{
							Description: "Flavor name",
							Computed:    true,
						},
						"reserved_capacity": schema.Int64Attribute{
							Description: "Number of available instances of given flavor from reservations",
							Computed:    true,
						},
						"supported_features": schema.SingleNestedAttribute{
							Description: "Set of enabled features based on the flavor's type and configuration",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterFlavorsSupportedFeaturesDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"security_groups": schema.BoolAttribute{
									Computed: true,
								},
							},
						},
						"price": schema.SingleNestedAttribute{
							Description: "Flavor price",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterFlavorsPriceDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"currency_code": schema.StringAttribute{
									Description: "Currency code. Shown if the `include_prices` query parameter if set to true",
									Computed:    true,
								},
								"price_per_hour": schema.Float64Attribute{
									Description: "Price per hour. Shown if the `include_prices` query parameter if set to true",
									Computed:    true,
								},
								"price_per_month": schema.Float64Attribute{
									Description: "Price per month. Shown if the `include_prices` query parameter if set to true",
									Computed:    true,
								},
								"price_status": schema.StringAttribute{
									Description: "Price status for the UI\nAvailable values: \"error\", \"hide\", \"show\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive(
											"error",
											"hide",
											"show",
										),
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

func (d *CloudGPUBaremetalClusterFlavorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudGPUBaremetalClusterFlavorsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

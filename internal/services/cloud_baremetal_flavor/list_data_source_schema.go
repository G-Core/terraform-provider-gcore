// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_flavor

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudBaremetalFlavorsDataSource)(nil)

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
			"disabled": schema.BoolAttribute{
				Description: "Flag for filtering disabled flavors in the region",
				Computed:    true,
				Optional:    true,
			},
			"exclude_linux": schema.BoolAttribute{
				Description: "Set to true to exclude flavors dedicated to linux images",
				Computed:    true,
				Optional:    true,
			},
			"exclude_windows": schema.BoolAttribute{
				Description: "Set to true to exclude flavors dedicated to windows images",
				Computed:    true,
				Optional:    true,
			},
			"include_capacity": schema.BoolAttribute{
				Description: "Set to true if the response should include flavor capacity",
				Computed:    true,
				Optional:    true,
			},
			"include_prices": schema.BoolAttribute{
				Description: "Set to true if the response should include flavor prices",
				Computed:    true,
				Optional:    true,
			},
			"include_reservation_stock": schema.BoolAttribute{
				Description: "Optional. Set to true if flavor listing should include count of reserved resources in stock.",
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
				CustomType:  customfield.NewNestedObjectListType[CloudBaremetalFlavorsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"architecture": schema.StringAttribute{
							Description: "Flavor architecture type",
							Computed:    true,
						},
						"disabled": schema.BoolAttribute{
							Description: "Disabled flavor flag",
							Computed:    true,
						},
						"flavor_id": schema.StringAttribute{
							Description: "Flavor ID is the same as name",
							Computed:    true,
						},
						"flavor_name": schema.StringAttribute{
							Description: "Flavor name",
							Computed:    true,
						},
						"os_type": schema.StringAttribute{
							Description: "Flavor operating system",
							Computed:    true,
						},
						"ram": schema.Int64Attribute{
							Description: "RAM size in MiB",
							Computed:    true,
						},
						"resource_class": schema.StringAttribute{
							Description: "Flavor resource class for mapping to hardware capacity",
							Computed:    true,
						},
						"vcpus": schema.Int64Attribute{
							Description: "Virtual CPU count. For bare metal flavors, it's a physical CPU count",
							Computed:    true,
						},
						"capacity": schema.Int64Attribute{
							Description: "Number of available instances of given configuration",
							Computed:    true,
						},
						"currency_code": schema.StringAttribute{
							Description: "Currency code. Shown if the `include_prices` query parameter if set to true",
							Computed:    true,
						},
						"hardware_description": schema.MapAttribute{
							Description: "Additional hardware description",
							Computed:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
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
						"reserved_capacity": schema.Int64Attribute{
							Description: "Number of available instances of given flavor from reservations",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudBaremetalFlavorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudBaremetalFlavorsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

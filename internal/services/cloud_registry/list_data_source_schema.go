// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudRegistriesDataSource)(nil)

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
				CustomType:  customfield.NewNestedObjectListType[CloudRegistriesItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Registry ID",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Registry creation date-time",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"name": schema.StringAttribute{
							Description: "Registry name",
							Computed:    true,
						},
						"repo_count": schema.Int64Attribute{
							Description: "Number of repositories in the registry",
							Computed:    true,
						},
						"storage_limit": schema.Int64Attribute{
							Description: "Registry storage limit, GiB",
							Computed:    true,
						},
						"storage_used": schema.Int64Attribute{
							Description: "Registry storage used, bytes",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Registry modification date-time",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"url": schema.StringAttribute{
							Description: "Registry url",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudRegistriesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudRegistriesDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

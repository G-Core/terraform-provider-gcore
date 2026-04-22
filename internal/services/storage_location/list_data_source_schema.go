// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_location

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*StorageLocationsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Locations represent cloud regions where new storages can be created.",
		Attributes: map[string]schema.Attribute{
			"order_by": schema.StringAttribute{
				Computed: true,
				Optional: true,
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
				CustomType:  customfield.NewNestedObjectListType[StorageLocationsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Description: "Full hostname/address for accessing the storage endpoint.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Human-readable display name for the location.",
							Computed:    true,
						},
						"technical_name": schema.StringAttribute{
							Description: "Internal technical identifier for the location",
							Computed:    true,
						},
						"title": schema.StringAttribute{
							Description: "Display title for the location (English). Null if no title is set.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Storage type supported by this location\nAvailable values: \"s3_compatible\", \"sftp\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("s3_compatible", "sftp"),
							},
						},
					},
				},
			},
		},
	}
}

func (d *StorageLocationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *StorageLocationsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

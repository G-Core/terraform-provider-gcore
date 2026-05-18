// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*StorageObjectStoragesDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "S3-compatible object storages provide scalable cloud storage with S3 API compatibility. Each storage is provisioned in a specific location and exposes one or more access keys for authentication.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Filter by storage ID",
				Optional:    true,
			},
			"location_name": schema.StringAttribute{
				Description: "Filter by storage location/region",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Filter by storage name",
				Optional:    true,
			},
			"provisioning_status": schema.StringAttribute{
				Description: "Filter by provisioning status\nAvailable values: \"active\", \"creating\", \"updating\", \"deleting\", \"deleted\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"active",
						"creating",
						"updating",
						"deleting",
						"deleted",
					),
				},
			},
			"show_deleted": schema.BoolAttribute{
				Description: "Include deleted storages",
				Optional:    true,
			},
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
				CustomType:  customfield.NewNestedObjectListType[StorageObjectStoragesItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Unique identifier for the storage instance",
							Computed:    true,
						},
						"address": schema.StringAttribute{
							Description: "Full hostname/address for accessing the storage endpoint",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "ISO 8601 timestamp when the storage was created",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"full_name": schema.StringAttribute{
							Description: "Read-only internal full name of the storage, composed as \"{`client_id`}-{name}\".\nUsed internally by the backend. Clients should continue to identify the storage by `name`.",
							Computed:    true,
						},
						"location_name": schema.StringAttribute{
							Description: "Geographic location code where the storage is provisioned",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-defined name for the storage instance, as supplied at creation time.",
							Computed:    true,
						},
						"provisioning_status": schema.StringAttribute{
							Description: "Lifecycle status of the storage. Use this to check readiness before operations.\nAvailable values: \"creating\", \"active\", \"updating\", \"deleting\", \"deleted\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"creating",
									"active",
									"updating",
									"deleting",
									"deleted",
								),
							},
						},
					},
				},
			},
		},
	}
}

func (d *StorageObjectStoragesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *StorageObjectStoragesDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

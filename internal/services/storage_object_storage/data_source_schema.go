// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*StorageObjectStorageDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "S3-compatible object storages provide scalable cloud storage with S3 API compatibility. Each storage is provisioned in a specific location and exposes one or more access keys for authentication.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"storage_id": schema.Int64Attribute{
				Optional: true,
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
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
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
					"order_by": schema.StringAttribute{
						Computed: true,
						Optional: true,
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
				},
			},
		},
	}
}

func (d *StorageObjectStorageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *StorageObjectStorageDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("storage_id"), path.MatchRoot("find_one_by")),
	}
}

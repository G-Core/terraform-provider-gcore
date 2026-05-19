// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_sftp

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*StorageSftpsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "SFTP storages provide file transfer protocol access for securely uploading, downloading, and managing files over SSH.",
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
				CustomType:  customfield.NewNestedObjectListType[StorageSftpsItemsDataSourceModel](ctx),
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
						"expires": schema.StringAttribute{
							Description: "Duration when the storage will expire. Null if no expiration is set.",
							Computed:    true,
						},
						"full_name": schema.StringAttribute{
							Description: "Read-only internal full name of the storage, composed as \"{`client_id`}-{name}\".\nUsed by the SFTP backend as the login username. Clients should use this value when connecting\nbut should continue to identify the storage by `name` in their own configuration.",
							Computed:    true,
						},
						"has_custom_config_file": schema.BoolAttribute{
							Description: "Whether this storage uses a custom configuration file",
							Computed:    true,
						},
						"has_password": schema.BoolAttribute{
							Description: "Whether password authentication is configured for this storage",
							Computed:    true,
						},
						"is_http_disabled": schema.BoolAttribute{
							Description: "Whether HTTP access is disabled for this storage (HTTPS only)",
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
						"server_alias": schema.StringAttribute{
							Description: "Custom domain alias for accessing the storage. Null if no alias is configured.",
							Computed:    true,
						},
						"ssh_key_ids": schema.ListAttribute{
							Description: "IDs of SSH keys associated with this SFTP storage",
							Computed:    true,
							CustomType:  customfield.NewListType[types.Int64](ctx),
							ElementType: types.Int64Type,
						},
						"password": schema.StringAttribute{
							Description: "SFTP password. Only returned when newly generated or set (create/patch). Omitted in GET/list responses.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *StorageSftpsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *StorageSftpsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

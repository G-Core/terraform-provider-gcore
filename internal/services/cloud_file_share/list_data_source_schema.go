// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudFileSharesDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "File shares provide NFS-based shared storage that can be mounted by virtual machines and Kubernetes clusters for persistent data.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "File share name. Uses partial match.",
				Optional:    true,
			},
			"type_name": schema.StringAttribute{
				Description: "File share type name\nAvailable values: \"standard\", \"vast\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("standard", "vast"),
				},
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
				CustomType:  customfield.NewNestedObjectListType[CloudFileSharesItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "File share ID",
							Computed:    true,
						},
						"connection_point": schema.StringAttribute{
							Description: "Connection point. Can be null during File share creation",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Datetime when the file share was created",
							Computed:    true,
						},
						"creator_task_id": schema.StringAttribute{
							Description: "Task that created this entity",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "File share name",
							Computed:    true,
						},
						"network_id": schema.StringAttribute{
							Description: "Network ID.",
							Computed:    true,
						},
						"network_name": schema.StringAttribute{
							Description: "Network name.",
							Computed:    true,
						},
						"project_id": schema.Int64Attribute{
							Description: "Project ID",
							Computed:    true,
						},
						"protocol": schema.StringAttribute{
							Description: "File share protocol",
							Computed:    true,
						},
						"region": schema.StringAttribute{
							Description: "Region name",
							Computed:    true,
						},
						"region_id": schema.Int64Attribute{
							Description: "Region ID",
							Computed:    true,
						},
						"share_network_name": schema.StringAttribute{
							Description: "Share network name. May be null if the file share was created with volume type VAST",
							Computed:    true,
						},
						"share_settings": schema.SingleNestedAttribute{
							Description: "Share settings specific to the file share type",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudFileSharesShareSettingsDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"type_name": schema.StringAttribute{
									Description: "Standard file share type\nAvailable values: \"standard\", \"vast\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("standard", "vast"),
									},
								},
								"allowed_characters": schema.StringAttribute{
									Description: `Available values: "LCD", "NPL".`,
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("LCD", "NPL"),
									},
								},
								"path_length": schema.StringAttribute{
									Description: `Available values: "LCD", "NPL".`,
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("LCD", "NPL"),
									},
								},
								"root_squash": schema.BoolAttribute{
									Description: "Enables or disables root squash for NFS clients.\n- If `true`, root squash is enabled: the root user is mapped to nobody for all file and folder management operations on the export.\n- If `false`, root squash is disabled: the NFS client `root` user retains root privileges.",
									Computed:    true,
								},
							},
						},
						"size": schema.Int64Attribute{
							Description: "File share size in GiB",
							Computed:    true,
							Validators: []validator.Int64{
								int64validator.AtLeast(1),
							},
						},
						"status": schema.StringAttribute{
							Description: "File share status\nAvailable values: \"available\", \"awaiting_transfer\", \"backup_creating\", \"backup_restoring\", \"backup_restoring_error\", \"creating\", \"creating_from_snapshot\", \"deleted\", \"deleting\", \"ensuring\", \"error\", \"error_deleting\", \"extending\", \"extending_error\", \"inactive\", \"manage_error\", \"manage_starting\", \"migrating\", \"migrating_to\", \"replication_change\", \"reverting\", \"reverting_error\", \"shrinking\", \"shrinking_error\", \"shrinking_possible_data_loss_error\", \"unmanage_error\", \"unmanage_starting\", \"unmanaged\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"available",
									"awaiting_transfer",
									"backup_creating",
									"backup_restoring",
									"backup_restoring_error",
									"creating",
									"creating_from_snapshot",
									"deleted",
									"deleting",
									"ensuring",
									"error",
									"error_deleting",
									"extending",
									"extending_error",
									"inactive",
									"manage_error",
									"manage_starting",
									"migrating",
									"migrating_to",
									"replication_change",
									"reverting",
									"reverting_error",
									"shrinking",
									"shrinking_error",
									"shrinking_possible_data_loss_error",
									"unmanage_error",
									"unmanage_starting",
									"unmanaged",
								),
							},
						},
						"subnet_id": schema.StringAttribute{
							Description: "Subnet ID.",
							Computed:    true,
						},
						"subnet_name": schema.StringAttribute{
							Description: "Subnet name.",
							Computed:    true,
						},
						"tags": schema.ListNestedAttribute{
							Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudFileSharesTagsDataSourceModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Description: "Tag key. Maximum 255 characters. Cannot contain spaces, tabs, newlines, empty string or '=' character.",
										Computed:    true,
									},
									"read_only": schema.BoolAttribute{
										Description: "If true, the tag is read-only and cannot be modified by the user",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Tag value. Maximum 255 characters. Cannot contain spaces, tabs, newlines, empty string or '=' character.",
										Computed:    true,
									},
								},
							},
						},
						"type_name": schema.StringAttribute{
							Description: "File share type name\nAvailable values: \"standard\", \"vast\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("standard", "vast"),
							},
						},
					},
				},
			},
		},
	}
}

func (d *CloudFileSharesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudFileSharesDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

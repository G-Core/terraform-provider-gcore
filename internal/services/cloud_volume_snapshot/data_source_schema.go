// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume_snapshot

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudVolumeSnapshotDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the snapshot",
				Computed:    true,
			},
			"snapshot_id": schema.StringAttribute{
				Description: "Unique identifier of the snapshot",
				Optional:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the snapshot was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Snapshot description",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Snapshot name",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "Snapshot size, GiB",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Snapshot status\nAvailable values: \"available\", \"backing-up\", \"creating\", \"deleted\", \"deleting\", \"error\", \"error_deleting\", \"restoring\", \"unmanaging\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"available",
						"backing-up",
						"creating",
						"deleted",
						"deleting",
						"error",
						"error_deleting",
						"restoring",
						"unmanaging",
					),
				},
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the snapshot was last updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"volume_id": schema.StringAttribute{
				Description: "ID of the volume this snapshot was made from",
				Computed:    true,
			},
			"tags": schema.ListNestedAttribute{
				Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudVolumeSnapshotTagsDataSourceModel](ctx),
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
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"instance_id": schema.StringAttribute{
						Description: "Optional. Filter snapshots by the instance whose volumes were snapshotted",
						Optional:    true,
					},
					"lifecycle_policy_id": schema.Int64Attribute{
						Description: "Optional. Filter by lifecycle policy ID",
						Optional:    true,
					},
					"schedule_id": schema.StringAttribute{
						Description: "Optional. Filter by schedule ID",
						Optional:    true,
					},
					"volume_id": schema.StringAttribute{
						Description: "Optional. Filter by volume ID",
						Optional:    true,
					},
				},
			},
		},
	}
}

func (d *CloudVolumeSnapshotDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudVolumeSnapshotDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("snapshot_id"), path.MatchRoot("find_one_by")),
	}
}

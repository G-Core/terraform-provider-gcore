// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume

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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudVolumeDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Volumes are block storage devices that can be attached to instances as boot or data disks, with support for resizing and type changes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Volume ID",
				Computed:    true,
			},
			"volume_id": schema.StringAttribute{
				Description: "Volume ID",
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
			"bootable": schema.BoolAttribute{
				Description: "Indicates whether the volume is bootable.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The date and time when the volume was created.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "The ID of the task that created this volume.",
				Computed:    true,
			},
			"is_root_volume": schema.BoolAttribute{
				Description: "Indicates whether this is a root volume.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the volume.",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "The region where the volume is located.",
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "The size of the volume in gibibytes (GiB).",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The current status of the volume.\nAvailable values: \"attaching\", \"available\", \"awaiting-transfer\", \"backing-up\", \"creating\", \"deleting\", \"detaching\", \"downloading\", \"error\", \"error_backing-up\", \"error_deleting\", \"error_extending\", \"error_restoring\", \"extending\", \"in-use\", \"maintenance\", \"reserved\", \"restoring-backup\", \"retyping\", \"reverting\", \"uploading\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"attaching",
						"available",
						"awaiting-transfer",
						"backing-up",
						"creating",
						"deleting",
						"detaching",
						"downloading",
						"error",
						"error_backing-up",
						"error_deleting",
						"error_extending",
						"error_restoring",
						"extending",
						"in-use",
						"maintenance",
						"reserved",
						"restoring-backup",
						"retyping",
						"reverting",
						"uploading",
					),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "The date and time when the volume was last updated.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"volume_type": schema.StringAttribute{
				Description: "The type of volume storage.",
				Computed:    true,
			},
			"snapshot_ids": schema.ListAttribute{
				Description: "List of snapshot IDs associated with this volume.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"volume_image_metadata": schema.MapAttribute{
				Description: "Image metadata for volumes created from an image.",
				Computed:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"attachments": schema.ListNestedAttribute{
				Description: "List of attachments associated with the volume.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudVolumeAttachmentsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"attachment_id": schema.StringAttribute{
							Description: "The unique identifier of the attachment object.",
							Computed:    true,
						},
						"volume_id": schema.StringAttribute{
							Description: "The unique identifier of the attached volume.",
							Computed:    true,
						},
						"attached_at": schema.StringAttribute{
							Description: "The date and time when the attachment was created.",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"device": schema.StringAttribute{
							Description: "The block device name inside the guest instance.",
							Computed:    true,
						},
						"flavor_id": schema.StringAttribute{
							Description: "The flavor ID of the instance.",
							Computed:    true,
						},
						"instance_name": schema.StringAttribute{
							Description: "The name of the instance if attached and the server name is known.",
							Computed:    true,
						},
						"server_id": schema.StringAttribute{
							Description: "The unique identifier of the instance.",
							Computed:    true,
						},
					},
				},
			},
			"limiter_stats": schema.SingleNestedAttribute{
				Description: "Schema representing the Quality of Service (QoS) parameters for a volume.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudVolumeLimiterStatsDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"iops_base_limit": schema.Int64Attribute{
						Description: "The sustained IOPS (Input/Output Operations Per Second) limit.",
						Computed:    true,
					},
					"iops_burst_limit": schema.Int64Attribute{
						Description: "The burst IOPS limit.",
						Computed:    true,
					},
					"m_bps_base_limit": schema.Int64Attribute{
						Description: "The sustained bandwidth limit in megabytes per second (MBps).",
						Computed:    true,
					},
					"m_bps_burst_limit": schema.Int64Attribute{
						Description: "The burst bandwidth limit in megabytes per second (MBps).",
						Computed:    true,
					},
				},
			},
			"tags": schema.ListNestedAttribute{
				Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudVolumeTagsDataSourceModel](ctx),
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
					"bootable": schema.BoolAttribute{
						Description: "Filter by bootable field",
						Optional:    true,
					},
					"cluster_id": schema.StringAttribute{
						Description: "Filter volumes by k8s cluster ID",
						Optional:    true,
					},
					"has_attachments": schema.BoolAttribute{
						Description: "Filter by the presence of attachments",
						Optional:    true,
					},
					"id_part": schema.StringAttribute{
						Description: "Filter the volume list result by the ID part of the volume",
						Optional:    true,
					},
					"instance_id": schema.StringAttribute{
						Description: "Filter volumes by instance ID",
						Optional:    true,
					},
					"name_part": schema.StringAttribute{
						Description: "Filter volumes by `name_part` inclusion in volume name.Any substring can be used and volumes will be returned with names containing the substring.",
						Optional:    true,
					},
					"tag_key": schema.ListAttribute{
						Description: "Optional. Filter by tag keys. ?`tag_key`=key1&`tag_key`=key2",
						Optional:    true,
						ElementType: types.StringType,
					},
					"tag_key_value": schema.StringAttribute{
						Description: "Optional. Filter by tag key-value pairs.",
						Optional:    true,
					},
				},
			},
		},
	}
}

func (d *CloudVolumeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudVolumeDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("volume_id"), path.MatchRoot("find_one_by")),
	}
}

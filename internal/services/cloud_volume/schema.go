// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*CloudVolumeResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.Int64Attribute{
				Description:   "Project ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"region_id": schema.Int64Attribute{
				Description:   "Region ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"source": schema.StringAttribute{
				Description: "Volume source type\nAvailable values: \"image\", \"snapshot\", \"new-volume\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"image",
						"snapshot",
						"new-volume",
					),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"attachment_tag": schema.StringAttribute{
				Description:   "Block device attachment tag (not exposed in the user tags). Only used in conjunction with `instance_id_to_attach_to`",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"image_id": schema.StringAttribute{
				Description:   "Image ID",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"instance_id_to_attach_to": schema.StringAttribute{
				Description:   "`instance_id` to attach newly-created volume to",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"size": schema.Int64Attribute{
				Description: "Volume size in GiB",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"snapshot_id": schema.StringAttribute{
				Description:   "Snapshot ID",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"type_name": schema.StringAttribute{
				Description: "Volume type. Defaults to `standard`. If not specified for source `snapshot`, volume type will be derived from the snapshot volume.\nAvailable values: \"cold\", \"ssd_hiiops\", \"ssd_local\", \"ssd_lowlatency\", \"standard\", \"ultra\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"cold",
						"ssd_hiiops",
						"ssd_local",
						"ssd_lowlatency",
						"standard",
						"ultra",
					),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"lifecycle_policy_ids": schema.ListAttribute{
				Description:   "List of lifecycle policy IDs (snapshot creation schedules) to associate with the volume",
				Optional:      true,
				ElementType:   types.Int64Type,
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description: "Volume name",
				Required:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Optional:    true,
				ElementType: types.StringType,
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
			"region": schema.StringAttribute{
				Description: "The region where the volume is located.",
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
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
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
			"tasks": schema.ListAttribute{
				Description: "List of task IDs representing asynchronous operations. Use these IDs to monitor operation progress:\n\\* `GET /v1/tasks/{`task_id`}` - Check individual task status and details\nPoll task status until completion (`FINISHED`/`ERROR`) before proceeding with dependent operations.",
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
				CustomType:  customfield.NewNestedObjectListType[CloudVolumeAttachmentsModel](ctx),
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
				CustomType:  customfield.NewNestedObjectType[CloudVolumeLimiterStatsModel](ctx),
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
		},
	}
}

func (r *CloudVolumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudVolumeResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

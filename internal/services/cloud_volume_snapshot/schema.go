// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume_snapshot

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudVolumeSnapshotResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseNonNullStateForUnknown()},
			},
			"project_id": schema.Int64Attribute{
				Description:   "Project ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
			},
			"region_id": schema.Int64Attribute{
				Description:   "Region ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
			},
			"volume_id": schema.StringAttribute{
				Description:   "Volume ID to make snapshot of",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"description": schema.StringAttribute{
				Description:   "Snapshot description",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description: "Snapshot name",
				Required:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
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
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the snapshot was last updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *CloudVolumeSnapshotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudVolumeSnapshotResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

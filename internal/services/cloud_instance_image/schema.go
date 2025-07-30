// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance_image

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudInstanceImageResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "Image ID",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"image_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"region_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"hw_firmware_type": schema.StringAttribute{
				Description: "Specifies the type of firmware with which to boot the guest.\nAvailable values: \"bios\", \"uefi\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("bios", "uefi"),
				},
			},
			"hw_machine_type": schema.StringAttribute{
				Description: "A virtual chipset type.\nAvailable values: \"pc\", \"q35\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("pc", "q35"),
				},
			},
			"is_baremetal": schema.BoolAttribute{
				Description: "Set to true if the image will be used by bare metal servers.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Image display name",
				Optional:    true,
			},
			"os_type": schema.StringAttribute{
				Description: "The operating system installed on the image.\nAvailable values: \"linux\", \"windows\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("linux", "windows"),
				},
			},
			"ssh_key": schema.StringAttribute{
				Description: "Whether the image supports SSH key or not\nAvailable values: \"allow\", \"deny\", \"required\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"allow",
						"deny",
						"required",
					),
				},
			},
			"tags": schema.MapAttribute{
				Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"architecture": schema.StringAttribute{
				Description: "An image architecture type: aarch64, `x86_64`\nAvailable values: \"aarch64\", \"x86_64\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("aarch64", "x86_64"),
				},
				Default: stringdefault.StaticString("x86_64"),
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the image was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Image description",
				Computed:    true,
			},
			"disk_format": schema.StringAttribute{
				Description: "Disk format",
				Computed:    true,
			},
			"display_order": schema.Int64Attribute{
				Computed: true,
			},
			"min_disk": schema.Int64Attribute{
				Description: "Minimal boot volume required",
				Computed:    true,
			},
			"min_ram": schema.Int64Attribute{
				Description: "Minimal VM RAM required",
				Computed:    true,
			},
			"os_distro": schema.StringAttribute{
				Description: "OS Distribution, i.e. Debian, CentOS, Ubuntu, CoreOS etc.",
				Computed:    true,
			},
			"os_version": schema.StringAttribute{
				Description: "OS version, i.e. 19.04 (for Ubuntu) or 9.4 for Debian",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "Image size in bytes",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Image status, i.e. active",
				Computed:    true,
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the image was updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"visibility": schema.StringAttribute{
				Description: "Image visibility. Globally visible images are public",
				Computed:    true,
			},
		},
	}
}

func (r *CloudInstanceImageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudInstanceImageResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*CloudFileShareResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"file_share_id": schema.StringAttribute{
				Description:   "File Share ID",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
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
			"protocol": schema.StringAttribute{
				Description: "File share protocol\nAvailable values: \"NFS\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("NFS"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"size": schema.Int64Attribute{
				Description: "File share size",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"network": schema.SingleNestedAttribute{
				Description: "File share network configuration",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"network_id": schema.StringAttribute{
						Description: "Network ID.",
						Required:    true,
					},
					"subnet_id": schema.StringAttribute{
						Description: "Subnetwork ID. If the subnet is not selected, it will be selected automatically.",
						Optional:    true,
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
			},
			"type_name": schema.StringAttribute{
				Description: "Standard file share type\nAvailable values: \"standard\", \"vast\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("standard", "vast"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"volume_type": schema.StringAttribute{
				Description:        "Deprecated. Use `type_name` instead.\nAvailable values: \"default_share_type\", \"vast_share_type\".",
				Computed:           true,
				Optional:           true,
				DeprecationMessage: "This attribute is deprecated.",
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("default_share_type", "vast_share_type"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"access": schema.ListNestedAttribute{
				Description: "Access Rules",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudFileShareAccessModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_mode": schema.StringAttribute{
							Description: "Access mode\nAvailable values: \"ro\", \"rw\".",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("ro", "rw"),
							},
						},
						"ip_address": schema.StringAttribute{
							Description: "Source IP or network",
							Required:    true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplaceIfConfigured()},
			},
			"share_settings": schema.SingleNestedAttribute{
				Description: "Configuration settings for the share",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudFileShareShareSettingsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"root_squash": schema.BoolAttribute{
						Description: "Enables or disables root squash for NFS clients.\n- If `true` (default), root squash is enabled: the root user is mapped to nobody for all file and folder management operations on the export.\n- If `false`, root squash is disabled: the NFS client `root` user retains root privileges. Use this option if you trust the root user not to perform operations that will corrupt data.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(true),
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplaceIfConfigured()},
			},
			"name": schema.StringAttribute{
				Description: "File share name",
				Required:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Optional:    true,
				ElementType: types.StringType,
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
			"id": schema.StringAttribute{
				Description: "File share ID",
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
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"share_network_name": schema.StringAttribute{
				Description: "Share network name. May be null if the file share was created with volume type VAST",
				Computed:    true,
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
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"tasks": schema.ListAttribute{
				Description: "List of task IDs",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
		},
	}
}

func (r *CloudFileShareResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudFileShareResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

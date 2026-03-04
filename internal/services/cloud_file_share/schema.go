// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
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
)

var _ resource.ResourceWithConfigValidators = (*CloudFileShareResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "File shares provide NFS-based shared storage that can be mounted by virtual machines and Kubernetes clusters for persistent data.",
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
			"protocol": schema.StringAttribute{
				Description: "File share protocol\nAvailable values: \"NFS\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("NFS"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"size": schema.Int64Attribute{
				Description: "File share size in GiB",
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
			"name": schema.StringAttribute{
				Description: "File share name",
				Required:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"share_settings": schema.SingleNestedAttribute{
				Description: "Configuration settings for the share",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudFileShareShareSettingsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"allowed_characters": schema.StringAttribute{
						Description: "Determines which characters are allowed in file names. Choose between:\n- Lowest Common Denominator (LCD), allows only characters allowed by all VAST Cluster-supported protocols\n- Native Protocol Limit (NPL), imposes no limitation beyond that of the client protocol.\nAvailable values: \"LCD\", \"NPL\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("LCD", "NPL"),
						},
					},
					"path_length": schema.StringAttribute{
						Description: "Affects the maximum limit of file path component name length. Choose between:\n- Lowest Common Denominator (LCD), imposes the lowest common denominator file length limit of all VAST Cluster-supported protocols. With this (default) option, the limitation on the length of a single component of the path is 255 characters\n- Native Protocol Limit (NPL), imposes no limitation beyond that of the client protocol.\nAvailable values: \"LCD\", \"NPL\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("LCD", "NPL"),
						},
					},
					"root_squash": schema.BoolAttribute{
						Description: "Enables or disables root squash for NFS clients.\n- If `true` (default), root squash is enabled: the root user is mapped to nobody for all file and folder management operations on the export.\n- If `false`, root squash is disabled: the NFS client `root` user retains root privileges. Use this option if you trust the root user not to perform operations that will corrupt data.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(true),
					},
				},
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
				Description: "List of task IDs representing asynchronous operations. Use these IDs to monitor operation progress:\n- `GET /v1/tasks/{task_id}` - Check individual task status and details\nPoll task status until completion (`FINISHED`/`ERROR`) before proceeding with dependent operations.",
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

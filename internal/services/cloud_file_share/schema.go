// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudFileShareResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "File shares provide NFS-based shared storage that can be mounted by virtual machines and Kubernetes clusters for persistent data.",
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
			"protocol": schema.StringAttribute{
				Description: "File share protocol\nAvailable values: \"NFS\", \"LUSTRE\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("NFS", "LUSTRE"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"size": schema.Int64Attribute{
				Description: "File share size in GiB",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"type_name": schema.StringAttribute{
				Description: "Standard file share type\nAvailable values: \"standard\", \"ddn\", \"vast\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"standard",
						"ddn",
						"vast",
					),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"access_rule_ids": schema.ListAttribute{
				Description: "List of access rules IDs associated with the file share",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"network": schema.SingleNestedAttribute{
				Description: "File share network configuration",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudFileShareNetworkModel](ctx),
				Attributes: map[string]schema.Attribute{
					"network_id": schema.StringAttribute{
						Description:   "Network ID.",
						Computed:      true,
						Optional:      true,
						PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
					},
					"subnet_id": schema.StringAttribute{
						Description: "Subnetwork ID. If the subnet is not selected, it will be selected automatically.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
							stringplanmodifier.RequiresReplaceIfConfigured(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description: "File share name",
				Required:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"share_settings": schema.SingleNestedAttribute{
				Description: "Configuration settings for the share",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudFileShareShareSettingsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"gid": schema.Int64Attribute{
						Description: "When set created file share will be owned by user with given GID",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(0, 2147483647),
						},
					},
					"uid": schema.Int64Attribute{
						Description: "When set created file share will be owned by group with given UID",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(0, 2147483647),
						},
					},
					"allowed_characters": schema.StringAttribute{
						Description: "Determines which characters are allowed in file names. Choose between:\n  - Lowest Common Denominator (LCD), allows only characters allowed by all VAST Cluster-supported protocols\n  - Native Protocol Limit (NPL), imposes no limitation beyond that of the client protocol.\nAvailable values: \"LCD\", \"NPL\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("LCD", "NPL"),
						},
					},
					"path_length": schema.StringAttribute{
						Description: "Affects the maximum limit of file path component name length. Choose between:\n  - Lowest Common Denominator (LCD), imposes the lowest common denominator file length limit of all VAST Cluster-supported protocols. With this option, the limitation on the length of a single component of the path is 255 characters\n  - Native Protocol Limit (NPL), imposes no limitation beyond that of the client protocol.\nAvailable values: \"LCD\", \"NPL\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("LCD", "NPL"),
						},
					},
					"root_squash": schema.BoolAttribute{
						Description: "Enables or disables root squash for NFS clients.\n  - If `true`, root squash is enabled: the root user is mapped to nobody for all file and folder management operations on the export.\n  - If `false`, root squash is disabled: the NFS client `root` user retains root privileges. Use this option if you trust the root user not to perform operations that will corrupt data.",
						Computed:    true,
						Optional:    true,
					},
					"trash_access": schema.BoolAttribute{
						Description: "Enables or disables access to the hidden `.vast_trash` folder of the share.\n  - If `false`, the folder is not accessible.\n  - If `true`, clients within the share's access range can move files and folders to `.vast_trash` for asynchronous deletion, which is much faster than deleting large directory trees in place.\n\nRequires `root_squash` to be disabled, because VAST ignores trash folder access for hosts that are under root squash.",
						Computed:    true,
						Optional:    true,
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
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
				Description:        "Deprecated. Use `network.network_id` instead. Network ID.",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated.",
			},
			"network_name": schema.StringAttribute{
				Description:        "Deprecated. Use `network.network_name` instead. Network name.",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated.",
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
				Description:        "Deprecated. Use `network.subnet_id` instead. Subnet ID.",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated.",
			},
			"subnet_name": schema.StringAttribute{
				Description:        "Deprecated. Use `network.subnet_name` instead. Subnet name.",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated.",
			},
		},
	}
}

func (r *CloudFileShareResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudFileShareResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		&typeNameShareSettingsValidator{},
		&typeNameNetworkValidator{},
	}
}

type typeNameShareSettingsValidator struct{}

func (v *typeNameShareSettingsValidator) Description(_ context.Context) string {
	return "share_settings is only allowed when type_name is 'vast'"
}

func (v *typeNameShareSettingsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *typeNameShareSettingsValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config CloudFileShareModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If share_settings is provided but type_name is not "vast", add error
	if !config.ShareSettings.IsNull() && !config.ShareSettings.IsUnknown() {
		typeName := config.TypeName.ValueString()
		if typeName != "" && typeName != "vast" {
			resp.Diagnostics.AddAttributeError(
				path.Root("share_settings"),
				"Invalid Attribute Combination",
				"share_settings can only be provided when type_name is 'vast'",
			)
		}
	}
}

type typeNameNetworkValidator struct{}

func (v *typeNameNetworkValidator) Description(_ context.Context) string {
	return "network is only allowed when type_name is 'standard'"
}

func (v *typeNameNetworkValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *typeNameNetworkValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config CloudFileShareModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// An unknown value (e.g. from a variable) cannot be checked until apply
	if config.Network.IsNull() || config.Network.IsUnknown() {
		return
	}

	// network_id is computed so that a share without a network block keeps the
	// one the API returns, but a configured block still needs it
	var networkID types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("network").AtName("network_id"), &networkID)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if networkID.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("network").AtName("network_id"),
			"Missing Required Attribute",
			"network_id is required when network is provided",
		)
	}

	if config.TypeName.IsUnknown() {
		return
	}

	// If network is provided but type_name is not "standard", add error
	typeName := config.TypeName.ValueString()
	if typeName != "" && typeName != "standard" {
		resp.Diagnostics.AddAttributeError(
			path.Root("network"),
			"Invalid Attribute Combination",
			"network can only be provided when type_name is 'standard'",
		)
	}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
	t "github.com/stainless-sdks/gcore-terraform/internal/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudGPUVirtualClusterResource)(nil)

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
			"flavor": schema.StringAttribute{
				Description:   "Cluster flavor ID",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"servers_count": schema.Int64Attribute{
				Description: "Number of servers in the cluster",
				Required:    true,
			},
			"servers_settings": schema.SingleNestedAttribute{
				Description: "Configuration settings for the servers in the cluster",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"interfaces": schema.ListNestedAttribute{
						Description: "Subnet IPs and floating IPs",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: `Available values: "external", "subnet", "any_subnet".`,
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive(
											"external",
											"subnet",
											"any_subnet",
										),
									},
								},
								"ip_family": schema.StringAttribute{
									Description: "Which subnets should be selected: IPv4, IPv6, or use dual stack.\nAvailable values: \"dual\", \"ipv4\", \"ipv6\".",
									Computed:    true,
									Optional:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive(
											"dual",
											"ipv4",
											"ipv6",
										),
									},
									PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								},
								"name": schema.StringAttribute{
									Description:   "Interface name",
									Computed:      true,
									Optional:      true,
									PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								},
								"network_id": schema.StringAttribute{
									Description: "Network ID the subnet belongs to. Port will be plugged in this network",
									Optional:    true,
								},
								"subnet_id": schema.StringAttribute{
									Description: "Port is assigned an IP address from this subnet",
									Optional:    true,
								},
								"floating_ip": schema.SingleNestedAttribute{
									Description: "Floating IP config for this subnet attachment",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"source": schema.StringAttribute{
											Description: `Available values: "new".`,
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOfCaseInsensitive("new"),
											},
										},
									},
								},
							},
						},
					},
					"volumes": schema.ListNestedAttribute{
						Description: "List of volumes",
						Required:    true,
						Validators: []validator.List{
							volumeImageSourceValidator{},
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"boot_index": schema.Int64Attribute{
									Description: "Boot index of the volume",
									Required:    true,
									Validators: []validator.Int64{
										int64validator.AtLeast(0),
									},
								},
								"name": schema.StringAttribute{
									Description: "Volume name",
									Required:    true,
								},
								"size": schema.Int64Attribute{
									Description: "Volume size in GiB",
									Required:    true,
								},
								"source": schema.StringAttribute{
									Description: `Available values: "new", "image".`,
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("new", "image"),
									},
								},
								"type": schema.StringAttribute{
									Description: "Volume type\nAvailable values: \"cold\", \"ssd_hiiops\", \"ssd_local\", \"ssd_lowlatency\", \"standard\", \"ultra\".",
									Required:    true,
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
								},
								"delete_on_termination": schema.BoolAttribute{
									Description: "Flag indicating whether the volume is deleted on instance termination",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
								"tags": schema.MapAttribute{
									Description: "Tags associated with the volume",
									Optional:    true,
									ElementType: types.StringType,
								},
								"image_id": schema.StringAttribute{
									Description: "Image ID for the volume (required if source is 'image')",
									Optional:    true,
								},
							},
						},
					},
					"credentials": schema.SingleNestedAttribute{
						Description: "Optional server access credentials",
						Optional:    true,
						Validators: []validator.Object{
							credentialsValidator{},
						},
						Attributes: map[string]schema.Attribute{
							"password_wo": schema.StringAttribute{
								Description: "Used to set the password for the specified 'username' on Linux instances. If 'username' is not provided, the password is applied to the default user of the image. Mutually exclusive with '`user_data`' - only one can be specified.",
								Optional:    true,
								WriteOnly:   true,
							},
							"ssh_key_name": schema.StringAttribute{
								Description: "Specifies the name of the SSH keypair, created via the\n[/v1/`ssh_keys` endpoint](/docs/api-reference/cloud/ssh-keys/add-or-generate-ssh-key).",
								Optional:    true,
							},
							"username": schema.StringAttribute{
								Description: "The 'username' and 'password' fields create a new user on the system",
								Optional:    true,
							},
							"password_wo_version": schema.Int64Attribute{
								Description: "Version of the password write-only field. Increment this value to trigger an update when changing the password.",
								Optional:    true,
							},
						},
					},
					"file_shares": schema.ListNestedAttribute{
						Description:   "List of file shares to be mounted across the cluster.",
						Computed:      true,
						Optional:      true,
						CustomType:    customfield.NewNestedObjectListType[CloudGPUVirtualClusterServersSettingsFileSharesModel](ctx),
						PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "Unique identifier of the file share in UUID format.",
									Required:    true,
								},
								"mount_path": schema.StringAttribute{
									Description: "Absolute mount path inside the system where the file share will be mounted.",
									Required:    true,
								},
							},
						},
					},
					"security_groups": schema.ListNestedAttribute{
						Description:   "List of security groups UUIDs",
						Computed:      true,
						Optional:      true,
						CustomType:    customfield.NewNestedObjectListType[CloudGPUVirtualClusterServersSettingsSecurityGroupsModel](ctx),
						PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "Resource ID",
									Required:    true,
								},
							},
						},
					},
					"user_data": schema.StringAttribute{
						Description:   "Optional custom user data (Base64-encoded)",
						Computed:      true,
						Optional:      true,
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
				},
				PlanModifiers: []planmodifier.Object{serversSettingsRequiresReplaceModifier{}},
			},
			"tags": schema.StringAttribute{
				Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				Optional:    true,
				CustomType:  jsontypes.NormalizedType{},
			},
			"name": schema.StringAttribute{
				Description: "Cluster name",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Cluster creation date time",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"status": schema.StringAttribute{
				Description: "Cluster status\nAvailable values: \"active\", \"creating\", \"degraded\", \"deleting\", \"error\", \"new\", \"rebooting\", \"rebuilding\", \"resizing\", \"shutoff\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"active",
						"creating",
						"degraded",
						"deleting",
						"error",
						"new",
						"rebooting",
						"rebuilding",
						"resizing",
						"shutoff",
					),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "Cluster update date time",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"servers_ids": schema.ListAttribute{
				Description:   "List of cluster nodes",
				Computed:      true,
				CustomType:    customfield.NewListType[types.String](ctx),
				ElementType:   types.StringType,
				PlanModifiers: []planmodifier.List{serversIDsPlanModifier{}},
			},
		},
	}
}

func (r *CloudGPUVirtualClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudGPUVirtualClusterResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

type volumeImageSourceValidator struct{}

func (v volumeImageSourceValidator) Description(_ context.Context) string {
	return "validates that image_id is provided when source is 'image'"
}

func (v volumeImageSourceValidator) MarkdownDescription(_ context.Context) string {
	return "validates that image_id is provided when source is 'image'"
}

func (v volumeImageSourceValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var volumes []CloudGPUVirtualClusterServersSettingsVolumesModel

	resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &volumes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for i, vol := range volumes {
		if !vol.Source.IsNull() && !vol.Source.IsUnknown() &&
			vol.Source.ValueString() == "image" &&
			(vol.ImageID.IsNull() || vol.ImageID.ValueString() == "") {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(i).AtName("image_id"),
				"Missing required attribute",
				"image_id must be specified when source is 'image'",
			)
		}
	}
}

type credentialsValidator struct{}

func (v credentialsValidator) Description(_ context.Context) string {
	return "validates that either ssh_key_name is provided, or both username and password_wo (with password_wo_version) are provided"
}

func (v credentialsValidator) MarkdownDescription(_ context.Context) string {
	return "validates that either ssh_key_name is provided, or both username and password_wo (with password_wo_version) are provided"
}

func (v credentialsValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()

	sshKeyName := attrs["ssh_key_name"]
	username := attrs["username"]
	passwordWo := attrs["password_wo"]
	passwordWoVersion := attrs["password_wo_version"]

	hasSSHKey := !sshKeyName.IsNull() && !sshKeyName.IsUnknown()
	hasUsername := !username.IsNull() && !username.IsUnknown()
	hasPasswordWo := !passwordWo.IsNull() && !passwordWo.IsUnknown()
	hasPasswordWoVersion := !passwordWoVersion.IsNull() && !passwordWoVersion.IsUnknown()

	// Valid scenarios (mutually exclusive):
	// 1. Only ssh_key_name is provided (no username or password_wo or password_wo_version)
	// 2. Only username, password_wo, and password_wo_version are provided together (no ssh_key_name)

	if hasSSHKey && (hasUsername || hasPasswordWo || hasPasswordWoVersion) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Conflicting credentials configuration",
			"Cannot specify 'ssh_key_name' together with 'username', 'password_wo', or 'password_wo_version'. Only one authentication method can be used.",
		)
	} else if !hasSSHKey && !(hasUsername && hasPasswordWo) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid credentials configuration",
			"Either 'ssh_key_name' must be provided, or both 'username' and 'password_wo' (with 'password_wo_version') must be provided together.",
		)
	} else if hasPasswordWo && !hasPasswordWoVersion {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing password_wo_version",
			"When using 'password_wo', you must also provide 'password_wo_version'. This field is used to track password changes since write-only fields are not stored in state.",
		)
	}
}

// serversIDsPlanModifier preserves the state value for servers_ids unless:
// - The resource is being created (no prior state)
// - The resource is being replaced (id becomes unknown)
// - The servers_count attribute is changing (which affects servers_ids)
type serversIDsPlanModifier struct{}

func (m serversIDsPlanModifier) Description(_ context.Context) string {
	return "Preserves state value unless resource is replaced or servers_count changes"
}

func (m serversIDsPlanModifier) MarkdownDescription(_ context.Context) string {
	return "Preserves state value unless resource is replaced or servers_count changes"
}

func (m serversIDsPlanModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// If there's no state (new resource), nothing to preserve
	if req.State.Raw.IsNull() {
		return
	}

	// If the planned value is already known, don't override it
	if !resp.PlanValue.IsUnknown() {
		return
	}

	// Check if the resource is being replaced by checking if id is becoming unknown
	var planID types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("id"), &planID)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if planID.IsUnknown() {
		// Resource is being replaced, don't preserve state
		return
	}

	// Check if servers_count is changing
	var stateServersCount, planServersCount types.Int64
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("servers_count"), &stateServersCount)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("servers_count"), &planServersCount)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !stateServersCount.Equal(planServersCount) {
		// servers_count is changing, servers_ids will change too
		return
	}

	// Safe to preserve state value
	resp.PlanValue = req.StateValue
}

// serversSettingsRequiresReplaceModifier triggers replacement only when user-specified fields change,
// ignoring computed fields that may show as unknown during planning.
type serversSettingsRequiresReplaceModifier struct{}

func (m serversSettingsRequiresReplaceModifier) Description(_ context.Context) string {
	return "Requires replacement only when user-specified server settings change"
}

func (m serversSettingsRequiresReplaceModifier) MarkdownDescription(_ context.Context) string {
	return "Requires replacement only when user-specified server settings change"
}

func (m serversSettingsRequiresReplaceModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// If there's no state (new resource), no replacement needed
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	// If config is null but state exists, resource is being removed
	if req.ConfigValue.IsNull() {
		return
	}

	// Compare config value with state value to determine if user-specified fields changed.
	// We compare config (what user wrote) vs state (what was applied), not plan vs state,
	// because plan may have computed fields marked as unknown.
	// Use deep comparison that only checks fields explicitly set in config.
	if configValuesChangedFromState(req.ConfigValue, req.StateValue) {
		resp.RequiresReplace = true
	}
}

// configValuesChangedFromState performs a deep comparison between config and state values,
// only checking fields that are explicitly set in config (not null/unknown).
// This allows computed fields in state to differ without triggering replacement.
func configValuesChangedFromState(configVal, stateVal attr.Value) bool {
	// If config value is null or unknown, consider it unchanged (computed field)
	if configVal.IsNull() || configVal.IsUnknown() {
		return false
	}

	// If state is null but config is not, something changed
	if stateVal.IsNull() {
		return true
	}

	// Handle objects/maps - compare only attributes present in config
	if ok, configAttrs := t.ChildAttributes(configVal); ok {
		if ok, stateAttrs := t.ChildAttributes(stateVal); ok {
			for field, configFieldVal := range configAttrs {
				// Skip fields not specified in config
				if configFieldVal.IsNull() || configFieldVal.IsUnknown() {
					continue
				}

				stateFieldVal, exists := stateAttrs[field]
				if !exists {
					return true
				}

				if configValuesChangedFromState(configFieldVal, stateFieldVal) {
					return true
				}
			}
			return false
		}
		return true
	}

	// Handle lists/tuples/sets - compare elements by index
	if ok, configElems := t.ChildItems(configVal); ok {
		if ok, stateElems := t.ChildItems(stateVal); ok {
			if len(configElems) != len(stateElems) {
				return true
			}
			for i, configElem := range configElems {
				if configValuesChangedFromState(configElem, stateElems[i]) {
					return true
				}
			}
			return false
		}
		return true
	}

	// For primitive types, use direct comparison
	return !configVal.Equal(stateVal)
}

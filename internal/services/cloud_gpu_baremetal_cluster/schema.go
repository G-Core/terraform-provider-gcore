// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudGPUBaremetalClusterResource)(nil)

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
			"image_id": schema.StringAttribute{
				Description: "System image ID",
				Required:    true,
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
						CustomType:    customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersSettingsFileSharesModel](ctx),
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
						CustomType:    customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersSettingsSecurityGroupsModel](ctx),
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
				PlanModifiers: []planmodifier.Object{planmodifiers.RequiresReplaceOnConfigChange("credentials", "user_data")},
			},
			"name": schema.StringAttribute{
				Description: "Cluster name",
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
				Description: "Cluster creation date time",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"managed_by": schema.StringAttribute{
				Description: "User type managing the resource\nAvailable values: \"k8s\", \"user\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("k8s", "user"),
				},
			},
			"status": schema.StringAttribute{
				Description: "Cluster status\nAvailable values: \"active\", \"creating\", \"degraded\", \"deleting\", \"error\", \"rebooting\", \"rebuilding\", \"resizing\", \"shutoff\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"active",
						"creating",
						"degraded",
						"deleting",
						"error",
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
				PlanModifiers: []planmodifier.List{planmodifiers.UseStateUnlessCountChanges("servers_count")},
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

func (r *CloudGPUBaremetalClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudGPUBaremetalClusterResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
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

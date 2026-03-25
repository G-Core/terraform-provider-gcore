// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/G-Core/terraform-provider-gcore/internal/customvalidator"
	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudGPUVirtualClusterResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "GPU virtual clusters provide managed virtual GPU servers with auto-scaling for parallel computation workloads.",
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
							customvalidator.CredentialsValidator{},
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
				PlanModifiers: []planmodifier.Object{planmodifiers.RequiresReplaceOnConfigChange()},
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
			!vol.ImageID.IsUnknown() && (vol.ImageID.IsNull() || vol.ImageID.ValueString() == "") {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(i).AtName("image_id"),
				"Missing required attribute",
				"image_id must be specified when source is 'image'",
			)
		}
	}
}

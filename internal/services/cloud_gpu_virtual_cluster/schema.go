// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
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
				Description:   "Number of servers in the cluster",
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
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
									Default: stringdefault.StaticString("ipv4"),
								},
								"name": schema.StringAttribute{
									Description: "Interface name",
									Optional:    true,
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
									Description: "Image ID for the volume",
									Optional:    true,
								},
							},
						},
					},
					"credentials": schema.SingleNestedAttribute{
						Description: "Optional server access credentials",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"password": schema.StringAttribute{
								Description: "Used to set the password for the specified 'username' on Linux instances. If 'username' is not provided, the password is applied to the default user of the image. Mutually exclusive with '`user_data`' - only one can be specified.",
								Optional:    true,
							},
							"ssh_key_name": schema.StringAttribute{
								Description: "Specifies the name of the SSH keypair, created via the\n[/v1/`ssh_keys` endpoint](/docs/api-reference/cloud/ssh-keys/add-or-generate-ssh-key).",
								Optional:    true,
							},
							"username": schema.StringAttribute{
								Description: "The 'username' and 'password' fields create a new user on the system",
								Optional:    true,
							},
						},
					},
					"file_shares": schema.ListNestedAttribute{
						Description: "List of file shares to be mounted across the cluster.",
						Optional:    true,
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
						Description: "List of security groups UUIDs",
						Optional:    true,
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
						Description: "Optional custom user data (Base64-encoded)",
						Optional:    true,
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
			},
			"tags": schema.MapAttribute{
				Description:   "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Optional:      true,
				ElementType:   types.StringType,
				PlanModifiers: []planmodifier.Map{mapplanmodifier.RequiresReplace()},
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
				Description: "Cluster status\nAvailable values: \"active\", \"deleting\", \"error\", \"new\", \"resizing\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"active",
						"deleting",
						"error",
						"new",
						"resizing",
					),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "Cluster update date time",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"servers_ids": schema.ListAttribute{
				Description: "List of cluster nodes",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"tasks": schema.ListAttribute{
				Description: "List of task IDs representing asynchronous operations. Use these IDs to monitor operation progress:\n* `GET /v1/tasks/{task_id}` - Check individual task status and details\nPoll task status until completion (`FINISHED`/`ERROR`) before proceeding with dependent operations.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
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

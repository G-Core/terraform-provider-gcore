// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*CloudBaremetalServerResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
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
				Description:   "The flavor of the instance.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"interfaces": schema.ListNestedAttribute{
				Description: "A list of network interfaces for the server. You can create one or more interfaces - private, public, or both.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "A public IP address will be assigned to the instance.\nAvailable values: \"external\", \"subnet\", \"any_subnet\", \"reserved_fixed_ip\".",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"external",
									"subnet",
									"any_subnet",
									"reserved_fixed_ip",
								),
							},
						},
						"interface_name": schema.StringAttribute{
							Description: "Interface name. Defaults to `null` and is returned as `null` in the API response if not set.",
							Optional:    true,
						},
						"ip_family": schema.StringAttribute{
							Description: "Specify `ipv4`, `ipv6`, or `dual` to enable both.\nAvailable values: \"dual\", \"ipv4\", \"ipv6\".",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"dual",
									"ipv4",
									"ipv6",
								),
							},
						},
						"port_group": schema.Int64Attribute{
							Description: "Specifies the trunk group to which this interface belongs. Applicable only for bare metal servers. Each unique port group is mapped to a separate trunk port. Use this to control how interfaces are grouped across trunks.",
							Computed:    true,
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.Between(0, 1),
							},
							Default: int64default.StaticInt64(0),
						},
						"network_id": schema.StringAttribute{
							Description: "The network where the instance will be connected.",
							Optional:    true,
						},
						"subnet_id": schema.StringAttribute{
							Description: "The instance will get an IP address from this subnet.",
							Optional:    true,
						},
						"floating_ip": schema.SingleNestedAttribute{
							Description: "Allows the instance to have a public IP that can be reached from the internet.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"source": schema.StringAttribute{
									Description: "A new floating IP will be created and attached to the instance. A floating IP is a public IP that makes the instance accessible from the internet, even if it only has a private IP. It works like SNAT, allowing outgoing and incoming traffic.\nAvailable values: \"new\", \"existing\".",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("new", "existing"),
									},
								},
								"existing_floating_id": schema.StringAttribute{
									Description: "An existing available floating IP id must be specified if the source is set to `existing`",
									Optional:    true,
								},
							},
						},
						"ip_address": schema.StringAttribute{
							Description: "You can specify a specific IP address from your subnet.",
							Optional:    true,
						},
						"port_id": schema.StringAttribute{
							Description: "Network ID the subnet belongs to. Port will be plugged in this network.",
							Optional:    true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"apptemplate_id": schema.StringAttribute{
				Description:   "Apptemplate ID. Either `image_id` or `apptemplate_id` is required.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"image_id": schema.StringAttribute{
				Description:   "Image ID. Either `image_id` or `apptemplate_id` is required.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "Server name.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name_template": schema.StringAttribute{
				Description:   "If you want server names to be automatically generated based on IP addresses, you can provide a name template instead of specifying the name manually. The template should include a placeholder that will be replaced during provisioning. Supported placeholders are: `{`ip_octets`}` (last 3 octets of the IP), `{`two_ip_octets`}`, and `{`one_ip_octet`}`.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"password": schema.StringAttribute{
				Description:   "For Linux instances, 'username' and 'password' are used to create a new user. When only 'password' is provided, it is set as the password for the default user of the image. For Windows instances, 'username' cannot be specified. Use the 'password' field to set the password for the 'Admin' user on Windows. Use the '`user_data`' field to provide a script to create new users on Windows. The password of the Admin user cannot be updated via '`user_data`'.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"ssh_key_name": schema.StringAttribute{
				Description:   "Specifies the name of the SSH keypair, created via the\n[/v1/`ssh_keys` endpoint](/docs/api-reference/cloud/ssh-keys/add-or-generate-ssh-key).",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"user_data": schema.StringAttribute{
				Description:   "String in base64 format. For Linux instances, '`user_data`' is ignored when 'password' field is provided. For Windows instances, Admin user password is set by 'password' field and cannot be updated via '`user_data`'. Examples of the `user_data`: https://cloudinit.readthedocs.io/en/latest/topics/examples.html",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"username": schema.StringAttribute{
				Description:   "For Linux instances, 'username' and 'password' are used to create a new user. For Windows instances, 'username' cannot be specified. Use 'password' field to set the password for the 'Admin' user on Windows.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"tags": schema.MapAttribute{
				Description:   "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Optional:      true,
				ElementType:   types.StringType,
				PlanModifiers: []planmodifier.Map{mapplanmodifier.RequiresReplace()},
			},
			"ddos_profile": schema.SingleNestedAttribute{
				Description: "Enable advanced DDoS protection for the server",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"profile_template": schema.Int64Attribute{
						Description: "Unique identifier of the DDoS protection template to use for this profile",
						Required:    true,
					},
					"fields": schema.ListNestedAttribute{
						Description: "List of field configurations that customize the protection parameters for this profile",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"base_field": schema.Int64Attribute{
									Description: "Unique identifier of the DDoS protection field being configured",
									Optional:    true,
								},
								"field_name": schema.StringAttribute{
									Description: "Human-readable name of the DDoS protection field being configured",
									Optional:    true,
								},
								"field_value": schema.ListAttribute{
									Description: "Complex value. Only one of 'value' or '`field_value`' must be specified.",
									Optional:    true,
									ElementType: jsontypes.NormalizedType{},
								},
								"value": schema.StringAttribute{
									Description:        "Basic type value. Only one of 'value' or '`field_value`' must be specified.",
									Optional:           true,
									DeprecationMessage: "This attribute is deprecated.",
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
			},
			"app_config": schema.StringAttribute{
				Description: "Parameters for the application template if creating the instance from an `apptemplate`.",
				Optional:    true,
				CustomType:  jsontypes.NormalizedType{},
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

func (r *CloudBaremetalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudBaremetalServerResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

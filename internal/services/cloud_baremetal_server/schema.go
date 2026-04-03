// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_server

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
)

var _ resource.ResourceWithConfigValidators = (*CloudBaremetalServerResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Bare metal servers are dedicated physical machines with direct hardware access, supporting provisioning, rebuilding, and network configuration within a cloud region.",
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
								int64validator.Between(0, 3),
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
			"name_template": schema.StringAttribute{
				Description:   "If you want server names to be automatically generated based on IP addresses, you can provide a name template instead of specifying the name manually. The template should include a placeholder that will be replaced during provisioning. Supported placeholders are: `{ip_octets}` (last 3 octets of the IP), `{two_ip_octets}`, and `{one_ip_octet}`.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"password": schema.StringAttribute{
				Description:   "For Linux instances, 'username' and 'password' are used to create a new user. When only 'password' is provided, it is set as the password for the default user of the image. For Windows instances, 'username' cannot be specified. Use the 'password' field to set the password for the 'Admin' user on Windows. Use the 'user_data' field to provide a script to create new users on Windows. The password of the Admin user cannot be updated via 'user_data'.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"ssh_key_name": schema.StringAttribute{
				Description:   "Specifies the name of the SSH keypair, created via the\n[/v1/`ssh_keys` endpoint](/docs/api-reference/cloud/ssh-keys/add-or-generate-ssh-key).",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"user_data": schema.StringAttribute{
				Description:   "String in base64 format. For Linux instances, 'user_data' is ignored when 'password' field is provided. For Windows instances, Admin user password is set by 'password' field and cannot be updated via 'user_data'. Examples of the `user_data`: https://cloudinit.readthedocs.io/en/latest/topics/examples.html",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"username": schema.StringAttribute{
				Description:   "For Linux instances, 'username' and 'password' are used to create a new user. For Windows instances, 'username' cannot be specified. Use 'password' field to set the password for the 'Admin' user on Windows.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"app_config": schema.MapAttribute{
				Description:   "Parameters for the application template if creating the instance from an `apptemplate`.",
				Optional:      true,
				ElementType:   jsontypes.NormalizedType{},
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
								"field_value": schema.StringAttribute{
									Optional:   true,
									CustomType: jsontypes.NormalizedType{},
								},
								"value": schema.StringAttribute{
									Description:        "Basic type value. Only one of 'value' or 'field_value' must be specified.",
									Optional:           true,
									DeprecationMessage: "This attribute is deprecated.",
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description: "Server name.",
				Computed:    true,
				Optional:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when bare metal server was created",
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
			"status": schema.StringAttribute{
				Description: "Bare metal server status\nAvailable values: \"ACTIVE\", \"BUILD\", \"DELETED\", \"ERROR\", \"HARD_REBOOT\", \"MIGRATING\", \"PASSWORD\", \"PAUSED\", \"REBOOT\", \"REBUILD\", \"RESCUE\", \"RESIZE\", \"REVERT_RESIZE\", \"SHELVED\", \"SHELVED_OFFLOADED\", \"SHUTOFF\", \"SOFT_DELETED\", \"SUSPENDED\", \"UNKNOWN\", \"VERIFY_RESIZE\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"ACTIVE",
						"BUILD",
						"DELETED",
						"ERROR",
						"HARD_REBOOT",
						"MIGRATING",
						"PASSWORD",
						"PAUSED",
						"REBOOT",
						"REBUILD",
						"RESCUE",
						"RESIZE",
						"REVERT_RESIZE",
						"SHELVED",
						"SHELVED_OFFLOADED",
						"SHUTOFF",
						"SOFT_DELETED",
						"SUSPENDED",
						"UNKNOWN",
						"VERIFY_RESIZE",
					),
				},
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"task_state": schema.StringAttribute{
				Description: "Task state",
				Computed:    true,
			},
			"vm_state": schema.StringAttribute{
				Description: "Bare metal server state\nAvailable values: \"active\", \"building\", \"deleted\", \"error\", \"paused\", \"rescued\", \"resized\", \"shelved\", \"shelved_offloaded\", \"soft-deleted\", \"stopped\", \"suspended\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"active",
						"building",
						"deleted",
						"error",
						"paused",
						"rescued",
						"resized",
						"shelved",
						"shelved_offloaded",
						"soft-deleted",
						"stopped",
						"suspended",
					),
				},
			},
			"addresses": schema.MapAttribute{
				Description: "Map of `network_name` to list of addresses in that network",
				Computed:    true,
				CustomType:  customfield.NewMapType[customfield.NestedObjectList[CloudBaremetalServerAddressesModel]](ctx),
				ElementType: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{"addr": schema.StringAttribute{
							Description: "Address",
							Required:    true,
						}.GetType(), "type": schema.StringAttribute{
							Description: "Type of the address\nAvailable values: \"floating\", \"fixed\".",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("floating", "fixed"),
							},
						}.GetType(), "interface_name": schema.StringAttribute{
							Description: "Interface name. This field will be `null` if `with_interfaces_name=true` is not set in the request when listing servers. It will also be `null` if the `interface_name` was not specified during server creation or when attaching the interface.",
							Optional:    true,
						}.GetType(), "subnet_id": schema.StringAttribute{
							Description: "The unique identifier of the subnet associated with this address.",
							Optional:    true,
						}.GetType(), "subnet_name": schema.StringAttribute{
							Description: "The name of the subnet associated with this address.",
							Optional:    true,
						}.GetType()},
					},
				},
			},
			"tasks": schema.ListAttribute{
				Description: "List of task IDs representing asynchronous operations. Use these IDs to monitor operation progress:\n- `GET /v1/tasks/{task_id}` - Check individual task status and details\nPoll task status until completion (`FINISHED`/`ERROR`) before proceeding with dependent operations.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"blackhole_ports": schema.ListNestedAttribute{
				Description: "IP addresses of the instances that are blackholed by DDoS mitigation system",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudBaremetalServerBlackholePortsModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"alarm_end": schema.StringAttribute{
							Description: "A date-time string giving the time that the alarm ended. If not yet ended, time will be given as 0001-01-01T00:00:00Z",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"alarm_start": schema.StringAttribute{
							Description: "A date-time string giving the time that the alarm started",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"alarm_state": schema.StringAttribute{
							Description: "Current state of alarm\nAvailable values: \"ACK_REQ\", \"ALARM\", \"ARCHIVED\", \"CLEAR\", \"CLEARING\", \"CLEARING_FAIL\", \"END_GRACE\", \"END_WAIT\", \"MANUAL_CLEAR\", \"MANUAL_CLEARING\", \"MANUAL_CLEARING_FAIL\", \"MANUAL_MITIGATING\", \"MANUAL_STARTING\", \"MANUAL_STARTING_FAIL\", \"MITIGATING\", \"STARTING\", \"STARTING_FAIL\", \"START_WAIT\", \"ack_req\", \"alarm\", \"archived\", \"clear\", \"clearing\", \"clearing_fail\", \"end_grace\", \"end_wait\", \"manual_clear\", \"manual_clearing\", \"manual_clearing_fail\", \"manual_mitigating\", \"manual_starting\", \"manual_starting_fail\", \"mitigating\", \"start_wait\", \"starting\", \"starting_fail\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"ACK_REQ",
									"ALARM",
									"ARCHIVED",
									"CLEAR",
									"CLEARING",
									"CLEARING_FAIL",
									"END_GRACE",
									"END_WAIT",
									"MANUAL_CLEAR",
									"MANUAL_CLEARING",
									"MANUAL_CLEARING_FAIL",
									"MANUAL_MITIGATING",
									"MANUAL_STARTING",
									"MANUAL_STARTING_FAIL",
									"MITIGATING",
									"STARTING",
									"STARTING_FAIL",
									"START_WAIT",
									"ack_req",
									"alarm",
									"archived",
									"clear",
									"clearing",
									"clearing_fail",
									"end_grace",
									"end_wait",
									"manual_clear",
									"manual_clearing",
									"manual_clearing_fail",
									"manual_mitigating",
									"manual_starting",
									"manual_starting_fail",
									"mitigating",
									"start_wait",
									"starting",
									"starting_fail",
								),
							},
						},
						"alert_duration": schema.StringAttribute{
							Description: "Total alert duration",
							Computed:    true,
						},
						"destination_ip": schema.StringAttribute{
							Description: "Notification destination IP address",
							Computed:    true,
						},
						"id": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
			"fixed_ip_assignments": schema.ListNestedAttribute{
				Description: "Fixed IP assigned to instance",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudBaremetalServerFixedIPAssignmentsModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"external": schema.BoolAttribute{
							Description: "Is network external",
							Computed:    true,
						},
						"ip_address": schema.StringAttribute{
							Description: "Ip address",
							Computed:    true,
						},
						"subnet_id": schema.StringAttribute{
							Description: "Interface subnet id",
							Computed:    true,
						},
					},
				},
			},
			"instance_isolation": schema.SingleNestedAttribute{
				Description: "Instance isolation information",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudBaremetalServerInstanceIsolationModel](ctx),
				Attributes: map[string]schema.Attribute{
					"reason": schema.StringAttribute{
						Description: "The reason of instance isolation if it is isolated from external internet.",
						Computed:    true,
					},
				},
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

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudInstanceResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Instances are cloud virtual machines with configurable CPU, memory, storage, and networking, supporting various operating systems and workloads.",
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
				Description: "A list of network interfaces for the instance. You can create one or more interfaces - private, public, or both.",
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
						"security_groups": schema.ListNestedAttribute{
							Description: "Specifies security group UUIDs to be applied to the instance network interface.",
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
			"volumes": schema.ListNestedAttribute{
				Description: "List of volumes that will be attached to the instance.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"size": schema.Int64Attribute{
							Description: "Volume size in GiB.",
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.AtLeast(1),
							},
						},
						"source": schema.StringAttribute{
							Description: "New volume will be created from scratch and attached to the instance.\nAvailable values: \"new-volume\", \"image\", \"snapshot\", \"apptemplate\", \"existing-volume\".",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"new-volume",
									"image",
									"snapshot",
									"apptemplate",
									"existing-volume",
								),
							},
						},
						"attachment_tag": schema.StringAttribute{
							Description: "Block device attachment tag (not exposed in the normal tags)",
							Optional:    true,
						},
						"delete_on_termination": schema.BoolAttribute{
							Description: "Set to `true` to automatically delete the volume when the instance is deleted.",
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
						},
						"name": schema.StringAttribute{
							Description: "The name of the volume. If not specified, a name will be generated automatically.",
							Optional:    true,
						},
						"tags": schema.MapAttribute{
							Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"type_name": schema.StringAttribute{
							Description: "Volume type name. Supported values:\n- `standard` - Network SSD block storage offering stable performance with high random I/O and data reliability (6 IOPS per 1 GiB, 0.4 MB/s per 1 GiB). Max IOPS: 4500. Max bandwidth: 300 MB/s.\n- `ssd_hiiops` - High-performance SSD storage for latency-sensitive transactional workloads (60 IOPS per 1 GiB, 2.5 MB/s per 1 GiB). Max IOPS: 9000. Max bandwidth: 500 MB/s.\n- `ssd_lowlatency` - SSD storage optimized for low-latency and real-time processing. Max IOPS: 5000. Average latency: 300 µs. Snapshots and volume resizing are **not** supported for `ssd_lowlatency`.\nAvailable values: \"cold\", \"ssd_hiiops\", \"ssd_local\", \"ssd_lowlatency\", \"standard\", \"ultra\".",
							Optional:    true,
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
						"image_id": schema.StringAttribute{
							Description: "Image ID.",
							Optional:    true,
						},
						"boot_index": schema.Int64Attribute{
							Description: "- `0` means that this is the primary boot device;\n- A unique positive value is set for the secondary bootable devices;\n- A negative number means that the boot is prohibited.",
							Optional:    true,
						},
						"snapshot_id": schema.StringAttribute{
							Description: "Snapshot ID.",
							Optional:    true,
						},
						"apptemplate_id": schema.StringAttribute{
							Description: "App template ID.",
							Optional:    true,
						},
						"volume_id": schema.StringAttribute{
							Description: "Volume ID.",
							Optional:    true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"allow_app_ports": schema.BoolAttribute{
				Description:   "Set to `true` if creating the instance from an `apptemplate`. This allows application ports in the security group for instances created from a marketplace application template.",
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"name_template": schema.StringAttribute{
				Description:   "If you want the instance name to be automatically generated based on IP addresses, you can provide a name template instead of specifying the name manually. The template should include a placeholder that will be replaced during provisioning. Supported placeholders are: `{ip_octets}` (last 3 octets of the IP), `{two_ip_octets}`, and `{one_ip_octet}`.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"password": schema.StringAttribute{
				Description:   "For Linux instances, 'username' and 'password' are used to create a new user. When only 'password' is provided, it is set as the password for the default user of the image. For Windows instances, 'username' cannot be specified. Use the 'password' field to set the password for the 'Admin' user on Windows. Use the 'user_data' field to provide a script to create new users on Windows. The password of the Admin user cannot be updated via 'user_data'.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"servergroup_id": schema.StringAttribute{
				Description:   "Placement group ID for instance placement policy.\n\nSupported group types:\n- `anti-affinity`: Ensures instances are placed on different hosts for high availability.\n- `affinity`: Places instances on the same host for low-latency communication.\n- `soft-anti-affinity`: Tries to place instances on different hosts but allows sharing if needed.",
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
			"configuration": schema.MapAttribute{
				Description:   "Parameters for the application template if creating the instance from an `apptemplate`.",
				Optional:      true,
				ElementType:   jsontypes.NormalizedType{},
				PlanModifiers: []planmodifier.Map{mapplanmodifier.RequiresReplace()},
			},
			"security_groups": schema.ListNestedAttribute{
				Description: "Specifies security group UUIDs to be applied to all instance network interfaces.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Resource ID",
							Required:    true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description: "Instance name.",
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
				Description: "Datetime when instance was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"instance_description": schema.StringAttribute{
				Description: "Instance description",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Instance status\nAvailable values: \"ACTIVE\", \"BUILD\", \"DELETED\", \"ERROR\", \"HARD_REBOOT\", \"MIGRATING\", \"PASSWORD\", \"PAUSED\", \"REBOOT\", \"REBUILD\", \"RESCUE\", \"RESIZE\", \"REVERT_RESIZE\", \"SHELVED\", \"SHELVED_OFFLOADED\", \"SHUTOFF\", \"SOFT_DELETED\", \"SUSPENDED\", \"UNKNOWN\", \"VERIFY_RESIZE\".",
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
				Description: "Virtual machine state (active)\nAvailable values: \"active\", \"building\", \"deleted\", \"error\", \"paused\", \"rescued\", \"resized\", \"shelved\", \"shelved_offloaded\", \"soft-deleted\", \"stopped\", \"suspended\".",
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
				CustomType:  customfield.NewMapType[customfield.NestedObjectList[CloudInstanceAddressesModel]](ctx),
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
							Description: "Interface name. This field will be `null` if `with_interfaces_name=true` is not set in the request when listing instances. It will also be `null` if the `interface_name` was not specified during instance creation or when attaching the interface.",
							Optional:    true,
						}.GetType(), "subnet_id": schema.StringAttribute{
							Description: "The unique identifier of the subnet associated with this address. Included only in the response for a single-resource lookup (GET by ID). For the trunk subports, this field is always set.",
							Optional:    true,
						}.GetType(), "subnet_name": schema.StringAttribute{
							Description: "The name of the subnet associated with this address. Included only in the response for a single-resource lookup (GET by ID). For the trunk subports, this field is always set.",
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
				CustomType:  customfield.NewNestedObjectListType[CloudInstanceBlackholePortsModel](ctx),
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
			"ddos_profile": schema.SingleNestedAttribute{
				Description: "Advanced DDoS protection profile. It is always `null` if query parameter `with_ddos=true` is not set.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudInstanceDDOSProfileModel](ctx),
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "Unique identifier for the DDoS protection profile",
						Computed:    true,
					},
					"fields": schema.ListNestedAttribute{
						Description: "List of configured field values for the protection profile",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectListType[CloudInstanceDDOSProfileFieldsModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Description: "Unique identifier for the DDoS protection field",
									Computed:    true,
								},
								"base_field": schema.Int64Attribute{
									Description: "ID of DDoS profile field",
									Computed:    true,
								},
								"default": schema.StringAttribute{
									Description: "Predefined default value for the field if not specified",
									Computed:    true,
								},
								"description": schema.StringAttribute{
									Description: "Detailed description explaining the field's purpose and usage guidelines",
									Computed:    true,
								},
								"field_name": schema.StringAttribute{
									Description: "Name of DDoS profile field",
									Computed:    true,
								},
								"field_type": schema.StringAttribute{
									Description: "Data type classification of the field (e.g., string, integer, array)",
									Computed:    true,
								},
								"field_value": schema.StringAttribute{
									Description: "Complex value. Only one of 'value' or 'field_value' must be specified.",
									Computed:    true,
									CustomType:  jsontypes.NormalizedType{},
								},
								"name": schema.StringAttribute{
									Description: "Human-readable name of the protection field",
									Computed:    true,
								},
								"required": schema.BoolAttribute{
									Description: "Indicates whether this field must be provided when creating a protection profile",
									Computed:    true,
								},
								"validation_schema": schema.StringAttribute{
									Description: "JSON schema defining validation rules and constraints for the field value",
									Computed:    true,
									CustomType:  jsontypes.NormalizedType{},
								},
								"value": schema.StringAttribute{
									Description: "Basic type value. Only one of 'value' or 'field_value' must be specified.",
									Computed:    true,
								},
							},
						},
					},
					"options": schema.SingleNestedAttribute{
						Description: "Configuration options controlling profile activation and BGP routing",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudInstanceDDOSProfileOptionsModel](ctx),
						Attributes: map[string]schema.Attribute{
							"active": schema.BoolAttribute{
								Description: "Controls whether the DDoS protection profile is enabled and actively protecting the resource",
								Computed:    true,
							},
							"bgp": schema.BoolAttribute{
								Description: "Enables Border Gateway Protocol (BGP) routing for DDoS protection traffic",
								Computed:    true,
							},
						},
					},
					"profile_template": schema.SingleNestedAttribute{
						Description: "Complete template configuration data used for this profile",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudInstanceDDOSProfileProfileTemplateModel](ctx),
						Attributes: map[string]schema.Attribute{
							"id": schema.Int64Attribute{
								Description: "Unique identifier for the DDoS protection template",
								Computed:    true,
							},
							"description": schema.StringAttribute{
								Description: "Detailed description explaining the template's purpose and use cases",
								Computed:    true,
							},
							"fields": schema.ListNestedAttribute{
								Description: "List of configurable fields that define the template's protection parameters",
								Computed:    true,
								CustomType:  customfield.NewNestedObjectListType[CloudInstanceDDOSProfileProfileTemplateFieldsModel](ctx),
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.Int64Attribute{
											Description: "Unique identifier for the DDoS protection field",
											Computed:    true,
										},
										"default": schema.StringAttribute{
											Description: "Predefined default value for the field if not specified",
											Computed:    true,
										},
										"description": schema.StringAttribute{
											Description: "Detailed description explaining the field's purpose and usage guidelines",
											Computed:    true,
										},
										"field_type": schema.StringAttribute{
											Description: "Data type classification of the field (e.g., string, integer, array)",
											Computed:    true,
										},
										"name": schema.StringAttribute{
											Description: "Human-readable name of the protection field",
											Computed:    true,
										},
										"required": schema.BoolAttribute{
											Description: "Indicates whether this field must be provided when creating a protection profile",
											Computed:    true,
										},
										"validation_schema": schema.StringAttribute{
											Description: "JSON schema defining validation rules and constraints for the field value",
											Computed:    true,
											CustomType:  jsontypes.NormalizedType{},
										},
									},
								},
							},
							"name": schema.StringAttribute{
								Description: "Human-readable name of the protection template",
								Computed:    true,
							},
						},
					},
					"profile_template_description": schema.StringAttribute{
						Description: "Detailed description of the protection template used for this profile",
						Computed:    true,
					},
					"protocols": schema.ListNestedAttribute{
						Description: "List of network protocols and ports configured for protection",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectListType[CloudInstanceDDOSProfileProtocolsModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"port": schema.StringAttribute{
									Description: "Network port number for which protocols are configured",
									Computed:    true,
								},
								"protocols": schema.ListAttribute{
									Description: "List of network protocols enabled on the specified port",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
							},
						},
					},
					"site": schema.StringAttribute{
						Description: "Geographic site identifier where the protection is deployed",
						Computed:    true,
					},
					"status": schema.SingleNestedAttribute{
						Description: "Current operational status and any error information for the profile",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudInstanceDDOSProfileStatusModel](ctx),
						Attributes: map[string]schema.Attribute{
							"error_description": schema.StringAttribute{
								Description: "Detailed error message describing any issues with the profile operation",
								Computed:    true,
							},
							"status": schema.StringAttribute{
								Description: "Current operational status of the DDoS protection profile",
								Computed:    true,
							},
						},
					},
				},
			},
			"fixed_ip_assignments": schema.ListNestedAttribute{
				Description: "Fixed IP assigned to instance",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudInstanceFixedIPAssignmentsModel](ctx),
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
				CustomType:  customfield.NewNestedObjectType[CloudInstanceInstanceIsolationModel](ctx),
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

func (r *CloudInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudInstanceResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

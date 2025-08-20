// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster

import (
	"context"

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*CloudGPUBaremetalClusterResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"region_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"flavor": schema.StringAttribute{
				Description:   "Flavor name",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"image_id": schema.StringAttribute{
				Description:   "Image ID",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "GPU Cluster name",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"interfaces": schema.ListNestedAttribute{
				Description: "A list of network interfaces for the server. You can create one or more interfaces - private, public, or both.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "A public IP address will be assigned to the server.\nAvailable values: \"external\", \"subnet\", \"any_subnet\".",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"external",
									"subnet",
									"any_subnet",
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
						"network_id": schema.StringAttribute{
							Description: "The network where the server will be connected.",
							Optional:    true,
						},
						"subnet_id": schema.StringAttribute{
							Description: "The server will get an IP address from this subnet.",
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
						"ip_address": schema.StringAttribute{
							Description: "You can specify a specific IP address from your subnet.",
							Optional:    true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"password": schema.StringAttribute{
				Description:   `A password for a bare metal server. This parameter is used to set a password for the "Admin" user on a Windows instance, a default user or a new user on a Linux instance`,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"ssh_key_name": schema.StringAttribute{
				Description:   "Specifies the name of the SSH keypair, created via the\n[/v1/`ssh_keys` endpoint](/docs/api-reference/cloud/ssh-keys/add-or-generate-ssh-key).",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"user_data": schema.StringAttribute{
				Description:   "String in base64 format. Must not be passed together with 'username' or 'password'. Examples of the `user_data`: https://cloudinit.readthedocs.io/en/latest/topics/examples.html",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"username": schema.StringAttribute{
				Description:   "A name of a new user in the Linux instance. It may be passed with a 'password' parameter",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"tags": schema.MapAttribute{
				Description:   "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Optional:      true,
				ElementType:   types.StringType,
				PlanModifiers: []planmodifier.Map{mapplanmodifier.RequiresReplace()},
			},
			"security_groups": schema.ListNestedAttribute{
				Description: "Security group UUIDs",
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
			"instances_count": schema.Int64Attribute{
				Description: "Number of servers to create",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
				Default:       int64default.StaticInt64(1),
			},
			"cluster_name": schema.StringAttribute{
				Description: "GPU Cluster Name",
				Computed:    true,
			},
			"cluster_status": schema.StringAttribute{
				Description: "GPU Cluster status\nAvailable values: \"ACTIVE\", \"ERROR\", \"PENDING\", \"SUSPENDED\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"ACTIVE",
						"ERROR",
						"PENDING",
						"SUSPENDED",
					),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the cluster was created",
				Computed:    true,
			},
			"creator_task_id": schema.StringAttribute{
				Description:        "Task that created this entity",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated.",
			},
			"image_name": schema.StringAttribute{
				Description: "Image name",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"task_id": schema.StringAttribute{
				Description: "Task ID associated with the cluster",
				Computed:    true,
			},
			"task_status": schema.StringAttribute{
				Description: "Task status\nAvailable values: \"CLUSTER_CLEAN_UP\", \"CLUSTER_RESIZE\", \"CLUSTER_RESUME\", \"CLUSTER_SUSPEND\", \"ERROR\", \"FINISHED\", \"IPU_SERVERS\", \"NETWORK\", \"POPLAR_SERVERS\", \"POST_DEPLOY_SETUP\", \"VIPU_CONTROLLER\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"CLUSTER_CLEAN_UP",
						"CLUSTER_RESIZE",
						"CLUSTER_RESUME",
						"CLUSTER_SUSPEND",
						"ERROR",
						"FINISHED",
						"IPU_SERVERS",
						"NETWORK",
						"POPLAR_SERVERS",
						"POST_DEPLOY_SETUP",
						"VIPU_CONTROLLER",
					),
				},
			},
			"tasks": schema.ListAttribute{
				Description: "List of task IDs",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"servers": schema.ListNestedAttribute{
				Description: "GPU cluster servers",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "GPU server ID",
							Computed:    true,
						},
						"addresses": schema.MapAttribute{
							Description: "Map of `network_name` to list of addresses in that network",
							Computed:    true,
							CustomType:  customfield.NewMapType[customfield.NestedObjectList[CloudGPUBaremetalClusterServersAddressesModel]](ctx),
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
						"blackhole_ports": schema.ListNestedAttribute{
							Description: "IP addresses of the instances that are blackholed by DDoS mitigation system",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersBlackholePortsModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"alarm_end": schema.StringAttribute{
										Description: "A date-time string giving the time that the alarm ended",
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
						"created_at": schema.StringAttribute{
							Description: "Datetime when GPU server was created",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"creator_task_id": schema.StringAttribute{
							Description: "Task that created this entity",
							Computed:    true,
						},
						"ddos_profile": schema.SingleNestedAttribute{
							Description: "Advanced DDoS protection profile. It is always `null` if query parameter `with_ddos=true` is not set.",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersDDOSProfileModel](ctx),
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Description: "Unique identifier for the DDoS protection profile",
									Computed:    true,
								},
								"fields": schema.ListNestedAttribute{
									Description: "List of configured field values for the protection profile",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersDDOSProfileFieldsModel](ctx),
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
												Description: "Complex value. Only one of 'value' or '`field_value`' must be specified.",
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
												Description: "Basic type value. Only one of 'value' or '`field_value`' must be specified.",
												Computed:    true,
											},
										},
									},
								},
								"options": schema.SingleNestedAttribute{
									Description: "Configuration options controlling profile activation and BGP routing",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersDDOSProfileOptionsModel](ctx),
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
									CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersDDOSProfileProfileTemplateModel](ctx),
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
											CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersDDOSProfileProfileTemplateFieldsModel](ctx),
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
									CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersDDOSProfileProtocolsModel](ctx),
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
									CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersDDOSProfileStatusModel](ctx),
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
							CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersFixedIPAssignmentsModel](ctx),
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
						"flavor": schema.SingleNestedAttribute{
							Description: "Flavor",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersFlavorModel](ctx),
							Attributes: map[string]schema.Attribute{
								"architecture": schema.StringAttribute{
									Description: "CPU architecture",
									Computed:    true,
								},
								"flavor_id": schema.StringAttribute{
									Description: "Flavor ID is the same as name",
									Computed:    true,
								},
								"flavor_name": schema.StringAttribute{
									Description: "Flavor name",
									Computed:    true,
								},
								"hardware_description": schema.SingleNestedAttribute{
									Description: "Additional hardware description",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersFlavorHardwareDescriptionModel](ctx),
									Attributes: map[string]schema.Attribute{
										"cpu": schema.StringAttribute{
											Description: "Human-readable CPU description",
											Computed:    true,
										},
										"disk": schema.StringAttribute{
											Description: "Human-readable disk description",
											Computed:    true,
										},
										"gpu": schema.StringAttribute{
											Description: "Human-readable GPU description",
											Computed:    true,
										},
										"license": schema.StringAttribute{
											Description: "If the flavor is licensed, this field contains the license type",
											Computed:    true,
										},
										"network": schema.StringAttribute{
											Description: "Human-readable NIC description",
											Computed:    true,
										},
										"ram": schema.StringAttribute{
											Description: "Human-readable RAM description",
											Computed:    true,
										},
									},
								},
								"os_type": schema.StringAttribute{
									Description: "Operating system",
									Computed:    true,
								},
								"ram": schema.Int64Attribute{
									Description: "RAM size in MiB",
									Computed:    true,
								},
								"resource_class": schema.StringAttribute{
									Description: "Flavor resource class for mapping to hardware capacity",
									Computed:    true,
								},
								"vcpus": schema.Int64Attribute{
									Description: "Virtual CPU count. For bare metal flavors, it's a physical CPU count",
									Computed:    true,
								},
							},
						},
						"instance_description": schema.StringAttribute{
							Description: "Instance description",
							Computed:    true,
						},
						"instance_isolation": schema.SingleNestedAttribute{
							Description: "Instance isolation information",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersInstanceIsolationModel](ctx),
							Attributes: map[string]schema.Attribute{
								"reason": schema.StringAttribute{
									Description: "The reason of instance isolation if it is isolated from external internet.",
									Computed:    true,
								},
							},
						},
						"name": schema.StringAttribute{
							Description: "GPU server name",
							Computed:    true,
						},
						"project_id": schema.Int64Attribute{
							Description: "Project ID",
							Computed:    true,
						},
						"region": schema.StringAttribute{
							Description: "Region name",
							Computed:    true,
						},
						"region_id": schema.Int64Attribute{
							Description: "Region ID",
							Computed:    true,
						},
						"security_groups": schema.ListNestedAttribute{
							Description: "Security groups",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersSecurityGroupsModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Name.",
										Computed:    true,
									},
								},
							},
						},
						"ssh_key_name": schema.StringAttribute{
							Description: "SSH key name assigned to instance",
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
						"tags": schema.ListNestedAttribute{
							Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersTagsModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Description: "Tag key. The maximum size for a key is 255 bytes.",
										Computed:    true,
									},
									"read_only": schema.BoolAttribute{
										Description: "If true, the tag is read-only and cannot be modified by the user",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Tag value. The maximum size for a value is 1024 bytes.",
										Computed:    true,
									},
								},
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
					},
				},
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

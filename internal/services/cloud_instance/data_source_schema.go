// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInstanceDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Instances are cloud virtual machines with configurable CPU, memory, storage, and networking, supporting various operating systems and workloads.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Instance ID",
				Computed:    true,
			},
			"instance_id": schema.StringAttribute{
				Description: "Instance ID",
				Optional:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Optional:    true,
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
			"name": schema.StringAttribute{
				Description: "Instance name",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"ssh_key_name": schema.StringAttribute{
				Description: "SSH key assigned to instance",
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
				CustomType:  customfield.NewMapType[customfield.NestedObjectList[CloudInstanceAddressesDataSourceModel]](ctx),
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
				CustomType:  customfield.NewNestedObjectListType[CloudInstanceBlackholePortsDataSourceModel](ctx),
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
				CustomType:  customfield.NewNestedObjectType[CloudInstanceDDOSProfileDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "Unique identifier for the DDoS protection profile",
						Computed:    true,
					},
					"fields": schema.ListNestedAttribute{
						Description: "List of configured field values for the protection profile",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectListType[CloudInstanceDDOSProfileFieldsDataSourceModel](ctx),
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
								"field_type": schema.StringAttribute{
									Description: "Data type classification of the field (e.g., string, integer, array)",
									Computed:    true,
								},
								"field_value": schema.StringAttribute{
									Description: "Complex value for the DDoS profile field",
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
							},
						},
					},
					"options": schema.SingleNestedAttribute{
						Description: "Configuration options controlling profile activation and BGP routing",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudInstanceDDOSProfileOptionsDataSourceModel](ctx),
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
						CustomType:  customfield.NewNestedObjectType[CloudInstanceDDOSProfileProfileTemplateDataSourceModel](ctx),
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
								CustomType:  customfield.NewNestedObjectListType[CloudInstanceDDOSProfileProfileTemplateFieldsDataSourceModel](ctx),
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
						CustomType:  customfield.NewNestedObjectListType[CloudInstanceDDOSProfileProtocolsDataSourceModel](ctx),
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
						CustomType:  customfield.NewNestedObjectType[CloudInstanceDDOSProfileStatusDataSourceModel](ctx),
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
				CustomType:  customfield.NewNestedObjectListType[CloudInstanceFixedIPAssignmentsDataSourceModel](ctx),
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
				CustomType:  customfield.NewNestedObjectType[CloudInstanceFlavorDataSourceModel](ctx),
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
						CustomType:  customfield.NewNestedObjectType[CloudInstanceFlavorHardwareDescriptionDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"ram": schema.StringAttribute{
								Description: "RAM description",
								Computed:    true,
							},
							"vcpus": schema.StringAttribute{
								Description: "CPU description",
								Computed:    true,
							},
							"cpu": schema.StringAttribute{
								Description: "Human-readable CPU description",
								Computed:    true,
							},
							"disk": schema.StringAttribute{
								Description: "Human-readable disk description",
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
							"gpu": schema.StringAttribute{
								Description: "Human-readable GPU description",
								Computed:    true,
							},
						},
					},
					"os_type": schema.StringAttribute{
						Description: "Flavor operating system",
						Computed:    true,
					},
					"ram": schema.Int64Attribute{
						Description: "RAM size in MiB",
						Computed:    true,
					},
					"vcpus": schema.Int64Attribute{
						Description: "Virtual CPU count. For bare metal flavors, it's a physical CPU count",
						Computed:    true,
					},
					"resource_class": schema.StringAttribute{
						Description: "Flavor resource class for mapping to hardware capacity",
						Computed:    true,
					},
				},
			},
			"instance_isolation": schema.SingleNestedAttribute{
				Description: "Instance isolation information",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudInstanceInstanceIsolationDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"reason": schema.StringAttribute{
						Description: "The reason of instance isolation if it is isolated from external internet.",
						Computed:    true,
					},
				},
			},
			"security_groups": schema.ListNestedAttribute{
				Description: "Security groups",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudInstanceSecurityGroupsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name.",
							Computed:    true,
						},
					},
				},
			},
			"tags": schema.ListNestedAttribute{
				Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudInstanceTagsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "Tag key. Maximum 255 characters. Cannot contain spaces, tabs, newlines, empty string or '=' character.",
							Computed:    true,
						},
						"read_only": schema.BoolAttribute{
							Description: "If true, the tag is read-only and cannot be modified by the user",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "Tag value. Maximum 255 characters. Cannot contain spaces, tabs, newlines, empty string or '=' character.",
							Computed:    true,
						},
					},
				},
			},
			"volumes": schema.ListNestedAttribute{
				Description: "List of volumes",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudInstanceVolumesDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Volume ID",
							Computed:    true,
						},
						"delete_on_termination": schema.BoolAttribute{
							Description: "Whether the volume is deleted together with the VM",
							Computed:    true,
						},
					},
				},
			},
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"available_floating": schema.BoolAttribute{
						Description: "Only show instances which are able to handle floating address",
						Optional:    true,
					},
					"changes_before": schema.StringAttribute{
						Description: "Filters the instances by a date and time stamp when the instances last changed.",
						Optional:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"changes_since": schema.StringAttribute{
						Description: "Filters the instances by a date and time stamp when the instances last changed status.",
						Optional:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"exclude_flavor_prefix": schema.StringAttribute{
						Description: "Exclude instances with specified flavor prefix",
						Optional:    true,
					},
					"exclude_secgroup": schema.StringAttribute{
						Description: "Exclude instances with specified security group name",
						Optional:    true,
					},
					"flavor_id": schema.StringAttribute{
						Description: "Filter out instances by `flavor_id`. Flavor id must match exactly.",
						Optional:    true,
					},
					"flavor_prefix": schema.StringAttribute{
						Description: "Filter out instances by `flavor_prefix`.",
						Optional:    true,
					},
					"include_ai": schema.BoolAttribute{
						Description:        "Include GPU clusters' servers",
						Computed:           true,
						Optional:           true,
						DeprecationMessage: "This attribute is deprecated.",
					},
					"include_baremetal": schema.BoolAttribute{
						Description:        "Include bare metal servers. Please, use `GET /v1/bminstances/` instead",
						Computed:           true,
						Optional:           true,
						DeprecationMessage: "This attribute is deprecated.",
					},
					"include_k8s": schema.BoolAttribute{
						Description: "Include managed k8s worker nodes",
						Computed:    true,
						Optional:    true,
					},
					"ip": schema.StringAttribute{
						Description: "An IPv4 address to filter results by. Note: partial matches are allowed. For example, searching for 192.168.0.1 will return 192.168.0.1, 192.168.0.10, 192.168.0.110, and so on.",
						Optional:    true,
					},
					"limit": schema.Int64Attribute{
						Description: "Optional. Limit the number of returned items",
						Computed:    true,
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.AtMost(1000),
						},
					},
					"name": schema.StringAttribute{
						Description: "Filter instances by name. You can provide a full or partial name, instances with matching names will be returned. For example, entering 'test' will return all instances that contain 'test' in their name.",
						Optional:    true,
					},
					"only_isolated": schema.BoolAttribute{
						Description: "Include only isolated instances",
						Computed:    true,
						Optional:    true,
					},
					"only_with_fixed_external_ip": schema.BoolAttribute{
						Description: "Return bare metals only with external fixed IP addresses.",
						Optional:    true,
					},
					"order_by": schema.StringAttribute{
						Description: "Order by field and direction.\nAvailable values: \"created.asc\", \"created.desc\", \"name.asc\", \"name.desc\", \"status.asc\", \"status.desc\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"created.asc",
								"created.desc",
								"name.asc",
								"name.desc",
								"status.asc",
								"status.desc",
							),
						},
					},
					"profile_name": schema.StringAttribute{
						Description: "Filter result by ddos protection profile name. Effective only with `with_ddos` set to true.",
						Optional:    true,
					},
					"protection_status": schema.StringAttribute{
						Description: "Filter result by DDoS `protection_status`. if parameter is provided. Effective only with `with_ddos` set to true. (Active, Queued or Error)\nAvailable values: \"Active\", \"Queued\", \"Error\".",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"Active",
								"Queued",
								"Error",
							),
						},
					},
					"status": schema.StringAttribute{
						Description: "Filters instances by status.\nAvailable values: \"ACTIVE\", \"BUILD\", \"ERROR\", \"HARD_REBOOT\", \"MIGRATING\", \"PAUSED\", \"REBOOT\", \"REBUILD\", \"RESIZE\", \"REVERT_RESIZE\", \"SHELVED\", \"SHELVED_OFFLOADED\", \"SHUTOFF\", \"SOFT_DELETED\", \"SUSPENDED\", \"VERIFY_RESIZE\".",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"ACTIVE",
								"BUILD",
								"ERROR",
								"HARD_REBOOT",
								"MIGRATING",
								"PAUSED",
								"REBOOT",
								"REBUILD",
								"RESIZE",
								"REVERT_RESIZE",
								"SHELVED",
								"SHELVED_OFFLOADED",
								"SHUTOFF",
								"SOFT_DELETED",
								"SUSPENDED",
								"VERIFY_RESIZE",
							),
						},
					},
					"tag_key_value": schema.StringAttribute{
						Description: "Optional. Filter by tag key-value pairs.",
						Optional:    true,
					},
					"tag_value": schema.ListAttribute{
						Description: "Optional. Filter by tag values. ?`tag_value`=value1&`tag_value`=value2",
						Optional:    true,
						ElementType: types.StringType,
					},
					"type_ddos_profile": schema.StringAttribute{
						Description: "Return bare metals either only with advanced or only basic DDoS protection. Effective only with `with_ddos` set to true. (advanced or basic)\nAvailable values: \"basic\", \"advanced\".",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("basic", "advanced"),
						},
					},
					"uuid": schema.StringAttribute{
						Description: "Filter the server list result by the UUID of the server. Allowed UUID part",
						Optional:    true,
					},
					"with_ddos": schema.BoolAttribute{
						Description: "Include DDoS profile information in the response when set to `true`. Otherwise, the `ddos_profile` field in the response is `null` by default.",
						Computed:    true,
						Optional:    true,
					},
					"with_interfaces_name": schema.BoolAttribute{
						Description: "Include `interface_name` in the addresses",
						Computed:    true,
						Optional:    true,
					},
				},
			},
		},
	}
}

func (d *CloudInstanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudInstanceDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("instance_id"), path.MatchRoot("find_one_by")),
	}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_floating_ip

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudFloatingIPsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "A floating IP is a static IP address that points to one of your Instances. It allows you to redirect network traffic to any of your Instances in the same datacenter.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Filter by floating IP status. DOWN - unassigned (available). ACTIVE - attached to a port (in use). ERROR - error state.\nAvailable values: \"ACTIVE\", \"DOWN\", \"ERROR\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"ACTIVE",
						"DOWN",
						"ERROR",
					),
				},
			},
			"tag_key_value": schema.StringAttribute{
				Description: "Optional. Filter by tag key-value pairs.",
				Optional:    true,
			},
			"tag_key": schema.ListAttribute{
				Description: "Optional. Filter by tag keys. ?`tag_key`=key1&`tag_key`=key2",
				Optional:    true,
				ElementType: types.StringType,
			},
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Floating IP ID",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Datetime when the floating IP was created",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"creator_task_id": schema.StringAttribute{
							Description: "Task that created this entity",
							Computed:    true,
						},
						"fixed_ip_address": schema.StringAttribute{
							Description: "IP address of the port the floating IP is attached to",
							Computed:    true,
						},
						"floating_ip_address": schema.StringAttribute{
							Description: "IP Address of the floating IP",
							Computed:    true,
						},
						"instance": schema.SingleNestedAttribute{
							Description: "Instance the floating IP is attached to",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudFloatingIPsInstanceDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "Instance ID",
									Computed:    true,
								},
								"addresses": schema.MapAttribute{
									Description: "Map of `network_name` to list of addresses in that network",
									Computed:    true,
									CustomType:  customfield.NewMapType[customfield.NestedObjectList[CloudFloatingIPsInstanceAddressesDataSourceModel]](ctx),
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
								"created_at": schema.StringAttribute{
									Description: "Datetime when instance was created",
									Computed:    true,
									CustomType:  timetypes.RFC3339Type{},
								},
								"creator_task_id": schema.StringAttribute{
									Description: "Task that created this entity",
									Computed:    true,
								},
								"flavor": schema.SingleNestedAttribute{
									Description: "Flavor",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectType[CloudFloatingIPsInstanceFlavorDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"flavor_id": schema.StringAttribute{
											Description: "Flavor ID is the same as name",
											Computed:    true,
										},
										"flavor_name": schema.StringAttribute{
											Description: "Flavor name",
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
									},
								},
								"instance_description": schema.StringAttribute{
									Description: "Instance description",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Instance name",
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
									CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsInstanceSecurityGroupsDataSourceModel](ctx),
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
									CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsInstanceTagsDataSourceModel](ctx),
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
								"volumes": schema.ListNestedAttribute{
									Description: "List of volumes",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsInstanceVolumesDataSourceModel](ctx),
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
							},
						},
						"loadbalancer": schema.SingleNestedAttribute{
							Description: "Load balancer the floating IP is attached to",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudFloatingIPsLoadbalancerDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "Load balancer ID",
									Computed:    true,
								},
								"admin_state_up": schema.BoolAttribute{
									Description: "Administrative state of the resource. When set to true, the resource is enabled and operational. When set to false, the resource is disabled and will not process traffic. Defaults to true.",
									Computed:    true,
								},
								"created_at": schema.StringAttribute{
									Description: "Datetime when the load balancer was created",
									Computed:    true,
									CustomType:  timetypes.RFC3339Type{},
								},
								"name": schema.StringAttribute{
									Description: "Load balancer name",
									Computed:    true,
								},
								"operating_status": schema.StringAttribute{
									Description: "Load balancer operating status\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive(
											"DEGRADED",
											"DRAINING",
											"ERROR",
											"NO_MONITOR",
											"OFFLINE",
											"ONLINE",
										),
									},
								},
								"project_id": schema.Int64Attribute{
									Description: "Project ID",
									Computed:    true,
								},
								"provisioning_status": schema.StringAttribute{
									Description: "Load balancer lifecycle status\nAvailable values: \"ACTIVE\", \"DELETED\", \"ERROR\", \"PENDING_CREATE\", \"PENDING_DELETE\", \"PENDING_UPDATE\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive(
											"ACTIVE",
											"DELETED",
											"ERROR",
											"PENDING_CREATE",
											"PENDING_DELETE",
											"PENDING_UPDATE",
										),
									},
								},
								"region": schema.StringAttribute{
									Description: "Region name",
									Computed:    true,
								},
								"region_id": schema.Int64Attribute{
									Description: "Region ID",
									Computed:    true,
								},
								"tags_v2": schema.ListNestedAttribute{
									Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsLoadbalancerTagsV2DataSourceModel](ctx),
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
								"additional_vips": schema.ListNestedAttribute{
									Description: "List of additional IP addresses",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsLoadbalancerAdditionalVipsDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"ip_address": schema.StringAttribute{
												Description: "IP address",
												Computed:    true,
											},
											"subnet_id": schema.StringAttribute{
												Description: "Subnet UUID",
												Computed:    true,
											},
										},
									},
								},
								"creator_task_id": schema.StringAttribute{
									Description: "Task that created this entity",
									Computed:    true,
								},
								"ddos_profile": schema.SingleNestedAttribute{
									Description: "Loadbalancer advanced DDoS protection profile.",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectType[CloudFloatingIPsLoadbalancerDDOSProfileDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"id": schema.Int64Attribute{
											Description: "Unique identifier for the DDoS protection profile",
											Computed:    true,
										},
										"fields": schema.ListNestedAttribute{
											Description: "List of configured field values for the protection profile",
											Computed:    true,
											CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsLoadbalancerDDOSProfileFieldsDataSourceModel](ctx),
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
											CustomType:  customfield.NewNestedObjectType[CloudFloatingIPsLoadbalancerDDOSProfileOptionsDataSourceModel](ctx),
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
											CustomType:  customfield.NewNestedObjectType[CloudFloatingIPsLoadbalancerDDOSProfileProfileTemplateDataSourceModel](ctx),
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
													CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsLoadbalancerDDOSProfileProfileTemplateFieldsDataSourceModel](ctx),
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
											CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsLoadbalancerDDOSProfileProtocolsDataSourceModel](ctx),
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
											CustomType:  customfield.NewNestedObjectType[CloudFloatingIPsLoadbalancerDDOSProfileStatusDataSourceModel](ctx),
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
								"flavor": schema.SingleNestedAttribute{
									Description: "Load balancer flavor (if not default)",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectType[CloudFloatingIPsLoadbalancerFlavorDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"flavor_id": schema.StringAttribute{
											Description: "Flavor ID is the same as name",
											Computed:    true,
										},
										"flavor_name": schema.StringAttribute{
											Description: "Flavor name",
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
									},
								},
								"floating_ips": schema.ListNestedAttribute{
									Description: "List of assigned floating IPs",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsLoadbalancerFloatingIPsDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description: "Floating IP ID",
												Computed:    true,
											},
											"created_at": schema.StringAttribute{
												Description: "Datetime when the floating IP was created",
												Computed:    true,
												CustomType:  timetypes.RFC3339Type{},
											},
											"creator_task_id": schema.StringAttribute{
												Description: "Task that created this entity",
												Computed:    true,
											},
											"fixed_ip_address": schema.StringAttribute{
												Description: "IP address of the port the floating IP is attached to",
												Computed:    true,
											},
											"floating_ip_address": schema.StringAttribute{
												Description: "IP Address of the floating IP",
												Computed:    true,
											},
											"port_id": schema.StringAttribute{
												Description: "Port ID the floating IP is attached to. The `fixed_ip_address` is the IP address of the port.",
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
											"router_id": schema.StringAttribute{
												Description: "Router ID",
												Computed:    true,
											},
											"status": schema.StringAttribute{
												Description: "Floating IP status. DOWN - unassigned (available). ACTIVE - attached to a port (in use). ERROR - error state.\nAvailable values: \"ACTIVE\", \"DOWN\", \"ERROR\".",
												Computed:    true,
												Validators: []validator.String{
													stringvalidator.OneOfCaseInsensitive(
														"ACTIVE",
														"DOWN",
														"ERROR",
													),
												},
											},
											"tags": schema.ListNestedAttribute{
												Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
												Computed:    true,
												CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsLoadbalancerFloatingIPsTagsDataSourceModel](ctx),
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
											"task_id": schema.StringAttribute{
												Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
												Computed:    true,
											},
											"updated_at": schema.StringAttribute{
												Description: "Datetime when the floating IP was last updated",
												Computed:    true,
												CustomType:  timetypes.RFC3339Type{},
											},
										},
									},
								},
								"listeners": schema.ListNestedAttribute{
									Description: "Load balancer listeners",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsLoadbalancerListenersDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description: "Listener ID",
												Computed:    true,
											},
										},
									},
								},
								"logging": schema.SingleNestedAttribute{
									Description: "Logging configuration",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectType[CloudFloatingIPsLoadbalancerLoggingDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"destination_region_id": schema.Int64Attribute{
											Description: "ID of the region in which the logs will be stored",
											Computed:    true,
										},
										"enabled": schema.BoolAttribute{
											Description: "Indicates if log streaming is enabled or disabled",
											Computed:    true,
										},
										"topic_name": schema.StringAttribute{
											Description: "The topic name to stream logs to",
											Computed:    true,
										},
										"retention_policy": schema.SingleNestedAttribute{
											Description: "Logs retention policy",
											Computed:    true,
											CustomType:  customfield.NewNestedObjectType[CloudFloatingIPsLoadbalancerLoggingRetentionPolicyDataSourceModel](ctx),
											Attributes: map[string]schema.Attribute{
												"period": schema.Int64Attribute{
													Description: "Duration of days for which logs must be kept.",
													Computed:    true,
													Validators: []validator.Int64{
														int64validator.AtMost(1825),
													},
												},
											},
										},
									},
								},
								"preferred_connectivity": schema.StringAttribute{
									Description: "Preferred option to establish connectivity between load balancer and its pools members\nAvailable values: \"L2\", \"L3\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("L2", "L3"),
									},
								},
								"stats": schema.SingleNestedAttribute{
									Description: "Statistics of load balancer.",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectType[CloudFloatingIPsLoadbalancerStatsDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"active_connections": schema.Int64Attribute{
											Description: "Currently active connections",
											Computed:    true,
										},
										"bytes_in": schema.Int64Attribute{
											Description: "Total bytes received",
											Computed:    true,
										},
										"bytes_out": schema.Int64Attribute{
											Description: "Total bytes sent",
											Computed:    true,
										},
										"request_errors": schema.Int64Attribute{
											Description: "Total requests that were unable to be fulfilled",
											Computed:    true,
										},
										"total_connections": schema.Int64Attribute{
											Description: "Total connections handled",
											Computed:    true,
										},
									},
								},
								"task_id": schema.StringAttribute{
									Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
									Computed:    true,
								},
								"updated_at": schema.StringAttribute{
									Description: "Datetime when the load balancer was last updated",
									Computed:    true,
									CustomType:  timetypes.RFC3339Type{},
								},
								"vip_address": schema.StringAttribute{
									Description: "Load balancer IP address",
									Computed:    true,
								},
								"vip_fqdn": schema.StringAttribute{
									Description: "Fully qualified domain name for the load balancer VIP",
									Computed:    true,
								},
								"vip_ip_family": schema.StringAttribute{
									Description: "Load balancer IP family\nAvailable values: \"dual\", \"ipv4\", \"ipv6\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive(
											"dual",
											"ipv4",
											"ipv6",
										),
									},
								},
								"vip_port_id": schema.StringAttribute{
									Description: "The ID of the Virtual IP (VIP) port.",
									Computed:    true,
								},
								"vrrp_ips": schema.ListNestedAttribute{
									Description: "List of VRRP IP addresses",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsLoadbalancerVrrpIPsDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"ip_address": schema.StringAttribute{
												Description: "IP address",
												Computed:    true,
											},
											"role": schema.StringAttribute{
												Description: "LoadBalancer instance role to which VRRP IP belong\nAvailable values: \"BACKUP\", \"MASTER\", \"STANDALONE\".",
												Computed:    true,
												Validators: []validator.String{
													stringvalidator.OneOfCaseInsensitive(
														"BACKUP",
														"MASTER",
														"STANDALONE",
													),
												},
											},
											"subnet_id": schema.StringAttribute{
												Description: "Subnet UUID",
												Computed:    true,
											},
										},
									},
								},
							},
						},
						"port_id": schema.StringAttribute{
							Description: "Port ID",
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
						"router_id": schema.StringAttribute{
							Description: "Router ID",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Floating IP status. DOWN - unassigned (available). ACTIVE - attached to a port (in use). ERROR - error state.\nAvailable values: \"ACTIVE\", \"DOWN\", \"ERROR\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"ACTIVE",
									"DOWN",
									"ERROR",
								),
							},
						},
						"tags": schema.ListNestedAttribute{
							Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudFloatingIPsTagsDataSourceModel](ctx),
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
						"task_id": schema.StringAttribute{
							Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Datetime when the floating IP was last updated",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
					},
				},
			},
		},
	}
}

func (d *CloudFloatingIPsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudFloatingIPsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

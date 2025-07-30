// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudGPUBaremetalClusterDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.Int64Attribute{
				Required: true,
			},
			"region_id": schema.Int64Attribute{
				Required: true,
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
			"flavor": schema.StringAttribute{
				Description: "Flavor ID is the same as the name",
				Computed:    true,
			},
			"image_id": schema.StringAttribute{
				Description: "Image ID",
				Computed:    true,
			},
			"image_name": schema.StringAttribute{
				Description: "Image name",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: `A password for a bare metal server. This parameter is used to set a password for the "Admin" user on a Windows instance, a default user or a new user on a Linux instance`,
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"ssh_key_name": schema.StringAttribute{
				Description: "Keypair name to inject into new cluster(s)",
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
			"user_data": schema.StringAttribute{
				Description: "String in base64 format. Must not be passed together with 'username' or 'password'. Examples of the `user_data`: https://cloudinit.readthedocs.io/en/latest/topics/examples.html",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "A name of a new user in the Linux instance. It may be passed with a 'password' parameter",
				Computed:    true,
			},
			"interfaces": schema.ListNestedAttribute{
				Description: "Networks managed by user and associated with the cluster",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterInterfacesDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"network_id": schema.StringAttribute{
							Description: "Network ID",
							Computed:    true,
						},
						"port_id": schema.StringAttribute{
							Description: "Network ID the subnet belongs to. Port will be plugged in this network",
							Computed:    true,
						},
						"subnet_id": schema.StringAttribute{
							Description: "Port is assigned to IP address from the subnet",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Network type",
							Computed:    true,
						},
					},
				},
			},
			"servers": schema.ListNestedAttribute{
				Description: "GPU cluster servers",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "GPU server ID",
							Computed:    true,
						},
						"addresses": schema.MapAttribute{
							Description: "Map of `network_name` to list of addresses in that network",
							Computed:    true,
							CustomType:  customfield.NewMapType[customfield.NestedObjectList[CloudGPUBaremetalClusterServersAddressesDataSourceModel]](ctx),
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
							CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersBlackholePortsDataSourceModel](ctx),
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
							CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersDDOSProfileDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Description: "DDoS protection profile ID",
									Computed:    true,
								},
								"profile_template": schema.SingleNestedAttribute{
									Description: "Template data",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersDDOSProfileProfileTemplateDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"id": schema.Int64Attribute{
											Computed: true,
										},
										"name": schema.StringAttribute{
											Computed: true,
										},
										"description": schema.StringAttribute{
											Computed: true,
										},
										"fields": schema.ListNestedAttribute{
											Computed:   true,
											CustomType: customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersDDOSProfileProfileTemplateFieldsDataSourceModel](ctx),
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"id": schema.Int64Attribute{
														Computed: true,
													},
													"name": schema.StringAttribute{
														Computed: true,
													},
													"default": schema.StringAttribute{
														Computed: true,
													},
													"description": schema.StringAttribute{
														Computed: true,
													},
													"field_type": schema.StringAttribute{
														Computed: true,
													},
													"required": schema.BoolAttribute{
														Computed: true,
													},
													"validation_schema": schema.StringAttribute{
														Computed:   true,
														CustomType: jsontypes.NormalizedType{},
													},
												},
											},
										},
									},
								},
								"fields": schema.ListNestedAttribute{
									Computed:   true,
									CustomType: customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersDDOSProfileFieldsDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.Int64Attribute{
												Computed: true,
											},
											"default": schema.StringAttribute{
												Computed:   true,
												CustomType: jsontypes.NormalizedType{},
											},
											"description": schema.StringAttribute{
												Computed: true,
											},
											"field_value": schema.StringAttribute{
												Computed:   true,
												CustomType: jsontypes.NormalizedType{},
											},
											"name": schema.StringAttribute{
												Computed: true,
											},
											"base_field": schema.Int64Attribute{
												Computed: true,
											},
											"field_name": schema.StringAttribute{
												Computed: true,
											},
											"field_type": schema.StringAttribute{
												Computed: true,
											},
											"required": schema.BoolAttribute{
												Computed: true,
											},
											"validation_schema": schema.StringAttribute{
												Computed:   true,
												CustomType: jsontypes.NormalizedType{},
											},
											"value": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
								"options": schema.SingleNestedAttribute{
									Computed:   true,
									CustomType: customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersDDOSProfileOptionsDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"active": schema.BoolAttribute{
											Description: "Activate profile.",
											Computed:    true,
										},
										"bgp": schema.BoolAttribute{
											Description: "Activate BGP protocol.",
											Computed:    true,
										},
									},
								},
								"profile_template_description": schema.StringAttribute{
									Description: "DDoS profile template description",
									Computed:    true,
								},
								"protocols": schema.ListNestedAttribute{
									Description: "List of protocols",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersDDOSProfileProtocolsDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"port": schema.StringAttribute{
												Computed: true,
											},
											"protocols": schema.ListAttribute{
												Computed:    true,
												CustomType:  customfield.NewListType[types.String](ctx),
												ElementType: types.StringType,
											},
										},
									},
								},
								"site": schema.StringAttribute{
									Computed: true,
								},
								"status": schema.SingleNestedAttribute{
									Computed:   true,
									CustomType: customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersDDOSProfileStatusDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"error_description": schema.StringAttribute{
											Description: "Description of the error, if it exists",
											Computed:    true,
										},
										"status": schema.StringAttribute{
											Description: "Profile status",
											Computed:    true,
										},
									},
								},
							},
						},
						"fixed_ip_assignments": schema.ListNestedAttribute{
							Description: "Fixed IP assigned to instance",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersFixedIPAssignmentsDataSourceModel](ctx),
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
							CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersFlavorDataSourceModel](ctx),
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
									CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersFlavorHardwareDescriptionDataSourceModel](ctx),
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
							CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClusterServersInstanceIsolationDataSourceModel](ctx),
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
							CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersSecurityGroupsDataSourceModel](ctx),
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
							CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterServersTagsDataSourceModel](ctx),
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
			"tags": schema.ListNestedAttribute{
				Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterTagsDataSourceModel](ctx),
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
		},
	}
}

func (d *CloudGPUBaremetalClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudGPUBaremetalClusterDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

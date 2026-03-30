// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_server

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudBaremetalServersDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Bare metal servers are dedicated physical machines with direct hardware access, supporting provisioning, rebuilding, and network configuration within a cloud region.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
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
			"flavor_id": schema.StringAttribute{
				Description: "Filter out instances by `flavor_id`. Flavor id must match exactly.",
				Optional:    true,
			},
			"flavor_prefix": schema.StringAttribute{
				Description: "Filter out instances by `flavor_prefix`.",
				Optional:    true,
			},
			"ip": schema.StringAttribute{
				Description: "An IPv4 address to filter results by. Note: partial matches are allowed. For example, searching for 192.168.0.1 will return 192.168.0.1, 192.168.0.10, 192.168.0.110, and so on.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Filter instances by name. You can provide a full or partial name, instances with matching names will be returned. For example, entering 'test' will return all instances that contain 'test' in their name.",
				Optional:    true,
			},
			"only_with_fixed_external_ip": schema.BoolAttribute{
				Description: "Return bare metals only with external fixed IP addresses.",
				Optional:    true,
			},
			"profile_name": schema.StringAttribute{
				Description: "Filter result by ddos protection profile name. Effective only with `with_ddos` set to true.",
				Optional:    true,
			},
			"protection_status": schema.StringAttribute{
				Description: "Filter result by DDoS `protection_status`. Effective only with `with_ddos` set to true. (Active, Queued or Error)\nAvailable values: \"Active\", \"Queued\", \"Error\".",
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
				Description: "Filters instances by a server status, as a string.\nAvailable values: \"ACTIVE\", \"BUILD\", \"ERROR\", \"HARD_REBOOT\", \"REBOOT\", \"REBUILD\", \"RESCUE\", \"SHUTOFF\", \"SUSPENDED\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"ACTIVE",
						"BUILD",
						"ERROR",
						"HARD_REBOOT",
						"REBOOT",
						"REBUILD",
						"RESCUE",
						"SHUTOFF",
						"SUSPENDED",
					),
				},
			},
			"tag_key_value": schema.StringAttribute{
				Description: "Optional. Filter by tag key-value pairs.",
				Optional:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "Filter the server list result by the UUID of the server. Allowed UUID part",
				Optional:    true,
			},
			"tag_value": schema.ListAttribute{
				Description: "Optional. Filter by tag values. ?`tag_value`=value1&`tag_value`=value2",
				Optional:    true,
				ElementType: types.StringType,
			},
			"include_k8s": schema.BoolAttribute{
				Description: "Include managed k8s worker nodes",
				Computed:    true,
				Optional:    true,
			},
			"only_isolated": schema.BoolAttribute{
				Description: "Include only isolated instances",
				Computed:    true,
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
			"with_interfaces_name": schema.BoolAttribute{
				Description: "Include `interface_name` in the addresses",
				Computed:    true,
				Optional:    true,
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
				CustomType:  customfield.NewNestedObjectListType[CloudBaremetalServersItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Bare metal server ID",
							Computed:    true,
						},
						"addresses": schema.MapAttribute{
							Description: "Map of `network_name` to list of addresses in that network",
							Computed:    true,
							CustomType:  customfield.NewMapType[customfield.NestedObjectList[CloudBaremetalServersAddressesDataSourceModel]](ctx),
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
						"blackhole_ports": schema.ListNestedAttribute{
							Description: "IP addresses of the instances that are blackholed by DDoS mitigation system",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudBaremetalServersBlackholePortsDataSourceModel](ctx),
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
						"created_at": schema.StringAttribute{
							Description: "Datetime when bare metal server was created",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"creator_task_id": schema.StringAttribute{
							Description: "Task that created this entity",
							Computed:    true,
						},
						"fixed_ip_assignments": schema.ListNestedAttribute{
							Description: "Fixed IP assigned to instance",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudBaremetalServersFixedIPAssignmentsDataSourceModel](ctx),
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
							Description: "Flavor details",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudBaremetalServersFlavorDataSourceModel](ctx),
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
									CustomType:  customfield.NewNestedObjectType[CloudBaremetalServersFlavorHardwareDescriptionDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
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
						"instance_isolation": schema.SingleNestedAttribute{
							Description: "Instance isolation information",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudBaremetalServersInstanceIsolationDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"reason": schema.StringAttribute{
									Description: "The reason of instance isolation if it is isolated from external internet.",
									Computed:    true,
								},
							},
						},
						"name": schema.StringAttribute{
							Description: "Bare metal server name",
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
						"ssh_key_name": schema.StringAttribute{
							Description: "SSH key assigned to bare metal server",
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
						"tags": schema.ListNestedAttribute{
							Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudBaremetalServersTagsDataSourceModel](ctx),
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
					},
				},
			},
		},
	}
}

func (d *CloudBaremetalServersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudBaremetalServersDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster_interface

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudGPUVirtualClusterInterfacesDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Description: "Cluster unique identifier",
				Required:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
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
				CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClusterInterfacesItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"floating_ips": schema.ListNestedAttribute{
							Description: "Bodies of floatingips that are NAT-ing ips of this port",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClusterInterfacesFloatingIPsDataSourceModel](ctx),
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
									"router_id": schema.StringAttribute{
										Description: "Router ID",
										Computed:    true,
									},
									"status": schema.StringAttribute{
										Description: "Floating IP status\nAvailable values: \"ACTIVE\", \"DOWN\", \"ERROR\".",
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
										CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClusterInterfacesFloatingIPsTagsDataSourceModel](ctx),
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
									"updated_at": schema.StringAttribute{
										Description: "Datetime when the floating IP was last updated",
										Computed:    true,
										CustomType:  timetypes.RFC3339Type{},
									},
								},
							},
						},
						"ip_assignments": schema.ListNestedAttribute{
							Description: "IP addresses assigned to this port",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClusterInterfacesIPAssignmentsDataSourceModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip_address": schema.StringAttribute{
										Description: "The IP address assigned to the port from the specified subnet",
										Computed:    true,
									},
									"subnet_id": schema.StringAttribute{
										Description: "ID of the subnet that allocated the IP",
										Computed:    true,
									},
								},
							},
						},
						"mac_address": schema.StringAttribute{
							Description: "MAC address of the virtual port",
							Computed:    true,
						},
						"network": schema.SingleNestedAttribute{
							Description: "Body of the network this port is attached to",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudGPUVirtualClusterInterfacesNetworkDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "Network ID",
									Computed:    true,
								},
								"created_at": schema.StringAttribute{
									Description: "Datetime when the network was created",
									Computed:    true,
									CustomType:  timetypes.RFC3339Type{},
								},
								"external": schema.BoolAttribute{
									Description: "True if the network `router:external` attribute",
									Computed:    true,
								},
								"mtu": schema.Int64Attribute{
									Description: "MTU (maximum transmission unit)",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Network name",
									Computed:    true,
								},
								"port_security_enabled": schema.BoolAttribute{
									Description: "Indicates `port_security_enabled` status of all newly created in the network ports.",
									Computed:    true,
								},
								"segmentation_id": schema.Int64Attribute{
									Description: "Id of network segment",
									Computed:    true,
								},
								"shared": schema.BoolAttribute{
									Description: "True when the network is shared with your project by external owner",
									Computed:    true,
								},
								"subnets": schema.ListNestedAttribute{
									Description: "List of subnetworks",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClusterInterfacesNetworkSubnetsDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description: "Subnet id.",
												Computed:    true,
											},
											"available_ips": schema.Int64Attribute{
												Description: "Number of available ips in subnet",
												Computed:    true,
											},
											"cidr": schema.StringAttribute{
												Description: "CIDR",
												Computed:    true,
											},
											"created_at": schema.StringAttribute{
												Description: "Datetime when the subnet was created",
												Computed:    true,
												CustomType:  timetypes.RFC3339Type{},
											},
											"dns_nameservers": schema.ListAttribute{
												Description: "List IP addresses of a DNS resolver reachable from the network",
												Computed:    true,
												CustomType:  customfield.NewListType[types.String](ctx),
												ElementType: types.StringType,
											},
											"enable_dhcp": schema.BoolAttribute{
												Description: "Indicates whether DHCP is enabled for this subnet. If true, IP addresses will be assigned automatically",
												Computed:    true,
											},
											"gateway_ip": schema.StringAttribute{
												Description: "Default GW IPv4 address, advertised in DHCP routes of this subnet. If null, no gateway is advertised by this subnet.",
												Computed:    true,
											},
											"has_router": schema.BoolAttribute{
												Description:        "Deprecated. Always returns `false`.",
												Computed:           true,
												DeprecationMessage: "This attribute is deprecated.",
											},
											"host_routes": schema.ListNestedAttribute{
												Description: "List of custom static routes to advertise via DHCP.",
												Computed:    true,
												CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClusterInterfacesNetworkSubnetsHostRoutesDataSourceModel](ctx),
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"destination": schema.StringAttribute{
															Description: "CIDR of destination IPv4 or IPv6 subnet.",
															Computed:    true,
														},
														"nexthop": schema.StringAttribute{
															Description: "IPv4 or IPv6 address to forward traffic to if it's destination IP matches 'destination' CIDR.",
															Computed:    true,
														},
													},
												},
											},
											"ip_version": schema.Int64Attribute{
												Description: "IP version used by the subnet (IPv4 or IPv6)\nAvailable values: 4, 6.",
												Computed:    true,
												Validators: []validator.Int64{
													int64validator.OneOf(4, 6),
												},
											},
											"name": schema.StringAttribute{
												Description: "Subnet name",
												Computed:    true,
											},
											"network_id": schema.StringAttribute{
												Description: "Network ID",
												Computed:    true,
											},
											"tags": schema.ListNestedAttribute{
												Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
												Computed:    true,
												CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClusterInterfacesNetworkSubnetsTagsDataSourceModel](ctx),
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
											"total_ips": schema.Int64Attribute{
												Description: "Total number of ips in subnet",
												Computed:    true,
											},
											"updated_at": schema.StringAttribute{
												Description: "Datetime when the subnet was last updated",
												Computed:    true,
												CustomType:  timetypes.RFC3339Type{},
											},
										},
									},
								},
								"tags": schema.ListNestedAttribute{
									Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClusterInterfacesNetworkTagsDataSourceModel](ctx),
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
								"type": schema.StringAttribute{
									Description: "Network type (vlan, vxlan)",
									Computed:    true,
								},
								"updated_at": schema.StringAttribute{
									Description: "Datetime when the network was last updated",
									Computed:    true,
									CustomType:  timetypes.RFC3339Type{},
								},
							},
						},
						"network_id": schema.StringAttribute{
							Description: "ID of the network the port is attached to",
							Computed:    true,
						},
						"port_id": schema.StringAttribute{
							Description: "ID of virtual ethernet port object",
							Computed:    true,
						},
						"port_security_enabled": schema.BoolAttribute{
							Description: "Port security status",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudGPUVirtualClusterInterfacesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudGPUVirtualClusterInterfacesDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

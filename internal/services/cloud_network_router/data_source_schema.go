// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_router

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudNetworkRouterDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Routers interconnect subnets and manage network routing, including external gateway connectivity and static routes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"router_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.Int64Attribute{
				Optional: true,
			},
			"region_id": schema.Int64Attribute{
				Optional: true,
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the router was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"distributed": schema.BoolAttribute{
				Description: "Whether the router is distributed or centralized.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Router name",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the router.",
				Computed:    true,
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the router was last updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"external_gateway_info": schema.SingleNestedAttribute{
				Description: "State of this router's external gateway.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudNetworkRouterExternalGatewayInfoDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"enable_snat": schema.BoolAttribute{
						Description: "Is SNAT enabled.",
						Computed:    true,
					},
					"external_fixed_ips": schema.ListNestedAttribute{
						Description: "List of external IPs that emit SNAT-ed traffic.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectListType[CloudNetworkRouterExternalGatewayInfoExternalFixedIPsDataSourceModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"ip_address": schema.StringAttribute{
									Description: "IP address",
									Computed:    true,
								},
								"subnet_id": schema.StringAttribute{
									Description: "ID of the subnet that allocated the IP",
									Computed:    true,
								},
							},
						},
					},
					"network_id": schema.StringAttribute{
						Description: "Id of the external network.",
						Computed:    true,
					},
				},
			},
			"interfaces": schema.ListNestedAttribute{
				Description: "List of router interfaces.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudNetworkRouterInterfacesDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip_assignments": schema.ListNestedAttribute{
							Description: "IP addresses assigned to this port",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudNetworkRouterInterfacesIPAssignmentsDataSourceModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip_address": schema.StringAttribute{
										Description: "IP address",
										Computed:    true,
									},
									"subnet_id": schema.StringAttribute{
										Description: "ID of the subnet that allocated the IP",
										Computed:    true,
									},
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
						"mac_address": schema.StringAttribute{
							Description: "MAC address of the virtual port",
							Computed:    true,
						},
					},
				},
			},
			"routes": schema.ListNestedAttribute{
				Description: "List of custom routes.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudNetworkRouterRoutesDataSourceModel](ctx),
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
		},
	}
}

func (d *CloudNetworkRouterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudNetworkRouterDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

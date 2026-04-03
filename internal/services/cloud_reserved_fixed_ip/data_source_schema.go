// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_reserved_fixed_ip

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudReservedFixedIPDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Reserved fixed IPs are static IP addresses that persist independently of instances and can be used as virtual IPs (VIPs) for high availability.",
		Attributes: map[string]schema.Attribute{
			"port_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.Int64Attribute{
				Optional: true,
			},
			"region_id": schema.Int64Attribute{
				Optional: true,
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the reserved fixed IP was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"fixed_ip_address": schema.StringAttribute{
				Description: "IPv4 address of the reserved fixed IP",
				Computed:    true,
			},
			"fixed_ipv6_address": schema.StringAttribute{
				Description: "IPv6 address of the reserved fixed IP",
				Computed:    true,
			},
			"is_external": schema.BoolAttribute{
				Description: "If reserved fixed IP belongs to a public network",
				Computed:    true,
			},
			"is_vip": schema.BoolAttribute{
				Description: "If reserved fixed IP is a VIP",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Reserved fixed IP name",
				Computed:    true,
			},
			"network_id": schema.StringAttribute{
				Description: "ID of the network the port is attached to",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Underlying port status",
				Computed:    true,
			},
			"subnet_id": schema.StringAttribute{
				Description: "ID of the subnet that owns the IP address",
				Computed:    true,
			},
			"subnet_v6_id": schema.StringAttribute{
				Description: "ID of the subnet that owns the IPv6 address",
				Computed:    true,
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the reserved fixed IP was last updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"allowed_address_pairs": schema.ListNestedAttribute{
				Description: "Group of subnet masks and/or IP addresses that share the current IP as VIP",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudReservedFixedIPAllowedAddressPairsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip_address": schema.StringAttribute{
							Description: "Subnet mask or IP address of the port specified in `allowed_address_pairs`",
							Computed:    true,
						},
						"mac_address": schema.StringAttribute{
							Description: "MAC address of the port specified in `allowed_address_pairs`",
							Computed:    true,
						},
					},
				},
			},
			"attachments": schema.ListNestedAttribute{
				Description: "Reserved fixed IP attachment entities",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudReservedFixedIPAttachmentsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_id": schema.StringAttribute{
							Description: "Resource ID",
							Computed:    true,
						},
						"resource_type": schema.StringAttribute{
							Description: "Resource type",
							Computed:    true,
						},
					},
				},
			},
			"network": schema.SingleNestedAttribute{
				Description: "Network details",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudReservedFixedIPNetworkDataSourceModel](ctx),
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
					"creator_task_id": schema.StringAttribute{
						Description: "Task that created this entity",
						Computed:    true,
					},
					"default": schema.BoolAttribute{
						Description: "True if network has `is_default` attribute",
						Computed:    true,
					},
					"external": schema.BoolAttribute{
						Description: "True if the network `router:external` attribute",
						Computed:    true,
					},
					"mtu": schema.Int64Attribute{
						Description: "MTU (maximum transmission unit). Default value is 1450",
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
					"segmentation_id": schema.Int64Attribute{
						Description: "Id of network segment",
						Computed:    true,
					},
					"shared": schema.BoolAttribute{
						Description: "True when the network is shared with your project by external owner",
						Computed:    true,
					},
					"subnets": schema.ListAttribute{
						Description: "List of subnetworks",
						Computed:    true,
						CustomType:  customfield.NewListType[types.String](ctx),
						ElementType: types.StringType,
					},
					"tags": schema.ListNestedAttribute{
						Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectListType[CloudReservedFixedIPNetworkTagsDataSourceModel](ctx),
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
			"reservation": schema.SingleNestedAttribute{
				Description: "Reserved fixed IP status with resource type and ID it is attached to",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudReservedFixedIPReservationDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"resource_id": schema.StringAttribute{
						Description: "ID of the instance or load balancer the IP is attached to",
						Computed:    true,
					},
					"resource_type": schema.StringAttribute{
						Description: "Resource type of the resource the IP is attached to",
						Computed:    true,
					},
					"status": schema.StringAttribute{
						Description: "IP reservation status",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *CloudReservedFixedIPDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudReservedFixedIPDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_subnet

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudNetworkSubnetDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Required:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Required:    true,
			},
			"subnet_id": schema.StringAttribute{
				Description: "Subnet ID",
				Required:    true,
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
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"enable_dhcp": schema.BoolAttribute{
				Description: "True if DHCP should be enabled",
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
			"id": schema.StringAttribute{
				Description: "Subnet id.",
				Computed:    true,
			},
			"ip_version": schema.Int64Attribute{
				Description: "IP version\nAvailable values: 4, 6.",
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
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
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
			"dns_nameservers": schema.ListAttribute{
				Description: "List IP addresses of a DNS resolver reachable from the network",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"host_routes": schema.ListNestedAttribute{
				Description: "List of custom static routes to advertise via DHCP.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudNetworkSubnetHostRoutesDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"destination": schema.StringAttribute{
							Description: "CIDR of destination IPv4 subnet.",
							Computed:    true,
						},
						"nexthop": schema.StringAttribute{
							Description: "IPv4 address to forward traffic to if it's destination IP matches 'destination' CIDR.",
							Computed:    true,
						},
					},
				},
			},
			"tags": schema.ListNestedAttribute{
				Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudNetworkSubnetTagsDataSourceModel](ctx),
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

func (d *CloudNetworkSubnetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudNetworkSubnetDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

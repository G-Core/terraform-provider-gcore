// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_subnet

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudNetworkSubnetResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Subnets define IP address ranges within a network for instance connectivity, with support for DHCP and DNS configuration.",
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
			"cidr": schema.StringAttribute{
				Description:   "CIDR",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"network_id": schema.StringAttribute{
				Description:   "Network ID",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"router_id_to_connect": schema.StringAttribute{
				Description:   "ID of the router to connect to. Requires `connect_to_network_router` set to true. If not specified, attempts to find a router created during network creation.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"connect_to_network_router": schema.BoolAttribute{
				Description:   "True if the network's router should get a gateway in this subnet. Must be explicitly 'false' when `gateway_ip` is null.",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplaceIfConfigured()},
				Default:       booldefault.StaticBool(true),
			},
			"ip_version": schema.Int64Attribute{
				Description: "IP version\nAvailable values: 4, 6.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.OneOf(4, 6),
				},
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
			},
			"name": schema.StringAttribute{
				Description: "Subnet name",
				Required:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"enable_dhcp": schema.BoolAttribute{
				Description: "True if DHCP should be enabled",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"gateway_ip": schema.StringAttribute{
				Description: "Default GW IPv4 address to advertise in DHCP routes in this subnet. Omit this field to let the cloud backend allocate it automatically. Set to null if no gateway must be advertised by this subnet's DHCP (useful when attaching instances to multiple subnets in order to prevent default route conflicts).",
				Computed:    true,
				Optional:    true,
			},
			"dns_nameservers": schema.ListAttribute{
				Description: "List IP addresses of DNS servers to advertise via DHCP.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"host_routes": schema.ListNestedAttribute{
				Description: "List of custom static routes to advertise via DHCP.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudNetworkSubnetHostRoutesModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"destination": schema.StringAttribute{
							Description: "CIDR of destination IPv4 subnet.",
							Required:    true,
						},
						"nexthop": schema.StringAttribute{
							Description: "IPv4 address to forward traffic to if it's destination IP matches 'destination' CIDR.",
							Required:    true,
						},
					},
				},
			},
			"available_ips": schema.Int64Attribute{
				Description: "Number of available ips in subnet",
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
			"has_router": schema.BoolAttribute{
				Description:        "Deprecated. Always returns `false`.",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated.",
				Default:            booldefault.StaticBool(false),
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
			"tasks": schema.ListAttribute{
				Description: "List of task IDs representing asynchronous operations. Use these IDs to monitor operation progress:\n- `GET /v1/tasks/{task_id}` - Check individual task status and details\nPoll task status until completion (`FINISHED`/`ERROR`) before proceeding with dependent operations.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
		},
	}
}

func (r *CloudNetworkSubnetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudNetworkSubnetResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

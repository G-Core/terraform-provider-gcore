// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_member

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a load balancer pool member. **Important:** This resource requires manual task polling since the SDK doesn't support `AddAndPoll` methods yet. Members are created within a specific pool and inherit the pool's load balancing configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Member ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pool_id": schema.StringAttribute{
				Description: "ID of the load balancer pool to which this member belongs",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				Description: "Member IP address",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"protocol_port": schema.Int64Attribute{
				Description: "Member IP port",
				Required:    true,
			},
			"instance_id": schema.StringAttribute{
				Description: "Either `subnet_id` or `instance_id` should be provided",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"monitor_address": schema.StringAttribute{
				Description: "An alternate IP address used for health monitoring of a backend member. Default is null which monitors the member address.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"monitor_port": schema.Int64Attribute{
				Description: "An alternate protocol port used for health monitoring of a backend member. Default is null which monitors the member `protocol_port`.",
				Optional:    true,
			},
			"subnet_id": schema.StringAttribute{
				Description: "`subnet_id` in which `address` is present. Either `subnet_id` or `instance_id` should be provided",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"admin_state_up": schema.BoolAttribute{
				Description: "Administrative state of the resource. When set to true, the resource is enabled and operational. When set to false, the resource is disabled and will not process traffic. When null is passed, the value is skipped and defaults to true.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"backup": schema.BoolAttribute{
				Description: "Set to true if the member is a backup member, to which traffic will be sent exclusively when all non-backup members will be unreachable. It allows to realize ACTIVE-BACKUP load balancing without thinking about VRRP and VIP configuration. Default is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"weight": schema.Int64Attribute{
				Description: "Member weight. Valid values are 0 < `weight` <= 256, defaults to 1. Controls traffic distribution based on the pool's load balancing algorithm:\n* `ROUND_ROBIN`: Distributes connections to each member in turn according to weights. Higher weight = more turns in the cycle. Example: weights 3 vs 1 = ~75% vs ~25% of requests.\n* `LEAST_CONNECTIONS`: Sends new connections to the member with fewest active connections, performing round-robin within groups of the same normalized load. Higher weight = allowed to hold more simultaneous connections before being considered 'more loaded'. Example: weights 2 vs 1 means 20 vs 10 active connections is treated as balanced.\n* `SOURCE_IP`: Routes clients consistently to the same member by hashing client source IP; hash result is modulo total weight of running members. Higher weight = more hash buckets, so more client IPs map to that member. Example: weights 2 vs 1 = roughly two-thirds of distinct client IPs map to the higher-weight member.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
				Validators: []validator.Int64{
					int64validator.Between(1, 256),
				},
			},
			"operating_status": schema.StringAttribute{
				Description: "Operating status of this member",
				Computed:    true,
			},
		},
	}
}

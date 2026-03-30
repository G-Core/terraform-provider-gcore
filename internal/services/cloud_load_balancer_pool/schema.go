// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*CloudLoadBalancerPoolResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Load balancer pools group backend instances with a load balancing algorithm and health monitoring configuration.",
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
			"listener_id": schema.StringAttribute{
				Description:   "Listener ID",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"load_balancer_id": schema.StringAttribute{
				Description:   "Loadbalancer ID",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"lb_algorithm": schema.StringAttribute{
				Description: "Load balancer algorithm\nAvailable values: \"LEAST_CONNECTIONS\", \"ROUND_ROBIN\", \"SOURCE_IP\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"LEAST_CONNECTIONS",
						"ROUND_ROBIN",
						"SOURCE_IP",
					),
				},
			},
			"name": schema.StringAttribute{
				Description: "Pool name",
				Required:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "Protocol\nAvailable values: \"HTTP\", \"HTTPS\", \"PROXY\", \"PROXYV2\", \"TCP\", \"UDP\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"HTTP",
						"HTTPS",
						"PROXY",
						"PROXYV2",
						"TCP",
						"UDP",
					),
				},
			},
			"admin_state_up": schema.BoolAttribute{
				Description: "Administrative state of the resource. When set to true, the resource is enabled and operational. When set to false, the resource is disabled and will not process traffic. Defaults to true.",
				Optional:    true,
			},
			"ca_secret_id": schema.StringAttribute{
				Description: "Secret ID of CA certificate bundle",
				Optional:    true,
			},
			"crl_secret_id": schema.StringAttribute{
				Description: "Secret ID of CA revocation list file",
				Optional:    true,
			},
			"secret_id": schema.StringAttribute{
				Description: "Secret ID for TLS client authentication to the member servers",
				Optional:    true,
			},
			"timeout_member_connect": schema.Int64Attribute{
				Description: "Backend member connection timeout in milliseconds",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 86400000),
				},
			},
			"timeout_member_data": schema.Int64Attribute{
				Description: "Backend member inactivity timeout in milliseconds",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 86400000),
				},
			},
			"healthmonitor": schema.SingleNestedAttribute{
				Description: "Health monitor details",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"delay": schema.Int64Attribute{
						Description: "The time, in seconds, between sending probes to members",
						Required:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 2147483647),
						},
					},
					"max_retries": schema.Int64Attribute{
						Description: "Number of successes before the member is switched to ONLINE state",
						Required:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 10),
						},
					},
					"timeout": schema.Int64Attribute{
						Description: "The maximum time to connect. Must be less than the delay value",
						Required:    true,
						Validators: []validator.Int64{
							int64validator.AtMost(2147483),
						},
					},
					"type": schema.StringAttribute{
						Description: "Health monitor type. Once health monitor is created, cannot be changed.\nAvailable values: \"HTTP\", \"HTTPS\", \"K8S\", \"PING\", \"TCP\", \"TLS-HELLO\", \"UDP-CONNECT\".",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"HTTP",
								"HTTPS",
								"K8S",
								"PING",
								"TCP",
								"TLS-HELLO",
								"UDP-CONNECT",
							),
						},
					},
					"admin_state_up": schema.BoolAttribute{
						Description: "Administrative state of the resource. When set to true, the resource is enabled and operational. When set to false, the resource is disabled and will not process traffic. Defaults to true.",
						Optional:    true,
					},
					"domain_name": schema.StringAttribute{
						Description: "Domain name for HTTP host header. Can only be used together with `HTTP` or `HTTPS` health monitor type.",
						Optional:    true,
					},
					"expected_codes": schema.StringAttribute{
						Description: "Expected HTTP response codes. Can be a single code or a range of codes. Can only be used together with `HTTP` or `HTTPS` health monitor type. For example, 200,202,300-302,401,403,404,500-504. If not specified, the default is 200.",
						Optional:    true,
					},
					"http_method": schema.StringAttribute{
						Description: "HTTP method. Can only be used together with `HTTP` or `HTTPS` health monitor type.\nAvailable values: \"CONNECT\", \"DELETE\", \"GET\", \"HEAD\", \"OPTIONS\", \"PATCH\", \"POST\", \"PUT\", \"TRACE\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"CONNECT",
								"DELETE",
								"GET",
								"HEAD",
								"OPTIONS",
								"PATCH",
								"POST",
								"PUT",
								"TRACE",
							),
						},
					},
					"http_version": schema.StringAttribute{
						Description: "HTTP version. Can only be used together with `HTTP` or `HTTPS` health monitor type. Supported values: 1.0, 1.1.\nAvailable values: \"1.0\", \"1.1\".",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("1.0", "1.1"),
						},
					},
					"max_retries_down": schema.Int64Attribute{
						Description: "Number of failures before the member is switched to ERROR state.",
						Computed:    true,
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 10),
						},
					},
					"url_path": schema.StringAttribute{
						Description: "URL Path. Defaults to '/'. Can only be used together with `HTTP` or `HTTPS` health monitor type.",
						Optional:    true,
					},
				},
			},
			"session_persistence": schema.SingleNestedAttribute{
				Description: "Session persistence details",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Session persistence type\nAvailable values: \"APP_COOKIE\", \"HTTP_COOKIE\", \"SOURCE_IP\".",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"APP_COOKIE",
								"HTTP_COOKIE",
								"SOURCE_IP",
							),
						},
					},
					"cookie_name": schema.StringAttribute{
						Description: "Should be set if app cookie or http cookie is used",
						Optional:    true,
					},
					"persistence_granularity": schema.StringAttribute{
						Description: "Subnet mask if `source_ip` is used. For UDP ports only",
						Optional:    true,
					},
					"persistence_timeout": schema.Int64Attribute{
						Description: "Session persistence timeout. For UDP ports only",
						Optional:    true,
					},
				},
			},
			"members": schema.ListNestedAttribute{
				Description: "Pool members",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerPoolMembersModel](ctx),
				PlanModifiers: []planmodifier.List{
					// Use custom modifier that sets empty list when config is null
					// This ensures omitting `members` from .tf file removes all members
					// See: GCLOUD2-20778
					planmodifiers.UseEmptyListWhenConfigNull(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Description: "Member IP address",
							Required:    true,
						},
						"protocol_port": schema.Int64Attribute{
							Description: "Member IP port",
							Required:    true,
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"admin_state_up": schema.BoolAttribute{
							Description: "Administrative state of the resource. When set to true, the resource is enabled and operational. When set to false, the resource is disabled and will not process traffic. Defaults to true.",
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
						"instance_id": schema.StringAttribute{
							Description: "Either `subnet_id` or `instance_id` should be provided",
							Optional:    true,
						},
						"monitor_address": schema.StringAttribute{
							Description: "An alternate IP address used for health monitoring of a backend member. Default is null which monitors the member address.",
							Optional:    true,
						},
						"monitor_port": schema.Int64Attribute{
							Description: "An alternate protocol port used for health monitoring of a backend member. Default is null which monitors the member `protocol_port`.",
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"subnet_id": schema.StringAttribute{
							Description: "`subnet_id` in which `address` is present. Either `subnet_id` or `instance_id` should be provided",
							Optional:    true,
						},
						"weight": schema.Int64Attribute{
							Description: "Member weight. Valid values are 0 < `weight` <= 256, defaults to 1. Controls traffic distribution based on the pool's load balancing algorithm:\n- `ROUND_ROBIN`: Distributes connections to each member in turn according to weights. Higher weight = more turns in the cycle. Example: weights 3 vs 1 = ~75% vs ~25% of requests.\n- `LEAST_CONNECTIONS`: Sends new connections to the member with fewest active connections, performing round-robin within groups of the same normalized load. Higher weight = allowed to hold more simultaneous connections before being considered 'more loaded'. Example: weights 2 vs 1 means 20 vs 10 active connections is treated as balanced.\n- `SOURCE_IP`: Routes clients consistently to the same member by hashing client source IP; hash result is modulo total weight of running members. Higher weight = more hash buckets, so more client IPs map to that member. Example: weights 2 vs 1 = roughly two-thirds of distinct client IPs map to the higher-weight member.",
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.AtMost(256),
							},
						},
					},
				},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"operating_status": schema.StringAttribute{
				Description: "Pool operating status\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
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
			"provisioning_status": schema.StringAttribute{
				Description: "Pool lifecycle status\nAvailable values: \"ACTIVE\", \"DELETED\", \"ERROR\", \"PENDING_CREATE\", \"PENDING_DELETE\", \"PENDING_UPDATE\".",
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
			"listeners": schema.ListNestedAttribute{
				Description: "Listeners IDs",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerPoolListenersModel](ctx),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Resource ID",
							Computed:    true,
						},
					},
				},
			},
			"loadbalancers": schema.ListNestedAttribute{
				Description: "Load balancers IDs",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerPoolLoadbalancersModel](ctx),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Resource ID",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *CloudLoadBalancerPoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudLoadBalancerPoolResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

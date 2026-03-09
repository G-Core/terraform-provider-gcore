// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudLoadBalancerPoolDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Load balancer pools group backend instances with a load balancing algorithm and health monitoring configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Pool ID",
				Computed:    true,
			},
			"pool_id": schema.StringAttribute{
				Description: "Pool ID",
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
			"admin_state_up": schema.BoolAttribute{
				Description: "Administrative state of the resource. When set to true, the resource is enabled and operational. When set to false, the resource is disabled and will not process traffic. Defaults to true.",
				Computed:    true,
			},
			"ca_secret_id": schema.StringAttribute{
				Description: "Secret ID of CA certificate bundle",
				Computed:    true,
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"crl_secret_id": schema.StringAttribute{
				Description: "Secret ID of CA revocation list file",
				Computed:    true,
			},
			"lb_algorithm": schema.StringAttribute{
				Description: "Load balancer algorithm\nAvailable values: \"LEAST_CONNECTIONS\", \"ROUND_ROBIN\", \"SOURCE_IP\".",
				Computed:    true,
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
				Computed:    true,
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
			"protocol": schema.StringAttribute{
				Description: "Protocol\nAvailable values: \"HTTP\", \"HTTPS\", \"PROXY\", \"PROXYV2\", \"TCP\", \"UDP\".",
				Computed:    true,
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
			"secret_id": schema.StringAttribute{
				Description: "Secret ID for TLS client authentication to the member servers",
				Computed:    true,
			},
			"timeout_client_data": schema.Int64Attribute{
				Description:        "Frontend client inactivity timeout in milliseconds",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated.",
				Validators: []validator.Int64{
					int64validator.Between(0, 86400000),
				},
			},
			"timeout_member_connect": schema.Int64Attribute{
				Description: "Backend member connection timeout in milliseconds",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 86400000),
				},
			},
			"timeout_member_data": schema.Int64Attribute{
				Description: "Backend member inactivity timeout in milliseconds",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 86400000),
				},
			},
			"healthmonitor": schema.SingleNestedAttribute{
				Description: "Health monitor parameters",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerPoolHealthmonitorDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "Health monitor ID",
						Computed:    true,
					},
					"admin_state_up": schema.BoolAttribute{
						Description: "Administrative state of the resource. When set to true, the resource is enabled and operational. When set to false, the resource is disabled and will not process traffic. Defaults to true.",
						Computed:    true,
					},
					"delay": schema.Int64Attribute{
						Description: "The time, in seconds, between sending probes to members",
						Computed:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 2147483647),
						},
					},
					"max_retries": schema.Int64Attribute{
						Description: "Number of successes before the member is switched to ONLINE state",
						Computed:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 10),
						},
					},
					"max_retries_down": schema.Int64Attribute{
						Description: "Number of failures before the member is switched to ERROR state",
						Computed:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 10),
						},
					},
					"operating_status": schema.StringAttribute{
						Description: "Health Monitor operating status\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
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
						Description: "Health monitor lifecycle status\nAvailable values: \"ACTIVE\", \"DELETED\", \"ERROR\", \"PENDING_CREATE\", \"PENDING_DELETE\", \"PENDING_UPDATE\".",
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
					"timeout": schema.Int64Attribute{
						Description: "The maximum time to connect. Must be less than the delay value",
						Computed:    true,
						Validators: []validator.Int64{
							int64validator.AtMost(2147483),
						},
					},
					"type": schema.StringAttribute{
						Description: "Health monitor type. Once health monitor is created, cannot be changed.\nAvailable values: \"HTTP\", \"HTTPS\", \"K8S\", \"PING\", \"TCP\", \"TLS-HELLO\", \"UDP-CONNECT\".",
						Computed:    true,
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
					"expected_codes": schema.StringAttribute{
						Description: "Expected HTTP response codes. Can be a single code or a range of codes. Can only be used together with `HTTP` or `HTTPS` health monitor type. For example, 200,202,300-302,401,403,404,500-504. If not specified, the default is 200.",
						Computed:    true,
					},
					"http_method": schema.StringAttribute{
						Description: "HTTP method\nAvailable values: \"CONNECT\", \"DELETE\", \"GET\", \"HEAD\", \"OPTIONS\", \"PATCH\", \"POST\", \"PUT\", \"TRACE\".",
						Computed:    true,
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
					"url_path": schema.StringAttribute{
						Description: "URL Path. Defaults to '/'",
						Computed:    true,
					},
				},
			},
			"listeners": schema.ListNestedAttribute{
				Description: "Listeners IDs",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerPoolListenersDataSourceModel](ctx),
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
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerPoolLoadbalancersDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Resource ID",
							Computed:    true,
						},
					},
				},
			},
			"members": schema.ListNestedAttribute{
				Description: "Pool members",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerPoolMembersDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Member ID must be provided if an existing member is being updated",
							Computed:    true,
						},
						"address": schema.StringAttribute{
							Description: "Member IP address",
							Computed:    true,
						},
						"admin_state_up": schema.BoolAttribute{
							Description: "Administrative state of the resource. When set to true, the resource is enabled and operational. When set to false, the resource is disabled and will not process traffic. Defaults to true.",
							Computed:    true,
						},
						"backup": schema.BoolAttribute{
							Description: "Set to true if the member is a backup member, to which traffic will be sent exclusively when all non-backup members will be unreachable. It allows to realize ACTIVE-BACKUP load balancing without thinking about VRRP and VIP configuration. Default is false",
							Computed:    true,
						},
						"operating_status": schema.StringAttribute{
							Description: "Member operating status of the entity\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
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
						"protocol_port": schema.Int64Attribute{
							Description: "Member IP port",
							Computed:    true,
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"provisioning_status": schema.StringAttribute{
							Description: "Pool member lifecycle status\nAvailable values: \"ACTIVE\", \"DELETED\", \"ERROR\", \"PENDING_CREATE\", \"PENDING_DELETE\", \"PENDING_UPDATE\".",
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
						"subnet_id": schema.StringAttribute{
							Description: "`subnet_id` in which `address` is present.",
							Computed:    true,
						},
						"weight": schema.Int64Attribute{
							Description: "Member weight. Valid values are 0 < `weight` <= 256, defaults to 1. Controls traffic distribution based on the pool's load balancing algorithm:\n- `ROUND_ROBIN`: Distributes connections to each member in turn according to weights. Higher weight = more turns in the cycle. Example: weights 3 vs 1 = ~75% vs ~25% of requests.\n- `LEAST_CONNECTIONS`: Sends new connections to the member with fewest active connections, performing round-robin within groups of the same normalized load. Higher weight = allowed to hold more simultaneous connections before being considered 'more loaded'. Example: weights 2 vs 1 means 20 vs 10 active connections is treated as balanced.\n- `SOURCE_IP`: Routes clients consistently to the same member by hashing client source IP; hash result is modulo total weight of running members. Higher weight = more hash buckets, so more client IPs map to that member. Example: weights 2 vs 1 = roughly two-thirds of distinct client IPs map to the higher-weight member.",
							Computed:    true,
							Validators: []validator.Int64{
								int64validator.AtMost(256),
							},
						},
						"monitor_address": schema.StringAttribute{
							Description: "An alternate IP address used for health monitoring of a backend member. Default is null which monitors the member address.",
							Computed:    true,
						},
						"monitor_port": schema.Int64Attribute{
							Description: "An alternate protocol port used for health monitoring of a backend member. Default is null which monitors the member `protocol_port`.",
							Computed:    true,
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
					},
				},
			},
			"session_persistence": schema.SingleNestedAttribute{
				Description: "Session persistence parameters",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerPoolSessionPersistenceDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Session persistence type\nAvailable values: \"APP_COOKIE\", \"HTTP_COOKIE\", \"SOURCE_IP\".",
						Computed:    true,
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
						Computed:    true,
					},
					"persistence_granularity": schema.StringAttribute{
						Description: "Subnet mask if `source_ip` is used. For UDP ports only",
						Computed:    true,
					},
					"persistence_timeout": schema.Int64Attribute{
						Description: "Session persistence timeout. For UDP ports only",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *CloudLoadBalancerPoolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudLoadBalancerPoolDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*CloudLoadBalancerResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"region_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"flavor": schema.StringAttribute{
				Description:   "Load balancer flavor name",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "Load balancer name",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name_template": schema.StringAttribute{
				Description:   "Load balancer name which will be changed by template.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"vip_network_id": schema.StringAttribute{
				Description:   "Network ID for load balancer. If not specified, default external network will be used. Mutually exclusive with `vip_port_id`",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"vip_subnet_id": schema.StringAttribute{
				Description:   "Subnet ID for load balancer. If not specified, any subnet from `vip_network_id` will be selected. Ignored when `vip_network_id` is not specified.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"tags": schema.MapAttribute{
				Description:   "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Optional:      true,
				ElementType:   types.StringType,
				PlanModifiers: []planmodifier.Map{mapplanmodifier.RequiresReplace()},
			},
			"floating_ip": schema.SingleNestedAttribute{
				Description: "Floating IP configuration for assignment",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"source": schema.StringAttribute{
						Description: "A new floating IP will be created and attached to the instance. A floating IP is a public IP that makes the instance accessible from the internet, even if it only has a private IP. It works like SNAT, allowing outgoing and incoming traffic.\nAvailable values: \"new\", \"existing\".",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("new", "existing"),
						},
					},
					"existing_floating_id": schema.StringAttribute{
						Description: "An existing available floating IP id must be specified if the source is set to `existing`",
						Optional:    true,
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
			},
			"preferred_connectivity": schema.StringAttribute{
				Description: "Preferred option to establish connectivity between load balancer and its pools members. L2 provides best performance, L3 provides less IPs usage. It is taking effect only if `instance_id` + `ip_address` is provided, not `subnet_id` + `ip_address`, because we're considering this as intentional `subnet_id` specification.\nAvailable values: \"L2\", \"L3\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("L2", "L3"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"vip_ip_family": schema.StringAttribute{
				Description: "IP family for load balancer subnet auto-selection if `vip_network_id` is specified\nAvailable values: \"dual\", \"ipv4\", \"ipv6\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"dual",
						"ipv4",
						"ipv6",
					),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"vip_port_id": schema.StringAttribute{
				Description:   "Existing Reserved Fixed IP port ID for load balancer. Mutually exclusive with `vip_network_id`",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"listeners": schema.ListNestedAttribute{
				Description: "Load balancer listeners. Maximum 50 per LB (excluding Prometheus endpoint listener).",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerListenersModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Load balancer listener name",
							Required:    true,
						},
						"protocol": schema.StringAttribute{
							Description: "Load balancer listener protocol\nAvailable values: \"HTTP\", \"HTTPS\", \"PROMETHEUS\", \"TCP\", \"TERMINATED_HTTPS\", \"UDP\".",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"HTTP",
									"HTTPS",
									"PROMETHEUS",
									"TCP",
									"TERMINATED_HTTPS",
									"UDP",
								),
							},
						},
						"protocol_port": schema.Int64Attribute{
							Description: "Protocol port",
							Required:    true,
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"allowed_cidrs": schema.ListAttribute{
							Description: "Network CIDRs from which service will be accessible",
							Optional:    true,
							ElementType: types.StringType,
						},
						"connection_limit": schema.Int64Attribute{
							Description: "Limit of the simultaneous connections",
							Computed:    true,
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.Between(-1, 1000000),
							},
							Default: int64default.StaticInt64(100000),
						},
						"insert_x_forwarded": schema.BoolAttribute{
							Description: "Add headers X-Forwarded-For, X-Forwarded-Port, X-Forwarded-Proto to requests. Only used with HTTP or `TERMINATED_HTTPS` protocols.",
							Optional:    true,
						},
						"pools": schema.ListNestedAttribute{
							Description: "Member pools",
							Computed:    true,
							Optional:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerListenersPoolsModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
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
									"ca_secret_id": schema.StringAttribute{
										Description: "Secret ID of CA certificate bundle",
										Optional:    true,
									},
									"crl_secret_id": schema.StringAttribute{
										Description: "Secret ID of CA revocation list file",
										Optional:    true,
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
											"expected_codes": schema.StringAttribute{
												Description: "Expected HTTP response codes. Can be a single code or a range of codes. Can only be used together with `HTTP` or `HTTPS` health monitor type. For example, 200,202,300-302,401,403,404,500-504. If not specified, the default is 200.",
												Optional:    true,
											},
											"http_method": schema.StringAttribute{
												Description: "HTTP method. Can only be used together with `HTTP` or `HTTPS` health monitor type.\nAvailable values: \"CONNECT\", \"DELETE\", \"GET\", \"HEAD\", \"OPTIONS\", \"PATCH\", \"POST\", \"PUT\", \"TRACE\".",
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
											"max_retries_down": schema.Int64Attribute{
												Description: "Number of failures before the member is switched to ERROR state.",
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
									"listener_id": schema.StringAttribute{
										Description: "Listener ID",
										Optional:    true,
									},
									"load_balancer_id": schema.StringAttribute{
										Description: "Loadbalancer ID",
										Optional:    true,
									},
									"members": schema.ListNestedAttribute{
										Description: "Pool members",
										Computed:    true,
										Optional:    true,
										CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerListenersPoolsMembersModel](ctx),
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"address": schema.StringAttribute{
													Description: "Member IP address",
													Required:    true,
												},
												"protocol_port": schema.Int64Attribute{
													Description: "Member IP port",
													Required:    true,
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
												},
												"subnet_id": schema.StringAttribute{
													Description: "`subnet_id` in which `address` is present. Either `subnet_id` or `instance_id` should be provided",
													Optional:    true,
												},
												"weight": schema.Int64Attribute{
													Description: "Member weight. Valid values are 0 < `weight` <= 256, defaults to 1. Controls traffic distribution based on the pool's load balancing algorithm:\n\\* `ROUND_ROBIN`: Distributes connections to each member in turn according to weights. Higher weight = more turns in the cycle. Example: weights 3 vs 1 = ~75% vs ~25% of requests.\n\\* `LEAST_CONNECTIONS`: Sends new connections to the member with fewest active connections, performing round-robin within groups of the same normalized load. Higher weight = allowed to hold more simultaneous connections before being considered 'more loaded'. Example: weights 2 vs 1 means 20 vs 10 active connections is treated as balanced.\n\\* `SOURCE_IP`: Routes clients consistently to the same member by hashing client source IP; hash result is modulo total weight of running members. Higher weight = more hash buckets, so more client IPs map to that member. Example: weights 2 vs 1 = roughly two-thirds of distinct client IPs map to the higher-weight member.",
													Optional:    true,
													Validators: []validator.Int64{
														int64validator.AtMost(256),
													},
												},
											},
										},
									},
									"secret_id": schema.StringAttribute{
										Description: "Secret ID for TLS client authentication to the member servers",
										Optional:    true,
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
									"timeout_client_data": schema.Int64Attribute{
										Description: "Frontend client inactivity timeout in milliseconds",
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.Between(0, 86400000),
										},
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
								},
							},
						},
						"secret_id": schema.StringAttribute{
							Description: "ID of the secret where PKCS12 file is stored for `TERMINATED_HTTPS` or PROMETHEUS listener",
							Optional:    true,
						},
						"sni_secret_id": schema.ListAttribute{
							Description: "List of secrets IDs containing PKCS12 format certificate/key bundles for `TERMINATED_HTTPS` or PROMETHEUS listeners",
							Optional:    true,
							ElementType: types.StringType,
						},
						"timeout_client_data": schema.Int64Attribute{
							Description: "Frontend client inactivity timeout in milliseconds",
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.Between(0, 86400000),
							},
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
						"user_list": schema.ListNestedAttribute{
							Description: "Load balancer listener list of username and encrypted password items",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"encrypted_password": schema.StringAttribute{
										Description: "Encrypted password to auth via Basic Authentication",
										Required:    true,
									},
									"username": schema.StringAttribute{
										Description: "Username to auth via Basic Authentication",
										Required:    true,
									},
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplaceIfConfigured()},
			},
			"logging": schema.SingleNestedAttribute{
				Description: "Logging configuration",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerLoggingModel](ctx),
				Attributes: map[string]schema.Attribute{
					"destination_region_id": schema.Int64Attribute{
						Description: "Destination region id to which the logs will be written",
						Optional:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "Enable/disable forwarding logs to LaaS",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"retention_policy": schema.SingleNestedAttribute{
						Description: "The logs retention policy",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"period": schema.Int64Attribute{
								Description: "Duration of days for which logs must be kept.",
								Required:    true,
								Validators: []validator.Int64{
									int64validator.AtMost(1825),
								},
							},
						},
					},
					"topic_name": schema.StringAttribute{
						Description: "The topic name to which the logs will be written",
						Optional:    true,
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplaceIfConfigured()},
			},
			"tasks": schema.ListAttribute{
				Description: "List of task IDs representing asynchronous operations. Use these IDs to monitor operation progress:\n\\* `GET /v1/tasks/{`task_id`}` - Check individual task status and details\nPoll task status until completion (`FINISHED`/`ERROR`) before proceeding with dependent operations.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
		},
	}
}

func (r *CloudLoadBalancerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudLoadBalancerResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

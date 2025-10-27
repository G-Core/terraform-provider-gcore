// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
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
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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
				Description: "Load balancer flavor name",
				Optional:    true,
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
			"name": schema.StringAttribute{
				Description: "Load balancer name",
				Optional:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"preferred_connectivity": schema.StringAttribute{
				Description: "Preferred option to establish connectivity between load balancer and its pools members. L2 provides best performance, L3 provides less IPs usage. It is taking effect only if `instance_id` + `ip_address` is provided, not `subnet_id` + `ip_address`, because we're considering this as intentional `subnet_id` specification.\nAvailable values: \"L2\", \"L3\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("L2", "L3"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the load balancer was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"operating_status": schema.StringAttribute{
				Description: "Load balancer operating status\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
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
				Description: "Load balancer lifecycle status\nAvailable values: \"ACTIVE\", \"DELETED\", \"ERROR\", \"PENDING_CREATE\", \"PENDING_DELETE\", \"PENDING_UPDATE\".",
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
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the load balancer was last updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"vip_address": schema.StringAttribute{
				Description: "Load balancer IP address",
				Computed:    true,
			},
			"tasks": schema.ListAttribute{
				Description: "List of task IDs representing asynchronous operations. Use these IDs to monitor operation progress:\n\\* `GET /v1/tasks/{`task_id`}` - Check individual task status and details\nPoll task status until completion (`FINISHED`/`ERROR`) before proceeding with dependent operations.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"additional_vips": schema.ListNestedAttribute{
				Description: "List of additional IP addresses",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerAdditionalVipsModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip_address": schema.StringAttribute{
							Description: "IP address",
							Computed:    true,
						},
						"subnet_id": schema.StringAttribute{
							Description: "Subnet UUID",
							Computed:    true,
						},
					},
				},
			},
			"ddos_profile": schema.SingleNestedAttribute{
				Description: "Loadbalancer advanced DDoS protection profile.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerDDOSProfileModel](ctx),
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "Unique identifier for the DDoS protection profile",
						Computed:    true,
					},
					"fields": schema.ListNestedAttribute{
						Description: "List of configured field values for the protection profile",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerDDOSProfileFieldsModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Description: "Unique identifier for the DDoS protection field",
									Computed:    true,
								},
								"base_field": schema.Int64Attribute{
									Description: "ID of DDoS profile field",
									Computed:    true,
								},
								"default": schema.StringAttribute{
									Description: "Predefined default value for the field if not specified",
									Computed:    true,
								},
								"description": schema.StringAttribute{
									Description: "Detailed description explaining the field's purpose and usage guidelines",
									Computed:    true,
								},
								"field_name": schema.StringAttribute{
									Description: "Name of DDoS profile field",
									Computed:    true,
								},
								"field_type": schema.StringAttribute{
									Description: "Data type classification of the field (e.g., string, integer, array)",
									Computed:    true,
								},
								"field_value": schema.StringAttribute{
									Description: "Complex value. Only one of 'value' or '`field_value`' must be specified.",
									Computed:    true,
									CustomType:  jsontypes.NormalizedType{},
								},
								"name": schema.StringAttribute{
									Description: "Human-readable name of the protection field",
									Computed:    true,
								},
								"required": schema.BoolAttribute{
									Description: "Indicates whether this field must be provided when creating a protection profile",
									Computed:    true,
								},
								"validation_schema": schema.StringAttribute{
									Description: "JSON schema defining validation rules and constraints for the field value",
									Computed:    true,
									CustomType:  jsontypes.NormalizedType{},
								},
								"value": schema.StringAttribute{
									Description: "Basic type value. Only one of 'value' or '`field_value`' must be specified.",
									Computed:    true,
								},
							},
						},
					},
					"options": schema.SingleNestedAttribute{
						Description: "Configuration options controlling profile activation and BGP routing",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerDDOSProfileOptionsModel](ctx),
						Attributes: map[string]schema.Attribute{
							"active": schema.BoolAttribute{
								Description: "Controls whether the DDoS protection profile is enabled and actively protecting the resource",
								Computed:    true,
							},
							"bgp": schema.BoolAttribute{
								Description: "Enables Border Gateway Protocol (BGP) routing for DDoS protection traffic",
								Computed:    true,
							},
						},
					},
					"profile_template": schema.SingleNestedAttribute{
						Description: "Complete template configuration data used for this profile",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerDDOSProfileProfileTemplateModel](ctx),
						Attributes: map[string]schema.Attribute{
							"id": schema.Int64Attribute{
								Description: "Unique identifier for the DDoS protection template",
								Computed:    true,
							},
							"description": schema.StringAttribute{
								Description: "Detailed description explaining the template's purpose and use cases",
								Computed:    true,
							},
							"fields": schema.ListNestedAttribute{
								Description: "List of configurable fields that define the template's protection parameters",
								Computed:    true,
								CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerDDOSProfileProfileTemplateFieldsModel](ctx),
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.Int64Attribute{
											Description: "Unique identifier for the DDoS protection field",
											Computed:    true,
										},
										"default": schema.StringAttribute{
											Description: "Predefined default value for the field if not specified",
											Computed:    true,
										},
										"description": schema.StringAttribute{
											Description: "Detailed description explaining the field's purpose and usage guidelines",
											Computed:    true,
										},
										"field_type": schema.StringAttribute{
											Description: "Data type classification of the field (e.g., string, integer, array)",
											Computed:    true,
										},
										"name": schema.StringAttribute{
											Description: "Human-readable name of the protection field",
											Computed:    true,
										},
										"required": schema.BoolAttribute{
											Description: "Indicates whether this field must be provided when creating a protection profile",
											Computed:    true,
										},
										"validation_schema": schema.StringAttribute{
											Description: "JSON schema defining validation rules and constraints for the field value",
											Computed:    true,
											CustomType:  jsontypes.NormalizedType{},
										},
									},
								},
							},
							"name": schema.StringAttribute{
								Description: "Human-readable name of the protection template",
								Computed:    true,
							},
						},
					},
					"profile_template_description": schema.StringAttribute{
						Description: "Detailed description of the protection template used for this profile",
						Computed:    true,
					},
					"protocols": schema.ListNestedAttribute{
						Description: "List of network protocols and ports configured for protection",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerDDOSProfileProtocolsModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"port": schema.StringAttribute{
									Description: "Network port number for which protocols are configured",
									Computed:    true,
								},
								"protocols": schema.ListAttribute{
									Description: "List of network protocols enabled on the specified port",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
							},
						},
					},
					"site": schema.StringAttribute{
						Description: "Geographic site identifier where the protection is deployed",
						Computed:    true,
					},
					"status": schema.SingleNestedAttribute{
						Description: "Current operational status and any error information for the profile",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerDDOSProfileStatusModel](ctx),
						Attributes: map[string]schema.Attribute{
							"error_description": schema.StringAttribute{
								Description: "Detailed error message describing any issues with the profile operation",
								Computed:    true,
							},
							"status": schema.StringAttribute{
								Description: "Current operational status of the DDoS protection profile",
								Computed:    true,
							},
						},
					},
				},
			},
			"floating_ips": schema.ListNestedAttribute{
				Description: "List of assigned floating IPs",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerFloatingIPsModel](ctx),
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
						"creator_task_id": schema.StringAttribute{
							Description: "Task that created this entity",
							Computed:    true,
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
							CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerFloatingIPsTagsModel](ctx),
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
						"task_id": schema.StringAttribute{
							Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Datetime when the floating IP was last updated",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
					},
				},
			},
			"stats": schema.SingleNestedAttribute{
				Description: "Statistics of load balancer.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerStatsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"active_connections": schema.Int64Attribute{
						Description: "Currently active connections",
						Computed:    true,
					},
					"bytes_in": schema.Int64Attribute{
						Description: "Total bytes received",
						Computed:    true,
					},
					"bytes_out": schema.Int64Attribute{
						Description: "Total bytes sent",
						Computed:    true,
					},
					"request_errors": schema.Int64Attribute{
						Description: "Total requests that were unable to be fulfilled",
						Computed:    true,
					},
					"total_connections": schema.Int64Attribute{
						Description: "Total connections handled",
						Computed:    true,
					},
				},
			},
			"tags_v2": schema.ListNestedAttribute{
				Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerTagsV2Model](ctx),
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
			"vrrp_ips": schema.ListNestedAttribute{
				Description: "List of VRRP IP addresses",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerVrrpIPsModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip_address": schema.StringAttribute{
							Description: "IP address",
							Computed:    true,
						},
						"role": schema.StringAttribute{
							Description: "LoadBalancer instance role to which VRRP IP belong\nAvailable values: \"BACKUP\", \"MASTER\", \"STANDALONE\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"BACKUP",
									"MASTER",
									"STANDALONE",
								),
							},
						},
						"subnet_id": schema.StringAttribute{
							Description: "Subnet UUID",
							Computed:    true,
						},
					},
				},
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

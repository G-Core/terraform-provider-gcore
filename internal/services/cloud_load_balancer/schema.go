// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudLoadBalancerResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Load balancers distribute incoming traffic across multiple instances with support for listeners, pools, and health monitoring.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.Int64Attribute{
				Description:   "Project ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
			},
			"region_id": schema.Int64Attribute{
				Description:   "Region ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
			},
			"flavor": schema.StringAttribute{
				Description:   "Load balancer flavor name",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vip_port_id": schema.StringAttribute{
				Description: "Existing Reserved Fixed IP port ID for load balancer. Mutually exclusive with `vip_network_id`",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Load balancer name.",
				Required:    true,
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
			"tags": schema.MapAttribute{
				Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"logging": schema.SingleNestedAttribute{
				Description: "Logging configuration",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerLoggingModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
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
			"admin_state_up": schema.BoolAttribute{
				Description: "Administrative state of the resource. When set to true, the resource is enabled and operational. When set to false, the resource is disabled and will not process traffic. Defaults to true.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the load balancer was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the load balancer was last updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"vip_address": schema.StringAttribute{
				Description: "Load balancer IP address",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vip_fqdn": schema.StringAttribute{
				Description: "Fully qualified domain name for the load balancer VIP",
				Computed:    true,
			},
			"additional_vips": schema.ListNestedAttribute{
				Description: "List of additional IP addresses",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerAdditionalVipsModel](ctx),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
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
			"floating_ips": schema.ListNestedAttribute{
				Description: "List of assigned floating IPs",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerFloatingIPsModel](ctx),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
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
							Description: "Floating IP status. DOWN - unassigned (available). ACTIVE - attached to a port (in use). ERROR - error state.\nAvailable values: \"ACTIVE\", \"DOWN\", \"ERROR\".",
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
			"stats": schema.SingleNestedAttribute{
				Description: "Statistics of load balancer.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerStatsModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
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
				// NOTE: Removed UseStateForUnknown() to fix GCLOUD2-20778 tags inconsistency error
				// tags_v2 is derived from tags input, so it should always refresh from API
				PlanModifiers: []planmodifier.List{},
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

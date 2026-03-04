// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_listener

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudLoadBalancerListenerResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Load balancer listeners handle incoming traffic on specified protocols and ports, forwarding requests to backend pools.",
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
			"load_balancer_id": schema.StringAttribute{
				Description:   "ID of already existent Load Balancer.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
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
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"protocol_port": schema.Int64Attribute{
				Description: "Protocol port",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"default_pool_id": schema.StringAttribute{
				Description:   "ID of already existent Load Balancer Pool to attach listener to.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"insert_x_forwarded": schema.BoolAttribute{
				Description:   "Add headers X-Forwarded-For, X-Forwarded-Port, X-Forwarded-Proto to requests. Only used with HTTP or `TERMINATED_HTTPS` protocols.",
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description: "Load balancer listener name",
				Required:    true,
			},
			"admin_state_up": schema.BoolAttribute{
				Description: "Administrative state of the resource. When set to true, the resource is enabled and operational. When set to false, the resource is disabled and will not process traffic. Defaults to true.",
				Optional:    true,
			},
			"secret_id": schema.StringAttribute{
				Description: "ID of the secret where PKCS12 file is stored for `TERMINATED_HTTPS` or PROMETHEUS listener\nAvailable values: \"\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(""),
				},
			},
			"allowed_cidrs": schema.ListAttribute{
				Description: "Network CIDRs from which service will be accessible",
				Optional:    true,
				ElementType: types.StringType,
			},
			"connection_limit": schema.Int64Attribute{
				Description: "Limit of the simultaneous connections. If -1 is provided, it is translated to the default value 100000.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(-1, 1000000),
				},
				Default: int64default.StaticInt64(100000),
			},
			"timeout_client_data": schema.Int64Attribute{
				Description: "Frontend client inactivity timeout in milliseconds",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 86400000),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"sni_secret_id": schema.ListAttribute{
				Description: "List of secrets IDs containing PKCS12 format certificate/key bundles for `TERMINATED_HTTPS` or PROMETHEUS listeners",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"user_list": schema.ListNestedAttribute{
				Description: "Load balancer listener list of username and encrypted password items",
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
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
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"operating_status": schema.StringAttribute{
				Description: "Listener operating status\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pool_count": schema.Int64Attribute{
				Description: "Number of pools (for UI)",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"provisioning_status": schema.StringAttribute{
				Description: "Listener lifecycle status\nAvailable values: \"ACTIVE\", \"DELETED\", \"ERROR\", \"PENDING_CREATE\", \"PENDING_DELETE\", \"PENDING_UPDATE\".",
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"insert_headers": schema.StringAttribute{
				Description: "Dictionary of additional header insertion into HTTP headers. Only used with HTTP and `TERMINATED_HTTPS` protocols.",
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"stats": schema.SingleNestedAttribute{
				Description: "Statistics of the load balancer. It is available only in get functions by a flag.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerListenerStatsModel](ctx),
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
		},
	}
}

func (r *CloudLoadBalancerListenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudLoadBalancerListenerResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_health_monitor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*CloudLoadBalancerPoolHealthMonitorResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Pool health monitors probe backend members to decide whether they can receive traffic. A pool has at most one health monitor.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseNonNullStateForUnknown()},
			},
			"pool_id": schema.StringAttribute{
				Description:   "Pool ID",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
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
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
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
					int64validator.Between(0, 2147483),
				},
			},
			"domain_name": schema.StringAttribute{
				Description: "Domain name for HTTP host header. Can only be used together with `HTTP` or `HTTPS` health monitor type.",
				Optional:    true,
			},
			"http_version": schema.StringAttribute{
				Description: "HTTP version. Can only be used together with `HTTP` or `HTTPS` health monitor type. Supported values: 1.0, 1.1.\nAvailable values: \"1.0\", \"1.1\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("1.0", "1.1"),
				},
			},
			"admin_state_up": schema.BoolAttribute{
				Description: "Administrative state of the resource. When set to true, the resource is enabled and operational. When set to false, the resource is disabled and will not process traffic. Defaults to true.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"expected_codes": schema.StringAttribute{
				Description: "Expected HTTP response codes. Can be a single code, a comma-separated list of codes, or a single range of codes. Can only be used together with `HTTP` or `HTTPS` health monitor type. For example, 200, 200,202,401,403,404, or 200-204. If not specified, the default is 200.",
				Computed:    true,
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
			"max_retries_down": schema.Int64Attribute{
				Description: "Number of failures before the member is switched to ERROR state.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 10),
				},
				Default: int64default.StaticInt64(3),
			},
			"url_path": schema.StringAttribute{
				Description: "The HTTP path the health monitor requests on each member. Defaults to `/` if not set. Can only be used with `HTTP` or `HTTPS` health monitor type.\n\nMust start with `/` and contain only plain path segments. Query strings (`?`), fragments (`#`), percent-encoding (`%`), and consecutive slashes (`//`) are not allowed.\n\nExamples of valid paths:\n- `/` — check the root (most common, default)\n- `/healthz` — a dedicated health endpoint",
				Computed:    true,
				Optional:    true,
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
		},
	}
}

func (r *CloudLoadBalancerPoolHealthMonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudLoadBalancerPoolHealthMonitorResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

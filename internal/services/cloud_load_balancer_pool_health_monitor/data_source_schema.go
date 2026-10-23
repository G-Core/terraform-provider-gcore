// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_health_monitor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudLoadBalancerPoolHealthMonitorDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Pool health monitors probe backend members to decide whether they can receive traffic. A pool has at most one health monitor.",
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
			"delay": schema.Int64Attribute{
				Description: "The time, in seconds, between sending probes to members",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 2147483647),
				},
			},
			"domain_name": schema.StringAttribute{
				Description: "Domain name for HTTP host header. Can only be used together with `HTTP` or `HTTPS` health monitor type.",
				Computed:    true,
			},
			"expected_codes": schema.StringAttribute{
				Description: "Expected HTTP response codes. Can be a single code, a comma-separated list of codes, or a single range of codes. Can only be used together with `HTTP` or `HTTPS` health monitor type. For example, 200, 200,202,401,403,404, or 200-204. If not specified, the default is 200. Null for non-HTTP(S) health monitor types.",
				Computed:    true,
			},
			"http_method": schema.StringAttribute{
				Description: "HTTP method. Null for non-HTTP(S) health monitor types.\nAvailable values: \"CONNECT\", \"DELETE\", \"GET\", \"HEAD\", \"OPTIONS\", \"PATCH\", \"POST\", \"PUT\", \"TRACE\".",
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
			"http_version": schema.StringAttribute{
				Description: "HTTP version. Can only be used together with `HTTP` or `HTTPS` health monitor type.\nAvailable values: \"1.0\", \"1.1\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("1.0", "1.1"),
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
					int64validator.Between(0, 2147483),
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
			"url_path": schema.StringAttribute{
				Description: "URL Path. Defaults to '/'. Null for non-HTTP(S) health monitor types.",
				Computed:    true,
			},
		},
	}
}

func (d *CloudLoadBalancerPoolHealthMonitorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudLoadBalancerPoolHealthMonitorDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

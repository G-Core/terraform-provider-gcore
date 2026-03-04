// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_listener

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudLoadBalancerListenerDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Listener ID",
				Computed:    true,
			},
			"listener_id": schema.StringAttribute{
				Description: "Listener ID",
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
			"show_stats": schema.BoolAttribute{
				Description: "Show stats",
				Computed:    true,
				Optional:    true,
			},
			"admin_state_up": schema.BoolAttribute{
				Description: "Administrative state of the resource. When set to true, the resource is enabled and operational. When set to false, the resource is disabled and will not process traffic. Defaults to true.",
				Computed:    true,
			},
			"connection_limit": schema.Int64Attribute{
				Description: "Limit of simultaneous connections",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(-1, 1000000),
				},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"load_balancer_id": schema.StringAttribute{
				Description: "Load balancer ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Load balancer listener name",
				Computed:    true,
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
			},
			"pool_count": schema.Int64Attribute{
				Description: "Number of pools (for UI)",
				Computed:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "Load balancer protocol\nAvailable values: \"HTTP\", \"HTTPS\", \"PROMETHEUS\", \"TCP\", \"TERMINATED_HTTPS\", \"UDP\".",
				Computed:    true,
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
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
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
			},
			"secret_id": schema.StringAttribute{
				Description: "ID of the secret where PKCS12 file is stored for `TERMINATED_HTTPS` or PROMETHEUS load balancer",
				Computed:    true,
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"timeout_client_data": schema.Int64Attribute{
				Description: "Frontend client inactivity timeout in milliseconds",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 86400000),
				},
			},
			"timeout_member_connect": schema.Int64Attribute{
				Description:        "Backend member connection timeout in milliseconds",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated.",
				Validators: []validator.Int64{
					int64validator.Between(0, 86400000),
				},
			},
			"timeout_member_data": schema.Int64Attribute{
				Description:        "Backend member inactivity timeout in milliseconds",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated.",
				Validators: []validator.Int64{
					int64validator.Between(0, 86400000),
				},
			},
			"allowed_cidrs": schema.ListAttribute{
				Description: "Network CIDRs from which service will be accessible",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"insert_headers": schema.MapAttribute{
				Description: "Dictionary of additional header insertion into HTTP headers. Only used with HTTP and `TERMINATED_HTTPS` protocols.",
				Computed:    true,
				CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
				ElementType: jsontypes.NormalizedType{},
			},
			"sni_secret_id": schema.ListAttribute{
				Description: "List of secret's ID containing PKCS12 format certificate/key bundles for `TERMINATED_HTTPS` or PROMETHEUS listeners",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"stats": schema.SingleNestedAttribute{
				Description: "Statistics of the load balancer. It is available only in get functions by a flag.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerListenerStatsDataSourceModel](ctx),
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
			"user_list": schema.ListNestedAttribute{
				Description: "Load balancer listener users list",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerListenerUserListDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"encrypted_password": schema.StringAttribute{
							Description: "Encrypted password to auth via Basic Authentication",
							Computed:    true,
						},
						"username": schema.StringAttribute{
							Description: "Username to auth via Basic Authentication",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudLoadBalancerListenerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudLoadBalancerListenerDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

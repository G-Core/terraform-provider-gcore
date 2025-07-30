// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_l7_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudLoadBalancerL7PolicyDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"l7policy_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.Int64Attribute{
				Required: true,
			},
			"region_id": schema.Int64Attribute{
				Required: true,
			},
			"action": schema.StringAttribute{
				Description: "Action\nAvailable values: \"REDIRECT_PREFIX\", \"REDIRECT_TO_POOL\", \"REDIRECT_TO_URL\", \"REJECT\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"REDIRECT_PREFIX",
						"REDIRECT_TO_POOL",
						"REDIRECT_TO_URL",
						"REJECT",
					),
				},
			},
			"id": schema.StringAttribute{
				Description: "ID",
				Computed:    true,
			},
			"listener_id": schema.StringAttribute{
				Description: "Listener ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Human-readable name of the policy",
				Computed:    true,
			},
			"operating_status": schema.StringAttribute{
				Description: "L7 policy operating status\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
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
			"position": schema.Int64Attribute{
				Description: "The position of this policy on the listener. Positions start at 1.",
				Computed:    true,
			},
			"provisioning_status": schema.StringAttribute{
				Description: `Available values: "ACTIVE", "DELETED", "ERROR", "PENDING_CREATE", "PENDING_DELETE", "PENDING_UPDATE".`,
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
			"redirect_http_code": schema.Int64Attribute{
				Description: "Requests matching this policy will be redirected to the specified URL or Prefix URL with the HTTP response code. Valid if action is `REDIRECT_TO_URL` or `REDIRECT_PREFIX`. Valid options are 301, 302, 303, 307, or 308. Default is 302.",
				Computed:    true,
			},
			"redirect_pool_id": schema.StringAttribute{
				Description: "Requests matching this policy will be redirected to the pool with this ID. Only valid if action is `REDIRECT_TO_POOL`.",
				Computed:    true,
			},
			"redirect_prefix": schema.StringAttribute{
				Description: "Requests matching this policy will be redirected to this Prefix URL. Only valid if action is `REDIRECT_PREFIX`.",
				Computed:    true,
			},
			"redirect_url": schema.StringAttribute{
				Description: "Requests matching this policy will be redirected to this URL. Only valid if action is `REDIRECT_TO_URL`.",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				Description: "A list of simple strings assigned to the resource.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"rules": schema.ListNestedAttribute{
				Description: "Rules. All the rules associated with a given policy are logically ANDed together. A request must match all the policy’s rules to match the policy.If you need to express a logical OR operation between rules, then do this by creating multiple policies with the same action.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerL7PolicyRulesDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "L7Rule ID",
							Computed:    true,
						},
						"compare_type": schema.StringAttribute{
							Description: "The comparison type for the L7 rule\nAvailable values: \"CONTAINS\", \"ENDS_WITH\", \"EQUAL_TO\", \"REGEX\", \"STARTS_WITH\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"CONTAINS",
									"ENDS_WITH",
									"EQUAL_TO",
									"REGEX",
									"STARTS_WITH",
								),
							},
						},
						"invert": schema.BoolAttribute{
							Description: "When true the logic of the rule is inverted. For example, with invert true, 'equal to' would become 'not equal to'. Default is false.",
							Computed:    true,
						},
						"key": schema.StringAttribute{
							Description: "The key to use for the comparison. For example, the name of the cookie to evaluate.",
							Computed:    true,
						},
						"operating_status": schema.StringAttribute{
							Description: "L7 policy operating status\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
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
						"project_id": schema.Int64Attribute{
							Description: "Project ID",
							Computed:    true,
						},
						"provisioning_status": schema.StringAttribute{
							Description: `Available values: "ACTIVE", "DELETED", "ERROR", "PENDING_CREATE", "PENDING_DELETE", "PENDING_UPDATE".`,
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
						"region_id": schema.Int64Attribute{
							Description: "Region ID",
							Computed:    true,
						},
						"tags": schema.ListAttribute{
							Description: "A list of simple strings assigned to the l7 rule",
							Computed:    true,
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
						},
						"task_id": schema.StringAttribute{
							Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The L7 rule type\nAvailable values: \"COOKIE\", \"FILE_TYPE\", \"HEADER\", \"HOST_NAME\", \"PATH\", \"SSL_CONN_HAS_CERT\", \"SSL_DN_FIELD\", \"SSL_VERIFY_RESULT\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"COOKIE",
									"FILE_TYPE",
									"HEADER",
									"HOST_NAME",
									"PATH",
									"SSL_CONN_HAS_CERT",
									"SSL_DN_FIELD",
									"SSL_VERIFY_RESULT",
								),
							},
						},
						"value": schema.StringAttribute{
							Description: "The value to use for the comparison. For example, the file type to compare.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudLoadBalancerL7PolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudLoadBalancerL7PolicyDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_l7_policy_rule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudLoadBalancerL7PolicyRuleDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"l7policy_id": schema.StringAttribute{
				Required: true,
			},
			"l7rule_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.Int64Attribute{
				Required: true,
			},
			"region_id": schema.Int64Attribute{
				Required: true,
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
			"id": schema.StringAttribute{
				Description: "L7Rule ID",
				Computed:    true,
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
			"tags": schema.ListAttribute{
				Description: "A list of simple strings assigned to the l7 rule",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
		},
	}
}

func (d *CloudLoadBalancerL7PolicyRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudLoadBalancerL7PolicyRuleDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

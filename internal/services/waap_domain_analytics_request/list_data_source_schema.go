// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_analytics_request

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainAnalyticsRequestsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"start": schema.StringAttribute{
				Description: "Filter traffic starting from a specified date in ISO 8601 format",
				Required:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"end": schema.StringAttribute{
				Description: "Filter traffic up to a specified end date in ISO 8601 format. If not provided, defaults to the current date and time.",
				Optional:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"ip": schema.StringAttribute{
				Description: "Filter the response by IP.",
				Optional:    true,
			},
			"ordering": schema.StringAttribute{
				Description: "Sort the response by given field.",
				Optional:    true,
			},
			"reference_id": schema.StringAttribute{
				Description: "Filter the response by reference ID.",
				Optional:    true,
			},
			"security_rule_name": schema.StringAttribute{
				Description: "Filter the response by security rule name.",
				Optional:    true,
			},
			"status_code": schema.Int64Attribute{
				Description: "Filter the response by response code.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(100, 599),
				},
			},
			"actions": schema.ListAttribute{
				Description: "Filter the response by actions.",
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOfCaseInsensitive(
							"allow",
							"block",
							"captcha",
							"handshake",
						),
					),
				},
				ElementType: types.StringType,
			},
			"countries": schema.ListAttribute{
				Description: "Filter the response by country codes in ISO 3166-1 alpha-2 format.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"traffic_types": schema.ListAttribute{
				Description: "Filter the response by traffic types.",
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOfCaseInsensitive(
							"policy_allowed",
							"policy_blocked",
							"custom_rule_allowed",
							"custom_blocked",
							"legit_requests",
							"sanctioned",
							"dynamic",
							"api",
							"static",
							"ajax",
							"redirects",
							"monitor",
							"err_40x",
							"err_50x",
							"passed_to_origin",
							"timeout",
							"other",
							"ddos",
							"legit",
							"monitored",
						),
					),
				},
				ElementType: types.StringType,
			},
			"limit": schema.Int64Attribute{
				Description: "Number of items to return",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
			},
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[WaapDomainAnalyticsRequestsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Request's unique id",
							Computed:    true,
						},
						"action": schema.StringAttribute{
							Description: "Action of the triggered rule",
							Computed:    true,
						},
						"client_ip": schema.StringAttribute{
							Description: "Client's IP address.",
							Computed:    true,
						},
						"country": schema.StringAttribute{
							Description: "Country code",
							Computed:    true,
						},
						"domain": schema.StringAttribute{
							Description: "Domain name",
							Computed:    true,
						},
						"method": schema.StringAttribute{
							Description: "HTTP method",
							Computed:    true,
						},
						"organization": schema.StringAttribute{
							Description: "Organization",
							Computed:    true,
						},
						"path": schema.StringAttribute{
							Description: "Request path",
							Computed:    true,
						},
						"reference_id": schema.StringAttribute{
							Description: "The reference ID to a sanction that was given to a user.",
							Computed:    true,
						},
						"request_time": schema.Int64Attribute{
							Description: "The UNIX timestamp in ms of the date a set of traffic counters was recorded",
							Computed:    true,
						},
						"result": schema.StringAttribute{
							Description: `Available values: "passed", "blocked", "suppressed", "".`,
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"passed",
									"blocked",
									"suppressed",
									"",
								),
							},
						},
						"rule_id": schema.StringAttribute{
							Description: "The ID of the triggered rule.",
							Computed:    true,
						},
						"rule_name": schema.StringAttribute{
							Description: "Name of the triggered rule",
							Computed:    true,
						},
						"status_code": schema.Int64Attribute{
							Description: "Status code for http request",
							Computed:    true,
						},
						"traffic_types": schema.StringAttribute{
							Description: "Comma separated list of traffic types.",
							Computed:    true,
						},
						"user_agent": schema.StringAttribute{
							Description: "User agent",
							Computed:    true,
						},
						"user_agent_client": schema.StringAttribute{
							Description: "Client from parsed User agent header",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *WaapDomainAnalyticsRequestsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *WaapDomainAnalyticsRequestsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

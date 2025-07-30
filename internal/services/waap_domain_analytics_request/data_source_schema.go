// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_analytics_request

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainAnalyticsRequestDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"request_id": schema.StringAttribute{
				Description: "The request ID",
				Required:    true,
			},
			"action": schema.StringAttribute{
				Description: "Request action",
				Computed:    true,
			},
			"content_type": schema.StringAttribute{
				Description: "Content type of request",
				Computed:    true,
			},
			"domain": schema.StringAttribute{
				Description: "Domain name",
				Computed:    true,
			},
			"http_status_code": schema.Int64Attribute{
				Description: "Status code for http request",
				Computed:    true,
			},
			"http_version": schema.StringAttribute{
				Description: "HTTP version of request",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "Request ID",
				Computed:    true,
			},
			"incident_id": schema.StringAttribute{
				Description: "ID of challenge that was generated",
				Computed:    true,
			},
			"method": schema.StringAttribute{
				Description: "Request method",
				Computed:    true,
			},
			"path": schema.StringAttribute{
				Description: "Request path",
				Computed:    true,
			},
			"query_string": schema.StringAttribute{
				Description: "The query string of the request",
				Computed:    true,
			},
			"reference_id": schema.StringAttribute{
				Description: "Reference ID to identify user sanction",
				Computed:    true,
			},
			"request_time": schema.StringAttribute{
				Description: "The time of the request",
				Computed:    true,
			},
			"request_type": schema.StringAttribute{
				Description: "The type of the request that generated an event",
				Computed:    true,
			},
			"requested_domain": schema.StringAttribute{
				Description: "The real domain name",
				Computed:    true,
			},
			"response_time": schema.StringAttribute{
				Description: "Time took to process all request",
				Computed:    true,
			},
			"result": schema.StringAttribute{
				Description: "The result of a request\nAvailable values: \"passed\", \"blocked\", \"suppressed\", \"\".",
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
				Description: "ID of the triggered rule",
				Computed:    true,
			},
			"rule_name": schema.StringAttribute{
				Description: "Name of the triggered rule",
				Computed:    true,
			},
			"scheme": schema.StringAttribute{
				Description: "The HTTP scheme of the request that generated an event",
				Computed:    true,
			},
			"session_request_count": schema.StringAttribute{
				Description: "The number requests in session",
				Computed:    true,
			},
			"traffic_types": schema.ListAttribute{
				Description: "List of traffic types",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"common_tags": schema.ListNestedAttribute{
				Description: "List of common tags",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[WaapDomainAnalyticsRequestCommonTagsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Description: "Tag description information",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "The tag's display name",
							Computed:    true,
						},
						"tag": schema.StringAttribute{
							Description: "Tag name",
							Computed:    true,
						},
					},
				},
			},
			"network": schema.SingleNestedAttribute{
				Description: "Network details",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[WaapDomainAnalyticsRequestNetworkDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"client_ip": schema.StringAttribute{
						Description: "Client IP",
						Computed:    true,
					},
					"country": schema.StringAttribute{
						Description: "Country code",
						Computed:    true,
					},
					"organization": schema.SingleNestedAttribute{
						Description: "Organization details",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[WaapDomainAnalyticsRequestNetworkOrganizationDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Description: "Organization name",
								Computed:    true,
							},
							"subnet": schema.StringAttribute{
								Description: "Network range",
								Computed:    true,
							},
						},
					},
				},
			},
			"pattern_matched_tags": schema.ListNestedAttribute{
				Description: "List of shield tags",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[WaapDomainAnalyticsRequestPatternMatchedTagsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Description: "Tag description information",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "The tag's display name",
							Computed:    true,
						},
						"execution_phase": schema.StringAttribute{
							Description: "The phase in which the tag was triggered: access -> Request, `header_filter` -> `response_header`, `body_filter` -> `response_body`",
							Computed:    true,
						},
						"field": schema.StringAttribute{
							Description: "The entity to which the variable that triggered the tag belong to. For example: `request_headers`, uri, cookies etc.",
							Computed:    true,
						},
						"field_name": schema.StringAttribute{
							Description: "The name of the variable which holds the value that triggered the tag",
							Computed:    true,
						},
						"pattern_name": schema.StringAttribute{
							Description: "The name of the detected regexp pattern",
							Computed:    true,
						},
						"pattern_value": schema.StringAttribute{
							Description: "The pattern which triggered the tag",
							Computed:    true,
						},
						"tag": schema.StringAttribute{
							Description: "Tag name",
							Computed:    true,
						},
					},
				},
			},
			"user_agent": schema.SingleNestedAttribute{
				Description: "User agent details",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[WaapDomainAnalyticsRequestUserAgentDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"base_browser": schema.StringAttribute{
						Description: "User agent browser",
						Computed:    true,
					},
					"base_browser_version": schema.StringAttribute{
						Description: "User agent browser version",
						Computed:    true,
					},
					"client": schema.StringAttribute{
						Description: "Client from User agent header",
						Computed:    true,
					},
					"client_type": schema.StringAttribute{
						Description: "User agent client type",
						Computed:    true,
					},
					"client_version": schema.StringAttribute{
						Description: "User agent client version",
						Computed:    true,
					},
					"cpu": schema.StringAttribute{
						Description: "User agent cpu",
						Computed:    true,
					},
					"device": schema.StringAttribute{
						Description: "User agent device",
						Computed:    true,
					},
					"device_type": schema.StringAttribute{
						Description: "User agent device type",
						Computed:    true,
					},
					"full_string": schema.StringAttribute{
						Description: "User agent",
						Computed:    true,
					},
					"os": schema.StringAttribute{
						Description: "User agent os",
						Computed:    true,
					},
					"rendering_engine": schema.StringAttribute{
						Description: "User agent engine",
						Computed:    true,
					},
				},
			},
			"request_headers": schema.StringAttribute{
				Description: "HTTP request headers",
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
			},
		},
	}
}

func (d *WaapDomainAnalyticsRequestDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapDomainAnalyticsRequestDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

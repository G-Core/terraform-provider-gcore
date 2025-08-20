// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app_log

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*FastedgeAppLogsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"client_ip": schema.StringAttribute{
				Description: "Search by client IP",
				Optional:    true,
			},
			"edge": schema.StringAttribute{
				Description: "Edge name",
				Optional:    true,
			},
			"from": schema.StringAttribute{
				Description: "Reporting period start time, RFC3339 format. Default 1 hour ago.",
				Optional:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"limit": schema.Int64Attribute{
				Description: "Limit for pagination",
				Optional:    true,
			},
			"search": schema.StringAttribute{
				Description: "Search string",
				Optional:    true,
			},
			"sort": schema.StringAttribute{
				Description: "Sort order (default desc)\nAvailable values: \"desc\", \"asc\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("desc", "asc"),
				},
			},
			"to": schema.StringAttribute{
				Description: "Reporting period end time, RFC3339 format. Default current time in UTC.",
				Optional:    true,
				CustomType:  timetypes.RFC3339Type{},
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
				CustomType:  customfield.NewNestedObjectListType[FastedgeAppLogsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Id of the log",
							Computed:    true,
						},
						"app_name": schema.StringAttribute{
							Description: "Name of the application",
							Computed:    true,
						},
						"client_ip": schema.StringAttribute{
							Description: "Client IP",
							Computed:    true,
						},
						"edge": schema.StringAttribute{
							Description: "Edge name",
							Computed:    true,
						},
						"log": schema.StringAttribute{
							Description: "Log message",
							Computed:    true,
						},
						"timestamp": schema.StringAttribute{
							Description: "Timestamp of a log in RFC3339 format",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
					},
				},
			},
		},
	}
}

func (d *FastedgeAppLogsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *FastedgeAppLogsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_discovery_scan_result

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

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainAPIDiscoveryScanResultsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"message": schema.StringAttribute{
				Description: "Filter by the message of the scan. Supports '\\*' as a wildcard character",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "The different statuses a task result can have\nAvailable values: \"SUCCESS\", \"FAILURE\", \"IN_PROGRESS\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"SUCCESS",
						"FAILURE",
						"IN_PROGRESS",
					),
				},
			},
			"type": schema.StringAttribute{
				Description: "The different types of scans that can be performed\nAvailable values: \"TRAFFIC_SCAN\", \"API_DESCRIPTION_FILE_SCAN\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("TRAFFIC_SCAN", "API_DESCRIPTION_FILE_SCAN"),
				},
			},
			"limit": schema.Int64Attribute{
				Description: "Number of items to return",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
			},
			"ordering": schema.StringAttribute{
				Description: "Sort the response by given field.\nAvailable values: \"id\", \"type\", \"start_time\", \"end_time\", \"status\", \"message\", \"-id\", \"-type\", \"-start_time\", \"-end_time\", \"-status\", \"-message\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"id",
						"type",
						"start_time",
						"end_time",
						"status",
						"message",
						"-id",
						"-type",
						"-start_time",
						"-end_time",
						"-status",
						"-message",
					),
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
				CustomType:  customfield.NewNestedObjectListType[WaapDomainAPIDiscoveryScanResultsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The scan ID",
							Computed:    true,
						},
						"end_time": schema.StringAttribute{
							Description: "The date and time the scan ended",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"message": schema.StringAttribute{
							Description: "The message associated with the scan",
							Computed:    true,
						},
						"start_time": schema.StringAttribute{
							Description: "The date and time the scan started",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"status": schema.StringAttribute{
							Description: "The different statuses a task result can have\nAvailable values: \"SUCCESS\", \"FAILURE\", \"IN_PROGRESS\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"SUCCESS",
									"FAILURE",
									"IN_PROGRESS",
								),
							},
						},
						"type": schema.StringAttribute{
							Description: "The different types of scans that can be performed\nAvailable values: \"TRAFFIC_SCAN\", \"API_DESCRIPTION_FILE_SCAN\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("TRAFFIC_SCAN", "API_DESCRIPTION_FILE_SCAN"),
							},
						},
					},
				},
			},
		},
	}
}

func (d *WaapDomainAPIDiscoveryScanResultsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *WaapDomainAPIDiscoveryScanResultsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

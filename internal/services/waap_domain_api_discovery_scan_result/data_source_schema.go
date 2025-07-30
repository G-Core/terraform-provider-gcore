// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_discovery_scan_result

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainAPIDiscoveryScanResultDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"scan_id": schema.StringAttribute{
				Description: "The scan ID",
				Required:    true,
			},
			"end_time": schema.StringAttribute{
				Description: "The date and time the scan ended",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"id": schema.StringAttribute{
				Description: "The scan ID",
				Computed:    true,
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
	}
}

func (d *WaapDomainAPIDiscoveryScanResultDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapDomainAPIDiscoveryScanResultDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

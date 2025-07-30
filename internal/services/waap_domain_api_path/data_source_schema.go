// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_path

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainAPIPathDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"path_id": schema.StringAttribute{
				Description: "The path ID",
				Required:    true,
			},
			"api_version": schema.StringAttribute{
				Description: "The API version",
				Computed:    true,
			},
			"first_detected": schema.StringAttribute{
				Description: "The date and time in ISO 8601 format the API path was first detected.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"http_scheme": schema.StringAttribute{
				Description: "The different HTTP schemes an API path can have\nAvailable values: \"HTTP\", \"HTTPS\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("HTTP", "HTTPS"),
				},
			},
			"id": schema.StringAttribute{
				Description: "The path ID",
				Computed:    true,
			},
			"last_detected": schema.StringAttribute{
				Description: "The date and time in ISO 8601 format the API path was last detected.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"method": schema.StringAttribute{
				Description: "The different methods an API path can have\nAvailable values: \"GET\", \"POST\", \"PUT\", \"PATCH\", \"DELETE\", \"TRACE\", \"HEAD\", \"OPTIONS\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"GET",
						"POST",
						"PUT",
						"PATCH",
						"DELETE",
						"TRACE",
						"HEAD",
						"OPTIONS",
					),
				},
			},
			"path": schema.StringAttribute{
				Description: "The API path, locations that are saved for resource IDs will be put in curly brackets",
				Computed:    true,
			},
			"request_count": schema.Int64Attribute{
				Description: "The number of requests for this path in the last 24 hours",
				Computed:    true,
			},
			"source": schema.StringAttribute{
				Description: "The different sources an API path can have\nAvailable values: \"API_DESCRIPTION_FILE\", \"TRAFFIC_SCAN\", \"USER_DEFINED\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"API_DESCRIPTION_FILE",
						"TRAFFIC_SCAN",
						"USER_DEFINED",
					),
				},
			},
			"status": schema.StringAttribute{
				Description: "The different statuses an API path can have\nAvailable values: \"CONFIRMED_API\", \"POTENTIAL_API\", \"NOT_API\", \"DELISTED_API\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"CONFIRMED_API",
						"POTENTIAL_API",
						"NOT_API",
						"DELISTED_API",
					),
				},
			},
			"api_groups": schema.ListAttribute{
				Description: "An array of api groups associated with the API path",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"tags": schema.ListAttribute{
				Description: "An array of tags associated with the API path",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
		},
	}
}

func (d *WaapDomainAPIPathDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapDomainAPIPathDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

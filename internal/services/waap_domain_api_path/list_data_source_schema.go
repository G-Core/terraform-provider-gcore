// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_path

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

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainAPIPathsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"api_group": schema.StringAttribute{
				Description: "Filter by the API group associated with the API path",
				Optional:    true,
			},
			"api_version": schema.StringAttribute{
				Description: "Filter by the API version",
				Optional:    true,
			},
			"http_scheme": schema.StringAttribute{
				Description: "The different HTTP schemes an API path can have\nAvailable values: \"HTTP\", \"HTTPS\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("HTTP", "HTTPS"),
				},
			},
			"method": schema.StringAttribute{
				Description: "The different methods an API path can have\nAvailable values: \"GET\", \"POST\", \"PUT\", \"PATCH\", \"DELETE\", \"TRACE\", \"HEAD\", \"OPTIONS\".",
				Optional:    true,
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
			"ordering": schema.StringAttribute{
				Description: "Sort the response by given field.\nAvailable values: \"id\", \"path\", \"method\", \"api_version\", \"http_scheme\", \"first_detected\", \"last_detected\", \"status\", \"source\", \"-id\", \"-path\", \"-method\", \"-api_version\", \"-http_scheme\", \"-first_detected\", \"-last_detected\", \"-status\", \"-source\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"id",
						"path",
						"method",
						"api_version",
						"http_scheme",
						"first_detected",
						"last_detected",
						"status",
						"source",
						"-id",
						"-path",
						"-method",
						"-api_version",
						"-http_scheme",
						"-first_detected",
						"-last_detected",
						"-status",
						"-source",
					),
				},
			},
			"path": schema.StringAttribute{
				Description: "Filter by the path. Supports '\\*' as a wildcard character",
				Optional:    true,
			},
			"source": schema.StringAttribute{
				Description: "The different sources an API path can have\nAvailable values: \"API_DESCRIPTION_FILE\", \"TRAFFIC_SCAN\", \"USER_DEFINED\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"API_DESCRIPTION_FILE",
						"TRAFFIC_SCAN",
						"USER_DEFINED",
					),
				},
			},
			"ids": schema.ListAttribute{
				Description: "Filter by the path ID",
				Optional:    true,
				ElementType: types.StringType,
			},
			"status": schema.ListAttribute{
				Description: "Filter by the status of the discovered API path",
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOfCaseInsensitive(
							"CONFIRMED_API",
							"POTENTIAL_API",
							"NOT_API",
							"DELISTED_API",
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
					int64validator.Between(0, 10),
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
				CustomType:  customfield.NewNestedObjectListType[WaapDomainAPIPathsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The path ID",
							Computed:    true,
						},
						"api_groups": schema.ListAttribute{
							Description: "An array of api groups associated with the API path",
							Computed:    true,
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
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
						"tags": schema.ListAttribute{
							Description: "An array of tags associated with the API path",
							Computed:    true,
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *WaapDomainAPIPathsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *WaapDomainAPIPathsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

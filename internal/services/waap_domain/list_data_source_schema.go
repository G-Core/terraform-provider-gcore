// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filter domains based on the domain name. Supports '\\*' as a wildcard character",
				Optional:    true,
			},
			"ordering": schema.StringAttribute{
				Description: "Sort the response by given field.\nAvailable values: \"id\", \"name\", \"status\", \"created_at\", \"-id\", \"-name\", \"-status\", \"-created_at\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"id",
						"name",
						"status",
						"created_at",
						"-id",
						"-name",
						"-status",
						"-created_at",
					),
				},
			},
			"status": schema.StringAttribute{
				Description: "The different statuses a domain can have\nAvailable values: \"active\", \"bypass\", \"monitor\", \"locked\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"active",
						"bypass",
						"monitor",
						"locked",
					),
				},
			},
			"ids": schema.ListAttribute{
				Description: "Filter domains based on their IDs",
				Optional:    true,
				ElementType: types.Int64Type,
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
				CustomType:  customfield.NewNestedObjectListType[WaapDomainsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The domain ID",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "The date and time the domain was created in ISO 8601 format",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"custom_page_set": schema.Int64Attribute{
							Description: "The ID of the custom page set",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The domain name",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The different statuses a domain can have\nAvailable values: \"active\", \"bypass\", \"monitor\", \"locked\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"active",
									"bypass",
									"monitor",
									"locked",
								),
							},
						},
					},
				},
			},
		},
	}
}

func (d *WaapDomainsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *WaapDomainsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app

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

var _ datasource.DataSourceWithConfigValidators = (*FastedgeAppsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_type": schema.StringAttribute{
				Description: "API type:  \nwasi-http - WASI with HTTP entry point  \nproxy-wasm - Proxy-Wasm app, callable from CDN\nAvailable values: \"wasi-http\", \"proxy-wasm\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("wasi-http", "proxy-wasm"),
				},
			},
			"binary": schema.Int64Attribute{
				Description: "Binary ID",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the app",
				Optional:    true,
			},
			"ordering": schema.StringAttribute{
				Description: "Ordering\nAvailable values: \"name\", \"-name\", \"status\", \"-status\", \"id\", \"-id\", \"template\", \"-template\", \"binary\", \"-binary\", \"plan\", \"-plan\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"name",
						"-name",
						"status",
						"-status",
						"id",
						"-id",
						"template",
						"-template",
						"binary",
						"-binary",
						"plan",
						"-plan",
					),
				},
			},
			"plan": schema.Int64Attribute{
				Description: "Plan ID",
				Optional:    true,
			},
			"status": schema.Int64Attribute{
				Description: "Status code:  \n0 - draft (inactive)  \n1 - enabled  \n2 - disabled  \n3 - hourly call limit exceeded  \n4 - daily call limit exceeded  \n5 - suspended",
				Optional:    true,
			},
			"template": schema.Int64Attribute{
				Description: "Template ID",
				Optional:    true,
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
				CustomType:  customfield.NewNestedObjectListType[FastedgeAppsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "App ID",
							Computed:    true,
						},
						"api_type": schema.StringAttribute{
							Description: "Wasm API type",
							Computed:    true,
						},
						"binary": schema.Int64Attribute{
							Description: "Binary ID",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "App name",
							Computed:    true,
						},
						"plan_id": schema.Int64Attribute{
							Description: "Application plan ID",
							Computed:    true,
						},
						"status": schema.Int64Attribute{
							Description: "Status code:  \n0 - draft (inactive)  \n1 - enabled  \n2 - disabled  \n3 - hourly call limit exceeded  \n4 - daily call limit exceeded  \n5 - suspended",
							Computed:    true,
						},
						"comment": schema.StringAttribute{
							Description: "Description of the binary",
							Computed:    true,
						},
						"debug": schema.BoolAttribute{
							Description: "Switch on logging for 30 minutes (switched off by default)",
							Computed:    true,
						},
						"debug_until": schema.StringAttribute{
							Description: "When debugging finishes",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"networks": schema.ListAttribute{
							Description: "Networks",
							Computed:    true,
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
						},
						"plan": schema.StringAttribute{
							Description: "Application plan name",
							Computed:    true,
						},
						"template": schema.Int64Attribute{
							Description: "Template ID",
							Computed:    true,
						},
						"template_name": schema.StringAttribute{
							Description: "Template name",
							Computed:    true,
						},
						"upgradeable_to": schema.Int64Attribute{
							Description: "ID of the binary the app can be upgraded to",
							Computed:    true,
						},
						"url": schema.StringAttribute{
							Description: "App URL",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *FastedgeAppsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *FastedgeAppsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

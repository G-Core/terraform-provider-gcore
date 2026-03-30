// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_template

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*FastedgeTemplatesDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "FastEdge templates encapsulate reusable configurations for FastEdge applications, including a WebAssembly binary reference and configurable parameters.",
		Attributes: map[string]schema.Attribute{
			"api_type": schema.StringAttribute{
				Description: "API type:  \nwasi-http - WASI with HTTP entry point  \nproxy-wasm - Proxy-Wasm app, callable from CDN\nAvailable values: \"wasi-http\", \"proxy-wasm\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("wasi-http", "proxy-wasm"),
				},
			},
			"only_mine": schema.BoolAttribute{
				Description: "When true, returns only templates created by the client. When false, includes shared templates.",
				Computed:    true,
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
				CustomType:  customfield.NewNestedObjectListType[FastedgeTemplatesItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Template ID",
							Computed:    true,
						},
						"api_type": schema.StringAttribute{
							Description: "Wasm API type",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the template",
							Computed:    true,
						},
						"owned": schema.BoolAttribute{
							Description: "Is the template owned by user?",
							Computed:    true,
						},
						"long_descr": schema.StringAttribute{
							Description: "Long description of the template",
							Computed:    true,
						},
						"short_descr": schema.StringAttribute{
							Description: "Short description of the template",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *FastedgeTemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *FastedgeTemplatesDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

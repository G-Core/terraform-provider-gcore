// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_template

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*FastedgeTemplateDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "FastEdge templates encapsulate reusable configurations for FastEdge applications, including a WebAssembly binary reference and configurable parameters.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"template_id": schema.Int64Attribute{
				Optional: true,
			},
			"api_type": schema.StringAttribute{
				Description: "Wasm API type",
				Computed:    true,
			},
			"binary_id": schema.Int64Attribute{
				Description: "ID of the WebAssembly binary to use for this template",
				Computed:    true,
			},
			"long_descr": schema.StringAttribute{
				Description: "Detailed markdown description explaining template features and usage",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Unique name for the template (used for identification and searching)",
				Computed:    true,
			},
			"owned": schema.BoolAttribute{
				Description: "Is the template owned by user?",
				Computed:    true,
			},
			"short_descr": schema.StringAttribute{
				Description: "Brief one-line description displayed in template listings",
				Computed:    true,
			},
			"params": schema.ListNestedAttribute{
				Description: "Parameters",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[FastedgeTemplateParamsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"data_type": schema.StringAttribute{
							Description: "Parameter type determines validation and UI rendering:  \nstring - text input  \nnumber - numeric input  \ndate - date picker  \ntime - time picker  \nsecret - references a secret  \nstore - references an edge store  \nbool - boolean toggle  \njson - JSON editor or multiline text area with JSON validation  \nenum - dropdown/select with allowed values defined via parameter metadata\nAvailable values: \"string\", \"number\", \"date\", \"time\", \"secret\", \"store\", \"bool\", \"json\", \"enum\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"string",
									"number",
									"date",
									"time",
									"secret",
									"store",
									"bool",
									"json",
									"enum",
								),
							},
						},
						"mandatory": schema.BoolAttribute{
							Description: "If true, this parameter must be provided when instantiating the template",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Parameter name used as a placeholder in template (e.g., `API_KEY`, `DATABASE_URL`)",
							Computed:    true,
						},
						"descr": schema.StringAttribute{
							Description: "Human-readable explanation of what this parameter controls",
							Computed:    true,
						},
						"metadata": schema.StringAttribute{
							Description: "Optional JSON-encoded string for additional parameter metadata, such as allowed values for enum types",
							Computed:    true,
						},
					},
				},
			},
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"api_type": schema.StringAttribute{
						Description: "API type:  \nwasi-http - WASI with HTTP entry point  \nproxy-wasm - Proxy-Wasm app, callable from CDN\nAvailable values: \"wasi-http\", \"proxy-wasm\".",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("wasi-http", "proxy-wasm"),
						},
					},
					"limit": schema.Int64Attribute{
						Description: "Maximum number of results to return",
						Computed:    true,
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 1000),
						},
					},
					"only_mine": schema.BoolAttribute{
						Description: "When true, returns only templates created by the client. When false, includes shared templates.",
						Computed:    true,
						Optional:    true,
					},
				},
			},
		},
	}
}

func (d *FastedgeTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *FastedgeTemplateDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("template_id"), path.MatchRoot("find_one_by")),
	}
}

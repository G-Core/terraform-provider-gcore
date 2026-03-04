// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*FastedgeAppDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Apps are descriptions of edge apps, that reference the binary and may contain app-specific settings, such as environment variables.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Optional: true,
			},
			"api_type": schema.StringAttribute{
				Description: "Wasm API type",
				Computed:    true,
			},
			"binary": schema.Int64Attribute{
				Description: "Binary ID",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "App description",
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
			"log": schema.StringAttribute{
				Description: "Logging channel (by default - kafka, which allows exploring logs with API)\nAvailable values: \"kafka\", \"none\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("kafka", "none"),
				},
			},
			"name": schema.StringAttribute{
				Description: "App name",
				Computed:    true,
			},
			"plan": schema.StringAttribute{
				Description: "Plan name",
				Computed:    true,
			},
			"plan_id": schema.Int64Attribute{
				Description: "Plan ID",
				Computed:    true,
			},
			"status": schema.Int64Attribute{
				Description: "Status code:  \n0 - draft (inactive)  \n1 - enabled  \n2 - disabled  \n3 - hourly call limit exceeded  \n4 - daily call limit exceeded  \n5 - suspended",
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
			"url": schema.StringAttribute{
				Description: "App URL",
				Computed:    true,
			},
			"env": schema.MapAttribute{
				Description: "Environment variables",
				Computed:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"networks": schema.ListAttribute{
				Description: "Networks",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"rsp_headers": schema.MapAttribute{
				Description: "Extra headers to add to the response",
				Computed:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"secrets": schema.MapNestedAttribute{
				Description: "Application secrets",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectMapType[FastedgeAppSecretsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The unique identifier of the secret.",
							Computed:    true,
						},
						"comment": schema.StringAttribute{
							Description: "A description or comment about the secret.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The unique name of the secret.",
							Computed:    true,
						},
					},
				},
			},
			"stores": schema.MapNestedAttribute{
				Description: "Application edge stores",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectMapType[FastedgeAppStoresDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The identifier of the store",
							Computed:    true,
						},
						"comment": schema.StringAttribute{
							Description: "A description of the store",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the store",
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
				},
			},
		},
	}
}

func (d *FastedgeAppDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *FastedgeAppDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("find_one_by")),
	}
}

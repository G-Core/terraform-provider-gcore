// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_template

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*FastedgeTemplateDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"api_type": schema.StringAttribute{
				Description: "Wasm API type",
				Computed:    true,
			},
			"binary_id": schema.Int64Attribute{
				Description: "Binary ID",
				Computed:    true,
			},
			"long_descr": schema.StringAttribute{
				Description: "Long description of the template",
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
			"short_descr": schema.StringAttribute{
				Description: "Short description of the template",
				Computed:    true,
			},
			"params": schema.ListNestedAttribute{
				Description: "Parameters",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[FastedgeTemplateParamsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"data_type": schema.StringAttribute{
							Description: "Parameter type\nAvailable values: \"string\", \"number\", \"date\", \"time\", \"secret\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"string",
									"number",
									"date",
									"time",
									"secret",
								),
							},
						},
						"mandatory": schema.BoolAttribute{
							Description: "Is this field mandatory?",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Parameter name",
							Computed:    true,
						},
						"descr": schema.StringAttribute{
							Description: "Parameter description",
							Computed:    true,
						},
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
	return []datasource.ConfigValidator{}
}

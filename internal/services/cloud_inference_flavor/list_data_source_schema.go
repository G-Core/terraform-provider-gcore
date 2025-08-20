// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_flavor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceFlavorsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"limit": schema.Int64Attribute{
				Description: "Optional. Limit the number of returned items",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtMost(1000),
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
				CustomType:  customfield.NewNestedObjectListType[CloudInferenceFlavorsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cpu": schema.Float64Attribute{
							Description: "Inference flavor cpu count.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Inference flavor description.",
							Computed:    true,
						},
						"gpu": schema.Int64Attribute{
							Description: "Inference flavor gpu count.",
							Computed:    true,
						},
						"gpu_compute_capability": schema.StringAttribute{
							Description: "Inference flavor gpu compute capability.",
							Computed:    true,
						},
						"gpu_memory": schema.Float64Attribute{
							Description: "Inference flavor gpu memory in Gi.",
							Computed:    true,
						},
						"gpu_model": schema.StringAttribute{
							Description: "Inference flavor gpu model.",
							Computed:    true,
						},
						"is_gpu_shared": schema.BoolAttribute{
							Description: "Inference flavor is gpu shared.",
							Computed:    true,
						},
						"memory": schema.Float64Attribute{
							Description: "Inference flavor memory in Gi.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Inference flavor name.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudInferenceFlavorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudInferenceFlavorsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

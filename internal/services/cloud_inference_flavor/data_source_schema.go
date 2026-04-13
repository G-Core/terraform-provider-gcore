// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_flavor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceFlavorDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Inference flavors define the GPU and CPU resource configurations available for inference deployments.",
		Attributes: map[string]schema.Attribute{
			"flavor_name": schema.StringAttribute{
				Description: "Inference flavor name.",
				Required:    true,
			},
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
	}
}

func (d *CloudInferenceFlavorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudInferenceFlavorDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

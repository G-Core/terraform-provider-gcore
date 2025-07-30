// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_flavor

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudInferenceFlavorsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudInferenceFlavorsItemsDataSourceModel] `json:"results,computed"`
}

type CloudInferenceFlavorsDataSourceModel struct {
	Limit    types.Int64                                                             `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems types.Int64                                                             `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[CloudInferenceFlavorsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudInferenceFlavorsDataSourceModel) toListParams(_ context.Context) (params cloud.InferenceFlavorListParams, diags diag.Diagnostics) {
	params = cloud.InferenceFlavorListParams{}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}

	return
}

type CloudInferenceFlavorsItemsDataSourceModel struct {
	CPU                  types.Float64 `tfsdk:"cpu" json:"cpu,computed"`
	Description          types.String  `tfsdk:"description" json:"description,computed"`
	GPU                  types.Int64   `tfsdk:"gpu" json:"gpu,computed"`
	GPUComputeCapability types.String  `tfsdk:"gpu_compute_capability" json:"gpu_compute_capability,computed"`
	GPUMemory            types.Float64 `tfsdk:"gpu_memory" json:"gpu_memory,computed"`
	GPUModel             types.String  `tfsdk:"gpu_model" json:"gpu_model,computed"`
	IsGPUShared          types.Bool    `tfsdk:"is_gpu_shared" json:"is_gpu_shared,computed"`
	Memory               types.Float64 `tfsdk:"memory" json:"memory,computed"`
	Name                 types.String  `tfsdk:"name" json:"name,computed"`
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_flavor

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudInferenceFlavorDataSourceModel struct {
	FlavorName           types.String  `tfsdk:"flavor_name" path:"flavor_name,required"`
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

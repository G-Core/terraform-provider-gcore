// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_api_key

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudInferenceAPIKeyDataSourceModel struct {
	APIKeyName      types.String                   `tfsdk:"api_key_name" path:"api_key_name,required"`
	ProjectID       types.Int64                    `tfsdk:"project_id" path:"project_id,required"`
	CreatedAt       types.String                   `tfsdk:"created_at" json:"created_at,computed"`
	Description     types.String                   `tfsdk:"description" json:"description,computed"`
	ExpiresAt       types.String                   `tfsdk:"expires_at" json:"expires_at,computed"`
	Name            types.String                   `tfsdk:"name" json:"name,computed"`
	DeploymentNames customfield.List[types.String] `tfsdk:"deployment_names" json:"deployment_names,computed"`
}

func (m *CloudInferenceAPIKeyDataSourceModel) toReadParams(_ context.Context) (params cloud.InferenceAPIKeyGetParams, diags diag.Diagnostics) {
	params = cloud.InferenceAPIKeyGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

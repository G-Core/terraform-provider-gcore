// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_registry_credential

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudInferenceRegistryCredentialDataSourceModel struct {
	CredentialName types.String `tfsdk:"credential_name" path:"credential_name,required"`
	ProjectID      types.Int64  `tfsdk:"project_id" path:"project_id,required"`
	Name           types.String `tfsdk:"name" json:"name,computed"`
	RegistryURL    types.String `tfsdk:"registry_url" json:"registry_url,computed"`
	Username       types.String `tfsdk:"username" json:"username,computed"`
}

func (m *CloudInferenceRegistryCredentialDataSourceModel) toReadParams(_ context.Context) (params cloud.InferenceRegistryCredentialGetParams, diags diag.Diagnostics) {
	params = cloud.InferenceRegistryCredentialGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

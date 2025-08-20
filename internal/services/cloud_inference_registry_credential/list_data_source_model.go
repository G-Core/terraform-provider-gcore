// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_registry_credential

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudInferenceRegistryCredentialsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudInferenceRegistryCredentialsItemsDataSourceModel] `json:"results,computed"`
}

type CloudInferenceRegistryCredentialsDataSourceModel struct {
	ProjectID types.Int64                                                                         `tfsdk:"project_id" path:"project_id,optional"`
	Limit     types.Int64                                                                         `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems  types.Int64                                                                         `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudInferenceRegistryCredentialsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudInferenceRegistryCredentialsDataSourceModel) toListParams(_ context.Context) (params cloud.InferenceRegistryCredentialListParams, diags diag.Diagnostics) {
	params = cloud.InferenceRegistryCredentialListParams{}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

type CloudInferenceRegistryCredentialsItemsDataSourceModel struct {
	Name        types.String `tfsdk:"name" json:"name,computed"`
	ProjectID   types.Int64  `tfsdk:"project_id" json:"project_id,computed"`
	RegistryURL types.String `tfsdk:"registry_url" json:"registry_url,computed"`
	Username    types.String `tfsdk:"username" json:"username,computed"`
}

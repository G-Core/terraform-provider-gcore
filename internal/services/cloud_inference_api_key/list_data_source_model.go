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

type CloudInferenceAPIKeysResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudInferenceAPIKeysItemsDataSourceModel] `json:"results,computed"`
}

type CloudInferenceAPIKeysDataSourceModel struct {
	ProjectID types.Int64                                                             `tfsdk:"project_id" path:"project_id,optional"`
	Limit     types.Int64                                                             `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems  types.Int64                                                             `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudInferenceAPIKeysItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudInferenceAPIKeysDataSourceModel) toListParams(_ context.Context) (params cloud.InferenceAPIKeyListParams, diags diag.Diagnostics) {
	params = cloud.InferenceAPIKeyListParams{}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

type CloudInferenceAPIKeysItemsDataSourceModel struct {
	CreatedAt       types.String                   `tfsdk:"created_at" json:"created_at,computed"`
	DeploymentNames customfield.List[types.String] `tfsdk:"deployment_names" json:"deployment_names,computed"`
	Description     types.String                   `tfsdk:"description" json:"description,computed"`
	ExpiresAt       types.String                   `tfsdk:"expires_at" json:"expires_at,computed"`
	Name            types.String                   `tfsdk:"name" json:"name,computed"`
}

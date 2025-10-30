// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_secret

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudInferenceSecretsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudInferenceSecretsItemsDataSourceModel] `json:"results,computed"`
}

type CloudInferenceSecretsDataSourceModel struct {
	ProjectID types.Int64                                                             `tfsdk:"project_id" path:"project_id,optional"`
	MaxItems  types.Int64                                                             `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudInferenceSecretsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudInferenceSecretsDataSourceModel) toListParams(_ context.Context) (params cloud.InferenceSecretListParams, diags diag.Diagnostics) {
	params = cloud.InferenceSecretListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

type CloudInferenceSecretsItemsDataSourceModel struct {
	ID   types.String                                                       `tfsdk:"id" json:"name,computed"`
	Data customfield.NestedObject[CloudInferenceSecretsDataDataSourceModel] `tfsdk:"data" json:"data,computed"`
	Name types.String                                                       `tfsdk:"name" json:"name,computed"`
	Type types.String                                                       `tfsdk:"type" json:"type,computed"`
}

type CloudInferenceSecretsDataDataSourceModel struct {
	AwsAccessKeyID     types.String `tfsdk:"aws_access_key_id" json:"aws_access_key_id,computed"`
	AwsSecretAccessKey types.String `tfsdk:"aws_secret_access_key" json:"aws_secret_access_key,computed"`
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_secret

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudInferenceSecretDataSourceModel struct {
	ID         types.String                                                      `tfsdk:"id" path:"secret_name,computed"`
	SecretName types.String                                                      `tfsdk:"secret_name" path:"secret_name,optional"`
	ProjectID  types.Int64                                                       `tfsdk:"project_id" path:"project_id,optional"`
	Name       types.String                                                      `tfsdk:"name" json:"name,computed"`
	Type       types.String                                                      `tfsdk:"type" json:"type,computed"`
	Data       customfield.NestedObject[CloudInferenceSecretDataDataSourceModel] `tfsdk:"data" json:"data,computed"`
	FindOneBy  *CloudInferenceSecretFindOneByDataSourceModel                     `tfsdk:"find_one_by"`
}

func (m *CloudInferenceSecretDataSourceModel) toReadParams(_ context.Context) (params cloud.InferenceSecretGetParams, diags diag.Diagnostics) {
	params = cloud.InferenceSecretGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

func (m *CloudInferenceSecretDataSourceModel) toListParams(_ context.Context) (params cloud.InferenceSecretListParams, diags diag.Diagnostics) {
	params = cloud.InferenceSecretListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.FindOneBy.Limit.IsNull() {
		params.Limit = param.NewOpt(m.FindOneBy.Limit.ValueInt64())
	}

	return
}

type CloudInferenceSecretDataDataSourceModel struct {
	AwsAccessKeyID     types.String `tfsdk:"aws_access_key_id" json:"aws_access_key_id,computed"`
	AwsSecretAccessKey types.String `tfsdk:"aws_secret_access_key" json:"aws_secret_access_key,computed"`
}

type CloudInferenceSecretFindOneByDataSourceModel struct {
	Limit types.Int64 `tfsdk:"limit" query:"limit,computed_optional"`
}

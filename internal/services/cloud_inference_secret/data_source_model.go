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

type CloudInferenceSecretDataSourceModel struct {
	ProjectID  types.Int64                                                       `tfsdk:"project_id" path:"project_id,required"`
	SecretName types.String                                                      `tfsdk:"secret_name" path:"secret_name,required"`
	Name       types.String                                                      `tfsdk:"name" json:"name,computed"`
	Type       types.String                                                      `tfsdk:"type" json:"type,computed"`
	Data       customfield.NestedObject[CloudInferenceSecretDataDataSourceModel] `tfsdk:"data" json:"data,computed"`
}

func (m *CloudInferenceSecretDataSourceModel) toReadParams(_ context.Context) (params cloud.InferenceSecretGetParams, diags diag.Diagnostics) {
	params = cloud.InferenceSecretGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

type CloudInferenceSecretDataDataSourceModel struct {
	AwsAccessKeyID     types.String `tfsdk:"aws_access_key_id" json:"aws_access_key_id,computed"`
	AwsSecretAccessKey types.String `tfsdk:"aws_secret_access_key" json:"aws_secret_access_key,computed"`
}

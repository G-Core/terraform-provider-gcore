// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_secret

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudInferenceSecretModel struct {
	ID        types.String                   `tfsdk:"id" json:"-,computed"`
	Name      types.String                   `tfsdk:"name" json:"name,required"`
	ProjectID types.Int64                    `tfsdk:"project_id" path:"project_id,optional"`
	Type      types.String                   `tfsdk:"type" json:"type,required"`
	Data      *CloudInferenceSecretDataModel `tfsdk:"data" json:"data,required"`
}

func (m CloudInferenceSecretModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudInferenceSecretModel) MarshalJSONForUpdate(state CloudInferenceSecretModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type CloudInferenceSecretDataModel struct {
	AwsAccessKeyID     types.String `tfsdk:"aws_access_key_id" json:"aws_access_key_id,required"`
	AwsSecretAccessKey types.String `tfsdk:"aws_secret_access_key" json:"aws_secret_access_key,required"`
}

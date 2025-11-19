// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_registry_credential

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type CloudInferenceRegistryCredentialModel struct {
	ID          types.String `tfsdk:"id" json:"-,computed"`
	Name        types.String `tfsdk:"name" json:"name,required"`
	ProjectID   types.Int64  `tfsdk:"project_id" path:"project_id,optional"`
	Password    types.String `tfsdk:"password" json:"password,required,no_refresh"`
	RegistryURL types.String `tfsdk:"registry_url" json:"registry_url,required"`
	Username    types.String `tfsdk:"username" json:"username,required"`
}

func (m CloudInferenceRegistryCredentialModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudInferenceRegistryCredentialModel) MarshalJSONForUpdate(state CloudInferenceRegistryCredentialModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

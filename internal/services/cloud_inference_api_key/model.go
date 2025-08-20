// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_api_key

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudInferenceAPIKeyModel struct {
	APIKeyName      types.String                   `tfsdk:"api_key_name" path:"api_key_name,optional"`
	ProjectID       types.Int64                    `tfsdk:"project_id" path:"project_id,optional"`
	Name            types.String                   `tfsdk:"name" json:"name,required"`
	ExpiresAt       types.String                   `tfsdk:"expires_at" json:"expires_at,optional"`
	Description     types.String                   `tfsdk:"description" json:"description,optional"`
	CreatedAt       types.String                   `tfsdk:"created_at" json:"created_at,computed"`
	Secret          types.String                   `tfsdk:"secret" json:"secret,computed,no_refresh"`
	DeploymentNames customfield.List[types.String] `tfsdk:"deployment_names" json:"deployment_names,computed"`
}

func (m CloudInferenceAPIKeyModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudInferenceAPIKeyModel) MarshalJSONForUpdate(state CloudInferenceAPIKeyModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

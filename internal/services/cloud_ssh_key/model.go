// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_ssh_key

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudSSHKeyModel struct {
	ID              types.String      `tfsdk:"id" json:"id,computed"`
	ProjectID       types.Int64       `tfsdk:"project_id" path:"project_id,optional"`
	Name            types.String      `tfsdk:"name" json:"name,required"`
	PublicKey       types.String      `tfsdk:"public_key" json:"public_key,required"`
	SharedInProject types.Bool        `tfsdk:"shared_in_project" json:"shared_in_project,computed_optional"`
	CreatedAt       timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Fingerprint     types.String      `tfsdk:"fingerprint" json:"fingerprint,computed"`
	PrivateKey      types.String      `tfsdk:"private_key" json:"private_key,computed,no_refresh"`
	State           types.String      `tfsdk:"state" json:"state,computed"`
}

func (m CloudSSHKeyModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudSSHKeyModel) MarshalJSONForUpdate(state CloudSSHKeyModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

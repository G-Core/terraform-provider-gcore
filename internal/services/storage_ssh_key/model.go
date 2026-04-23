// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_ssh_key

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageSSHKeyModel struct {
	ID        types.Int64       `tfsdk:"id" json:"id,computed"`
	Name      types.String      `tfsdk:"name" json:"name,required"`
	PublicKey types.String      `tfsdk:"public_key" json:"public_key,required"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
}

func (m StorageSSHKeyModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StorageSSHKeyModel) MarshalJSONForUpdate(state StorageSSHKeyModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

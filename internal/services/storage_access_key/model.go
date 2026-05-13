// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_access_key

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageAccessKeyModel struct {
	ID        types.String      `tfsdk:"id" json:"-,computed"`
	AccessKey types.String      `tfsdk:"access_key" json:"access_key,computed"`
	StorageID types.Int64       `tfsdk:"storage_id" path:"storage_id,required"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	SecretKey types.String      `tfsdk:"secret_key" json:"secret_key,computed,no_refresh"`
}

func (m StorageAccessKeyModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StorageAccessKeyModel) MarshalJSONForUpdate(state StorageAccessKeyModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

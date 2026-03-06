// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_secret

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FastedgeSecretModel struct {
	ID          types.Int64                                                 `tfsdk:"id" json:"id,computed"`
	Name        types.String                                                `tfsdk:"name" json:"name,required"`
	Comment     types.String                                                `tfsdk:"comment" json:"comment,optional"`
	SecretSlots customfield.NestedObjectSet[FastedgeSecretSecretSlotsModel] `tfsdk:"secret_slots" json:"secret_slots,computed_optional"`
	AppCount    types.Int64                                                 `tfsdk:"app_count" json:"app_count,computed"`
}

func (m FastedgeSecretModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m FastedgeSecretModel) MarshalJSONForUpdate(state FastedgeSecretModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type FastedgeSecretSecretSlotsModel struct {
	Slot     types.Int64  `tfsdk:"slot" json:"slot,required"`
	Checksum types.String `tfsdk:"checksum" json:"checksum,computed"`
	Value    types.String `tfsdk:"value" json:"value,optional"`
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_secret

import (
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FastedgeSecretDataSourceModel struct {
	ID          types.Int64                                                           `tfsdk:"id" path:"secret_id,computed"`
	SecretID    types.Int64                                                           `tfsdk:"secret_id" path:"secret_id,required"`
	AppCount    types.Int64                                                           `tfsdk:"app_count" json:"app_count,computed"`
	Comment     types.String                                                          `tfsdk:"comment" json:"comment,computed"`
	Name        types.String                                                          `tfsdk:"name" json:"name,computed"`
	SecretSlots customfield.NestedObjectSet[FastedgeSecretSecretSlotsDataSourceModel] `tfsdk:"secret_slots" json:"secret_slots,computed"`
}

type FastedgeSecretSecretSlotsDataSourceModel struct {
	Slot     types.Int64  `tfsdk:"slot" json:"slot,computed"`
	Checksum types.String `tfsdk:"checksum" json:"checksum,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

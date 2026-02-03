// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_secret

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type FastedgeSecretDataSourceModel struct {
	ID          types.Int64                                                            `tfsdk:"id" path:"id,required"`
	AppCount    types.Int64                                                            `tfsdk:"app_count" json:"app_count,computed"`
	Comment     types.String                                                           `tfsdk:"comment" json:"comment,computed"`
	Name        types.String                                                           `tfsdk:"name" json:"name,computed"`
	SecretSlots customfield.NestedObjectList[FastedgeSecretSecretSlotsDataSourceModel] `tfsdk:"secret_slots" json:"secret_slots,computed"`
}

type FastedgeSecretSecretSlotsDataSourceModel struct {
	Slot     types.Int64  `tfsdk:"slot" json:"slot,computed"`
	Checksum types.String `tfsdk:"checksum" json:"checksum,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

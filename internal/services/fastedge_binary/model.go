// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_binary

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type FastedgeBinaryModel struct {
	ID         types.Int64  `tfsdk:"id" json:"id,computed"`
	Body       types.String `tfsdk:"body" json:"body,required,no_refresh"`
	APIType    types.String `tfsdk:"api_type" json:"api_type,computed"`
	Checksum   types.String `tfsdk:"checksum" json:"checksum,computed"`
	Source     types.Int64  `tfsdk:"source" json:"source,computed"`
	Status     types.Int64  `tfsdk:"status" json:"status,computed"`
	UnrefSince types.String `tfsdk:"unref_since" json:"unref_since,computed"`
}

func (m FastedgeBinaryModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m.Body)
}

func (m FastedgeBinaryModel) MarshalJSONForUpdate(state FastedgeBinaryModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m.Body, state.Body)
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_applied_preset

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNAppliedPresetModel struct {
	ID       types.Int64 `tfsdk:"id" json:"-,computed"`
	ObjectID types.Int64 `tfsdk:"object_id" json:"object_id,required"`
	PresetID types.Int64 `tfsdk:"preset_id" path:"preset_id,required"`
}

func (m CDNAppliedPresetModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CDNAppliedPresetModel) MarshalJSONForUpdate(state CDNAppliedPresetModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

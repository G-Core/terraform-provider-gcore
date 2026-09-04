// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_preset

import (
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNPresetDataSourceModel struct {
	PresetID       types.Int64                           `tfsdk:"preset_id" path:"preset_id,required"`
	ID             types.Int64                           `tfsdk:"id" json:"id,computed"`
	Name           types.String                          `tfsdk:"name" json:"name,computed"`
	ObjectType     types.String                          `tfsdk:"object_type" json:"object_type,computed"`
	Service        types.String                          `tfsdk:"service" json:"service,computed"`
	PresetSettings customfield.Map[jsontypes.Normalized] `tfsdk:"preset_settings" json:"preset_settings,computed"`
}

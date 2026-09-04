// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_preset

import (
	"context"

	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNPresetsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CDNPresetsItemsDataSourceModel] `json:"results,computed"`
}

type CDNPresetsDataSourceModel struct {
	MaxItems types.Int64                                                  `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[CDNPresetsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CDNPresetsDataSourceModel) toListParams(_ context.Context) (params cdn.PresetListParams, diags diag.Diagnostics) {
	params = cdn.PresetListParams{}

	return
}

type CDNPresetsItemsDataSourceModel struct {
	ID             types.Int64                           `tfsdk:"id" json:"id,computed"`
	Name           types.String                          `tfsdk:"name" json:"name,computed"`
	ObjectType     types.String                          `tfsdk:"object_type" json:"object_type,computed"`
	PresetSettings customfield.Map[jsontypes.Normalized] `tfsdk:"preset_settings" json:"preset_settings,computed"`
	Service        types.String                          `tfsdk:"service" json:"service,computed"`
}

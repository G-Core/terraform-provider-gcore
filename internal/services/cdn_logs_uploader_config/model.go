// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_config

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNLogsUploaderConfigModel struct {
	ID              types.Int64                   `tfsdk:"id" json:"id,computed"`
	Name            types.String                  `tfsdk:"name" json:"name,required"`
	Policy          types.Int64                   `tfsdk:"policy" json:"policy,required"`
	Target          types.Int64                   `tfsdk:"target" json:"target,required"`
	Resources       customfield.List[types.Int64] `tfsdk:"resources" json:"resources,computed_optional"`
	Enabled         types.Bool                    `tfsdk:"enabled" json:"enabled,computed_optional"`
	ForAllResources types.Bool                    `tfsdk:"for_all_resources" json:"for_all_resources,computed_optional"`
	ClientID        types.Int64                   `tfsdk:"client_id" json:"client_id,computed"`
	Created         timetypes.RFC3339             `tfsdk:"created" json:"created,computed" format:"date-time"`
	Updated         timetypes.RFC3339             `tfsdk:"updated" json:"updated,computed" format:"date-time"`
	Status          jsontypes.Normalized          `tfsdk:"status" json:"status,computed"`
}

func (m CDNLogsUploaderConfigModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CDNLogsUploaderConfigModel) MarshalJSONForUpdate(state CDNLogsUploaderConfigModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_config

import (
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNLogsUploaderConfigDataSourceModel struct {
	ID              types.Int64                   `tfsdk:"id" path:"id,required"`
	ClientID        types.Int64                   `tfsdk:"client_id" json:"client_id,computed"`
	Created         timetypes.RFC3339             `tfsdk:"created" json:"created,computed" format:"date-time"`
	Enabled         types.Bool                    `tfsdk:"enabled" json:"enabled,computed"`
	ForAllResources types.Bool                    `tfsdk:"for_all_resources" json:"for_all_resources,computed"`
	Name            types.String                  `tfsdk:"name" json:"name,computed"`
	Policy          types.Int64                   `tfsdk:"policy" json:"policy,computed"`
	Target          types.Int64                   `tfsdk:"target" json:"target,computed"`
	Updated         timetypes.RFC3339             `tfsdk:"updated" json:"updated,computed" format:"date-time"`
	Resources       customfield.List[types.Int64] `tfsdk:"resources" json:"resources,computed"`
	Status          jsontypes.Normalized          `tfsdk:"status" json:"status,computed"`
}

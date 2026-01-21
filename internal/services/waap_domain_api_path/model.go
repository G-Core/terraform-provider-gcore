// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_path

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type WaapDomainAPIPathModel struct {
	ID            types.String      `tfsdk:"id" json:"id,computed"`
	DomainID      types.Int64       `tfsdk:"domain_id" path:"domain_id,required"`
	HTTPScheme    types.String      `tfsdk:"http_scheme" json:"http_scheme,required"`
	Method        types.String      `tfsdk:"method" json:"method,required"`
	APIVersion    types.String      `tfsdk:"api_version" json:"api_version,optional"`
	Path          types.String      `tfsdk:"path" json:"path,required"`
	Status        types.String      `tfsdk:"status" json:"status,optional"`
	APIGroups     *[]types.String   `tfsdk:"api_groups" json:"api_groups,optional"`
	Tags          *[]types.String   `tfsdk:"tags" json:"tags,optional"`
	FirstDetected timetypes.RFC3339 `tfsdk:"first_detected" json:"first_detected,computed" format:"date-time"`
	LastDetected  timetypes.RFC3339 `tfsdk:"last_detected" json:"last_detected,computed" format:"date-time"`
	RequestCount  types.Int64       `tfsdk:"request_count" json:"request_count,computed"`
	Source        types.String      `tfsdk:"source" json:"source,computed"`
}

func (m WaapDomainAPIPathModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m WaapDomainAPIPathModel) MarshalJSONForUpdate(state WaapDomainAPIPathModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

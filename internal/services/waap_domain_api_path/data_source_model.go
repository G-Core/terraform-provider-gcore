// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_path

import (
	"context"

	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainAPIPathDataSourceModel struct {
	DomainID      types.Int64                    `tfsdk:"domain_id" path:"domain_id,required"`
	PathID        types.String                   `tfsdk:"path_id" path:"path_id,required"`
	APIVersion    types.String                   `tfsdk:"api_version" json:"api_version,computed"`
	FirstDetected timetypes.RFC3339              `tfsdk:"first_detected" json:"first_detected,computed" format:"date-time"`
	HTTPScheme    types.String                   `tfsdk:"http_scheme" json:"http_scheme,computed"`
	ID            types.String                   `tfsdk:"id" json:"id,computed"`
	LastDetected  timetypes.RFC3339              `tfsdk:"last_detected" json:"last_detected,computed" format:"date-time"`
	Method        types.String                   `tfsdk:"method" json:"method,computed"`
	Path          types.String                   `tfsdk:"path" json:"path,computed"`
	RequestCount  types.Int64                    `tfsdk:"request_count" json:"request_count,computed"`
	Source        types.String                   `tfsdk:"source" json:"source,computed"`
	Status        types.String                   `tfsdk:"status" json:"status,computed"`
	APIGroups     customfield.List[types.String] `tfsdk:"api_groups" json:"api_groups,computed"`
	Tags          customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
}

func (m *WaapDomainAPIPathDataSourceModel) toReadParams(_ context.Context) (params waap.DomainAPIPathGetParams, diags diag.Diagnostics) {
	params = waap.DomainAPIPathGetParams{
		DomainID: m.DomainID.ValueInt64(),
	}

	return
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_discovery_scan_result

import (
	"context"

	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WaapDomainAPIDiscoveryScanResultDataSourceModel struct {
	DomainID  types.Int64       `tfsdk:"domain_id" path:"domain_id,required"`
	ScanID    types.String      `tfsdk:"scan_id" path:"scan_id,required"`
	EndTime   timetypes.RFC3339 `tfsdk:"end_time" json:"end_time,computed" format:"date-time"`
	ID        types.String      `tfsdk:"id" json:"id,computed"`
	Message   types.String      `tfsdk:"message" json:"message,computed"`
	StartTime timetypes.RFC3339 `tfsdk:"start_time" json:"start_time,computed" format:"date-time"`
	Status    types.String      `tfsdk:"status" json:"status,computed"`
	Type      types.String      `tfsdk:"type" json:"type,computed"`
}

func (m *WaapDomainAPIDiscoveryScanResultDataSourceModel) toReadParams(_ context.Context) (params waap.DomainAPIDiscoveryScanResultGetParams, diags diag.Diagnostics) {
	params = waap.DomainAPIDiscoveryScanResultGetParams{
		DomainID: m.DomainID.ValueInt64(),
	}

	return
}

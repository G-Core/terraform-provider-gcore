// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_insight

import (
	"context"

	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainInsightDataSourceModel struct {
	DomainID         types.Int64                   `tfsdk:"domain_id" path:"domain_id,required"`
	InsightID        types.String                  `tfsdk:"insight_id" path:"insight_id,required"`
	Description      types.String                  `tfsdk:"description" json:"description,computed"`
	FirstSeen        timetypes.RFC3339             `tfsdk:"first_seen" json:"first_seen,computed" format:"date-time"`
	ID               types.String                  `tfsdk:"id" json:"id,computed"`
	InsightType      types.String                  `tfsdk:"insight_type" json:"insight_type,computed"`
	LastSeen         timetypes.RFC3339             `tfsdk:"last_seen" json:"last_seen,computed" format:"date-time"`
	LastStatusChange timetypes.RFC3339             `tfsdk:"last_status_change" json:"last_status_change,computed" format:"date-time"`
	Recommendation   types.String                  `tfsdk:"recommendation" json:"recommendation,computed"`
	Status           types.String                  `tfsdk:"status" json:"status,computed"`
	Labels           customfield.Map[types.String] `tfsdk:"labels" json:"labels,computed"`
}

func (m *WaapDomainInsightDataSourceModel) toReadParams(_ context.Context) (params waap.DomainInsightGetParams, diags diag.Diagnostics) {
	params = waap.DomainInsightGetParams{
		DomainID: m.DomainID.ValueInt64(),
	}

	return
}

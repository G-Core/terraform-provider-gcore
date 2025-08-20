// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_insight_silence

import (
	"context"

	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainInsightSilenceDataSourceModel struct {
	DomainID    types.Int64                   `tfsdk:"domain_id" path:"domain_id,required"`
	SilenceID   types.String                  `tfsdk:"silence_id" path:"silence_id,required"`
	Author      types.String                  `tfsdk:"author" json:"author,computed"`
	Comment     types.String                  `tfsdk:"comment" json:"comment,computed"`
	ExpireAt    timetypes.RFC3339             `tfsdk:"expire_at" json:"expire_at,computed" format:"date-time"`
	ID          types.String                  `tfsdk:"id" json:"id,computed"`
	InsightType types.String                  `tfsdk:"insight_type" json:"insight_type,computed"`
	Labels      customfield.Map[types.String] `tfsdk:"labels" json:"labels,computed"`
}

func (m *WaapDomainInsightSilenceDataSourceModel) toReadParams(_ context.Context) (params waap.DomainInsightSilenceGetParams, diags diag.Diagnostics) {
	params = waap.DomainInsightSilenceGetParams{
		DomainID: m.DomainID.ValueInt64(),
	}

	return
}
